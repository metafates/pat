package stack

type Stack[T any] struct {
	items []T
	len   int
}

func New[T any]() *Stack[T] {
	return &Stack[T]{items: []T{}}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
	s.len++
}

func (s *Stack[T]) Pop() (item T) {
	if s.len == 0 {
		return
	}

	item = s.items[s.len-1]
	s.items = s.items[:s.len-1]
	s.len--

	return
}

func (s *Stack[T]) Len() int {
	return s.len
}

func (s *Stack[T]) Peek() (t T) {
	if s.len == 0 {
		return t
	}

	return s.items[s.len-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}
