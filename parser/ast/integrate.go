package ast

// 为了方便输出AST，我们需要为AST中的每个结构体实现MarshalJSON()方法
// 输出指定内容，略去位置，null值，以及一些不必要的信息
const EMPTY = 0

const (
	PLUS = 1 + iota
	MINUS
	TIMES
	DIVIDE
)

const (
	LESS = 5 + iota
	LESSEQUAL
	GREATER
	GREATEREQUAL
	EQUAL
	DIAMOND
)

const (
	OR = 11 + iota
	AND
)

type TypeIDSequence struct {
	Seq []TypeIDPair
}

func (paramList *ParamList) Integrate() (*TypeIDSequence, error) {
	if paramList.ParamList == nil {
		return nil, nil
	}

	tuple := make([]TypeIDPair, 0)
	tuple = append(tuple, TypeIDPair{
		Type: paramList.ParamList.Type,
		ID:   paramList.ParamList.ID,
	})

	if paramList.ParamList.ParamListRest != nil {
		for _, elem := range *paramList.ParamList.ParamListRest {
			tuple = append(tuple, TypeIDPair{
				Type: elem.Type,
				ID:   elem.ID,
			})
		}
	}

	return &TypeIDSequence{
		Seq: tuple,
	}, nil
}

func (l *LocalVariableDeclaration) Integrate() (TypeIDSequence, error) {
	tuple := make([]TypeIDPair, 0)
	tuple = append(tuple, TypeIDPair{
		Type: l.Type,
		ID:   l.ID,
	})

	if l.LocalVariableDeclarationRest == nil {
		return TypeIDSequence{
			Seq: tuple,
		}, nil
	}

	for _, rest := range *l.LocalVariableDeclarationRest {
		tuple = append(tuple, TypeIDPair{
			Type: l.Type,
			ID:   rest.ID,
		})
	}
	return TypeIDSequence{
		Seq: tuple,
	}, nil
}

func (t Term) Integrate() []Factor {
	if t.TermRest == nil {
		return []Factor{t.Factor}
	}
	return []Factor{t.Factor, t.TermRest.Factor}
}

func (e Exp) Integrate() []Term {
	if e.ExpRest == nil {
		return []Term{e.Term}
	}
	return []Term{e.Term, e.ExpRest.Term}
}

func (conditionalExp ConditionalExp) Integrate() []RelationExp {
	if conditionalExp.ConditionalExpRest == nil {
		return []RelationExp{conditionalExp.RelationExp}
	}
	return []RelationExp{conditionalExp.RelationExp, conditionalExp.ConditionalExpRest.RelationExp}
}

func (r RelationExp) Integrate() []CompExp {
	if r.RelationExpRest == nil {
		return []CompExp{r.CompExp}
	}
	return []CompExp{r.CompExp, r.RelationExpRest.CompExp}
}

func (c CompExp) Integrate() []Exp {
	return []Exp{c.LExp, c.RExp}
}

func (a ActParamList) Integrate() ([]Exp, error) {
	if a.ActParamList == nil {
		return []Exp{}, nil
	}
	if a.ActParamList.ActParamListRest == nil {
		return []Exp{a.ActParamList.Exp}, nil
	}
	exps := make([]Exp, 0)
	exps = append(exps, a.ActParamList.Exp)
	for _, elem := range *a.ActParamList.ActParamListRest {
		exps = append(exps, elem.Exp)
	}
	return exps, nil
}
