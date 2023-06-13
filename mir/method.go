package mir

import (
	"CompilerInGo/hir"
)

// generateMethod 生成方法
func (g *MIRGenerator) generateMethod(method hir.Method) ([]Statement, []int, int) {
	// 语句序列
	var stmtSeq []Statement
	// 参数变量
	var paramIDs []int
	// 返回值变量
	returnVar := g.NewAnonymousVar()

	// 解析参数
	for _, param := range method.Params {
		// 形参局部变量声明语句
		stmt := g.NewLocalVariableDeclaration(param.Type, param.ID)
		stmtSeq = append(stmtSeq, *stmt)
		// 添加到形参列表
		paramIDs = append(paramIDs, hir.StrToVar(stmt.Res.Str()))
	}

	// 添加到方法列表
	methodInfo := g.Methods[method.Name]
	methodInfo.ReturnVar = returnVar
	methodInfo.ActParams = paramIDs
	g.Methods[method.Name] = methodInfo

	// 解析方法体
	stmtSeq = append(stmtSeq, g.generateStatement(*method.Body)...)

	// 使用注释标注方法开始位置
	stmtSeq[0].Comment = stmtSeq[0].Comment + " # method: " + method.Name

	return stmtSeq, paramIDs, returnVar
}
