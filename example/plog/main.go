package main

import (
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/slog"
)

var (
	logConfFlag = pflags.Struct("log", (*plog.LogConfig)(nil), "")
)

func main() {
	plog.SetSlog()
	pflags.Parse()

	logConf := new(plog.LogConfig)
	plog.PanicError(logConfFlag(logConf))

	plog.Infof("%v", logConf)
	plog.EnableLogToFile(logConf)

	plog.Infof("this is info log")
	plog.Errorf("this is error log")

	logger := slog.New()
	logger.Infof("this is single logger")
}
