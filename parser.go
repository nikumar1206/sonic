package main

import (
	"fmt"
	"io"
	"iter"
	"strconv"
)

// basically keep a stack
// stack item can be a fully valid json (nested json exists)
// stack item is popped everytime we encounter a closing bracket or a closing ], depending on stack opener
// stack item is added when we encounter an opening bracket or opening [

type Parser struct {
	lexer *lexer
	stack *stack
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{
		lexer: newLexer(rd),
		stack: newStack(),
	}
}

// Parse needs to be an iterable
func (p *Parser) Parse(dropinMap map[string]any) {
	next, stop := iter.Pull(p.lexer.tokens())
	fmt.Println("started parsing")
	defer stop()

	for {
		parsedToken := getToken(next)
		fmt.Println("token we got", parsedToken)
		switch tokenType := parsedToken.getType(); tokenType {
		case TOKENEOF:
			if p.stack.len() != 1 {
				panic("something is wrong")
			}
			dropinMap = p.stack.pop().value.(map[string]any)
			return
		case TokenLBracket:
			// array
			p.stack.push(stackItem{
				kind: tokenType, value: []any{},
			})
		case TokenLBrace:
			p.stack.push(stackItem{
				kind: tokenType, value: map[any]any{},
			})
		case TokenString, TokenNumber, TokenFalseBool, TokenTrueBool, TokenNull:
			value, err := p.parseValue(parsedToken)
			if err != nil {
				panic(err)
			}

			p.pushVal(stackItem{kind: tokenType, value: value})

		case TokenRBrace, TokenRBracket:
			// is a closer token

			if p.stack.len() < 2 {
				panic("bad json i think. u might have too many closers")
			}
			p.pushVal(p.stack.pop())

		}

	}
}

func (p *Parser) pushVal(s stackItem) {
	if p.stack.len() == 0 {
		p.stack.push(s)
		return
	}

	lastItem := p.stack.peak()

	switch lastItemVal := lastItem.value.(type) {
	case map[string]any:
		p.Parse(lastItemVal)
		key, valid := s.value.(string)
		if !valid {
			panic("i gues it wasnt a string oops gl!")
		}
		lastItemVal[key] = s.value
	case []any:
		lastItemVal = append(lastItemVal, s.value)
	}
}

// getToken just calls the next on manual iterator.
// but who cares about the valid value??
func getToken(next func() (parsedToken, bool)) parsedToken {
	token, valid := next()

	if !valid {
		panic("invalid? should never happen")
	}
	return token
}

func (p *Parser) parseValue(pt parsedToken) (any, error) {
	switch pt.getType() {
	case TokenString:
		return string(pt.getVal()), nil
	case TokenNumber:
		isInt, intVal, floatVal, err := parseNumber(pt.getVal())
		if err != nil {
			panic("wow what a number")
		}

		if isInt {
			return intVal, nil
		}
		return floatVal, nil

	case TokenTrueBool:
		return true, nil
	case TokenFalseBool:
		return false, nil
	case TokenNull:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", pt.getVal())
	}
}

func parseNumber(b []byte) (isInt bool, intVal int64, floatVal float64, err error) {
	if intVal, err = strconv.ParseInt(string(b), 10, 64); err == nil {
		return true, intVal, 0, nil
	}

	if floatVal, err = strconv.ParseFloat(string(b), 64); err == nil {
		return false, 0, floatVal, nil
	}

	// Not a valid number
	return false, 0, 0, err
}
