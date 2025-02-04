package main

import (
	"io"
)

// basically keep a stack
// stack item can be a fully valid json (nested json exists)
// stack item is popped everytime we encounter a closing bracket or a closing ], depending on stack opener
// stack item is added when we encounter an opening bracket or opening [

type Parser struct {
	lexer *lexer
	stack *stack
}

func New(rd io.Reader) *Parser {
	return &Parser{
		lexer: newLexer(rd),
		stack: newStack(),
	}
}

// Parse needs to be an iterable
func (p *Parser) Parse() {}

func (p *Parser) handleMap() {
}

func (p *Parser) handleArr() {
}

func (p *Parser) handleVal() {
}
