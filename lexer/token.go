package lexer

import (
	"CompilerInGo/utils"
	"fmt"
	"github.com/kpango/glg"
)

// TokenPool Token池
type TokenPool struct {
	Pool []Token
}

type TokenCategory uint
type TokenType uint

type Token struct {
	Category TokenCategory      // 分类
	Type     TokenType          // 类型
	Literal  any                // 字面量
	Pos      utils.PositionPair // 位置
}

// TokenCategory
const (
	EOF     = iota //0 EOF
	KEYWORD        //1 关键字
	IDENT          //2 标识符
	INTEGER        //3 整数
	DECIMAL        //4 小数
	STR            //5 字符串
	CHAR           //6 字符
	DELIM          //7 分隔符
	OPERA          //8 运算符
	COMMENT        //9 注释
)

// TokenType
// 关键字
const (
	VOID   = 10 + iota //10 void
	VAR                //11 var
	INT                //12 int
	FLOAT              //13 float
	STRING             //14 string
	BEGIN              //15 begin
	END                //16 end
	IF                 //17 if
	THEN               //18 then
	ELSE               //19 else
	WHILE              //20 while
	DO                 //21 do
	CALL               //22 call
	READ               //23 read
	WRITE              //24 write
	AND                //25 and
	OR                 //26 or
)

// 分隔符
const (
	LBRACE    = 27 + iota //27 {
	RBRACE                //28 }
	LPAREN                //29 (
	RPAREN                //30 )
	SEMICOLON             //31 ;
	SPACE                 //32 空格
)

// 运算符
const (
	EQUAL        = 33 + iota //33 ==
	ASSIGN                   //34 =
	LESS                     //35 <
	LESSEQUAL                //36 <=
	GREATER                  //37 >
	GREATEREQUAL             //38 >=
	DIAMOND                  //39 <>
	PLUS                     //40 +
	MINUS                    //41 -
	TIMES                    //42 *
	DIVIDE                   //43 /
)

// 字面量
const (
	INTEGER_LITERAL            = 44 + iota //44 整数字面量
	DECIMAL_LITERAL                        //45 小数字面量
	STRING_LITERAL                         //46 字符串字面量
	CHAR_LITERAL                           //47 字符字面量
	EOF_LITERAL                            //48 EOF字面量
	SINGLELINE_COMMENT_LITERAL             //49 注释字面量
	MULTILINE_COMMENT_LITERAL              //50 注释字面量
)

// 标识符
const (
	IDENTIFIER = 51 + iota //51 标识符
)

// tokenType Token类型对应的字符串，输出时使用
var tokenType = map[TokenType]string{
	VOID:   "void",
	VAR:    "var",
	INT:    "int",
	FLOAT:  "float",
	STRING: "string",
	BEGIN:  "begin",
	END:    "end",
	IF:     "if",
	THEN:   "then",
	ELSE:   "else",
	WHILE:  "while",
	DO:     "do",
	CALL:   "call",
	READ:   "read",
	WRITE:  "write",
	AND:    "and",
	OR:     "or",

	LBRACE:    "LBRACE {",
	RBRACE:    "RBRACE }",
	LPAREN:    "LPAREN (",
	RPAREN:    "RPAREN )",
	SEMICOLON: "SEMICOLON ;",
	SPACE:     "SPACE",

	EQUAL:        "EQUAL ==",
	ASSIGN:       "ASSIGN =",
	LESS:         "LESS <",
	LESSEQUAL:    "LESSEQUAL <=",
	GREATER:      "GREATER >",
	GREATEREQUAL: "GREATEREQUAL >=",
	DIAMOND:      "DIAMOND <>",
	PLUS:         "PLUS +",
	MINUS:        "MINUS -",
	TIMES:        "TIMES *",
	DIVIDE:       "DIVIDE /",

	INTEGER_LITERAL:            "INTEGER_LITERAL",
	DECIMAL_LITERAL:            "DECIMAL_LITERAL",
	STRING_LITERAL:             "STRING_LITERAL",
	CHAR_LITERAL:               "CHAR_LITERAL",
	EOF_LITERAL:                "EOF_LITERAL",
	SINGLELINE_COMMENT_LITERAL: "SINGLELINE_COMMENT_LITERAL",
	MULTILINE_COMMENT_LITERAL:  "MULTILINE_COMMENT_LITERAL",

	IDENTIFIER: "IDENTIFIER",
}

// tokenCategory Token类别对应的字符串，输出时使用
var tokenCategory = map[TokenCategory]string{
	EOF:     "EOF",
	KEYWORD: "KEYWORD",
	IDENT:   "IDENTIFIER",
	INTEGER: "INTEGER",
	DECIMAL: "DECIMAL",
	STR:     "STRING",
	CHAR:    "CHAR",
	DELIM:   "DELIMITER",
	OPERA:   "OPERATOR",
	COMMENT: "COMMENT",
}

// IsKeyword 判断是否为关键字
func IsKeyword(s string) bool {
	switch s {
	case "void", "var", "int", "float", "string", "begin", "end", "if", "then", "else", "while", "do", "call", "read", "write", "and", "or":
		return true
	default:
		return false
	}
}

// IsDelim 判断是否为分隔符
func IsDelim(s string) bool {
	switch s {
	case "{", "}", "(", ")", ";", " ":
		return true
	default:
		return false
	}
}

// IsOpera 判断是否为运算符
func IsOpera(s string) bool {
	switch s {
	case "==", "=", "<", "<=", ">", ">=", "<>", "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

// CategoryName 获取Token类别对应的字符串
func (t *Token) CategoryName() string {
	return tokenCategory[t.Category]
}

// String 获取Token的字符串表示
func (t *Token) String() string {
	switch t.Type {
	case CHAR_LITERAL:
		// 字符字面量
		if t.Literal == "" {
			// 空字符
			return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), tokenType[t.Type], "")
		}
		switch t.Literal.(type) {
		case string:
			// 字符非空且为字符串（经过转义）
			return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), tokenType[t.Type], t.Literal)
		default:
			// 字符非空且未经转义
			return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), tokenType[t.Type], string(t.Literal.(rune)))
		}
	default:
		// 其他类型
		return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), tokenType[t.Type], t.Literal)
	}
}

// setType 设置Token的类型
func (t *Token) setType() {
	for k, v := range tokenType {
		if v == t.Literal {
			t.Type = k
		}
	}
}

// setCategory 设置Token的分类
func (t *Token) setCategory() {
	switch t.Type {
	case VOID, VAR, INT, FLOAT, STRING, BEGIN, END, IF, THEN, ELSE, WHILE, DO, CALL, READ, WRITE, AND, OR:
		t.Category = KEYWORD
	case LBRACE, RBRACE, LPAREN, RPAREN, SEMICOLON, SPACE:
		t.Category = DELIM
	case EQUAL, ASSIGN, LESS, LESSEQUAL, GREATER, GREATEREQUAL, DIAMOND, PLUS, MINUS, TIMES, DIVIDE:
		t.Category = OPERA
	case INTEGER_LITERAL:
		t.Category = INTEGER
	case DECIMAL_LITERAL:
		t.Category = DECIMAL
	case STRING_LITERAL:
		t.Category = STR
	case CHAR_LITERAL:
		t.Category = CHAR
	case EOF_LITERAL:
		t.Category = EOF
	case SINGLELINE_COMMENT_LITERAL, MULTILINE_COMMENT_LITERAL:
		t.Category = COMMENT
	default:
		t.Category = IDENT
	}
}

// NewToken 创建Token
func NewToken(literal any, position utils.PositionPair, tokenType TokenType) Token {
	switch literal.(type) {
	// 按类型创建
	case int64, float64:
		// 整数和小数
		token := Token{
			Literal: literal,
			Pos:     position,
			Type:    tokenType,
		}
		token.setCategory()
		return token
	case string:
		// 字符串
		token := Token{
			Literal: literal,
			Pos:     position,
			Type:    tokenType,
		}
		token.Literal = utils.EscapeString(token.Literal.(string)) // 转义字符串
		token.setCategory()
		return token
	case rune:
		// 字符
		token := Token{
			Literal: literal,
			Pos:     position,
			Type:    tokenType,
		}
		token.Literal = utils.EscapeRune(token.Literal.(rune)) // 转义字符
		token.setCategory()
		return token
	default:
		// 其他类型
		glg.Fatalln("Token Literal Type Error")
	}

	return Token{}
}

// NewTokenPool 创建Token池
func NewTokenPool() *TokenPool {
	return &TokenPool{
		Pool: make([]Token, 0),
	}
}

// Add 向Token池中添加Token
func (pl *TokenPool) Add(token Token) {
	pl.Pool = append(pl.Pool, token)
}

// Get 获取Token池中位置为index的Token
func (pl *TokenPool) Get(index int) Token {
	return pl.Pool[index]
}

// Len 获取Token池中Token的数量
func (pl *TokenPool) Len() int {
	return len(pl.Pool)
}

// Last 获取Token池中最后一个Token
func (pl *TokenPool) Last() Token {
	if pl.Len() == 0 {
		// Token池为空
		return Token{}
	}
	return pl.Pool[pl.Len()-1]
}
