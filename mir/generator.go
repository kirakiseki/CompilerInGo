package mir

import (
	"CompilerInGo/hir"
	"fmt"
	"github.com/kpango/glg"
)

type Context struct {
	MethodIn      MethodInfo
	LoopCondLabel int
	LoopEndLabel  int
	CallLabel     int
}

type MethodInfo struct {
	Name      string
	Pos       int
	ActParams []int
	ReturnVar int
}

type MIRGenerator struct {
	Program    *Program
	HIRProgram *hir.Program
	Vars       map[string]int
	Methods    map[string]MethodInfo
	Context    Context
	CtxStack   *Stack[Context]
	MethodSeq  []Statement
}

func NewMIRGenerator() *MIRGenerator {
	context := Context{
		LoopCondLabel: 0,
		LoopEndLabel:  0,
		MethodIn: MethodInfo{
			Name: "main",
		},
		CallLabel: 0,
	}
	ctxStack := NewStack[Context]()
	ctxStack.Push(context)
	return &MIRGenerator{
		Program:  NewProgram(),
		Vars:     make(map[string]int),
		Methods:  make(map[string]MethodInfo),
		Context:  context,
		CtxStack: ctxStack,
	}
}

func (g *MIRGenerator) Generate(program *hir.Program) *Program {
	g.HIRProgram = program
	mainMethod := program.GetMethod("main")

	if mainMethod == nil {
		glg.Fatal("No main method found")
	}

	g.Program.StmtSeq = append(g.Program.StmtSeq, g.generateStatement(*mainMethod.Body)...)

	offset := len(g.Program.StmtSeq)

	g.Program.StmtSeq = append(g.Program.StmtSeq, g.MethodSeq...)

	for idx, stmt := range g.Program.StmtSeq {
		arg1, arg2, res := stmt.Arg1, stmt.Arg2, stmt.Res
		if arg1.Str() == "_T_HERE_TO_JMP+1" {
			g.Program.StmtSeq[idx].Arg1 = IntParam(idx + 2)
		}
		if arg2.Str() == "_T_HERE_TO_JMP+1" {
			g.Program.StmtSeq[idx].Arg2 = IntParam(idx + 2)
		}
		if res.Str() == "_T_HERE_TO_JMP+1" {
			g.Program.StmtSeq[idx].Res = IntParam(idx + 2)
		}

		var method string
		if cnt, err := fmt.Sscanf(arg1.Str(), "_T_JMP_%s", &method); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Arg1 = IntParam(g.Methods[method].Pos + offset)
		}
		if cnt, err := fmt.Sscanf(arg2.Str(), "_T_JMP_%s", &method); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Arg2 = IntParam(g.Methods[method].Pos + offset)
		}
		if cnt, err := fmt.Sscanf(res.Str(), "_T_JMP_%s", &method); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Res = IntParam(g.Methods[method].Pos + offset)
		}
	}

	return g.Program
}

func (g *MIRGenerator) NewVar(name string) int {
	g.Vars[name] = len(g.Vars) + 1
	return g.Vars[name]
}

func (g *MIRGenerator) NewAnonymousVar() int {
	return g.NewVar(fmt.Sprintf("%d", len(g.Vars)+2))
}

func (g *MIRGenerator) GetVar(name string) int {
	return g.Vars[name]
}

func (g *MIRGenerator) Print() {
	for idx, stmt := range g.Program.StmtSeq {
		_ = glg.Infof("%3d| %s", idx, stmt.Str())
	}
}

func (g *MIRGenerator) NewMethod(name string) {
	g.Methods[name] = MethodInfo{
		Pos: len(g.Methods),
	}
}
