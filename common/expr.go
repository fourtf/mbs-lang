package common

type Expr interface {
	Print() string
	Eval()
}

type Block struct {
}
