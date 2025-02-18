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
func (p *Parser) Parse() any {
	next, stop := iter.Pull(p.lexer.tokens())
	defer stop()
	for {
		parsedToken := getToken(next)
		fmt.Println("recieved token", parsedToken.getType())
		switch tokenType := parsedToken.getType(); tokenType {
		case TokenEOF:
			fmt.Println("at EOF")
			if p.stack.len() != 1 {
				panic("something is wrong")
			}
			val := p.stack.pop().value

			return val
		case TokenLBracket:
			// array
			p.stack.push(stackItem{
				kind: tokenType, value: []any{},
			})
			fmt.Println("instantiating a new array for u")
		case TokenLBrace:
			p.stack.push(stackItem{
				kind: tokenType, value: map[string]any{},
			})

			fmt.Println("instantiating new stack object")
		case TokenString, TokenNumber, TokenFalseBool, TokenTrueBool, TokenNull:
			value, err := p.parseValue(parsedToken)
			if err != nil {
				panic(err)
			}

			p.pushVal(stackItem{kind: tokenType, value: value})
		case TokenRBrace, TokenRBracket:
			// is a closer token

			if p.stack.len() < 1 {
				fmt.Println("what was in the stack at the end")
				p.stack.debug()
				panic("bad json i think. u might have too many closers")
			}
			lastItem := p.stack.pop()
			fmt.Println("we are popping", lastItem)
			p.pushVal(lastItem)
		}

	}
}

func (p *Parser) pushVal(s stackItem) {
	if p.stack.len() == 0 {
		p.stack.push(s)
		return
	}

	lastItem := p.stack.peak()
	fmt.Println("peakedvalue from stack", lastItem.value)
	switch lastItemVal := lastItem.value.(type) {
	case map[string]any:
		if p.stack.getLastUndefinedKey() == nil {
			key, valid := s.value.(string)
			if !valid {
				fmt.Println("pre panic: value was ", s.kind, s.value)
				panic("key should be a string")
			}
			lastItemVal[key] = nil
			p.stack.setLastUndefinedKey(&key)
			fmt.Println("added key ", key)
		} else {
			lastItemVal[*p.stack.getLastUndefinedKey()] = s.value
			p.stack.setLastUndefinedKey(nil)
		}
	case []any:
		lastItem.value = append(lastItemVal, s.value)
	default: // Handle the case where the top of the stack is NOT a map or array.
		panic("Invalid stack state. Expected map or array.")
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
			fmt.Println(err)
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
