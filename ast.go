package trompe

type Token struct {
	Loc  Loc
	Text string
}

type CommentNode struct {
	Begin Loc
	Text  Token
}

type Node interface {
	Loc() *Loc
}

type ChunkNode struct {
	loc   Loc
	Block BlockNode
}

type BlockNode struct {
	loc   Loc
	Stats []Node
}

type StatNode interface {
	Node
}

type LetStatNode struct {
	Let Loc
	Ptn PtnNode
	Eq  Loc
	Exp ExpNode
}

type DefStatNode struct {
	Def   Loc
	Name  Token
	Open  Loc
	Args  ArgListNode
	Close Loc
	Block BlockNode
	End   Loc
}

type ShortDefStatNode struct {
	Def   Loc
	Name  Token
	Open  Loc
	Args  ArgListNode
	Close Loc
	Eq    Loc
	Exp   ExpNode
}

type ArgListNode struct {
	Names []Token
	Sep   []Loc
}

func (args *ArgListNode) NameStrs() []string {
	nameStrs := make([]string, len(args.Names))
	for _, tok := range args.Names {
		nameStrs = append(nameStrs, tok.Text)
	}
	return nameStrs
}

type IfStatNode struct {
	Cond []IfCondNode
	Else *ElseStatNode
	End  Loc
}

type IfCondNode struct {
	If     Loc
	Cond   ExpNode
	Then   Loc
	Action BlockNode
}

type ElseStatNode struct {
	Else   Loc
	Action BlockNode
}

type CaseStatNode struct {
	Case  Loc
	Cond  ExpNode
	Claus []CaseClauNode
	Else  *ElseStatNode
}

type CaseClauNode struct {
	Ptn    PtnNode
	In     Loc
	Action *BlockNode
}

type RetStatNode struct {
	Ret   Loc
	Value ExpNode
}

type FunCallStatNode struct {
	Exp *FunCallExpNode
}

type ExpNode interface {
	Node
}

type FunCallExpNode struct {
	Prefix Node
	Args   EltListNode
}

type CondOpExpNode struct {
	Colon Loc
	Q     Loc
	Cond  ExpNode
	True  ExpNode
	False ExpNode
}

type VarExpNode struct {
	Name Token
}

type UnitExpNode struct {
	Open  Loc
	Close Loc
}

type BoolExpNode struct {
	loc   Loc
	Value bool
}

type IntExpNode struct {
	loc   Loc
	Value string
}

type StrExpNode struct {
	loc   Loc
	Value string
}

type ListExpNode struct {
	Elts EltListNode
}

type EltListNode struct {
	Open  Loc
	Close Loc
	Elts  []Node
	Seps  []Loc
}

type TupleExpNode struct {
	Elts EltListNode
}

type SomeExpNode struct {
	SomeLoc Loc
	Value   ExpNode
}

type NoneExpNode struct {
	loc Loc
}

type AnonFunExpNode struct {
	Open    Loc
	Close   Loc
	Args    []Node
	ArgSeps []Loc
	In      Loc
	Block   BlockNode
}

type PtnNode interface {
	Node
}

type UnitPtnNode struct {
	Open  Loc
	Close Loc
}

type BoolPtnNode struct {
	Text Token
}

type IntPtnNode struct {
	Text Token
}

type StrPtnNode struct {
	Text Token
}

type ListPtnNode struct {
	Elts EltPtnListNode
}

type ConsPtnNode struct {
	Left  PtnNode
	Sep   Loc
	Right PtnNode
}

type TuplePtnNode struct {
	Elts EltPtnListNode
}

type EltPtnListNode struct {
	Open  Loc
	Close Loc
	Elts  []PtnNode
	Seps  []Loc
}

func (chunk *ChunkNode) Loc() *Loc {
	return &chunk.loc
}

func (block *BlockNode) Loc() *Loc {
	return &block.loc
}

func (stat *LetStatNode) Loc() *Loc {
	return &stat.Let
}

func (stat *DefStatNode) Loc() *Loc {
	return &stat.Def
}

func (stat *ShortDefStatNode) Loc() *Loc {
	return &stat.Def
}

func (stat *IfStatNode) Loc() *Loc {
	return &stat.Cond[0].If
}

func (stat *CaseStatNode) Loc() *Loc {
	return &stat.Case
}

func (stat *RetStatNode) Loc() *Loc {
	return &stat.Ret
}

func (stat *FunCallStatNode) Loc() *Loc {
	return stat.Exp.Loc()
}

func (exp *FunCallExpNode) Loc() *Loc {
	return exp.Prefix.Loc()
}

func (exp *CondOpExpNode) Loc() *Loc {
	return &exp.Colon
}

func (exp *VarExpNode) Loc() *Loc {
	return &exp.Name.Loc
}

func (exp *UnitExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *BoolExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *IntExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *StrExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *ListExpNode) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *TupleExpNode) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *SomeExpNode) Loc() *Loc {
	return &exp.SomeLoc
}

func (exp *NoneExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *AnonFunExpNode) Loc() *Loc {
	return &exp.Open
}

func (ptn *UnitPtnNode) Loc() *Loc {
	return &ptn.Open
}

func (ptn *BoolPtnNode) Loc() *Loc {
	return &ptn.Text.Loc
}

func (ptn *IntPtnNode) Loc() *Loc {
	return &ptn.Text.Loc
}
func (ptn *StrPtnNode) Loc() *Loc {
	return &ptn.Text.Loc
}
func (ptn *ListPtnNode) Loc() *Loc {
	return &ptn.Elts.Open
}

func (ptn *ConsPtnNode) Loc() *Loc {
	return ptn.Left.Loc()
}

func (ptn *TuplePtnNode) Loc() *Loc {
	return &ptn.Elts.Open
}
