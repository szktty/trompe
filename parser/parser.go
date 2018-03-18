package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	. "github.com/szktty/trompe"
)

func ToLoc(tok antlr.Token) Loc {
	line := tok.GetLine()
	col := tok.GetColumn()
	offset := tok.GetStart()
	len_ := tok.GetStop() - offset
	start := Pos{Line: line, Col: col, Offset: offset}
	end := Pos{Line: line, Col: col, Offset: offset + len_}
	return Loc{start, end}
}

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
	Args []ExpNode
}

func NewFuncallListener() *FuncallListener {
	return new(FuncallListener)
}

func (l *FuncallListener) EnterFuncall(ctx *FuncallContext) {
	// TODO: '(', ')'
	if argsCtx := ctx.Arglist(); argsCtx != nil {
		args := NewArglistListener()
		argsCtx.EnterRule(args)
	}
}

type ArglistListener struct {
	*BaseTrompeListener
	Node ArgListNode
}

func NewArglistListener() *ArglistListener {
	return new(ArglistListener)
}

func (l *ArglistListener) EnterArglist(ctx *ArglistContext) {
	// TODO: '(', ')'
	if expsCtx := ctx.Explist(); expsCtx != nil {
		exps := NewExplistListener()
		expsCtx.EnterRule(exps)
	}
}

type ExplistListener struct {
	*BaseTrompeListener
	Exps []ExpNode
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
	// TODO
	//l.Exps = Node{Exps: exps}
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
	fmt.Printf("enter exp\n")
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
