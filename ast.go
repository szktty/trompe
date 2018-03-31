package trompe

import (
	"bytes"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
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
	WriteTo(*bytes.Buffer)
}

type ChunkNode struct {
	loc   Loc
	Block *BlockNode
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
	Params *ParamListNode
	Close  Loc
	Eq     Loc
	Exp    ExpNode
}

type ParamListNode struct {
	Names []Token
	Sep   []Loc
}

type IfStatNode struct {
	Cond       []IfCondNode
	Else       *Loc
	ElseAction *BlockNode
	End        Loc
}

type IfCondNode struct {
	If     Loc
	Cond   ExpNode
	Then   Loc
	Action BlockNode
}

type CaseStatNode struct {
	Case       Loc
	Cond       ExpNode
	Claus      []CaseClauNode
	Else       *Loc
	ElseAction *BlockNode
}

type CaseClauNode struct {
	When   Loc
	Ptn    PtnNode
	In     *Loc
	Guard  ExpNode
	Then   Loc
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
	Params *ParamListNode
	In     Loc
	Stats  []StatNode
	Exp    ExpNode
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

func NewToken(loc Loc, text string) Token {
	return Token{loc, text}
}

func NewTokenAntlr(tok antlr.Token) Token {
	fmt.Printf("antlr token %s\n", tok.GetText())
	return Token{Loc: NewLocAntlr(tok), Text: tok.GetText()}
}

func NodeDesc(node Node) string {
	buf := bytes.NewBuffer(nil)
	node.WriteTo(buf)
	return buf.String()
}

func (chunk *ChunkNode) Loc() *Loc {
	return &chunk.loc
}

func (chunk *ChunkNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(chunk ")
	if block := chunk.Block; block != nil {
		block.WriteTo(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (block *BlockNode) Loc() *Loc {
	return &block.loc
}

func (block *BlockNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(block [")
	fmt.Println(block.Stats)
	for _, stat := range block.Stats {
		stat.WriteTo(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("])")
}

func (stat *LetStatNode) Loc() *Loc {
	return &stat.Let
}

func (stat *LetStatNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(let ")
	stat.Ptn.WriteTo(buf)
	buf.WriteString(" ")
	stat.Exp.WriteTo(buf)
	buf.WriteString(")")
}

func (stat *DefStatNode) Loc() *Loc {
	return &stat.Def
}

func (stat *DefStatNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(def \"%s\" [", stat.Name.Text))
	if stat.Params != nil {
		stat.Params.WriteTo(buf)
	}
	buf.WriteString("])")
}

func (stat *ShortDefStatNode) Loc() *Loc {
	return &stat.Def
}

func (stat *ShortDefStatNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(shortdef \"%s\" ", stat.Name.Text))
	stat.Exp.WriteTo(buf)
	buf.WriteString(")")
}

func (params *ParamListNode) Loc() *Loc {
	return &params.Names[0].Loc
}

func (params *ParamListNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(params [")
	for _, param := range params.Names {
		buf.WriteString(fmt.Sprintf("\"%s\" ", param.Text))
	}
	buf.WriteString("])")
}

func (params *ParamListNode) NameStrs() []string {
	nameStrs := make([]string, len(params.Names))
	for _, tok := range params.Names {
		nameStrs = append(nameStrs, tok.Text)
	}
	return nameStrs
}

func (stat *IfStatNode) Loc() *Loc {
	return &stat.Cond[0].If
}

func (stat *IfStatNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(if \"%s\" [")
	for _, cond := range stat.Cond {
		cond.WriteTo(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("] ")
	if else_ := stat.ElseAction; else_ == nil {
		else_.WriteTo(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (cond *IfCondNode) Loc() *Loc {
	return &cond.If
}

func (cond *IfCondNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(ifcond ")
	cond.Cond.WriteTo(buf)
	buf.WriteString(" ")
	cond.Action.WriteTo(buf)
	buf.WriteString(")")
}

func (stat *CaseStatNode) Loc() *Loc {
	return &stat.Case
}

func (stat *CaseStatNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(case ")
	stat.Cond.WriteTo(buf)
	buf.WriteString(" ")
	for _, clau := range stat.Claus {
		clau.WriteTo(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("] ")
	if else_ := stat.ElseAction; else_ != nil {
		else_.WriteTo(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (clau *CaseClauNode) Loc() *Loc {
	return &clau.When
}

func (clau *CaseClauNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(caseclau ")
	clau.Ptn.WriteTo(buf)
	buf.WriteString(" ")
	if clau.Guard != nil {
		clau.Guard.WriteTo(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(" ")
	clau.Action.WriteTo(buf)
	buf.WriteString(")")
}

func (stat *RetStatNode) Loc() *Loc {
	return &stat.Ret
}

func (stat *RetStatNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(return ")
	if exp := stat.Exp; exp != nil {
		exp.WriteTo(buf)
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
}

func (exp *ParenExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *ParenExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(paren ")
	exp.Exp.WriteTo(buf)
	buf.WriteString(")")
}

func (exp *FunCallExpNode) Loc() *Loc {
	return exp.Callable.Loc()
}

func (exp *FunCallExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(funcall ")
	exp.Callable.WriteTo(buf)
	buf.WriteString(" ")
	exp.Args.WriteTo(buf)
	buf.WriteString(")")
}

func (exp *CondOpExpNode) Loc() *Loc {
	return &exp.Colon
}

func (exp *CondOpExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(condop")
	exp.Cond.WriteTo(buf)
	buf.WriteString("  ")
	exp.True.WriteTo(buf)
	buf.WriteString("  ")
	exp.False.WriteTo(buf)
	buf.WriteString(")")
}

func NewVarExpNode(name Token) VarExpNode {
	return VarExpNode{name}
}

func (exp *VarExpNode) Loc() *Loc {
	return &exp.Name.Loc
}

func (exp *VarExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(var \"%s\")", exp.Name.Text))
}

func (exp *UnitExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *UnitExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(unit)")
}

func (exp *BoolExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *BoolExpNode) WriteTo(buf *bytes.Buffer) {
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

func (exp *IntExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(int \"%s\")", exp.Value.Text))
}

func (exp *StrExpNode) Loc() *Loc {
	return &exp.Value.Loc
}

func (exp *StrExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(str \"%s\")", exp.Value.Text))
}

func (exp *ListExpNode) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *ListExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(list ")
	exp.Elts.WriteTo(buf)
	buf.WriteString(")")
}

func (exp *TupleExpNode) Loc() *Loc {
	return &exp.Elts.Open
}

func (exp *TupleExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(tuple ")
	exp.Elts.WriteTo(buf)
	buf.WriteString(")")
}

func (exp *EltListNode) Loc() *Loc {
	return &exp.Open
}

func (elts *EltListNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(eltlist [")
	for _, elt := range elts.Elts {
		elt.WriteTo(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("])")
}

func (exp *SomeExpNode) Loc() *Loc {
	return &exp.SomeLoc
}

func (exp *SomeExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(some ")
	exp.Value.WriteTo(buf)
	buf.WriteString(")")
}

func (exp *NoneExpNode) Loc() *Loc {
	return &exp.loc
}

func (exp *NoneExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(none)")
}

func (exp *AnonFunExpNode) Loc() *Loc {
	return &exp.Open
}

func (exp *AnonFunExpNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(anonfun ")
	exp.Params.WriteTo(buf)
	buf.WriteString(" [")
	fmt.Println(exp.Stats)
	for _, stat := range exp.Stats {
		stat.WriteTo(buf)
		buf.WriteString(" ")
	}
	exp.Exp.WriteTo(buf)
	buf.WriteString("])")
}

func (ptn *UnitPtnNode) Loc() *Loc {
	return &ptn.Open
}

func (ptn *UnitPtnNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(unitptn)")
}

func (ptn *BoolPtnNode) Loc() *Loc {
	return &ptn.loc
}

func (ptn *BoolPtnNode) WriteTo(buf *bytes.Buffer) {
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

func (ptn *IntPtnNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(intptn \"%s\")", ptn.Value.Text))
}

func (ptn *StrPtnNode) Loc() *Loc {
	return &ptn.Value.Loc
}

func (ptn *StrPtnNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("(strptn \"%s\")", ptn.Value.Text))
}

func (ptn *ListPtnNode) Loc() *Loc {
	return &ptn.Elts.Open
}

func (ptn *ListPtnNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(listptn ")
	ptn.Elts.WriteTo(buf)
	buf.WriteString(")")
}

func (ptn *ConsPtnNode) Loc() *Loc {
	return ptn.Left.Loc()
}

func (ptn *ConsPtnNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(consptn ")
	ptn.Left.WriteTo(buf)
	buf.WriteString(" ")
	ptn.Right.WriteTo(buf)
	buf.WriteString(")")
}

func (ptn *TuplePtnNode) Loc() *Loc {
	return &ptn.Elts.Open
}

func (ptn *TuplePtnNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(tupleptn ")
	ptn.Elts.WriteTo(buf)
	buf.WriteString(")")
}

func (elts *EltPtnListNode) Loc() *Loc {
	return &elts.Open
}

func (elts *EltPtnListNode) WriteTo(buf *bytes.Buffer) {
	buf.WriteString("(eltptnlist [")
	for _, elt := range elts.Elts {
		elt.WriteTo(buf)
		buf.WriteString(" ")
	}
	buf.WriteString("])")
}
