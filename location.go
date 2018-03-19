package trompe

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Pos struct {
	Line   int
	Col    int
	Offset int
}

type Loc struct {
	Start Pos
	End   Pos
}

func TokenLoc(tok antlr.Token) Loc {
	line := tok.GetLine()
	col := tok.GetColumn()
	offset := tok.GetStart()
	len_ := tok.GetStop() - offset
	start := Pos{Line: line, Col: col, Offset: offset}
	end := Pos{Line: line, Col: col, Offset: offset + len_}
	return Loc{start, end}
}
