package utils

import (
	"github.com/kpango/glg"
	"runtime/debug"
)

func MustValue[T any](value T, err error) T {
	if err != nil {
		_ = glg.Fail("MustValue errorDisplay detected! Stacktrace:", string(debug.Stack()))
		glg.Fatal(err)
	}
	return value
}
