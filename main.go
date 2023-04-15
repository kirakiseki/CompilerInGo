package main

import (
	"CompilerInGo/lexer"
	"CompilerInGo/utils"
	"github.com/kpango/glg"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {

	// 设置运行模式
	mode := "DEBUG"

	// 设置CPU Profiling
	if mode == "DEBUG" {
		// 创建CPU Profiling文件
		cpuProf, err := os.Create("./test/benchmark/cpu.prof")
		if err != nil {
			glg.Fatalln(err)
		}

		defer cpuProf.Close()

		// 使用pprof进行CPU Profiling
		runtime.SetCPUProfileRate(10000)
		_ = pprof.StartCPUProfile(cpuProf)
		defer pprof.StopCPUProfile()
	}

	// 初始化logger
	//utils.InitLogger(mode)
	utils.InitLogger("DEBUG")

	// 初始化lexer
	lex := lexer.NewLexer("test/fail.program")
	// 初始化token池
	tokenPool := lexer.NewTokenPool()

	// Lexer计时开始
	startTime := time.Now()

	// 读取第一个Token
	// IfTokenError 检查Token是否出错，若出错则输出错误信息并退出程序
	token := lexer.IfTokenError(lex.ScanToken())
	tokenPool.Add(token)
	// 若未读到EOF则继续读取
	for tokenPool.Last().Category != lexer.EOF {
		token := lexer.IfTokenError(lex.ScanToken())
		tokenPool.Add(token)
	}

	// Lexer计时结束
	elapsedTime := time.Since(startTime)

	// 输出Token池
	_ = glg.Info("Token Pool:")
	_ = glg.Infof("%3s:%3s to %3s:%3s %12s %27s (%v)", "Row", "Col", "Row", "Col", "Category", "Type", "Literal")

	for _, token := range tokenPool.Pool {
		_ = glg.Info(token.String())
	}
	// 显示Lexer运行时间
	_ = glg.Info("Lexing finished in ", elapsedTime)
}
