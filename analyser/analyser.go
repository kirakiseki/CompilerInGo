package analyser

import (
	"CompilerInGo/analyser/symbol"
	"CompilerInGo/lexer"
	"CompilerInGo/parser/ast"
	"errors"
	"fmt"
	"github.com/kpango/glg"
)

type Analyser struct {
	methods       *symbol.SymbolTable[ast.Method]
	scope         *symbol.SymbolTable[ast.Type]
	unusedVars    *symbol.SymbolTable[string]
	unusedMethods *symbol.SymbolTable[ast.Method]
	methodName    string
}

func NewAnalyser() *Analyser {
	return &Analyser{
		methods:       symbol.NewSymbolTable[ast.Method](),
		scope:         symbol.NewSymbolTable[ast.Type](),
		unusedVars:    symbol.NewSymbolTable[string](),
		unusedMethods: symbol.NewSymbolTable[ast.Method](),
	}
}

func (a *Analyser) ScopeInit() {
	a.scope = symbol.NewSymbolTable[ast.Type]()
	a.unusedVars = symbol.NewSymbolTable[string]()
}

func (a *Analyser) Analyse(AST *ast.Program) int {

	errs := 0

	if AST == nil {
		_ = glg.Error("AST is nil")
		errs++
		return errs
	}
	if AST.Method == nil || len(AST.Method) == 0 {
		_ = glg.Error("AST.Methods is nil/empty")
		errs++
		return errs
	}

	for _, method := range AST.Method {
		err := a.analyseMethod(method)
		if err != nil {
			_ = glg.Error(err)
			errs++
		}
	}

	if !a.methods.HasSymbol("main") {
		_ = glg.Error("no entrypoint for program: no valid main method")
		errs++
	}

	a.unusedMethods.RemoveSymbol("main")
	if a.unusedMethods.Size() > 0 {
		for _, unusedMethod := range a.unusedMethods.Symbols {
			_ = glg.Warnf("Analyser: unused method %s", unusedMethod.ID.Literal.(string))
		}
	}

	return errs
}

func (a *Analyser) analyseMethod(method ast.Method) error {
	a.methodName = method.ID.Literal.(string)
	if a.methods.HasSymbol(a.methodName) {
		return errors.New(fmt.Sprintf("method name %s is duplicated with another method", a.methodName))
	}

	a.ScopeInit()
	paramsSeq, _ := method.ParamList.Integrate()
	a.scope.AddSymbol(a.methodName, ast.Type{})

	for _, param := range paramsSeq.Seq {
		if a.scope.HasSymbol(param.ID.Literal.(string)) {
			return errors.New(fmt.Sprintf("param name %s is duplicated", param.ID.Literal.(string)))
		}
		a.scope.AddSymbol(param.ID.Literal.(string), param.Type)
	}

	err := a.analyseBlock(method.Block)
	if err != nil {
		return err
	}

	a.methods.AddSymbol(a.methodName, method)
	a.unusedMethods.AddSymbol(a.methodName, method)

	if a.unusedVars.Size() > 0 {
		for _, unusedVar := range a.unusedVars.Symbols {
			_ = glg.Warnf("Analyser: unused variable %s in method %s", unusedVar, a.methodName)
		}
	}

	return nil
}

func (a *Analyser) analyseBlock(block ast.Block) error {
	if block.Statements == nil || len(*block.Statements) == 0 {
		return nil
	}
	for _, stmts := range *block.Statements {
		err := a.analyseStmt(stmts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseStmt(stmts ast.Statement) error {
	switch stmts.Type {
	case ast.CONDITIONALSTATEMENT:
		return a.analyseConditionalStmt(*((stmts.Statement).(*ast.ConditionalStatement)))
	case ast.LOOPSTATEMENT:
		return a.analyseLoopStmt(*((stmts.Statement).(*ast.LoopStatement)))
	case ast.CALLSTATEMENT:
		return a.analyseCallStmt(*((stmts.Statement).(*ast.CallStatement)))
	case ast.ASSIGNMENTSTATEMENT:
		return a.analyseAssignmentStmt(*((stmts.Statement).(*ast.AssignmentStatement)))
	case ast.RETURNSTATEMENT:
		return a.analyseReturnStmt(*((stmts.Statement).(*ast.ReturnStatement)))
	case ast.BREAKSTATEMENT:
		return a.analyseBreakStmt(*((stmts.Statement).(*ast.BreakStatement)))
	case ast.CONTINUESTATEMENT:
		return a.analyseContinueStmt(*((stmts.Statement).(*ast.ContinueStatement)))
	case ast.LOCALVARIABLEDECLARATION:
		return a.analyseLocalVarDecl(*((stmts.Statement).(*ast.LocalVariableDeclaration)))
	case ast.BLOCK:
		return a.analyseBlock(*(stmts.Statement.(*ast.Block)))
	default:
		if stmts == (ast.Statement{}) {
			return nil
		}
		return errors.New("unknown statement type")
	}
}

func (a *Analyser) analyseConditionalStmt(statement ast.ConditionalStatement) error {
	err := a.analyseConditionalExp(statement.ConditionalExp)
	if err != nil {
		return err
	}

	err = a.analyseStmt(statement.Statement)
	if err != nil {
		return err
	}

	if statement.ElseStatement != nil {
		err = a.analyseStmt(*statement.ElseStatement)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Analyser) analyseLoopStmt(statement ast.LoopStatement) error {
	err := a.analyseConditionalExp(statement.ConditionalExp)
	if err != nil {
		return err
	}

	err = a.analyseStmt(statement.Statement)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyser) analyseCallStmt(statement ast.CallStatement) error {
	if !a.methods.HasSymbol(statement.ID.Literal.(string)) {
		return errors.New(fmt.Sprintf("method %s is not defined", statement.ID.Literal.(string)))
	}

	method, _ := a.methods.GetSymbol(statement.ID.Literal.(string))
	methodParams, _ := method.ParamList.Integrate()
	actParams, _ := statement.ActParamList.Integrate()

	for _, exp := range actParams {
		err := a.analyseExp(exp)
		if err != nil {
			return err
		}
	}

	if len(actParams) != len(methodParams.Seq) {
		return errors.New(fmt.Sprintf("method %s is called with wrong number of parameters", statement.ID.Literal.(string)))
	}

	a.unusedMethods.RemoveSymbol(statement.ID.Literal.(string))

	return nil
}

func (a *Analyser) analyseAssignmentStmt(statement ast.AssignmentStatement) error {
	if !a.scope.HasSymbol(statement.ID.Literal.(string)) {
		return errors.New(fmt.Sprintf("variable %s is not defined in method %s", statement.ID.Literal.(string), a.methodName))
	}

	err := a.analyseExp(statement.Exp)
	if err != nil {
		return err
	}

	a.unusedVars.RemoveSymbol(statement.ID.Literal.(string))

	return nil
}

func (a *Analyser) analyseReturnStmt(statement ast.ReturnStatement) error {
	if statement.Exp == nil {
		return nil
	}
	err := a.analyseExp(*statement.Exp)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyser) analyseBreakStmt(statement ast.BreakStatement) error {
	return nil
}

func (a *Analyser) analyseContinueStmt(statement ast.ContinueStatement) error {
	return nil
}

func (a *Analyser) analyseLocalVarDecl(declaration ast.LocalVariableDeclaration) error {
	decls, _ := declaration.Integrate()
	for _, decl := range decls.Seq {
		if a.scope.HasSymbol(decl.ID.Literal.(string)) {
			return errors.New(fmt.Sprintf("variable %s is duplicated in method %s", decl.ID.Literal.(string), a.methodName))
		}
		a.scope.AddSymbol(decl.ID.Literal.(string), decl.Type)
		a.unusedVars.AddSymbol(decl.ID.Literal.(string), decl.ID.Literal.(string))
	}
	return nil
}

func (a *Analyser) analyseConditionalExp(exp ast.ConditionalExp) error {
	relationExps := exp.Integrate()
	for _, relationExp := range relationExps {
		err := a.analyseRelationExp(relationExp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseRelationExp(exp ast.RelationExp) error {
	compExps := exp.Integrate()
	for _, compExp := range compExps {
		err := a.analyseCompExp(compExp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseCompExp(exp ast.CompExp) error {
	exps := exp.Integrate()
	for _, expElem := range exps {
		err := a.analyseExp(expElem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseExp(exp ast.Exp) error {
	terms := exp.Integrate()
	for _, term := range terms {
		err := a.analyseTerm(term)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseTerm(term ast.Term) error {
	factors := term.Integrate()
	for _, factor := range factors {
		err := a.analyseFactor(factor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseFactor(factor ast.Factor) error {
	switch factor.Factor.(type) {
	case ast.FactorTuple:
		return a.analyseExp(*factor.Factor.(ast.FactorTuple).Exp)
	case lexer.Token:
		if factor.Factor.(lexer.Token).Type == lexer.IDENTIFIER {
			if a.methods.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return errors.New(fmt.Sprintf("%s is a method, but used as a variable", factor.Factor.(lexer.Token).Literal.(string)))
			}
			if !a.scope.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return errors.New(fmt.Sprintf("variable %s is not defined in method %s", factor.Factor.(lexer.Token).Literal.(string), a.methodName))
			} else {
				a.unusedVars.RemoveSymbol(factor.Factor.(lexer.Token).Literal.(string))
			}
		}

		return nil
	default:
		return nil
	}
}
