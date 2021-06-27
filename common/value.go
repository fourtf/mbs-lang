package common

import "strconv"

/*In here are all the primitive data types that our language supports*/

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

func (b Boolean) Eval() interface{} {
	return b.Data
}

func (b Boolean) Type() Type {
	return BooleanType
}

type String struct {
	Data string
}

func (s String) Print() string { return `"` + s.Data + `"` }
func (s String) Eval() interface{} {
	return s.Data
}

func (s String) Type() Type {
	return StringType
}

type Integer struct {
	Data int64
}

func (i Integer) Print() string { return strconv.FormatInt(i.Data, 10) }
func (i Integer) Eval() interface{} {
	return i.Data
}
func (i Integer) Type() Type {
	return IntegerType
}

type Float struct {
	Data float64
}

func (f Float) Print() string { return strconv.FormatFloat(f.Data, 'f', 5, 64) }
func (f Float) Eval() interface{} {
	return f.Data
}
func (f Float) Type() Type {
	return FloatType
}
