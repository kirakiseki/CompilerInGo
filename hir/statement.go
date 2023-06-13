package hir

// 将ast中的statement整合转换为hir中的statement

type Statement interface {
	stmt()
}

type ConditionalStatement struct {
	Condition ConditionalExp
	IfBody    *Statement
	ElseBody  *Statement
}

type LoopStatement struct {
	Condition ConditionalExp
	Body      *Statement
}

type CallStatement struct {
	Method   string
	ActParam []Exp
}

type AssignStatement struct {
	Target string
	Exp    Exp
}

type ReturnStatement struct {
	Exp Exp
}

type BreakStatement struct{}

type ContinueStatement struct{}

type LocalVariableDeclaration struct {
	TypeIDPair []TypeIDPair
}

type Block struct {
	Statements []*Statement
}

func (c ConditionalStatement) stmt() {}

func NewConditionalStatement(condition ConditionalExp, ifBody, elseBody *Statement) ConditionalStatement {
	return ConditionalStatement{
		Condition: condition,
		IfBody:    ifBody,
		ElseBody:  elseBody,
	}
}

func (l LoopStatement) stmt() {}

func NewLoopStatement(condition ConditionalExp, body *Statement) LoopStatement {
	return LoopStatement{
		Condition: condition,
		Body:      body,
	}
}

func (c CallStatement) stmt() {}

func NewCallStatement(method string, actParam []Exp) CallStatement {
	return CallStatement{
		Method:   method,
		ActParam: actParam,
	}
}

func (a AssignStatement) stmt() {}

func NewAssignStatement(target string, exp Exp) AssignStatement {
	return AssignStatement{
		Target: target,
		Exp:    exp,
	}
}

func (r ReturnStatement) stmt() {}

func NewReturnStatement(exp Exp) ReturnStatement {
	return ReturnStatement{
		Exp: exp,
	}
}

func (b BreakStatement) stmt() {}

func NewBreakStatement() BreakStatement {
	return BreakStatement{}
}

func (c ContinueStatement) stmt() {}

func NewContinueStatement() ContinueStatement {
	return ContinueStatement{}
}

func (l LocalVariableDeclaration) stmt() {}

func NewLocalVariableDeclaration(typeIDPair []TypeIDPair) LocalVariableDeclaration {
	return LocalVariableDeclaration{
		TypeIDPair: typeIDPair,
	}
}

func (b Block) stmt() {}

func NewBlock(statements []*Statement) Block {
	return Block{
		Statements: statements,
	}
}
