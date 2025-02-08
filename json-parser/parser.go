package jsonparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	lexer  *lexer
	tokens []token
	pos    int
}

func newParser(input []byte) *parser {
	p := parser{lexer: newLexer(input)}

	for {
		token, ok := p.lexer.nextToken()
		if !ok {
			break
		}

		p.tokens = append(p.tokens, token)
	}

	return &p
}

func Parse(input []byte) (interface{}, error) {
	p := newParser(input)

	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("unnexpected end of json input")
	}

	fmt.Println(p.tokens)

	return p.parseValue()
}

func (p *parser) parseValue() (interface{}, error) {
	token := p.tokens[p.pos]
	p.pos++

	switch token.tokenType {
	case TokenString:
		return token.value, nil
	case TokenNumber:
		if strings.Contains(token.value, ".") {
			return strconv.ParseFloat(token.value, 64)
		}
		return strconv.Atoi(token.value)
	case TokenBool:
		return token.value == "true", nil
	case TokenNull:
		return nil, nil
	case TokenObjectStart:
		return p.parseObject()
	case TokenArrayStart:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("unnexpected token '%s' at position %v", token.value, p.pos-1)
	}

}

func (p *parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	for p.pos < len(p.tokens) {

		if p.tokens[p.pos].tokenType == TokenObjectEnd {
			p.pos++
			break
		}

		if p.tokens[p.pos].tokenType != TokenString {
			return nil, errors.New("expected string key")
		}

		key := p.tokens[p.pos].value
		p.pos++

		if p.tokens[p.pos].tokenType != TokenColon {
			return nil, errors.New("expected : after key")
		}

		p.pos++

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		if p.tokens[p.pos].tokenType == TokenObjectEnd {
			p.pos++
			break
		}

		if p.tokens[p.pos].tokenType != TokenComma {
			return nil, errors.New("expected ',' or '}'")
		}

		p.pos++

		if p.tokens[p.pos].tokenType == TokenObjectEnd {
			return nil, fmt.Errorf("unnexpected token ','")
		}
	}

	return obj, nil
}

func (p *parser) parseArray() ([]interface{}, error) {
	arr := make([]interface{}, 0)

	for p.pos < len(p.tokens) {
		if p.tokens[p.pos].tokenType == TokenArrayEnd {
			p.pos++
			break
		}

		item, err := p.parseValue()

		if err != nil {
			return nil, err
		}

		arr = append(arr, item)

		if p.tokens[p.pos].tokenType == TokenArrayEnd {
			p.pos++
			break
		}

		if p.tokens[p.pos].tokenType != TokenComma {
			return nil, fmt.Errorf("expected ',' or ]")
		}

		p.pos++

	}

	return arr, nil
}
