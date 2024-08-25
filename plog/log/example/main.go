package main

import (
	"context"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/level"
	"github.com/go-puzzles/puzzles/plog/log"
)

var (
	fileLog = &plog.LogConfig{}
)

func main() {
	fileLog.SetDefault()
	logger := log.New()
	logger.SetOutput(fileLog)
	
	logger.Enable(level.LevelDebug)
	ctx := plog.With(context.Background(), "group")
	logger.Infoc(ctx, "this is a log: %v name=%v", "super", "yong")
	
	ctx = plog.With(ctx, "key", "value")
	logger.Infoc(ctx, "this is key-paire log")
	logger.Debugc(ctx, "this is key-paire debug log")
	logger.Warnc(ctx, "this is key-paire warn log")
	logger.Errorc(ctx, "this is key-paire error log")
	
	logger.Infof("this is a log: %v name=%v", "super", "yong")
	logger.Debugf("this is a log: %v name=%v", "super", "yong")
	logger.Errorf("this is a log: %v name=%v", "super", "yong")
	logger.Warnf("this is a log: %v name=%v", "super", "yong")
	
	logger.Enable(level.LevelError)
	logger.Infof("this is a log: %v name=%v", "super", "yong")
	logger.Debugf("this is a log: %v name=%v", "super", "yong")
	logger.Errorf("this is a log: %v name=%v", "super", "yong")
	logger.Warnf("this is a log: %v name=%v", "super", "yong")
}
