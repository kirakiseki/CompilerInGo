package lexer

import (
	"CompilerInGo/utils"
	"testing"
)

func BenchmarkEscapeRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.EscapeRune('\'')
	}
}

func BenchmarkEscapeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.EscapeString("a\n")
	}
}
