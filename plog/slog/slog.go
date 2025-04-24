package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-puzzles/puzzles/plog/level"
	logctx "github.com/go-puzzles/puzzles/plog/log-ctx"
	"github.com/go-puzzles/puzzles/plog/parser"
	"github.com/go-puzzles/puzzles/plog/slog/handler"
)

type logable func(ctx context.Context, msg string, args ...any)

type Logger struct {
	logger          *slog.Logger
	logLevel        level.Level
	slogLoggerLevel *slog.LevelVar
	slogOpt         *slog.HandlerOptions
	callDepth       int
}

type Option func(*Logger)

func WithCalldepth(d int) Option {
	return func(l *Logger) {
		l.callDepth = d
	}
}

func New(opts ...Option) *Logger {
	return NewSlogPrettyLogger(os.Stdout, opts...)
}

func NewWithHandler(handler slog.Handler, opt *slog.HandlerOptions, opts ...Option) *Logger {
	lev := &slog.LevelVar{}
	lev.Set(slog.LevelInfo)

	opt.Level = lev
	logger := slog.New(handler)
	l := &Logger{
		logger:          logger,
		logLevel:        level.LevelInfo,
		slogLoggerLevel: lev,
		slogOpt:         opt,
		callDepth:       4,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func NewSlogPrettyLogger(w io.Writer, opts ...Option) *Logger {
	slogOpt := &slog.HandlerOptions{}
	handler := handler.NewPrettyHandler(w, slogOpt)
	return NewWithHandler(handler, slogOpt, opts...)
}

func NewSlogTextLogger(w io.Writer, opts ...Option) *Logger {
	slogOpt := &slog.HandlerOptions{}
	handler := slog.NewTextHandler(w, slogOpt)
	return NewWithHandler(handler, slogOpt, opts...)
}

func NewSlogJsonLogger(w io.Writer, opts ...Option) *Logger {
	slogOpt := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(w, slogOpt)
	return NewWithHandler(handler, slogOpt, opts...)
}

func relativeToGOROOT(path string) string {
	gopath := os.Getenv("GOPATH")
	rp, err := filepath.Rel(gopath, path)
	if err != nil {
		return path
	}
	return rp
}

func (sl *Logger) getSrouce() string {
	_, file, line, _ := runtime.Caller(sl.callDepth)
	fs := strings.Split(file, "/")
	return fmt.Sprintf("%s:%d", fs[len(fs)-1], line)
}

func (l *Logger) Enable(lev level.Level) {
	l.logLevel = lev
	l.slogLoggerLevel.Set(slog.Level(lev))
	return
}

func (l *Logger) IsDebug() bool {
	return l.logLevel == level.LevelDebug
}

func (l *Logger) SetOutput(w io.Writer) {
	handler := slog.NewJSONHandler(w, l.slogOpt)
	l.logger = slog.New(handler)
}

func (l *Logger) logc(c context.Context, lb logable, msg string, v ...any) {
	lc := logctx.GetLogContext(c)
	if lc == nil {
		lc = &logctx.LogContext{}
	}

	msg, args := l.logFmt(lc, msg, v...)

	lb(c, msg, args...)
}

func (l *Logger) logFmt(lc *logctx.LogContext, msg string, v ...any) (string, []any) {
	s, keys, values, err := parser.ParseFmtKeyValue(msg, v...)
	if err != nil {
		return msg + " " + err.Error(), nil
	}

	msg = s
	keys = append(keys, lc.Keys...)
	values = append(values, lc.Values...)
	var args []any

	for idx, key := range keys {
		args = append(args, key, values[idx])
	}
	args = append(args, "source", l.getSrouce())

	for _, g := range lc.Group {
		msg = fmt.Sprintf("[%s] %s", g, msg)
	}

	return msg, args
}

func (l *Logger) Infoc(ctx context.Context, msg string, v ...any) {
	l.logc(ctx, l.logger.InfoContext, msg, v...)
}

func (l *Logger) Debugc(ctx context.Context, msg string, v ...any) {
	l.logc(ctx, l.logger.DebugContext, msg, v...)
}

func (l *Logger) Warnc(ctx context.Context, msg string, v ...any) {
	l.logc(ctx, l.logger.WarnContext, msg, v...)
}

func (l *Logger) Errorc(ctx context.Context, msg string, v ...any) {
	l.logc(ctx, l.logger.ErrorContext, msg, v...)
}

func (l *Logger) Fatalc(ctx context.Context, msg string, v ...any) {
	l.logc(ctx, l.logger.ErrorContext, msg, v...)
	os.Exit(1)
}

func (l *Logger) Infof(msg string, v ...any) {
	l.logc(context.TODO(), l.logger.InfoContext, msg, v...)
}

func (l *Logger) Debugf(msg string, v ...any) {
	l.logc(context.TODO(), l.logger.DebugContext, msg, v...)
}

func (l *Logger) Warnf(msg string, v ...any) {
	l.logc(context.TODO(), l.logger.WarnContext, msg, v...)
}

func (l *Logger) Errorf(msg string, v ...any) {
	l.logc(context.TODO(), l.logger.ErrorContext, msg, v...)
}

func (l *Logger) Fatalf(msg string, v ...any) {
	l.logc(context.TODO(), l.logger.ErrorContext, msg, v...)
	os.Exit(1)

}

func (l *Logger) PanicError(err error, v ...any) {
	if err == nil {
		return
	}

	var s string
	if len(v) > 0 {
		s = err.Error() + ":" + fmt.Sprint(v...)
	} else {
		s = err.Error()
	}

	l.logc(context.TODO(), l.logger.ErrorContext, s)
	panic(err)
}
