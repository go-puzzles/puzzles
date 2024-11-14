package pgorm

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	prefix                    string
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
	logMod                    logger.LogLevel
	logger                    plog.Logger
	
	traceStr, traceErrStr, traceWarnStr string
}

type GormLoggerOption func(g *gormLogger)

func WithPrefix(prefix string) GormLoggerOption {
	return func(g *gormLogger) {
		g.prefix = prefix
	}
}

func WithIgnoreRecordNotFound() GormLoggerOption {
	return func(g *gormLogger) {
		g.ignoreRecordNotFoundError = true
	}
}

func WithSlowThreshold(dur time.Duration) GormLoggerOption {
	return func(g *gormLogger) {
		g.slowThreshold = dur
	}
}

func NewGormLogger(opts ...GormLoggerOption) *gormLogger {
	var (
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)
	
	l := &gormLogger{
		logger:       log.New(),
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
	
	for _, opt := range opts {
		opt(l)
	}
	
	return l
}

func (gl *gormLogger) wrapPrefix(ctx context.Context) context.Context {
	if gl.prefix != "" {
		ctx = plog.With(ctx, "[%v]", gl.prefix)
	}
	
	return ctx
}

func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *gl
	newLogger.logMod = level
	return &newLogger
}

func (gl *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	gl.logger.Infoc(gl.wrapPrefix(ctx), msg, data)
}

func (gl *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	gl.logger.Warnc(gl.wrapPrefix(ctx), msg, data)
}

func (gl *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	gl.logger.Errorc(gl.wrapPrefix(ctx), msg, data)
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if !plog.IsDebug() {
		return
	}
	
	ctx = gl.wrapPrefix(ctx)
	elapsed := time.Since(begin)
	
	switch {
	case err != nil && (!errors.Is(err, logger.ErrRecordNotFound) || !gl.ignoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			gl.logger.Errorc(ctx, gl.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			gl.logger.Errorc(ctx, gl.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > gl.slowThreshold && gl.slowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", gl.slowThreshold)
		if rows == -1 {
			gl.logger.Warnc(ctx, gl.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			gl.logger.Warnc(ctx, gl.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case gl.logMod == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			gl.logger.Infoc(ctx, gl.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			gl.logger.Infoc(ctx, gl.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
