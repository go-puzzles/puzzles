package main

import (
	"context"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/level"
	"github.com/go-puzzles/puzzles/plog/slog"
)

var (
	fileLog = &plog.LogConfig{}
)

func main() {
	fileLog.SetDefault()
	
	logger := slog.New()
	logger.SetOutput(fileLog)
	logger.Enable(level.LevelDebug)
	
	logger.Infoc(context.Background(), "this is a message")
	logger.Debugc(context.Background(), "this is a message")
	logger.Warnc(context.Background(), "this is a message")
	logger.Errorc(context.Background(), "this is a message")
	
	ctx := plog.With(context.Background(), "group")
	ctx = plog.With(ctx, "handler")
	logger.Infoc(ctx, "this is a ctx message")
	logger.Debugc(ctx, "this is a ctx message")
	logger.Warnc(ctx, "this is a ctx message")
	logger.Errorc(ctx, "this is a ctx message")
	
	ctx = plog.With(ctx, "name", "hoven")
	logger.Infoc(ctx, "this is a key-value message")
	ctx = plog.With(ctx, "city", "shenzhen")
	logger.Debugc(ctx, "this is a key-value message")
	ctx = plog.With(ctx, "age", 16, "phone", 12345)
	logger.Warnc(ctx, "this is a key-value message")
	
	ctx = plog.With(ctx, "email", "xxxx", "somethingBad")
	logger.Errorc(ctx, "this is a key-value with bad key message")
}
