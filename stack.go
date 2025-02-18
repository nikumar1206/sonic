package main

import (
	"encoding/json"
	"fmt"
)

// stack methods are push, pop, peak

type stackItem struct {
	kind  tokenType
	value any
}

type stack struct {
	data              []stackItem
	lastUndefinedKeys []*string
}

func newStack() *stack {
	return &stack{
		data: []stackItem{},
	}
}

func (s *stack) push(i stackItem) {
	s.data = append(s.data, i)
	s.lastUndefinedKeys = append(s.lastUndefinedKeys, nil)
}

func (s *stack) pop() stackItem {
	lastItem := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]

	s.lastUndefinedKeys = s.lastUndefinedKeys[0 : len(s.lastUndefinedKeys)-1]

	return lastItem
}

func (s *stack) peak() *stackItem {
	return &s.data[len(s.data)-1]
}

func (s *stack) len() int {
	return len(s.data)
}

func (s *stack) getLastUndefinedKey() *string {
	return s.lastUndefinedKeys[s.len()-1]
}

func (s *stack) setLastUndefinedKey(key *string) {
	s.lastUndefinedKeys[s.len()-1] = key
}

func (s *stack) debug() {
	if len(s.data) == 0 {
		fmt.Println("stack is empty")
		return
	}
	fmt.Println("Stack Debug:")
	for i := len(s.data) - 1; i >= 0; i-- {
		val := s.data[i].value
		fmt.Printf("[%d] kind: %v, value: ", i, s.data[i].kind)
		switch v := val.(type) {
		case string:
			// Print the string as-is (without quotes)
			fmt.Println(v)
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64, bool:
			// Print the primitive type with %v
			fmt.Println(v)
		default:
			// For other types, fall back to JSON pretty printing
			valueJSON, err := json.MarshalIndent(v, " ", "  ")
			if err != nil {
				fmt.Printf("<error: %v>\n", err)
			} else {
				fmt.Println(string(valueJSON))
			}
		}
	}
}
