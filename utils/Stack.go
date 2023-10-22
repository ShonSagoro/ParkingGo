package utils

type Stack struct {
	elements []interface{}
}

func NewStack() *Stack {
	return &Stack{nil}
}

func (s *Stack) Pop() interface{} {
	if s.isEmpty() {
		return nil
	}

	n := s.elements[len(s.elements)-1]
	s.elements = s.elements[len(s.elements)-1:]

	return n
}

func (s *Stack) Push(value interface{}) {
	s.elements = append(s.elements, value)
}

func (s *Stack) isEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack) Peek() interface{} {
	if s.isEmpty() {
		return nil
	}
	return s.elements[len(s.elements)-1]
}

func (s *Stack) GetElements() []interface{} {
	return s.elements
}
