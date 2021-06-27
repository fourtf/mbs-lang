package parser

import (
	. "mbs/common"
	"strings"
)

/*
	This file contains a number of parsing functions which can be conveniently combined. Less manual error handling is
	needed if the parser combinator are used because the small building blocks take care of it. Returning values is
	handled using "out" parameters. This means that the results (e.g. the value read in the name function) are written
	to the pointer parameter.
*/

// Parser is a function type which can be combined using other parser functions.
type Parser func(string) (string, error)

// name reads a Name (alphanumeric sequence) and writes the result into the adress `out`.
func name(out *string) Parser {
	return func(code string) (string, error) {
		code, name, err := ParseName(code)
		if err == nil {
			*out = name
		}
		return code, err
	}
}

// name reads an Expression using the ParserExpression function and writes the result into the adress `out`.
func expr(out *Expr) Parser {
	return func(code string) (string, error) {
		code, expr, err := ParseExpression(code)
		if err == nil {
			*out = expr
		}
		return code, err
	}
}

// pfunc parsing function which also returns an Expression into a one compatible with the `Parser` interface.
func pfunc(out *Expr, fn func(string) (string, Expr, error)) Parser {
	return func(code string) (string, error) {
		code, expr, err := fn(code)
		if err == nil {
			*out = expr
		}
		return code, err
	}
}

// name reads a Block using the ParserBlock function and writes the result into the adress `out`.
func block(out *Block) Parser {
	return func(code string) (string, error) {
		code, expr, err := ParseBlock(code)
		if err == nil {
			*out = expr
		}
		return code, err
	}
}

// sequence runs the parser in sequnce using the previous parsers code. If any parsers fail the combined parser also fails.
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

// alternative picks runs every parser but stops when one of them worked. If no parser works then an error is returned.
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

// token reads a specific sequence of characters but doesn't return the read value.
func token(t string) Parser {
	return func(code string) (string, error) {
		code = stripWhitespaceLeft(code)
		if !strings.HasPrefix(code, t) {
			return code, &ParseError{Message: "Couldn't match token '" + t + "'"}
		}

		return code[len(t):], nil
	}
}

// opt runs the provided parser. If the provided parser fails the new parser still succeeds and consumed no input.
func opt(p Parser) Parser {
	return func(code string) (string, error) {
		tmp, err := p(code)

		if err != nil {
			return code, nil
		}

		return tmp, nil
	}
}
