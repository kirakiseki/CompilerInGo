package utils

import "github.com/kpango/glg"

func InitLogger(level string) {
	switch level {
	case "DEBUG":
		glg.Get().SetMode(glg.STD).SetLevel(glg.DEBG)
		_ = glg.Debug("Logger initialized at DEBUG level.")
	case "INFO":
		glg.Get().SetMode(glg.STD).SetLevel(glg.INFO)
		_ = glg.Info("Logger initialized at INFO level.")
	}
}
