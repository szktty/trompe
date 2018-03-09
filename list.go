package trompe

type List struct {
	Value interface{}
	Next  *List
}

var ListNil = &List{nil, nil}

func CreateList(value interface{}) *List {
	return &List{Value: value, Next: nil}
}

func (l *List) Cons(value interface{}) *List {
	return &List{Value: value, Next: l}
}
