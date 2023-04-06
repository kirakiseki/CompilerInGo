package lexer

import "CompilerInGo/utils"

type TokenCategory uint
type TokenType uint

type Token struct {
	Category TokenCategory
	Type     TokenType
	Literal  string
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
)

// TokenType
// 关键字
const (
	VOID   = 9 + iota //9 void
	VAR               //10 var
	INT               //11 int
	FLOAT             //12 float
	STRING            //13 string
	BEGIN             //14 begin
	END               //15 end
	IF                //16 if
	THEN              //17 then
	ELSE              //18 else
	WHILE             //19 while
	DO                //20 do
	CALL              //21 call
	READ              //22 read
	WRITE             //23 write
	AND               //24 and
	OR                //25 or
)

// 分隔符
const (
	LBRACE    = 26 + iota //26 {
	RBRACE                //27 }
	LPAREN                //28 (
	RPAREN                //29 )
	SEMICOLON             //30 ;
	SPACE                 //31 空格
)

// 运算符
const (
	EQUAL        = 32 + iota //32 ==
	ASSIGN                   //33 =
	LESS                     //34 <
	LESSEQUAL                //35 <=
	GREATER                  //36 >
	GREATEREQUAL             //37 >=
	LESSGREATER              //38 <>
	PLUS                     //39 +
	MINUS                    //40 -
	TIMES                    //41 *
	DIVIDE                   //42 /
)

// 字面量
const (
	INTEGER_LITERAL = 43 + iota //43 整数字面量
	DECIMAL_LITERAL             //44 小数字面量
	STRING_LITERAL              //45 字符串字面量
	CHAR_LITERAL                //46 字符字面量
)

var tokenString = map[TokenType]string{
	EOF: "EOF",

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
	SPACE:     " ",

	EQUAL:        "==",
	ASSIGN:       "=",
	LESS:         "<",
	LESSEQUAL:    "<=",
	GREATER:      ">",
	GREATEREQUAL: ">=",
	LESSGREATER:  "<>",
	PLUS:         "+",
	MINUS:        "-",
	TIMES:        "*",
	DIVIDE:       "/",

	INTEGER_LITERAL: "INTEGER_LITERAL",
	DECIMAL_LITERAL: "DECIMAL_LITERAL",
	STRING_LITERAL:  "STRING_LITERAL",
	CHAR_LITERAL:    "CHAR_LITERAL",
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

func (t Token) String() string {
	return tokenString[t.Type]
}

func (t Token) setType() {
	for k, v := range tokenString {
		if v == t.Literal {
			t.Type = k
		}
	}
}

func (t Token) setCategory() {
	switch t.Type {
	case VOID, VAR, INT, FLOAT, STRING, BEGIN, END, IF, THEN, ELSE, WHILE, DO, CALL, READ, WRITE, AND, OR:
		t.Category = KEYWORD
	case LBRACE, RBRACE, LPAREN, RPAREN, SEMICOLON, SPACE:
		t.Category = DELIM
	case EQUAL, ASSIGN, LESS, LESSEQUAL, GREATER, GREATEREQUAL, LESSGREATER, PLUS, MINUS, TIMES, DIVIDE:
		t.Category = OPERA
	case EOF:
		t.Category = EOF
	case INTEGER_LITERAL:
		t.Category = INTEGER
	case DECIMAL_LITERAL:
		t.Category = DECIMAL
	case STRING_LITERAL:
		t.Category = STR
	case CHAR_LITERAL:
		t.Category = CHAR
	default:
		t.Category = IDENT
	}
}

func NewToken(literal string, position utils.PositionPair) Token {
	token := Token{
		Literal: literal,
		Pos:     position,
	}
	token.setType()
	token.setCategory()
	return token
}

func NewLiteralToken(literal string, position utils.PositionPair, tokenType TokenType) Token {
	token := Token{
		Literal: literal,
		Pos:     position,
		Type:    tokenType,
	}
	token.setCategory()
	return token
}
