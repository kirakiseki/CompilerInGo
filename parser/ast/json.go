package ast

import (
	"CompilerInGo/lexer"
	"encoding/json"
)

// 为了方便输出AST，我们需要为AST中的每个结构体实现MarshalJSON()方法
// 输出指定内容，略去位置，null值，以及一些不必要的信息

type TypeIDPair struct {
	Type Type
	ID   ID
}

func (resultType *ResultType) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Literal  any
		Type     string
		Category string
	}{
		Literal:  resultType.Literal,
		Type:     lexer.TokenTypeString[resultType.Type],
		Category: lexer.TokenCategoryString[resultType.Category],
	})
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Literal  any
		Type     string
		Category string
	}{
		Literal:  id.Literal,
		Type:     lexer.TokenTypeString[id.Type],
		Category: lexer.TokenCategoryString[id.Category],
	})
}

func (typ Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Literal  any
		Type     string
		Category string
	}{
		Literal:  typ.Literal,
		Type:     lexer.TokenTypeString[typ.Type],
		Category: lexer.TokenCategoryString[typ.Category],
	})
}

func (paramList *ParamList) MarshalJSON() ([]byte, error) {
	if paramList.ParamList == nil {
		return json.Marshal([]struct{}{})
	}

	tuple := make([]TypeIDPair, 0)
	tuple = append(tuple, TypeIDPair{
		Type: paramList.ParamList.Type,
		ID:   paramList.ParamList.ID,
	})

	for _, elem := range *paramList.ParamList.ParamListRest {
		tuple = append(tuple, TypeIDPair{
			Type: elem.Type,
			ID:   elem.ID,
		})
	}

	return json.Marshal(tuple)
}

func (l *LocalVariableDeclaration) MarshalJSON() ([]byte, error) {
	pair := make([]any, 0)
	pair = append(pair, struct {
		Type Type
	}{
		Type: l.Type,
	}, struct {
		ID ID
	}{
		ID: l.ID,
	})

	if l.LocalVariableDeclarationRest == nil {
		return json.Marshal(pair)
	}

	for _, elem := range *l.LocalVariableDeclarationRest {
		pair = append(pair, struct {
			ID ID
		}{
			ID: elem.ID,
		})
	}
	return json.Marshal(pair)
}

func (s Statement) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string
		Statement any
	}{
		Type:      TypeString[s.Type],
		Statement: s.Statement,
	})
}

func (f Factor) MarshalJSON() ([]byte, error) {
	switch f.Factor.(type) {
	case FactorTuple:
		return json.Marshal(f.Factor)
	case ID, lexer.Token:
		return json.Marshal(struct {
			Literal  any
			Type     string
			Category string
		}{
			Literal:  f.Factor.(lexer.Token).Literal,
			Type:     lexer.TokenTypeString[f.Factor.(lexer.Token).Type],
			Category: lexer.TokenCategoryString[f.Factor.(lexer.Token).Category],
		})
	default:
		return json.Marshal(f.Factor)
	}
}

func (t Term) MarshalJSON() ([]byte, error) {
	if t.TermRest == nil {
		return json.Marshal(struct {
			Factor Factor
		}{
			Factor: t.Factor,
		})
	}
	return json.Marshal(struct {
		LFactor Factor
		Op      lexer.Token
		RFactor Factor
	}{
		LFactor: t.Factor,
		Op:      t.TermRest.MulOrDiv,
		RFactor: t.TermRest.Factor,
	})
}

func (e Exp) MarshalJSON() ([]byte, error) {
	if e.ExpRest == nil {
		return json.Marshal(struct {
			Term Term
		}{
			Term: e.Term,
		})
	}
	return json.Marshal(struct {
		LTerm Term
		Op    lexer.Token
		RTerm Term
	}{
		LTerm: e.Term,
		Op:    e.ExpRest.PlusOrMinus,
		RTerm: e.ExpRest.Term,
	})
}

func (c CmpOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Literal  string
		Type     string
		Category string
	}{
		Literal:  c.Literal.(string),
		Type:     lexer.TokenTypeString[c.Type],
		Category: lexer.TokenCategoryString[c.Category],
	})
}

func (conditionalExp ConditionalExp) MarshalJSON() ([]byte, error) {
	if conditionalExp.ConditionalExpRest == nil {
		return json.Marshal(struct {
			RelationExp RelationExp
		}{
			RelationExp: conditionalExp.RelationExp,
		})
	}
	return json.Marshal(struct {
		LRelationExp RelationExp
		Or           lexer.Token
		RRelationExp RelationExp
	}{
		LRelationExp: conditionalExp.RelationExp,
		Or:           conditionalExp.ConditionalExpRest.Or,
		RRelationExp: conditionalExp.ConditionalExpRest.RelationExp,
	})
}

func (r RelationExp) MarshalJSON() ([]byte, error) {
	if r.RelationExpRest == nil {
		return json.Marshal(struct {
			CompExp CompExp
		}{
			CompExp: r.CompExp,
		})
	}
	return json.Marshal(struct {
		LCompExp CompExp
		And      lexer.Token
		RCompExp CompExp
	}{
		LCompExp: r.CompExp,
		And:      r.RelationExpRest.And,
		RCompExp: r.RelationExpRest.CompExp,
	})
}

func (a ActParamList) MarshalJSON() ([]byte, error) {
	if a.ActParamList == nil {
		return json.Marshal(nil)
	}
	if a.ActParamList.ActParamListRest == nil {
		return json.Marshal(struct {
			Exp Exp
		}{
			Exp: a.ActParamList.Exp,
		})
	}
	exps := make([]Exp, 0)
	exps = append(exps, a.ActParamList.Exp)
	for _, elem := range *a.ActParamList.ActParamListRest {
		exps = append(exps, elem.Exp)
	}
	return json.Marshal(exps)
}
