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

func (g *MIRGenerator) generateCompExp(compExp hir.CompExp) ([]Statement, int) {
	if compExp.Op == ast.EMPTY {
		return g.generateExp(compExp.LExp)
	}
	var stmtSeq []Statement
	lExpStmtSeq, lExpResultID := g.generateExp(compExp.LExp)
	rExpStmtSeq, rExpResultID := g.generateExp(compExp.RExp)
	stmtSeq = append(stmtSeq, lExpStmtSeq...)
	stmtSeq = append(stmtSeq, rExpStmtSeq...)
	resultID := g.NewAnonymousVar()
	switch compExp.Op {
	case ast.EQUAL:
		stmtSeq = append(stmtSeq, *NewStatement(JNEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s == %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s == %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("equal: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s == %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.DIAMOND:
		stmtSeq = append(stmtSeq, *NewStatement(JEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s != %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s != %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("notEqual: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s != %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.GREATER:
		stmtSeq = append(stmtSeq, *NewStatement(JLESSEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s > %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s > %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("greater: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s > %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.GREATEREQUAL:
		stmtSeq = append(stmtSeq, *NewStatement(JLESS, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s >= %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s >= %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("greaterEqual: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s >= %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.LESS:
		stmtSeq = append(stmtSeq, *NewStatement(JGREATEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s < %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s < %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("less: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s < %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.LESSEQUAL:
		stmtSeq = append(stmtSeq, *NewStatement(JGREAT, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s <= %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s <= %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("lessEqual: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s <= %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	}
	return stmtSeq, resultID
}

func (g *MIRGenerator) generateRelationalExp(relationalExp hir.RelationExp) ([]Statement, int) {
	if relationalExp.Op == ast.EMPTY {
		return g.generateCompExp(relationalExp.LExp)
	}
	var stmtSeq []Statement
	lExpStmtSeq, lExpResultID := g.generateCompExp(relationalExp.LExp)
	stmtSeq = append(stmtSeq, lExpStmtSeq...)
	resultID := g.NewAnonymousVar()
	stmtSeq = append(stmtSeq, *NewStatement(JNZERO, StrParam(hir.VarToStr(lExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s true: goto here+3", hir.VarToStr(lExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 5)), fmt.Sprintf("goto here+len(rExpStmtSeq)+5)")))
	rExpStmtSeq, rExpResultID := g.generateCompExp(relationalExp.RExp)
	stmtSeq = append(stmtSeq, rExpStmtSeq...)
	stmtSeq = append(stmtSeq, *NewStatement(JZERO, StrParam(hir.VarToStr(rExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s false: goto here+3", hir.VarToStr(rExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s true: %s = 1", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("goto here+2")))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s false: %s = 0", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	offset := len(rExpStmtSeq)
	stmtSeq[len(lExpStmtSeq)+4].Res = StrParam(fmt.Sprintf("_T_JMP_REF_%d", offset+5))
	return stmtSeq, resultID
}

func (g *MIRGenerator) generateConditionalExp(conditionalExp hir.ConditionalExp) ([]Statement, int) {
	if conditionalExp.Op == ast.EMPTY {
		return g.generateRelationalExp(conditionalExp.LExp)
	}

	var stmtSeq []Statement
	lExpStmtSeq, lExpResultID := g.generateRelationalExp(conditionalExp.LExp)
	stmtSeq = append(stmtSeq, lExpStmtSeq...)
	resultID := g.NewAnonymousVar()
	stmtSeq = append(stmtSeq, *NewStatement(JZERO, StrParam(hir.VarToStr(lExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s false: goto here+3", hir.VarToStr(lExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 5)), fmt.Sprintf("goto here+len(rExpStmtSeq)+5)")))
	rExpStmtSeq, rExpResultID := g.generateRelationalExp(conditionalExp.RExp)
	stmtSeq = append(stmtSeq, rExpStmtSeq...)
	stmtSeq = append(stmtSeq, *NewStatement(JZERO, StrParam(hir.VarToStr(rExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s false: goto here+3", hir.VarToStr(rExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s true: %s = 1", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("goto here+2")))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s false: %s = 0", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	offset := len(rExpStmtSeq)
	stmtSeq[len(lExpStmtSeq)+4].Res = StrParam(fmt.Sprintf("_T_JMP_REF_%d", offset+5))
	return stmtSeq, resultID
}
