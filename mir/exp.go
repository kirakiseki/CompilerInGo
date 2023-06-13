package mir

import (
	"CompilerInGo/hir"
	"CompilerInGo/parser/ast"
	"fmt"
)

// generateExp 生成算术表达式
func (g *MIRGenerator) generateExp(exp hir.Exp) ([]Statement, int) {
	// 右项为空
	if exp.Op == ast.EMPTY {
		return g.generateTerm(exp.LTerm)
	}

	// 语句序列
	var stmtSeq []Statement
	// 左项语句序列，左项结果变量
	lTermStmtSeq, lTermResultID := g.generateTerm(exp.LTerm)
	// 右项语句序列，右项结果变量
	rTermStmtSeq, rTermResultID := g.generateTerm(exp.RTerm)
	// 将左项语句序列和右项语句序列添加到语句序列中
	stmtSeq = append(stmtSeq, lTermStmtSeq...)
	stmtSeq = append(stmtSeq, rTermStmtSeq...)

	// 根据运算符生成语句序列
	switch exp.Op {
	case ast.PLUS:
		// 算术表达式结果变量
		resultID := g.NewAnonymousVar()
		// 生成算术表达式语句
		stmtSeq = append(stmtSeq, *NewStatement(PLUS, StrParam(hir.VarToStr(lTermResultID)), StrParam(hir.VarToStr(rTermResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s + %s", hir.VarToStr(resultID), hir.VarToStr(lTermResultID), hir.VarToStr(rTermResultID))))
		return stmtSeq, resultID
	case ast.MINUS:
		// 算术表达式结果变量
		resultID := g.NewAnonymousVar()
		// 生成算术表达式语句
		stmtSeq = append(stmtSeq, *NewStatement(MINUS, StrParam(hir.VarToStr(lTermResultID)), StrParam(hir.VarToStr(rTermResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s - %s", hir.VarToStr(resultID), hir.VarToStr(lTermResultID), hir.VarToStr(rTermResultID))))
		return stmtSeq, resultID
	}

	return nil, 0
}

// generateTerm 生成项
func (g *MIRGenerator) generateTerm(term hir.Term) ([]Statement, int) {
	// 右因子为空
	if term.Op == ast.EMPTY {
		return g.generateFactor(term.LFactor)
	}

	// 语句序列
	var stmtSeq []Statement
	// 左因子语句序列，左因子结果变量
	lFactorStmtSeq, lFactorResultID := g.generateFactor(term.LFactor)
	// 右因子语句序列，右因子结果变量
	rFactorStmtSeq, rFactorResultID := g.generateFactor(term.RFactor)
	// 将左因子语句序列和右因子语句序列添加到语句序列中
	stmtSeq = append(stmtSeq, lFactorStmtSeq...)
	stmtSeq = append(stmtSeq, rFactorStmtSeq...)

	// 根据运算符生成语句序列
	switch term.Op {
	case ast.TIMES:
		// 项结果变量
		resultID := g.NewAnonymousVar()
		// 生成项计算语句
		stmtSeq = append(stmtSeq, *NewStatement(TIMES, StrParam(hir.VarToStr(lFactorResultID)), StrParam(hir.VarToStr(rFactorResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s * %s", hir.VarToStr(resultID), hir.VarToStr(lFactorResultID), hir.VarToStr(rFactorResultID))))
		return stmtSeq, resultID
	case ast.DIVIDE:
		// 项结果变量
		resultID := g.NewAnonymousVar()
		// 生成项计算语句
		stmtSeq = append(stmtSeq, *NewStatement(DIVIDE, StrParam(hir.VarToStr(lFactorResultID)), StrParam(hir.VarToStr(rFactorResultID)), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("%s = %s / %s", hir.VarToStr(resultID), hir.VarToStr(lFactorResultID), hir.VarToStr(rFactorResultID))))
		return stmtSeq, resultID
	}

	return nil, 0
}

// generateFactor 生成因子
func (g *MIRGenerator) generateFactor(factor hir.Factor) ([]Statement, int) {
	// 按照因子类型生成语句序列
	switch factor.(type) {
	case *hir.Exp:
		// (Exp)
		return g.generateExp(*factor.(*hir.Exp))
	case hir.ID:
		// ID
		return nil, g.GetVar(string(factor.(hir.ID)))
	case *hir.Integer:
		// INTC
		intVar := g.NewAnonymousVar()
		return []Statement{*NewStatement(ASSIGN, StrParam(hir.VarToStr(intVar)), IntParam(factor.(*hir.Integer).Val), StrParam(hir.VarToStr(intVar)), fmt.Sprintf("%s = %s", StrParam(hir.VarToStr(intVar)), IntParam(factor.(*hir.Integer).Val).Str()))}, intVar
	case *hir.Float:
		// DECI
		floatVar := g.NewAnonymousVar()
		return []Statement{*NewStatement(ASSIGN, StrParam(hir.VarToStr(floatVar)), FloatParam(factor.(*hir.Float).Val), StrParam(hir.VarToStr(floatVar)), fmt.Sprintf("%s = %s", StrParam(hir.VarToStr(floatVar)), FloatParam(factor.(*hir.Float).Val).Str()))}, floatVar
	default:
		return nil, 0
	}
}

// generateCompExp 生成比较表达式
func (g *MIRGenerator) generateCompExp(compExp hir.CompExp) ([]Statement, int) {
	// 右算术表达式为空
	if compExp.Op == ast.EMPTY {
		return g.generateExp(compExp.LExp)
	}

	// 语句序列
	var stmtSeq []Statement
	// 左算术表达式语句序列，左算术表达式结果变量
	lExpStmtSeq, lExpResultID := g.generateExp(compExp.LExp)
	// 右算术表达式语句序列，右算术表达式结果变量
	rExpStmtSeq, rExpResultID := g.generateExp(compExp.RExp)
	// 将左算术表达式语句序列和右算术表达式语句序列添加到语句序列中
	stmtSeq = append(stmtSeq, lExpStmtSeq...)
	stmtSeq = append(stmtSeq, rExpStmtSeq...)
	// 比较表达式结果变量
	resultID := g.NewAnonymousVar()
	// 按照运算符生成语句序列
	switch compExp.Op {
	case ast.EQUAL:
		// = 等于
		// 1 左右不等，跳转到 4
		// 2 结果为1
		// 3 跳转到 5
		// 4 结果为0
		// 5 （判断后语句）
		stmtSeq = append(stmtSeq, *NewStatement(JNEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s == %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s == %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("equal: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s == %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.DIAMOND:
		// <> 不等于
		// 1 左右相等，跳转到 4
		// 2 结果为0
		// 3 跳转到 5
		// 4 结果为1
		// 5 （判断后语句）
		stmtSeq = append(stmtSeq, *NewStatement(JEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s != %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s != %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("notEqual: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s != %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.GREATER:
		// > 大于
		// 1 左小于等于右，跳转到 4
		// 2 结果为0
		// 3 跳转到 5
		// 4 结果为1
		// 5 （判断后语句）
		stmtSeq = append(stmtSeq, *NewStatement(JLESSEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s > %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s > %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("greater: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s > %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.GREATEREQUAL:
		// >= 大于等于
		// 1 左小于右，跳转到 4
		// 2 结果为0
		// 3 跳转到 5
		// 4 结果为1
		// 5 （判断后语句）
		stmtSeq = append(stmtSeq, *NewStatement(JLESS, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s >= %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s >= %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("greaterEqual: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s >= %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.LESS:
		// < 小于
		// 1 左大于等于右，跳转到 4
		// 2 结果为0
		// 3 跳转到 5
		// 4 结果为1
		// 5 （判断后语句）
		stmtSeq = append(stmtSeq, *NewStatement(JGREATEQUAL, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s < %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s < %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("less: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s < %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	case ast.LESSEQUAL:
		// <= 小于等于
		// 1 左大于右，跳转到 4
		// 2 结果为0
		// 3 跳转到 5
		// 4 结果为1
		// 5 （判断后语句）
		stmtSeq = append(stmtSeq, *NewStatement(JGREAT, StrParam(hir.VarToStr(lExpResultID)), StrParam(hir.VarToStr(rExpResultID)), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s <= %s false: goto here+3", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID))))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s <= %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("lessEqual: goto here+2")))
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s <= %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	}

	return stmtSeq, resultID
}

// generateRelationalExp 生成关系表达式
func (g *MIRGenerator) generateRelationalExp(relationalExp hir.RelationExp) ([]Statement, int) {
	// 只有左比较表达式，没有右比较表达式，直接返回左比较表达式
	if relationalExp.Op == ast.EMPTY {
		return g.generateCompExp(relationalExp.LExp)
	}

	// 语句序列
	var stmtSeq []Statement
	// 由于短路运算，需要先计算左表达式
	lExpStmtSeq, lExpResultID := g.generateCompExp(relationalExp.LExp)
	stmtSeq = append(stmtSeq, lExpStmtSeq...)

	// 结果变量
	resultID := g.NewAnonymousVar()
	// 1 左非0，跳转到 4
	// 2 结果为0
	// 3 跳转到右+5
	// 4 （右表达式判断）
	stmtSeq = append(stmtSeq, *NewStatement(JNZERO, StrParam(hir.VarToStr(lExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s true: goto here+3", hir.VarToStr(lExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s false: %s = 0", hir.VarToStr(lExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 5)), fmt.Sprintf("goto here+len(rExpStmtSeq)+5)")))

	// 若左表达式为1，判断右表达式
	rExpStmtSeq, rExpResultID := g.generateCompExp(relationalExp.RExp)
	stmtSeq = append(stmtSeq, rExpStmtSeq...)

	// 右+1 右为0,跳转到右+4
	// 右+2 结果为1
	// 右+3 跳转到右+5
	// 右+4 结果为0
	// 右+5 （判断后语句）
	stmtSeq = append(stmtSeq, *NewStatement(JZERO, StrParam(hir.VarToStr(rExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s false: goto here+3", hir.VarToStr(rExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s true: %s = 1", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("goto here+2")))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s false: %s = 0", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))

	// 语句序列偏移量
	offset := len(rExpStmtSeq)
	stmtSeq[len(lExpStmtSeq)+4].Res = StrParam(fmt.Sprintf("_T_JMP_REF_%d", offset+5))

	return stmtSeq, resultID
}

// generateConditionalExp 生成条件表达式
func (g *MIRGenerator) generateConditionalExp(conditionalExp hir.ConditionalExp) ([]Statement, int) {
	// 只有左比较表达式，没有右比较表达式，直接返回左比较表达式
	if conditionalExp.Op == ast.EMPTY {
		return g.generateRelationalExp(conditionalExp.LExp)
	}

	// 语句序列
	var stmtSeq []Statement
	// 由于短路运算，需要先计算左表达式
	lExpStmtSeq, lExpResultID := g.generateRelationalExp(conditionalExp.LExp)
	stmtSeq = append(stmtSeq, lExpStmtSeq...)
	resultID := g.NewAnonymousVar()
	// 1 左为0，跳转到 4
	// 2 结果为0
	// 3 跳转到右+5
	// 4 （右表达式判断）
	stmtSeq = append(stmtSeq, *NewStatement(JZERO, StrParam(hir.VarToStr(lExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s false: goto here+3", hir.VarToStr(lExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s true: %s = 1", hir.VarToStr(lExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 5)), fmt.Sprintf("goto here+len(rExpStmtSeq)+5)")))

	// 若左表达式为0，判断右表达式
	rExpStmtSeq, rExpResultID := g.generateRelationalExp(conditionalExp.RExp)
	stmtSeq = append(stmtSeq, rExpStmtSeq...)
	// 右+1 右为0,跳转到右+4
	// 右+2 结果为1
	// 右+3 跳转到右+5
	// 右+4 结果为0
	// 右+5 （判断后语句）
	stmtSeq = append(stmtSeq, *NewStatement(JZERO, StrParam(hir.VarToStr(rExpResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 3)), fmt.Sprintf("if %s false: goto here+3", hir.VarToStr(rExpResultID))))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(1), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s true: %s = 1", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", 2)), fmt.Sprintf("goto here+2")))
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(resultID)), IntParam(0), StrParam(hir.VarToStr(resultID)), fmt.Sprintf("if %s false: %s = 0", hir.VarToStr(rExpResultID), hir.VarToStr(resultID))))

	// 语句序列偏移量
	offset := len(rExpStmtSeq)
	stmtSeq[len(lExpStmtSeq)+4].Res = StrParam(fmt.Sprintf("_T_JMP_REF_%d", offset+5))
	return stmtSeq, resultID
}
