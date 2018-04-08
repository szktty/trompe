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
	return ValueTypeRange
}

func (r *Range) Desc() string {
	if r.Close {
		return fmt.Sprintf("%d..<%d", r.Start, r.End)
	} else {
		return fmt.Sprintf("%d...%d", r.Start, r.End)
	}
}

type RangeIter struct {
	cur int
	end int
}

func (r *Range) NewIter() Iter {
	if r.Close {
		return &RangeIter{r.Start, r.End}
	} else {
		return &RangeIter{r.Start, r.End - 1}
	}
}

func (iter *RangeIter) Type() int {
	return ValueTypeIter
}

func (iter *RangeIter) Desc() string {
	return fmt.Sprintf("Iter(%d...%d)", iter.cur, iter.end)
}

func (iter *RangeIter) Next() Value {
	iter.cur++
	if iter.cur > iter.end {
		return NewInt(iter.cur)
	} else {
		return nil
	}
}
