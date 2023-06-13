package analyser

import (
	"CompilerInGo/analyser/symbol"
	"CompilerInGo/hir"
	"CompilerInGo/lexer"
	"CompilerInGo/parser/ast"
	"errors"
	"fmt"
	"github.com/kpango/glg"
)

// Analyser 语义分析器
type Analyser struct {
	methods       *symbol.SymbolTable[hir.Method] // 方法表
	scope         *symbol.SymbolTable[ast.Type]   // 作用域内的变量表
	unusedVars    *symbol.SymbolTable[hir.ID]     // 未使用的变量表
	unusedMethods *symbol.SymbolTable[hir.ID]     // 未使用的方法表
	methodIn      ast.Method                      // 当前分析的方法
}

// NewAnalyser 新建语义分析器
func NewAnalyser() *Analyser {
	return &Analyser{
		methods:       symbol.NewSymbolTable[hir.Method](),
		scope:         symbol.NewSymbolTable[ast.Type](),
		unusedVars:    symbol.NewSymbolTable[hir.ID](),
		unusedMethods: symbol.NewSymbolTable[hir.ID](),
	}
}

// ScopeInit 对新的作用域进行初始化
func (a *Analyser) ScopeInit() {
	a.scope = symbol.NewSymbolTable[ast.Type]()
	a.unusedVars = symbol.NewSymbolTable[hir.ID]()
}

// Analyse 对AST进行语义分析
func (a *Analyser) Analyse(AST *ast.Program) (*hir.Program, int) {
	// 错误计数
	errs := 0

	// 检查AST是否为空
	if AST == nil {
		_ = glg.Error("AST is nil")
		errs++
		return nil, errs
	}
	if AST.Method == nil || len(AST.Method) == 0 {
		_ = glg.Error("AST.Methods is nil/empty")
		errs++
		return nil, errs
	}

	//遍历分析AST中的每个方法
	for _, method := range AST.Method {
		// 在子程序中进行分析
		resMethod, err := a.analyseMethod(method)
		if err != nil {
			_ = glg.Error(err)
			errs++
		} else {
			// 分析成功，将方法添加到方法表、未使用的寒暑表中
			a.methods.AddSymbol(a.methodIn.GetMethodName(), *resMethod)
			a.unusedMethods.AddSymbol(a.methodIn.GetMethodName(), hir.ID(resMethod.Name))
		}
	}

	// 检查是否有main方法
	if !a.methods.HasSymbol("main") {
		_ = glg.Error("no entrypoint for program: no valid main method")
		errs++
	}

	// 除main方法外，检查是否有未使用的方法
	a.unusedMethods.RemoveSymbol("main")
	if a.unusedMethods.Size() > 0 {
		for _, unusedMethod := range a.unusedMethods.Symbols {
			_ = glg.Warnf("Analyser: unused method %s", unusedMethod)
		}
	}

	// 返回HIR和错误计数
	return hir.NewProgram(a.methods.ToArray()), errs
}

// analyseMethod 对方法进行语义分析
func (a *Analyser) analyseMethod(method ast.Method) (*hir.Method, error) {
	// 切换当前分析的方法
	a.methodIn = method
	// 检查方法名是否重复
	if a.methods.HasSymbol(a.methodIn.GetMethodName()) {
		return nil, errors.New(fmt.Sprintf("method name %s is duplicated with another method", a.methodIn.GetMethodName()))
	}

	// 初始化作用域
	a.ScopeInit()

	// 分析方法的参数并添加到作用域中
	paramsSeq, _ := method.ParamList.Integrate()
	a.scope.AddSymbol(a.methodIn.GetMethodName(), ast.Type{})

	// 将方法添加到方法表中
	a.methods.AddSymbol(a.methodIn.GetMethodName(), hir.Method{})

	// 分析参数表
	if paramsSeq != nil {
		for _, param := range paramsSeq.Seq {
			// 作用域中已经存在同名的变量
			if a.scope.HasSymbol(param.ID.Literal.(string)) {
				a.methods.RemoveSymbol(a.methodIn.GetMethodName())
				return nil, errors.New(fmt.Sprintf("param name %s is duplicated", param.ID.Literal.(string)))
			}
			// 将参数添加到作用域中
			a.scope.AddSymbol(param.ID.Literal.(string), param.Type)
		}
	}

	// 分析方法体
	stmts, err := a.analyseBlock(method.Block)
	if err != nil {
		// 分析失败，从方法表中移除方法
		a.methods.RemoveSymbol(a.methodIn.GetMethodName())
		return nil, err
	}

	// 检查是否有未使用的变量
	if a.unusedVars.Size() > 0 {
		for _, unusedVar := range a.unusedVars.Symbols {
			_ = glg.Warnf("Analyser: unused variable %s in method %s", unusedVar, a.methodIn.GetMethodName())
		}
	}

	// 整合为HIR的方法
	resultType := hir.AstResultType(method.ResultType)
	paramList := hir.AstParamList(method.ParamList)

	resMethod := hir.NewMethod(resultType.ToHIR(), a.methodIn.GetMethodName(), paramList.ToHIR(), &stmts)

	return resMethod, nil
}

// analyseBlock 对块进行语义分析
func (a *Analyser) analyseBlock(block ast.Block) (hir.Statement, error) {
	// 检查块是否为空
	if block.Statements == nil || len(*block.Statements) == 0 {
		return nil, nil
	}

	// 分析块中的每个语句
	stmts := make([]*hir.Statement, 0)
	for _, stmt := range *block.Statements {
		// 在子程序中分析每条语句
		resStmt, err := a.analyseStmt(stmt)
		if err != nil {
			return nil, err
		}
		// 将分析结果添加到块中
		stmts = append(stmts, resStmt)
	}

	return hir.NewBlock(stmts), nil
}

// analyseStmt 对语句进行语义分析
func (a *Analyser) analyseStmt(stmts ast.Statement) (*hir.Statement, error) {
	// 根据语句类型，调用相应的子程序
	switch stmts.Type {
	case ast.CONDITIONALSTATEMENT:
		condStmt, err := a.analyseConditionalStmt(*((stmts.Statement).(*ast.ConditionalStatement)))
		return &condStmt, err
	case ast.LOOPSTATEMENT:
		loopStmt, err := a.analyseLoopStmt(*((stmts.Statement).(*ast.LoopStatement)))
		return &loopStmt, err
	case ast.CALLSTATEMENT:
		callStmt, err := a.analyseCallStmt(*((stmts.Statement).(*ast.CallStatement)))
		return &callStmt, err
	case ast.ASSIGNMENTSTATEMENT:
		assignStmt, err := a.analyseAssignmentStmt(*((stmts.Statement).(*ast.AssignmentStatement)))
		return &assignStmt, err
	case ast.RETURNSTATEMENT:
		returnStmt, err := a.analyseReturnStmt(*((stmts.Statement).(*ast.ReturnStatement)))
		return &returnStmt, err
	case ast.BREAKSTATEMENT:
		breakStmt, err := a.analyseBreakStmt(*((stmts.Statement).(*ast.BreakStatement)))
		return &breakStmt, err
	case ast.CONTINUESTATEMENT:
		contStmt, err := a.analyseContinueStmt(*((stmts.Statement).(*ast.ContinueStatement)))
		return &contStmt, err
	case ast.LOCALVARIABLEDECLARATION:
		varDeclStmt, err := a.analyseLocalVarDecl(*((stmts.Statement).(*ast.LocalVariableDeclaration)))
		return &varDeclStmt, err
	case ast.BLOCK:
		blockStmt, err := a.analyseBlock(*(stmts.Statement.(*ast.Block)))
		return &blockStmt, err
	default:
		if stmts == (ast.Statement{}) {
			// 空语句
			return nil, nil
		}
		// 未知语句类型
		return nil, errors.New("unknown statement type")
	}
}

// analyseConditionalStmt 对条件语句进行语义分析
func (a *Analyser) analyseConditionalStmt(statement ast.ConditionalStatement) (hir.Statement, error) {
	// 分析条件表达式
	condExp, err := a.analyseConditionalExp(statement.ConditionalExp)
	if err != nil {
		return nil, err
	}

	// 分析if语句
	ifStmt, err := a.analyseStmt(statement.Statement)
	if err != nil {
		return nil, err
	}

	// 存在else语句
	if statement.ElseStatement != nil {
		// 分析else语句
		elseStmt, err := a.analyseStmt(*statement.ElseStatement)
		if err != nil {
			return nil, err
		}
		//返回if-else语句
		return hir.NewConditionalStatement(*condExp, ifStmt, elseStmt), nil
	}

	// 返回if语句
	return hir.NewConditionalStatement(*condExp, ifStmt, nil), nil
}

// analyseLoopStmt 对循环语句进行语义分析
func (a *Analyser) analyseLoopStmt(statement ast.LoopStatement) (hir.Statement, error) {
	// 分析条件表达式
	condExp, err := a.analyseConditionalExp(statement.ConditionalExp)
	if err != nil {
		return nil, err
	}

	// 分析循环体
	whileStmt, err := a.analyseStmt(statement.Statement)
	if err != nil {
		return nil, err
	}

	return hir.NewLoopStatement(*condExp, whileStmt), nil
}

// analyseCallStmt 对调用语句进行语义分析
func (a *Analyser) analyseCallStmt(statement ast.CallStatement) (hir.Statement, error) {
	// 检查方法是否存在
	if !a.methods.HasSymbol(statement.ID.Literal.(string)) {
		return nil, errors.New(fmt.Sprintf("method %s is not defined", statement.ID.Literal.(string)))
	}

	// 获取方法、参数列表
	targetMethod, _ := a.methods.GetSymbol(statement.ID.Literal.(string))
	methodParams := targetMethod.Params

	// 获取实参列表
	actParams, _ := statement.ActParamList.Integrate()

	// 不能调用main方法
	if targetMethod.Name == "main" {
		return nil, errors.New("main method is not callable")
	}

	// 分析实参列表
	resExps := make([]hir.Exp, 0)
	for _, exp := range actParams {
		// 分析表达式
		resExp, err := a.analyseExp(exp)
		if err != nil {
			return nil, err
		}
		resExps = append(resExps, *resExp)
	}

	// 实参与形参个数不匹配
	if (methodParams == nil && len(actParams) != 0) || (methodParams != nil && len(methodParams) != 0 && len(actParams) == 0) || len(actParams) != len(methodParams) {
		return nil, errors.New(fmt.Sprintf("method %s is called with wrong number of parameters", statement.ID.Literal.(string)))
	}

	// 被调用，不再是未使用方法
	a.unusedMethods.RemoveSymbol(statement.ID.Literal.(string))

	return hir.NewCallStatement(statement.ID.Literal.(string), resExps), nil
}

// analyseAssignmentStmt 对赋值语句进行语义分析
func (a *Analyser) analyseAssignmentStmt(statement ast.AssignmentStatement) (hir.Statement, error) {
	// 作用域中是否存在变量
	if !a.scope.HasSymbol(statement.ID.Literal.(string)) {
		return nil, errors.New(fmt.Sprintf("variable %s is not defined in method %s", statement.ID.Literal.(string), a.methodIn.GetMethodName()))
	}

	// 分析表达式
	resExp, err := a.analyseExp(statement.Exp)
	if err != nil {
		return nil, err
	}

	// 被赋值，不再是未使用变量
	a.unusedVars.RemoveSymbol(statement.ID.Literal.(string))

	return hir.NewAssignStatement(statement.ID.Literal.(string), *resExp), nil
}

// analyseReturnStmt 对返回语句进行语义分析
func (a *Analyser) analyseReturnStmt(statement ast.ReturnStatement) (hir.Statement, error) {
	// 返回值为空
	if statement.Exp == nil {
		return nil, nil
	}

	// 分析表达式
	resExp, err := a.analyseExp(*statement.Exp)
	if err != nil {
		return nil, err
	}

	return hir.NewReturnStatement(*resExp), nil
}

// analyseBreakStmt 对break语句进行语义分析
func (a *Analyser) analyseBreakStmt(statement ast.BreakStatement) (hir.Statement, error) {
	return hir.NewBreakStatement(), nil
}

// analyseContinueStmt 对continue语句进行语义分析
func (a *Analyser) analyseContinueStmt(statement ast.ContinueStatement) (hir.Statement, error) {
	return hir.NewContinueStatement(), nil
}

// analyseExpStmt 对表达式语句进行语义分析
func (a *Analyser) analyseLocalVarDecl(declaration ast.LocalVariableDeclaration) (hir.Statement, error) {
	// 分析变量声明 type-ID对
	decls, _ := declaration.Integrate()
	declTypeIDArray := make([]hir.TypeIDPair, 0)

	for _, decl := range decls.Seq {
		// 变量重复声明
		if a.scope.HasSymbol(decl.ID.Literal.(string)) {
			return nil, errors.New(fmt.Sprintf("variable %s is duplicated in method %s", decl.ID.Literal.(string), a.methodIn.GetMethodName()))
		}

		// 添加到作用域
		a.scope.AddSymbol(decl.ID.Literal.(string), decl.Type)
		// 添加到未使用变量
		a.unusedVars.AddSymbol(decl.ID.Literal.(string), hir.ID(decl.ID.Literal.(string)))

		// 转换为HIR
		declHIR := hir.AstTypeIDPair(decl).ToHIR()
		declTypeIDArray = append(declTypeIDArray, *declHIR)
	}

	return hir.NewLocalVariableDeclaration(declTypeIDArray), nil
}

// analyseConditionalExp 对条件表达式进行语义分析
func (a *Analyser) analyseConditionalExp(exp ast.ConditionalExp) (*hir.ConditionalExp, error) {
	// 分析左右关系表达式
	relationExps := exp.Integrate()

	resRelationExps := make([]hir.RelationExp, 0)
	for _, relationExp := range relationExps {
		// 遍历分析左右关系表达式
		resRelationExp, err := a.analyseRelationExp(relationExp)
		if err != nil {
			return nil, err
		}
		resRelationExps = append(resRelationExps, *resRelationExp)
	}

	// 只有一个关系表达式（左侧）
	if len(resRelationExps) == 1 {
		condExp := hir.NewConditionalExp(&resRelationExps[0], nil)
		return &condExp, nil
	}
	// 有两个关系表达式（左右）
	condExp := hir.NewConditionalExp(&resRelationExps[0], &resRelationExps[1])
	return &condExp, nil
}

// analyseRelationExp 对关系表达式进行语义分析
func (a *Analyser) analyseRelationExp(exp ast.RelationExp) (*hir.RelationExp, error) {
	// 分析左右比较表达式
	compExps := exp.Integrate()

	resCompExps := make([]hir.CompExp, 0)
	for _, compExp := range compExps {
		// 遍历分析左右比较表达式
		resCompExp, err := a.analyseCompExp(compExp)
		if err != nil {
			return nil, err
		}
		resCompExps = append(resCompExps, *resCompExp)
	}

	// 只有一个比较表达式（左侧）
	if len(resCompExps) == 1 {
		relationExp := hir.NewRelationExp(&resCompExps[0], nil)
		return &relationExp, nil
	}
	// 有两个比较表达式（左右）
	relationExp := hir.NewRelationExp(&resCompExps[0], &resCompExps[1])
	return &relationExp, nil
}

// analyseCompExp 对比较表达式进行语义分析
func (a *Analyser) analyseCompExp(exp ast.CompExp) (*hir.CompExp, error) {
	// 分析左右算术表达式
	exps := exp.Integrate()

	resExps := make([]hir.Exp, 0)
	for _, expElem := range exps {
		// 遍历分析左右算术表达式
		resExp, err := a.analyseExp(expElem)
		if err != nil {
			return nil, err
		}
		resExps = append(resExps, *resExp)
	}

	// 只有一个算术表达式（左侧）
	if len(resExps) == 1 {
		compExp := hir.NewCompExp(&resExps[0], ast.EMPTY, nil)
		return &compExp, nil
	}
	// 有两个算术表达式（左右）
	// 按照比较运算符类型构造比较表达式
	switch exp.CmpOp.Type {
	case lexer.LESS:
		compExp := hir.NewCompExp(&resExps[0], ast.LESS, &resExps[1])
		return &compExp, nil
	case lexer.LESSEQUAL:
		compExp := hir.NewCompExp(&resExps[0], ast.LESSEQUAL, &resExps[1])
		return &compExp, nil
	case lexer.GREATER:
		compExp := hir.NewCompExp(&resExps[0], ast.GREATER, &resExps[1])
		return &compExp, nil
	case lexer.GREATEREQUAL:
		compExp := hir.NewCompExp(&resExps[0], ast.GREATEREQUAL, &resExps[1])
		return &compExp, nil
	case lexer.EQUAL:
		compExp := hir.NewCompExp(&resExps[0], ast.EQUAL, &resExps[1])
		return &compExp, nil
	case lexer.DIAMOND:
		compExp := hir.NewCompExp(&resExps[0], ast.DIAMOND, &resExps[1])
		return &compExp, nil
	default:
		// 未知比较运算符
		return nil, errors.New(fmt.Sprintf("unknown CmpOp %s", exp.CmpOp.Literal))
	}
}

// analyseExp 对算术表达式进行语义分析
func (a *Analyser) analyseExp(exp ast.Exp) (*hir.Exp, error) {
	// 分析左右项
	terms := exp.Integrate()

	resTerms := make([]hir.Term, 0)
	for _, term := range terms {
		// 遍历分析左右项
		resTerm, err := a.analyseTerm(term)
		if err != nil {
			return nil, err
		}
		resTerms = append(resTerms, *resTerm)
	}

	// 只有一个项（左侧）
	if len(resTerms) == 1 {
		exp := hir.NewExp(&resTerms[0], ast.EMPTY, nil)
		return &exp, nil
	}

	// 有两个项（左右）
	// 按照加减运算符类型构造算术表达式
	switch exp.ExpRest.PlusOrMinus.Type {
	case lexer.PLUS:
		exp := hir.NewExp(&resTerms[0], ast.PLUS, &resTerms[1])
		return &exp, nil
	case lexer.MINUS:
		exp := hir.NewExp(&resTerms[0], ast.MINUS, &resTerms[1])
		return &exp, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown PlusOrMinus %s", exp.ExpRest.PlusOrMinus.Literal))
	}
}

// analyseTerm 对项进行语义分析
func (a *Analyser) analyseTerm(term ast.Term) (*hir.Term, error) {
	// 分析左右因子
	factors := term.Integrate()

	resFactors := make([]hir.Factor, 0)
	for _, factor := range factors {
		// 遍历分析左右因子
		resFactor, err := a.analyseFactor(factor)
		if err != nil {
			return nil, err
		}
		resFactors = append(resFactors, resFactor)
	}

	// 只有一个因子（左侧）
	if len(resFactors) == 1 {
		resTerm := hir.NewTerm(&resFactors[0], ast.EMPTY, nil)
		return &resTerm, nil
	}

	// 有两个因子（左右）
	// 按照乘除运算符类型构造项
	switch term.TermRest.MulOrDiv.Type {
	case lexer.TIMES:
		term := hir.NewTerm(&resFactors[0], ast.TIMES, &resFactors[1])
		return &term, nil
	case lexer.DIVIDE:
		term := hir.NewTerm(&resFactors[0], ast.DIVIDE, &resFactors[1])
		return &term, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown MulOrDiv %s", term.TermRest.MulOrDiv.Literal))
	}
}

// analyseFactor 对因子进行语义分析
func (a *Analyser) analyseFactor(factor ast.Factor) (hir.Factor, error) {
	// 按照因子类型进行分析
	switch factor.Factor.(type) {
	case ast.FactorTuple:
		// (Exp)
		return a.analyseExp(*factor.Factor.(ast.FactorTuple).Exp)
	case lexer.Token:
		// ID| INTC | DECI
		if factor.Factor.(lexer.Token).Type == lexer.IDENTIFIER {
			// 如果是方法名，报错
			if a.methods.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return nil, errors.New(fmt.Sprintf("%s is a method, but used as a variable", factor.Factor.(lexer.Token).Literal.(string)))
			}

			// 如果是未定义的变量，报错
			if !a.scope.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return nil, errors.New(fmt.Sprintf("variable %s is not defined in method %s", factor.Factor.(lexer.Token).Literal.(string), a.methodIn.GetMethodName()))
			} else {
				// 使用了变量，从未使用变量列表中删除
				a.unusedVars.RemoveSymbol(factor.Factor.(lexer.Token).Literal.(string))
			}
			ID := hir.ID(factor.Factor.(lexer.Token).Literal.(string))
			return ID, nil
		} else if factor.Factor.(lexer.Token).Type == lexer.INTEGER_LITERAL {
			return hir.NewInteger(factor.Factor.(lexer.Token).Literal.(int64)), nil
		} else if factor.Factor.(lexer.Token).Type == lexer.DECIMAL_LITERAL {
			return hir.NewFloat(factor.Factor.(lexer.Token).Literal.(float64)), nil
		} else {
			return nil, errors.New(fmt.Sprintf("unknown factor %s", factor.Factor.(lexer.Token).Literal))
		}
	default:
		return nil, errors.New(fmt.Sprintf("unknown factor %s", factor.Factor))
	}
}
