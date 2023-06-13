package mir

import (
	"fmt"
	"github.com/kpango/glg"
	"strconv"
)

// 操作符
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
	JZERO
	JNZERO
	STOP
)

// OpString 操作符字符串，输出用
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
	JZERO:       "j0",
	JNZERO:      "j!0",
	STOP:        "STOP",
}

// Program 输出的中间代码，四元式格式
type Program struct {
	StmtSeq []Statement
}

func NewProgram() *Program {
	return &Program{
		StmtSeq: make([]Statement, 0),
	}
}

// Param 参数，可以是IntParam, FloatParam, StrParam， *Statement
type Param interface {
	p()
	Str() string
	Int() int
}

type StrParam string

func (s StrParam) p() {}

func (s StrParam) Str() string {
	return string(s)
}

func (s StrParam) Int() int {
	i, err := strconv.Atoi(string(s))
	if err != nil {
		glg.Fatal(err)
	}
	return i
}

func (i IntParam) Str() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i IntParam) Int() int {
	return int(i)
}

func (f FloatParam) Str() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (f FloatParam) Int() int {
	glg.Fatal("Cannot convert float to int")
	return 0
}

type IntParam int64

func (i IntParam) p() {}

type FloatParam float64

func (f FloatParam) p() {}

func (s *Statement) p() {}

func (s *Statement) Str() string {
	str := fmt.Sprintf("( %4s, %4s, %4s, %4s)", OpString[s.Op], s.Arg1.Str(), s.Arg2.Str(), s.Res.Str())
	if s.Comment == "" {
		return str
	}
	return fmt.Sprintf("%-30s   # %s", str, s.Comment)
}

func (s *Statement) Int() int {
	glg.Fatal("Cannot convert statement to int")
	return 0
}

// Statement 四元式语句
type Statement struct {
	Op      int    // 操作符
	Arg1    Param  // 参数1
	Arg2    Param  // 参数2
	Res     Param  // 结果
	Comment string // 注释
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

// Label 生成一个新的标签
func (p *Program) Label() int {
	return len(p.StmtSeq) - 1
}
