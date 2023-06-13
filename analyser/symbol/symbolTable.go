package symbol

// SymbolTable 使用泛型实现的符号表
type SymbolTable[T any] struct {
	Symbols map[string]T
}

// NewSymbolTable 创建一个新的符号表
func NewSymbolTable[T any]() *SymbolTable[T] {
	return &SymbolTable[T]{
		Symbols: make(map[string]T),
	}
}

// HasSymbol 判断符号表中是否存在某个符号
func (s *SymbolTable[T]) HasSymbol(name string) bool {
	_, ok := s.Symbols[name]
	return ok
}

// AddSymbol 向符号表中添加一个符号
func (s *SymbolTable[T]) AddSymbol(name string, symbol T) {
	s.Symbols[name] = symbol
}

// GetSymbol 获取符号表中的某个符号
func (s *SymbolTable[T]) GetSymbol(name string) (T, bool) {
	symbol, ok := s.Symbols[name]
	return symbol, ok
}

// RemoveSymbol 从符号表中移除某个符号
func (s *SymbolTable[T]) RemoveSymbol(name string) {
	if !s.HasSymbol(name) {
		return
	}

	delete(s.Symbols, name)
}

// Size 获取符号表的大小
func (s *SymbolTable[T]) Size() int {
	return len(s.Symbols)
}

// ToArray 将符号表转换为数组
func (s *SymbolTable[T]) ToArray() []T {
	array := make([]T, 0)
	for _, symbol := range s.Symbols {
		array = append(array, symbol)
	}
	return array
}
