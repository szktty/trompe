package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/szktty/trompe"
)

type NodeBuilder struct {
	*BaseTrompeListener
}

func NewNodeBuilder() *NodeBuilder {
	return new(NodeBuilder)
}

func (l *NodeBuilder) EnterEveryRule(ctx antlr.ParserRuleContext) {
	//fmt.Println(ctx.GetText())
}

func (l *NodeBuilder) EnterChunk(ctx *ChunkContext) {
	fmt.Println("enter chunk")
}

func (l *NodeBuilder) ExitChunk(ctx *ChunkContext) {
	fmt.Println("exit chunk")
}

func (l *NodeBuilder) EnterBlock(ctx *BlockContext) {
	fmt.Println("enter block")
}

func (l *NodeBuilder) ExitBlock(ctx *BlockContext) {
	fmt.Println("exit block")
}

func (l *NodeBuilder) ExitFuncall(ctx *FuncallContext) {
	fmt.Printf("exit funcall: %s\n", ctx.GetText())
}

func ptoken(tok antlr.Token) {
	fmt.Printf("token = line %d, col %d, %d-%d\n", tok.GetLine(), tok.GetColumn(), tok.GetStart(), tok.GetStop())
}

func toLoc(tok antlr.Token) trompe.Loc {
	line := tok.GetLine()
	col := tok.GetColumn()
	offset := tok.GetStart()
	len_ := tok.GetStop() - offset
	start := trompe.Pos{Line: line, Col: col, Offset: offset}
	end := trompe.Pos{Line: line, Col: col, Offset: offset + len_}
	return trompe.Loc{start, end}
}

func (l *NodeBuilder) ExitArgs(ctx *ArgsContext) {
	fmt.Printf("exit args: %s\n", ctx.GetText())
	o := ctx.GetO()
	ptoken(o)
	//loc := toLoc(o)
}

func Parse(file string) {
	input, _ := antlr.NewFileStream(file)
	lexer := NewTrompeLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewTrompeParser(stream)
	p.BuildParseTrees = true
	tree := p.Chunk()
	antlr.ParseTreeWalkerDefault.Walk(NewNodeBuilder(), tree)
}
