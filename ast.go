package trompe

type Comment struct {
	Loc  Loc
	Text string
}

type Node interface {
	Loc() *Loc
}

type Chunk struct {
	loc   Loc
	Block Block
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
	Exp Exp
}

type DefStat struct {
	Def   Loc
	Name  StrTok
	Open  Loc
	Args  ArgNodeList
	Close Loc
	Block Block
	End   Loc
}

type ShortDefStat struct {
	Def   Loc
	Name  StrTok
	Open  Loc
	Args  ArgNodeList
	Close Loc
	Eq    Loc
	Exp   Exp
}

type ArgNodeList struct {
	Names []StrTok
	Sep   []Loc
}

func (args *ArgNodeList) NameStrs() []string {
	nameStrs := make([]string, len(args.Names))
	for _, tok := range args.Names {
		nameStrs = append(nameStrs, tok.Value)
	}
	return nameStrs
}

type IfStat struct {
	Cond []IfCond
	Else *ElseStat
	End  Loc
}

type IfCond struct {
	If     Loc
	Cond   Exp
	Then   Loc
	Action Block
}

type ElseStat struct {
	Else   Loc
	Action Block
}

type CaseStat struct {
	Case  Loc
	Cond  Exp
	Claus []CaseClau
	Else  *ElseStat
}

type CaseClau struct {
	Ptn    PtnExp
	In     Loc
	Action *Block
}

type RetStat struct {
	Ret   Loc
	Value Exp
}

type FunCallStat struct {
	Exp *FunCallExp
}

type Exp interface {
	Node
}

type FunCallExp struct {
	Prefix Node
	Args   EltList
}

type CondOpExp struct {
	Colon Loc
	Q     Loc
	Cond  Exp
	True  Exp
	False Exp
}

type VarExp struct {
	loc  Loc
	Name string
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
	Node
}

func (chunk *Chunk) Loc() *Loc {
	return &chunk.loc
}

func (block *Block) Loc() *Loc {
	return &block.loc
}

func (stat *LetStat) Loc() *Loc {
	return &stat.Let
}

func (stat *DefStat) Loc() *Loc {
	return &stat.Def
}

func (stat *ShortDefStat) Loc() *Loc {
	return &stat.Def
}

func (stat *IfStat) Loc() *Loc {
	return &stat.Cond[0].If
}

func (stat *CaseStat) Loc() *Loc {
	return &stat.Case
}

func (stat *RetStat) Loc() *Loc {
	return &stat.Ret
}

func (stat *FunCallStat) Loc() *Loc {
	return stat.Exp.Loc()
}

func (exp *FunCallExp) Loc() *Loc {
	return exp.Prefix.Loc()
}

func (exp *CondOpExp) Loc() *Loc {
	return &exp.Colon
}

func (exp *VarExp) Loc() *Loc {
	return &exp.loc
}

func (exp *UnitExp) Loc() *Loc {
	return &exp.loc
}

func (exp *BoolExp) Loc() *Loc {
	return &exp.loc
}

func (exp *IntExp) Loc() *Loc {
	return &exp.loc
}

func (exp *StrExp) Loc() *Loc {
	return &exp.loc
}

func (exp *ListExp) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *TupleExp) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *SomeExp) Loc() *Loc {
	return &exp.SomeLoc
}

func (exp *NoneExp) Loc() *Loc {
	return &exp.loc
}

func (exp *AnonFunExp) Loc() *Loc {
	return &exp.Open
}
