package common

import (
	"fmt"
	"strings"
)

/*In here all the non primitive types that our AST can contain are stored.
These types are also defining the code execution by implementing the "Expr"-Interface*/

type Type string

// all of the expressions that can occur in our AST
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

// stores all variables and their values that can be accessed in the current scope
var variables map[string]interface{} = make(map[string]interface{})

// the interface that every expression that can occur in our AST implements
type Expr interface {
	Print() string
	Eval() interface{} // used to execute the code the AST represents
	Type() Type        // the typechecker uses this to easily access the type of an expression
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

func (b Block) Eval() interface{} {
	// saving the variables of the outer scope
	outerscopeVars := map[string]interface{}{}
	for k, v := range variables {
		outerscopeVars[k] = v
	}
	// executing the code inside the block
	for _, expr := range b.Statements {
		expr.Eval()
	}
	// restoring the variables after exiting the block to "delete" the variables defined in the scope of the current block
	variables = outerscopeVars
	return nil
}

func (b Block) Type() Type {
	return BlockType
}

type ReadVar struct {
	Name string
}

func (v ReadVar) Print() string {
	return v.Name
}

func (v ReadVar) Eval() interface{} {
	return variables[v.Name]
}

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

func (v WriteVar) Eval() interface{} {
	variables[v.Name] = v.Expr.Eval()
	return nil
}

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
func (op Operator) Eval() interface{} {
	// getting the primitive value of both expressions
	firstExp := op.FirstExp.Eval()
	secondExp := op.SecondExp.Eval()

	// performing the operation
	switch operator := op.Symbol; operator {
	case "+":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) + secondExp.(int64)
			case float64:
				return firstExp.(float64) + secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) + secondExp.(float64)
		case string:
			return firstExp.(string) + secondExp.(string)
		}
	case "-":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) - secondExp.(int64)
			case float64:
				return firstExp.(float64) - secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) - secondExp.(float64)
		}
	case "*":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) * secondExp.(int64)
			case float64:
				return firstExp.(float64) * secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) * secondExp.(float64)
		}

	case "/":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) / secondExp.(int64)
			case float64:
				return firstExp.(float64) / secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) / secondExp.(float64)
		}
	case "==":
		return firstExp == secondExp
	case "!=":
		return firstExp != secondExp
	case ">":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) > secondExp.(int64)
			case float64:
				return firstExp.(float64) > secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) > secondExp.(float64)
		}
	case "<":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) < secondExp.(int64)
			case float64:
				return firstExp.(float64) < secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) < secondExp.(float64)
		}
	case ">=":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) >= secondExp.(int64)
			case float64:
				return firstExp.(float64) >= secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) >= secondExp.(float64)
		}
	case "<=":
		switch firstExp.(type) {
		case int64:
			switch secondExp.(type) {
			case int64:
				return firstExp.(int64) <= secondExp.(int64)
			case float64:
				return firstExp.(float64) <= secondExp.(float64)
			}
		case float64:
			return firstExp.(float64) <= secondExp.(float64)
		}
	case "&&":
		return firstExp.(bool) && secondExp.(bool)
	case "||":
		return firstExp.(bool) && secondExp.(bool)
	}
	return nil
}

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

func (f FunctionCall) Eval() interface{} {
	//we only allow the functions println and readln
	if f.Name == "println" {
		println(f.Argument.Eval().(string))
	} else if f.Name == "readln" {
		var input string
		fmt.Scanf("%s", &input)
		return input
	}
	return nil
}

func (f FunctionCall) Type() Type {
	return FunctionCallType
}

type If struct {
	Condition Expr
	Body      Block
}

func (i If) Print() string {
	return "if (" + i.Condition.Print() + ") {\n" + i.Body.Print() + "}\n"
}

func (i If) Eval() interface{} {
	if i.Condition.Eval().(bool) {
		i.Body.Eval()
	}
	return nil
}

func (i If) Type() Type {
	return IfType
}

type For struct {
	Init        Expr
	Condition   Expr
	Advancement Expr
	Body        Block
}

func (f For) Print() string {
	return fmt.Sprintf("for (%s; %s; %s) {\n%s}", f.Init.Print(), f.Condition.Print(), f.Advancement.Print(), f.Body.Print())
}

func (f For) Eval() interface{} {
	for f.Init.Eval(); f.Condition.Eval().(bool); f.Advancement.Eval() {
		f.Body.Eval()
	}
	return nil
}

func (f For) Type() Type {
	return ForType
}

// Nop is used whenever a statement or expression doesn't do anything e.g. empty values in a for-loop (for (;;)).
type Nop struct{}

func (i Nop) Print() string {
	return "nop"
}

func (i Nop) Eval() interface{} {
	return nil
}

func (i Nop) Type() Type {
	return NopType
}
