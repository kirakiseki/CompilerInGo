package lexer

import "io"

type Reader interface {
	io.RuneReader
	io.Seeker
}
