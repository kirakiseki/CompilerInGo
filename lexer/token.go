package lexer

import (
	"CompilerInGo/utils"
	"fmt"
	"github.com/kpango/glg"
)

type TokenPool struct {
	Pool []Token
}

type TokenCategory uint
type TokenType uint

type Token struct {
	Category TokenCategory
	Type     TokenType
	Literal  any
	Pos      utils.PositionPair
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
	VOID   = 10 + iota //9 void
	VAR                //10 var
	INT                //11 int
	FLOAT              //12 float
	STRING             //13 string
	BEGIN              //14 begin
	END                //15 end
	IF                 //16 if
	THEN               //17 then
	ELSE               //18 else
	WHILE              //19 while
	DO                 //20 do
	CALL               //21 call
	READ               //22 read
	WRITE              //23 write
	AND                //24 and
	OR                 //25 or
)

// 分隔符
const (
	LBRACE    = 27 + iota //26 {
	RBRACE                //27 }
	LPAREN                //28 (
	RPAREN                //29 )
	SEMICOLON             //30 ;
	SPACE                 //31 空格
)

// 运算符
const (
	EQUAL        = 33 + iota //32 ==
	ASSIGN                   //33 =
	LESS                     //34 <
	LESSEQUAL                //35 <=
	GREATER                  //36 >
	GREATEREQUAL             //37 >=
	DIAMOND                  //38 <>
	PLUS                     //39 +
	MINUS                    //40 -
	TIMES                    //41 *
	DIVIDE                   //42 /
)

// 字面量
const (
	INTEGER_LITERAL            = 44 + iota //43 整数字面量
	DECIMAL_LITERAL                        //44 小数字面量
	STRING_LITERAL                         //45 字符串字面量
	CHAR_LITERAL                           //46 字符字面量
	EOF_LITERAL                            //47 EOF字面量
	SINGLELINE_COMMENT_LITERAL             //48 注释字面量
	MULTILINE_COMMENT_LITERAL              //49 注释字面量
)

// 标识符
const (
	IDENTIFIER = 51 + iota //48 标识符
)

var tokenString = map[TokenType]string{
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

	LBRACE:    "{",
	RBRACE:    "}",
	LPAREN:    "(",
	RPAREN:    ")",
	SEMICOLON: ";",
	SPACE:     "SPACE",

	EQUAL:        "==",
	ASSIGN:       "=",
	LESS:         "<",
	LESSEQUAL:    "<=",
	GREATER:      ">",
	GREATEREQUAL: ">=",
	DIAMOND:      "<>",
	PLUS:         "+",
	MINUS:        "-",
	TIMES:        "*",
	DIVIDE:       "/",

	INTEGER_LITERAL:            "INTEGER_LITERAL",
	DECIMAL_LITERAL:            "DECIMAL_LITERAL",
	STRING_LITERAL:             "STRING_LITERAL",
	CHAR_LITERAL:               "CHAR_LITERAL",
	EOF_LITERAL:                "EOF_LITERAL",
	SINGLELINE_COMMENT_LITERAL: "SINGLELINE_COMMENT_LITERAL",
	MULTILINE_COMMENT_LITERAL:  "MULTILINE_COMMENT_LITERAL",

	IDENTIFIER: "IDENTIFIER",
}

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

func IsKeyword(s string) bool {
	switch s {
	case "void", "var", "int", "float", "string", "begin", "end", "if", "then", "else", "while", "do", "call", "read", "write", "and", "or":
		return true
	default:
		return false
	}
}

func IsDelim(s string) bool {
	switch s {
	case "{", "}", "(", ")", ";", " ":
		return true
	default:
		return false
	}
}

func IsOpera(s string) bool {
	switch s {
	case "==", "=", "<", "<=", ">", ">=", "<>", "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func IsLiteral(s string) bool {
	switch s {
	case "INTEGER_LITERAL", "DECIMAL_LITERAL", "STRING_LITERAL", "CHAR_LITERAL":
		return true
	default:
		return false
	}
}

func (t *Token) CategoryName() string {
	return tokenCategory[t.Category]
}

func (t *Token) String() string {
	switch t.Type {
	case CHAR_LITERAL:
		if t.Literal == "" {
			return fmt.Sprintf("%v %s(%v) <Begin:Row:%d,Col:%d  End:Row:%d,Col:%d>", t.CategoryName(), tokenString[t.Type], "", t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col)
		}
		switch t.Literal.(type) {
		case string:
			return fmt.Sprintf("%v %s(%v) <Begin:Row:%d,Col:%d  End:Row:%d,Col:%d>", t.CategoryName(), tokenString[t.Type], t.Literal, t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col)
		default:
			return fmt.Sprintf("%v %s(%v) <Begin:Row:%d,Col:%d  End:Row:%d,Col:%d>", t.CategoryName(), tokenString[t.Type], string(t.Literal.(rune)), t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col)
		}
	default:
		return fmt.Sprintf("%v %s(%v) <Begin:Row:%d,Col:%d  End:Row:%d,Col:%d>", t.CategoryName(), tokenString[t.Type], t.Literal, t.Pos.Begin.Row, t.Pos.Begin.Col, t.Pos.End.Row, t.Pos.End.Col)
	}
}

func (t *Token) setType() {
	for k, v := range tokenString {
		if v == t.Literal {
			t.Type = k
		}
	}
}

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

func NewToken(literal any, position utils.PositionPair, tokenType TokenType) Token {
	switch literal.(type) {
	case int64, float64:
		token := Token{
			Literal: literal,
			Pos:     position,
			Type:    tokenType,
		}
		token.setCategory()
		return token
	case string:
		token := Token{
			Literal: literal,
			Pos:     position,
			Type:    tokenType,
		}
		token.Literal = utils.EscapeString(token.Literal.(string))
		token.setCategory()
		return token
	case rune:
		token := Token{
			Literal: literal,
			Pos:     position,
			Type:    tokenType,
		}
		token.Literal = utils.EscapeRune(token.Literal.(rune))
		token.setCategory()
		return token
	default:
		glg.Fatalln("Token Literal Type Error")
	}

	return Token{}
}

func NewTokenPool() *TokenPool {
	return &TokenPool{
		Pool: make([]Token, 0),
	}
}

func (pl *TokenPool) Add(token Token) {
	pl.Pool = append(pl.Pool, token)
}

func (pl *TokenPool) Get(index int) Token {
	return pl.Pool[index]
}

func (pl *TokenPool) Len() int {
	return len(pl.Pool)
}

func (pl *TokenPool) Last() Token {
	if pl.Len() == 0 {
		return Token{}
	}
	return pl.Pool[pl.Len()-1]
}
