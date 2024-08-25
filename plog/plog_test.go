package plog

import (
	"context"
	"testing"
	
	"github.com/go-puzzles/puzzles/plog/level"
	"github.com/go-puzzles/puzzles/plog/slog"
)

func TestPlog(t *testing.T) {
	ctx := context.Background()
	ctx = With(ctx, "group")
	ctx = With(ctx, "name", "hoven")
	Infoc(ctx, "this is context log")
	Debugc(ctx, "this is context log")
	Warnc(ctx, "this is context log")
	Errorc(ctx, "this is context log")
	
	Infoc(ctx, "this is key-value log", "age", 16)
	Debugc(ctx, "this is key-value log", "age", 16)
	Errorc(ctx, "this is key-value log", "age", 16)
	Warnc(ctx, "this is key-value log", "age", 16, "city")
	
	Enable(level.LevelDebug)
	Infoc(ctx, "this is key-value log with format. format: %v", "formatMsg", "age", 16)
	Debugc(ctx, "this is key-value log with format. format: %v", "formatMsg", "age", 16)
	Errorc(ctx, "this is key-value log with format. format: %v", "formatMsg", "age", 16, "city")
	Warnc(ctx, "this is key-value log with format. format: %v", "formatMsg", "age", 16)
	
	SetLogger(slog.New())
	
	Infoc(ctx, "this is slog context log")
	Debugc(ctx, "this is slog context log")
	Warnc(ctx, "this is slog context log")
	Errorc(ctx, "this is slog context log")
	
	Infoc(ctx, "this is slog key-value log", "age", 16)
	Debugc(ctx, "this is slog key-value log", "age", 16)
	Errorc(ctx, "this is slog key-value log", "age", 16)
	Warnc(ctx, "this is slog key-value log", "age", 16, "city")
	
	Enable(level.LevelDebug)
	Infoc(ctx, "this is slog key-value log with format. format: %v", "formatMsg", "age", 16)
	Debugc(ctx, "this is slog key-value log with format. format: %v", "formatMsg", "age", 16)
	Errorc(ctx, "this is slog key-value log with format. format: %v", "formatMsg", "age", 16, "city")
	Warnc(ctx, "this is slog key-value log with format. format: %v", "formatMsg", "age", 12)
}
