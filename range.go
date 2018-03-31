package trompe

import (
	"fmt"
)

type Range struct {
	Start int
	End   int
	Close bool
}

func (r *Range) Type() int {
	return ValRangeType
}

func (r *Range) Desc() string {
	if r.Close {
		return fmt.Sprintf("%d..<%d", r.Start, r.End)
	} else {
		return fmt.Sprintf("%d...%d", r.Start, r.End)
	}
}

func (r *Range) Bool() bool {
	panic("range")
}

func (r *Range) Int() int {
	panic("range")
}

func (r *Range) String() string {
	panic("range")
}

func (r *Range) Closure() Closure {
	panic("range")
}

func (r *Range) List() *List {
	panic("range")
}

func (r *Range) Tuple() []Value {
	panic("range")
}
