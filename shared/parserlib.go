package trompe

import (
	"fmt"
)

type ParseError struct {
	Parser *Parser
	Msg    string
}

type Token struct {
	tok int
	lit string
	Loc *Loc
}

type Parser struct {
	s         *Lexer
	recentLit string
	recentPos *Pos
	unit      *Node
	err       *ParseError
}

func (l *Parser) Lex(lval *yySymType) int {
	// TODO: error
	tok, lit, loc := l.s.Scan()
	if tok == EOF {
		return 0
	}
	switch tok {
	case EOF:
		return 0
	case LIDENT, UIDENT, INT, FLOAT, STRING, CHAR, REGEXP, KEYWORD:
		lval.word = &Word{Loc: loc, Value: lit}
	default:
		lval.tok = Token{tok: tok, lit: lit, Loc: loc}
	}
	l.recentLit = lit
	l.recentPos = l.s.position()
	return tok
}

func (l *Parser) Error(e string) {
	Debugf("terms = %s", l.s.terms)
	if len(l.s.terms) > 0 {
		var term string
		last := l.s.terms[len(l.s.terms)-1]
		switch last {
		case "sig", "begin", "if", "function", "fun":
			term = "end"
		case "while", "for", "match", "try":
			term = "done"
		default:
			goto ret
		}

		if term != "" {
			e += fmt.Sprintf(" (Hint: Did you forget keyword `%s' for expression `%s'?)",
				term, last)

		}
	}

ret:
	l.err = &ParseError{Parser: l, Msg: e}
}

func (e *ParseError) Error() string {
	p := e.Parser
	return fmt.Sprintf("File %s, line %d, column %d:\nError: %q %s\n",
		p.s.path, p.recentPos.Line+1, p.recentPos.Col, p.recentLit, e.Msg)
}

func Parse(s *Lexer) (*Node, error) {
	l := Parser{s: s}
	if yyParse(&l) != 0 {
		if l.err != nil {
			return nil, l.err
		} else {
			panic("Parse error")
		}
	}
	if l.unit == nil {
		panic("Compiled unit is none")
	}
	return l.unit, nil
}

func WrapWithPolyNode(node *Node) *Node {
	vars := checkTyvars(make([]*Word, 0), node)
	if len(vars) == 0 {
		return node
	} else {
		return newNode(node.Loc, &TypePolyNode{Vars: vars, App: node})
	}
}

func checkTyvars(vars []*Word, node *Node) []*Word {
	switch desc := node.Desc.(type) {
	case *TypeAliasNode:
		return checkTyvars(vars, desc.Exp)
	case *TypeVarNode:
		return addTyvar(vars, desc.Name)
	case *TypeArrowNode:
		vars = checkTyvars(vars, desc.Left)
		return checkTyvars(vars, desc.Right)
	case *TypeConstrAppNode:
		for _, e := range desc.Exps {
			vars = checkTyvars(vars, e)
		}
		return vars
	case *TypeTupleNode:
		for _, e := range desc.Comps {
			vars = checkTyvars(vars, e)
		}
		return vars
	default:
		return vars
	}
}

func addTyvar(ws []*Word, w *Word) []*Word {
	for _, e := range ws {
		if e.Value == w.Value {
			return ws
		}
	}
	return append(ws, w)
}
