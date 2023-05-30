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
	methods *symbol.SymbolTable[ast.Method]
}

func NewAnalyser() *Analyser {
	return &Analyser{
		methods: symbol.NewSymbolTable[ast.Method](),
	}
}

func (a *Analyser) Analyse(AST *ast.Program) {
	if AST == nil {
		_ = glg.Error("AST is nil")
	}
	if AST.Method == nil || len(AST.Method) == 0 {
		_ = glg.Error("AST.Methods is nil/empty")
	}

	for _, method := range AST.Method {
		err := a.analyseMethod(method)
		if err != nil {
			_ = glg.Error(err)
		}
	}
}

func (a *Analyser) analyseMethod(method ast.Method) error {
	if a.methods.HasSymbol(method.ID.Literal.(string)) {
		return errors.New(fmt.Sprintf("method name %s is duplicated with another method", method.ID.Literal.(string)))
	}

	paramsTable := symbol.NewSymbolTable[ast.Type]()
	paramsSeq, _ := method.ParamList.Integrate()
	paramsTable.AddSymbol(method.ID.Literal.(string), ast.Type{})

	for _, param := range paramsSeq.Seq {
		if paramsTable.HasSymbol(param.ID.Literal.(string)) {
			return errors.New(fmt.Sprintf("param name %s is duplicated", param.ID.Literal.(string)))
		}
		paramsTable.AddSymbol(param.ID.Literal.(string), param.Type)
	}

	err := a.analyseBlock(method.Block, paramsTable)
	if err != nil {
		_ = glg.Error(err)
	}

	a.methods.AddSymbol(method.ID.Literal.(string), method)
	return nil
}

func (a *Analyser) analyseBlock(block ast.Block, params *symbol.SymbolTable[ast.Type]) error {
	if block.Statements == nil || len(*block.Statements) == 0 {
		return nil
	}
	for _, stmts := range *block.Statements {
		err := a.analyseStmt(stmts, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseStmt(stmts ast.Statement, params *symbol.SymbolTable[ast.Type]) error {
	switch stmts.Type {
	case ast.CONDITIONALSTATEMENT:
		return a.analyseConditionalStmt(*((stmts.Statement).(*ast.ConditionalStatement)), params)
	case ast.LOOPSTATEMENT:
		return a.analyseLoopStmt(*((stmts.Statement).(*ast.LoopStatement)), params)
	case ast.CALLSTATEMENT:
		return a.analyseCallStmt(*((stmts.Statement).(*ast.CallStatement)), params)
	case ast.ASSIGNMENTSTATEMENT:
		return a.analyseAssignmentStmt(*((stmts.Statement).(*ast.AssignmentStatement)), params)
	case ast.RETURNSTATEMENT:
		return a.analyseReturnStmt(*((stmts.Statement).(*ast.ReturnStatement)), params)
	case ast.BREAKSTATEMENT:
		return a.analyseBreakStmt(*((stmts.Statement).(*ast.BreakStatement)), params)
	case ast.CONTINUESTATEMENT:
		return a.analyseContinueStmt(*((stmts.Statement).(*ast.ContinueStatement)), params)
	case ast.LOCALVARIABLEDECLARATION:
		return a.analyseLocalVarDecl(*((stmts.Statement).(*ast.LocalVariableDeclaration)), params)
	case ast.BLOCK:
		return a.analyseBlock(*(stmts.Statement.(*ast.Block)), params)
	default:
		if stmts == (ast.Statement{}) {
			return nil
		}
		return errors.New("unknown statement type")
	}
}

func (a *Analyser) analyseConditionalStmt(statement ast.ConditionalStatement, params *symbol.SymbolTable[ast.Type]) error {
	err := a.analyseConditionalExp(statement.ConditionalExp, params)
	if err != nil {
		return err
	}

	err = a.analyseStmt(statement.Statement, params)
	if err != nil {
		return err
	}

	if statement.ElseStatement != nil {
		err = a.analyseStmt(*statement.ElseStatement, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Analyser) analyseLoopStmt(statement ast.LoopStatement, params *symbol.SymbolTable[ast.Type]) error {
	err := a.analyseConditionalExp(statement.ConditionalExp, params)
	if err != nil {
		return err
	}

	err = a.analyseStmt(statement.Statement, params)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyser) analyseCallStmt(statement ast.CallStatement, params *symbol.SymbolTable[ast.Type]) error {
	if !a.methods.HasSymbol(statement.ID.Literal.(string)) {
		return errors.New(fmt.Sprintf("method %s is not defined", statement.ID.Literal.(string)))
	}

	method, _ := a.methods.GetSymbol(statement.ID.Literal.(string))
	methodParams, _ := method.ParamList.Integrate()
	actParams, _ := statement.ActParamList.Integrate()

	for _, exp := range actParams {
		err := a.analyseExp(exp, params)
		if err != nil {
			return err
		}
	}

	if len(actParams) != len(methodParams.Seq) {
		return errors.New(fmt.Sprintf("method %s is called with wrong number of parameters", statement.ID.Literal.(string)))
	}

	return nil
}

func (a *Analyser) analyseAssignmentStmt(statement ast.AssignmentStatement, params *symbol.SymbolTable[ast.Type]) error {
	if !params.HasSymbol(statement.ID.Literal.(string)) {
		return errors.New(fmt.Sprintf("variable %s is not defined", statement.ID.Literal.(string)))
	}

	err := a.analyseExp(statement.Exp, params)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyser) analyseReturnStmt(statement ast.ReturnStatement, params *symbol.SymbolTable[ast.Type]) error {
	if statement.Exp == nil {
		return nil
	}
	err := a.analyseExp(*statement.Exp, params)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyser) analyseBreakStmt(statement ast.BreakStatement, params *symbol.SymbolTable[ast.Type]) error {
	return nil
}

func (a *Analyser) analyseContinueStmt(statement ast.ContinueStatement, params *symbol.SymbolTable[ast.Type]) error {
	return nil
}

func (a *Analyser) analyseLocalVarDecl(declaration ast.LocalVariableDeclaration, params *symbol.SymbolTable[ast.Type]) error {
	decls, _ := declaration.Integrate()
	for _, decl := range decls.Seq {
		if params.HasSymbol(decl.ID.Literal.(string)) {
			return errors.New(fmt.Sprintf("variable %s is duplicated", decl.ID.Literal.(string)))
		}
		params.AddSymbol(decl.ID.Literal.(string), decl.Type)
	}
	return nil
}

func (a *Analyser) analyseConditionalExp(exp ast.ConditionalExp, params *symbol.SymbolTable[ast.Type]) error {
	relationExps := exp.Integrate()
	for _, relationExp := range relationExps {
		err := a.analyseRelationExp(relationExp, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseRelationExp(exp ast.RelationExp, params *symbol.SymbolTable[ast.Type]) error {
	compExps := exp.Integrate()
	for _, compExp := range compExps {
		err := a.analyseCompExp(compExp, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseCompExp(exp ast.CompExp, params *symbol.SymbolTable[ast.Type]) error {
	exps := exp.Integrate()
	for _, expElem := range exps {
		err := a.analyseExp(expElem, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseExp(exp ast.Exp, params *symbol.SymbolTable[ast.Type]) error {
	terms := exp.Integrate()
	for _, term := range terms {
		err := a.analyseTerm(term, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseTerm(term ast.Term, params *symbol.SymbolTable[ast.Type]) error {
	factors := term.Integrate()
	for _, factor := range factors {
		err := a.analyseFactor(factor, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyser) analyseFactor(factor ast.Factor, params *symbol.SymbolTable[ast.Type]) error {
	switch factor.Factor.(type) {
	case ast.FactorTuple:
		return a.analyseExp(*factor.Factor.(ast.FactorTuple).Exp, params)
	case lexer.Token:
		if factor.Factor.(lexer.Token).Type == lexer.IDENTIFIER {
			if a.methods.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return errors.New(fmt.Sprintf("%s is a method, but used as a variable", factor.Factor.(lexer.Token).Literal.(string)))
			}
			if !params.HasSymbol(factor.Factor.(lexer.Token).Literal.(string)) {
				return errors.New(fmt.Sprintf("variable %s is not defined", factor.Factor.(lexer.Token).Literal.(string)))
			}
		}
		return nil
	default:
		return nil
	}
}
