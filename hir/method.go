package hir

// Program AST树的HIR表示，由多个Method组成
type Program struct {
	Methods []Method
}

// Method HIR中的方法，由返回类型，方法名，参数列表和方法体组成
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

// GetMethod 根据方法名获取方法
func (p Program) GetMethod(name string) *Method {
	for _, method := range p.Methods {
		if method.Name == name {
			return &method
		}
	}
	return nil
}
