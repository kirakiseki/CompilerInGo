package lexer

import (
	"CompilerInGo/utils"
	"errors"
	"github.com/kpango/glg"
	"io"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	reader utils.Reader   // 读取器
	Pos    utils.Position // 当前位置
	File   []byte         // 文件内容
}

var Lex *Lexer      // 全局Lex变量
var Pool *TokenPool // 全局TokenPool变量

// NewLexer 创建一个新的词法分析器
func NewLexer(file string) *Lexer {
	lexer := &Lexer{
		Pos:  utils.Position{Row: 1},             // 设置初始位置
		File: utils.MustValue(os.ReadFile(file)), // 读取文件并检查读取状态
	}

	lexer.reader = strings.NewReader(string(lexer.File)) // 设置读取器
	lexer.NextRune()                                     // 读取第一个字符

	Lex = lexer // 设置全局Lex变量
	return lexer
}

// NextRune 读取下一个字符
func (l *Lexer) NextRune() rune {
	ch, size, err := l.reader.ReadRune()

	if errors.Is(err, io.EOF) { // EOF 文件末尾
		l.Pos.Ch = 0
		return 0
	} else if err != nil { // 读错误
		glg.Fatalln(err)
	}

	// 更新位置信息
	l.Pos.Ch = ch
	l.Pos.FilePos += uint(size)

	// 更新行列信息
	if ch == '\n' {
		l.Pos.Row++
		l.Pos.Col = 0
	} else {
		l.Pos.Col++
	}

	return ch
}

// peek 查看下一个字符
func (l *Lexer) peek() rune {
	ch, _, err := l.reader.ReadRune() // 读取下一个字符
	if errors.Is(err, io.EOF) {       // EOF 文件末尾
		return 0
	} else if err != nil { // 读错误
		glg.Fatalln(err)
	}

	l.unread() // 回退一个字符

	return ch
}

// unread 回退一个字符
func (l *Lexer) unread() {
	// 使用Seek从当前位置回退一个字符，且MustValue检查错误
	utils.MustValue(l.reader.Seek(-1, io.SeekCurrent))
}

// skipBlank 跳过空白字符
func (l *Lexer) skipBlank() {
	for {
		ch := l.Pos.Ch
		if ch != '\t' && ch != '\r' && ch != '\n' {
			break
		}

		// 读取下一个字符
		l.NextRune()
	}
}

// scanString 扫描字符串
func (l *Lexer) scanString() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos} // 记录开始位置

	var str strings.Builder // 使用strings.Builder拼接字符串

	for {
		ch := l.peek() // 查看下一个字符
		if ch == '"' { // 字符串结束
			tokenPos.End = l.Pos // 更新结束位置
			break
		} else if ch == 0 { // EOF
			// 未结束字符串
			return Token{}, errors.New("unterminated string literal")
		}

		// 字符串未结束 继续读取
		l.NextRune()
		str.WriteRune(l.Pos.Ch)

		// 更新结束位置
		tokenPos.End = l.Pos
	}

	// 读取结束双引号
	l.NextRune()
	// 记录结束位置
	tokenPos.End = l.Pos
	// 读取指针后移
	l.NextRune()

	return NewToken(str.String(), tokenPos, STRING_LITERAL), nil
}

// scanChar 扫描字符
func (l *Lexer) scanChar() (Token, error) {
	// 记录开始位置
	tokenPos := utils.PositionPair{Begin: l.Pos}

	// 读取下一个字符
	ch := l.NextRune()
	if ch == 0 {
		// 未结束字符
		return Token{}, errors.New("unterminated char literal")
	}
	if ch == '\'' {
		// 空字符，读取指针后移
		l.NextRune()
		return NewToken("", tokenPos, CHAR_LITERAL), nil
	}

	if endCh := l.NextRune(); endCh != '\'' {
		// 判断是否单字符
		if endCh == 0 || endCh == '\n' {
			// 未结束字符
			return Token{}, errors.New("unterminated char literal")
		}
		// 非单字符，非法
		return Token{}, errors.New("invalid char literal")
	}

	// 记录结束位置
	tokenPos.End = l.Pos
	// 读取指针后移
	l.NextRune()

	return NewToken(ch, tokenPos, CHAR_LITERAL), nil
}

// scanIdentifier 扫描标识符
func (l *Lexer) scanIdentifier() (Token, error) {
	// 记录开始位置
	tokenPos := utils.PositionPair{Begin: l.Pos}

	// 使用strings.Builder拼接字符串
	var str strings.Builder

	// 读取标识符 由$ 字母 数字组成 开头在ScanToken中已判断
	for l.Pos.Ch == '$' || unicode.IsLetter(l.Pos.Ch) || unicode.IsDigit(l.Pos.Ch) {
		// 写入字符
		str.WriteRune(l.Pos.Ch)
		// 更新结束位置
		tokenPos.End = l.Pos
		// 读取下一个字符
		l.NextRune()
	}

	return NewToken(str.String(), tokenPos, IDENTIFIER), nil
}

// scanNumber 扫描数字
func (l *Lexer) scanNumber() (Token, error) {
	// 记录开始位置
	tokenPos := utils.PositionPair{Begin: l.Pos}

	// 使用strings.Builder拼接字符串
	var str strings.Builder
	// 记录小数点数量
	var pointCount int

	if l.Pos.Ch == '-' {
		// 如果为负数，写入字符
		str.WriteRune(l.Pos.Ch)
		// 更新结束位置
		tokenPos.End = l.Pos
		// 读取下一个字符
		l.NextRune()
	}

	// 如果为数字或小数点，继续读取
	for unicode.IsDigit(l.Pos.Ch) || l.Pos.Ch == '.' {
		// 如果为小数点，记录数量
		if l.Pos.Ch == '.' {
			pointCount++
			if pointCount > 1 {
				// 如果小数点数量大于1，非法
				return IfTokenError(Token{}, errors.New("invalid number literal")), nil
			}
		}

		// 写入字符
		str.WriteRune(l.Pos.Ch)
		// 更新结束位置
		tokenPos.End = l.Pos
		// 读取下一个字符
		l.NextRune()
	}

	if pointCount == 0 {
		// 如果为整数
		num := utils.MustValue(utils.ParseInt(str.String()))
		return NewToken(num, tokenPos, INTEGER_LITERAL), nil
	} else {
		// 如果为小数
		num := utils.MustValue(utils.ParseFloat(str.String()))
		return NewToken(num, tokenPos, DECIMAL_LITERAL), nil
	}
}

// scanComment 扫描注释
func (l *Lexer) scanComment() (Token, error) {
	// 记录开始位置
	tokenPos := utils.PositionPair{Begin: l.Pos}

	// 使用strings.Builder拼接字符串
	var str strings.Builder
	// 注释字面量
	var result string
	// 注释类型(单行SINGLELINE_COMMENT_LITERAL，多行MULTILINE_COMMENT_LITERAL)
	var commentType TokenType

	// 读取下一个字符
	switch l.peek() {
	case '/':
		// 单行注释
		l.NextRune()
		commentType = SINGLELINE_COMMENT_LITERAL

		for {
			ch := l.NextRune()
			// 扫描到换行或EOF结束
			if ch == 0 || ch == '\n' {
				break
			}

			str.WriteRune(ch)
			// 更新结束位置
			tokenPos.End = l.Pos
		}

		result = str.String()
	case '*':
		// 多行注释
		l.NextRune()
		commentType = MULTILINE_COMMENT_LITERAL

		for {
			ch := l.NextRune()
			// 未配对的注释
			if ch == 0 {
				return Token{}, errors.New("unterminated comment")
			}

			str.WriteRune(ch)
			// 更新结束位置
			tokenPos.End = l.Pos

			// 与结束的*/配对
			if ch == '*' && l.peek() == '/' {
				// 读取指针后移到/
				l.NextRune()
				// 更新结束位置
				tokenPos.End = l.Pos
				break
			}
		}
		result = str.String()[0 : str.Len()-1] // 去掉最后一个 *
	}

	// 读取指针右移
	l.NextRune()

	return NewToken(result, tokenPos, commentType), nil
}

// scanOperator 扫描操作符
func (l *Lexer) scanOperator() (Token, error) {
	// 记录开始位置
	tokenPos := utils.PositionPair{Begin: l.Pos}

	// 使用strings.Builder拼接字符串
	var str strings.Builder
	// 写入当前字符
	str.WriteRune(l.Pos.Ch)
	// 更新结束位置
	tokenPos.End = l.Pos

	// 判断是否为双字符操作符
	if ch := l.peek(); (str.String() == "=" && ch == '=') || (str.String() == "<" && ch == '>') || (str.String() == "<" && ch == '=') || (str.String() == ">" && ch == '=') {
		// 如果是==/<>/<=/>=，写入下一个字符
		str.WriteRune(ch)
		// 更新结束位置
		tokenPos.End = l.Pos
		// 读取指针右移到双字符操作符的第二个字符
		l.NextRune()
	}

	// 读取指针后移
	l.NextRune()

	// 判断操作符类型
	switch str.String() {
	case "==":
		return NewToken("==", tokenPos, EQUAL), nil
	case "=":
		return NewToken("=", tokenPos, ASSIGN), nil
	case "<":
		return NewToken("<", tokenPos, LESS), nil
	case "<=":
		return NewToken("<=", tokenPos, LESSEQUAL), nil
	case ">":
		return NewToken(">", tokenPos, GREATER), nil
	case ">=":
		return NewToken(">=", tokenPos, GREATEREQUAL), nil
	case "<>":
		return NewToken("<>", tokenPos, DIAMOND), nil
	case "+":
		return NewToken("+", tokenPos, PLUS), nil
	case "-":
		return NewToken("-", tokenPos, MINUS), nil
	case "*":
		return NewToken("*", tokenPos, TIMES), nil
	case "/":
		return NewToken("/", tokenPos, DIVIDE), nil
	}

	// 未匹配到操作符
	return Token{}, nil
}

// scanDelim 扫描分隔符
func (l *Lexer) scanDelim() (Token, error) {
	// 记录开始位置
	tokenPos := utils.PositionPair{
		Begin: l.Pos,
		End:   l.Pos,
	}

	ch := l.Pos.Ch

	// 读取指针后移
	l.NextRune()

	// 判断分隔符类型
	switch ch {
	case '(':
		return NewToken("(", tokenPos, LPAREN), nil
	case ')':
		return NewToken(")", tokenPos, RPAREN), nil
	case '{':
		return NewToken("{", tokenPos, LBRACE), nil
	case '}':
		return NewToken("}", tokenPos, RBRACE), nil
	case ';':
		return NewToken(";", tokenPos, SEMICOLON), nil
	case ' ':
		return NewToken(" ", tokenPos, SPACE), nil
	}

	// 未匹配到分隔符
	return Token{}, nil
}

func (l *Lexer) ScanToken() (Token, error) {
	// 跳过 '\t' '\r' '\n'
	l.skipBlank()

	if l.Pos.Ch == 0 {
		// 扫描到EOF
		_ = glg.Debug("Scan Completed")
		return NewToken("EOF_LITERAL", utils.PositionPair{Begin: l.Pos, End: l.Pos}, EOF_LITERAL), nil
	} else if IsDelim(string(l.Pos.Ch)) {
		// 扫描分隔符
		return l.scanDelim()
	} else if l.Pos.Ch == '"' {
		// 扫描字符串
		return l.scanString()
	} else if l.Pos.Ch == '\'' {
		// 扫描字符
		return l.scanChar()
	} else if (l.Pos.Ch == '/' && l.peek() == '/') || (l.Pos.Ch == '/' && l.peek() == '*') {
		// 扫描注释
		return l.scanComment()
	} else if unicode.IsLetter(l.Pos.Ch) || l.Pos.Ch == '$' {
		// 扫描标识符或关键字
		if l.Pos.Ch == '$' {
			// $开头 只能为标识符
			return l.scanIdentifier()
		}
		// 按照标识符扫描
		token := IfTokenError(l.scanIdentifier())
		// 判断标识符内容是否为关键字
		if IsKeyword(token.Literal.(string)) {
			// 是关键字
			// 寻找关键字对应的类型
			for k, v := range tokenType {
				if v == token.Literal.(string) {
					token.Type = k
				}
			}
			// 设置类别为关键字
			token.Category = KEYWORD
		}
		return token, nil
	} else if unicode.IsDigit(l.Pos.Ch) || (l.Pos.Ch == '-' && unicode.IsDigit(l.peek())) {
		// 扫描数字
		// 以数字开头或者以负号开头后面跟着数字
		return l.scanNumber()
	} else if IsOpera(string(l.Pos.Ch)) {
		// 扫描操作符
		return l.scanOperator()
	} else {
		// 无法识别的字符
		return IfTokenError(Token{}, errors.New("unrecognized token")), nil
	}
}
