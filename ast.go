package trompe

import (
	"bytes"
	"fmt"
)

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
	Write(*bytes.Buffer)
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
	Def    Loc
	Name   Token
	Open   Loc
	Params *ParamListNode
	Close  Loc
	Block  BlockNode
	End    Loc
}

type ShortDefStatNode struct {
	Def    Loc
	Name   Token
	Open   Loc
	Params ParamListNode
	Close  Loc
	Eq     Loc
	Exp    ExpNode
}

type ParamListNode struct {
	Names []Token
	Sep   []Loc
}

func (params *ParamListNode) NameStrs() []string {
	nameStrs := make([]string, len(params.Names))
	for _, tok := range params.Names {
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
	Ret Loc
	Exp ExpNode
}

type ExpNode interface {
	Node
}

type ParenExpNode struct {
	Open  Loc
	Close Loc
	Exp   Node
}

type FunCallExpNode struct {
	Callable Node
	Args     EltListNode
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
	Value Token
}

type StrExpNode struct {
	Value Token
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
	Open   Loc
	Close  Loc
	Params ParamListNode
	In     Loc
	Block  BlockNode
}

type PtnNode interface {
	Node
}

type UnitPtnNode struct {
	Open  Loc
	Close Loc
}

type BoolPtnNode struct {
	loc   Loc
	Value bool
}

type IntPtnNode struct {
	Value Token
}

type StrPtnNode struct {
	Value Token
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

func (chunk *ChunkNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(chunk ")
	chunk.Block.Write(buf)
	buf.WriteString(")")
}

func (block *BlockNode) Loc() *Loc {
	return &block.loc
}

func (block *BlockNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(block [")
	for _, stat := range block.Stats {
		stat.Write(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("])")
}

func (stat *LetStatNode) Loc() *Loc {
	return &stat.Let
}

func (stat *LetStatNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(let ")
	stat.Ptn.Write(buf)
	buf.WriteString(" ")
	stat.Exp.Write(buf)
	buf.WriteString(")")
}

func (stat *DefStatNode) Loc() *Loc {
	return &stat.Def
}

func (stat *DefStatNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(def \"%s\" [", stat.Name.Text))
	if stat.Params != nil {
		stat.Params.Write(buf)
	}
	buf.WriteString("])")
}

func (stat *ShortDefStatNode) Loc() *Loc {
	return &stat.Def
}

func (stat *ShortDefStatNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(shortdef \"%s\" ", stat.Name.Text))
	stat.Exp.Write(buf)
	buf.WriteString(")")
}

func (params *ParamListNode) Loc() *Loc {
	return &params.Names[0].Loc
}

func (params *ParamListNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(params [")
	for _, param := range params.Names {
		buf.WriteString(fmt.Sprintf("\"%s\" ", param.Text))
	}
	buf.WriteString("])")
}

func (elts *EltListNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(eltlist [")
	for _, elt := range elts.Elts {
		elt.Write(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("])")
}

func (stat *IfStatNode) Loc() *Loc {
	return &stat.Cond[0].If
}

func (stat *IfStatNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(if \"%s\" [")
	for _, cond := range stat.Cond {
		cond.Write(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("] ")
	if else_ := stat.Else; else_ != nil {
		else_.Write(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (stat *CaseStatNode) Loc() *Loc {
	return &stat.Case
}

func (stat *CaseStatNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(case ")
	stat.Cond.Write(buf)
	buf.WriteString(" ")
	for _, clau := range stat.Claus {
		clau.Write(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("] ")
	if else_ := stat.Else; else_ != nil {
		else_.Write(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (stat *RetStatNode) Loc() *Loc {
	return &stat.Ret
}

func (stat *RetStatNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(return ")
	if exp := stat.Exp; exp != nil {
		exp.Write(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (exp *ParenExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *ParenExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(paren ")
	exp.Exp.Write(buf)
	buf.WriteString(")")
}

func (exp *FunCallExpNode) Loc() *Loc {
	return exp.Callable.Loc()
}

func (exp *FunCallExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(funcall ")
	exp.Callable.Write(buf)
	buf.WriteString("  ")
	exp.Args.Write(buf)
	buf.WriteString(")")
}

func (exp *CondOpExpNode) Loc() *Loc {
	return &exp.Colon
}

func (exp *CondOpExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(condop")
	exp.Cond.Write(buf)
	buf.WriteString("  ")
	exp.True.Write(buf)
	buf.WriteString("  ")
	exp.False.Write(buf)
	buf.WriteString(")")
}

func (exp *VarExpNode) Loc() *Loc {
	return &exp.Name.Loc
}

func (exp *VarExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(var \"%s\")", exp.Name.Text))
}

func (exp *UnitExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *UnitExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(unit)")
}

func (exp *BoolExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *BoolExpNode) Write(buf *bytes.Buffer) {
	var value string
	if exp.Value {
		value = "true"
	} else {
		value = "false"
	}
	buf.WriteString(fmt.Sprintf("(bool \"%s\")", value))
}

func (exp *IntExpNode) Loc() *Loc {
	return &exp.Value.Loc
}

func (exp *IntExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(int \"%s\")", exp.Value.Text))
}

func (exp *StrExpNode) Loc() *Loc {
	return &exp.Value.Loc
}

func (exp *StrExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(str \"%s\")", exp.Value.Text))
}

func (exp *ListExpNode) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *ListExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(list ")
	exp.Elts.Write(buf)
	buf.WriteString(")")
}

func (exp *TupleExpNode) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *TupleExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(tuple ")
	exp.Elts.Write(buf)
	buf.WriteString(")")
}

func (exp *EltListNode) Loc() *Loc {
	return &exp.Open
}

func (elts *EltListNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(eltlist [")
	for _, elt := range elts.Elts {
		elt.Write(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("])")
}

func (exp *SomeExpNode) Loc() *Loc {
	return &exp.SomeLoc
}

func (exp *SomeExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(some ")
	exp.Value.Write(buf)
	buf.WriteString(")")
}

func (exp *NoneExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *NoneExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(none)")
}

func (exp *AnonFunExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *AnonFunExpNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(anonfun ")
	exp.Params.Write(buf)
	buf.WriteString(" ")
	exp.Block.Write(buf)
	buf.WriteString(")")
}

func (ptn *UnitPtnNode) Loc() *Loc {
	return &ptn.Open
}

func (ptn *UnitPtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(unitptn)")
}

func (ptn *BoolPtnNode) Loc() *Loc {
	return &ptn.loc
}

func (ptn *BoolPtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(bool ")
	if ptn.Value {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	buf.WriteString(")")
}

func (ptn *IntPtnNode) Loc() *Loc {
	return &ptn.Value.Loc
}

func (ptn *IntPtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(intptn \"%s\")", ptn.Value.Text))
}

func (ptn *StrPtnNode) Loc() *Loc {
	return &ptn.Value.Loc
}

func (ptn *StrPtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(strptn \"%s\")", ptn.Value.Text))
}

func (ptn *ListPtnNode) Loc() *Loc {
	return &ptn.Elts.Open
}

func (ptn *ListPtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(listptn ")
	ptn.Elts.Write(buf)
	buf.WriteString(")")
}

func (ptn *ConsPtnNode) Loc() *Loc {
	return ptn.Left.Loc()
}

func (ptn *ConsPtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(consptn ")
	ptn.Left.Write(buf)
	buf.WriteString(" ")
	ptn.Right.Write(buf)
	buf.WriteString(")")
}

func (ptn *TuplePtnNode) Loc() *Loc {
	return &ptn.Elts.Open
}

func (ptn *TuplePtnNode) Write(buf *bytes.Buffer) {
	buf.WriteString("(tupleptn ")
	ptn.Elts.Write(buf)
	buf.WriteString(")")
}
