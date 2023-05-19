package main

import (
	"CompilerInGo/lexer"
	"CompilerInGo/parser"
	"CompilerInGo/utils"
	"bytes"
	"encoding/json"
	"flag"
	"github.com/kpango/glg"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	// 解析命令行参数
	filepath := flag.String("f", "./test.program", "input source program")
	mode := flag.String("m", "DEBUG", "logger mode (DEBUG, INFO, CLOSE)")
	flag.Parse()

	// 设置CPU Profiling
	if *mode == "DEBUG" {
		// 创建CPU Profiling文件
		cpuProf, err := os.Create("./cpu.prof")
		if err != nil {
			glg.Fatalln(err)
		}

		defer cpuProf.Close()

		// 使用pprof进行CPU Profiling
		runtime.SetCPUProfileRate(3000)
		_ = pprof.StartCPUProfile(cpuProf)
		defer pprof.StopCPUProfile()
	}

	// 初始化logger
	utils.InitLogger(*mode)

	// 初始化lexer
	lex := lexer.NewLexer(*filepath)
	_ = glg.Info("Lexer initialized")

	// 初始化token池
	lexer.Pool = lexer.NewTokenPool()

	// Lexer计时开始
	startTime := time.Now()

	// 读取第一个Token
	// IfTokenError 检查Token是否出错，若出错则输出错误信息并退出程序
	token := lexer.IfTokenError(lex.ScanToken())
	lexer.Pool.PushBack(token)
	// 若未读到EOF则继续读取
	for lexer.Pool.Last().Category != lexer.EOF {
		token := lexer.IfTokenError(lex.ScanToken())
		lexer.Pool.PushBack(token)
	}

	// Lexer计时结束
	elapsedTime := time.Since(startTime)

	// 输出Token池
	_ = glg.Info("Token Pool:")
	_ = glg.Infof("%3s:%3s to %3s:%3s %12s %27s (%v)", "Row", "Col", "Row", "Col", "Category", "Type", "Literal")

	for _, token := range lexer.Pool.Pool {
		_ = glg.Info(token.String())
	}
	// 显示Lexer运行时间
	_ = glg.Info("Lexing finished in ", elapsedTime)

	// 初始化parser
	pser := parser.NewParser()
	_ = glg.Info("Parser initialized")

	// Parser计时开始
	startTime = time.Now()

	// 开始parse
	program, err := pser.Parse()
	if err != nil {
		glg.Fatal(err)
	}

	// Parser计时结束
	elapsedTime = time.Since(startTime)

	// 输出AST
	// 将AST转换为JSON格式
	marshaled, _ := json.Marshal(program)
	// 在Debug模式下输出原始JSON
	_ = glg.Debug("Raw JSON:", string(marshaled))

	// 在Info模式下输出格式化后的JSON
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, marshaled, "", "....")
	_ = glg.Info("Parser Result:")
	_ = glg.Info("AST in JSON format:\n" + prettyJSON.String())

	// 显示Parser运行时间
	_ = glg.Info("Parsing finished in ", elapsedTime)
}
