package main

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"strconv"
	"unsafe"
)

// basically keep a stack
// stack item can be a fully valid json (nested json exists)
// stack item is popped everytime we encounter a closing bracket or a closing ], depending on stack opener
// stack item is added when we encounter an opening bracket or opening [

type Parser interface {
	ParseToken(Token) *any
	Parse() any
}
type StackParser struct {
	lexer *lexer
	stack *stack
}

type RecursiveParser struct {
	lexer *lexer
}

func NewParser(rd io.Reader, _type string) Parser {
	switch _type {
	case "stack":
		return &StackParser{
			lexer: newLexer(rd),
			stack: newStack(),
		}
	case "recursive":
		return &RecursiveParser{
			lexer: newLexer(rd),
		}
	default:
		panic("either recursive or stack parsing")

	}
}

func (p *StackParser) ParseToken(t Token) *any {
	switch t._type {
	case TokenEOF:
		if p.stack.len() != 1 {
			panic("something is wrong")
		}
		val := p.stack.pop()
		return &val

	case TokenLBracket:
		// array
		p.stack.push([]any{})
	case TokenLBrace:
		p.stack.push(map[string]any{})

	case TokenString, TokenNumber, TokenFalseBool, TokenTrueBool, TokenNull:
		value, err := parseValue(t)
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
	return nil
}

func (p *StackParser) pushVal(s any) {
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

func parseValue(pt Token) (any, error) {
	switch pt._type {
	case TokenString:
		v := pt.value
		return v, nil
	case TokenNumber:
		v := pt.value
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Println("what was pt", string(pt.value))
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
		return nil, fmt.Errorf("unexpected token: %s", pt.value)
	}
}

func (p *StackParser) Parse() any {
	for {
		t := p.lexer.nextToken()
		val := p.ParseToken(t)

		if val != nil {
			return *val
		}
	}
}

// logic completely incorrect here lol
func (p *RecursiveParser) Parse() any {
	ending_val := make(map[string]any)
	for {
		t := p.lexer.nextToken()
		if t == tokenEOF {
			break
		}

		val := p.ParseToken(t)

		switch v := (*val).(type) {
		case map[string]any:
			maps.Copy(ending_val, v)
		case string:
			ending_val["string_value"] = v
		case float64:
			ending_val["number_value"] = v
		case bool:
			ending_val["bool_value"] = v
		case []any:
			ending_val["array_value"] = v
		default:
			fmt.Println("hit default")
		}
	}
	var result any = ending_val
	return &result
}

func (p *RecursiveParser) ParseToken(t Token) *any {
	var val any = nil
	var err error
	switch t._type {
	case TokenLBracket:
		val, err = p.parseArray()
		if err != nil {
			panic(val)
		}
	case TokenLBrace:
		val, err = p.parseObject()
		if err != nil {
			panic(err)
		}
	case TokenString, TokenNumber, TokenFalseBool, TokenTrueBool, TokenNull:
		val, err = parseValue(t)
		if err != nil {
			panic(err)
		}

	}
	return &val
}

func (p *RecursiveParser) parseObject() (map[string]any, error) {
	obj := make(map[string]any)

	t := p.lexer.nextToken()

	for {
		switch t._type {
		case TokenRBrace:
			p.lexer.nextToken()
			return obj, nil
		case TokenString:
			p.lexer.nextToken() // likely comma?

			value, err := parseValue(p.lexer.nextToken())
			if err != nil {
				return nil, err
			}
			obj[t.value] = value

			switch t._type {
			case TokenComma:
				p.lexer.nextToken()
			case TokenRBrace:
				continue
			default:
				return nil, errors.New("expected comma or closing brace")
			}
		default:
			return nil, errors.New("expected string key or closing brace")
		}
	}
}

func (p *RecursiveParser) parseArray() ([]any, error) {
	var arr []any
	t := p.lexer.nextToken()

	for {
		switch t._type {
		case TokenRBracket:
			p.lexer.nextToken()
			return arr, nil
		default:
			value := p.ParseToken(t)
			arr = append(arr, value)

			switch t._type {
			case TokenComma:
				p.lexer.nextToken()
			case TokenRBracket:
				continue
			default:
				return nil, errors.New("expected comma or closing bracket")
			}
		}
	}
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
