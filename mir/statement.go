package mir

import (
	"CompilerInGo/hir"
	"fmt"
	"github.com/kpango/glg"
)

// generateStatement 生成语句
func (g *MIRGenerator) generateStatement(stmt hir.Statement) []Statement {
	// 按语句类型分别处理
	switch stmt.(type) {
	case hir.ConditionalStatement:
		return g.generateConditionalStatement(stmt.(hir.ConditionalStatement))
	case hir.LoopStatement:
		return g.generateLoopStatement(stmt.(hir.LoopStatement))
	case hir.CallStatement:
		return g.generateCallStatement(stmt.(hir.CallStatement))
	case hir.AssignStatement:
		return g.generateAssignStatement(stmt.(hir.AssignStatement))
	case hir.ReturnStatement:
		return g.generateReturnStatement(stmt.(hir.ReturnStatement))
	case hir.LocalVariableDeclaration:
		return g.generateLocalVariableDeclaration(stmt.(hir.LocalVariableDeclaration))
	case hir.BreakStatement:
		return g.generateBreakStatement(stmt.(hir.BreakStatement))
	case hir.ContinueStatement:
		return g.generateContinueStatement(stmt.(hir.ContinueStatement))
	case hir.Block:
		return g.generateBlock(stmt.(hir.Block))
	default:
		glg.Fatal("Unknown statement type")
	}
	return nil
}

// generateBlock 生成语句块
func (g *MIRGenerator) generateBlock(block hir.Block) []Statement {
	var stmtSeq []Statement
	for _, stmt := range block.Statements {
		// 遍历语句块中的语句
		stmtSeq = append(stmtSeq, g.generateStatement(*stmt)...)
	}
	return stmtSeq
}

// generateLocalVariableDeclaration 生成局部变量声明语句
func (g *MIRGenerator) generateLocalVariableDeclaration(stmt hir.LocalVariableDeclaration) []Statement {
	var stmtSeq []Statement
	for _, pair := range stmt.TypeIDPair {
		// 遍历局部变量声明语句中的变量
		stmtSeq = append(stmtSeq, *g.NewLocalVariableDeclaration(pair.Type, pair.ID))
	}
	return stmtSeq
}

// NewLocalVariableDeclaration 新局部变量声明语句
func (g *MIRGenerator) NewLocalVariableDeclaration(t hir.Type, id hir.ID) *Statement {
	// 定义新变量
	varID := g.NewVar(string(id))
	// 赋值语句
	return NewStatement(ASSIGN, StrParam(hir.VarToStr(varID)), StrParam(id), StrParam(hir.VarToStr(varID)), fmt.Sprintf("%s = %s", hir.VarToStr(varID), id))
}

// generateAssignStatement 生成赋值语句
func (g *MIRGenerator) generateAssignStatement(stmt hir.AssignStatement) []Statement {
	var stmtSeq []Statement

	// 待赋值的变量
	varID := g.GetVar(stmt.Target)

	// 解析表达式语句和表达式值的结果变量
	expStmtSeq, expResultID := g.generateExp(stmt.Exp)
	stmtSeq = append(stmtSeq, expStmtSeq...)
	// 表达式结果变量赋值给待赋值的变量
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(varID)), StrParam(hir.VarToStr(expResultID)), StrParam(hir.VarToStr(varID)), fmt.Sprintf("%s = %s", hir.VarToStr(varID), hir.VarToStr(expResultID))))
	return stmtSeq
}

// generateReturnStatement 生成返回语句
func (g *MIRGenerator) generateReturnStatement(stmt hir.ReturnStatement) []Statement {
	// 解析表达式语句和表达式值的结果变量
	stmtSeq, expResultID := g.generateExp(stmt.Exp)

	// 退出当前方法，恢复上下文
	g.Context = g.CtxStack.Top()
	g.CtxStack.Pop()

	// 若为main方法，则生成STOP语句
	if g.Context.MethodIn.Name != "main" {
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(hir.VarToStr(g.Methods[g.Context.MethodIn.Name].ReturnVar)), fmt.Sprintf("method %s return value %s : goto %s", g.Context.MethodIn.Name, hir.VarToStr(expResultID), hir.VarToStr(g.Methods[g.Context.MethodIn.Name].ReturnVar))))
	} else {
		stmtSeq = append(stmtSeq, *NewStatement(STOP, StrParam("_"), StrParam("_"), StrParam(hir.VarToStr(expResultID)), fmt.Sprintf("main return value: %s : STOP", hir.VarToStr(expResultID))))
	}
	return stmtSeq
}

// generateCallStatement 生成调用语句
func (g *MIRGenerator) generateCallStatement(stmt hir.CallStatement) []Statement {
	var stmtSeq []Statement

	// 解析新的函数，生成新的上下文
	g.Context = Context{
		MethodIn:      MethodInfo{Name: stmt.Method},
		LoopCondLabel: -1,
		LoopEndLabel:  -1,
	}
	g.CtxStack.Push(g.Context)

	// 解析实参
	var actParams []int
	for _, exp := range stmt.ActParam {
		// 将表达式的结果赋值给实参，并加入到实参列表中
		expStmtSeq, expResultID := g.generateExp(exp)
		stmtSeq = append(stmtSeq, expStmtSeq...)
		actParams = append(actParams, expResultID)
	}

	// 若方法未定义，则生成新方法
	method, ok := g.Methods[stmt.Method]
	if !ok {
		// 生成新方法
		g.NewMethod(stmt.Method)
		method = g.Methods[stmt.Method]
		methodStmtSeq, formalParams, returnLabel := g.generateMethod(*g.HIRProgram.GetMethod(stmt.Method))
		method.Pos = len(g.MethodSeq)
		g.MethodSeq = append(g.MethodSeq, methodStmtSeq...)
		method.ActParams = formalParams
		method.ReturnVar = returnLabel
		g.Methods[stmt.Method] = method
	}

	// 将实参赋值给形参
	for i := 0; i < len(actParams); i++ {
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(method.ActParams[i])), StrParam(hir.VarToStr(actParams[i])), StrParam(hir.VarToStr(method.ActParams[i])), fmt.Sprintf("call param: %s = %s", hir.VarToStr(method.ActParams[i]), hir.VarToStr(actParams[i]))))
	}

	// 生成执行完后跳转位置的标签
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(method.ReturnVar)), StrParam("_T_HERE_TO_JMP+1"), StrParam(hir.VarToStr(method.ReturnVar)), fmt.Sprintf("call returnTo: %s", hir.VarToStr(method.ReturnVar))))
	// 生成跳转语句
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_METHOD_%s", stmt.Method)), fmt.Sprintf("call method: %s", stmt.Method)))

	return stmtSeq
}

// generateConditionalStatement 生成条件语句
func (g *MIRGenerator) generateConditionalStatement(stmt hir.ConditionalStatement) []Statement {
	var stmtSeq []Statement

	// 解析条件语句中的条件表达式
	expStmtSeq, expResultID := g.generateConditionalExp(stmt.Condition)
	stmtSeq = append(stmtSeq, expStmtSeq...)

	// 若条件为假，则跳转到else语句块/if语句块结束
	expFalseStmt := *NewStatement(JZERO, StrParam(hir.VarToStr(expResultID)), StrParam("_"), StrParam("_"), "_")

	// 若条件为真，则执行if语句块
	trueSeq := g.generateStatement(*stmt.IfBody)
	// 设置跳转语句的跳转位置
	expFalseStmt.Res = StrParam(fmt.Sprintf("_T_JMP_REF_%d", len(trueSeq)+2))
	expFalseStmt.Comment = fmt.Sprintf("if condition false: goto here+%d", len(trueSeq)+2)
	stmtSeq = append(stmtSeq, expFalseStmt)

	// 添加if语句块
	trueSeq[0].Comment = fmt.Sprintf("true block: %s", trueSeq[0].Comment)
	stmtSeq = append(stmtSeq, trueSeq...)

	// 添加else语句块
	if stmt.ElseBody != nil {
		falseSeq := g.generateStatement(*stmt.ElseBody)
		falseSeq[0].Comment = fmt.Sprintf("false block: %s", falseSeq[0].Comment)
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", len(falseSeq)+1)), fmt.Sprintf("goto here+%d", len(falseSeq)+1)))
		stmtSeq = append(stmtSeq, falseSeq...)
	}

	return stmtSeq
}

// generateLoopStatement 生成循环语句
func (g *MIRGenerator) generateLoopStatement(stmt hir.LoopStatement) []Statement {
	var stmtSeq []Statement

	// 解析循环语句中的条件表达式
	expStmtSeq, expResultID := g.generateConditionalExp(stmt.Condition)
	stmtSeq = append(stmtSeq, expStmtSeq...)

	// 循环体
	bodySeq := g.generateStatement(*stmt.Body)

	// 跳转到循环结束
	skipLoopStmt := *NewStatement(JZERO, StrParam(hir.VarToStr(expResultID)), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", len(bodySeq)+2)), fmt.Sprintf("while condition %s false : skip loop: goto here+%d", hir.VarToStr(expResultID), len(bodySeq)+2))
	// 跳转到条件表达式判断
	nextLoopStmt := *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", -(len(bodySeq)+len(expStmtSeq)+1))), fmt.Sprintf("next loop: goto here+%d", -(len(bodySeq)+len(expStmtSeq)+1)))

	// 语句拼接
	stmtSeq = append(stmtSeq, skipLoopStmt)
	stmtSeq = append(stmtSeq, bodySeq...)
	stmtSeq = append(stmtSeq, nextLoopStmt)

	return stmtSeq
}

// genraateBreakStatement 生成break语句
func (g *MIRGenerator) generateBreakStatement(stmt hir.BreakStatement) []Statement {
	var stmtSeq []Statement
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam("_"), "_T_BREAK"))
	return stmtSeq
}

// generateContinueStatement 生成continue语句
func (g *MIRGenerator) generateContinueStatement(stmt hir.ContinueStatement) []Statement {
	var stmtSeq []Statement
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam("_"), "_T_CONTINUE"))
	return stmtSeq
}
