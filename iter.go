package trompe

type Iter interface {
	Value
	Next() Value
}

func NewIter(val Value) Iter {
	switch val := val.(type) {
	case *Range:
		return val.NewIter()
	default:
		panic("unsupported")
	}
}
