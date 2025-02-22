package main

import (
	"fmt"
	"io"
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
func (p *Parser) Parse() any {
	for {
		parsedToken := p.lexer.nextToken()
		switch tokenType := parsedToken.getType(); tokenType {
		case TokenEOF:
			if p.stack.len() != 1 {
				panic("something is wrong")
			}
			val := p.stack.pop()

			return val
		case TokenLBracket:
			// array
			p.stack.push([]any{})
		case TokenLBrace:
			p.stack.push(map[string]any{})

		case TokenString, TokenNumber, TokenFalseBool, TokenTrueBool, TokenNull:
			value, err := p.parseValue(parsedToken)
			if err != nil {
				panic(err)
			}

			p.pushVal(value)
		case TokenRBrace, TokenRBracket:
			// is a closer token
			if p.stack.len() < 1 {
				panic("bad json i think. u might have too many closers")
			}
			lastItem := p.stack.pop()
			p.pushVal(lastItem)
		}

	}
}

func (p *Parser) pushVal(s any) {
	if p.stack.len() == 0 {
		p.stack.push(s)
		return
	}

	lastItemPtr := p.stack.peak()
	if lastItemPtr == nil {
		panic("Unexpected nil stack peak")
	}

	// Get the actual value
	lastItem := *lastItemPtr

	switch lastItemVal := (lastItem).(type) {
	case map[string]any:
		if p.stack.getLastUndefinedKey() == nil {
			key, valid := s.(string)
			if !valid {
				panic("key should be a string")
			}
			lastItemVal[key] = nil
			p.stack.setLastUndefinedKey(&key)
		} else {
			lastItemVal[*p.stack.getLastUndefinedKey()] = s
			p.stack.setLastUndefinedKey(nil)
		}
	case []any:
		lastItemVal = append(lastItemVal, s)
		*p.stack.peak() = lastItemVal

	default: // Handle the case where the top of the stack is NOT a map or array.
		panic("Invalid stack state. Expected map or array. received")
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
		val, err := parseNumber(pt.getVal())
		if err != nil {
			panic("wow what a number")
		}

		return val, nil

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

func parseNumber(b []byte) (v any, err error) {
	if v, err = strconv.ParseInt(string(b), 10, 64); err == nil {
		return
	}

	if v, err = strconv.ParseFloat(string(b), 64); err == nil {
		return
	}

	// Not a valid number
	return 0, ErrNaN
}
