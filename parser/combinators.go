package parser

import (
	. "mbs/common"
	"strings"
)

type Parser func(string) (string, error)

func name(out *string) Parser {
	return func(code string) (string, error) {
		code, name, err := ParseName(code)
		if err == nil {
			*out = name
		}
		return code, err
	}
}

func expr(out *Expr) Parser {
	return func(code string) (string, error) {
		code, expr, err := ParseExpression(code)
		if err == nil {
			*out = expr
		}
		return code, err
	}
}

func pfunc(out *Expr, fn func(string) (string, Expr, error)) Parser {
	return func(code string) (string, error) {
		code, expr, err := fn(code)
		if err == nil {
			*out = expr
		}
		return code, err
	}
}

func block(out *Block) Parser {
	return func(code string) (string, error) {
		code, expr, err := ParseBlock(code)
		if err == nil {
			*out = expr
		}
		return code, err
	}
}

func sequence(parsers ...Parser) Parser {
	return func(code string) (string, error) {
		for _, p := range parsers {
			var err error
			code, err = p(code)

			if err != nil {
				return code, err
			}
		}
		return code, nil
	}
}

func alternative(parsers ...Parser) Parser {
	return func(code string) (string, error) {
		for _, p := range parsers {
			tmp, err := p(code)

			if err == nil {
				return tmp, nil
			}
		}
		return code, &ParseError{Message: "Couldn't match any alternative"}
	}
}

func token(t string) Parser {
	return func(code string) (string, error) {
		code = stripWhitespaceLeft(code)
		if !strings.HasPrefix(code, t) {
			return code, &ParseError{Message: "Couldn't match token '" + t + "'"}
		}

		return code[len(t):], nil
	}
}

func opt(p Parser) Parser {
	return func(code string) (string, error) {
		tmp, err := p(code)

		if err != nil {
			return code, nil
		}

		return tmp, nil
	}
}
