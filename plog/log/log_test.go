package log

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	l := New()

	l.Infoc(context.Background(), "this is a log")
}
