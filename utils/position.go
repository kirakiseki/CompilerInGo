package utils

type Position struct {
	Row     uint
	Col     uint
	FilePos uint
	Ch      rune
}

type PositionPair struct {
	Begin Position
	End   Position
}
