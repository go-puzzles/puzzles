// File:		pprof.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pprofpuzzle

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"

	basepuzzle "github.com/go-puzzles/puzzles/cores/puzzles/base"
)

var (
	pp = &pprofPuzzles{
		BasePuzzle: &basepuzzle.BasePuzzle{
			PuzzleName: "PprofPuzzle",
		},
	}
	pprofUrl = "/debug/pprof/"
)

type pprofPuzzles struct {
	*basepuzzle.BasePuzzle
	lis net.Listener
}

func WithCorePprof() cores.ServiceOption {
	return func(o *cores.Options) {
		o.WaitPprof = make(chan struct{})
		o.RegisterPuzzle(pp)
	}
}

func (p *pprofPuzzles) StartPuzzle(ctx context.Context, opts *cores.Options) error {
	_, port, _ := net.SplitHostPort(opts.ListenerAddr)
	target := fmt.Sprintf("localhost:%s", port)

	router := mux.NewRouter()

	registerHandler := func(path string) *mux.Route {
		return router.Path(path)
	}

	registerHandler("/").HandlerFunc(pprof.Index)
	registerHandler("/allocs").Handler(pprof.Handler("allocs"))
	registerHandler("/block").Handler(pprof.Handler("block"))
	registerHandler("/cmdline").HandlerFunc(pprof.Cmdline)
	registerHandler("/goroutine").Handler(pprof.Handler("goroutine"))
	registerHandler("/heap").Handler(pprof.Handler("heap"))
	registerHandler("/mutex").Handler(pprof.Handler("mutex"))
	registerHandler("/profile").HandlerFunc(pprof.Profile)
	registerHandler("/threadcreate").Handler(pprof.Handler("threadcreate"))
	registerHandler("/trace").HandlerFunc(pprof.Trace)
	registerHandler("/symbol").HandlerFunc(pprof.Symbol)

	opts.HttpMux.Handle(pprofUrl, http.StripPrefix("/debug/pprof", router))
	opts.WaitPprof <- struct{}{}

	plog.Debugc(ctx, "PprofPuzzle enabled. URL=%s", fmt.Sprintf("http://%s%s", target, pprofUrl))
	return nil
}

func (p *pprofPuzzles) Stop() (err error) {
	defer func() {
		plog.Debugf("pprof puzzle stopped...")
		if errors.Is(err, net.ErrClosed) {
			err = nil
		}
	}()
	return p.lis.Close()
}
