package plog

import (
	"context"
	"io"
	
	"github.com/go-puzzles/puzzles/plog/level"
)

type baseLogger interface {
	Enable(level.Level)
	IsDebug() bool
	SetOutput(io.Writer)
}

type Logger interface {
	baseLogger
	ContextLogger
	Infof(string, ...any)
	Debugf(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
	PanicError(error, ...any)
}

type ContextLogger interface {
	Infoc(context.Context, string, ...any)
	Debugc(context.Context, string, ...any)
	Warnc(context.Context, string, ...any)
	Errorc(context.Context, string, ...any)
}
