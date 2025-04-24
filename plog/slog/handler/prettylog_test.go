// File:		prettylog_test.go
// Created by:	Hoven
// Created on:	2025-04-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"log/slog"
	"os"
	"testing"
)

func TestPrettyHandler(t *testing.T) {
	handler := NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)
	logger.WithGroup("Group").WithGroup("SubGroup").Info("this is info msg with group", "name", "hoven")
	logger.Info("this is info msg", "name", "hoven")
}
