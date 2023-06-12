package mir

import (
	"CompilerInGo/hir"
)

func (g *MIRGenerator) generateMethod(method hir.Method) ([]Statement, []int, int) {
	var stmtSeq []Statement
	var paramIDs []int
	returnLabel := g.NewAnonymousVar()

	for _, param := range method.Params {
		stmt := g.NewLocalVariableDeclaration(param.Type, param.ID)
		stmtSeq = append(stmtSeq, *stmt)
		paramIDs = append(paramIDs, hir.StrToVar(stmt.Res.Str()))
	}

	methodInfo := g.Methods[method.Name]
	methodInfo.ReturnVar = returnLabel
	methodInfo.ActParams = paramIDs
	g.Methods[method.Name] = methodInfo

	stmtSeq = append(stmtSeq, g.generateStatement(*method.Body)...)
	stmtSeq[0].Comment = stmtSeq[0].Comment + " # method: " + method.Name
	return stmtSeq, paramIDs, returnLabel
}
