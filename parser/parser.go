package parser

import (
	"CompilerInGo/lexer"
	"CompilerInGo/parser/ast"
	"fmt"
	"github.com/kpango/glg"
)

type Parser struct {
	program *ast.Program
	token   *lexer.Token
	*TokenStream
	err error
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse() (program *ast.Program, err error) {
	defer func() {
		if r := recover(); r != p.err {
			_ = glg.Error(r)
		}
		program, err = p.program, p.err
	}()

	p.TokenStream = &TokenStream{
		pool: *lexer.Pool,
	}

	p.parse()

	return p.program, p.err
}

func (p *Parser) parse() {
	program, _ := ast.NewProgram()
	p.program = &program

	for {
		token := p.ReadToken()
		p.token = &token
		switch {
		case token.Type == lexer.EOF_LITERAL:
			return
		case ast.IsResultType(token):
			p.program.Method = append(p.program.Method, *p.parseMethod())
		default:
			p.Errorf("Unexpected token %s", token.String())
		}
	}
}

func (p *Parser) parseMethod() *ast.Method {
	resuleType := (ast.ResultType)(*p.token)
	ident := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)
	paramList := p.parseParamList()
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)
	block := p.parseBlock()

	method, err := ast.NewMethod(resuleType, ident, lParen, *paramList, rParen, *block)
	if err != nil {
		panic(err)
	}

	return &method
}

func (p *Parser) parseParamList() *ast.ParamList {
	typ, isType := p.OptionalAcceptTokenByFunc(ast.IsType)
	if !isType {
		paramList, _ := ast.NewParamList()
		return &paramList
	}

	typToken := (ast.Type)(typ)
	id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))

	commaTypeIDTuple := []any{typToken, id}

	for {
		comma, isComma := p.OptionalAcceptTokenByType(lexer.COMMA)
		if !isComma {
			paramList, _ := ast.NewParamList(commaTypeIDTuple...)
			return &paramList
		}
		typ := (ast.Type)(p.MustAcceptTokenByFunc(ast.IsType))
		id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))
		commaTypeIDTuple = append(commaTypeIDTuple, comma, typ, id)
	}
}

func (p *Parser) parseBlock() *ast.Block {
	lBrace := p.MustAcceptTokenByType(lexer.LBRACE)
	statements := p.parseStmtList()
	rBrace := p.MustAcceptTokenByType(lexer.RBRACE)
	block, err := ast.NewBlock(lBrace, statements, rBrace)
	if err != nil {
		panic(err)
	}

	return &block
}

func (p *Parser) parseStmtList() []ast.Statement {
	statements := make([]ast.Statement, 0)
	for {
		statement := p.parseStmt()
		if statement == (ast.Statement{}) {
			return statements
		}
		statements = append(statements, statement)
	}
}

func (p *Parser) parseStmt() ast.Statement {
	token := p.MustAcceptTokenByFunc(func(token lexer.Token) bool {
		return ast.IsID(token) ||
			ast.IsType(token) ||
			token.Type == lexer.CALL ||
			token.Type == lexer.IF ||
			token.Type == lexer.WHILE ||
			token.Type == lexer.RETURN ||
			token.Type == lexer.BREAK ||
			token.Type == lexer.CONTINUE ||
			token.Type == lexer.LBRACE ||
			token.Type == lexer.RBRACE ||
			token.Type == lexer.SEMICOLON
	})
	p.token = &token

	switch token.Type {
	case lexer.CALL:
		return ast.Statement{
			Statement: p.parseCallStmt(),
			Type:      ast.CALLSTATEMENT,
		}
	case lexer.IF:
		return ast.Statement{
			Statement: p.parseIfStmt(),
			Type:      ast.CONDITIONALSTATEMENT,
		}
	case lexer.WHILE:
		return ast.Statement{
			Statement: p.parseLoopStmt(),
			Type:      ast.LOOPSTATEMENT,
		}
	case lexer.RETURN:
		return ast.Statement{
			Statement: p.parseReturnStmt(),
			Type:      ast.RETURNSTATEMENT,
		}
	case lexer.BREAK:
		return ast.Statement{
			Statement: p.parseBreakStmt(),
			Type:      ast.BREAKSTATEMENT,
		}
	case lexer.CONTINUE:
		return ast.Statement{
			Statement: p.parseContinueStmt(),
			Type:      ast.CONTINUESTATEMENT,
		}
	case lexer.LBRACE:
		p.UnreadToken()
		return ast.Statement{
			Statement: p.parseBlock(),
			Type:      ast.BLOCK,
		}
	case lexer.RBRACE:
		p.UnreadToken()
		return ast.Statement{}
	case lexer.SEMICOLON:
		return ast.Statement{}
	default:
		if ast.IsID(token) {
			return ast.Statement{
				Statement: p.parseAssignStmt(),
				Type:      ast.ASSIGNMENTSTATEMENT,
			}
		} else if ast.IsType(token) {
			return ast.Statement{
				Statement: p.parseLocalVariableDeclarationStmt(),
				Type:      ast.LOCALVARIABLEDECLARATION,
			}
		} else {
			p.Errorf("Unexpected token %s", token.String())
			return ast.Statement{}
		}
	}
}

func (p *Parser) Errorf(format string, args ...interface{}) {
	p.err = fmt.Errorf(format, args...)
	panic(p.err)
}

func (p *Parser) ErrorToken(token lexer.Token, format string, args ...interface{}) {
	p.err = fmt.Errorf("%s: %s", token.String(), fmt.Sprintf(format, args...))
	panic(p.err)
}

func (p *Parser) parseBreakStmt() *ast.BreakStatement {
	token := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewBreakStatement(*p.token, token)
	return &stmt
}

func (p *Parser) parseContinueStmt() *ast.ContinueStatement {
	token := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewContinueStatement(*p.token, token)
	return &stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStatement {
	exp := p.parseExp()
	semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewReturnStatement(*p.token, exp, semicolon)
	return &stmt
}

func (p *Parser) parseCallStmt() *ast.CallStatement {
	token := *p.token
	id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)
	actParamList := p.parseActParamList()
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)
	semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)

	stmt, _ := ast.NewCallStatement(token, id, lParen, actParamList, rParen, semicolon)
	return &stmt
}

func (p *Parser) parseIfStmt() *ast.ConditionalStatement {
	token := *p.token
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)
	conditionalExp := p.parseConditionalExp()
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)
	stmt := p.parseStmt()
	elseToken, hasElse := p.OptionalAcceptTokenByType(lexer.ELSE)
	if hasElse {
		elseStmt := p.parseStmt()
		conditionalStmt, _ := ast.NewConditionalStatement(token, lParen, *conditionalExp, rParen, stmt, &elseToken, &elseStmt)
		return &conditionalStmt
	} else {
		conditionalStmt, _ := ast.NewConditionalStatement(token, lParen, *conditionalExp, rParen, stmt, nil, nil)
		return &conditionalStmt
	}
}

func (p *Parser) parseLoopStmt() *ast.LoopStatement {
	token := *p.token
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)
	conditionalExp := p.parseConditionalExp()
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)
	stmt := p.parseStmt()
	loopStmt, _ := ast.NewLoopStatement(token, lParen, *conditionalExp, rParen, stmt)
	return &loopStmt
}

func (p *Parser) parseAssignStmt() *ast.AssignmentStatement {
	equal := p.MustAcceptTokenByType(lexer.ASSIGN)
	exp := p.parseExp()
	semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewAssignmentStatement(ast.ID(*p.token), equal, *exp, semicolon)
	return &stmt
}

func (p *Parser) parseLocalVariableDeclarationStmt() *ast.LocalVariableDeclaration {
	id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))

	commaIDPair := make([]any, 0)

	for {
		comma, isComma := p.OptionalAcceptTokenByType(lexer.COMMA)
		if !isComma {
			semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)
			paramList, _ := ast.NewLocalDeclarationStatement(ast.Type(*p.token), id, commaIDPair, semicolon)
			return &paramList
		}
		id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))
		commaIDPair = append(commaIDPair, comma, id)
	}
}

func (p *Parser) parseExp() *ast.Exp {
	lTerm := p.parseTerm()
	plusOrMinus, isPlusOrMinus := p.OptionalAcceptTokenByType(lexer.PLUS, lexer.MINUS)
	if !isPlusOrMinus {
		exp, _ := ast.NewExp(*lTerm, nil, nil)
		return &exp
	}
	rTerm := p.parseTerm()
	exp, _ := ast.NewExp(*lTerm, &plusOrMinus, rTerm)
	return &exp
}

func (p *Parser) parseTerm() *ast.Term {
	lFactor := p.parseFactor()
	mulOrDiv, isMulOrDiv := p.OptionalAcceptTokenByType(lexer.TIMES, lexer.DIVIDE)
	if !isMulOrDiv {
		term, _ := ast.NewTerm(*lFactor, nil, nil)
		return &term
	}
	rFactor := p.parseFactor()
	term, _ := ast.NewTerm(*lFactor, &mulOrDiv, rFactor)
	return &term
}

func (p *Parser) parseFactor() *ast.Factor {
	token := p.MustAcceptTokenByFunc(func(token lexer.Token) bool {
		return ast.IsID(token) || token.Type == lexer.INTEGER_LITERAL || token.Type == lexer.DECIMAL_LITERAL || token.Type == lexer.LPAREN
	})

	if token.Type == lexer.INTEGER_LITERAL || token.Type == lexer.DECIMAL_LITERAL || ast.IsID(token) {
		factor, _ := ast.NewFactor(token)
		return &factor
	} else if token.Type == lexer.LPAREN {
		exp := p.parseExp()
		rParen := p.MustAcceptTokenByType(lexer.RPAREN)
		factor, _ := ast.NewFactor(token, exp, rParen)
		return &factor
	} else {
		p.Errorf("Unexpected token %s", token.String())
		return nil
	}
}

func (p *Parser) parseConditionalExp() *ast.ConditionalExp {
	lRelationExp := p.parseRelationExp()
	or, hasOr := p.OptionalAcceptTokenByType(lexer.OR)
	if !hasOr {
		conditionalExp, _ := ast.NewConditionalExp(*lRelationExp, nil, nil)
		return &conditionalExp
	}
	rRelationExp := p.parseRelationExp()
	conditionalExp, _ := ast.NewConditionalExp(*lRelationExp, &or, rRelationExp)
	return &conditionalExp
}

func (p *Parser) parseRelationExp() *ast.RelationExp {
	lCompExp := p.parseCompExp()
	and, hasAnd := p.OptionalAcceptTokenByType(lexer.AND)
	if !hasAnd {
		relationExp, _ := ast.NewRelationExp(*lCompExp, nil, nil)
		return &relationExp
	}
	rCompExp := p.parseCompExp()
	relationExp, _ := ast.NewRelationExp(*lCompExp, &and, rCompExp)
	return &relationExp
}

func (p *Parser) parseCompExp() *ast.CompExp {
	lExp := p.parseExp()
	cmpOp := p.MustAcceptTokenByType(lexer.LESS, lexer.LESSEQUAL, lexer.GREATER, lexer.GREATEREQUAL, lexer.EQUAL, lexer.DIAMOND)
	rExp := p.parseExp()
	compExp, _ := ast.NewCompExp(*lExp, ast.CmpOp(cmpOp), *rExp)
	return &compExp
}

func (p *Parser) parseActParamList() *ast.ActParamList {
	_, hasExp := p.OptionalAcceptTokenByFunc(func(token lexer.Token) bool {
		return ast.IsID(token) || token.Type == lexer.INTEGER_LITERAL || token.Type == lexer.DECIMAL_LITERAL || token.Type == lexer.LPAREN
	})
	if !hasExp {
		actParamList, _ := ast.NewActParamList(nil)
		return &actParamList
	}
	p.UnreadToken()
	exp := p.parseExp()
	commaExpPair := make([]any, 0)
	commaExpPair = append(commaExpPair, *exp)
	for {
		comma, isComma := p.OptionalAcceptTokenByType(lexer.COMMA)
		if !isComma {
			actParamList, _ := ast.NewActParamList(commaExpPair)
			return &actParamList
		}
		exp := p.parseExp()
		commaExpPair = append(commaExpPair, comma, *exp)
	}
}
