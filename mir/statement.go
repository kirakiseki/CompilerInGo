package mir

import (
	"CompilerInGo/hir"
	"fmt"
)

func (g *MIRGenerator) generateStatement(stmt hir.Statement) []Statement {
	switch stmt.(type) {
	case hir.ConditionalStatement:
		return g.generateConditionalStatement(stmt.(hir.ConditionalStatement))
	//case hir.LoopStatement:
	//	return g.generateLoopStatement(stmt.(hir.LoopStatement))
	case hir.CallStatement:
		return g.generateCallStatement(stmt.(hir.CallStatement))
	case hir.AssignStatement:
		return g.generateAssignStatement(stmt.(hir.AssignStatement))
	case hir.ReturnStatement:
		return g.generateReturnStatement(stmt.(hir.ReturnStatement))
	//case hir.BreakStatement:
	//	return g.generateBreakStatement(stmt.(hir.BreakStatement))
	//case hir.ContinueStatement:
	//	return g.generateContinueStatement(stmt.(hir.ContinueStatement))
	case hir.LocalVariableDeclaration:
		return g.generateLocalVariableDeclaration(stmt.(hir.LocalVariableDeclaration))
	case hir.Block:
		return g.generateBlock(stmt.(hir.Block))
	default:
		//glg.Fatal("Unknown statement type")
	}
	return nil
}

func (g *MIRGenerator) generateBlock(block hir.Block) []Statement {
	var stmtSeq []Statement
	for _, stmt := range block.Statements {
		stmtSeq = append(stmtSeq, g.generateStatement(*stmt)...)
	}
	return stmtSeq
}

func (g *MIRGenerator) generateLocalVariableDeclaration(stmt hir.LocalVariableDeclaration) []Statement {
	var stmtSeq []Statement
	for _, pair := range stmt.TypeIDPair {
		stmtSeq = append(stmtSeq, *g.NewLocalVariableDeclaration(pair.Type, pair.ID))
	}
	return stmtSeq
}

func (g *MIRGenerator) NewLocalVariableDeclaration(t hir.Type, id hir.ID) *Statement {
	varID := g.NewVar(string(id))
	return NewStatement(ASSIGN, StrParam(hir.VarToStr(varID)), StrParam(id), StrParam(hir.VarToStr(varID)), fmt.Sprintf("%s = %s", hir.VarToStr(varID), id))
}

func (g *MIRGenerator) generateAssignStatement(stmt hir.AssignStatement) []Statement {
	var stmtSeq []Statement
	varID := g.GetVar(stmt.Target)
	expStmtSeq, expResultID := g.generateExp(stmt.Exp)
	stmtSeq = append(stmtSeq, expStmtSeq...)
	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(varID)), StrParam(hir.VarToStr(expResultID)), StrParam(hir.VarToStr(varID)), fmt.Sprintf("%s = %s", hir.VarToStr(varID), hir.VarToStr(expResultID))))
	return stmtSeq
}

func (g *MIRGenerator) generateReturnStatement(stmt hir.ReturnStatement) []Statement {
	stmtSeq, expResultID := g.generateExp(stmt.Exp)
	g.Context = g.CtxStack.Top()
	g.CtxStack.Pop()
	if g.Context.MethodIn.Name != "main" {
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(hir.VarToStr(g.Methods[g.Context.MethodIn.Name].ReturnVar)), fmt.Sprintf("method %s return value %s : goto %s", g.Context.MethodIn.Name, hir.VarToStr(expResultID), hir.VarToStr(g.Methods[g.Context.MethodIn.Name].ReturnVar))))
	} else {
		stmtSeq = append(stmtSeq, *NewStatement(STOP, StrParam("_"), StrParam("_"), StrParam(hir.VarToStr(expResultID)), fmt.Sprintf("main return value: %s : STOP", hir.VarToStr(expResultID))))
	}
	return stmtSeq
}

func (g *MIRGenerator) generateCallStatement(stmt hir.CallStatement) []Statement {
	var stmtSeq []Statement
	g.Context = Context{
		MethodIn: MethodInfo{Name: stmt.Method},
	}
	g.CtxStack.Push(g.Context)

	var actParams []int
	for _, exp := range stmt.ActParam {
		expStmtSeq, expResultID := g.generateExp(exp)
		stmtSeq = append(stmtSeq, expStmtSeq...)
		actParams = append(actParams, expResultID)
	}

	method, ok := g.Methods[stmt.Method]
	if !ok {
		g.NewMethod(stmt.Method)
		method = g.Methods[stmt.Method]
		methodStmtSeq, formalParams, returnLabel := g.generateMethod(*g.HIRProgram.GetMethod(stmt.Method))
		method.Pos = len(g.MethodSeq)
		g.MethodSeq = append(g.MethodSeq, methodStmtSeq...)
		method.ActParams = formalParams
		method.ReturnVar = returnLabel
		g.Methods[stmt.Method] = method
	}

	for i := 0; i < len(actParams); i++ {
		stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(method.ActParams[i])), StrParam(hir.VarToStr(actParams[i])), StrParam(hir.VarToStr(method.ActParams[i])), fmt.Sprintf("call param: %s = %s", hir.VarToStr(method.ActParams[i]), hir.VarToStr(actParams[i]))))
	}

	stmtSeq = append(stmtSeq, *NewStatement(ASSIGN, StrParam(hir.VarToStr(method.ReturnVar)), StrParam("_T_HERE_TO_JMP+1"), StrParam(hir.VarToStr(method.ReturnVar)), fmt.Sprintf("call returnTo: %s", hir.VarToStr(method.ReturnVar))))
	stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_METHOD_%s", stmt.Method)), fmt.Sprintf("call method: %s", stmt.Method)))

	return stmtSeq
}

func (g *MIRGenerator) generateConditionalStatement(stmt hir.ConditionalStatement) []Statement {
	var stmtSeq []Statement

	expStmtSeq, expResultID := g.generateConditionalExp(stmt.Condition)
	stmtSeq = append(stmtSeq, expStmtSeq...)
	expFalseStmt := *NewStatement(JZERO, StrParam(hir.VarToStr(expResultID)), StrParam("_"), StrParam("_"), "_")

	trueSeq := g.generateStatement(*stmt.IfBody)
	expFalseStmt.Res = StrParam(fmt.Sprintf("_T_JMP_REF_%d", len(trueSeq)+2))
	expFalseStmt.Comment = fmt.Sprintf("if false: goto here+%d", len(trueSeq)+2)
	stmtSeq = append(stmtSeq, expFalseStmt)

	trueSeq[0].Comment = fmt.Sprintf("true block: %s", trueSeq[0].Comment)
	stmtSeq = append(stmtSeq, trueSeq...)
	if stmt.ElseBody != nil {
		falseSeq := g.generateStatement(*stmt.ElseBody)
		falseSeq[0].Comment = fmt.Sprintf("false block: %s", falseSeq[0].Comment)
		stmtSeq = append(stmtSeq, *NewStatement(JMP, StrParam("_"), StrParam("_"), StrParam(fmt.Sprintf("_T_JMP_REF_%d", len(falseSeq)+1)), fmt.Sprintf("goto here+%d", len(falseSeq)+1)))
		stmtSeq = append(stmtSeq, falseSeq...)
	}
	return stmtSeq
}
