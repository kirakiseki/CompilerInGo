package hir

import (
	"CompilerInGo/lexer"
	"CompilerInGo/parser/ast"
)

type AstType ast.Type
type AstResultType ast.ResultType
type AstTypeIDPair ast.TypeIDPair
type AstParamList ast.ParamList

func (typ AstType) ToHIR() Type {
	switch typ.Type {
	case lexer.INT:
		return TInteger
	case lexer.FLOAT:
		return TFloat
	case lexer.STRING:
		return TString
	case lexer.CHAR:
		return TChar
	default:
		return TErr
	}
}

func (resultType *AstResultType) ToHIR() ResultType {
	switch resultType.Type {
	case lexer.INT:
		return TInteger
	case lexer.FLOAT:
		return TFloat
	case lexer.STRING:
		return TString
	case lexer.CHAR:
		return TChar
	case lexer.VOID:
		return TVoid
	default:
		return TErr
	}
}

func (typeIDPair AstTypeIDPair) ToHIR() *TypeIDPair {
	return NewTypeIDPair(AstType(typeIDPair.Type).ToHIR(), typeIDPair.ID.Literal.(string))
}

func (paramList *AstParamList) ToHIR() []*TypeIDPair {
	list := ast.ParamList(*paramList)
	paramsSeq, _ := list.Integrate()
	params := make([]*TypeIDPair, 0)
	if paramsSeq == nil {
		return params
	}
	for _, param := range paramsSeq.Seq {
		params = append(params, AstTypeIDPair(param).ToHIR())
	}
	return params
}
