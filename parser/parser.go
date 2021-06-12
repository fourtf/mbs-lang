package parser

import (
	. "mbs/common"
	"regexp"
)

func ParseWriteVar(code []byte) ([]byte, Expr, error) {
	// In sequence:
	// - Name
	// - "="
	// - Expr

	code, name, err := ParseName(code)
	if err != nil {
		return nil, nil, err
	}

	// TODO
	name = name

	return nil, nil, nil
}

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*`)
)

type ParseError struct{}

func (m *ParseError) Error() string {
	return "parse error"
}

// ParseName takes an input and returns one of:
// - (the code without the name, the name, nil)
// - (nil, nil, the error)
func ParseName(code []byte) ([]byte, []byte, error) {
	name := nameRegex.Find(code)

	if name != nil {
		return code[len(name):], name, nil
	}

	return nil, nil, &ParseError{}
}

func ParseCode(code string) (*Block, error) {
	// Either:
	// - WriteVar
	// - Function
	// - If
	// - For

	return nil, nil
}
