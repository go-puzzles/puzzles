package main

import (
	"net/http"

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
		httppuzzle.WithCoreHttpPuzzle("/api", router),
	)

	plog.PanicError(cores.Start(core, port()))
}
