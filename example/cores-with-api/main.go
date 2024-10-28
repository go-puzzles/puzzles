package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-puzzles/puzzles/cores"
	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
	pprofpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/pprof-puzzle"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"
)

var (
	port = pflags.IntRequired("port", "service run port")
)

func main() {
	pflags.Parse()
	router := mux.NewRouter()

	router.Path("/hello").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	core := cores.NewPuzzleCore(
		pprofpuzzle.WithCorePprof(),
		cores.WithNameWorker("worker1", func(ctx context.Context) error {
			time.Sleep(time.Second * 5)
			plog.Infoc(ctx, "worker stop")
			return nil
		}),
		cores.WithCronWorker("@every 2s", func(ctx context.Context) error {
			plog.Infoc(ctx, "worker start")
			return nil
		}),
		httppuzzle.WithCoreHttpPuzzle("/api", router),
	)

	plog.PanicError(cores.Start(core, port()))
}
