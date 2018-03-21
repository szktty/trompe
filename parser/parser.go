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
	if funcallCtx := ctx.Funcall(); funcallCtx != nil {
		funcall := NewFuncallListener()
		funcallCtx.EnterRule(funcall)
		l.Node = &funcall.Node
	}
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
	callable := NewCallableListener()
	ctx.Callable().EnterRule(callable)

	// TODO: '(', ')'
	args := NewArglistListener()
	if argsCtx := ctx.Arglist(); argsCtx != nil {
		argsCtx.EnterRule(args)
	}

	l.Node = FunCallExpNode{Callable: callable.Node, Args: args.Node}
}

type CallableListener struct {
	*BaseTrompeListener
	Node ExpNode
}

func NewCallableListener() *CallableListener {
	return new(CallableListener)
}

func (l *CallableListener) EnterCallable(ctx *CallableContext) {
	fmt.Printf("enter callable\n")
	if varCtx := ctx.Var_(); varCtx != nil {
		var_ := NewVarExpListener()
		ctx.Var_().EnterRule(var_)
		l.Node = &var_.Node
	} else {
		exp := NewParenexpListener()
		ctx.Parenexp().EnterRule(exp)
		l.Node = exp.Node
	}
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

	if strCtx := ctx.String_(); strCtx != nil {
		str := NewStringListener()
		strCtx.EnterRule(str)
		l.Node = &str.Node
	}
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
