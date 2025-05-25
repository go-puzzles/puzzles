package cores

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-puzzles/puzzles/perror"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/robfig/cron/v3"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

type listenerGetter func() net.Listener

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
	cmux     cmux.CMux

	cron     *cron.Cron
	mountFns []mountFn
}

type Options struct {
	Ctx          context.Context
	ListenerAddr string
	ServiceName  string
	Tags         []string
	WaitPprof    chan struct{}

	GrpcListener listenerGetter
	HttpListener listenerGetter

	HttpMux     *http.ServeMux
	HttpHandler http.Handler

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

	coreOpts := &Options{
		Ctx:         ctx,
		HttpMux:     http.NewServeMux(),
		ServiceName: pflags.GetServiceName(),
		puzzles:     make(map[string]Puzzle),
		workers:     make([]Worker, 0),
	}
	coreOpts.HttpHandler = coreOpts.HttpMux

	for _, opt := range opts {
		opt(coreOpts)
	}

	return &CoreService{
		ctx:      ctx,
		opts:     coreOpts,
		cancel:   cancel,
		mountFns: make([]mountFn, 0),
	}
}

func (c *CoreService) options() *Options {
	return c.opts
}

func (c *CoreService) isHttpRun() bool {
	return c.listener != nil
}

func (c *CoreService) listenerMatch(matcher ...cmux.Matcher) listenerGetter {
	var lst net.Listener
	var once sync.Once

	return func() net.Listener {
		once.Do(func() {
			lst = c.cmux.Match(matcher...)
		})

		return lst
	}
}

func (c *CoreService) listenerMatchWithWriter(matcher ...cmux.MatchWriter) listenerGetter {
	var lst net.Listener
	var once sync.Once

	return func() net.Listener {
		once.Do(func() {
			lst = c.cmux.MatchWithWriters(matcher...)
		})

		return lst
	}
}

func (c *CoreService) runMountFn() error {
	grp, ctx := errgroup.WithContext(c.ctx)

	for _, mount := range c.mountFns {
		mf := mount
		cc := plog.With(ctx, "Worker", mf.name)
		grp.Go(func() (err error) {
			err = waitContext(cc, false, mf.fn)
			if err != nil {
				plog.Errorc(cc, "handle mount error: %v", err)
				if mf.daemon {
					return perror.WrapError(500, err, "waitContext")
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
		c.cmux = cmux.New(c.listener)
		c.opts.GrpcListener = c.listenerMatchWithWriter(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
		c.opts.HttpListener = c.listenerMatchWithWriter(
			HTTP1FastMatchWriter(),
			HTTP2MatchWithHeaderExclude("content-type", "application/grpc"),
		)

		c.mountFns = append(c.mountFns, c.listenHttp())
		c.mountFns = append(c.mountFns, c.listenCmux())
	}
	c.wrapPuzzles()
	c.wrapWorker()
	c.welcome()
	return c.runMountFn()
}

func (c *CoreService) listenHttp() mountFn {
	return mountFn{
		fn: func(ctx context.Context) (err error) {
			defer func() {
				if errors.Is(err, net.ErrClosed) {
					err = nil
				}

				if err != nil {
					plog.Errorc(ctx, "HttpListener serve error: %v", err)
				}
			}()

			err = http.Serve(c.opts.HttpListener(), c.opts.HttpHandler)
			return
		},
		name:   "HttpListener",
		daemon: true,
	}
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

			return c.cmux.Serve()
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
		return perror.PackError(500, fmt.Sprintf("unsupported address type: %T", v))
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return perror.WrapError(500, err, "failed to start listener")
	}

	srv.listener = listener
	return srv.runWithListener(listener)
}

func Run(srv *CoreService) error {
	return srv.run()
}
