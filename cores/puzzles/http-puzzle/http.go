package httppuzzle

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	basepuzzle "github.com/go-puzzles/puzzles/cores/puzzles/base"
)

type httpPuzzles struct {
	*basepuzzle.BasePuzzle
	httpCors bool
	pattern  string
	router   *mux.Router
}

var (
	hp = &httpPuzzles{
		BasePuzzle: &basepuzzle.BasePuzzle{
			PuzzleName: "HttpPuzzle",
		},
		router: mux.NewRouter(),
	}
)

func WithCoreHttpCORS() cores.ServiceOption {
	return func(o *cores.Options) {
		hp.httpCors = true
		plog.Debugf("Http enable CORS")
		o.RegisterPuzzle(hp)
	}
}

func WithCoreHttpPuzzle(pattern string, handler http.Handler) cores.ServiceOption {
	return func(o *cores.Options) {
		if !strings.HasPrefix(pattern, "/") {
			pattern = "/" + pattern
		}

		hp.pattern = pattern
		hp.router.PathPrefix(pattern).Handler(http.StripPrefix(pattern, handler))
		o.RegisterPuzzle(hp)
	}
}

func (h *httpPuzzles) waitForOtherPuzzles(opt *cores.Options) {
	if opt.WaitPprof != nil {
		select {
		case <-opt.WaitPprof:
		}
	}
}

func (h *httpPuzzles) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	h.waitForOtherPuzzles(opt)

	var handler http.Handler = h.router
	if h.httpCors {
		handler = cors.AllowAll().Handler(handler)
	}

	opt.HttpMux.Handle(hp.pattern+"/", h.router)

	_, port, _ := net.SplitHostPort(opt.ListenerAddr)
	target := fmt.Sprintf("127.0.0.1:%s", port)
	plog.Debugc(ctx, "HttpPuzzle enabled. URL=%s", fmt.Sprintf("http://%s%s", target, h.pattern))

	return nil
}

func (h *httpPuzzles) Stop() (err error) {
	return nil
}
