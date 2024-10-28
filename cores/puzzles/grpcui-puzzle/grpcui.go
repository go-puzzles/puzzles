package grpcuipuzzle

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/go-puzzles/puzzles/cores"
	grpcpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/grpc-puzzle"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	gup       = &grpcUiPuzzles{}
	grpcuiUrl = "/debug/grpc/ui/"
)

type grpcUiPuzzles struct {
	grpcUILis    net.Listener
	grpcSelfConn *grpc.ClientConn
}

func WithCoreGrpcUI() cores.ServiceOption {
	return func(o *cores.Options) {
		o.RegisterPuzzle(gup)
	}
}

func (g *grpcUiPuzzles) setupGrpcUIRouter(ctx context.Context, serviceName string) (*mux.Router, error) {
	router := mux.NewRouter()

	handler, err := standalone.HandlerViaReflection(ctx, g.grpcSelfConn, serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "start grpcUI")
	}
	router.PathPrefix(grpcuiUrl).Handler(http.StripPrefix(strings.TrimSuffix(grpcuiUrl, "/"), handler))

	return router, nil
}

func (g *grpcUiPuzzles) prepareSelfConnect() error {
	_, port, _ := net.SplitHostPort(grpcpuzzle.GrpcSrvListener().Addr().String())
	target := fmt.Sprintf("127.0.0.1:%s", port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16 * 1024 * 1024)),
	}
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return errors.Wrap(err, "self connect to grpc")
	}
	g.grpcSelfConn = conn
	return nil
}

func (g *grpcUiPuzzles) Name() string {
	return "GrpcUIHandler"
}

func (g *grpcUiPuzzles) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	ctx, cancel := context.WithTimeout(opt.Ctx, time.Second*5)
	defer cancel()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "waitForGrpcServerInit")
	case <-grpcpuzzle.IsGrpcServerInit():
	}

	glis := opt.Cmux.Match(cores.HttpPrefixMatcher(strings.TrimSuffix(grpcuiUrl, "/")))
	g.grpcUILis = glis

	if err := g.prepareSelfConnect(); err != nil {
		return errors.Wrap(err, "prepareSelfConnect")
	}
	router, err := g.setupGrpcUIRouter(opt.Ctx, opt.ServiceName)
	if err != nil {
		return errors.Wrap(err, "setupGrpcUIRouter")
	}

	plog.Infoc(ctx, "GRPCUI enabled. URL=%s", fmt.Sprintf("http://%s%s", g.grpcSelfConn.Target(), grpcuiUrl))
	return http.Serve(glis, router)
}

func (g *grpcUiPuzzles) Stop() error {
	defer plog.Debugf("grpcui puzzle stopped...")
	g.grpcSelfConn.Close()
	return g.grpcUILis.Close()
}
