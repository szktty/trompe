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

func (l *ChunkListener) ExitChunk(ctx *ChunkContext) {
	fmt.Println("exit chunk")
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

func (l *BlockListener) ExitBlock(ctx *BlockContext) {
	fmt.Println("exit block")
}

type StatListener struct {
	*BaseTrompeListener
	Node StatNode
}

func NewStatListener() *StatListener {
	return new(StatListener)
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
