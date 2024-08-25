package plog

import (
	"testing"
	
	logctx "github.com/go-puzzles/puzzles/plog/log-ctx"
	"github.com/stretchr/testify/assert"
)

func TestParseFmtKeyValue(t *testing.T) {
	lc := &logctx.LogContext{
		Keys:   make([]string, 0),
		Values: make([]string, 0),
	}
	nlc, err := parseFmtKeyValue(lc, "group")
	assert.Nil(t, err)
	
	t.Logf("newlc: %v", nlc)
	t.Log("==============================")
	
	lc = &logctx.LogContext{
		Keys:   make([]string, 0),
		Values: make([]string, 0),
	}
	nlc, err = parseFmtKeyValue(lc, "key1", "value1")
	assert.Nil(t, err)
	
	t.Logf("newlc: %v", nlc)
	t.Log("==============================")
	
	lc = &logctx.LogContext{
		Keys:   make([]string, 0),
		Values: make([]string, 0),
	}
	nlc, err = parseFmtKeyValue(lc, "key1", "value1", "key2", "value2")
	assert.Nil(t, err)
	
	t.Logf("newlc: %v", nlc)
	t.Log("==============================")
	
	lc = &logctx.LogContext{
		Keys:   make([]string, 0),
		Values: make([]string, 0),
	}
	nlc, err = parseFmtKeyValue(lc, "key1", "value1", "key2")
	assert.Nil(t, err)
	
	t.Logf("newlc: %v", nlc)
	t.Log("==============================")
	
}
