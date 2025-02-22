package main

import (
	"bufio"
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
	if err != nil {
		return tokenEOF
	}

	for isWhiteSpace(ch) {
		ch, err = l.reader.ReadByte()
		if err != nil {
			return tokenEOF
		}
	}

	switch ch {
	case '{':
		return tokenLBrace
	case '}':
		return tokenRBrace
	case ':':
		return tokenColon
	case ',':
		return tokenComma
	case '[':
		return tokenLBracket
	case ']':
		return tokenRBracket
	case 'f', 'n', 't':
		l.reader.UnreadByte()
		return TokenIdent.NewParsedTokenFromString(string(l.readValue(keepReadingIdent, false, 5)))
	case '"':
		return TokenString.NewParsedTokenFromBytes(l.readDoubleQuoteString())
	case '\'':
		return TokenString.NewParsedTokenFromBytes(l.readSingleQuoteString())
	default:
		if isNumeric(ch) {
			l.reader.UnreadByte()
			return TokenNumber.NewParsedTokenFromBytes(l.readNumber())
		} else {
			return tokenIllegal
		}
	}
}

// tokens enables the iter Pattern for consuming tokens.
// can we read and parse concurrently with this? idk
func (l *lexer) Tokens() iter.Seq[parsedToken] {
	return func(yield func(parsedToken) bool) {
		for {
			if token := l.nextToken(); !yield(token) {
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
	return l.readValue(keepReadingDoubleQuoteString, true, 256)
}

func (l *lexer) readSingleQuoteString() []byte {
	return l.readValue(keepReadingSingleQuoteString, true, 256)
}

func keepReadingDoubleQuoteString(b byte) bool { return b != '"' }
func keepReadingIdent(b byte) bool             { return isAlpha(b) && !isWhiteSpace(b) }
func keepReadingSingleQuoteString(b byte) bool { return b != '\'' }

func (l *lexer) readNumber() []byte {
	return l.readValue(isNumeric, false, 8)
}

// readValue reads until EOF or continueFunc returns False
func (l *lexer) readValue(continueFunc func(byte) bool, hasCloser bool, bufCap int) []byte {
	buf := make([]byte, 0, 12)
	for {
		ch, err := l.reader.ReadByte()
		if err != nil {
			slog.Debug(err.Error())
			break
		}

		if !continueFunc(ch) {
			if !hasCloser {
				l.reader.UnreadByte()
			}
			break
		} else {
			buf = append(buf, ch)
		}
	}

	return buf
}

func (l *lexer) NewParsedToken() parsedToken {
	val := string(l.readValue(keepReadingIdent, false, 5))

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

func isAlpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isNumeric(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ch == 'e' || ch == '.' || ch == '-'
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
