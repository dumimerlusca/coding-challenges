package jsonparser

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	pos   int
	input []byte
}

type tokenType int

const (
	TokenObjectStart tokenType = iota
	TokenObjectEnd
	TokenArrayStart
	TokenArrayEnd
	TokenColon
	TokenComma
	TokenString
	TokenNumber
	TokenBool
	TokenNull
)

type token struct {
	tokenType tokenType
	value     string
}

func newLexer(input []byte) *lexer {
	return &lexer{input: input}
}

func (l *lexer) nextToken() (token, bool) {
	l.skipWhiteSpace()

	if l.pos >= len(l.input) {
		return token{}, false
	}

	r, size := utf8.DecodeRune(l.input[l.pos:])

	switch r {
	case '{':
		l.pos += size
		return token{tokenType: TokenObjectStart, value: "{"}, true
	case '}':
		l.pos += size
		return token{tokenType: TokenObjectEnd, value: "}"}, true
	case '[':
		l.pos += size
		return token{tokenType: TokenArrayStart, value: "["}, true
	case ']':
		l.pos += size
		return token{tokenType: TokenArrayEnd, value: "]"}, true
	case ':':
		l.pos += size
		return token{tokenType: TokenColon, value: ":"}, true
	case ',':
		l.pos += size
		return token{tokenType: TokenComma, value: ","}, true
	case '"':
		l.pos += size
		return l.readString()
	default:
		if bytes.HasPrefix(l.input[l.pos:], []byte("true")) {
			l.pos += 4
			return token{tokenType: TokenBool, value: "true"}, true
		}
		if bytes.HasPrefix(l.input[l.pos:], []byte("false")) {
			l.pos += 5
			return token{tokenType: TokenBool, value: "false"}, true
		}
		if bytes.HasPrefix(l.input[l.pos:], []byte("null")) {
			l.pos += 4
			return token{tokenType: TokenNull, value: "null"}, true
		}

		if unicode.IsNumber(r) {
			return l.readNumber()
		}

		l.pos += size
		return token{tokenType: -1, value: string(r)}, true
	}
}

func (l *lexer) readNumber() (token, bool) {
	tok := token{tokenType: TokenNumber}

	for {
		if l.pos >= len(l.input) {
			return token{}, false
		}

		r, size := utf8.DecodeRune(l.input[l.pos:])

		if unicode.IsNumber(r) || r == '.' {
			tok.value += string(r)
			l.pos += size
		} else {
			break
		}
	}

	return tok, true
}

func (l *lexer) readString() (token, bool) {
	tok := token{tokenType: TokenString}
	for {
		if l.pos >= len(l.input) {
			return token{}, false
		}

		r, size := utf8.DecodeRune(l.input[l.pos:])
		l.pos += size

		if r == '"' {
			break
		} else {
			tok.value += string(r)
		}
	}

	return tok, true
}

func (l *lexer) skipWhiteSpace() {
	for {
		if l.pos >= len(l.input) {
			break
		}

		r, size := utf8.DecodeRune(l.input[l.pos:])

		if unicode.IsSpace(r) {
			l.pos += size
		} else {
			break
		}

	}

}
