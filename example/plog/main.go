package main

import (
	"context"

	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
)

var (
	logConfFlag = pflags.Struct("log", (*plog.LogConfig)(nil), "")
)

func main() {
	// plog.SetSlog()
	pflags.Parse()

	// logConf := new(plog.LogConfig)
	// plog.PanicError(logConfFlag(logConf))

	// plog.Infof("%v", logConf)
	// plog.EnableLogToFile(logConf)

	// plog.Infof("this is info log")
	// plog.Errorf("this is error log")

	// logger := slog.New()
	// logger.Infof("this is single logger")

	plog.Debugf("this is format log")
	plog.Debugc(context.Background(), "this is context log: %v", 111)
}
