package ast

import (
	"CompilerInGo/lexer"
	"errors"
)

// AST中的结点类型
const (
	PROGRAM = iota
	METHOD
	RESULTTYPE
	IDTYPE
	PARAMLIST
	TYPE
	BLOCK
	STATEMENT
	CONDITIONALSTATEMENT
	LOOPSTATEMENT
	CALLSTATEMENT
	ASSIGNMENTSTATEMENT
	RETURNSTATEMENT
	BREAKSTATEMENT
	CONTINUESTATEMENT
	LOCALVARIABLEDECLARATION
	ACTPARAMLIST
	EXP
	CONDITIONALEXP
	TERM
	FACTOR
	RELATIONEXP
	COMPEXP
	CMPOP
)

// TypeString 结点类型对应的字符串
var TypeString = map[uint]string{
	PROGRAM:                  "Program",
	METHOD:                   "Method",
	RESULTTYPE:               "ResultType",
	IDTYPE:                   "ID",
	PARAMLIST:                "ParamList",
	TYPE:                     "Type",
	BLOCK:                    "Block",
	STATEMENT:                "Statement",
	CONDITIONALSTATEMENT:     "ConditionalStatement",
	LOOPSTATEMENT:            "LoopStatement",
	CALLSTATEMENT:            "CallStatement",
	ASSIGNMENTSTATEMENT:      "AssignStatement",
	RETURNSTATEMENT:          "ReturnStatement",
	BREAKSTATEMENT:           "BreakStatement",
	CONTINUESTATEMENT:        "ContinueStatement",
	LOCALVARIABLEDECLARATION: "LocalVariableDeclaration",
	ACTPARAMLIST:             "ActParamList",
	EXP:                      "Exp",
	CONDITIONALEXP:           "ConditionalExp",
	TERM:                     "Term",
	FACTOR:                   "Factor",
	RELATIONEXP:              "RelationExp",
	COMPEXP:                  "CompExp",
	CMPOP:                    "CmpOP",
}

// Program AST根结点
// Program→ Method { Method }
type Program struct {
	Method []Method
}

// Method 方法结点
// Method→ ResultType  ID  '('  ParamList  ')'  Block
type Method struct {
	ResultType ResultType
	ID         ID
	LParen     lexer.Token
	ParamList  ParamList
	RParen     lexer.Token
	Block      Block
}

// ResultType 方法返回值类型
// ResultType → 'integer' | 'float' | 'char' | 'string' | 'void'
type ResultType lexer.Token

// ID 标识符
type ID lexer.Token

// ParamListRest 参数列表可选部分
type ParamListRest struct {
	Comma lexer.Token
	Type  Type
	ID    ID
}

// ParamListField 参数列表不为空的必选部分
type ParamListField struct {
	Type          Type
	ID            ID
	ParamListRest *[]ParamListRest
}

// ParamList 参数列表
// ParamList→ Type ID { ',' Type ID } | ε
type ParamList struct {
	ParamList *ParamListField
}

// Type 变量类型
// Type→ 'integer' | 'float' | 'char' | 'string'
type Type lexer.Token

// Block 语句块
// Block→ '{' Statement { Statement } '}' | '{'  '}'
type Block struct {
	LBrace     lexer.Token
	Statements *[]Statement
	RBrace     lexer.Token
}

// Statement 语句
// Statement→ ConditionalStatement
//
//				  | LoopStatement
//	              | CallStatement
//	 		      | AssignmentStatement
//				  | ReturnStatement
//				  | BreakStatement
//				  | ContinueStatement
//				  | LocalVariableDeclaration
//				  | Block
//				  | ';'
//
// Statement的类型约束在创建AST时进行
type Statement struct {
	Type      uint
	Statement any
}

// LocalVariableDeclarationRest 局部变量声明可选部分
type LocalVariableDeclarationRest struct {
	Comma lexer.Token
	ID    ID
}

// LocalVariableDeclaration 局部变量声明
// LocalVariableDeclaration→Type ID { ',' ID } ';'
type LocalVariableDeclaration struct {
	Type                         Type
	ID                           ID
	LocalVariableDeclarationRest *[]LocalVariableDeclarationRest
	Semicolon                    lexer.Token
}

// CallStatement 函数调用语句
// CallStatement→ 'call' ID '(' ActParamList ')' ';'
type CallStatement struct {
	Call         lexer.Token
	ID           ID
	LParen       lexer.Token
	ActParamList ActParamList
	RParen       lexer.Token
	Semicolon    lexer.Token
}

// ActParamListRest 实参列表可选部分
type ActParamListRest struct {
	Comma lexer.Token
	Exp   Exp
}

// ActParamListField 实参列表不为空的必选部分
type ActParamListField struct {
	Exp              Exp
	ActParamListRest *[]ActParamListRest
}

// ActParamList 实参列表
// ActParamList→ ε | Exp { ',' Exp }
type ActParamList struct {
	ActParamList *ActParamListField
}

// AssignmentStatement 赋值语句
// AssignmentStatement→ ID '=' Exp ';'
type AssignmentStatement struct {
	ID        ID
	Assign    lexer.Token
	Exp       Exp
	Semicolon lexer.Token
}

// ConditionalStatement 条件语句
// ConditionalStatement→'if' '(' ConditionalExp ')' Statement [ 'else' Statement ]
type ConditionalStatement struct {
	If             lexer.Token
	LParen         lexer.Token
	ConditionalExp ConditionalExp
	RParen         lexer.Token
	Statement      Statement
	Else           *lexer.Token
	ElseStatement  *Statement
}

// LoopStatement 循环语句
// LoopStatement→'while' '(' ConditionalExp ')' Statement
type LoopStatement struct {
	While          lexer.Token
	LParen         lexer.Token
	ConditionalExp ConditionalExp
	RParen         lexer.Token
	Statement      Statement
}

// ReturnStatement 返回语句
// ReturnStatement→ 'return'  [ Exp ]  ';'
type ReturnStatement struct {
	Return    lexer.Token
	Exp       *Exp
	Semicolon lexer.Token
}

// BreakStatement 跳出语句
// BreakStatement→ 'break' ';'
type BreakStatement struct {
	Break     lexer.Token
	Semicolon lexer.Token
}

// ContinueStatement 继续语句
// ContinueStatement→ 'continue' ';'
type ContinueStatement struct {
	Continue  lexer.Token
	Semicolon lexer.Token
}

// ExpRest 表达式可选部分
type ExpRest struct {
	PlusOrMinus lexer.Token
	Term        Term
}

// Exp 表达式
// Exp→ Term  { '+' | '-'  Term }
type Exp struct {
	Term    Term
	ExpRest *ExpRest
}

// TermRest 单项可选部分
type TermRest struct {
	MulOrDiv lexer.Token
	Factor   Factor
}

// Term 单项
// Term→ Factor { '*' | '/'  Factor }
type Term struct {
	Factor   Factor
	TermRest *TermRest
}

// FactorTuple 括号表达式
type FactorTuple struct {
	LParen lexer.Token
	Exp    *Exp
	RParen lexer.Token
}

// Factor 单因子
// Factor→ ID | INTC | DECI | '(' Exp ')'
// Factor的类型约束在创建AST时进行
type Factor struct {
	Factor any
}

// ConditionalExpRest 条件表达式可选部分
type ConditionalExpRest struct {
	Or          lexer.Token
	RelationExp RelationExp
}

// ConditionalExp 条件表达式
// ConditionalExp→RelationExp { 'or' RelationExp }
type ConditionalExp struct {
	RelationExp        RelationExp
	ConditionalExpRest *ConditionalExpRest
}

// RelationExpRest 关系表达式可选部分
type RelationExpRest struct {
	And     lexer.Token
	CompExp CompExp
}

// RelationExp 关系表达式
// RelationExp→ CompExp { 'and' CompExp }
type RelationExp struct {
	CompExp         CompExp
	RelationExpRest *RelationExpRest
}

// CompExp 比较表达式
// CompExp→ Exp CmpOp Exp
type CompExp struct {
	LExp  Exp
	CmpOp CmpOp
	RExp  Exp
}

// CmpOp 比较运算符
// CmpOp→'<' | '<=' | '>' | '>=| '==' | '<>'
type CmpOp lexer.Token

// IsResultType 判断是否为返回值类型
func IsResultType(token lexer.Token) bool {
	return token.Type == lexer.INT || token.Type == lexer.FLOAT || token.Type == lexer.CHAR || token.Type == lexer.STRING || token.Type == lexer.VOID
}

// IsType 判断是否为变量类型
func IsType(token lexer.Token) bool {
	return token.Type == lexer.INT || token.Type == lexer.FLOAT || token.Type == lexer.CHAR || token.Type == lexer.STRING
}

// IsID 判断是否为标识符
func IsID(token lexer.Token) bool {
	return token.Type == lexer.IDENTIFIER
}

// NewProgram 创建Program
func NewProgram() (Program, error) {
	return Program{
		Method: make([]Method, 0), // 初始化Method数组
	}, nil
}

// NewMethod 创建Method
// resultType: 返回值类型
// idToken: 函数名
// lParen: 左括号
// paramList: 参数列表
// rParen: 右括号
// block: 函数体
func NewMethod(resultType ResultType, idToken ID, lParen lexer.Token, paramList ParamList, rParen lexer.Token, block Block) (Method, error) {
	// 检查参数是否合法
	if idToken.Type != lexer.IDENTIFIER {
		return Method{}, errors.New("Method: invalid idToken")

	}
	if lParen.Type != lexer.LPAREN {
		return Method{}, errors.New("Method: invalid lParen")
	}
	if rParen.Type != lexer.RPAREN {
		return Method{}, errors.New("Method: invalid rParen")
	}
	return Method{
		ResultType: resultType,
		ID:         idToken,
		LParen:     lParen,
		ParamList:  paramList,
		RParen:     rParen,
		Block:      block,
	}, nil
}

// NewResultType 创建返回值类型
func NewResultType(typeToken lexer.Token) (ResultType, error) {
	switch typeToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.CHAR, lexer.STRING, lexer.VOID:
		// 检查参数是否合法
		return ResultType(typeToken), nil
	default:
		return ResultType{}, errors.New("ResultType: invalid type")
	}
}

// NewParamList 创建参数列表
// param: 不定长参数列表
func NewParamList(param ...any) (ParamList, error) {
	// 没有参数
	if len(param) == 0 {
		return ParamList{}, nil
	}
	// 需要Type ID对，返回错误
	if len(param) == 1 {
		return ParamList{}, errors.New("ParamList: expect Type ID pair")
	}
	// 只有一个Type ID对
	if len(param) == 2 {
		switch param[0].(type) {
		// 第一个参数为Type
		case Type:
			switch param[1].(type) {
			// 第二个参数为ID
			case ID:
				return ParamList{
					ParamList: &ParamListField{
						Type: param[0].(Type),
						ID:   param[1].(ID),
					},
				}, nil
			default:
				return ParamList{}, errors.New("ParamList: expect Type ID pair")
			}
		default:
			return ParamList{}, errors.New("ParamList: expect Type ID pair")
		}
	}
	// 有多个Type ID对
	if (len(param)-2)%3 == 0 {
		// 使用ParamListRest存储除第一个外的Type ID对
		paramListRest := make([]ParamListRest, 0)
		// 剩下的Type ID对
		paramRest := param[2:]
		// 用于存储单个Type ID对
		var paramListRestSingle ParamListRest
		// 计数存储的Type ID对个数
		restCounter := 0
		// 遍历剩下的Type ID对
		for index, elem := range paramRest {
			switch elem.(type) {
			// 第k个参数为逗号
			case lexer.Token:
				if index%3 == 0 && elem.(lexer.Token).Type == lexer.COMMA {
					paramListRestSingle.Comma = elem.(lexer.Token)
				} else {
					return ParamList{}, errors.New("ParamList: expect comma")
				}
			// 第k+1个参数为Type
			case Type:
				if index%3 == 1 {
					paramListRestSingle.Type = elem.(Type)
				} else {
					return ParamList{}, errors.New("ParamList: expect Type")
				}
			// 第k+2个参数为ID
			case ID:
				if index%3 == 2 {
					paramListRestSingle.ID = elem.(ID)
					// 将单个Type ID对存入paramListRest
					paramListRest = append(paramListRest, paramListRestSingle)
					restCounter++
				} else {
					return ParamList{}, errors.New("ParamList: expect ID")
				}
			}
		}
		// 检查存储的Type ID对个数是否正确
		if restCounter != (len(param)-2)/3 {
			return ParamList{}, errors.New("ParamList: expect Type ID pair")
		}
		return ParamList{
			ParamList: &ParamListField{
				Type:          (param[0]).(Type),
				ID:            (param[1]).(ID),
				ParamListRest: &paramListRest,
			},
		}, nil
	} else {
		return ParamList{}, errors.New("ParamList: expect Type ID pair")
	}
}

// NewType 创建类型
func NewType(typeToken lexer.Token) (Type, error) {
	switch typeToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.CHAR, lexer.STRING:
		return Type(typeToken), nil
	default:
		return Type{}, errors.New("Type: invalid type")
	}
}

// NewBlock 创建代码块
// lBrace: 左大括号
// statements: 语句列表
// rBrace: 右大括号
func NewBlock(lBrace lexer.Token, statements []Statement, rBrace lexer.Token) (Block, error) {
	// 检查大括号是否合法
	if lBrace.Type != lexer.LBRACE {
		return Block{}, errors.New("Block: invalid lBrace")
	}
	if rBrace.Type != lexer.RBRACE {
		return Block{}, errors.New("Block: invalid rBrace")
	}
	// 是否有语句
	if len(statements) == 0 {
		return Block{
			LBrace: lBrace,
			RBrace: rBrace,
		}, nil
	}
	// 语句列表不为空
	return Block{
		LBrace:     lBrace,
		Statements: &statements,
		RBrace:     rBrace,
	}, nil
}

// NewLocalVariableDeclarationStatement  创建局部变量声明
// typeToken: 类型
// idToken: 标识符
// localVariableDeclarationRest: 局部变量声明的后续部分
// semicolonToken: 分号
func NewLocalVariableDeclarationStatement(typeToken Type, idToken ID, localVariableDeclarationRest []any, semicolonToken lexer.Token) (LocalVariableDeclaration, error) {
	// 检查类型、分号是否合法
	if idToken.Type != lexer.IDENTIFIER {
		return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid id token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid semicolon token")
	}

	// 检查局部变量声明的后续部分是否合法
	// 逗号 +  标识符 必须成对出现
	if len(localVariableDeclarationRest)%2 != 0 {
		return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid number of tokens")
	}

	// 没有后续部分
	if len(localVariableDeclarationRest) == 0 {
		return LocalVariableDeclaration{
			Type:      typeToken,
			ID:        idToken,
			Semicolon: semicolonToken,
		}, nil
	}

	// 有后续部分 用rest保存后续部分
	rest := make([]LocalVariableDeclarationRest, 0)

	// 用于存储单个逗号 + 标识符
	var restSingle LocalVariableDeclarationRest
	restCounter := 0
	// 遍历后续部分
	for index, elem := range localVariableDeclarationRest {
		// 清空restSingle
		if index%2 == 0 {
			restSingle = LocalVariableDeclarationRest{}
		}

		switch elem.(type) {
		case lexer.Token:
			// 第k个元素为逗号
			if index%2 != 0 {
				return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid comma token")
			} else {
				restSingle.Comma = elem.(lexer.Token)
			}
		case ID:
			// 第k+1个元素为标识符
			if index%2 == 0 {
				return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid id token")
			} else {
				restSingle.ID = elem.(ID)
				// 将单个逗号 + 标识符对存入rest
				rest = append(rest, restSingle)
				restCounter++
			}

		}
	}

	return LocalVariableDeclaration{
		Type:                         typeToken,
		ID:                           idToken,
		LocalVariableDeclarationRest: &rest,
	}, nil
}

// NewCallStatement 创建调用语句
// callToken: call
// idToken: 标识符
// lParen: 左括号
// actParamList: 实参列表
// rParen: 右括号
// semicolonToken: 分号
func NewCallStatement(callToken lexer.Token, idToken ID, lParen lexer.Token, actParamList *ActParamList, rParen, semicolonToken lexer.Token) (CallStatement, error) {
	// 检查call、标识符、左右括号、分号是否合法
	if callToken.Type != lexer.CALL {
		return CallStatement{}, errors.New("CallStatement: invalid call token")
	}
	if idToken.Type != lexer.IDENTIFIER {
		return CallStatement{}, errors.New("CallStatement: invalid id token")
	}
	if lParen.Type != lexer.LPAREN {
		return CallStatement{}, errors.New("CallStatement: invalid lParen token")
	}
	if rParen.Type != lexer.RPAREN {
		return CallStatement{}, errors.New("CallStatement: invalid rParen token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return CallStatement{}, errors.New("CallStatement: invalid semicolon token")
	}
	return CallStatement{
		Call:         callToken,
		ID:           idToken,
		LParen:       lParen,
		ActParamList: *actParamList,
		RParen:       rParen,
		Semicolon:    semicolonToken,
	}, nil
}

// NewActParamList 创建实参列表
// token: 实参列表的token(空或Exp表达式)
func NewActParamList(token []any) (ActParamList, error) {
	// 实参列表为空
	if token == nil {
		return ActParamList{}, nil
	}
	// 实参列表不为空
	switch token[0].(type) {
	case Exp:
		// 只有一个表达式
		if len(token) == 1 {
			return ActParamList{
				ActParamList: &ActParamListField{
					Exp:              token[0].(Exp),
					ActParamListRest: nil,
				},
			}, nil
		}
		// 有多个表达式
		// 逗号 + 表达式 必须成对出现
		if (len(token)-1)%2 == 0 {
			// 用于存储逗号 + 表达式对
			actParamListRest := make([]ActParamListRest, 0)
			// 用于存储单个逗号 + 表达式对
			var actParamListRestSingle ActParamListRest
			// 遍历逗号 + 表达式对
			for index, elem := range token[1:] {
				//第k个元素为逗号
				if index%2 == 0 {
					actParamListRestSingle.Comma = elem.(lexer.Token)
					//第k+1个元素为表达式
				} else if index%2 == 1 {
					actParamListRestSingle.Exp = elem.(Exp)
					// 将单个逗号 + 表达式对存入actParamListRest
					actParamListRest = append(actParamListRest, actParamListRestSingle)
				}
			}
			return ActParamList{
				ActParamList: &ActParamListField{
					Exp:              token[0].(Exp),
					ActParamListRest: &actParamListRest,
				},
			}, nil
		}
	default:
		// 实参列表不为空，但是第一个元素不是Exp表达式
		return ActParamList{}, errors.New("ActParamList: invalid token")
	}
	return ActParamList{}, errors.New("ActParamList: invalid token")
}

// NewAssignmentStatement 创建赋值语句
// idToken: 标识符
// equalToken: 等号
// exp: 表达式
// semicolonToken: 分号
func NewAssignmentStatement(idToken ID, equalToken lexer.Token, exp Exp, semicolonToken lexer.Token) (AssignmentStatement, error) {
	// 检查标识符、等号、分号是否合法
	if idToken.Type != lexer.IDENTIFIER {
		return AssignmentStatement{}, errors.New("AssignmentStatement: invalid id token")
	}
	if equalToken.Type != lexer.ASSIGN {
		return AssignmentStatement{}, errors.New("AssignmentStatement: invalid equal token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return AssignmentStatement{}, errors.New("AssignmentStatement: invalid semicolon token")
	}

	return AssignmentStatement{
		ID:        idToken,
		Assign:    equalToken,
		Exp:       exp,
		Semicolon: semicolonToken,
	}, nil
}

// NewConditionalStatement 创建条件语句
// ifToken: if
// lParen: 左括号
// conditionalExp: 条件表达式
// rParen: 右括号
// statement: 语句
// elseToken: else（可选）
// elseStatement: else语句（可选）
func NewConditionalStatement(ifToken, lParen lexer.Token, conditionalExp ConditionalExp, rParen lexer.Token, statement Statement, elseToken *lexer.Token, elseStatement *Statement) (ConditionalStatement, error) {
	// 检查if、左右括号
	if ifToken.Type != lexer.IF {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid if token")
	}
	if lParen.Type != lexer.LPAREN {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid lParen token")
	}
	if rParen.Type != lexer.RPAREN {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid rParen token")
	}

	// 是否有else
	if elseToken != nil && elseToken.Type != lexer.ELSE {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid else token")
	}
	// else与elseStatement必须同时存在或同时不存在
	if (elseStatement != nil && elseToken == nil) || (elseToken != nil && elseStatement == nil) {
		return ConditionalStatement{}, errors.New("ConditionalStatement: expected else token and else statement to be both nil or both not nil")
	}

	return ConditionalStatement{
		If:             ifToken,
		LParen:         lParen,
		ConditionalExp: conditionalExp,
		RParen:         rParen,
		Statement:      statement,
		Else:           elseToken,
		ElseStatement:  elseStatement,
	}, nil
}

// NewLoopStatement 创建循环语句
// whileToken: while
// lParen: 左括号
// conditionalExp: 条件表达式
// rParen: 右括号
// statement: 语句
func NewLoopStatement(whileToken, lParen lexer.Token, conditionalExp ConditionalExp, rParen lexer.Token, statement Statement) (LoopStatement, error) {
	// 检查while、左右括号
	if whileToken.Type != lexer.WHILE {
		return LoopStatement{}, errors.New("LoopStatement: invalid while token")
	}
	if lParen.Type != lexer.LPAREN {
		return LoopStatement{}, errors.New("LoopStatement: invalid lParen token")
	}
	if rParen.Type != lexer.RPAREN {
		return LoopStatement{}, errors.New("LoopStatement: invalid rParen token")
	}
	return LoopStatement{
		While:          whileToken,
		LParen:         lParen,
		ConditionalExp: conditionalExp,
		RParen:         rParen,
		Statement:      statement,
	}, nil
}

// NewCmpOp 创建比较运算符
func NewCmpOp(cmpOpToken lexer.Token) (CmpOp, error) {
	switch cmpOpToken.Type {
	case lexer.LESS, lexer.LESSEQUAL, lexer.GREATER, lexer.GREATEREQUAL, lexer.EQUAL, lexer.DIAMOND:
		return CmpOp(cmpOpToken), nil
	default:
		return CmpOp{}, errors.New("CmpOp: invalid type")
	}
}

// NewReturnStatement 创建返回语句
// returnToken: return
// exp: 表达式（可选）
// semicolonToken: 分号
func NewReturnStatement(returnToken lexer.Token, exp *Exp, semicolonToken lexer.Token) (ReturnStatement, error) {
	// 检查return、分号
	if returnToken.Type != lexer.RETURN {
		return ReturnStatement{}, errors.New("ReturnStatement: invalid return token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return ReturnStatement{}, errors.New("ReturnStatement: invalid semicolon token")
	}

	// 是否返回表达式
	if exp != nil {
		return ReturnStatement{
			Return:    returnToken,
			Exp:       exp,
			Semicolon: semicolonToken,
		}, nil
	}

	return ReturnStatement{
		Return:    returnToken,
		Semicolon: semicolonToken,
	}, nil
}

// NewBreakStatement 创建break语句
// breakToken: break
// semicolonToken: 分号
func NewBreakStatement(breakToken, semicolonToken lexer.Token) (BreakStatement, error) {
	// 检查break、分号
	if breakToken.Type != lexer.BREAK {
		return BreakStatement{}, errors.New("BreakStatement: invalid break token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return BreakStatement{}, errors.New("BreakStatement: invalid semicolon token")
	}

	return BreakStatement{
		Break:     breakToken,
		Semicolon: semicolonToken,
	}, nil
}

// NewContinueStatement 创建continue语句
// continueToken: continue
// semicolonToken: 分号
func NewContinueStatement(continueToken, semicolonToken lexer.Token) (ContinueStatement, error) {
	// 检查continue、分号
	if continueToken.Type != lexer.CONTINUE {
		return ContinueStatement{}, errors.New("ContinueStatement: invalid continue token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return ContinueStatement{}, errors.New("ContinueStatement: invalid semicolon token")
	}

	return ContinueStatement{
		Continue:  continueToken,
		Semicolon: semicolonToken,
	}, nil
}

// NewExp 创建表达式
// lTerm: 左项
// plusOrMinus: 加减号（可选）
// rTerm: 右项（可选）
func NewExp(lTerm Term, plusOrMinus *lexer.Token, rTerm *Term) (Exp, error) {
	// 加减号和右项必须同时存在或同时不存在
	if (plusOrMinus == nil && rTerm != nil) || (plusOrMinus != nil && rTerm == nil) {
		return Exp{}, errors.New("Exp: expected plusOrMinus and rTerm to be both nil or both not nil")
	}
	// 检查加减号
	if plusOrMinus != nil && plusOrMinus.Type != lexer.PLUS && plusOrMinus.Type != lexer.MINUS {
		return Exp{}, errors.New("Exp: invalid plusOrMinus token")
	}
	// 如果加减号和右项都不存在，则返回左项
	if plusOrMinus == nil && rTerm == nil {
		return Exp{
			Term: lTerm,
		}, nil
	}
	return Exp{
		Term: lTerm,
		ExpRest: &ExpRest{
			PlusOrMinus: *plusOrMinus,
			Term:        *rTerm,
		},
	}, nil
}

// NewTerm 创建项
// lFactor: 左因子
// mulOrDiv: 乘除号（可选）
// rFactor: 右因子（可选）
func NewTerm(lFactor Factor, mulOrDiv *lexer.Token, rFactor *Factor) (Term, error) {
	// 乘除号和右因子必须同时存在或同时不存在
	if (mulOrDiv == nil && rFactor != nil) || (mulOrDiv != nil && rFactor == nil) {
		return Term{}, errors.New("Term: expected mulOrDiv and rFactor to be both nil or both not nil")
	}
	// 检查乘除号
	if mulOrDiv != nil && mulOrDiv.Type != lexer.TIMES && mulOrDiv.Type != lexer.DIVIDE {
		return Term{}, errors.New("Term: invalid mulOrDiv token")
	}
	// 如果乘除号和右因子都不存在，则返回左因子
	if mulOrDiv == nil && rFactor == nil {
		return Term{
			Factor: lFactor,
		}, nil
	}
	return Term{
		Factor: lFactor,
		TermRest: &TermRest{
			MulOrDiv: *mulOrDiv,
			Factor:   *rFactor,
		},
	}, nil
}

// NewFactor 创建因子
// factor: 不定长度因子
func NewFactor(factor ...any) (Factor, error) {
	// 空因子，返回错误
	if len(factor) == 0 {
		return Factor{}, errors.New("Factor: expected at least one factor")
	} else if len(factor) == 1 {
		// 一个因子，判断类型
		switch factor[0].(type) {
		case lexer.Token, ID:
			// 如果是标识符或者token，但不是IDENTIFIER、INTEGER_LITERAL、DECIMAL_LITERAL，返回错误
			if factor[0].(lexer.Token).Type != lexer.IDENTIFIER && factor[0].(lexer.Token).Type != lexer.INTEGER_LITERAL && factor[0].(lexer.Token).Type != lexer.DECIMAL_LITERAL {
				return Factor{}, errors.New("Factor: invalid token")
			}
			// 返回ID、INTEGER_LITERAL、DECIMAL_LITERAL类型的因子
			return Factor{
				Factor: factor[0].(lexer.Token),
			}, nil
		}
	} else if len(factor) == 3 {
		// '(' Exp ')'类型的因子
		switch factor[0].(type) {
		case lexer.Token:
			// 检查左括号
			if factor[0].(lexer.Token).Type != lexer.LPAREN {
				return Factor{}, errors.New("Factor: invalid lParen token")
			}
			switch factor[2].(type) {
			case lexer.Token:
				// 检查右括号
				if factor[2].(lexer.Token).Type != lexer.RPAREN {
					return Factor{}, errors.New("Factor: invalid rParen token")
				}
				switch factor[1].(type) {
				case *Exp:
					// 中间为表达式 合法
					return Factor{
						Factor: FactorTuple{
							LParen: factor[0].(lexer.Token),
							Exp:    factor[1].(*Exp),
							RParen: factor[2].(lexer.Token),
						},
					}, nil
				default:
					return Factor{}, errors.New("Factor: invalid exp")
				}
			default:
				return Factor{}, errors.New("Factor: invalid rParen token")
			}
		default:
			return Factor{}, errors.New("Factor: invalid lParen token")
		}
	}
	return Factor{}, nil
}

// NewConditionalExp 创建条件表达式
// lExp: 左表达式
// Or: or（可选）
// rExp: 右表达式（可选）
func NewConditionalExp(lExp RelationExp, Or *lexer.Token, rExp *RelationExp) (ConditionalExp, error) {
	// or和右表达式必须同时存在或同时不存在
	if (Or == nil && rExp != nil) || (Or != nil && rExp == nil) {
		return ConditionalExp{}, errors.New("NewConditionalExp: expected Or and rExp to be both nil or both not nil")
	}
	// 检查or
	if Or != nil && Or.Type != lexer.OR {
		return ConditionalExp{}, errors.New("NewConditionalExp: invalid Or token")
	}
	// 如果or和右表达式都不存在，则返回左表达式
	if Or == nil && rExp == nil {
		return ConditionalExp{
			RelationExp: lExp,
		}, nil
	}
	return ConditionalExp{
		RelationExp: lExp,
		ConditionalExpRest: &ConditionalExpRest{
			Or:          *Or,
			RelationExp: *rExp,
		},
	}, nil
}

// NewRelationExp 创建关系表达式
// lExp: 左表达式
// And: and（可选）
// rExp: 右表达式（可选）
func NewRelationExp(lExp CompExp, And *lexer.Token, rExp *CompExp) (RelationExp, error) {
	// and和右表达式必须同时存在或同时不存在
	if (And == nil && rExp != nil) || (And != nil && rExp == nil) {
		return RelationExp{}, errors.New("NewRelationExp: expected And and rExp to be both nil or both not nil")
	}
	// 检查and
	if And != nil && And.Type != lexer.OR {
		return RelationExp{}, errors.New("NewRelationExp: invalid And token")
	}
	// 如果and和右表达式都不存在，则返回左表达式
	if And == nil && rExp == nil {
		return RelationExp{
			CompExp: lExp,
		}, nil
	}
	return RelationExp{
		CompExp: lExp,
		RelationExpRest: &RelationExpRest{
			And:     *And,
			CompExp: *rExp,
		},
	}, nil
}

// NewCompExp 创建比较表达式
// lExp: 左表达式
// cmpOp: 比较运算符
// rExp: 右表达式
func NewCompExp(lExp Exp, cmpOp CmpOp, rExp Exp) (CompExp, error) {
	// < | <= | > | >= | == | <>
	if cmpOp.Type != lexer.LESS && cmpOp.Type != lexer.LESSEQUAL && cmpOp.Type != lexer.GREATER && cmpOp.Type != lexer.GREATEREQUAL && cmpOp.Type != lexer.EQUAL && cmpOp.Type != lexer.DIAMOND {
		return CompExp{}, errors.New("NewCompExp: invalid cmpOp token")
	}
	return CompExp{
		LExp:  lExp,
		CmpOp: cmpOp,
		RExp:  rExp,
	}, nil
}
