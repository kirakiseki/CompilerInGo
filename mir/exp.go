package mir

import (
	"CompilerInGo/hir"
	"CompilerInGo/parser/ast"
	"fmt"
)

func (g *MIRGenerator) generateExp(exp hir.Exp) ([]Statement, int) {
	if exp.Op == ast.EMPTY {
		return g.generateTerm(exp.LTerm)
	}
	var stmtSeq []Statement
	lTermStmtSeq, lTermResultID := g.generateTerm(exp.LTerm)
	rTermStmtSeq, rTermResultID := g.generateTerm(exp.RTerm)
	stmtSeq = append(stmtSeq, lTermStmtSeq...)
	stmtSeq = append(stmtSeq, rTermStmtSeq...)
	switch exp.Op {
	case ast.PLUS:
		resultID := g.NewAnonymousVar()
		stmtSeq = append(stmtSeq, *NewStatement(PLUS, StrParam(hir.VarToStr(lTermResultID)), StrParam(hir.VarToStr(rTermResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s + %s", hir.VarToStr(resultID), hir.VarToStr(lTermResultID), hir.VarToStr(rTermResultID))))
		return stmtSeq, resultID
	case ast.MINUS:
		resultID := g.NewAnonymousVar()
		stmtSeq = append(stmtSeq, *NewStatement(MINUS, StrParam(hir.VarToStr(lTermResultID)), StrParam(hir.VarToStr(rTermResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s - %s", hir.VarToStr(resultID), hir.VarToStr(lTermResultID), hir.VarToStr(rTermResultID))))
		return stmtSeq, resultID
	}
	return nil, 0
}

func (g *MIRGenerator) generateTerm(term hir.Term) ([]Statement, int) {
	if term.Op == ast.EMPTY {
		return g.generateFactor(term.LFactor)
	}
	var stmtSeq []Statement
	lFactorStmtSeq, lFactorResultID := g.generateFactor(term.LFactor)
	rFactorStmtSeq, rFactorResultID := g.generateFactor(term.RFactor)
	stmtSeq = append(stmtSeq, lFactorStmtSeq...)
	stmtSeq = append(stmtSeq, rFactorStmtSeq...)
	switch term.Op {
	case ast.TIMES:
		resultID := g.NewAnonymousVar()
		stmtSeq = append(stmtSeq, *NewStatement(TIMES, StrParam(hir.VarToStr(lFactorResultID)), StrParam(hir.VarToStr(rFactorResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s * %s", hir.VarToStr(resultID), hir.VarToStr(lFactorResultID), hir.VarToStr(rFactorResultID))))
		return stmtSeq, resultID
	case ast.DIVIDE:
		resultID := g.NewAnonymousVar()
		stmtSeq = append(stmtSeq, *NewStatement(DIVIDE, StrParam(hir.VarToStr(lFactorResultID)), StrParam(hir.VarToStr(rFactorResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s / %s", hir.VarToStr(resultID), hir.VarToStr(lFactorResultID), hir.VarToStr(rFactorResultID))))
		return stmtSeq, resultID
	}
	return nil, 0
}

func (g *MIRGenerator) generateFactor(factor hir.Factor) ([]Statement, int) {
	switch factor.(type) {
	case *hir.Exp:
		return g.generateExp(*factor.(*hir.Exp))
	case hir.ID:
		return nil, g.GetVar(string(factor.(hir.ID)))
	case *hir.Integer:
		intVar := g.NewAnonymousVar()
		return []Statement{*NewStatement(ASSIGN, StrParam(hir.VarToStr(intVar)), IntParam(factor.(*hir.Integer).Val), StrParam(hir.VarToStr(intVar)), fmt.Sprintf("%s = %s", StrParam(hir.VarToStr(intVar)), IntParam(factor.(*hir.Integer).Val).Str()))}, intVar
	case *hir.Float:
		floatVar := g.NewAnonymousVar()
		return []Statement{*NewStatement(ASSIGN, StrParam(hir.VarToStr(floatVar)), FloatParam(factor.(*hir.Float).Val), StrParam(hir.VarToStr(floatVar)), fmt.Sprintf("%s = %s", StrParam(hir.VarToStr(floatVar)), FloatParam(factor.(*hir.Float).Val).Str()))}, floatVar
	default:
		return nil, 0
	}
}
