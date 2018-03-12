package trompe

type Pos struct {
	Line   int
	Col    int
	Offset int
}

type Loc struct {
	Start Pos
	End   Pos
}
