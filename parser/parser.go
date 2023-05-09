package parser

import (
	"CompilerInGo/lexer"
	"github.com/kpango/glg"
)

type Parser struct {
	pool      lexer.TokenPool
	nextToken func() (lexer.Token, error)
	now       *lexer.Token
	next      *lexer.Token
}

func NewParser() *Parser {
	parser := &Parser{
		pool: *lexer.Pool,
		now:  &lexer.Token{},
		next: &lexer.Token{},
	}

	parser.nextToken = parser.pool.Pop()
	parser.Next()
	_ = glg.Debug("Parser initialized")
	return parser
}

func (p *Parser) Next() {
	if *p.next != (lexer.Token{}) {
		p.now = p.next
	} else {
		token, err := p.nextToken()
		if err != nil {
			panic(err)
		}
		p.now = &token
	}
	token, _ := p.nextToken()
	p.next = &token
	_ = glg.Debug(p.now.String())
}
