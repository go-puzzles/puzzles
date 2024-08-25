package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
)

var (
	port = pflags.Int("port", 0, "service run port")
	dura = pflags.Duration("durations", time.Second*5, "worker loop duration")
)

func main() {
	pflags.Parse()
	
	core := cores.NewPuzzleCore(
		cores.WithWorker(func(ctx context.Context) error {
			t := time.NewTicker(dura())
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-t.C:
				}
				
				fmt.Println("running")
			}
		}),
	)
	
	plog.PanicError(cores.Run(core))
}
