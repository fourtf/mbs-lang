package parser

import (
	. "mbs/common"
	"regexp"
	"strconv"
	"strings"
)

func ParseWriteVar(code []byte) ([]byte, Expr, error) {
	code, name, err := ParseName(code)
	if err != nil {
		return nil, nil, err
	}

	equalsPos := strings.Index(string(code), "=")
	semicolonPos := strings.Index(string(code), ";")
	expressionString := strings.ReplaceAll(string(code[equalsPos+1:semicolonPos]), " ", "")
	expr, err := ParseExpression(expressionString)
	if err != nil {
		return nil, nil, err
	}
	code = code[semicolonPos+1:]

	return code, WriteVar{Name: string(name), Expr: expr}, nil
}

func ParseExpression(expression string) (Expr, error) {
	if exp, err := ParseString(expression); err == nil {
		return exp, nil
	}
	if exp, err := ParseBoolean(expression); err == nil {
		return exp, nil
	}
	if exp, err := ParseInteger(expression); err == nil {
		return exp, nil
	}
	if exp, err := ParseFloat(expression); err == nil {
		return exp, nil
	}
	if exp, err := ParseFunction(expression); err == nil {
		return exp, nil
	}
	if exp, err := ParseOperator(expression); err == nil {
		return exp, nil
	}
	return nil, &ParseError{Message: "Couldn't parse to any expression"}
}

func ParseString(expression string) (Expr, error) {
	if strings.HasPrefix(expression, "\"") && strings.HasSuffix(expression, "\"") {
		return String{Data: expression[1 : len(expression)-1]}, nil
	}
	return nil, &ParseError{Message: "Couldn't parse the expression to a String"}
}

func ParseBoolean(expression string) (Expr, error) {
	if expression == "true" {
		return Boolean{Data: true}, nil
	} else if expression == "false" {
		return Boolean{Data: false}, nil
	}
	return nil, &ParseError{Message: "Couldn't parse the expression to a Boolean"}
}

func ParseInteger(expression string) (Expr, error) {
	if integer, err := strconv.ParseInt(expression, 10, 64); err == nil {
		return Integer{Data: integer}, nil
	}
	return nil, &ParseError{Message: "Couldn't parse the expression to an Integer"}
}

func ParseFloat(expression string) (Expr, error) {
	if float, err := strconv.ParseFloat(expression, 64); err == nil {
		return Float{Data: float}, nil
	}
	return nil, &ParseError{Message: "Couldn't parse the expression to a Float"}
}

func ParseFunction(expression string) (Expr, error) {
	// TODO
	return nil, &ParseError{Message: "Couldn't parse the expression to a function"}
}

var (
	operators = []string{"+", "-", "*", "/", ">", "<", "==", "!=", ">=", "<="}
)

func ParseOperator(expression string) (Expr, error) {
	for _, operator := range operators {
		pos := strings.Index(expression, operator)
		if pos != -1 {
			firstExp, err := ParseExpression(expression[:pos])
			if err != nil {
				return nil, &ParseError{Message: "Couldn't parse first expression!"}
			}
			secondExp, err := ParseExpression(expression[pos+1:])
			if err != nil {
				return nil, &ParseError{Message: "Couldn't parse second expression!"}
			}
			return Operator{Symbol: operator, FirstExp: firstExp, SecondExp: secondExp}, nil
		}
	}
	return nil, &ParseError{Message: "Couldn't parse the expression to an Operator"}
}

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*`)
)

type ParseError struct {
	Message string
}

func (m *ParseError) Error() string {
	return m.Message
}

// ParseName takes an input and returns one of:
// - (the code without the name, the name, nil)
// - (nil, nil, the error)
func ParseName(code []byte) ([]byte, []byte, error) {
	codeWithoutWhitespace := []byte(strings.ReplaceAll(string(code), " ", ""))
	name := nameRegex.Find(codeWithoutWhitespace)

	if name != nil {
		return codeWithoutWhitespace[len(name):], name, nil
	}

	return nil, nil, &ParseError{Message: "Couldn't parse the name"}
}

func ParseCode(code string) (*Block, error) {
	code = strings.ReplaceAll(code, " ", "") // to prevent filtering whitespace over and over again
	// Either:
	// - WriteVar
	// - Function
	// - If
	// - For

	return nil, nil
}
