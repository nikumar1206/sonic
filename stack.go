package main

// stack methods are push, pop, peak

type stackItem struct {
	kind  tokenType
	value any
}

type stack struct {
	data []stackItem
}

func newStack() *stack {
	return &stack{
		data: []stackItem{},
	}
}

func (s *stack) push(i stackItem) {
	s.data = append(s.data, i)
}

func (s *stack) pop() stackItem {
	lastItem := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]

	return lastItem
}

func (s *stack) peak() *stackItem {
	return &s.data[len(s.data)-1]
}

func (s *stack) len() int {
	return len(s.data)
}
