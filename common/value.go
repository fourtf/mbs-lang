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
