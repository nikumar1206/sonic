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
	tokenComma     = &tokenWOVal{t: TokenComma}
	tokenLBracket  = &tokenWOVal{t: TokenLBracket}
	tokenRBracket  = &tokenWOVal{t: TokenRBracket}
	tokenColon     = &tokenWOVal{t: TokenColon}
	tokenLBrace    = &tokenWOVal{t: TokenLBrace}
	tokenRBrace    = &tokenWOVal{t: TokenRBrace}
	tokenTrueBool  = &tokenWOVal{t: TokenTrueBool}
	tokenFalseBool = &tokenWOVal{t: TokenFalseBool}
	tokenNull      = &tokenWOVal{t: TokenNull}
	// tokenDoubleQuote = &tokenWOVal{t: TokenDoubleQuote}
	// tokenSingleQuote = &tokenWOVal{t: TokenSingleQuote}
	tokenIllegal = &tokenWOVal{t: TokenIllegal}
	tokenEOF     = &tokenWOVal{t: TokenEOF}
)

type parsedToken interface {
	getType() tokenType
	getVal() []byte
}

// tokenWVal should be used only for numbers
type tokenWVal struct {
	t tokenType
	v []byte
}

// tokenWOVal should be used for all other tokens
type tokenWOVal struct {
	t tokenType
}

func (twol *tokenWOVal) getVal() []byte {
	return nil
}

func (twol *tokenWOVal) getType() tokenType {
	return twol.t
}

func (twl *tokenWVal) getType() tokenType {
	return twl.t
}

func (twl *tokenWVal) getVal() []byte {
	return twl.v
}

func (t tokenType) NewParsedTokenFromBytes(d []byte) parsedToken {
	return &tokenWVal{t: t, v: d}
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
