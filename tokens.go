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
	tokenComma       = &tokenWOVal{t: TokenComma}
	tokenLBracket    = &tokenWOVal{t: TokenLBracket}
	tokenRBracket    = &tokenWOVal{t: TokenRBracket}
	tokenColon       = &tokenWOVal{t: TokenColon}
	tokenLBrace      = &tokenWOVal{t: TokenLBrace}
	tokenRBrace      = &tokenWOVal{t: TokenRBrace}
	tokenTrueBool    = &tokenWOVal{t: TokenTrueBool}
	tokenFalseBool   = &tokenWOVal{t: TokenFalseBool}
	tokenNull        = &tokenWOVal{t: TokenNull}
	tokenDoubleQuote = &tokenWOVal{t: TokenDoubleQuote}
	tokenSingleQuote = &tokenWOVal{t: TokenSingleQuote}
	tokenILLEGAL     = &tokenWOVal{t: TokenIllegal}
	tokenEOF         = &tokenWOVal{t: TokenEOF}
)

type parsedToken interface {
	getType() tokenType
	getVal() []byte
	String() string
	// isOpener() bool
	// isCloser() bool
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

func (twol *tokenWOVal) String() string {
	return twol.t.String()
}

func (twl *tokenWVal) getType() tokenType {
	return twl.t
}

func (twl *tokenWVal) getVal() []byte {
	return twl.v
}

func (twl *tokenWVal) String() string {
	return fmt.Sprintf("%s (%s)", string(twl.v), twl.t)
}

// NewParsedToken returns a parsedToken from
func (t tokenType) NewParsedToken() parsedToken {
	switch t {
	case TokenComma:
		return tokenComma
	case TokenLBracket:
		return tokenLBracket
	case TokenRBracket:
		return tokenRBracket
	case TokenColon:
		return tokenColon
	case TokenLBrace:
		return tokenLBrace
	case TokenRBrace:
		return tokenRBrace
	case TokenNull:
		return tokenNull
	case TokenDoubleQuote:
		return tokenDoubleQuote
	case TokenSingleQuote:
		return tokenSingleQuote
	case TokenFalseBool:
		return tokenFalseBool
	case TokenTrueBool:
		return tokenTrueBool
	case TokenIllegal:
		return tokenILLEGAL
	case TokenEOF:
		return tokenEOF
	default:
		panic(ErrLol)
	}
}

func (t tokenType) NewParsedTokenFromBytes(d []byte) parsedToken {
	return &tokenWVal{t: t, v: d}
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
		return "("
	case TokenRBracket:
		return ")"
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

// func (t tokenType) isOpener() bool {
// 	return t == TokenLBrace || t == TokenLBracket
// }
