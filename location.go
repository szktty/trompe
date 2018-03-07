package trompe

type Position struct {
	Line   int
	Col    int
	Offset int
}

type Location struct {
	Start Position
	End   Position
}
