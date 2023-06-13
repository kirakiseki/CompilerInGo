package mir

import "github.com/kpango/glg"

// Stack 使用泛型实现的栈
type Stack[T any] struct {
	Elems []T
	Size  int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Elems: make([]T, 0),
		Size:  0,
	}
}

func (s *Stack[T]) Push(elem T) {
	s.Elems = append(s.Elems, elem)
	s.Size++
}

func (s *Stack[T]) Top() T {
	if s.Size == 0 {
		glg.Fatal("Out of stack")
	}
	return s.Elems[s.Size-1]
}

func (s *Stack[T]) Pop() T {
	e := s.Top()
	s.Size--
	return e
}
