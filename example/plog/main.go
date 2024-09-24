package main

import (
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
)

var (
	logConfFlag = pflags.Struct("log", (*plog.LogConfig)(nil), "")
)

func main() {
	pflags.Parse()

	logConf := new(plog.LogConfig)
	plog.PanicError(logConfFlag(logConf))

	plog.Infof("%v", logConf)
	plog.EnableLogToFile(logConf)

	plog.Infof("this is info log")
	plog.Errorf("this is error log")
}
