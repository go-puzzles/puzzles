package logctx

import (
	"context"
	"io"
)

type ContextKey string

var (
	LogContextKey ContextKey = "plogContext"
	LoggerKey     ContextKey = "plogLogger"
)

type LogContext struct {
	Group  []string
	Keys   []string
	Values []string
}

func GetLogContext(c context.Context) *LogContext {
	if c == nil {
		return nil
	}

	val := c.Value(LogContextKey)
	lc, ok := val.(*LogContext)
	if !ok {
		return nil
	}
	return lc
}

func ExtractLogger(ctx context.Context) io.Writer {
	if ctx == nil {
		return nil
	}
	i := ctx.Value(LoggerKey)
	if i == nil {
		return nil
	}
	w, ok := i.(io.Writer)
	if !ok {
		return nil
	}
	return w
}
