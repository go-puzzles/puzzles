package predis

import "errors"

var (
	ErrLockFailed = errors.New("redis lock failed")
	ErrDuplicated = errors.New("task duplicated")
)
