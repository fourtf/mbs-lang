package common

type Value interface {
	Print() string
}

type Boolean struct {
	Data bool
}

func (b Boolean) Print() string {
	if b.Data {
		return "true"
	} else {
		return "false"
	}
}

func (boolean Boolean) Eval() {}

type String struct {
	Data string
}

func (str String) Print() string { return "" }
func (str String) Eval()         {}

type Integer struct {
	Data int64
}

func (integer Integer) Print() string { return "" }
func (integer Integer) Eval()         {}

type Float struct {
	Data float64
}

func (float Float) Print() string { return "" }
func (float Float) Eval()         {}
