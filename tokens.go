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
	tokenComma     = parsedToken{t: TokenComma}
	tokenLBracket  = parsedToken{t: TokenLBracket}
	tokenRBracket  = parsedToken{t: TokenRBracket}
	tokenColon     = parsedToken{t: TokenColon}
	tokenLBrace    = parsedToken{t: TokenLBrace}
	tokenRBrace    = parsedToken{t: TokenRBrace}
	tokenTrueBool  = parsedToken{t: TokenTrueBool}
	tokenFalseBool = parsedToken{t: TokenFalseBool}
	tokenNull      = parsedToken{t: TokenNull}
	// tokenDoubleQuote = &parsedToken{t: TokenDoubleQuote}
	// tokenSingleQuote = &parsedToken{t: TokenSingleQuote}
	tokenIllegal = parsedToken{t: TokenIllegal}
	tokenEOF     = parsedToken{t: TokenEOF}
)

// tokenWVal should be used only for numbers
type parsedToken struct {
	t tokenType
	v string
}

func (twl parsedToken) getType() tokenType {
	return twl.t
}

func (twl parsedToken) getVal() string {
	return twl.v
}

func (t tokenType) NewParsedTokenFromBytes(d []byte) parsedToken {
	return parsedToken{t: t, v: bytesToString(d)}
}

func (t tokenType) NewParsedTokenFromString(val string) parsedToken {
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
