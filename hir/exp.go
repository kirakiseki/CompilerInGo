package hir

import (
	"CompilerInGo/parser/ast"
)

type ConditionalExp struct {
	LExp RelationExp
	Op   int
	RExp RelationExp
}

type RelationExp struct {
	LExp CompExp
	Op   int
	RExp CompExp
}

type CompExp struct {
	LExp Exp
	Op   int
	RExp Exp
}

type Exp struct {
	LTerm Term
	Op    int
	RTerm Term
}

type Term struct {
	LFactor Factor
	Op      int
	RFactor Factor
}

type Factor interface {
	factor()
}

func NewConditionalExp(lExp *RelationExp, rExp *RelationExp) ConditionalExp {
	if rExp == nil {
		return ConditionalExp{
			LExp: *lExp,
			Op:   ast.EMPTY,
		}
	}

	return ConditionalExp{
		LExp: *lExp,
		Op:   ast.OR,
		RExp: *rExp,
	}
}

func NewRelationExp(lExp *CompExp, rExp *CompExp) RelationExp {
	if rExp == nil {
		return RelationExp{
			LExp: *lExp,
			Op:   ast.EMPTY,
		}
	}

	return RelationExp{
		LExp: *lExp,
		Op:   ast.AND,
		RExp: *rExp,
	}
}

func NewCompExp(lExp *Exp, op int, rExp *Exp) CompExp {
	if op == ast.EMPTY || rExp == nil {
		return CompExp{
			LExp: *lExp,
			Op:   ast.EMPTY,
		}
	}

	return CompExp{
		LExp: *lExp,
		Op:   op,
		RExp: *rExp,
	}
}

func NewExp(lTerm *Term, op int, rTerm *Term) Exp {
	if op == ast.EMPTY || rTerm == nil {
		return Exp{
			LTerm: *lTerm,
			Op:    ast.EMPTY,
		}
	}

	return Exp{
		LTerm: *lTerm,
		Op:    op,
		RTerm: *rTerm,
	}
}

func NewTerm(lFactor *Factor, op int, rFactor *Factor) Term {
	if op == ast.EMPTY || rFactor == nil {
		return Term{
			LFactor: *lFactor,
			Op:      ast.EMPTY,
		}
	}
	return Term{
		LFactor: *lFactor,
		Op:      op,
		RFactor: *rFactor,
	}
}

func (e Exp) factor() {}
