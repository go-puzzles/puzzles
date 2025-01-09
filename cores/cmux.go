// File:		cmux.go
// Created by:	Hoven
// Created on:	2025-01-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package cores

import (
	"io"

	"github.com/soheilhy/cmux"
)

func HTTP1FastMatchWriter() cmux.MatchWriter {
	return func(_ io.Writer, r io.Reader) bool {
		matchHttp1 := cmux.HTTP1Fast()(r)
		return matchHttp1
	}
}

func HTTP2MatchWriter() cmux.MatchWriter {
	return func(_ io.Writer, r io.Reader) bool {
		matchHttp2 := cmux.HTTP2()(r)
		return matchHttp2
	}
}

func HTTP2MatchWithHeaderExclude(name, value string) cmux.MatchWriter {
	return func(w io.Writer, r io.Reader) bool {
		// if match http2 and match header exclude, return false
		matchExclude := cmux.HTTP2MatchHeaderFieldSendSettings(name, value)(w, r)
		return !matchExclude
	}
}
