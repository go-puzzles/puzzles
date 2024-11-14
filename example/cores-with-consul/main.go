package main

import (
	"github.com/go-puzzles/puzzles/cores"
	consulpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/consul-puzzle"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
)

var (
	port = pflags.Int("port", 0, "service run port")
)

func main() {
	pflags.Parse(
		// pflags.WithConsulEnable(),
	)
	
	core := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		consulpuzzle.WithConsulRegister(),
	)
	
	plog.PanicError(cores.Start(core, port()))
}
