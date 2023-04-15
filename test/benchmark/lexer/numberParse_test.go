package lexer

import (
	"CompilerInGo/utils"
	"strconv"
	"testing"
)

func BenchmarkParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = utils.ParseInt("123")
	}
}

func BenchmarkParseFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = utils.ParseFloat("123.456")
	}
}

func BenchmarkParseIntBuiltIn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = strconv.ParseInt("123", 10, 64)
	}
}

func BenchmarkParseFloatBuiltIn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = strconv.ParseFloat("123.456", 64)
	}
}
