// File:		prettylog.go
// Created by:	Hoven
// Created on:	2025-04-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97

	timeFormat = "2006/01/02 15:04:05"
)

var levelWidth = 15

func colorize(colorCode int, v string, enableColor bool) string {
	if !enableColor {
		return v
	}
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

type PrettyHandler struct {
	attrs  []slog.Attr
	groups []string
	w      io.Writer
	m      *sync.Mutex
	opts   *slog.HandlerOptions
	isTerm bool
}

type OptionFunc func(*PrettyHandler)

func checkIsTerm(w io.Writer) bool {
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		return false
	}

	isTerm := true

	if w, ok := w.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	return isTerm
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions, options ...OptionFunc) *PrettyHandler {
	h := &PrettyHandler{
		m:      &sync.Mutex{},
		w:      w,
		attrs:  []slog.Attr{},
		groups: []string{},
		opts:   opts,
		isTerm: true,
	}

	if h.w == nil {
		h.w = os.Stdout
	}

	h.isTerm = checkIsTerm(h.w)
	if !h.isTerm {
		levelWidth = 6
	}

	for _, option := range options {
		option(h)
	}

	if h.opts == nil {
		h.opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}

	return h
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.opts.Level == nil {
		return true
	}
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	h2 := h.clone()
	for _, a := range attrs {
		if a.Equal(slog.Attr{}) {
			continue
		}
		h2.attrs = append(h2.attrs, a)
	}
	return h2
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *PrettyHandler) clone() *PrettyHandler {
	h2 := &PrettyHandler{
		w:      h.w,
		m:      h.m,
		opts:   h.opts,
		isTerm: h.isTerm,
	}

	h2.attrs = make([]slog.Attr, len(h.attrs))
	copy(h2.attrs, h.attrs)

	h2.groups = make([]string, len(h.groups))
	copy(h2.groups, h.groups)

	return h2
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	levelStr := h.formatLevel(r.Level)

	attrs, err := h.formatAttrs(r)
	if err != nil {
		return err
	}

	var attrsBytes []byte
	if len(attrs) == 0 {
		return h.writeLog(r.Time, levelStr, r.Message, attrsBytes)
	}

	if h.isTerm {
		attrsBytes, err = h.formatTermAttrs(attrs)
	} else {
		attrsBytes, err = h.formatTextAttrs(attrs)
	}

	if err != nil {
		return fmt.Errorf("error when formatting Term(%v)attrs: %w", h.isTerm, err)
	}

	return h.writeLog(r.Time, levelStr, r.Message, attrsBytes)
}

func (h *PrettyHandler) formatTextAttrs(attrs map[string]any) ([]byte, error) {
	var buf bytes.Buffer

	for k, v := range attrs {
		buf.WriteString(fmt.Sprintf("%s=%v ", k, v))
	}

	return buf.Bytes(), nil
}

func (h *PrettyHandler) formatTermAttrs(attrs map[string]any) ([]byte, error) {
	attrsBytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error when marshaling attrs: %w", err)
	}
	return attrsBytes, nil
}

// formatLevel will colorize the level string
func (h *PrettyHandler) formatLevel(level slog.Level) string {
	levelStr := level.String()

	switch level {
	case slog.LevelDebug:
		return colorize(darkGray, levelStr, h.isTerm)
	case slog.LevelInfo:
		return colorize(cyan, levelStr, h.isTerm)
	case slog.LevelWarn:
		return colorize(lightYellow, levelStr, h.isTerm)
	case slog.LevelError:
		return colorize(lightRed, levelStr, h.isTerm)
	default:
		return levelStr
	}
}

func (h *PrettyHandler) formatAttrs(r slog.Record) (map[string]any, error) {
	attrs := make(map[string]any)

	h.addPredefAttrs(attrs)
	h.addRecordAttrs(attrs, r)

	return attrs, nil
}

func (h *PrettyHandler) addPredefAttrs(attrs map[string]any) {
	for _, attr := range h.attrs {
		if attr.Equal(slog.Attr{}) {
			continue
		}
		attrs[attr.Key] = attr.Value.Any()
	}
}

// addRecordAttrs adds record attributes
func (h *PrettyHandler) addRecordAttrs(attrs map[string]any, r slog.Record) {
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Equal(slog.Attr{}) {
			return true
		}

		if len(h.groups) > 0 {
			h.addAttrToGroup(attrs, attr)
		} else {
			attrs[attr.Key] = attr.Value.Any()
		}

		return true
	})
}

// addAttrToGroup adds attribute to the appropriate group
func (h *PrettyHandler) addAttrToGroup(attrs map[string]any, attr slog.Attr) {
	current := attrs

	for i, g := range h.groups {
		if i == len(h.groups)-1 {
			h.addAttrToFinalGroup(current, g, attr)
		} else {
			current = h.ensureNestedGroup(current, g)
		}
	}
}

// addAttrToFinalGroup adds attribute to the final group
func (h *PrettyHandler) addAttrToFinalGroup(current map[string]any, group string, attr slog.Attr) {
	if grp, ok := current[group].(map[string]any); ok {
		grp[attr.Key] = attr.Value.Any()
	} else {
		current[group] = map[string]any{attr.Key: attr.Value.Any()}
	}
}

// ensureNestedGroup ensures the nested group exists and returns it
func (h *PrettyHandler) ensureNestedGroup(current map[string]any, group string) map[string]any {
	if grp, ok := current[group].(map[string]any); ok {
		return grp
	}

	newGrp := make(map[string]any)
	current[group] = newGrp
	return newGrp
}

// // addSourceLocation adds source code location information
// func (h *PrettyHandler) addSourceLocation(attrs map[string]any, r slog.Record) {
// 	if h.opts.AddSource && r.PC != 0 {
// 		fs := runtime.CallersFrames([]uintptr{r.PC})
// 		frame, _ := fs.Next()
// 		source := map[string]any{
// 			"function": frame.Function,
// 			"file":     frame.File,
// 			"line":     frame.Line,
// 		}
// 		attrs["source"] = source
// 	}
// }

// writeLog will write the formatted log to the output
func (h *PrettyHandler) writeLog(t time.Time, level, msg string, attrsBytes []byte) error {
	h.m.Lock()
	defer h.m.Unlock()

	var buf strings.Builder

	buf.WriteString(colorize(lightGray, t.Format(timeFormat), h.isTerm))
	buf.WriteString(" ")

	buf.WriteString(level)

	if len(level) <= levelWidth {
		buf.WriteString(strings.Repeat(" ", levelWidth-len(level)))
	}

	buf.WriteString(colorize(white, msg, h.isTerm))

	if attrsBytes != nil {
		buf.WriteString(" ")
		buf.WriteString(colorize(darkGray, string(attrsBytes), h.isTerm))
	}

	_, err := fmt.Fprintln(h.w, buf.String())
	return err
}
