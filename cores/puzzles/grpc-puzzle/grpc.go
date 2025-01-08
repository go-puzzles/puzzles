package grpcpuzzle

import (
	"context"
	"net"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcLis  net.Listener
	grpcInit = make(chan struct{}, 1)

	gp = &grpcPuzzles{
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

func GrpcSrvListener() net.Listener {
	return grpcLis
}

type grpcPuzzles struct {
	grpcSrv            *grpc.Server
	grpcServersFunc    []func(*grpc.Server)
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	opts               []grpc.ServerOption
}

type grpcPuzzlesOption func(gp *grpcPuzzles)

func WithCoreGrpcUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpcPuzzlesOption {
	return func(gp *grpcPuzzles) {
		gp.unaryInterceptors = append(gp.unaryInterceptors, interceptors...)
		plog.Debugf("Grpc unary interceptors enabled.")
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

func (g *grpcPuzzles) Name() string {
	return "GrpcHandler"
}

func (g *grpcPuzzles) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	if opt.Cmux == nil {
		return errors.New("no http listener specify. please run service with cores.Start")
	}

	if len(g.unaryInterceptors) != 0 {
		g.opts = append(g.opts, grpc.ChainUnaryInterceptor(g.unaryInterceptors...))
	}

	if len(g.streamInterceptors) != 0 {
		g.opts = append(g.opts, grpc.ChainStreamInterceptor(g.streamInterceptors...))
	}

	grpcLis = opt.Cmux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	g.grpcSrv = grpc.NewServer(g.opts...)

	for _, fn := range g.grpcServersFunc {
		fn(g.grpcSrv)
	}

	reflection.Register(g.grpcSrv)
	grpcInit <- struct{}{}
	return g.grpcSrv.Serve(grpcLis)
}

func (g *grpcPuzzles) Stop() error {
	defer plog.Debugf("grpc puzzle stopped...")
	g.grpcSrv.Stop()
	return grpcLis.Close()
}
