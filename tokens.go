package main

import (
	"errors"
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
	TokenLParen
	// TokenRParen is the token type for a right parenthesis.
	TokenRParen
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
	TOKENILLEGAL

	TOKENEOF
)

// just initialize these at the beginning, even if the tokens may not be used
// will allow fewer allocations during json stream
// maybe lazy init? but thats limited benefit for more complexity
var (
	tokenComma       = &tokenWOVal{t: TokenComma}
	tokenLParen      = &tokenWOVal{t: TokenLParen}
	tokenRParen      = &tokenWOVal{t: TokenRParen}
	tokenColon       = &tokenWOVal{t: TokenColon}
	tokenLBrace      = &tokenWOVal{t: TokenLBrace}
	tokenRBrace      = &tokenWOVal{t: TokenRBrace}
	tokenTrueBool    = &tokenWOVal{t: TokenTrueBool}
	tokenFalseBool   = &tokenWOVal{t: TokenFalseBool}
	tokenNull        = &tokenWOVal{t: TokenNull}
	tokenDoubleQuote = &tokenWOVal{t: TokenDoubleQuote}
	tokenSingleQuote = &tokenWOVal{t: TokenSingleQuote}
	tokenILLEGAL     = &tokenWOVal{t: TOKENILLEGAL}
	tokenEOF         = &tokenWOVal{t: TOKENEOF}
)

type parsedToken interface {
	getType() tokenType
	getVal() []byte
	String() string
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
	case TokenLParen:
		return tokenLParen
	case TokenRParen:
		return tokenRParen
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
	case TOKENILLEGAL:
		return tokenILLEGAL
	case TOKENEOF:
		return tokenEOF
	default:
		panic(errors.New("for strings numbers, use NewParsedTokenFromBytes"))
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
	case TokenLParen:
		return "("
	case TokenRParen:
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
	case TOKENILLEGAL:
		return "ILLEGAL"
	case TOKENEOF:
		return "EOF"
	case TokenNull:
		return "null"
	default:
		return fmt.Sprintf("UnknownToken(%d)", t)
	}
}
