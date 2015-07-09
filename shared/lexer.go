package trompe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

var keywords = map[string]int{
	"_":          WILDCARD,
	"abstract":   ABSTRACT,
	"and":        AND,
	"as":         AS,
	"assert":     ASSERT,
	"begin":      BEGIN,
	"constraint": CONSTRAINT,
	"do":         DO,
	"done":       DONE,
	"downto":     DOWNTO,
	"else":       ELSE,
	"end":        END,
	"exception":  EXCEPTION,
	"external":   EXTERNAL,
	"false":      FALSE,
	"for":        FOR,
	"fun":        FUN,
	"function":   FUNCTION,
	"goto":       GOTO,
	"if":         IF,
	"import":     IMPORT,
	"in":         IN,
	"let":        LET,
	"land":       BAND,
	"lor":        BOR,
	"lxor":       BXOR,
	"lsr":        RSHIFT,
	"lsl":        LSHIFT,
	"match":      MATCH,
	"mod":        MOD,
	"module":     MODULE,
	"mutable":    MUTABLE,
	"not":        NOT,
	"of":         OF,
	"open":       OPEN,
	"or":         OR,
	"override":   OVERRIDE,
	"partial":    PARTIAL,
	"rec":        REC,
	"sig":        SIG,
	"struct":     STRUCT,
	"then":       THEN,
	"to":         TO,
	"trait":      TRAIT,
	"true":       TRUE,
	"try":        TRY,
	"type":       TYPE,
	"use":        USE,
	"val":        VAL,
	"when":       WHEN,
	"while":      WHILE,
	"with":       WITH,
	"without":    WITHOUT,
	"Some":       SOME,
	"None":       NONE,
	"+":          ADD,
	"+.":         ADD_DOT,
	"-":          SUB,
	"-.":         SUB_DOT,
	"*":          MUL,
	"*.":         MUL_DOT,
	"/":          DIV,
	"/.":         DIV_DOT,
	"%":          REM,
	"**":         POW,
	"^":          CONCAT,
	"<>":         NE,
	"<=":         LE,
	">=":         GE,
	"<":          LT,
	">":          GT,
	"=":          EQ,
	"(":          LPAREN,
	")":          RPAREN,
	"{":          LBRACE,
	"}":          RBRACE,
	"[":          LBRACK,
	"]":          RBRACK,
	";":          SEMI,
	";;":         SEMI2,
	":":          COLON,
	"::":         COLON2,
	",":          COMMA,
	".":          DOT,
	".(":         DOT_LPAREN,
	"<-":         LARROW,
	"->":         RARROW,
	"|":          PIPE,
	"$":          DOL,
	"&":          AMP,
	"?":          Q,
	"!":          EP,
	"~":          TILDA,
	"'":          QUOTE,
}

type Lexer struct {
	path              string
	src               []rune
	line, col, offset int
	start, end        Pos
	terms             []string // terms of open and close
	nextIdent         bool
	indent            *Indent
	numDedents        int
}

type Indent struct {
	Size int
	Prev *Indent
}

func NewLexerFromFile(path string) *Lexer {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		// TODO: error
	}
	l := new(Lexer)
	l.src = bytes.Runes(b)
	l.path = path
	l.terms = make([]string, 0)
	l.indent = &Indent{}
	return l
}

func NewLexer(path string, src string) *Lexer {
	l := new(Lexer)
	reader := strings.NewReader(src)
	for r, _, e := reader.ReadRune(); e == nil; r, _, e = reader.ReadRune() {
		l.src = append(l.src, r)
	}
	l.path = path
	l.terms = make([]string, 0)
	return l
}

func (l *Lexer) addTerm(term string) {
	switch term {
	case "end":
		if len(l.terms) > 0 {
			switch l.terms[len(l.terms)-1] {
			case "sig", "begin", "if", "function", "fun":
				l.terms = l.terms[0 : len(l.terms)-1]
			default:
				// ignore (delegate to parser)
			}
		}
		// ignore (delegate to parser)
	case "done":
		if len(l.terms) > 0 {
			switch l.terms[len(l.terms)-1] {
			case "while", "for", "match", "try":
				l.terms = l.terms[0 : len(l.terms)-1]
			default:
				// ignore (delegate to parser)
			}
			// ignore (delegate to parser)
		}
	default:
		l.terms = append(l.terms, term)
	}
}

func (s *Lexer) Scan() (tok int, lit string, loc *Loc) {
dedent:
	if s.numDedents > 0 {
		s.numDedents--
		s.indent = s.indent.Prev
		tok = DEDENT
		pos := s.position()
		loc = NewLoc(pos, pos)
		loc.File = s.path
		return
	}

indent:
	if s.nextIdent {
		s.skipNewlines()
		size := s.scanIndent()
		if s.indent.Size < size {
			s.nextIdent = false
			s.indent = &Indent{Size: size, Prev: s.indent}
			tok = INDENT
			return
		} else if s.indent.Size == size {
			s.nextIdent = false
		} else {
			n := 0
			in := s.indent
			for in != nil {
				if in.Size == size {
					s.numDedents = n
					s.nextIdent = false
					goto dedent
				}
				in = in.Prev
				n++
			}
			panic(fmt.Errorf("Line %d: invalid indent size %d",
				s.position().Line+1, size))
		}
	}

start:
	s.skipWhitespace()
	var start = s.position()
	ch := s.peek()
	switch ch {
	case -1:
		lit = "<EOF>"
		tok = EOF
	case '\n':
		s.next()
		s.nextIdent = true
		goto indent
	case ' ':
		s.skipWhitespace()
		goto start
	case '#':
		s.skipSingleLineComment()
		goto start
	case '(':
		s.next()
		if ok := s.skipComment(); ok {
			goto start
		} else {
			lit = "("
			tok = LPAREN
		}
	case '.':
		s.next()
		if s.peek() == '(' {
			s.next()
			lit = ".("
			tok = DOT_LPAREN
		} else {
			lit = "."
			tok = DOT
		}
	case ')', '[', ']', '{', '}', ',', '?', '\'':
		lit = string(ch)
		tok = keywords[lit]
		s.next()
	case '+', '-', '*', '/', '%', '=', '<', '>', '~', '^', '&', '|', '$':
		lit = s.scanOperator()
		tok = keywords[lit]
	case ':':
		s.next()
		ch := s.peek()
		switch {
		case ch == ':':
			s.next()
			tok = COLON2
		case isLetter(ch):
			lit = s.scanIdent()
			tok = LABELL
		default:
			tok = COLON
		}
	case ';':
		s.next()
		if s.peek() == ';' {
			s.next()
			tok = SEMI2
			lit = ";;"
		} else {
			tok = SEMI
			lit = ";"
		}
	case '"':
		s.next()
		tok = STRING
		lit = s.scanString(ch)
	case '_':
		s.next()
		tok = WILDCARD
		lit = "_"
	default:
		switch {
		case isLetter(ch):
			lit = s.scanIdent()
			if keyword, ok := keywords[lit]; ok {
				tok = keyword
				switch tok {
				case SIG, BEGIN, IF, WHILE, FOR, MATCH, TRY,
					FUNCTION, FUN, END, DONE:
					s.addTerm(lit)
				}
			} else if lit == strings.ToLower(lit) {
				if s.peek() == ':' {
					if s.npeek(2) == ':' {
						tok = LIDENT
					} else {
						s.next()
						tok = LABELR
					}
				} else {
					tok = LIDENT
				}
			} else {
				tok = UIDENT
			}
		case isDigit(ch):
			if v, ok := s.scanInt(); ok {
				tok = INT
				lit = v
			} else if v, ok := s.scanFloat(); ok {
				tok = FLOAT
				lit = v
			}
		default:
			panic(fmt.Errorf("invalid character '%c'", ch))
		}
	}
	loc = NewLoc(start, s.position())
	loc.File = s.path
	return
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' ||
		'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

func isWhitespace(ch rune) bool {
	return ch == ' '
}

func isNewline(ch rune) bool {
	return ch == '\r' || ch == '\n'
}

func (s *Lexer) peek() rune {
	if !s.isAtEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *Lexer) npeek(n int) rune {
	if s.offset+n+1 < len(s.src) {
		return s.src[s.offset+n-1]
	} else {
		return -1
	}
}

func (s *Lexer) next() {
	if !s.isAtEOF() {
		if s.peek() == '\n' {
			s.col = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *Lexer) isAtEOF() bool {
	return len(s.src) <= s.offset
}

func (s *Lexer) position() *Pos {
	return &Pos{Line: s.line, Col: s.offset - s.col, Offset: s.offset}
}

func (s *Lexer) scanIndent() int {
	size := 0
	for {
		switch s.peek() {
		case -1:
			return size
		case '\n':
			s.next()
			size = 0
		case '#':
			s.skipSingleLineComment()
			size = 0
		case ' ':
			s.next()
			size++
		default:
			return size
		}
	}
}

func (s *Lexer) skipWhitespace() {
	for isWhitespace(s.peek()) {
		s.next()
	}
}

func (s *Lexer) skipNewlines() {
	for isNewline(s.peek()) {
		s.next()
	}
}

func (s *Lexer) skipUpToNewline() {
	for !isNewline(s.peek()) {
		s.next()
	}
}

func (s *Lexer) skipSingleLineComment() {
	s.skipUpToNewline()
}

func (s *Lexer) skipComment() bool {
	if s.peek() != '*' {
		return false
	}
	s.next()
	flag := false
	for {
		switch s.peek() {
		case -1:
			return false
		case ')':
			if flag {
				s.next()
				return true
			} else {
				flag = false
			}
		case '*':
			flag = !flag
			s.next()
		default:
			s.next()
		}
	}
}

func (s *Lexer) scanIdent() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}

func (s *Lexer) scanInt() (string, bool) {
	var ret []rune
	if s.peek() == '0' {
		ret = append(ret, s.peek())
		s.next()
		switch s.peek() {
		case 'x', 'X':
			ret = append(ret, s.peek())
			s.next()
			for isHexDigit(s.peek()) {
				ret = append(ret, s.peek())
				s.next()
			}
			if s.peek() == '.' {
				return "", false
			}
		}
	} else {
		for isDigit(s.peek()) {
			ret = append(ret, s.peek())
			s.next()
		}
		if s.peek() == '.' {
			return "", false
		}
	}

	return string(ret), true
}

func (s *Lexer) scanFloat() (string, bool) {
	var ret []rune
	if s.peek() == '0' {
		ret = append(ret, s.peek())
		s.next()
		switch s.peek() {
		case 'x', 'X':
			ret = append(ret, s.peek())
			s.next()
			for isHexDigit(s.peek()) {
				ret = append(ret, s.peek())
				s.next()
			}
			if s.peek() == '.' {
				ret = append(ret, s.peek())
				s.next()
				for isHexDigit(s.peek()) {
					ret = append(ret, s.peek())
					s.next()
				}
			}
		}
	} else {
		for isDigit(s.peek()) {
			ret = append(ret, s.peek())
			s.next()
		}
		if s.peek() == '.' {
			ret = append(ret, s.peek())
			s.next()
			for isDigit(s.peek()) {
				ret = append(ret, s.peek())
				s.next()
			}
		}
	}

	switch s.peek() {
	case 'e', 'E', 'p', 'P':
		ret = append(ret, s.peek())
		s.next()
		switch s.peek() {
		case '-', '+':
			ret = append(ret, s.peek())
			s.next()
			for isDigit(s.peek()) {
				ret = append(ret, s.peek())
				s.next()
			}
		default:
			return "", false
		}
	}
	return string(ret), true
}

func (s *Lexer) scanChar() (rune, bool) {
	c := s.peek()
	s.next()
	if c == '\\' {
		c := s.peek()
		s.next()
		switch c {
		case '\\', '\'', '"':
			return c, true
		case 't':
			return '\t', true
		case 'r':
			return '\r', true
		case 'n':
			return '\n', true
		case '0':
			return '\x00', true
		default:
			return c, true
		}
	} else {
		return c, false
	}
}

func (s *Lexer) scanString(p rune) string {
	var ret = bytes.NewBuffer(nil)
	c, esc := s.scanChar()
	for ; !(c == p && !esc); c, esc = s.scanChar() {
		ret.WriteRune(c)
	}
	return ret.String()
}

func (s *Lexer) scanLongString() string {
	ret := bytes.NewBuffer(nil)
	lv := 0
	for s.peek() == '=' {
		lv++
		s.next()
	}
	if s.peek() != '[' {
		panic("invalid long string")
	}
	s.next()

	end := bytes.NewBuffer(nil)
	for true {
		c := s.peek()
		s.next()
		switch c {
		case ']':
			end.WriteRune(c)
			cont := true
			for i := 0; i < lv; i++ {
				c := s.peek()
				end.WriteRune(c)
				s.next()
				if c != '=' {
					cont = false
					break
				}
			}
			if cont && s.peek() == ']' {
				s.next()
				return ret.String()
			}
			ret.WriteString(end.String())
			end.Reset()
		default:
			ret.WriteRune(c)
		}
	}
	panic("error")
}

func (s *Lexer) scanOperator() string {
	buf := bytes.NewBuffer(nil)
loop:
	for !s.isAtEOF() {
		c := s.peek()
		switch c {
		case '+', '-', '*', '/', '%', '=', '<', '>', '~', '@', '^', '&', '|', ':', '$':
			s.next()
			buf.WriteRune(c)
		default:
			break loop
		}
	}
	return buf.String()
}
