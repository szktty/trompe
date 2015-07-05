package trompe

import (
	"fmt"
)

type Pos struct {
	Line   int
	Col    int
	Offset int
}

type Loc struct {
	File  string
	Start *Pos
	End   *Pos
	Len   int
}

type WithLoc interface {
	Loc() *Loc
}

func NewLoc(start *Pos, end *Pos) *Loc {
	l := &Loc{}
	l.Init(start, end)
	return l
}

func ZeroLoc() *Loc {
	return NewLoc(&Pos{}, &Pos{})
}

func (l *Loc) Init(start *Pos, end *Pos) {
	l.Start = start
	l.End = end
	l.Len = end.Offset - start.Offset
}

func (l *Loc) StartString() string {
	return fmt.Sprintf("line %d, col %d", l.Start.Line+1, l.Start.Col+1)
}

func (l *Loc) String() string {
	return fmt.Sprintf("{%d,%d,%d,%d,%d,%d,%d}",
		l.Start.Line, l.Start.Col, l.Start.Offset,
		l.End.Line, l.End.Col, l.End.Offset, l.Len)
}

func (l1 *Loc) Union(l2 *Loc) *Loc {
	// TODO
	return l1
}
