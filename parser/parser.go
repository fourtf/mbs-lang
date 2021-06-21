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
	code = stripWhitespaceLeft(code)
	code, name, err := ParseName(code)
	if err != nil {
		return "", nil, err
	}

	code = stripWhitespaceLeft(code)
	code, ok := stripPrefix(code, "=")
	if !ok {
		return "", nil, NewParseErrorExpected("=")
	}

	code = stripWhitespaceLeft(code)
	code, expr, err := ParseExpression(code)
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
	if remainingCode, exp, err := ParseParentheses(code); err == nil {
		return remainingCode, exp, nil
	}
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
	if remainingCode, exp, err := ParseFunctionCall(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseReadVar(code); err == nil {
		return remainingCode, exp, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse to any expression"}
}

var (
	stringRegex       = regexp.MustCompile(`^"((\\"|[^"])*)"`)
	stringEscapeRegex = regexp.MustCompile(`\\"`)
)

func ParseString(code string) (string, Expr, error) {
	match := stringRegex.FindStringIndex(code)
	if match == nil {
		return code, nil, &ParseError{Message: "Couldn't parse a string"}
	}

	data := stringEscapeRegex.ReplaceAllStringFunc(code[match[0]+1:match[1]-1], escapeStringRepl)

	return code[match[1]:], String{Data: data}, nil
}

func escapeStringRepl(match string) string {
	switch match[1] {
	case 'r':
		return "\r"
	case 'n':
		return "\n"
	case 't':
		return "\t"
	default:
		return match[1:2]
	}
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

func ParseFunctionCall(code string) (string, Expr, error) {
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
	code, expr, err := parseOrNop(code, ParseExpression)
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
	code = stripWhitespaceLeft(code)
	code, ok := stripPrefix(code, "(")
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
	code, ok = stripPrefix(code, ")")
	if !ok {
		return code, nil, NewParseErrorExpected(")")
	}

	return code, expr, nil
}

var (
	operators = []string{"+", "-", "*", "/", ">", "<", "==", "!=", ">=", "<=", "&&", "||"}
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
	if operator == "" {
		return code, nil, &ParseError{Message: "Couldn't parse the expression to an Operator"}
	}

	code, secondExp, err := ParseExpressionWithoutOperator(stripWhitespaceLeft(code[len(operator):]))
	if err != nil {
		return code, nil, &ParseError{Message: "Couldn't parse second expression!"}
	}

	return code, Operator{Symbol: operator, FirstExp: firstExp, SecondExp: secondExp}, nil
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

func ParseIf(code string) (string, Expr, error) {
	// if
	code = stripWhitespaceLeft(code)
	code, ok := stripPrefix(code, "if")
	if !ok {
		return code, nil, NewParseErrorExpected("if")
	}

	// (
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, "(")
	if !ok {
		return code, nil, NewParseErrorExpected("(")
	}

	// expr
	code, condition, err := ParseExpression(code)
	if err != nil {
		return code, nil, err
	}

	// )
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, ")")
	if !ok {
		return code, nil, NewParseErrorExpected(")")
	}

	// {
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, "{")
	if !ok {
		return code, nil, NewParseErrorExpected("{")
	}

	// expr
	code, body, err := ParseBlock(code)
	if err != nil {
		return code, nil, err
	}

	// }
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, "}")
	if !ok {
		return code, nil, NewParseErrorExpected("}")
	}

	return code, &If{Condition: condition, Body: body}, nil
}

func ParseFor(code string) (string, Expr, error) {
	// multi(
	// 	stripPrefix("for"),
	// 	parentheses("(", multi(opt(stmt), opt(expr), opt(stmt))), ")"),
	// 	parentheses("{", block, "}"))

	// for
	code = stripWhitespaceLeft(code)
	code, ok := stripPrefix(code, "for")
	if !ok {
		return code, nil, NewParseErrorExpected("for")
	}

	// (
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, "(")
	if !ok {
		return code, nil, NewParseErrorExpected("(")
	}

	// a = 123 or Nop
	code, initExpr, _ := parseOrNop(code, ParseWriteVar)

	// ;
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, ";")
	if !ok {
		return "", nil, NewParseErrorExpected(";")
	}

	// expr
	code, conditionExpr, _ := parseOrNop(code, ParseExpression)

	// ;
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, ";")
	if !ok {
		return "", nil, NewParseErrorExpected(";")
	}

	// a = a + 1 or Nop
	code, advancementExpr, _ := parseOrNop(code, ParseWriteVar)

	// )
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, ")")
	if !ok {
		return code, nil, NewParseErrorExpected(")")
	}

	// {
	code, ok = stripPrefix(stripWhitespaceLeft(code), "{")
	if !ok {
		return code, nil, NewParseErrorExpected("{")
	}

	// block
	code, bodyExpr, err := ParseBlock(code)
	if err != nil {
		return code, nil, err
	}

	// }
	code = stripWhitespaceLeft(code)
	code, ok = stripPrefix(code, "}")
	if !ok {
		return code, nil, NewParseErrorExpected("}")
	}

	return code, &For{Init: initExpr, Condition: conditionExpr, Advancement: advancementExpr, Body: bodyExpr}, nil
}

func ParseBlock(code string) (string, *Block, error) {
	// Either:
	// - WriteVar
	// - FunctionCall
	// - If
	// - For

	stmts := make([]Expr, 0)

	// update if below if you update this!
	opts := [](func(code string) (string, Expr, error)){ParseWriteVar, ParseFunctionCall, ParseIf, ParseFor}

outer:
	for {
		for i, opt := range opts {
			remainingCode, expr, err := opt(code)
			if err != nil {
				continue
			}

			// ParseIf or ParseFor don't have a ;
			if i < 2 {
				var ok bool
				remainingCode, ok = stripPrefix(remainingCode, ";")
				if !ok {
					continue
				}
			}

			stmts = append(stmts, expr)
			code = remainingCode
			continue outer
		}

		break
	}

	return code, &Block{Statements: stmts}, nil
}

func ParseCode(code string) (*Block, error) {
	code, blk, err := ParseBlock(code)
	if err != nil {
		return nil, err
	}

	code = stripWhitespaceLeft(code)
	if code != "" {
		return nil, &ParseError{Message: "Couldn't continue parsing after: `" + code + "`"}
	}

	return blk, nil
}

// parseOrNop returns a "Nop" if the function fails to parse the code.
func parseOrNop(code string, fn func(code string) (string, Expr, error)) (string, Expr, error) {
	remainingCode, expr, err := fn(code)

	if err != nil {
		return code, &Nop{}, nil
	}

	return remainingCode, expr, nil
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
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}

	return s, false
}
