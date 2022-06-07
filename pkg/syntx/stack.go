package syntx

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{make([]T, 0)}
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var result T

	l := len(s.data)

	if l == 0 {
		return result, false
	}

	result = s.data[l-1]
	s.data = s.data[:l-1]
	return result, true
}

func (s *Stack[T]) Top() (T, bool) {
	var result T

	l := len(s.data)

	if l == 0 {
		return result, false
	}

	return s.data[l-1], true
}

func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}
