package mir

const (
	ERROR = iota
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
)

var OpString = map[int]string{
	ERROR:       "ERROR",
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
}

type Program struct {
	StmtSeq []Statement
}

type Param interface {
	p()
}

type Statement struct {
	Op   int
	Arg1 Param
	Arg2 string
}
