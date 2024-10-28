// File:		matcher.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package cores

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/soheilhy/cmux"
)

// HttpPrefixMatcher returns a cmux.MatcherFunc that matches based on HTTP path prefix.
func HttpPrefixMatcher(prefix string) cmux.Matcher {
	return func(r io.Reader) bool {
		var buf bytes.Buffer
		tee := io.TeeReader(r, &buf)

		req, err := http.ReadRequest(bufio.NewReader(tee))
		if err != nil {
			plog.Errorf("failed to read http request: %v", err)
			return false
		}
		// Check if the request URL path has the specified prefix.
		return strings.HasPrefix(req.URL.Path, prefix)
	}
}
