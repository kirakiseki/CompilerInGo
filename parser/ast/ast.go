package ast

import (
	"CompilerInGo/lexer"
	"errors"
)

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

type Program struct {
	Method []Method
}

type Method struct {
	ResultType ResultType
	ID         ID
	LParen     lexer.Token
	ParamList  ParamList
	RParen     lexer.Token
	Block      Block
}
type ResultType lexer.Token
type ID lexer.Token

type ParamListRest struct {
	Comma lexer.Token
	Type  Type
	ID    ID
}

type ParamListField struct {
	Type          Type
	ID            ID
	ParamListRest *[]ParamListRest
}

type ParamList struct {
	ParamList *ParamListField
}

type Type lexer.Token

type Block struct {
	LBrace     lexer.Token
	Statements *[]Statement
	RBrace     lexer.Token
}

type Statement struct {
	Type      uint
	Statement any
}

type LocalVariableDeclarationRest struct {
	Comma lexer.Token
	ID    ID
}

type LocalVariableDeclaration struct {
	Type                         Type
	ID                           ID
	LocalVariableDeclarationRest *[]LocalVariableDeclarationRest
	Semicolon                    lexer.Token
}

type CallStatement struct {
	Call         lexer.Token
	ID           ID
	LParen       lexer.Token
	ActParamList ActParamList
	RParen       lexer.Token
	Semicolon    lexer.Token
}

type ActParamListRest struct {
	Comma lexer.Token
	Exp   Exp
}

type ActParamListField struct {
	Exp              Exp
	ActParamListRest *[]ActParamListRest
}

type ActParamList struct {
	ActParamList *ActParamListField
}

type AssignmentStatement struct {
	ID        ID
	Assign    lexer.Token
	Exp       Exp
	Semicolon lexer.Token
}

type ConditionalStatement struct {
	If             lexer.Token
	LParen         lexer.Token
	ConditionalExp ConditionalExp
	RParen         lexer.Token
	Statement      Statement
	Else           *lexer.Token
	ElseStatement  *Statement
}

type LoopStatement struct {
	While          lexer.Token
	LParen         lexer.Token
	ConditionalExp ConditionalExp
	RParen         lexer.Token
	Statement      Statement
}

type ReturnStatement struct {
	Return    lexer.Token
	Exp       *Exp
	Semicolon lexer.Token
}

type BreakStatement struct {
	Break     lexer.Token
	Semicolon lexer.Token
}

type ContinueStatement struct {
	Continue  lexer.Token
	Semicolon lexer.Token
}

type ExpRest struct {
	PlusOrMinus lexer.Token
	Term        Term
}

type Exp struct {
	Term    Term
	ExpRest *ExpRest
}

type TermRest struct {
	MulOrDiv lexer.Token
	Factor   Factor
}

type Term struct {
	Factor   Factor
	TermRest *TermRest
}

type FactorTuple struct {
	LParen lexer.Token
	Exp    *Exp
	RParen lexer.Token
}

type Factor struct {
	Factor any
}

type ConditionalExpRest struct {
	Or          lexer.Token
	RelationExp RelationExp
}

type ConditionalExp struct {
	RelationExp        RelationExp
	ConditionalExpRest *ConditionalExpRest
}

type RelationExpRest struct {
	And     lexer.Token
	CompExp CompExp
}

type RelationExp struct {
	CompExp         CompExp
	RelationExpRest *RelationExpRest
}

type CompExp struct {
	LExp  Exp
	CmpOp CmpOp
	RExp  Exp
}

type CmpOp lexer.Token

func IsResultType(token lexer.Token) bool {
	return token.Type == lexer.INT || token.Type == lexer.FLOAT || token.Type == lexer.CHAR || token.Type == lexer.STRING || token.Type == lexer.VOID
}

func IsType(token lexer.Token) bool {
	return token.Type == lexer.INT || token.Type == lexer.FLOAT || token.Type == lexer.CHAR || token.Type == lexer.STRING
}

func IsID(token lexer.Token) bool {
	return token.Type == lexer.IDENTIFIER
}

func NewProgram() (Program, error) {
	return Program{
		Method: make([]Method, 0),
	}, nil
}

func NewMethod(resultType ResultType, idToken ID, lParen lexer.Token, paramList ParamList, rParen lexer.Token, block Block) (Method, error) {
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

func NewResultType(typeToken lexer.Token) (ResultType, error) {
	switch typeToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.CHAR, lexer.STRING, lexer.VOID:
		return ResultType(typeToken), nil
	default:
		return ResultType{}, errors.New("ResultType: invalid type")
	}
}

func NewParamList(param ...any) (ParamList, error) {
	if len(param) == 0 {
		return ParamList{}, nil
	}
	if len(param) == 1 {
		return ParamList{}, errors.New("ParamList: expect Type ID pair")
	}
	if len(param) == 2 {
		switch param[0].(type) {
		case Type:
			switch param[1].(type) {
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
	if (len(param)-2)%3 == 0 {
		paramListRest := make([]ParamListRest, 0)
		paramRest := param[2:]
		var paramListRestSingle ParamListRest
		restCounter := 0
		for index, elem := range paramRest {
			switch elem.(type) {
			case lexer.Token:
				if index%3 == 0 && elem.(lexer.Token).Type == lexer.COMMA {
					paramListRestSingle.Comma = elem.(lexer.Token)
				} else {
					return ParamList{}, errors.New("ParamList: expect comma")
				}
			case Type:
				if index%3 == 1 {
					paramListRestSingle.Type = elem.(Type)
				} else {
					return ParamList{}, errors.New("ParamList: expect Type")
				}
			case ID:
				if index%3 == 2 {
					paramListRestSingle.ID = elem.(ID)
					paramListRest = append(paramListRest, paramListRestSingle)
					restCounter++
				} else {
					return ParamList{}, errors.New("ParamList: expect ID")
				}
			}
		}
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

func NewType(typeToken lexer.Token) (Type, error) {
	switch typeToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.CHAR, lexer.STRING:
		return Type(typeToken), nil
	default:
		return Type{}, errors.New("Type: invalid type")
	}
}

func NewBlock(lBrace lexer.Token, statements []Statement, rBrace lexer.Token) (Block, error) {
	if lBrace.Type != lexer.LBRACE {
		return Block{}, errors.New("Block: invalid lBrace")
	}
	if rBrace.Type != lexer.RBRACE {
		return Block{}, errors.New("Block: invalid rBrace")
	}
	if len(statements) == 0 {
		return Block{
			LBrace: lBrace,
			RBrace: rBrace,
		}, nil
	}
	return Block{
		LBrace:     lBrace,
		Statements: &statements,
		RBrace:     rBrace,
	}, nil
}

func NewLocalDeclarationStatement(typeToken Type, idToken ID, localVariableDeclarationRest []any, semicolonToken lexer.Token) (LocalVariableDeclaration, error) {
	if idToken.Type != lexer.IDENTIFIER {
		return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid id token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid semicolon token")
	}

	if len(localVariableDeclarationRest)%2 != 0 {
		return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid number of tokens")
	}

	if len(localVariableDeclarationRest) == 0 {
		return LocalVariableDeclaration{
			Type:      typeToken,
			ID:        idToken,
			Semicolon: semicolonToken,
		}, nil
	}

	rest := make([]LocalVariableDeclarationRest, 0)

	var restSingle LocalVariableDeclarationRest
	restCounter := 0
	for index, elem := range localVariableDeclarationRest {
		if index%2 == 0 {
			restSingle = LocalVariableDeclarationRest{}
		}

		switch elem.(type) {
		case lexer.Token:
			if index%2 != 0 {
				return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid comma token")
			} else {
				restSingle.Comma = elem.(lexer.Token)
			}
		case ID:
			if index%2 == 0 {
				return LocalVariableDeclaration{}, errors.New("LocalVariableDeclaration: invalid id token")
			} else {
				restSingle.ID = elem.(ID)
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

func NewCallStatement(callToken lexer.Token, idToken ID, lParen lexer.Token, actParamList *ActParamList, rParen, semicolonToken lexer.Token) (CallStatement, error) {
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

func NewActParamList(token []any) (ActParamList, error) {
	if token == nil {
		return ActParamList{}, nil
	}
	switch token[0].(type) {
	case Exp:
		if len(token) == 1 {
			return ActParamList{
				ActParamList: &ActParamListField{
					Exp:              token[0].(Exp),
					ActParamListRest: nil,
				},
			}, nil
		}
		if (len(token)-1)%2 == 0 {
			actParamListRest := make([]ActParamListRest, 0)
			var actParamListRestSingle ActParamListRest
			for index, elem := range token[1:] {
				if index%2 == 0 {
					actParamListRestSingle.Comma = elem.(lexer.Token)
				} else if index%2 == 1 {
					actParamListRestSingle.Exp = elem.(Exp)
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
		return ActParamList{}, errors.New("ActParamList: invalid token")
	}
	return ActParamList{}, errors.New("ActParamList: invalid token")
}

func NewAssignmentStatement(idToken ID, equalToken lexer.Token, exp Exp, semicolonToken lexer.Token) (AssignmentStatement, error) {
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

func NewConditionalStatement(ifToken, lParen lexer.Token, conditionalExp ConditionalExp, rParen lexer.Token, statement Statement, elseToken *lexer.Token, elseStatement *Statement) (ConditionalStatement, error) {
	if ifToken.Type != lexer.IF {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid if token")
	}
	if lParen.Type != lexer.LPAREN {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid lParen token")
	}
	if rParen.Type != lexer.RPAREN {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid rParen token")
	}
	if elseToken != nil && elseToken.Type != lexer.ELSE {
		return ConditionalStatement{}, errors.New("ConditionalStatement: invalid else token")
	}
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

func NewLoopStatement(whileToken, lParen lexer.Token, conditionalExp ConditionalExp, rParen lexer.Token, statement Statement) (LoopStatement, error) {
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

func NewCmpOp(cmpOpToken lexer.Token) (CmpOp, error) {
	switch cmpOpToken.Type {
	case lexer.LESS, lexer.LESSEQUAL, lexer.GREATER, lexer.GREATEREQUAL, lexer.EQUAL, lexer.DIAMOND:
		return CmpOp(cmpOpToken), nil
	default:
		return CmpOp{}, errors.New("CmpOp: invalid type")
	}
}

func NewReturnStatement(returnToken lexer.Token, exp *Exp, semicolonToken lexer.Token) (ReturnStatement, error) {
	if returnToken.Type != lexer.RETURN {
		return ReturnStatement{}, errors.New("ReturnStatement: invalid return token")
	}
	if semicolonToken.Type != lexer.SEMICOLON {
		return ReturnStatement{}, errors.New("ReturnStatement: invalid semicolon token")
	}

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

func NewBreakStatement(breakToken, semicolonToken lexer.Token) (BreakStatement, error) {
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

func NewContinueStatement(continueToken, semicolonToken lexer.Token) (ContinueStatement, error) {
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

func NewExp(lTerm Term, plusOrMinus *lexer.Token, rTerm *Term) (Exp, error) {
	if (plusOrMinus == nil && rTerm != nil) || (plusOrMinus != nil && rTerm == nil) {
		return Exp{}, errors.New("Exp: expected plusOrMinus and rTerm to be both nil or both not nil")
	}
	if plusOrMinus != nil && plusOrMinus.Type != lexer.PLUS && plusOrMinus.Type != lexer.MINUS {
		return Exp{}, errors.New("Exp: invalid plusOrMinus token")
	}
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

func NewTerm(lFactor Factor, mulOrDiv *lexer.Token, rFactor *Factor) (Term, error) {
	if (mulOrDiv == nil && rFactor != nil) || (mulOrDiv != nil && rFactor == nil) {
		return Term{}, errors.New("Term: expected mulOrDiv and rFactor to be both nil or both not nil")
	}
	if mulOrDiv != nil && mulOrDiv.Type != lexer.TIMES && mulOrDiv.Type != lexer.DIVIDE {
		return Term{}, errors.New("Term: invalid mulOrDiv token")
	}
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

func NewFactor(factor ...any) (Factor, error) {
	if len(factor) == 0 {
		return Factor{}, errors.New("Factor: expected at least one factor")
	} else if len(factor) == 1 {
		switch factor[0].(type) {
		case lexer.Token, ID:
			if factor[0].(lexer.Token).Type != lexer.IDENTIFIER && factor[0].(lexer.Token).Type != lexer.INTEGER_LITERAL && factor[0].(lexer.Token).Type != lexer.DECIMAL_LITERAL {
				return Factor{}, errors.New("Factor: invalid token")
			}
			return Factor{
				Factor: factor[0].(lexer.Token),
			}, nil
		}
	} else if len(factor) == 3 {
		switch factor[0].(type) {
		case lexer.Token:
			if factor[0].(lexer.Token).Type != lexer.LPAREN {
				return Factor{}, errors.New("Factor: invalid lParen token")
			}
			switch factor[2].(type) {
			case lexer.Token:
				if factor[2].(lexer.Token).Type != lexer.RPAREN {
					return Factor{}, errors.New("Factor: invalid rParen token")
				}
				switch factor[1].(type) {
				case *Exp:
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

func NewConditionalExp(lExp RelationExp, Or *lexer.Token, rExp *RelationExp) (ConditionalExp, error) {
	if (Or == nil && rExp != nil) || (Or != nil && rExp == nil) {
		return ConditionalExp{}, errors.New("NewConditionalExp: expected Or and rExp to be both nil or both not nil")
	}
	if Or != nil && Or.Type != lexer.OR {
		return ConditionalExp{}, errors.New("NewConditionalExp: invalid Or token")
	}
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

func NewRelationExp(lExp CompExp, And *lexer.Token, rExp *CompExp) (RelationExp, error) {
	if (And == nil && rExp != nil) || (And != nil && rExp == nil) {
		return RelationExp{}, errors.New("NewRelationExp: expected And and rExp to be both nil or both not nil")
	}
	if And != nil && And.Type != lexer.OR {
		return RelationExp{}, errors.New("NewRelationExp: invalid And token")
	}
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

func NewCompExp(lExp Exp, cmpOp CmpOp, rExp Exp) (CompExp, error) {
	if cmpOp.Type != lexer.LESS && cmpOp.Type != lexer.LESSEQUAL && cmpOp.Type != lexer.GREATER && cmpOp.Type != lexer.GREATEREQUAL && cmpOp.Type != lexer.EQUAL && cmpOp.Type != lexer.DIAMOND {
		return CompExp{}, errors.New("NewCompExp: invalid cmpOp token")
	}
	return CompExp{
		LExp:  lExp,
		CmpOp: cmpOp,
		RExp:  rExp,
	}, nil
}
