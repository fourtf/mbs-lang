package common

type Expr interface {
	Print() string
	Eval()
}

type Block struct {
	Statements []Expr
}

type ReadVar struct {
	Name string
}

func (v ReadVar) Print() string {
	return v.Name
}
func (v ReadVar) Eval() {}

type WriteVar struct {
	Name string
	Expr Expr
}

func (v WriteVar) Print() string {
	return v.Name + " = " + v.Expr.Print()
}
func (v WriteVar) Eval() {}

type Operator struct {
	Symbol    string
	FirstExp  Expr
	SecondExp Expr
}

func (op Operator) Print() string {
	return op.FirstExp.Print() + " " + op.Symbol + " " + op.SecondExp.Print()
}
func (op Operator) Eval() {}

type FunctionCall struct {
	Name     string
	Argument Expr // TODO: we only allow one argument
}

func (f FunctionCall) Print() string {
	return f.Name + "(" + f.Argument.Print() + ")"
}
func (f FunctionCall) Eval() {}
