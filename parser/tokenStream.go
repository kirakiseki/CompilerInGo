package parser

import (
	"CompilerInGo/lexer"
	"CompilerInGo/utils"
	"github.com/kpango/glg"
)

type TokenStream struct {
	pool  lexer.TokenPool
	pos   int
	width int
}

func NewTokenStream() TokenStream {
	return TokenStream{
		pool:  *lexer.Pool,
		pos:   0,
		width: 0,
	}
}

func (ts *TokenStream) ReadToken() lexer.Token {
	if ts.pos >= len(ts.pool.Pool) {
		ts.width = 0
		return lexer.NewToken("EOF_LITERAL", utils.PositionPair{}, lexer.EOF_LITERAL)
	}
	token := ts.pool.Pool[ts.pos]
	for token.Type == lexer.SPACE || token.Type == lexer.SINGLELINE_COMMENT_LITERAL || token.Type == lexer.MULTILINE_COMMENT_LITERAL {
		ts.pos++
		token = ts.pool.Pool[ts.pos]
	}
	ts.width = 1
	ts.pos++

	return token
}

func (ts *TokenStream) UnreadToken() {
	ts.pos -= ts.width
}

func (ts *TokenStream) PeekToken() lexer.Token {
	token := ts.ReadToken()
	ts.UnreadToken()
	return token
}

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

func (ts *TokenStream) AcceptTokenByFunc(f func(token lexer.Token) bool) (lexer.Token, error) {
	token := ts.ReadToken()
	if f(token) {
		return token, nil
	}
	ts.UnreadToken()
	return token, utils.NewErrorf("Token %v does not meet the expectation", token)
}

func (ts *TokenStream) MustAcceptTokenByType(exceptedType ...lexer.TokenType) lexer.Token {
	token, err := ts.AcceptTokenByType(exceptedType...)
	if err != nil {
		glg.Fatalf("Expect token type %v, but got %#v", exceptedType, token)
	}
	return token
}

func (ts *TokenStream) MustAcceptTokenByFunc(f func(token lexer.Token) bool) lexer.Token {
	token, err := ts.AcceptTokenByFunc(f)
	if err != nil {
		glg.Fatalf("Token %v does not meet the expectation", token)
	}
	return token
}

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

func (ts *TokenStream) OptionalAcceptTokenByFunc(f func(token lexer.Token) bool) (lexer.Token, bool) {
	token := ts.ReadToken()
	if f(token) {
		return token, true
	}
	ts.UnreadToken()
	return token, false
}
