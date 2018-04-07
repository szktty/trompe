package trompe

import (
	"fmt"
)

type Range struct {
	Start int
	End   int
	Close bool
}

func NewRange(start int, end int, close bool) *Range {
	return &Range{start, end, close}
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

func (r *Range) Iter() ValIter {
	return nil
}

type RangeIter struct {
	cur int
	end int
}

func (r *Range) NewValIter() ValIter {
	if r.Close {
		return &RangeIter{r.Start, r.End}
	} else {
		return &RangeIter{r.Start, r.End - 1}
	}
}

func (iter *RangeIter) Type() int {
	return ValRangeType
}

func (iter *RangeIter) Desc() string {
	return fmt.Sprintf("Iter(%d...%d)", iter.cur, iter.end)
}

func (iter *RangeIter) Bool() bool {
	panic("range")
}

func (iter *RangeIter) Int() int {
	panic("range")
}

func (iter *RangeIter) String() string {
	panic("range")
}

func (iter *RangeIter) Closure() Closure {
	panic("range")
}

func (iter *RangeIter) List() *List {
	panic("range")
}

func (iter *RangeIter) Tuple() []Value {
	panic("range")
}

func (iter *RangeIter) Iter() ValIter {
	return iter
}

func (iter *RangeIter) Next() Value {
	iter.cur++
	if iter.cur > iter.end {
		return &ValInt{iter.cur}
	} else {
		return nil
	}
}
