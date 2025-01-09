package grpcpuzzle

import (
	"context"
	"net"

	"github.com/go-puzzles/puzzles/cores"
	basepuzzle "github.com/go-puzzles/puzzles/cores/puzzles/base"
	"github.com/go-puzzles/puzzles/plog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcLis  net.Listener
	grpcInit = make(chan struct{}, 1)

	gp = &grpcPuzzles{
		BasePuzzle: &basepuzzle.BasePuzzle{
			PuzzleName: "GrpcHandler",
		},
		opts: make([]grpc.ServerOption, 0),
		unaryInterceptors: []grpc.UnaryServerInterceptor{
			unaryServerLoggerInterceptor,
		},
		streamInterceptors: []grpc.StreamServerInterceptor{
			StreamServerLoggerInterceptor,
		},
	}
)

func IsGrpcServerInit() <-chan struct{} {
	return grpcInit
}

type grpcPuzzles struct {
	*basepuzzle.BasePuzzle
	grpcSrv            *grpc.Server
	grpcServersFunc    []func(*grpc.Server)
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	opts               []grpc.ServerOption
}

type grpcPuzzlesOption func(gp *grpcPuzzles)

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpcPuzzlesOption {
	return func(gp *grpcPuzzles) {
		gp.unaryInterceptors = append(gp.unaryInterceptors, interceptors...)
		plog.Debugf("Grpc unary interceptors enabled.")
	}
}

func WithServerOptions(opts ...grpc.ServerOption) grpcPuzzlesOption {
	return func(gp *grpcPuzzles) {
		gp.opts = append(gp.opts, opts...)
	}
}

func WithCoreGrpcPuzzle(grpcSrv func(srv *grpc.Server), opts ...grpcPuzzlesOption) cores.ServiceOption {
	return func(o *cores.Options) {
		gp.grpcServersFunc = append(gp.grpcServersFunc, grpcSrv)
		o.Tags = append(o.Tags, "grpc")

		for _, opt := range opts {
			opt(gp)
		}

		o.RegisterPuzzle(gp)
		plog.Debugf("Grpc server enabled.")
	}
}

func (g *grpcPuzzles) Before(opt *cores.Options) error {
	if len(g.unaryInterceptors) != 0 {
		g.opts = append(g.opts, grpc.ChainUnaryInterceptor(g.unaryInterceptors...))
	}

	if len(g.streamInterceptors) != 0 {
		g.opts = append(g.opts, grpc.ChainStreamInterceptor(g.streamInterceptors...))
	}

	g.grpcSrv = grpc.NewServer(g.opts...)
	for _, fn := range g.grpcServersFunc {
		fn(g.grpcSrv)
	}

	return nil
}

func (g *grpcPuzzles) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	grpcInit <- struct{}{}

	reflection.Register(g.grpcSrv)
	return g.grpcSrv.Serve(opt.GrpcListener())
}

func (g *grpcPuzzles) Stop() error {
	defer plog.Debugf("grpc puzzle stopped...")
	g.grpcSrv.Stop()
	return grpcLis.Close()
}
