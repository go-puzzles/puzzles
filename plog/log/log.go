package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-logfmt/logfmt"
	"github.com/go-puzzles/puzzles/plog/level"
	logctx "github.com/go-puzzles/puzzles/plog/log-ctx"
	"github.com/go-puzzles/puzzles/plog/parser"
)

type logable interface {
	Output(calldepth int, s string) error
}

type Logger struct {
	callDepth int
	logLevel  level.Level

	stdout io.Writer
	stderr io.Writer

	infoLog  *log.Logger
	debugLog *log.Logger
	warnLog  *log.Logger
	errLog   *log.Logger
}

type Option func(*Logger)

func WithWriter(stdout, stderr io.Writer) Option {
	return func(l *Logger) {
		l.stdout = stdout
		l.stderr = stderr
	}
}

func WithCalldepth(d int) Option {
	return func(l *Logger) {
		l.callDepth = d
	}
}

func New(opts ...Option) *Logger {
	l := &Logger{
		callDepth: 3,
		stdout:    os.Stdout,
		stderr:    os.Stderr,
	}
	l.Enable(level.LevelInfo)

	for _, opt := range opts {
		opt(l)
	}

	l.initLogger()
	return l
}

func (l *Logger) initLogger() {
	stdout := l.stdout
	stderr := l.stderr

	logFlagNoFile := log.LstdFlags | log.Ldate | log.Ltime
	logFlagWithFile := log.LstdFlags | log.Ldate | log.Ltime | log.Lshortfile

	l.infoLog = log.New(stdout, formatPrefix("[INFO]"), logFlagNoFile)
	l.debugLog = log.New(stdout, formatPrefix("[DEBUG]"), logFlagWithFile)
	l.errLog = log.New(stderr, formatPrefix("[ERROR]"), logFlagWithFile)
	l.warnLog = log.New(stdout, formatPrefix("[WARN]"), logFlagNoFile)
}

// SetOutput set log output position
func (l *Logger) SetOutput(w io.Writer) {
	l.infoLog.SetOutput(w)
	l.debugLog.SetOutput(w)
	l.errLog.SetOutput(w)
	l.warnLog.SetOutput(w)
}

func formatPrefix(prefix string) string {
	return fmt.Sprintf("%-8s", prefix)
}

func (l *Logger) logFmt(lc *logctx.LogContext, msg string, v ...any) string {
	s, keys, values, err := parser.ParseFmtKeyValue(msg, v...)
	if err != nil {
		return msg + " " + err.Error()
	}

	msg = s
	keys = append(keys, lc.Keys...)
	values = append(values, lc.Values...)

	if len(lc.Group) != 0 {
		grpMsg := strings.Join(lc.Group, ":")
		msg = fmt.Sprintf("%s: %s", grpMsg, msg)
	}

	var buf bytes.Buffer

	encoder := logfmt.NewEncoder(&buf)

	for i := 0; i < len(keys); i++ {
		encoder.EncodeKeyval(keys[i], values[i])
	}
	str := buf.String()
	if str == "" {
		return msg
	}

	return msg + " " + str
}

func (l *Logger) logc(c context.Context, lb logable, msg string, v ...any) {
	lc := logctx.GetLogContext(c)
	if lc == nil {
		lc = &logctx.LogContext{}
	}
	msg = l.logFmt(lc, msg, v...)
	for _, line := range strings.Split(msg, "\n") {
		lb.Output(l.callDepth, line)
	}
}

func (l *Logger) Enable(lev level.Level) {
	l.logLevel = lev
}

func (l *Logger) IsDebug() bool {
	return l.logLevel == level.LevelDebug
}

func (l *Logger) checkLevel(lev level.Level) bool {
	return lev >= l.logLevel
}

func (l *Logger) Infoc(ctx context.Context, msg string, v ...any) {
	if !l.checkLevel(level.LevelInfo) {
		return
	}

	l.logc(ctx, l.infoLog, msg, v...)
}

func (l *Logger) Debugc(ctx context.Context, msg string, v ...any) {
	if !l.checkLevel(level.LevelDebug) {
		return
	}

	l.logc(ctx, l.debugLog, msg, v...)
}

func (l *Logger) Warnc(ctx context.Context, msg string, v ...any) {
	if !l.checkLevel(level.LevelWarn) {
		return
	}

	l.logc(ctx, l.warnLog, msg, v...)
}

func (l *Logger) Errorc(ctx context.Context, msg string, v ...any) {
	if !l.checkLevel(level.LevelError) {
		return
	}

	l.logc(ctx, l.errLog, msg, v...)
}

func (l *Logger) Fatalc(ctx context.Context, msg string, v ...any) {
	l.logc(ctx, l.errLog, msg, v...)
	os.Exit(1)
}

func (l *Logger) Infof(msg string, v ...any) {
	l.Infoc(context.Background(), msg, v...)
}

func (l *Logger) Debugf(msg string, v ...any) {
	l.Debugc(context.Background(), msg, v...)
}

func (l *Logger) Warnf(msg string, v ...any) {
	l.Warnc(context.Background(), msg, v...)
}

func (l *Logger) Errorf(msg string, v ...any) {
	l.Errorc(context.Background(), msg, v...)
}

func (l *Logger) Fatalf(msg string, v ...any) {
	l.Fatalc(context.Background(), msg, v...)
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

	l.Errorc(context.Background(), s)
	panic(err)
}
