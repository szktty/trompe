package trompe

type Comment struct {
	Loc  Loc
	Text string
}

type Node interface {
	Loc() Loc
}

type Chunk struct {
	loc   Loc
	Block Node
}

type Block struct {
	loc   Loc
	Stats []Node
}

type Stat interface {
	Node
}

type LetStat struct {
	Let Loc
	Ptn PtnExp
	Eq  Loc
}

type IfStat struct {
	Cond []IfCond
	End  Loc
}

type IfCond struct {
	If   Loc
	Cond Node
	Then Loc
}

type RetStat struct {
	Ret   Loc
	Value Node
}

type FunCall struct {
	Prefix Node
	Args   EltList
}

type Exp interface {
	Node
}

type UnitExp struct {
	loc Loc
}

type BoolExp struct {
	loc   Loc
	Value bool
}

type IntExp struct {
	loc   Loc
	Value string
}

type StrExp struct {
	loc   Loc
	Value string
}

type ListExp struct {
	Elts EltList
}

type EltList struct {
	Open  Loc
	Close Loc
	Elts  []Node
	Seps  []Loc
}

type TupleExp struct {
	Elts EltList
}

type SomeExp struct {
	SomeLoc Loc
	Value   Exp
}

type NoneExp struct {
	loc Loc
}

type AnonFunExp struct {
	Open    Loc
	Close   Loc
	Args    []Node
	ArgSeps []Loc
	In      Loc
	Block   Block
}

type PtnExp interface {
}

func (exp *UnitExp) Loc() *Loc {
	return &exp.loc
}
