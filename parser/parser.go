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
	wv := WriteVar{}
	code, err := sequence(name(&wv.Name), token("="), expr(&wv.Expr))(code)

	if err != nil {
		return code, nil, err
	}

	return code, wv, err
}

func ParseExpression(code string) (string, Expr, error) {
	if remainingCode, exp, err := ParseOperator(code); err == nil {
		return remainingCode, exp, nil
	}
	if remainingCode, exp, err := ParseExpressionWithoutOperator(code); err == nil {
		return remainingCode, exp, nil
	}
	return code, nil, &ParseError{Message: "Couldn't parse any expression"}
}

func ParseExpressionWithoutOperator(code string) (string, Expr, error) {
	var e Expr = nil

	code, err := alternative(
		pfunc(&e, ParseParentheses),
		pfunc(&e, ParseString),
		pfunc(&e, ParseFloat),
		pfunc(&e, ParseInteger),
		pfunc(&e, ParseFunctionCall),
		pfunc(&e, ParseBoolean),
		pfunc(&e, ParseReadVar),
	)(code)

	if err != nil {
		return code, nil, &ParseError{Message: "Couldn't parse any expression"}
	}

	return code, e, nil
}

var (
	stringRegex       = regexp.MustCompile(`^"((\\"|[^"])*)"`)
	stringEscapeRegex = regexp.MustCompile(`\\"`)
)

func ParseString(code string) (string, Expr, error) {
	code = stripWhitespaceLeft(code)
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
	code = stripWhitespaceLeft(code)
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
	code = stripWhitespaceLeft(code)
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
	code = stripWhitespaceLeft(code)
	match := floatRegex.FindString(code)
	if float, err := strconv.ParseFloat(match, 64); err == nil {
		return code[len(match):], Float{Data: float}, nil
	}

	return code, nil, &ParseError{Message: "Couldn't parse the expression to a Float"}
}

func ParseFunctionCall(code string) (string, Expr, error) {
	fn := FunctionCall{Argument: Nop{}}
	code, err := sequence(name(&fn.Name), token("("), opt(expr(&fn.Argument)), token(")"))(code)

	return code, fn, err
}

func ParseParentheses(code string) (string, Expr, error) {
	var res Expr = nil
	code, err := sequence(token("("), expr(&res), token(")"))(code)

	return code, res, err
}

var (
	operators = []string{"+", "-", "*", "/", ">", "<", "==", "!=", ">=", "<=", "&&", "||"}
)

func ParseOperator(code string) (string, Expr, error) {
	code = stripWhitespaceLeft(code)
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
	if_ := If{}
	code, err := sequence(token("if"), token("("), expr(&if_.Condition), token(")"), token("{"), block(if_.Body), token("}"))(code)

	if err != nil {
		return code, nil, err
	}

	return code, if_, nil
}

func ParseFor(code string) (string, Expr, error) {
	for_ := For{Init: &Nop{}, Condition: &Nop{}, Advancement: &Nop{}}

	code, err := sequence(
		token("for"),
		token("("),
		opt(pfunc(&for_.Init, ParseWriteVar)),
		token(";"),
		expr(&for_.Condition),
		token(";"),
		opt(pfunc(&for_.Advancement, ParseWriteVar)),
		token(")"),
		token("{"),
		block(for_.Body),
		token("}"))(code)

	if err != nil {
		return code, nil, err
	}

	return code, &for_, nil
}

func ParseBlock(code string) (string, *Block, error) {
	// Either:
	// - WriteVar
	// - FunctionCall
	// - If
	// - For

	stmts := make([]Expr, 0)

	for {
		var e Expr = nil

		tmp, err := alternative(
			sequence(pfunc(&e, ParseWriteVar), token(";")),
			sequence(pfunc(&e, ParseFunctionCall), token(";")),
			pfunc(&e, ParseIf),
			pfunc(&e, ParseFor),
		)(code)

		if err != nil {
			return code, &Block{Statements: stmts}, nil
		}

		code = tmp
		stmts = append(stmts, e)
	}
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

var stripWhitespaceRegex = regexp.MustCompile(`^\s+`)

// stripWhitespaceLeft strips all whitespace on the left of the string and returns a string without it.
func stripWhitespaceLeft(s string) string {
	loc := stripWhitespaceRegex.FindStringIndex(s)

	if loc == nil {
		return s
	}

	return s[loc[1]:]
}
