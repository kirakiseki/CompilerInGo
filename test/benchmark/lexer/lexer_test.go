package lexer

import (
	"CompilerInGo/lexer"
	"testing"
)

func BenchmarkLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lex := lexer.NewLexer("../../long.program")
		// 初始化token池
		tokenPool := lexer.NewTokenPool()

		// 读取第一个Token
		// IfTokenError 检查Token是否出错，若出错则输出错误信息并退出程序
		token := lexer.IfTokenError(lex.ScanToken())
		tokenPool.Add(token)
		// 若未读到EOF则继续读取
		for tokenPool.Last().Category != lexer.EOF {
			token := lexer.IfTokenError(lex.ScanToken())
			tokenPool.Add(token)
		}
	}
}
