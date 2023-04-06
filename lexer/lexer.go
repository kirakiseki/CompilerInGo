package lexer

import (
	"CompilerInGo/utils"
	"errors"
	"github.com/kpango/glg"
	"io"
	"os"
	"strings"
)

type Lexer struct {
	reader Reader
	pos    utils.Position
	file   string
}

func NewLexer(file string) *Lexer {
	lexer := &Lexer{
		reader: strings.NewReader(string(utils.MustValue(os.ReadFile(file)))), // 读取文件并检查读取状态
		pos:    utils.Position{Row: 1},                                        // 设置初始位置
		file:   file,                                                          //文件路径
	}
	lexer.nextRune() // 读取第一个字符
	return lexer
}

func (l *Lexer) nextRune() rune {
	// 读取下一个字符
	ch, size, err := l.reader.ReadRune()

	if errors.Is(err, io.EOF) { // EOF 文件末尾
		l.pos.Ch = 0
		return 0
	} else if err != nil { // 读错误
		glg.Fatalln(err)
	}

	l.pos.Ch = ch
	l.pos.FilePos += uint(size)

	if ch == '\n' {
		l.pos.Row++
		l.pos.Col = 0
	} else {
		l.pos.Col++
	}

	return ch
}

func (l *Lexer) skipSpace() {
	for {
		ch := l.pos.Ch
		if ch != ' ' && ch != '\t' && ch != '\r' {
			break
		}
		l.nextRune()
	}
}

func (l *Lexer) TraverseRune() {
	for {
		l.skipSpace()
		_ = glg.Debugf("%c %+v", l.pos.Ch, l.pos)
		ch := l.nextRune()
		if ch == 0 {
			break
		}
	}
}

func (l *Lexer) scanString() Token {
	tokenPos := utils.PositionPair{Begin: l.pos}
	// TODO: scan string
	return NewLiteralToken("", tokenPos, STRING_LITERAL)
}
