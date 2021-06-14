package common

type Expr interface {
	Print() string
	Eval()
}

type Block struct {
}

type WriteVar struct {
	Name string
	Expr Expr
}

func (writeVar WriteVar) Print() string { return "" }
func (WriteVar WriteVar) Eval()         {}

type Operator struct {
	Symbol    string
	FirstExp  Expr
	SecondExp Expr
}

func (op Operator) Print() string { return "" }
func (op Operator) Eval()         {}
