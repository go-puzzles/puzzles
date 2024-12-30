package plog

import (
	"context"
	"io"
	"time"

	"github.com/go-puzzles/puzzles/plog/level"
	"github.com/go-puzzles/puzzles/plog/log"
	"github.com/go-puzzles/puzzles/plog/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger Logger
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600)
	logger = log.New(log.WithCalldepth(4))
}

func SetSlog() {
	logger = slog.New(slog.WithCalldepth(5))
}

func SetLogger(l Logger) {
	logger = l
}

func IsDebug() bool {
	return logger.IsDebug()
}

func EnableLogToFile(jackLog *LogConfig) {
	logger.SetOutput((*lumberjack.Logger)(jackLog))
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func Enable(l level.Level) {
	logger.Enable(l)
}

func Errorf(msg string, v ...any) {
	logger.Errorf(msg, v...)
}

func Warnf(msg string, v ...any) {
	logger.Warnf(msg, v...)
}

func Infof(msg string, v ...any) {
	logger.Infof(msg, v...)
}

func Debugf(msg string, v ...any) {
	logger.Debugf(msg, v...)
}

func Fatalf(msg string, v ...any) {
	logger.Fatalf(msg, v...)
}

func Infoc(ctx context.Context, msg string, v ...any) {
	logger.Infoc(ctx, msg, v...)
}

func Debugc(ctx context.Context, msg string, v ...any) {
	logger.Debugc(ctx, msg, v...)
}

func Warnc(ctx context.Context, msg string, v ...any) {
	logger.Warnc(ctx, msg, v...)
}

func Errorc(ctx context.Context, msg string, v ...any) {
	logger.Errorc(ctx, msg, v...)
}

func Fatalc(ctx context.Context, msg string, v ...any) {
	logger.Fatalc(ctx, msg, v...)
}

func PanicError(err error, v ...any) {
	if err == nil {
		return
	}

	logger.PanicError(err, v...)
}
