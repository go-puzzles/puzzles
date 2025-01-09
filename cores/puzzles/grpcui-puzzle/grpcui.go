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
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	basepuzzle "github.com/go-puzzles/puzzles/cores/puzzles/base"
	grpcpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/grpc-puzzle"
)

var (
	gup = &grpcUiPuzzles{
		BasePuzzle: &basepuzzle.BasePuzzle{
			PuzzleName: "GrpcuiPuzzle",
		},
	}
	grpcuiUrl = "/debug/grpc/ui/"
)

type grpcUiPuzzles struct {
	*basepuzzle.BasePuzzle
	grpcSelfConn *grpc.ClientConn
}

func WithCoreGrpcUI() cores.ServiceOption {
	return func(o *cores.Options) {
		o.RegisterPuzzle(gup)
	}
}

func (g *grpcUiPuzzles) prepareSelfConnect(lisAddr string) error {
	_, port, _ := net.SplitHostPort(lisAddr)
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

func (g *grpcUiPuzzles) Before(opt *cores.Options) error {
	if err := g.prepareSelfConnect(opt.ListenerAddr); err != nil {
		return errors.Wrap(err, "prepareSelfConnect")
	}

	return nil
}

func (g *grpcUiPuzzles) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "waitForGrpcServerInit")
	case <-grpcpuzzle.IsGrpcServerInit():
	}

	handler, err := standalone.HandlerViaReflection(ctx, g.grpcSelfConn, opt.ServiceName)
	if err != nil {
		return errors.Wrap(err, "start grpcUI")
	}

	opt.HttpMux.Handle(grpcuiUrl, http.StripPrefix(strings.TrimSuffix(grpcuiUrl, "/"), handler))

	plog.Debugc(ctx, "GrpcuiPuzzle enabled. URL=%s", fmt.Sprintf("http://%s%s", g.grpcSelfConn.Target(), grpcuiUrl))
	return nil
}

func (g *grpcUiPuzzles) Stop() error {
	defer plog.Debugf("grpcui puzzle stopped...")
	return g.grpcSelfConn.Close()
}
