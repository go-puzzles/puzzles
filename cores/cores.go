package cores

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/robfig/cron/v3"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

type listenEntry interface {
	~int | ~string
}

type mountFn struct {
	name   string
	daemon bool

	fn func(ctx context.Context) error

	cron    string
	running bool

	sched cron.Schedule
}

type CoreService struct {
	ctx      context.Context
	cancel   func()
	opts     *Options
	listener net.Listener

	cron     *cron.Cron
	mountFns []mountFn
}

type Options struct {
	Ctx          context.Context
	ListenerAddr string
	ServiceName  string
	Tags         []string
	Cmux         cmux.CMux
	WaitPprof    chan struct{}

	puzzles map[string]Puzzle
	workers []Worker
}

// RegisterPuzzle will put Puzzle into cores service
// It replaces the previously registered Puzzle with the same name
func (o *Options) RegisterPuzzle(p Puzzle) {
	o.puzzles[p.Name()] = p
}

type ServiceOption func(*Options)

func NewPuzzleCore(opts ...ServiceOption) *CoreService {
	ctx, cancel := context.WithCancel(context.TODO())

	c := &CoreService{
		ctx:    ctx,
		cancel: cancel,
		opts: &Options{
			Ctx:     ctx,
			puzzles: make(map[string]Puzzle),
			workers: make([]Worker, 0),
		},
		mountFns: make([]mountFn, 0),
	}

	for _, opt := range opts {
		opt(c.opts)
	}
	return c
}

func (c *CoreService) options() *Options {
	return c.opts
}

func (c *CoreService) isHttpRun() bool {
	return c.listener != nil
}

func (c *CoreService) runMountFn() error {
	grp, ctx := errgroup.WithContext(c.ctx)

	for _, mount := range c.mountFns {
		mf := mount
		cc := plog.With(ctx, "Worker", mf.name)
		grp.Go(func() (err error) {
			err = waitContext(cc, false, mf.fn)
			if err != nil {
				if mf.daemon {
					return err
				}
				return nil
			}
			return nil
		})
	}

	return grp.Wait()
}

func (c *CoreService) serve() error {
	c.mountFns = []mountFn{
		c.gracefulKill(),
	}

	if c.listener != nil {
		c.opts.Cmux = cmux.New(c.listener)
		c.mountFns = append(c.mountFns, c.listenCmux())
	}

	c.wrapWorker()
	c.wrapPuzzles()
	c.welcome()
	return c.runMountFn()
}

func (c *CoreService) listenCmux() mountFn {
	return mountFn{
		fn: func(ctx context.Context) (err error) {
			defer func() {
				if errors.Is(err, net.ErrClosed) {
					err = nil
				}

				if err != nil {
					plog.Errorc(ctx, "Cmux serve error: %v", err)
				}
			}()
			return c.opts.Cmux.Serve()
		},
		name:   "CmuxListener",
		daemon: true,
	}
}

func (c *CoreService) runWithListener(lis net.Listener) error {
	c.opts.ListenerAddr = lis.Addr().String()
	return c.serve()
}

func (c *CoreService) run() error {
	return c.serve()
}

func Start[T listenEntry](srv *CoreService, addr T) error {
	var address string
	switch v := any(addr).(type) {
	case int:
		address = ":" + strconv.Itoa(v)
	case string:
		address = v
	default:
		return fmt.Errorf("unsupported address type: %T", v)
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	srv.listener = listener
	return srv.runWithListener(listener)
}

func Run(srv *CoreService) error {
	return srv.run()
}
