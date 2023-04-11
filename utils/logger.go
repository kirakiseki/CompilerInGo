package utils

import (
	"github.com/kpango/glg"
)

func InitLogger(level string) {
	// 按照level日志等级初始化logger
	switch level {
	case "INFO":
		// 设置日志等级为INFO
		glg.Get().SetMode(glg.STD).SetLevel(glg.INFO)
		_ = glg.Info("Logger initialized at INFO level.")
	case "DEBUG":
		fallthrough
	default:
		// 设置日志默认等级为DEBUG
		glg.Get().SetMode(glg.STD).SetLevel(glg.DEBG)
		_ = glg.Debug("Logger initialized at DEBUG level.")
	}
}
