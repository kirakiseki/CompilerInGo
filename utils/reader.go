package utils

import "io"

// Reader 读取器接口
type Reader interface {
	io.RuneReader
	io.Seeker
}
