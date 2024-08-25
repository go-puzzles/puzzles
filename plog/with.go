package plog

import (
	"context"
	"fmt"
	"io"
	
	logctx "github.com/go-puzzles/puzzles/plog/log-ctx"
)

var (
	badKey = "!BADKEY"
)

func sliceClone(strSlice []string) []string {
	if strSlice == nil {
		return strSlice
	}
	return append(strSlice[:0:0], strSlice...)
}

func cloneLogContext(c *logctx.LogContext) *logctx.LogContext {
	if c == nil {
		return nil
	}
	clone := &logctx.LogContext{
		Group:  c.Group,
		Keys:   sliceClone(c.Keys),
		Values: sliceClone(c.Values),
	}
	
	return clone
}

func getKVParis(kvs []any) (string, string, []any) {
	switch x := kvs[0].(type) {
	case string:
		if len(kvs) == 1 {
			return badKey, x, nil
		}
		return x, fmt.Sprintf("%v", kvs[1]), kvs[2:]
	default:
		return badKey, fmt.Sprintf("%v", x), kvs[1:]
	}
}

func parseRemains(keys, values []string, remains []any) ([]string, []string) {
	var key, val string
	for len(remains) > 0 {
		key, val, remains = getKVParis(remains)
		keys = append(keys, key)
		values = append(values, val)
	}
	
	return keys, values
}

func parseFmtKeyValue(logCtx *logctx.LogContext, msg string, v ...any) (c *logctx.LogContext, err error) {
	if len(v) == 0 {
		logCtx.Group = append(logCtx.Group, msg)
		return logCtx, nil
	}
	
	logCtx.Keys = append(logCtx.Keys, msg)
	logCtx.Values = append(logCtx.Values, fmt.Sprintf("%v", v[0]))
	if len(v) == 1 {
		return logCtx, nil
	}
	
	logCtx.Keys, logCtx.Values = parseRemains(logCtx.Keys, logCtx.Values, v[1:])
	
	return logCtx, nil
}

// With used to store some data in log-ctx
// Ie supports two forms of writing
// 1. With(log-ctx, "group")
// When only msg has a value, it is used as a group
// 2. With(log-ctx, "key1", "value1") or With(log-ctx, "key1", "value1", "key2", "value2")
// When msg and v both have values, With resolves them into key-value pairs
// The parsed results are stored in the LogContext for use by the Logger
func With(c context.Context, msg string, v ...any) context.Context {
	if c == nil {
		c = context.Background()
	}
	
	if msg == "" && len(v) == 0 {
		return c
	}
	
	lc := logctx.GetLogContext(c)
	if lc == nil {
		lc = &logctx.LogContext{}
	}
	
	newLc := cloneLogContext(lc)
	
	newLc, err := parseFmtKeyValue(newLc, msg, v...)
	if err != nil {
		Errorf("with parse error: %v", err)
		return c
	}
	
	return context.WithValue(c, logctx.LogContextKey, newLc)
}

func WithLogger(c context.Context, w io.Writer) context.Context {
	return context.WithValue(c, logctx.LoggerKey, w)
}
