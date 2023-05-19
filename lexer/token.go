package lexer

import (
	"CompilerInGo/utils"
	"encoding/json"
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
	Pos      utils.PositionPair `json:"-"` // 位置
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
	VOID     = 10 + iota //10 void
	VAR                  //11 var
	INT                  //12 int
	FLOAT                //13 float
	STRING               //14 string
	BEGIN                //15 begin
	END                  //16 end
	IF                   //17 if
	THEN                 //18 then
	ELSE                 //19 else
	WHILE                //20 while
	DO                   //21 do
	CALL                 //22 call
	READ                 //23 read
	WRITE                //24 write
	AND                  //25 and
	OR                   //26 or
	RETURN               //27 return
	CONTINUE             //28 continue
	BREAK                //29 break
)

// 分隔符
const (
	LBRACE    = 30 + iota //30 {
	RBRACE                //31 }
	LPAREN                //32 (
	RPAREN                //33 )
	SEMICOLON             //34 ;
	SPACE                 //35 空格
	COMMA                 //36 ,
)

// 运算符
const (
	EQUAL        = 37 + iota //37 ==
	ASSIGN                   //38 =
	LESS                     //39 <
	LESSEQUAL                //40 <=
	GREATER                  //41 >
	GREATEREQUAL             //42 >=
	DIAMOND                  //43 <>
	PLUS                     //44 +
	MINUS                    //45 -
	TIMES                    //46 *
	DIVIDE                   //47 /
)

// 字面量
const (
	INTEGER_LITERAL            = 48 + iota //48 整数字面量
	DECIMAL_LITERAL                        //49 小数字面量
	STRING_LITERAL                         //50 字符串字面量
	CHAR_LITERAL                           //51 字符字面量
	EOF_LITERAL                            //52 EOF字面量
	SINGLELINE_COMMENT_LITERAL             //53 注释字面量
	MULTILINE_COMMENT_LITERAL              //54 注释字面量
)

// 标识符
const (
	IDENTIFIER = 55 + iota //55 标识符
)

// TokenTypeString Token类型对应的字符串，输出时使用
var TokenTypeString = map[TokenType]string{
	VOID:     "void",
	VAR:      "var",
	INT:      "int",
	FLOAT:    "float",
	STRING:   "string",
	BEGIN:    "begin",
	END:      "end",
	IF:       "if",
	THEN:     "then",
	ELSE:     "else",
	WHILE:    "while",
	DO:       "do",
	CALL:     "call",
	READ:     "read",
	WRITE:    "write",
	AND:      "and",
	OR:       "or",
	RETURN:   "return",
	CONTINUE: "continue",
	BREAK:    "break",

	LBRACE:    "LBRACE {",
	RBRACE:    "RBRACE }",
	LPAREN:    "LPAREN (",
	RPAREN:    "RPAREN )",
	SEMICOLON: "SEMICOLON ;",
	SPACE:     "SPACE",
	COMMA:     "COMMA ,",

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

// TokenCategoryString Token类别对应的字符串，输出时使用
var TokenCategoryString = map[TokenCategory]string{
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
	case "void", "var", "int", "float", "string", "begin", "end", "if", "then", "else", "while", "do", "call", "read", "write", "and", "or", "return", "continue", "break":
		return true
	default:
		return false
	}
}

// IsDelim 判断是否为分隔符
func IsDelim(s string) bool {
	switch s {
	case "{", "}", "(", ")", ";", " ", ",":
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
	return TokenCategoryString[t.Category]
}

// String 获取Token的字符串表示
func (t *Token) String() string {
	switch t.Type {
	case CHAR_LITERAL:
		// 字符字面量
		if t.Literal == "" {
			// 空字符
			return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), TokenTypeString[t.Type], "")
		}
		switch t.Literal.(type) {
		case string:
			// 字符非空且为字符串（经过转义）
			return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), TokenTypeString[t.Type], t.Literal)
		default:
			// 字符非空且未经转义
			return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), TokenTypeString[t.Type], string(t.Literal.(rune)))
		}
	default:
		// 其他类型
		return fmt.Sprintf("%3d:%3d to %3d:%3d %12s %27s (%v)", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col, t.CategoryName(), TokenTypeString[t.Type], t.Literal)
	}
}

// setType 设置Token的类型
func (t *Token) setType() {
	for k, v := range TokenTypeString {
		if v == t.Literal {
			t.Type = k
		}
	}
}

// setCategory 设置Token的分类
func (t *Token) setCategory() {
	switch t.Type {
	case VOID, VAR, INT, FLOAT, STRING, BEGIN, END, IF, THEN, ELSE, WHILE, DO, CALL, READ, WRITE, AND, OR, CONTINUE, BREAK, RETURN:
		t.Category = KEYWORD
	case LBRACE, RBRACE, LPAREN, RPAREN, SEMICOLON, SPACE, COMMA:
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

func (t Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Literal  any
		Type     string
		Category string
	}{
		Literal:  t.Literal,
		Type:     TokenTypeString[t.Type],
		Category: TokenCategoryString[t.Category],
	})
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

// PushBack 向Token池后添加Token
func (pl *TokenPool) PushBack(token Token) {
	pl.Pool = append(pl.Pool, token)
}

// PushFront 向Token池钱添加Token
func (pl *TokenPool) PushFront(token Token) {
	pl.Pool = append([]Token{token}, pl.Pool...)
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

func (pl *TokenPool) Pop() func() (Token, error) {
	index := 0
	return func() (Token, error) {
		if index >= pl.Len() {
			return Token{}, utils.NewError("TokenPool Pop Out of range")
		}
		token := pl.Pool[index]
		index++
		return token, nil
	}
}
