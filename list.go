package trompe

import (
	"fmt"
)

type List struct {
	Value Value
	Next  *List
}

var ListNil = &List{nil, nil}

func NewList(value Value) *List {
	return &List{Value: value, Next: nil}
}

func (l *List) Len() int {
	i := 0
	for l.Next != nil {
		l = l.Next
		i++
	}
	return i
}

func (l *List) Cons(value Value) *List {
	return &List{Value: value, Next: l}
}

// interface Value

func (l *List) Type() int {
	return ValueTypeList
}

func (l *List) Desc() string {
	// TODO
	return fmt.Sprintf("list")
}
