package httppuzzle

import (
	"context"
	"net"
	"net/http"
	"strings"
	
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
)

type httpPuzzles struct {
	lis      net.Listener
	httpCors bool
	router   *mux.Router
}

var (
	hp = &httpPuzzles{
		router: mux.NewRouter(),
	}
)

func WithCoreHttpCORS() cores.ServiceOption {
	return func(o *cores.Options) {
		hp.httpCors = true
		plog.Infof("Http enable CORS")
		o.RegisterPuzzle(hp)
	}
}

func WithCoreHttpPuzzle(pattern string, handler http.Handler) cores.ServiceOption {
	return func(o *cores.Options) {
		if !strings.HasPrefix(pattern, "/") {
			pattern = "/" + pattern
		}
		defer plog.Infof("Registered http endpoint prefix. Prefix=%s", pattern)
		
		hp.router.PathPrefix(pattern).Handler(http.StripPrefix(strings.TrimSuffix(pattern, "/"), handler))
		o.RegisterPuzzle(hp)
	}
}

func (h *httpPuzzles) Name() string {
	return "HttpHandler"
}

func (h *httpPuzzles) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	if opt.Cmux == nil {
		return errors.New("no http listener specify. please run service with cores.Start")
	}
	
	lis := opt.Cmux.Match(cmux.HTTP1Fast(), cmux.HTTP2())
	h.lis = lis
	
	var handler http.Handler = h.router
	if h.httpCors {
		handler = cors.AllowAll().Handler(handler)
	}
	return http.Serve(lis, handler)
}

func (h *httpPuzzles) Stop() error {
	plog.Debugf("http puzzles stop...")
	return h.lis.Close()
}
