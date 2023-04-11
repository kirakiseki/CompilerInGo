package lexer

import (
	"CompilerInGo/utils"
	"errors"
	"github.com/kpango/glg"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	reader utils.Reader
	Pos    utils.Position
	File   []byte
}

var Lex *Lexer

func NewLexer(file string) *Lexer {
	lexer := &Lexer{
		Pos:  utils.Position{Row: 1},             // 设置初始位置
		File: utils.MustValue(os.ReadFile(file)), // 读取文件并检查读取状态
	}

	lexer.reader = strings.NewReader(string(lexer.File)) // 设置读取器
	lexer.nextRune()                                     // 读取第一个字符

	Lex = lexer
	return lexer
}

func (l *Lexer) nextRune() rune {
	// 读取下一个字符
	ch, size, err := l.reader.ReadRune()

	if errors.Is(err, io.EOF) { // EOF 文件末尾
		l.Pos.Ch = 0
		return 0
	} else if err != nil { // 读错误
		glg.Fatalln(err)
	}

	l.Pos.Ch = ch
	l.Pos.FilePos += uint(size)

	if ch == '\n' {
		l.Pos.Row++
		l.Pos.Col = 0
	} else {
		l.Pos.Col++
	}

	return ch
}

func (l *Lexer) peek() rune {
	ch, _, err := l.reader.ReadRune()
	if errors.Is(err, io.EOF) { // EOF 文件末尾
		return 0
	} else if err != nil { // 读错误
		glg.Fatalln(err)
	}
	l.unread()
	return ch
}

func (l *Lexer) unread() {
	utils.MustValue(l.reader.Seek(-1, io.SeekCurrent))
}

func (l *Lexer) skipBlank() {
	for {
		ch := l.Pos.Ch
		if ch != '\t' && ch != '\r' && ch != '\n' {
			break
		}
		l.nextRune()
	}
}

func (l *Lexer) TraverseRune() {
	for {
		l.skipBlank()
		_ = glg.Debugf("%c %+v", l.Pos.Ch, l.Pos)
		ch := l.nextRune()
		if ch == 0 {
			break
		}
	}
}

func (l *Lexer) scanString() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	var str strings.Builder
	for {
		ch := l.peek()
		if ch == '"' {
			str.WriteRune(ch)
			tokenPos.End = l.Pos
			break
		} else if ch == 0 { // EOF
			return Token{}, errors.New("unterminated string literal")
		}
		l.nextRune()
		str.WriteRune(l.Pos.Ch)
		tokenPos.End = l.Pos
	}
	l.nextRune()
	l.nextRune()
	return NewToken(str.String()[0:str.Len()-1], tokenPos, STRING_LITERAL), nil // 去掉最后一个双引号
}

func (l *Lexer) scanChar() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	ch := l.nextRune()
	if ch == 0 {
		return Token{}, errors.New("unterminated char literal")
	}
	if ch == '\'' {
		l.nextRune()
		return NewToken("", tokenPos, CHAR_LITERAL), nil
	}
	if endCh := l.nextRune(); endCh != '\'' {
		if endCh == 0 || endCh == '\n' {
			return Token{}, errors.New("unterminated char literal")
		}
		return Token{}, errors.New("invalid char literal")
	}
	tokenPos.End = l.Pos
	l.nextRune()

	return NewToken(ch, tokenPos, CHAR_LITERAL), nil
}

func (l *Lexer) scanIdentifier() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	var str strings.Builder
	for l.Pos.Ch == '$' || unicode.IsLetter(l.Pos.Ch) || unicode.IsDigit(l.Pos.Ch) {
		str.WriteRune(l.Pos.Ch)
		tokenPos.End = l.Pos
		l.nextRune()
	}

	return NewToken(str.String(), tokenPos, IDENTIFIER), nil
}

func (l *Lexer) scanNumber() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	var str strings.Builder
	var pointCount int
	if l.Pos.Ch == '-' {
		str.WriteRune(l.Pos.Ch)
		tokenPos.End = l.Pos
		l.nextRune()
	}
	for unicode.IsDigit(l.Pos.Ch) || l.Pos.Ch == '.' {
		if l.Pos.Ch == '.' {
			pointCount++
			if pointCount > 1 {
				return IfTokenError(Token{}, errors.New("invalid number literal")), nil
			}
		}
		str.WriteRune(l.Pos.Ch)
		tokenPos.End = l.Pos
		l.nextRune()
	}

	if pointCount == 0 {
		num := utils.MustValue(strconv.ParseInt(str.String(), 10, 64))
		return NewToken(num, tokenPos, INTEGER_LITERAL), nil
	} else {
		num := utils.MustValue(strconv.ParseFloat(str.String(), 64))
		return NewToken(num, tokenPos, DECIMAL_LITERAL), nil
	}
}

func (l *Lexer) scanComment() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	var str strings.Builder
	var result string
	var commentType TokenType
	switch l.peek() {
	case '/':
		l.nextRune()
		commentType = SINGLELINE_COMMENT_LITERAL
		for {
			ch := l.nextRune()
			if ch == 0 || ch == '\n' {
				break
			}
			str.WriteRune(ch)
			tokenPos.End = l.Pos
		}
		result = str.String()
	case '*':
		l.nextRune()
		commentType = MULTILINE_COMMENT_LITERAL
		for {
			ch := l.nextRune()
			if ch == 0 {
				return Token{}, errors.New("unterminated comment")
			}
			str.WriteRune(ch)
			tokenPos.End = l.Pos
			if ch == '*' && l.peek() == '/' {
				l.nextRune()
				tokenPos.End = l.Pos
				break
			}
		}
		result = str.String()[0 : str.Len()-1] // 去掉最后一个 *
	}
	l.nextRune()
	return NewToken(result, tokenPos, commentType), nil
}

func (l *Lexer) scanOperator() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	var str strings.Builder
	str.WriteRune(l.Pos.Ch)
	tokenPos.End = l.Pos
	if ch := l.peek(); (str.String() == "=" && ch == '=') || (str.String() == "<" && ch == '>') {
		str.WriteRune(ch)
		tokenPos.End = l.Pos
		l.nextRune()
	}
	l.nextRune()

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
	return Token{}, nil
}

func (l *Lexer) scanDelim() (Token, error) {
	tokenPos := utils.PositionPair{Begin: l.Pos}

	var str strings.Builder
	str.WriteRune(l.Pos.Ch)
	tokenPos.End = l.Pos
	l.nextRune()

	switch str.String() {
	case "(":
		return NewToken("(", tokenPos, LPAREN), nil
	case ")":
		return NewToken(")", tokenPos, RPAREN), nil
	case "{":
		return NewToken("{", tokenPos, LBRACE), nil
	case "}":
		return NewToken("}", tokenPos, RBRACE), nil
	case ";":
		return NewToken(";", tokenPos, SEMICOLON), nil
	case " ":
		return NewToken(" ", tokenPos, SPACE), nil
	}
	return Token{}, nil
}

func (l *Lexer) ScanToken() (Token, error) {
	l.skipBlank()

	if l.Pos.Ch == 0 {
		_ = glg.Debug("Scan Completed")
		return NewToken("EOF_LITERAL", utils.PositionPair{Begin: l.Pos, End: l.Pos}, EOF_LITERAL), nil
	} else if IsDelim(string(l.Pos.Ch)) {
		return l.scanDelim()
	} else if l.Pos.Ch == '"' {
		return l.scanString()
	} else if l.Pos.Ch == '\'' {
		return l.scanChar()
	} else if (l.Pos.Ch == '/' && l.peek() == '/') || (l.Pos.Ch == '/' && l.peek() == '*') {
		return l.scanComment()
	} else if unicode.IsLetter(l.Pos.Ch) || l.Pos.Ch == '$' {
		if l.Pos.Ch == '$' {
			return l.scanIdentifier()
		}
		token := IfTokenError(l.scanIdentifier())
		if IsKeyword(token.Literal.(string)) {
			for k, v := range tokenString {
				if v == token.Literal.(string) {
					token.Type = k
				}
			}
			token.Category = KEYWORD
		}
		return token, nil
	} else if unicode.IsDigit(l.Pos.Ch) || (l.Pos.Ch == '-' && unicode.IsDigit(l.peek())) {
		return l.scanNumber()
	} else if IsOpera(string(l.Pos.Ch)) {
		return l.scanOperator()
	} else {
		return IfTokenError(Token{}, errors.New("unrecognized token")), nil
	}
	return Token{}, nil
}
