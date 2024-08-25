package grpcuipuzzle

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/soheilhy/cmux"
)

// httpPrefixMatcher returns a cmux.MatcherFunc that matches based on HTTP path prefix.
func httpPrefixMatcher(prefix string) cmux.Matcher {
	return func(r io.Reader) bool {
		var buf bytes.Buffer
		tee := io.TeeReader(r, &buf)

		req, err := http.ReadRequest(bufio.NewReader(tee))
		if err != nil {
			return false
		}
		// Check if the request URL path has the specified prefix.
		return strings.HasPrefix(req.URL.Path, prefix)
	}
}
