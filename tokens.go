package main

import (
	"fmt"
)

type tokenType int8

const (
	// TokenInt is the token type for an integer.
	TokenNumber tokenType = iota
	// TokenString is the token type for a string.
	TokenString
	// TokenComma is the token type for a comma.
	TokenComma
	TokenLBracket
	// TokenRBracket is the token type for a right bracket.
	TokenRBracket
	// TokenLBrace is the token type for a left brace.
	TokenIdent // reserved for true, false, null
	// TokenColon is the token type for a colon.
	TokenColon
	// TokenRBrace is the token type for a right brace.
	TokenLBrace
	TokenRBrace
	TokenTrueBool
	TokenFalseBool
	TokenNull
	TokenDoubleQuote
	TokenSingleQuote
	// TokenIllegal is anything that doesn't belong in JSON.
	TokenIllegal
	TokenEOF
)

// just initialize these at the beginning, even if the tokens may not be used
// will allow fewer allocations during json stream
// maybe lazy init? but thats limited benefit for more complexity
var (
	tokenComma     = Token{_type: TokenComma}
	tokenLBracket  = Token{_type: TokenLBracket}
	tokenRBracket  = Token{_type: TokenRBracket}
	tokenColon     = Token{_type: TokenColon}
	tokenLBrace    = Token{_type: TokenLBrace}
	tokenRBrace    = Token{_type: TokenRBrace}
	tokenTrueBool  = Token{_type: TokenTrueBool}
	tokenFalseBool = Token{_type: TokenFalseBool}
	tokenNull      = Token{_type: TokenNull}
	tokenIllegal   = Token{_type: TokenIllegal}
	tokenEOF       = Token{_type: TokenEOF}
)

// tokenWVal should be used only for numbers
type Token struct {
	_type tokenType
	value string
}

func (t tokenType) NewTokenFromBytes(d []byte) Token {
	return Token{_type: t, value: bytesToString(d)}
}

func (t tokenType) NewTokenFromString(val string) Token {
	if len(val) > 5 {
		return tokenIllegal
	}

	switch val {

	case "null":
		return tokenNull
	case "false":
		return tokenFalseBool

	case "true":
		return tokenTrueBool

	default:
		return tokenIllegal
	}
}

func (t tokenType) String() string {
	switch t {
	case TokenNumber:
		return "Number"
	case TokenString:
		return "String"
	case TokenComma:
		return "Comma"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenColon:
		return ":"
	case TokenDoubleQuote:
		return "\""
	case TokenSingleQuote:
		return "'"
	case TokenIllegal:
		return "ILLEGAL"
	case TokenEOF:
		return "EOF"
	case TokenNull:
		return "null"
	default:
		return fmt.Sprintf("UnknownToken(%d)", t)
	}
}
