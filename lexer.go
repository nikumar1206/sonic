package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"unicode"
)

type lexer struct {
	reader *bytes.Reader
}

func newLexer(input []byte) *lexer {
	return &lexer{
		reader: bytes.NewReader(input),
	}
}

func (l *lexer) nextToken() parsedToken {
	ch, err := l.reader.ReadByte()
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
	case '(':
		token = TokenLParen
	case ')':
		token = TokenRParen
	case 'f', 'n', 't':
		fmt.Println("found ident")
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
		l.reader.Seek(-1, io.SeekCurrent)
		return l.getIdentTokenType().NewParsedToken()
	case TokenNumber:
		l.reader.Seek(-1, io.SeekCurrent)
		return TokenNumber.NewParsedTokenFromBytes(l.readNumber())
	default:
		return token.NewParsedToken()
	}
}

func (l *lexer) readDoubleQuoteString() []byte {
	return l.readValue(keepReadingDoubleQuoteString)
}

func (l *lexer) readSingleQuoteString() []byte {
	return l.readValue(keepReadingSingleQuoteString)
}

func keepReadingDoubleQuoteString(b byte) bool { return b != '"' }
func keepReadingIdent(b byte) bool             { return isAlpha(b) || !isWhiteSpace(b) }
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
	fmt.Println("read ident value as", string(val))

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
