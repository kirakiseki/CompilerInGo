package parser

import (
	"CompilerInGo/lexer"
	"CompilerInGo/parser/ast"
	"fmt"
	"github.com/kpango/glg"
)

type Parser struct {
	// AST根结点
	program *ast.Program
	// 当前token
	token *lexer.Token
	// token流
	*TokenStream
	err error
}

// NewParser 创建一个新的Parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse 开始解析
func (p *Parser) Parse() (program *ast.Program, err error) {
	// 错误处理
	defer func() {
		// 如果有错误，打印错误
		if r := recover(); r != p.err {
			_ = glg.Error(r)
		}
		// 返回结果
		program, err = p.program, p.err
	}()

	// 初始化token流
	ts := NewTokenStream()
	p.TokenStream = &ts

	// 开始解析
	p.parse()

	return p.program, p.err
}

// parse 解析
func (p *Parser) parse() {
	// 创建一个新的AST根结点
	program, _ := ast.NewProgram()
	p.program = &program

	// 不断读取token，直到EOF
	for {
		token := p.ReadToken()
		p.token = &token

		switch {
		case token.Type == lexer.EOF_LITERAL:
			//读到EOF，解析结束
			return
		case ast.IsResultType(token):
			// 读到返回值类型，解析函数
			p.program.Method = append(p.program.Method, *p.parseMethod())
		default:
			// 读到其他类型的token，报错
			p.Errorf("Unexpected token %s", token.String())
		}
	}
}

// parseMethod 解析函数
func (p *Parser) parseMethod() *ast.Method {
	resultType := (ast.ResultType)(*p.token)                     //返回值类型
	ident := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER)) //函数名
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)              //左括号
	paramList := p.parseParamList()                              //参数列表
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)              //右括号
	block := p.parseBlock()                                      //函数体

	method, err := ast.NewMethod(resultType, ident, lParen, *paramList, rParen, *block)
	if err != nil {
		panic(err)
	}

	return &method
}

// parseParamList 解析参数列表
func (p *Parser) parseParamList() *ast.ParamList {
	// 可选参数列表，可能为空，使用Optional判断是否接受
	typ, notEmpty := p.OptionalAcceptTokenByFunc(ast.IsType)
	if !notEmpty {
		// 如果为空，返回空的参数列表
		paramList, _ := ast.NewParamList()
		return &paramList
	}

	typToken := (ast.Type)(typ)                               //参数类型
	id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER)) //参数名

	// 逗号+类型+参数名的tuple数组
	commaTypeIDTuple := []any{typToken, id}

	for {
		// 逗号，可能为空，使用Optional判断是否接受
		comma, isComma := p.OptionalAcceptTokenByType(lexer.COMMA)
		if !isComma {
			// 不为逗号，返回参数列表
			paramList, _ := ast.NewParamList(commaTypeIDTuple...)
			return &paramList
		}
		// 逗号后面必须是类型和参数名
		typ := (ast.Type)(p.MustAcceptTokenByFunc(ast.IsType))
		id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))
		// 添加到tuple数组
		commaTypeIDTuple = append(commaTypeIDTuple, comma, typ, id)
	}
}

// parseBlock 解析代码块
func (p *Parser) parseBlock() *ast.Block {
	lBrace := p.MustAcceptTokenByType(lexer.LBRACE) //左大括号
	statements := p.parseStmtList()                 //语句列表
	rBrace := p.MustAcceptTokenByType(lexer.RBRACE) //右大括号

	block, err := ast.NewBlock(lBrace, statements, rBrace)
	if err != nil {
		panic(err)
	}

	return &block
}

// parseStmtList 解析语句列表
func (p *Parser) parseStmtList() []ast.Statement {
	statements := make([]ast.Statement, 0)

	// 不断解析语句，直到返回空语句
	for {
		statement := p.parseStmt()
		if statement == (ast.Statement{}) {
			return statements
		}
		statements = append(statements, statement)
	}
}

// parseStmt 解析单条语句
func (p *Parser) parseStmt() ast.Statement {
	// 读取一个token，根据token类型解析语句
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

	// 根据token类型判断语句类型
	switch token.Type {
	case lexer.CALL:
		// 函数调用语句
		return ast.Statement{
			Statement: p.parseCallStmt(),
			Type:      ast.CALLSTATEMENT,
		}
	case lexer.IF:
		// 条件语句
		return ast.Statement{
			Statement: p.parseIfStmt(),
			Type:      ast.CONDITIONALSTATEMENT,
		}
	case lexer.WHILE:
		// 循环语句
		return ast.Statement{
			Statement: p.parseLoopStmt(),
			Type:      ast.LOOPSTATEMENT,
		}
	case lexer.RETURN:
		// 返回语句
		return ast.Statement{
			Statement: p.parseReturnStmt(),
			Type:      ast.RETURNSTATEMENT,
		}
	case lexer.BREAK:
		// 跳出语句
		return ast.Statement{
			Statement: p.parseBreakStmt(),
			Type:      ast.BREAKSTATEMENT,
		}
	case lexer.CONTINUE:
		// 继续语句
		return ast.Statement{
			Statement: p.parseContinueStmt(),
			Type:      ast.CONTINUESTATEMENT,
		}
	case lexer.LBRACE:
		// 代码块
		// 回退{，由parseBlock解析
		p.UnreadToken()
		return ast.Statement{
			Statement: p.parseBlock(),
			Type:      ast.BLOCK,
		}
	case lexer.RBRACE:
		// 代码块结束
		p.UnreadToken()
		return ast.Statement{}
	case lexer.SEMICOLON:
		// 空语句
		return ast.Statement{}
	default:
		if ast.IsID(token) {
			// 赋值语句
			return ast.Statement{
				Statement: p.parseAssignStmt(),
				Type:      ast.ASSIGNMENTSTATEMENT,
			}
		} else if ast.IsType(token) {
			// 变量声明语句
			return ast.Statement{
				Statement: p.parseLocalVariableDeclarationStmt(),
				Type:      ast.LOCALVARIABLEDECLARATION,
			}
		} else {
			// 未知语句
			p.Errorf("Unexpected token %s", token.String())
			return ast.Statement{}
		}
	}
}

// Errorf 存储错误并panic
func (p *Parser) Errorf(format string, args ...interface{}) {
	p.err = fmt.Errorf(format, args...)
	panic(p.err)
}

// ErrorToken 存储错误并panic
func (p *Parser) ErrorToken(token lexer.Token, format string, args ...interface{}) {
	p.err = fmt.Errorf("%s: %s", token.String(), fmt.Sprintf(format, args...))
	panic(p.err)
}

// parseBreakStmt 解析跳出语句
func (p *Parser) parseBreakStmt() *ast.BreakStatement {
	token := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewBreakStatement(*p.token, token)
	return &stmt
}

// parseContinueStmt 解析继续语句
func (p *Parser) parseContinueStmt() *ast.ContinueStatement {
	token := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewContinueStatement(*p.token, token)
	return &stmt
}

// parseLoopStmt 解析循环语句
func (p *Parser) parseReturnStmt() *ast.ReturnStatement {
	exp := p.parseExp()
	semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)
	stmt, _ := ast.NewReturnStatement(*p.token, exp, semicolon)
	return &stmt
}

// parseCallStmt 解析函数调用语句
func (p *Parser) parseCallStmt() *ast.CallStatement {
	token := *p.token                                         // call
	id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER)) // 函数名
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)           // (
	actParamList := p.parseActParamList()                     // 实参列表
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)           // )
	semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)     // ;

	stmt, _ := ast.NewCallStatement(token, id, lParen, actParamList, rParen, semicolon)
	return &stmt
}

// parseIfStmt 解析条件语句
func (p *Parser) parseIfStmt() *ast.ConditionalStatement {
	token := *p.token                                             // if
	lParen := p.MustAcceptTokenByType(lexer.LPAREN)               // (
	conditionalExp := p.parseConditionalExp()                     // 条件表达式
	rParen := p.MustAcceptTokenByType(lexer.RPAREN)               // )
	stmt := p.parseStmt()                                         // if语句
	elseToken, hasElse := p.OptionalAcceptTokenByType(lexer.ELSE) // 是否有else
	if hasElse {
		// 有else
		elseStmt := p.parseStmt()
		conditionalStmt, _ := ast.NewConditionalStatement(token, lParen, *conditionalExp, rParen, stmt, &elseToken, &elseStmt)
		return &conditionalStmt
	} else {
		conditionalStmt, _ := ast.NewConditionalStatement(token, lParen, *conditionalExp, rParen, stmt, nil, nil)
		return &conditionalStmt
	}
}

// parseLoopStmt 解析循环语句
func (p *Parser) parseLoopStmt() *ast.LoopStatement {
	token := *p.token                               // while
	lParen := p.MustAcceptTokenByType(lexer.LPAREN) // (
	conditionalExp := p.parseConditionalExp()       // 条件表达式
	rParen := p.MustAcceptTokenByType(lexer.RPAREN) // )
	stmt := p.parseStmt()                           // while语句

	loopStmt, _ := ast.NewLoopStatement(token, lParen, *conditionalExp, rParen, stmt)
	return &loopStmt
}

// parseExpStmt 解析表达式语句
func (p *Parser) parseAssignStmt() *ast.AssignmentStatement {
	token := *p.token                                     // id
	equal := p.MustAcceptTokenByType(lexer.ASSIGN)        // =
	exp := p.parseExp()                                   // 表达式
	semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON) // ;

	stmt, _ := ast.NewAssignmentStatement(ast.ID(token), equal, *exp, semicolon)
	return &stmt
}

// parseLocalVariableDeclarationStmt 解析变量声明语句
func (p *Parser) parseLocalVariableDeclarationStmt() *ast.LocalVariableDeclaration {
	token := *p.token                                         // 类型
	id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER)) // 变量名

	commaIDPair := make([]any, 0) // 逗号-变量名对

	// 不断解析逗号-变量名对
	for {
		// 检查是否有逗号
		comma, isComma := p.OptionalAcceptTokenByType(lexer.COMMA)
		if !isComma {
			// 没有逗号，结束
			// 解析分号
			semicolon := p.MustAcceptTokenByType(lexer.SEMICOLON)

			paramList, _ := ast.NewLocalVariableDeclarationStatement(ast.Type(token), id, commaIDPair, semicolon)
			return &paramList
		}
		// 解析变量名
		id := (ast.ID)(p.MustAcceptTokenByType(lexer.IDENTIFIER))
		// 添加逗号-变量名对到commaIDPair列表
		commaIDPair = append(commaIDPair, comma, id)
	}
}

// parseExp 解析表达式
func (p *Parser) parseExp() *ast.Exp {
	// 左项
	lTerm := p.parseTerm()
	// 检查是否有加号或减号
	plusOrMinus, isPlusOrMinus := p.OptionalAcceptTokenByType(lexer.PLUS, lexer.MINUS)
	if !isPlusOrMinus {
		// 没有加号或减号，结束
		exp, _ := ast.NewExp(*lTerm, nil, nil)
		return &exp
	}
	// 有加减号，解析右项
	rTerm := p.parseTerm()

	exp, _ := ast.NewExp(*lTerm, &plusOrMinus, rTerm)
	return &exp
}

// parseTerm 解析项
func (p *Parser) parseTerm() *ast.Term {
	// 左因子
	lFactor := p.parseFactor()
	// 检查是否有乘号或除号
	mulOrDiv, isMulOrDiv := p.OptionalAcceptTokenByType(lexer.TIMES, lexer.DIVIDE)
	if !isMulOrDiv {
		// 没有乘号或除号，结束
		term, _ := ast.NewTerm(*lFactor, nil, nil)
		return &term
	}
	// 有乘除号，解析右因子
	rFactor := p.parseFactor()

	term, _ := ast.NewTerm(*lFactor, &mulOrDiv, rFactor)
	return &term
}

// parseFactor 解析因子
func (p *Parser) parseFactor() *ast.Factor {
	// 判断单token或(Exp)
	token := p.MustAcceptTokenByFunc(func(token lexer.Token) bool {
		return ast.IsID(token) || token.Type == lexer.INTEGER_LITERAL || token.Type == lexer.DECIMAL_LITERAL || token.Type == lexer.LPAREN
	})

	if token.Type == lexer.INTEGER_LITERAL || token.Type == lexer.DECIMAL_LITERAL || ast.IsID(token) {
		// 单token
		factor, _ := ast.NewFactor(token)
		return &factor
	} else if token.Type == lexer.LPAREN {
		// (Exp)
		exp := p.parseExp()                             // Exp
		rParen := p.MustAcceptTokenByType(lexer.RPAREN) // )

		factor, _ := ast.NewFactor(token, exp, rParen)
		return &factor
	} else {
		// 未知token
		p.Errorf("Unexpected token %s", token.String())
		return nil
	}
}

// parseConditionalExp 解析条件表达式
func (p *Parser) parseConditionalExp() *ast.ConditionalExp {
	// 左关系表达式
	lRelationExp := p.parseRelationExp()
	// 检查是否有or
	or, hasOr := p.OptionalAcceptTokenByType(lexer.OR)
	if !hasOr {
		// 没有or，结束
		conditionalExp, _ := ast.NewConditionalExp(*lRelationExp, nil, nil)
		return &conditionalExp
	}

	// 有or，解析右关系表达式
	rRelationExp := p.parseRelationExp()
	conditionalExp, _ := ast.NewConditionalExp(*lRelationExp, &or, rRelationExp)

	return &conditionalExp
}

// parseRelationExp 解析关系表达式
func (p *Parser) parseRelationExp() *ast.RelationExp {
	// 左比较表达式
	lCompExp := p.parseCompExp()
	// 检查是否有and
	and, hasAnd := p.OptionalAcceptTokenByType(lexer.AND)
	if !hasAnd {
		// 没有and，结束
		relationExp, _ := ast.NewRelationExp(*lCompExp, nil, nil)
		return &relationExp
	}

	// 有and，解析右比较表达式
	rCompExp := p.parseCompExp()
	relationExp, _ := ast.NewRelationExp(*lCompExp, &and, rCompExp)

	return &relationExp
}

// parseCompExp 解析比较表达式
func (p *Parser) parseCompExp() *ast.CompExp {
	// 左表达式
	lExp := p.parseExp()
	// 比较运算符
	cmpOp := p.MustAcceptTokenByType(lexer.LESS, lexer.LESSEQUAL, lexer.GREATER, lexer.GREATEREQUAL, lexer.EQUAL, lexer.DIAMOND)
	// 右表达式
	rExp := p.parseExp()

	compExp, _ := ast.NewCompExp(*lExp, ast.CmpOp(cmpOp), *rExp)
	return &compExp
}

// parseActParamList 解析实参列表
func (p *Parser) parseActParamList() *ast.ActParamList {
	// 检查是否有实参
	_, hasExp := p.OptionalAcceptTokenByFunc(func(token lexer.Token) bool {
		return ast.IsID(token) || token.Type == lexer.INTEGER_LITERAL || token.Type == lexer.DECIMAL_LITERAL || token.Type == lexer.LPAREN
	})
	if !hasExp {
		// 没有实参，结束
		actParamList, _ := ast.NewActParamList(nil)
		return &actParamList
	}
	// 回退一个Exp中的token
	p.UnreadToken()

	// 解析实参
	exp := p.parseExp()

	// 解析剩下的逗号和实参
	commaExpPair := make([]any, 0)
	// 保存第一个实参
	commaExpPair = append(commaExpPair, *exp)
	for {
		// 检查是否有逗号
		comma, isComma := p.OptionalAcceptTokenByType(lexer.COMMA)
		if !isComma {
			// 没有逗号，结束
			actParamList, _ := ast.NewActParamList(commaExpPair)
			return &actParamList
		}
		// 解析实参
		exp := p.parseExp()
		// 保存到commaExpPair列表
		commaExpPair = append(commaExpPair, comma, *exp)
	}
}
