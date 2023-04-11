package utils

// Position 位置
type Position struct {
	Row     uint
	Col     uint
	FilePos uint
	Ch      rune
}

// PositionPair 位置对（开始+结束）
type PositionPair struct {
	Begin Position
	End   Position
}
