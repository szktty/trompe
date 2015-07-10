package trompe

type TypedNode struct {
	Desc interface{}
	Loc  *Loc
	Type Type
}

type TypedProgramNode struct {
	Items []*TypedNode
}

type TypedLetNode struct {
	Public   bool
	Rec      bool
	Bindings []*TypedNode
	Body     *TypedNode
}

type TypedLetBindingNode struct {
	Ptn, Body *TypedNode
}

type TypedBlockNode struct {
	Name   *TypedNode
	Params []*TypedNode
	Body   *TypedNode
}

type TypedConstrAppNode struct {
	Constr *TypedNode
	Exp    *TypedNode
}

type TypedConstrNode struct {
	Path *NamePath
}

type TypedIfNode struct {
	Cond  *TypedNode
	True  *TypedNode
	False *TypedNode
}

type TypedWhileNode struct {
	Cond *TypedNode
	Body *TypedNode
}

type TypedForNode struct {
	Name  *TypedNode
	Init  *TypedNode
	Dir   int
	Limit *TypedNode
	Body  *TypedNode
}

type TypedCaseNode struct {
	Exp   *TypedNode
	Match []*TypedNode
}

type TypedTryNode struct {
	Exp   *TypedNode
	Match []*TypedNode
}

type TypedFunctionNode struct {
	Match []*TypedNode
}

type TypedMatchNode struct {
	Ptn  *TypedNode
	Cond *TypedNode
	Body *TypedNode
}

type TypedFunNode struct {
	MultiMatch *TypedNode
}

type TypedMultiMatchNode struct {
	Params []*TypedNode
	Cond   *TypedNode
	Body   *TypedNode
}

type TypedWildcardNode struct{}

type TypedPtnIdentNode struct {
	Name string
}

type TypedPtnConstNode struct {
	Value *TypedNode
}

type TypedPtnVarNode struct {
	Ptn  *TypedNode
	Name string
}

type TypedPtnTypeNode struct {
	Ptn     *TypedNode
	TypeExp *TypedNode
}

type TypedSeqPtnNode struct {
	Left, Right *TypedNode
}

type TypedPtnTupleNode struct {
	Comps []*TypedNode
}

type TypedPtnListNode struct {
	Elts []*TypedNode
}

type TypedPtnArrayNode struct {
	Elts []*TypedNode
}

type TypedPtnListConsNode struct {
	Head, Tail *TypedNode
}

type TypedLabeledParamNode struct {
	Name string
	Ptn  *TypedNode
}

type TypedAppNode struct {
	Exp  *TypedNode
	Args []*TypedNode
}

type TypedSeqExpNode struct {
	Exps []*TypedNode
}

type TypedValuePathNode struct {
	Path *NamePath
}

type TypedIdentNode struct {
	Name string
}

type TypedLabeledArgNode struct {
	Name string
	Exp  *TypedNode
}

type TypedUnitNode struct{}

type TypedBoolNode struct {
	Value bool
}

type TypedIntNode struct {
	Value    string
	Decimal  bool
	SmallInt int64
}

type TypedFloatNode struct {
	Value string
}

type TypedCharNode struct {
	Value rune
}

type TypedStringNode struct {
	Value string
}

type TypedFormatNode struct {
	Format *Format
}

type TypedListNode struct {
	Elts []*TypedNode
}

type TypedListConsNode struct {
	Head, Tail *TypedNode
}

type TypedArrayNode struct {
	Elts []*TypedNode
}

type TypedArrayAccessNode struct {
	Array *TypedNode
	Index *TypedNode
	Set   *TypedNode
}

type TypedTupleNode struct {
	Comps []*TypedNode
}

type TypedOptionNode struct {
	Value *TypedNode
}

type TypedAddNode struct {
	Left, Right *TypedNode
}

type TypedSubNode struct {
	Left, Right *TypedNode
}

type TypedMulNode struct {
	Left, Right *TypedNode
}

type TypedDivNode struct {
	Left, Right *TypedNode
}

type TypedModNode struct {
	Left, Right *TypedNode
}

type TypedFAddNode struct {
	Left, Right *TypedNode
}

type TypedFSubNode struct {
	Left, Right *TypedNode
}

type TypedFMulNode struct {
	Left, Right *TypedNode
}

type TypedFDivNode struct {
	Left, Right *TypedNode
}

type TypedEqNode struct {
	Left, Right *TypedNode
}

type TypedNeNode struct {
	Left, Right *TypedNode
}

type TypedLtNode struct {
	Left, Right *TypedNode
}

type TypedLeNode struct {
	Left, Right *TypedNode
}

type TypedGtNode struct {
	Left, Right *TypedNode
}

type TypedGeNode struct {
	Left, Right *TypedNode
}

type TypedBandNode struct {
	Left, Right *Node
}

type TypedBorNode struct {
	Left, Right *Node
}

type TypedBxorNode struct {
	Left, Right *Node
}

type TypedLshiftNode struct {
	Left, Right *Node
}

type TypedRshiftNode struct {
	Left, Right *Node
}

type TypedNegNode struct {
	Exp *Node
}

type TypedNotNode struct {
	Exp *Node
}

type TypedConcatNode struct {
	Left, Right *Node
}

func IsConstTypedNode(node *TypedNode) bool {
	switch node.Desc.(type) {
	case *TypedUnitNode, *TypedBoolNode, *TypedIntNode, *TypedFloatNode,
		*TypedCharNode, *TypedStringNode:
		return true
	default:
		return false
	}
}

func IsConstAllTypedNodes(nodes []*TypedNode) bool {
	for _, node := range nodes {
		if !IsConstTypedNode(node) {
			return false
		}
	}
	return true
}

func NewTypedNode(loc *Loc, ty Type, desc interface{}) *TypedNode {
	return &TypedNode{Desc: desc, Loc: loc, Type: ty}
}

func (node *TypedNode) Name() (string, bool) {
	switch desc := node.Desc.(type) {
	case *TypedIdentNode:
		return desc.Name, true
	default:
		return "", false
	}
}

func (node *TypedNode) NameExn() string {
	if name, ok := node.Name(); ok {
		return name
	} else {
		panic("node does not have Name member")
	}
}

func TypesOfTypedNodes(nodes []*TypedNode) []Type {
	tys := make([]Type, len(nodes))
	for i, node := range nodes {
		tys[i] = node.Type
	}
	return tys
}
