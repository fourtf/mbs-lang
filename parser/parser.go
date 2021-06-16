package parser

import (
	. "mbs/common"
	"regexp"
	"strconv"
	"strings"
)

func ParseReadVar(code string) (string, Expr, error) {
	code, name, err := ParseName(code)
	if err != nil {
		return "", nil, err
	}

	return code, ReadVar{Name: name}, nil
}

func ParseWriteVar(code string) (string, Expr, error) {
	code, name, err := ParseName(code)
	if err != nil {
		return "", nil, err
	}
	equalsPos := strings.Index(code, "=")
	code, expr, err := ParseExpression(stripWhitespaceLeft(code[equalsPos+1:]))
	if err != nil {
		return "", nil, err
	}

	return code, WriteVar{Name: name, Expr: expr}, nil
}

func ParseExpression(code string) (string, Expr, error) {
	if remainingCode, exp, err := ParseOperator(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseExpressionWithoutOperator(code); err == nil {
		return remainingCode, exp, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse to any expression"}
}

func ParseExpressionWithoutOperator(code string) (string, Expr, error) {
	if remainingCode, exp, err := ParseString(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseBoolean(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseFloat(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseInteger(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseFunction(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseReadVar(code); err == nil {
		return remainingCode, exp, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse to any expression"}
}

func ParseString(code string) (string, Expr, error) {
	firstQuotes := strings.Index(code, "\"")
	secondQuotes := strings.Index(code[firstQuotes+1:], "\"")
	if firstQuotes != -1 && secondQuotes != -1 {
		return code[secondQuotes+2:], String{Data: code[firstQuotes+1 : secondQuotes+1]}, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse the expression to a String"}
}

func ParseBoolean(code string) (string, Expr, error) {
	if strings.HasPrefix(code, "true") {
		return code[4:], Boolean{Data: true}, nil
	} else if strings.HasPrefix(code, "false") {
		return code[5:], Boolean{Data: false}, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse the expression to a Boolean"}
}

var (
	intRegex = regexp.MustCompile(`^-?\d+`)
)

func ParseInteger(code string) (string, Expr, error) {
	match := intRegex.FindString(code)
	if integer, err := strconv.ParseInt(match, 10, 64); err == nil {
		return code[len(match):], Integer{Data: integer}, nil
	}

	return code, nil, &ParseError{Message: "Couldn't parse the expression to an Integer"}
}

var (
	floatRegex = regexp.MustCompile(`^-?\d+\.\d+`)
)

func ParseFloat(code string) (string, Expr, error) {
	match := floatRegex.FindString(code)
	if float, err := strconv.ParseFloat(match, 64); err == nil {
		return code[len(match):], Float{Data: float}, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse the expression to a Float"}
}

func ParseFunction(code string) (string, Expr, error) {
	// name
	code = stripWhitespaceLeft(code)
	code, name, err := ParseName(code)
	if err != nil {
		return "", nil, err
	}

	// (
	code = stripWhitespaceLeft(code)
	code, ok := stripPrefix(code, "(")
	if !ok {
		return "", nil, NewParseErrorExpected("(")
	}

	// expr
	// TODO: we only allow a single argument; not a list divided with ","
	code = stripWhitespaceLeft(code)
	code, expr, err := ParseExpression(code)
	if err != nil {
		return "", nil, err
	}

	// )
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, ")")
	if !ok {
		return "", nil, NewParseErrorExpected(")")
	}

	return code, FunctionCall{Name: name, Argument: expr}, nil
}

// ParseParentheses:  ( print ( "( )()()((" + "asd" ) )
// ParseFunction:       print ( "( )()()((" + "asd" ) ) -> )
// ParseOperator:               "( )()()((" + "asd" ) ) -> )

// "( )()()((" + "asd" ) )

// ParseParentheses -> ParseFunction -> ParseString

func ParseParentheses(code string) (string, Expr, error) {
	// (
	code, ok := stripPrefix(stripWhitespaceLeft(code[:0]), "(")
	if !ok {
		return code, nil, NewParseErrorExpected("(")
	}

	// expr
	code, expr, err := ParseExpression(code)
	if err != nil {
		return code, nil, err
	}

	// )
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code[:0], ")")
	if !ok {
		return code, nil, NewParseErrorExpected(")")
	}

	return code, expr, nil
}

var (
	operators = []string{"+", "-", "*", "/", ">", "<", "==", "!=", ">=", "<="}
)

func ParseOperator(code string) (string, Expr, error) {
	code, firstExp, err := ParseExpressionWithoutOperator(code)
	if err != nil {
		return code, nil, &ParseError{Message: "Couldn't parse first expression!"}
	}
	code = stripWhitespaceLeft(code)
	operator := ""
	for _, op := range operators {
		if strings.HasPrefix(code, op) {
			operator = op
			break
		}
	}
	if operator != "" {
		code, secondExp, err := ParseExpressionWithoutOperator(stripWhitespaceLeft(code[1:]))
		if err != nil {
			return code, nil, &ParseError{Message: "Couldn't parse second expression!"}
		}
		return code, Operator{Symbol: operator, FirstExp: firstExp, SecondExp: secondExp}, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse the expression to an Operator"}
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

func NewParseErrorExpected(expected string) *ParseError {
	return &ParseError{Message: "Expected '" + expected + "'"}
}

var whitespaceRegex = regexp.MustCompile(`\s+`)

// ParseName takes an input and returns one of:
// - (the code without the name, the name, nil)
// - (nil, nil, the error)
func ParseName(code string) (string, string, error) {

	codeWithoutWhitespace := stripWhitespaceLeft(code)
	name := nameRegex.FindString(codeWithoutWhitespace)

	if name == "" {
		return "", "", &ParseError{Message: "Couldn't parse the name"}
	}

	return codeWithoutWhitespace[len(name):], name, nil
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

var stripWhitespaceRegex = regexp.MustCompile(`^\s+`)

// stripWhitespaceLeft strips all whitespace on the left of the string and returns a string without it.
func stripWhitespaceLeft(s string) string {
	loc := stripWhitespaceRegex.FindStringIndex(s)

	if loc == nil {
		return s
	}

	return s[loc[1]:]
}

// stripPrefix tries to remove a prefix from a string. Returns the stripped string and a bool indicating if it was
// successful. Returns the original string if the prefix wasn't found.
func stripPrefix(s, prefix string) (string, bool) {
	// if strings.HasPrefix()

	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}

	return s, false
}
