package cores

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type WorkerFunc func(ctx context.Context) error

type Worker interface {
	Name() string
	Fn(ctx context.Context) error
}

type base struct {
	name string
	fn   WorkerFunc
}

type simpleWorker struct {
	*base
	daemon bool
}

func (w *simpleWorker) Name() string {
	return w.name
}

func (w *simpleWorker) Fn(ctx context.Context) error {
	return w.fn(ctx)
}

type cronWorker struct {
	*base
	cron    string
	running bool
}

func (w *cronWorker) Name() string {
	return w.name
}

func (w *cronWorker) Fn(ctx context.Context) error {
	return w.fn(ctx)
}

func withSimpleWorker(name string, daemon bool, fn WorkerFunc) ServiceOption {
	return func(cs *Options) {
		cs.workers = append(cs.workers, &simpleWorker{
			base: &base{
				name: name,
				fn:   fn,
			},
			daemon: daemon,
		})
	}
}

func WithWorker(fn WorkerFunc) ServiceOption {
	return withSimpleWorker(funcName(fn), false, fn)
}

func WithNameWorker(name string, fn WorkerFunc) ServiceOption {
	return withSimpleWorker(name, false, fn)
}

func WithDaemonWorker(fn WorkerFunc) ServiceOption {
	return withSimpleWorker(funcName(fn), true, fn)
}

func WithDaemonNameWorker(name string, fn WorkerFunc) ServiceOption {
	return withSimpleWorker(name, true, fn)
}

func WithCronWorker(cron string, fn WorkerFunc) ServiceOption {
	return func(cs *Options) {
		cs.workers = append(cs.workers, &cronWorker{
			base: &base{
				name: funcName(fn),
				fn:   fn,
			},
			cron: cron,
		})
	}
}

func funcName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func StructName(i any) string {
	tPtr := reflect.TypeOf(i)

	if tPtr.Kind() == reflect.Ptr {
		tPtr = tPtr.Elem()
	}

	return tPtr.Name()
}

func (c *CoreService) mountSimpleWorker(worker *simpleWorker) {
	fn := func(ctx context.Context) error {
		var c context.Context
		if worker.daemon {
			c = ctx
		} else {
			var cancel func()
			c, cancel = context.WithCancel(context.TODO())
			go func() {
				select {
				case <-ctx.Done():
					cancel()
				}
			}()
		}
		c = plog.With(c, "Worker", worker.name)

		if err := worker.fn(c); err != nil {
			plog.Errorc(c, "worker: %v run error: %v", worker.name, err)
			return errors.Wrap(err, worker.name)
		}
		return nil
	}

	mf := mountFn{
		fn:     fn,
		name:   worker.name,
		daemon: worker.daemon,
	}

	c.mountFns = append(c.mountFns, mf)
}

func (c *CoreService) mountCronWorker(worker *cronWorker) {
	if c.cron == nil {
		c.cron = cron.New()
		c.mountFns = append(c.mountFns, mountFn{
			fn: func(ctx context.Context) error {
				defer c.cron.Stop()
				err := waitContext(ctx, true, func(_ context.Context) error {
					c.cron.Run()
					return nil
				})

				if err != nil {
					plog.Errorc(ctx, "Run cron error: %v", err)
					return errors.Wrap(err, "Run cron error")
				}
				return nil
			},
			name:   "CronWorker",
			daemon: true,
		})
	}

	sched, err := cron.ParseStandard(worker.cron)
	if err != nil {
		plog.Fatalf("cron worker cron invalid. err: %v", err)
	}

	fn := func() {
		ctx := plog.With(c.ctx, "Worker", worker.name)
		defer plog.Debugc(ctx, "Cron Worker: %v Next scheduler time: %v", worker.name, sched.Next(time.Now()))

		err := waitContext(ctx, false, func(ctx context.Context) error {
			if worker.running {
				return errors.New("job still running")
			}
			worker.running = true
			defer func() { worker.running = false }()

			return worker.fn(ctx)
		})

		if err != nil {
			plog.Errorc(ctx, "Run cron worker: %v error: %v", worker.name, err)
		}
	}

	c.cron.AddFunc(worker.cron, fn)
}

func (c *CoreService) wrapWorker() {
	for _, worker := range c.opts.workers {
		switch w := worker.(type) {
		case *simpleWorker:
			if w.name == "" {
				w.name = plog.GetFuncName(w)
			}
			c.mountSimpleWorker(w)
		case *cronWorker:
			if w.name == "" {
				w.name = plog.GetFuncName(w)
			}
			c.mountCronWorker(w)
		default:
			plog.Warnc(c.ctx, "Unknown worker type. worker: %v", plog.GetFuncName(w))
		}
	}
}

func waitContext(ctx context.Context, quitImmediately bool, fn func(context.Context) error) error {
	stop := make(chan error)

	go func() {
		stop <- fn(ctx)
	}()

	select {
	case err := <-stop:
		return err
	case <-ctx.Done():
		if quitImmediately {
			return ctx.Err()
		}

		ticket := time.NewTicker(time.Second * 5)
		defer ticket.Stop()

		select {
		case err := <-stop:
			return err
		case <-ticket.C:
			plog.Infoc(ctx, "Force closing worker")
			return errors.Wrap(ctx.Err(), "Force close")
		}
	}
}
