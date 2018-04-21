package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	. "github.com/szktty/trompe"
)

type ChunkListener struct {
	*BaseTrompeListener
	Node ChunkNode
}

func NewChunkListener() *ChunkListener {
	return new(ChunkListener)
}

func (l *ChunkListener) EnterChunk(ctx *ChunkContext) {
	fmt.Println("enter chunk")
	blockCtx := ctx.Block()
	if blockCtx != nil {
		block := NewBlockListener()
		blockCtx.EnterRule(block)
		l.Node = ChunkNode{Block: &block.Node}
	}
}

type BlockListener struct {
	*BaseTrompeListener
	Node BlockNode
}

func NewBlockListener() *BlockListener {
	return new(BlockListener)
}

func (l *BlockListener) EnterBlock(ctx *BlockContext) {
	fmt.Println("enter block")
	var stats []Node
	for _, statCtx := range ctx.AllStat() {
		stat := NewStatListener()
		statCtx.EnterRule(stat)
		if stat.Node == nil {
			panic(fmt.Sprintf("stat not found: %s", statCtx.GetText()))
		}
		stats = append(stats, stat.Node)
	}
	l.Node = BlockNode{Stats: stats}
}

type StatListener struct {
	*BaseTrompeListener
	Node StatNode
}

func NewStatListener() *StatListener {
	return new(StatListener)
}

func (l *StatListener) EnterStat(ctx *StatContext) {
	fmt.Printf("enter stat\n")
	if funcallCtx := ctx.Funcall(); funcallCtx != nil {
		funcall := NewFuncallListener()
		funcallCtx.EnterRule(funcall)
		l.Node = &funcall.Node
	} else if forCtx := ctx.For_(); forCtx != nil {
		for_ := NewForStatListener()
		forCtx.EnterRule(for_)
		l.Node = &for_.Node
	} else {
		panic("not impl")
	}
}

type ForStatListener struct {
	*BaseTrompeListener
	Node ForStatNode
}

func NewForStatListener() *ForStatListener {
	return new(ForStatListener)
}

func (l *ForStatListener) EnterFor_(ctx *For_Context) {
	fmt.Printf("enter for\n")
	ptn := NewPatternListener()
	ctx.Pattern().EnterRule(ptn)

	exp := NewExpListener()
	ctx.Exp().EnterRule(exp)

	block := NewBlockListener()
	ctx.Block().EnterRule(block)

	// TODO: pattern
	l.Node = ForStatNode{Ptn: ptn.Node, Exp: exp.Node, Block: block.Node}
}

type FuncallListener struct {
	*BaseTrompeListener
	Node FunCallExpNode
}

func NewFuncallListener() *FuncallListener {
	return new(FuncallListener)
}

func (l *FuncallListener) EnterFuncall(ctx *FuncallContext) {
	fmt.Printf("enter funcall\n")
	exp := NewSimpleExpListener()
	ctx.Simpleexp().EnterRule(exp)

	// TODO: '(', ')'
	fmt.Printf("arglist\n")
	args := NewArglistListener()
	if argsCtx := ctx.Arglist(); argsCtx != nil {
		argsCtx.EnterRule(args)
	}

	l.Node = FunCallExpNode{Callable: exp.Node, Args: args.Node}
}

type ArglistListener struct {
	*BaseTrompeListener
	Node EltListNode
}

func NewArglistListener() *ArglistListener {
	return new(ArglistListener)
}

func (l *ArglistListener) EnterArglist(ctx *ArglistContext) {
	// TODO: '(', ')'
	if expsCtx := ctx.Explist(); expsCtx != nil {
		exps := NewExplistListener()
		expsCtx.EnterRule(exps)
		l.Node = exps.Node
	}
}

type ExplistListener struct {
	*BaseTrompeListener
	Node EltListNode
}

func NewExplistListener() *ExplistListener {
	return new(ExplistListener)
}

func (l *ExplistListener) EnterExplist(ctx *ExplistContext) {
	fmt.Printf("enter explist\n")
	var exps []Node
	for _, expCtx := range ctx.AllExp() {
		exp := NewExpListener()
		expCtx.EnterRule(exp)
		exps = append(exps, exp.Node)
	}
	l.Node = EltListNode{Elts: exps}
}

type ExpListener struct {
	*BaseTrompeListener
	Node ExpNode
}

func NewExpListener() *ExpListener {
	return new(ExpListener)
}

func (l *ExpListener) EnterExp(ctx *ExpContext) {
	// TODO
	fmt.Printf("enter exp: %s\n", ctx.GetText())
	if expCtx := ctx.Simpleexp(); expCtx != nil {
		exp := NewSimpleExpListener()
		expCtx.EnterRule(exp)
		l.Node = exp.Node
	} else if opCtx := ctx.Rangeop(); opCtx != nil {
		fmt.Printf("enter range\n")
		op := NewRangeOpListener()
		opCtx.EnterRule(op)
		left := NewExpListener()
		leftCtx := ctx.GetLeft()
		leftCtx.EnterRule(left)
		rightCtx := ctx.GetRight()
		right := NewExpListener()
		rightCtx.EnterRule(right)
		l.Node = &RangeExpNode{Left: left.Node,
			Op:    op.Token,
			Close: op.Close,
			Right: right.Node}
	} else {
		panic("not impl")
	}
}

type RangeOpListener struct {
	*BaseTrompeListener
	Token Token
	Close bool
}

func NewRangeOpListener() *RangeOpListener {
	return new(RangeOpListener)
}

func (l *RangeOpListener) EnterRangeop(ctx *RangeopContext) {
	switch ctx.GetText() {
	case "...":
		l.Close = true
	case "..<":
		l.Close = false
	default:
		panic("invalid token")
	}
	l.Token = NewTokenAntlr(ctx.GetStart())
}

type ParenexpListener struct {
	*BaseTrompeListener
	Node ExpNode
}

func NewParenexpListener() *ParenexpListener {
	return new(ParenexpListener)
}

func (l *ParenexpListener) EnterParenexp(ctx *ParenexpContext) {
	// TODO: open, close
	exp := NewExpListener()
	ctx.Exp().EnterRule(exp)
	l.Node = exp.Node
}

type SimpleExpListener struct {
	*BaseTrompeListener
	Node ExpNode
}

func NewSimpleExpListener() *SimpleExpListener {
	return new(SimpleExpListener)
}

func (l *SimpleExpListener) EnterSimpleexp(ctx *SimpleexpContext) {
	// TODO
	fmt.Printf("enter simpleexp: %s\n", ctx.GetText())
	fmt.Printf("int %s\n", ctx.Int_())
	fmt.Printf("int %s\n", ctx.Hexint())
	fmt.Printf("float %s\n", ctx.Float_())
	fmt.Printf("float %s\n", ctx.Hexfloat())

	if parenCtx := ctx.Parenexp(); parenCtx != nil {
		exp := NewParenexpListener()
		ctx.Parenexp().EnterRule(exp)
		l.Node = exp.Node
	} else if varCtx := ctx.Var_(); varCtx != nil {
		var_ := NewVarExpListener()
		ctx.Var_().EnterRule(var_)
		l.Node = &var_.Node
	} else if intCtx := ctx.Int_(); intCtx != nil {
		int_ := NewIntListener()
		intCtx.EnterRule(int_)
		l.Node = &int_.Node
	} else if strCtx := ctx.String_(); strCtx != nil {
		str := NewStringListener()
		strCtx.EnterRule(str)
		l.Node = &str.Node
	} else {
		panic("not impl")
	}
}

type VarExpListener struct {
	*BaseTrompeListener
	Node VarExpNode
}

func NewVarExpListener() *VarExpListener {
	return new(VarExpListener)
}

func (l *VarExpListener) EnterVar_(ctx *Var_Context) {
	// TODO
	fmt.Printf("enter var\n")
	l.Node = NewVarExpNode(NewTokenAntlr(ctx.GetStart()))
}

type IntListener struct {
	*BaseTrompeListener
	Node IntExpNode
}

func NewIntListener() *IntListener {
	return new(IntListener)
}

func (l *IntListener) EnterInt_(ctx *Int_Context) {
	fmt.Printf("enter int\n")
	l.Node = IntExpNode{Value: NewTokenAntlr(ctx.GetStart())}
}

type StringListener struct {
	*BaseTrompeListener
	Node StrExpNode
}

func NewStringListener() *StringListener {
	return new(StringListener)
}

func (l *StringListener) EnterString_(ctx *String_Context) {
	fmt.Printf("enter string\n")
	l.Node = StrExpNode{Value: NewTokenAntlr(ctx.GetStart())}
}

type PatternListener struct {
	*BaseTrompeListener
	Node PtnNode
}

func NewPatternListener() *PatternListener {
	return new(PatternListener)
}

func (l *PatternListener) EnterPattern(ctx *PatternContext) {
	fmt.Printf("enter pattern: %s\n", ctx.GetText())

	if varCtx := ctx.NAME(); varCtx != nil {
		l.Node = &VarPtnNode{NewTokenAntlr(ctx.GetStart())}
	} else if intCtx := ctx.Int_(); intCtx != nil {
		int_ := NewIntListener()
		intCtx.EnterRule(int_)
		l.Node = &int_.Node
	} else if strCtx := ctx.String_(); strCtx != nil {
		str := NewStringListener()
		strCtx.EnterRule(str)
		l.Node = &str.Node
	} else {
		panic("not impl")
	}
}

func Parse(file string) Node {
	input, _ := antlr.NewFileStream(file)
	lexer := NewTrompeLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewTrompeParser(stream)
	p.BuildParseTrees = true
	tree := p.Chunk()
	listener := NewChunkListener()
	antlr.ParseTreeWalkerDefault.EnterRule(listener, tree)
	fmt.Printf("%s\n", NodeDesc(&listener.Node))
	return &listener.Node
}
