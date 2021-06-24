package common

import (
	"fmt"
	"strings"
)

type Type string

const (
	BlockType        Type = "Block"
	ReadVarType      Type = "ReadVar"
	WriteVarType     Type = "WriteVar"
	OperatorType     Type = "Operator"
	FunctionCallType Type = "FunctionCall"
	IfType           Type = "If"
	ForType          Type = "For"
	NopType          Type = "Nop"
	BooleanType      Type = "Boolean"
	IntegerType      Type = "Integer"
	FloatType        Type = "Float"
	StringType       Type = "String"
)

type Expr interface {
	Print() string
	Eval()
	Type() Type
}

type Block struct {
	Statements []Expr
}

func (b Block) Print() string {
	bld := strings.Builder{}

	for _, stmt := range b.Statements {

		bld.WriteString(stmt.Print())
		bld.WriteString("\n")
	}

	return bld.String()
}

func (b Block) Eval() {}

func (b Block) Type() Type {
	return BlockType
}

type ReadVar struct {
	Name string
}

func (v ReadVar) Print() string {
	return v.Name
}
func (v ReadVar) Eval() {}

func (v ReadVar) Type() Type {
	return ReadVarType
}

type WriteVar struct {
	Name string
	Expr Expr
}

func (v WriteVar) Print() string {
	return v.Name + " = " + v.Expr.Print()
}
func (v WriteVar) Eval() {}

func (v WriteVar) Type() Type {
	return WriteVarType
}

type Operator struct {
	Symbol    string
	FirstExp  Expr
	SecondExp Expr
}

func (op Operator) Print() string {
	return op.FirstExp.Print() + " " + op.Symbol + " " + op.SecondExp.Print()
}
func (op Operator) Eval() {}

func (op Operator) Type() Type {
	return OperatorType
}

type FunctionCall struct {
	Name     string
	Argument Expr // TODO: we only allow one argument
}

func (f FunctionCall) Print() string {
	return f.Name + "(" + f.Argument.Print() + ")"
}
func (f FunctionCall) Eval() {}

func (f FunctionCall) Type() Type {
	return FunctionCallType
}

type If struct {
	Condition Expr
	Body      *Block
}

func (i If) Print() string {
	return "if (" + i.Condition.Print() + ") {\n" + i.Body.Print() + "}\n"
}

func (i If) Eval() {}

func (i If) Type() Type {
	return IfType
}

type For struct {
	Init        Expr
	Condition   Expr
	Advancement Expr
	Body        *Block
}

func (i For) Print() string {
	return fmt.Sprintf("for (%s; %s; %s) {\n%s}", i.Init.Print(), i.Condition.Print(), i.Advancement.Print(), i.Body.Print())
}

func (i For) Eval() {}

func (i For) Type() Type {
	return ForType
}

// Nop is used whenever a statement or expression doesn't do anything e.g. empty values in a for-loop (for (;;)).
type Nop struct{}

func (i Nop) Print() string {
	return "nop"
}

func (i Nop) Eval() {}

func (i Nop) Type() Type {
	return NopType
}
