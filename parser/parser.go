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
		l.Node = ChunkNode{Block: block.Node}
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
	statLen := len(ctx.AllStat()) + 1
	stats := make([]Node, statLen)
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
	}
}

type FuncallListener struct {
	*BaseTrompeListener
	Exp FunCallExpNode
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

	l.Exp = FunCallExpNode{Callable: callable.Exp, Args: args.Exps}
}

type CallableListener struct {
	*BaseTrompeListener
	Exp ExpNode
}

func NewCallableListener() *CallableListener {
	return new(CallableListener)
}

func (l *CallableListener) EnterCallable(ctx *CallableContext) {
	fmt.Printf("enter callable\n")
	if varCtx := ctx.Var_(); varCtx != nil {
		var_ := NewVarExpListener()
		ctx.Var_().EnterRule(var_)
		l.Exp = &var_.Exp
	} else {
		exp := NewParenexpListener()
		ctx.Parenexp().EnterRule(exp)
		l.Exp = exp.Exp
	}
}

type ArglistListener struct {
	*BaseTrompeListener
	Exps EltListNode
}

func NewArglistListener() *ArglistListener {
	return new(ArglistListener)
}

func (l *ArglistListener) EnterArglist(ctx *ArglistContext) {
	// TODO: '(', ')'
	if expsCtx := ctx.Explist(); expsCtx != nil {
		exps := NewExplistListener()
		expsCtx.EnterRule(exps)
		l.Exps = exps.Exps
	}
}

type ExplistListener struct {
	*BaseTrompeListener
	Exps EltListNode
}

func NewExplistListener() *ExplistListener {
	return new(ExplistListener)
}

func (l *ExplistListener) EnterExplist(ctx *ExplistContext) {
	expLen := len(ctx.AllExp())
	exps := make([]Node, expLen)
	for _, expCtx := range ctx.AllExp() {
		exp := NewExpListener()
		expCtx.EnterRule(exp)
		exps = append(exps, exp.Exp)
	}
	l.Exps = EltListNode{Elts: exps}
}

type ExpListener struct {
	*BaseTrompeListener
	Exp ExpNode
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
		l.Exp = &str.String
	}
}

type ParenexpListener struct {
	*BaseTrompeListener
	Exp ExpNode
}

func NewParenexpListener() *ParenexpListener {
	return new(ParenexpListener)
}

func (l *ParenexpListener) EnterParenexp(ctx *ParenexpContext) {
	// TODO: open, close
	exp := NewExpListener()
	ctx.Exp().EnterRule(exp)
	l.Exp = exp.Exp
}

type VarExpListener struct {
	*BaseTrompeListener
	Exp VarExpNode
}

func NewVarExpListener() *VarExpListener {
	return new(VarExpListener)
}

func (l *VarExpListener) EnterVar_(ctx *Var_Context) {
	// TODO
	fmt.Printf("enter var\n")
}

type StringListener struct {
	*BaseTrompeListener
	String StrExpNode
}

func NewStringListener() *StringListener {
	return new(StringListener)
}

func (l *StringListener) EnterString(ctx *String_Context) {
	fmt.Printf("enter string\n")
	l.String = StrExpNode{Value: ctx.GetText()}
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
	return &listener.Node
}
