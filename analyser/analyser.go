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

type Analyser struct {
	methods       *symbol.SymbolTable[hir.Method]
	scope         *symbol.SymbolTable[ast.Type]
	unusedVars    *symbol.SymbolTable[hir.ID]
	unusedMethods *symbol.SymbolTable[hir.ID]
	methodIn      ast.Method
}

func NewAnalyser() *Analyser {
	return &Analyser{
		methods:       symbol.NewSymbolTable[hir.Method](),
		scope:         symbol.NewSymbolTable[ast.Type](),
		unusedVars:    symbol.NewSymbolTable[hir.ID](),
		unusedMethods: symbol.NewSymbolTable[hir.ID](),
	}
}

func (a *Analyser) ScopeInit() {
	a.scope = symbol.NewSymbolTable[ast.Type]()
	a.unusedVars = symbol.NewSymbolTable[hir.ID]()
}

func (a *Analyser) Analyse(AST *ast.Program) (*hir.Program, int) {

	errs := 0

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

	for _, method := range AST.Method {
		resMethod, err := a.analyseMethod(method)
		if err != nil {
			_ = glg.Error(err)
			errs++
		} else {
			a.methods.AddSymbol(a.methodIn.GetMethodName(), *resMethod)
			a.unusedMethods.AddSymbol(a.methodIn.GetMethodName(), hir.ID(resMethod.Name))
		}
	}

	if !a.methods.HasSymbol("main") {
		_ = glg.Error("no entrypoint for program: no valid main method")
		errs++
	}

	if a.unusedMethods.Size() > 0 {
		for _, unusedMethod := range a.unusedMethods.Symbols {
			_ = glg.Warnf("Analyser: unused method %s", unusedMethod)
		}
	}

	return hir.NewProgram(a.methods.ToArray()), errs
}

func (a *Analyser) analyseMethod(method ast.Method) (*hir.Method, error) {
	a.methodIn = method
	if a.methods.HasSymbol(a.methodIn.GetMethodName()) {
		return nil, errors.New(fmt.Sprintf("method name %s is duplicated with another method", a.methodIn.GetMethodName()))
	}

	a.ScopeInit()
	paramsSeq, _ := method.ParamList.Integrate()
	a.scope.AddSymbol(a.methodIn.GetMethodName(), ast.Type{})
	a.methods.AddSymbol(a.methodIn.GetMethodName(), hir.Method{})

	if paramsSeq != nil {
		for _, param := range paramsSeq.Seq {
			if a.scope.HasSymbol(param.ID.Literal.(string)) {
				a.methods.RemoveSymbol(a.methodIn.GetMethodName())
				return nil, errors.New(fmt.Sprintf("param name %s is duplicated", param.ID.Literal.(string)))
			}
			a.scope.AddSymbol(param.ID.Literal.(string), param.Type)
		}
	}

	stmts, err := a.analyseBlock(method.Block)
	if err != nil {
		a.methods.RemoveSymbol(a.methodIn.GetMethodName())
		return nil, err
	}

	if a.unusedVars.Size() > 0 {
		for _, unusedVar := range a.unusedVars.Symbols {
			_ = glg.Warnf("Analyser: unused variable %s in method %s", unusedVar, a.methodIn.GetMethodName())
		}
	}

	resultType := hir.AstResultType(method.ResultType)
	paramList := hir.AstParamList(method.ParamList)

	resMethod := hir.NewMethod(resultType.ToHIR(), a.methodIn.GetMethodName(), paramList.ToHIR(), &stmts)

	return resMethod, nil
}

func (a *Analyser) analyseBlock(block ast.Block) (hir.Statement, error) {
	if block.Statements == nil || len(*block.Statements) == 0 {
		return nil, nil
	}

	stmts := make([]*hir.Statement, 0)
	for _, stmt := range *block.Statements {
		resStmt, err := a.analyseStmt(stmt)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, resStmt)
	}

	return hir.NewBlock(stmts), nil
}

func (a *Analyser) analyseStmt(stmts ast.Statement) (*hir.Statement, error) {
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
			return nil, nil
		}
		return nil, errors.New("unknown statement type")
	}
}

func (a *Analyser) analyseConditionalStmt(statement ast.ConditionalStatement) (hir.Statement, error) {
	condExp, err := a.analyseConditionalExp(statement.ConditionalExp)
	if err != nil {
		return nil, err
	}

	ifStmt, err := a.analyseStmt(statement.Statement)
	if err != nil {
		return nil, err
	}

	if statement.ElseStatement != nil {
		elseStmt, err := a.analyseStmt(*statement.ElseStatement)
		if err != nil {
			return nil, err
		}
		return hir.NewConditionalStatement(*condExp, ifStmt, elseStmt), nil
	}
	return hir.NewConditionalStatement(*condExp, ifStmt, nil), nil
}

func (a *Analyser) analyseLoopStmt(statement ast.LoopStatement) (hir.Statement, error) {
	condExp, err := a.analyseConditionalExp(statement.ConditionalExp)
	if err != nil {
		return nil, err
	}

	whileStmt, err := a.analyseStmt(statement.Statement)
	if err != nil {
		return nil, err
	}

	return hir.NewLoopStatement(*condExp, whileStmt), nil
}

func (a *Analyser) analyseCallStmt(statement ast.CallStatement) (hir.Statement, error) {
	if !a.methods.HasSymbol(statement.ID.Literal.(string)) {
		return nil, errors.New(fmt.Sprintf("method %s is not defined", statement.ID.Literal.(string)))
	}

	targetMethod, _ := a.methods.GetSymbol(statement.ID.Literal.(string))
	methodParams := targetMethod.Params
	actParams, _ := statement.ActParamList.Integrate()

	resExps := make([]hir.Exp, 0)
	for _, exp := range actParams {
		resExp, err := a.analyseExp(exp)
		if err != nil {
			return nil, err
		}
		resExps = append(resExps, *resExp)
	}

	if (methodParams == nil && len(actParams) != 0) || (methodParams != nil && len(actParams) == 0) || len(actParams) != len(methodParams) {
		return nil, errors.New(fmt.Sprintf("method %s is called with wrong number of parameters", statement.ID.Literal.(string)))
	}

	a.unusedMethods.RemoveSymbol(statement.ID.Literal.(string))

	return hir.NewCallStatement(statement.ID.Literal.(string), resExps), nil
}

func (a *Analyser) analyseAssignmentStmt(statement ast.AssignmentStatement) (hir.Statement, error) {
	if !a.scope.HasSymbol(statement.ID.Literal.(string)) {
		return nil, errors.New(fmt.Sprintf("variable %s is not defined in method %s", statement.ID.Literal.(string), a.methodIn.GetMethodName()))
	}

	resExp, err := a.analyseExp(statement.Exp)
	if err != nil {
		return nil, err
	}

	a.unusedVars.RemoveSymbol(statement.ID.Literal.(string))

	return hir.NewAssignStatement(statement.ID.Literal.(string), *resExp), nil
}

func (a *Analyser) analyseReturnStmt(statement ast.ReturnStatement) (hir.Statement, error) {
	if statement.Exp == nil {
		return nil, nil
	}
	resExp, err := a.analyseExp(*statement.Exp)
	if err != nil {
		return nil, err
	}

	return hir.NewReturnStatement(*resExp), nil
}

func (a *Analyser) analyseBreakStmt(statement ast.BreakStatement) (hir.Statement, error) {
	return hir.NewBreakStatement(), nil
}

func (a *Analyser) analyseContinueStmt(statement ast.ContinueStatement) (hir.Statement, error) {
	return hir.NewContinueStatement(), nil
}

func (a *Analyser) analyseLocalVarDecl(declaration ast.LocalVariableDeclaration) (hir.Statement, error) {
	decls, _ := declaration.Integrate()
	declTypeIDArray := make([]hir.TypeIDPair, 0)
	for _, decl := range decls.Seq {
		if a.scope.HasSymbol(decl.ID.Literal.(string)) {
			return nil, errors.New(fmt.Sprintf("variable %s is duplicated in method %s", decl.ID.Literal.(string), a.methodIn.GetMethodName()))
		}
		a.scope.AddSymbol(decl.ID.Literal.(string), decl.Type)
		a.unusedVars.AddSymbol(decl.ID.Literal.(string), hir.ID(decl.ID.Literal.(string)))

		declHIR := hir.AstTypeIDPair(decl).ToHIR()
		declTypeIDArray = append(declTypeIDArray, *declHIR)
	}

	return hir.NewLocalVariableDeclaration(declTypeIDArray), nil
}

func (a *Analyser) analyseConditionalExp(exp ast.ConditionalExp) (*hir.ConditionalExp, error) {
	relationExps := exp.Integrate()
	resRelationExps := make([]hir.RelationExp, 0)
	for _, relationExp := range relationExps {
		resRelationExp, err := a.analyseRelationExp(relationExp)
		if err != nil {
			return nil, err
		}
		resRelationExps = append(resRelationExps, *resRelationExp)
	}

	if len(resRelationExps) == 1 {
		condExp := hir.NewConditionalExp(&resRelationExps[0], nil)
		return &condExp, nil
	}
	condExp := hir.NewConditionalExp(&resRelationExps[0], &resRelationExps[1])
	return &condExp, nil
}

func (a *Analyser) analyseRelationExp(exp ast.RelationExp) (*hir.RelationExp, error) {
	compExps := exp.Integrate()

	resCompExps := make([]hir.CompExp, 0)
	for _, compExp := range compExps {
		resCompExp, err := a.analyseCompExp(compExp)
		if err != nil {
			return nil, err
		}
		resCompExps = append(resCompExps, *resCompExp)
	}

	if len(resCompExps) == 1 {
		relationExp := hir.NewRelationExp(&resCompExps[0], nil)
		return &relationExp, nil
	}
	relationExp := hir.NewRelationExp(&resCompExps[0], &resCompExps[1])
	return &relationExp, nil
}

func (a *Analyser) analyseCompExp(exp ast.CompExp) (*hir.CompExp, error) {
	exps := exp.Integrate()
	resExps := make([]hir.Exp, 0)
	for _, expElem := range exps {
		resExp, err := a.analyseExp(expElem)
		if err != nil {
			return nil, err
		}
		resExps = append(resExps, *resExp)
	}

	if len(resExps) == 1 {
		compExp := hir.NewCompExp(&resExps[0], ast.EMPTY, nil)
		return &compExp, nil
	}
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
		return nil, errors.New(fmt.Sprintf("unknown CmpOp %s", exp.CmpOp.Literal))
	}
}

func (a *Analyser) analyseExp(exp ast.Exp) (*hir.Exp, error) {
	terms := exp.Integrate()
	resTerms := make([]hir.Term, 0)
	for _, term := range terms {
		resTerm, err := a.analyseTerm(term)
		if err != nil {
			return nil, err
		}
		resTerms = append(resTerms, *resTerm)
	}

	if len(resTerms) == 1 {
		exp := hir.NewExp(&resTerms[0], ast.EMPTY, nil)
		return &exp, nil
	}

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

func (a *Analyser) analyseTerm(term ast.Term) (*hir.Term, error) {
	factors := term.Integrate()
	resFactors := make([]hir.Factor, 0)
	for _, factor := range factors {
		resFactor, err := a.analyseFactor(factor)
		if err != nil {
			return nil, err
		}
		resFactors = append(resFactors, resFactor)
	}
	if len(resFactors) == 1 {
		resTerm := hir.NewTerm(&resFactors[0], ast.EMPTY, nil)
		return &resTerm, nil
	}
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

func (a *Analyser) analyseFactor(factor ast.Factor) (hir.Factor, error) {
	switch factor.Factor.(type) {
	case ast.FactorTuple:
		return a.analyseExp(*factor.Factor.(ast.FactorTuple).Exp)
	case lexer.Token:
		if factor.Factor.(lexer.Token).Type == lexer.IDENTIFIER {
			if a.methods.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return nil, errors.New(fmt.Sprintf("%s is a method, but used as a variable", factor.Factor.(lexer.Token).Literal.(string)))
			}
			if !a.scope.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return nil, errors.New(fmt.Sprintf("variable %s is not defined in method %s", factor.Factor.(lexer.Token).Literal.(string), a.methodIn.GetMethodName()))
			} else {
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
