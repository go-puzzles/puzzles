package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	badKey = "<BadKey>"
)

func ParseFmtStr(format string) (msg string, isKV []bool, keys, descs []string) {
	if json.Valid([]byte(format)) {
		return format, nil, nil, nil
	}

	var msgs []string
	for _, s := range strings.Split(format, " ") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		idx := strings.Index(s, "=%")
		if idx == -1 || strings.Contains(s[:idx], "=") {
			re, _ := regexp.Compile("%[^%]+")
			matches := re.FindAllStringIndex(s, -1)
			for i := 0; i < len(matches); i++ {
				isKV = append(isKV, false)
			}
			msgs = append(msgs, s)
			continue
		}
		keys = append(keys, s[:idx])
		descs = append(descs, s[idx+1:])
		isKV = append(isKV, true)
	}
	msg = strings.Join(msgs, " ")
	return
}

func ParseFmtKeyValue(msg string, v ...any) (m string, keys, values []string, err error) {
	/*
		msg= hello a=%s world %d        v = &a, 1
		hello world %d
		[true, false]
		[a]
		[%s]

		it will parse the message like a=%s out of the message and put it to the end of msg
		hello world 1 a=%s
	*/
	msgTmpl, isKV, keys, desc := ParseFmtStr(msg)
	var msgV []any
	var objV []any
	for i, kv := range isKV {
		var val any
		if i >= len(v) {
			val = badKey
		} else {
			val = v[i]
		}
		if kv {
			objV = append(objV, val)
		} else {
			msgV = append(msgV, val)
		}
	}

	msg = fmt.Sprintf(msgTmpl, msgV...)
	if len(objV) != len(desc) {
		return "", nil, nil, errors.New("invalid numbers of keys and values")
	}

	for i := range desc {
		values = append(values, fmt.Sprintf(desc[i], objV[i]))
	}

	if len(isKV) < len(v) {
		keys, values = parseRemains(v[len(isKV):], keys, values)
	}

	return msg, keys, values, nil
}

func parseRemains(remains []any, keys, vals []string) ([]string, []string) {
	var key, val string
	for len(remains) > 0 {
		key, val, remains = getKVParis(remains)
		keys = append(keys, key)
		vals = append(vals, val)
	}

	return keys, vals
}

func getKVParis(kvs []any) (string, string, []any) {
	switch x := kvs[0].(type) {
	case string:
		if len(kvs) == 1 {
			return badKey, x, nil
		}
		return x, fmt.Sprintf("%v", kvs[1]), kvs[2:]
	default:
		return badKey, fmt.Sprintf("%v", x), kvs[1:]
	}
}
