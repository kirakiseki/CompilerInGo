package lexer

import (
	"CompilerInGo/utils"
	"github.com/kpango/glg"
	"strings"
)

func IfTokenError(token Token, err error) Token {
	// Token解析是否出错
	if err != nil {
		_ = glg.Fail("Error while scanning Token: ", err)

		// 获取lexer中状态
		lex := Lex

		// 获取文件出错行内容
		errorLine := utils.GetLine(lex.File, lex.Pos)

		// 显示错误信息
		_ = glg.Failf("Position: Line %d, Column %d", lex.Pos.Row, lex.Pos.Col)
		_ = glg.Fail(errorLine)

		// 构造错误位置指示器
		var str strings.Builder
		for i := 0; i < int(lex.Pos.Col-1); i++ {
			str.WriteRune('-')
		}
		str.WriteRune('^')

		// 显示错误位置指示器
		_ = glg.Fail(str.String())

		// 退出程序
		glg.Fatal("Error while scanning Token")
	}
	return token
}
