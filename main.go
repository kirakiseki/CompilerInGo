package main

import (
	"CompilerInGo/lexer"
	"CompilerInGo/utils"
)

func main() {
	utils.InitLogger("DEBUG")
	lex := lexer.NewLexer("./test/sample.program")
	lex.TraverseRune()
}
