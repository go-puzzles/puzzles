package plog

import (
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func GetStructName(i any) string {
	tPtr := reflect.TypeOf(i)

	if tPtr.Kind() == reflect.Ptr {
		tPtr = tPtr.Elem()
	}

	return tPtr.Name()
}

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// TimeFuncDuration returns the duration consumed by function.
// It has specified usage like:
//
//	    f := TimeFuncDuration()
//		   DoSomething()
//		   duration := f()
func TimeFuncDuration() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}

func TimeDurationDefer(prefix ...string) func() {
	ps := "operation"
	if len(prefix) != 0 {
		ps = strings.Join(prefix, ", ")
	}
	start := time.Now()

	return func() {
		Infof("%v elapsed time: %v", ps, time.Since(start))
	}
}

func Jsonify(v any) string {
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Errorf("jsonify error: %v", err)
		panic(err)
	}
	return string(d)
}

func JsonifyNoIndent(v interface{}) string {
	d, err := json.Marshal(v)
	PanicError(err)
	return string(d)
}
