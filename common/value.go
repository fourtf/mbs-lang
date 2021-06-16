package common

import "strconv"

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

func (s String) Print() string { return `"` + s.Data + `"` }
func (s String) Eval()         {}

type Integer struct {
	Data int64
}

func (i Integer) Print() string { return strconv.FormatInt(i.Data, 10) }
func (i Integer) Eval()         {}

type Float struct {
	Data float64
}

func (f Float) Print() string { return strconv.FormatFloat(f.Data, 'f', 5, 64) }
func (f Float) Eval()         {}
