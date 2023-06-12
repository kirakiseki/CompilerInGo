package mir

import (
	"fmt"
	"strconv"
)

const (
	ERROR = iota
	ASSIGN
	PLUS
	MINUS
	TIMES
	DIVIDE
	JMP
	JEQUAL
	JNEQUAL
	JGREAT
	JGREATEQUAL
	JLESS
	JLESSEQUAL
	STOP
)

var OpString = map[int]string{
	ERROR:       "ERROR",
	ASSIGN:      "=",
	PLUS:        "+",
	MINUS:       "-",
	TIMES:       "*",
	DIVIDE:      "/",
	JMP:         "j",
	JEQUAL:      "j=",
	JNEQUAL:     "j!=",
	JGREAT:      "j>",
	JGREATEQUAL: "j>=",
	JLESS:       "j<",
	JLESSEQUAL:  "j<=",
	STOP:        "STOP",
}

type Program struct {
	StmtSeq []Statement
}

func NewProgram() *Program {
	return &Program{
		StmtSeq: make([]Statement, 0),
	}
}

type Param interface {
	p()
	Str() string
}

type StrParam string

func (s StrParam) p() {}

func (s StrParam) Str() string {
	return string(s)
}

func (i IntParam) Str() string {
	return strconv.FormatInt(int64(i), 10)
}

func (f FloatParam) Str() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

type IntParam int64

func (i IntParam) p() {}

type FloatParam float64

func (f FloatParam) p() {}

func (s *Statement) p() {}

func (s *Statement) Str() string {
	str := fmt.Sprintf("(%s, %s, %s, %s)", OpString[s.Op], s.Arg1.Str(), s.Arg2.Str(), s.Res.Str())
	if s.Comment == "" {
		return str
	}
	return fmt.Sprintf("%-30s   # %s", str, s.Comment)
}

type Statement struct {
	Op      int
	Arg1    Param
	Arg2    Param
	Res     Param
	Comment string
}

func NewStatement(op int, arg1, arg2, res Param, comm string) *Statement {
	return &Statement{
		Op:      op,
		Arg1:    arg1,
		Arg2:    arg2,
		Res:     res,
		Comment: comm,
	}
}

func (p *Program) Label() int {
	return len(p.StmtSeq) - 1
}
