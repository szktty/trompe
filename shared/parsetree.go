package trompe

import (
	"fmt"
)

type Node struct {
	Desc interface{}
	Loc  *Loc
}

func newNode(loc *Loc, desc interface{}) *Node {
	return &Node{Loc: loc, Desc: desc}
}

func (n *Node) unionLoc(other *Node) *Loc {
	return n.Loc.Union(other.Loc)
}

func (n *Node) isWildcard() bool {
	if desc, ok := n.Desc.(*IdentNode); ok {
		return desc.Name == "_"
	} else {
		return false
	}
}

func NewNodeList(es ...*Node) []*Node {
	ns := make([]*Node, len(es))
	for i, e := range es {
		ns[i] = e
	}
	return ns
}

func ConsNodeList(head *Node, tail []*Node) []*Node {
	nodes := make([]*Node, len(tail)+1)
	nodes[0] = head
	for i, e := range tail {
		nodes[i+1] = e
	}
	Debugf("nodes = %s, %s, %s", nodes, head, tail)
	return nodes
}

func LocOfNodeList(ns []*Node) *Loc {
	return ns[0].unionLoc(ns[len(ns)-1])
}

type ProgramNode struct {
	Items []*Node
}

type ModuleDefNode struct {
	Name    *Word
	TypeExp *Node
}

type TraitDefNode struct {
	Name    *Word
	TypeExp *Node
}

type ModulePathNode struct {
	Path []*Word
}

func (node *ModulePathNode) NamePath() *NamePath {
	return NamePathOfWordList(node.Path)
}

type ImportNode struct {
	Open  bool
	Path  *Node // ModulePathNode
	Alias *Word
}

type ValueDefNode struct {
	Name    *Node
	TypeExp *Node
}

type TraitNode struct {
}

type LetNode struct {
	Public   bool
	Rec      bool
	Bindings []*Node
	Body     *Node
}

type LetBindingNode struct {
	Rec     bool
	TypeExp *Node
	Ptn     *Node
	Body    *Node
}

type LabelNode struct {
	Name *Word
	Opt  bool
}

type ExtNode struct {
	Name    string
	TypeExp *Node
	Prim    string
}

type UseNode struct {
	Trait  *Word
	Params []*Node
	Dir    int
	Vals   []*Word
}

const (
	TraitDirection = iota
	TraitInclude
	TraitExclude
)

type TraitParam struct {
	Name    *Word
	TypeExp *Node
}

type BlockNode struct {
	Rec    bool
	Name   *Node
	Params []*Node
	Body   *Node
}

type TypeSpecifiedExpNode struct {
	TypeExp *Node
	Exp     *Node
}

type PrefixExpNode struct {
	Prefix *Word
	Exp    *Node
}

type ConstrAppNode struct {
	Constr *Node
	Exp    *Node
}

type ConstrNode struct {
	Path *Node
	Name *Word
}

type ListConsNode struct {
	Head, Tail *Node
}

type IfNode struct {
	Cond  *Node
	True  *Node
	False *Node
}

type WhileNode struct {
	Cond *Node
	Body *Node
}

type ForNode struct {
	Name  *Node
	Init  *Node
	Dir   int
	Limit *Node
	Body  *Node
}

const (
	_ForDirType = iota
	ForDirTo
	ForDirDownTo
)

type CaseNode struct {
	Exp   *Node
	Match []*Node
}

type TryNode struct {
	Exp   *Node
	Match []*Node
}

type FunctionNode struct {
	Match []*Node
}

type FunNode struct {
	MultiMatch *Node
}

type MatchNode struct {
	Ptn  *Node
	Cond *Node
	Body *Node
}

type MultiMatchNode struct {
	Params []*Node
	Cond   *Node
	Body   *Node
}

type WildcardNode struct{}

type PtnIdentNode struct {
	Name string
}

type PtnConstNode struct {
	Value *Node
}

type PtnVarNode struct {
	Ptn  *Node
	Name string
}

type PtnTypeNode struct {
	Ptn     *Node
	TypeExp *Node
}

type SeqPtnNode struct {
	Left, Right *Node
}

type PtnListConsNode struct {
	Head, Tail *Node
}

type PtnListNode struct {
	Elts []*Node
}

type PtnArrayNode struct {
	Elts []*Node
}

type PtnTupleNode struct {
	Comps []*Node
}

type PtnConstrAppNode struct {
	Constr *Node
	Ptn    *Node
}

type PtnSomeNode struct {
	Ptn *Node
}

type LabeledParamNode struct {
	Name *Word
	Ptn  *Node
}

type AppNode struct {
	Exp  *Node
	Args []*Node
}

type LabeledArgNode struct {
	Name *Word
	Exp  *Node
}

type SeqExpNode struct {
	Exps []*Node
}

type IdentNode struct {
	Name string
}

type ValuePathNode struct {
	Path *Node // ModulePathNode
	Name *Node // IdentNode
}

func (node *ValuePathNode) NamePath() *NamePath {
	path := node.Path.Desc.(*ModulePathNode).NamePath()
	path.AddName(node.Name.Desc.(*IdentNode).Name)
	return path
}

type UnitNode struct{}

type BoolNode struct {
	Value bool
}

type IntNode struct {
	Value string
}

type FloatNode struct {
	Value string
}

type CharNode struct {
	Value rune
}

type StringNode struct {
	Value string
}

type RegexpNode struct {
	Value string
}

type ListNode struct {
	Elts []*Node
}

type ArrayNode struct {
	Elts []*Node
}

type ArrayAccessNode struct {
	Array *Node
	Index *Node
	Set   *Node
}

type TupleNode struct {
	Comps []*Node
}

type OptionNode struct {
	Value *Node
}

type AddNode struct {
	Left, Right *Node
}

type SubNode struct {
	Left, Right *Node
}

type MulNode struct {
	Left, Right *Node
}

type DivNode struct {
	Left, Right *Node
}

type ModNode struct {
	Left, Right *Node
}

type FAddNode struct {
	Left, Right *Node
}

type FSubNode struct {
	Left, Right *Node
}

type FMulNode struct {
	Left, Right *Node
}

type FDivNode struct {
	Left, Right *Node
}

type EqNode struct {
	Left, Right *Node
}

type NeNode struct {
	Left, Right *Node
}

type LtNode struct {
	Left, Right *Node
}

type LeNode struct {
	Left, Right *Node
}

type GtNode struct {
	Left, Right *Node
}

type GeNode struct {
	Left, Right *Node
}

type BandNode struct {
	Left, Right *Node
}

type BorNode struct {
	Left, Right *Node
}

type BxorNode struct {
	Left, Right *Node
}

type LshiftNode struct {
	Left, Right *Node
}

type RshiftNode struct {
	Left, Right *Node
}

type NegNode struct {
	Exp *Node
}

type NotNode struct {
	Exp *Node
}

type ConcatNode struct {
	Left, Right *Node
}

type TypePolyNode struct {
	Vars []*Word
	App  *Node
}

type TypeExpNode struct {
	Left, Right *Node
}

type TypeVarNode struct {
	Name *Word
}

type TypeAliasNode struct {
	Exp  *Node
	Name *Word
}

type TypeTupleNode struct {
	Comps []*Node
}

type TypeArrowNode struct {
	Label *Node
	Left  *Node
	Right *Node
}

type TypeConstrAppNode struct {
	Exps   []*Node
	Constr *Node // TypeConstrNode
}

type TypeConstrNode struct {
	Path *Node // ModulePathNode
	Name *Word
}

func (node *TypeConstrNode) NamePath() *NamePath {
	path := node.Path.Desc.(*ModulePathNode).NamePath()
	return path.AddName(node.Name.Value)
}

type TypeParamConstrNode struct {
	Exps   []*Node
	Constr *Node
}

type Word struct {
	Loc   *Loc
	Value string
}

func NewWordList(es ...*Word) []*Word {
	ns := make([]*Word, len(es))
	for i, e := range es {
		ns[i] = e
	}
	return ns
}

func (word *Word) String() string {
	return fmt.Sprintf("{loc=%s, contents=\"%s\"}",
		word.Loc.String(), word.Value)
}

func StringsOfWords(ws []*Word) []string {
	ret := make([]string, len(ws))
	for i, w := range ws {
		ret[i] = w.Value
	}
	return ret
}
