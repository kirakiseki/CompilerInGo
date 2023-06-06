package symbol

type SymbolTable[T any] struct {
	Symbols map[string]T
}

func NewSymbolTable[T any]() *SymbolTable[T] {
	return &SymbolTable[T]{
		Symbols: make(map[string]T),
	}
}

func (s *SymbolTable[T]) HasSymbol(name string) bool {
	_, ok := s.Symbols[name]
	return ok
}

func (s *SymbolTable[T]) AddSymbol(name string, symbol T) {
	s.Symbols[name] = symbol
}

func (s *SymbolTable[T]) GetSymbol(name string) (T, bool) {
	symbol, ok := s.Symbols[name]
	return symbol, ok
}

func (s *SymbolTable[T]) RemoveSymbol(name string) {
	if !s.HasSymbol(name) {
		return
	}

	delete(s.Symbols, name)
}

func (s *SymbolTable[T]) Size() int {
	return len(s.Symbols)
}

func (s *SymbolTable[T]) ToArray() []T {
	array := make([]T, 0)
	for _, symbol := range s.Symbols {
		array = append(array, symbol)
	}
	return array
}