package utils

import (
	"github.com/kpango/glg"
	"runtime/debug"
)

func MustValue[T any](value T, err error) T {
	if err != nil {
		// 检查是否出错
		// 输出Stacktrace
		_ = glg.Fail("MustValue error detected! Stacktrace:", string(debug.Stack()))
		glg.Fatal(err)
	}

	// 没有出错则返回值
	return value
}
