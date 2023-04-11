package utils

import (
	"bufio"
	"github.com/kpango/glg"
	"strings"
)

// GetLine 获取文件中指定行的内容
func GetLine(file []byte, Pos Position) string {
	// 创建读取器
	reader := bufio.NewReader(strings.NewReader(string(file)))

	var line []byte
	var err error

	// 读取第一行
	line, _, err = reader.ReadLine()
	for i := 0; i < int(Pos.Row-1); i++ {
		// 读取到指定行
		line, _, err = reader.ReadLine()
	}
	if err != nil {
		// 读取错误
		if err.Error() == "EOF" {
			// 读取到文件末尾
			return "EOF"
		}
		// 读取时发生错误
		glg.Fatalln(err)
	}

	// 返回指定行内容
	return string(line)
}
