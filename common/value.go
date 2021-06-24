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

func (b Boolean) Type() Type {
	return BooleanType
}

type String struct {
	Data string
}

func (s String) Print() string { return `"` + s.Data + `"` }
func (s String) Eval()         {}

func (s String) Type() Type {
	return StringType
}

type Integer struct {
	Data int64
}

func (i Integer) Print() string { return strconv.FormatInt(i.Data, 10) }
func (i Integer) Eval()         {}
func (i Integer) Type() Type {
	return IntegerType
}

type Float struct {
	Data float64
}

func (f Float) Print() string { return strconv.FormatFloat(f.Data, 'f', 5, 64) }
func (f Float) Eval()         {}
func (f Float) Type() Type {
	return FloatType
}
