package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"unicode"
)

type lexer struct {
	reader *bufio.Reader
}

func newLexer(rd io.Reader) *lexer {
	return &lexer{
		reader: bufio.NewReader(rd),
	}
}

func (l *lexer) nextToken() parsedToken {
	ch, err := l.reader.ReadByte()
	fmt.Println("yielding new token")
	if err != nil {
		return tokenEOF
	}
	token := TOKENILLEGAL

	for isWhiteSpace(ch) {
		ch, err = l.reader.ReadByte()
		if err != nil {
			return tokenEOF
		}
	}

	switch ch {
	case '{':
		token = TokenLBrace
	case '}':
		token = TokenRBrace
	case ':':
		token = TokenColon
	case ',':
		token = TokenComma
	case '[':
		token = TokenLBracket
	case ']':
		token = TokenRBracket
	case 'f', 'n', 't':
		token = TokenIdent
	case '"':
		token = TokenDoubleQuote
	case '\'':
		token = TokenSingleQuote
	default:
		if isNumeric(ch) {
			token = TokenNumber
		}
	}

	switch token {
	case TokenSingleQuote:
		return TokenString.NewParsedTokenFromBytes(l.readSingleQuoteString())
	case TokenDoubleQuote:
		return TokenString.NewParsedTokenFromBytes(l.readDoubleQuoteString())
	case TokenIdent:
		l.reader.UnreadByte()
		return l.getIdentTokenType().NewParsedToken()
	case TokenNumber:
		l.reader.UnreadByte()
		return TokenNumber.NewParsedTokenFromBytes(l.readNumber())
	default:
		return token.NewParsedToken()
	}
}

// tokens enables the iter Pattern for consuming tokens.
// can we read and parse concurrently with this? idk
func (l *lexer) tokens() iter.Seq[parsedToken] {
	fmt.Println("started yielding")
	return func(yield func(parsedToken) bool) {
		for {
			if token := l.nextToken(); token == tokenEOF || !yield(token) {
				return
			}
		}
	}
}

// some function like Send that sends it to a provided channel, and another side receives it?
func (l *lexer) sendTokens(c chan parsedToken) {
	for {
		t := l.nextToken()
		c <- t
		if t == tokenEOF {
			return
		}
	}
}

func (l *lexer) readDoubleQuoteString() []byte {
	return l.readValue(keepReadingDoubleQuoteString)
}

func (l *lexer) readSingleQuoteString() []byte {
	return l.readValue(keepReadingSingleQuoteString)
}

func keepReadingDoubleQuoteString(b byte) bool { return b != '"' }
func keepReadingIdent(b byte) bool             { return isAlpha(b) && !isWhiteSpace(b) }
func keepReadingSingleQuoteString(b byte) bool { return b != '\'' }

func (l *lexer) readNumber() []byte {
	return l.readValue(isNumeric)
}

// readValue reads until EOF or continueFunc returns False
func (l *lexer) readValue(continueFunc func(byte) bool) []byte {
	var acc bytes.Buffer

	for {
		ch, err := l.reader.ReadByte()
		if err != nil {
			slog.Debug(err.Error())
			break
		}

		if !continueFunc(ch) {
			break
		} else {
			acc.WriteByte(ch)
		}
	}
	return acc.Bytes()
}

func (l *lexer) getIdentTokenType() tokenType {
	val := string(l.readValue(keepReadingIdent))

	if val == "null" {
		return TokenNull
	}
	if val == "false" {
		return TokenFalseBool
	}
	if val == "true" {
		return TokenTrueBool
	}
	return TOKENILLEGAL
}

func isAlpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isNumeric(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isNumeric(ch)
}

func isValidWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f'
}

func isInvalidWhitespace(b byte) bool {
	if unicode.IsSpace(rune(b)) {
		return !(b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f')
	}
	return false
}

func isWhiteSpace(b byte) bool {
	return isValidWhitespace(b) || isInvalidWhitespace(b)
}
