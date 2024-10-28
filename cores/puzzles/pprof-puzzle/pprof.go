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
	"net"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"
)

var (
	pp       = &pprofPuzzles{}
	pprofUrl = "/debug/pprof/"
)

type pprofPuzzles struct {
	lis net.Listener
}

func WithCorePprof() cores.ServiceOption {
	return func(o *cores.Options) {
		plog.Infof("Pprof eanble. Prefix=%s", pprofUrl)
		o.WaitPprof = make(chan struct{})
		o.RegisterPuzzle(pp)
	}
}

func (p *pprofPuzzles) Name() string {
	return "pprofPuzzle"
}

func (p *pprofPuzzles) StartPuzzle(ctx context.Context, opts *cores.Options) error {
	if opts.Cmux == nil {
		return errors.New("no http listener specify. please run service with cores.Start")
	}

	pLis := opts.Cmux.Match(cores.HttpPrefixMatcher(strings.TrimSuffix(pprofUrl, "/")))
	p.lis = pLis

	handler := mux.NewRouter()
	handler.HandleFunc(pprofUrl, pprof.Index)
	handler.Handle(pprofUrl+"allocs", pprof.Handler("allocs"))
	handler.Handle(pprofUrl+"block", pprof.Handler("block"))
	handler.HandleFunc(pprofUrl+"cmdline", pprof.Cmdline)
	handler.Handle(pprofUrl+"goroutine", pprof.Handler("goroutine"))
	handler.Handle(pprofUrl+"heap", pprof.Handler("heap"))
	handler.Handle(pprofUrl+"mutex", pprof.Handler("mutex"))
	handler.HandleFunc(pprofUrl+"profile", pprof.Profile)
	handler.Handle(pprofUrl+"threadcreate", pprof.Handler("threadcreate"))
	handler.HandleFunc(pprofUrl+"trace", pprof.Trace)
	handler.HandleFunc(pprofUrl+"symbol", pprof.Symbol)

	opts.WaitPprof <- struct{}{}
	return http.Serve(pLis, handler)
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
