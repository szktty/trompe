package trompe

type Token interface {
	Loc() *Loc
}

type StrTok struct {
	loc   Loc
	Value string
}

func (tok *StrTok) Loc() *Loc {
	return &tok.loc
}
