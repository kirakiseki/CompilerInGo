package utils

import (
	"bufio"
	"github.com/kpango/glg"
	"strings"
)

func GetLine(file []byte, Pos Position) string {
	reader := bufio.NewReader(strings.NewReader(string(file)))

	var line []byte
	var err error
	line, _, err = reader.ReadLine()
	for i := 0; i < int(Pos.Row-1); i++ {
		line, _, err = reader.ReadLine()
	}
	if err != nil {
		if err.Error() == "EOF" {
			return "EOF"
		}
		glg.Fatalln(err)
	}
	return string(line)
}
