package hir

type Program struct {
	Methods []Method
}

type Method struct {
	ReturnType ResultType
	Name       string
	Params     []*TypeIDPair
	Body       *Statement
}

func NewProgram(methods []Method) *Program {
	return &Program{
		Methods: methods,
	}
}

func NewMethod(t ResultType, name string, params []*TypeIDPair, body *Statement) *Method {
	return &Method{
		ReturnType: t,
		Name:       name,
		Params:     params,
		Body:       body,
	}
}
