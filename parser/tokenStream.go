package parser

import (
	"CompilerInGo/lexer"
	"CompilerInGo/utils"
	"github.com/kpango/glg"
)

// TokenStream Token流
type TokenStream struct {
	pool  lexer.TokenPool
	pos   int
	width int
}

// NewTokenStream 创建一个新的Token流
func NewTokenStream() TokenStream {
	return TokenStream{
		pool:  *lexer.Pool,
		pos:   0,
		width: 0,
	}
}

// ReadToken 读取一个Token
func (ts *TokenStream) ReadToken() lexer.Token {
	// 越界，返回EOF
	if ts.pos >= len(ts.pool.Pool) {
		ts.width = 0
		return lexer.NewToken("EOF_LITERAL", utils.PositionPair{}, lexer.EOF_LITERAL)
	}
	// 读取Token
	token := ts.pool.Pool[ts.pos]
	// 跳过空格和注释
	for token.Type == lexer.SPACE || token.Type == lexer.SINGLELINE_COMMENT_LITERAL || token.Type == lexer.MULTILINE_COMMENT_LITERAL {
		ts.pos++
		token = ts.pool.Pool[ts.pos]
	}
	// 更新宽度
	ts.width = 1
	// 更新位置
	ts.pos++

	return token
}

// UnreadToken 回退一个Token
func (ts *TokenStream) UnreadToken() {
	ts.pos -= ts.width
}

// PeekToken 预读一个Token
func (ts *TokenStream) PeekToken() lexer.Token {
	token := ts.ReadToken()
	ts.UnreadToken()
	return token
}

// AcceptTokenByType 读取一个指定类型的Token
// 如果读取到的Token不是指定类型，那么回退Token并返回错误
func (ts *TokenStream) AcceptTokenByType(exceptedType ...lexer.TokenType) (lexer.Token, error) {
	token := ts.ReadToken()
	for _, t := range exceptedType {
		if token.Type == t {
			return token, nil
		}
	}
	ts.UnreadToken()
	return token, utils.NewErrorf("Expect token type %v, but got %v", exceptedType, token)
}

// AcceptTokenByFunc 读取一个满足条件的Token
// 如果读取到的Token不满足条件，那么回退Token并返回错误
func (ts *TokenStream) AcceptTokenByFunc(f func(token lexer.Token) bool) (lexer.Token, error) {
	token := ts.ReadToken()
	if f(token) {
		return token, nil
	}
	ts.UnreadToken()
	return token, utils.NewErrorf("Token %v does not meet the expectation", token)
}

// MustAcceptTokenByType 必须满足指定类型的Token
// 如果读取到的Token不是指定类型，那么直接退出程序
func (ts *TokenStream) MustAcceptTokenByType(exceptedType ...lexer.TokenType) lexer.Token {
	token, err := ts.AcceptTokenByType(exceptedType...)
	if err != nil {
		glg.Fatalf("Expect token type %v, but got %#v", exceptedType, token)
	}
	return token
}

// MustAcceptTokenByFunc 必须满足条件的Token
// 如果读取到的Token不满足条件，那么直接退出程序
func (ts *TokenStream) MustAcceptTokenByFunc(f func(token lexer.Token) bool) lexer.Token {
	token, err := ts.AcceptTokenByFunc(f)
	if err != nil {
		glg.Fatalf("Token %v does not meet the expectation", token)
	}
	return token
}

// OptionalAcceptTokenByType 可选满足指定类型的Token
// 如果读取到的Token不是指定类型，那么回退Token并返回false
func (ts *TokenStream) OptionalAcceptTokenByType(exceptedType ...lexer.TokenType) (lexer.Token, bool) {
	token := ts.ReadToken()
	for _, t := range exceptedType {
		if token.Type == t {
			return token, true
		}
	}
	ts.UnreadToken()
	return token, false
}

// OptionalAcceptTokenByFunc 可选满足条件的Token
// 如果读取到的Token不满足条件，那么回退Token并返回false
func (ts *TokenStream) OptionalAcceptTokenByFunc(f func(token lexer.Token) bool) (lexer.Token, bool) {
	token := ts.ReadToken()
	if f(token) {
		return token, true
	}
	ts.UnreadToken()
	return token, false
}
