package mir

import (
	"CompilerInGo/hir"
	"fmt"
	"github.com/kpango/glg"
	"strings"
)

// Context 上下文
type Context struct {
	MethodIn      MethodInfo // 当前所在方法
	LoopCondLabel int        // 循环条件标签
	LoopEndLabel  int        // 循环结束标签
	CallLabel     int        // 方法调用标签
}

// MethodInfo 方法信息
type MethodInfo struct {
	Name      string // 方法名
	Pos       int    // 方法在程序中的位置
	ActParams []int  // 方法实参
	ReturnVar int    // 方法返回值
}

// MIRGenerator 中间代码生成器
type MIRGenerator struct {
	Program    *Program              // 中间代码
	HIRProgram *hir.Program          // HIR程序
	Vars       map[string]int        // 变量表
	Labels     map[int]int           // 标签表
	Methods    map[string]MethodInfo // 方法表
	Context    Context               // 当前上下文
	CtxStack   *Stack[Context]       // 上下文栈
	MethodSeq  []Statement           // 方法序列
}

// NewMIRGenerator 新建中间代码生成器
func NewMIRGenerator() *MIRGenerator {
	context := Context{
		LoopCondLabel: -1,
		LoopEndLabel:  -1,
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

// Generate 生成中间代码
func (g *MIRGenerator) Generate(program *hir.Program) *Program {
	// 初始化
	g.HIRProgram = program

	// 从main方法开始生成
	mainMethod := program.GetMethod("main")

	// 检查main方法是否存在
	if mainMethod == nil {
		glg.Fatal("No main method found")
	}

	// 生成main方法
	g.Program.StmtSeq = append(g.Program.StmtSeq, g.generateStatement(*mainMethod.Body)...)

	// 将其他方法添加到方法序列，并拼接到main方法后
	offset := len(g.Program.StmtSeq)
	g.Program.StmtSeq = append(g.Program.StmtSeq, g.MethodSeq...)

	// 对预先定义的跳转位置进行修正
	// 遍历所有语句
	for idx, stmt := range g.Program.StmtSeq {
		arg1, arg2, res := stmt.Arg1, stmt.Arg2, stmt.Res

		// 跳转到当前语句+2
		if arg1.Str() == "_T_HERE_TO_JMP+1" {
			g.Program.StmtSeq[idx].Arg1 = IntParam(idx + 2)
		}
		if arg2.Str() == "_T_HERE_TO_JMP+1" {
			g.Program.StmtSeq[idx].Arg2 = IntParam(idx + 2)
		}
		if res.Str() == "_T_HERE_TO_JMP+1" {
			g.Program.StmtSeq[idx].Res = IntParam(idx + 2)
		}

		// 按照当前语句的位置进行相对ref跳转
		var ref int
		if cnt, err := fmt.Sscanf(arg1.Str(), "_T_JMP_REF_%d", &ref); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Arg1 = IntParam(idx + ref)
		}
		if cnt, err := fmt.Sscanf(arg2.Str(), "_T_JMP_REF_%d", &ref); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Arg2 = IntParam(idx + ref)
		}
		if cnt, err := fmt.Sscanf(res.Str(), "_T_JMP_REF_%d", &ref); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Res = IntParam(idx + ref)
		}

		// 跳转到方法开始位置
		var method string
		if cnt, err := fmt.Sscanf(arg1.Str(), "_T_JMP_METHOD_%s", &method); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Arg1 = IntParam(g.Methods[method].Pos + offset)
		}
		if cnt, err := fmt.Sscanf(arg2.Str(), "_T_JMP_METHOD_%s", &method); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Arg2 = IntParam(g.Methods[method].Pos + offset)
		}
		if cnt, err := fmt.Sscanf(res.Str(), "_T_JMP_METHOD_%s", &method); err == nil && cnt == 1 {
			g.Program.StmtSeq[idx].Res = IntParam(g.Methods[method].Pos + offset)
		}
	}

	// 处理continue和break
	for idx, stmt := range g.Program.StmtSeq {
		comment := stmt.Comment
		if comment == "_T_CONTINUE" {
			for i := idx; i < len(g.Program.StmtSeq); i++ {
				if strings.HasPrefix(g.Program.StmtSeq[i].Comment, "next loop:") {
					_ = glg.Warn(g.Program.StmtSeq[i].Str())
					g.Program.StmtSeq[idx].Res = IntParam(g.Program.StmtSeq[i].Res.Int())
					break
				}
				if i == len(g.Program.StmtSeq)-1 {
					glg.Fatal("No loop end found")
				}
			}
		} else if comment == "_T_BREAK" {
			for i := idx; i >= 0; i-- {
				if strings.HasPrefix(g.Program.StmtSeq[i].Comment, "while condition") {
					g.Program.StmtSeq[idx].Res = IntParam(g.Program.StmtSeq[i].Res.Int())
					break
				}
				if i == 0 {
					glg.Fatal("No loop context found")
				}
			}
		}
	}

	return g.Program
}

// NewVar 生成新的变量
func (g *MIRGenerator) NewVar(name string) int {
	g.Vars[name] = len(g.Vars) + 1
	return g.Vars[name]
}

// NewAnonymousVar 生成匿名变量
func (g *MIRGenerator) NewAnonymousVar() int {
	return g.NewVar(fmt.Sprintf("%d", len(g.Vars)+2))
}

// GetVar 获取变量
func (g *MIRGenerator) GetVar(name string) int {
	return g.Vars[name]
}

// Print 打印中间代码
func (g *MIRGenerator) Print() {
	for idx, stmt := range g.Program.StmtSeq {
		_ = glg.Infof("%3d| %s", idx, stmt.Str())
	}
}

// NewMethod 添加方法定义
func (g *MIRGenerator) NewMethod(name string) {
	g.Methods[name] = MethodInfo{
		Pos: len(g.Methods),
	}
}

// NewLabel 生成新的标签
func (g *MIRGenerator) NewLabel() int {
	return len(g.Labels) + 1
}

// SetLabel 设置标签
func (g *MIRGenerator) SetLabel(label, value int) {
	g.Labels[label] = value
}

// GetLabel 获取标签
func (g *MIRGenerator) GetLabel(label int) int {
	return g.Labels[label]
}
