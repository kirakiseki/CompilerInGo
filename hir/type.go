package hir

import "fmt"

type Literal interface {
	lit()
}

type Type int

type ResultType int

const (
	TErr = iota
	TInteger
	TFloat
	TChar
	TString
	TVoid
)

type ID string

type TypeIDPair struct {
	Type Type
	ID   ID
}

type Integer struct {
	Val int64
}

type Float struct {
	Val float64
}

type Char struct {
	Val rune
}

type String struct {
	Val string
}

func NewTypeIDPair(t Type, id string) *TypeIDPair {
	return &TypeIDPair{Type: t, ID: ID(id)}
}

func NewInteger(val int64) *Integer {
	return &Integer{Val: val}
}

func (i ID) factor() {}

func (i Integer) GetVal() int64 {
	return i.Val
}

func (i Integer) lit() {}

func (i Integer) factor() {}

func NewFloat(val float64) *Float {
	return &Float{Val: val}
}

func (f Float) GetVal() float64 {
	return f.Val
}

func (f Float) lit() {}

func (f Float) factor() {}

func NewChar(val rune) *Char {
	return &Char{Val: val}
}

func (c Char) GetVal() rune {
	return c.Val
}

func (c Char) lit() {}

func NewString(val string) *String {
	return &String{Val: val}
}

func (s String) GetVal() string {
	return s.Val
}

func (s String) lit() {}

func VarToStr(id int) string {
	return fmt.Sprintf("_T%d", id)
}

func VarToStrWithSuffix(id int, suffix string) string {
	return fmt.Sprintf("_T%d_%s", id, suffix)
}

func StrToVar(s string) int {
	var id int
	_, err := fmt.Sscanf(s, "_T%d", &id)
	if err != nil {
		return 0
	}
	return id
}
