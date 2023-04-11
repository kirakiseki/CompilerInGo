package main

import (
	"CompilerInGo/lexer"
	"CompilerInGo/utils"
	"github.com/kpango/glg"
)

func main() {
	utils.InitLogger("DEBUG")

	lex := lexer.NewLexer("./test/sample.program")
	tokenPool := lexer.NewTokenPool()

	token := lexer.IfTokenError(lex.ScanToken())
	tokenPool.Add(token)
	for tokenPool.Last().Category != lexer.EOF {
		token := lexer.IfTokenError(lex.ScanToken())
		tokenPool.Add(token)
	}

	for _, token := range tokenPool.Pool {
		_ = glg.Info(token.String())
		//_ = glg.Debugf("%#v", token)
	}
}
