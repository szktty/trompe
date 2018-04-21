package trompe

import (
	"fmt"
	"strings"
)

var ptnWildcard = &ptnVar{"_"}

type Pattern struct {
	Comp ptnComp
}

type ptnComp interface {
	Eval(*Env, Value) bool
	Desc() string
}

func newPattern(c ptnComp) *Pattern {
	return &Pattern{c}
}

func NewPatternFromNode(n PtnNode) *Pattern {
	comp := parsePtnNode(n)
	return newPattern(comp)
}

func parsePtnNode(n PtnNode) ptnComp {
	switch n := n.(type) {
	case *UnitPtnNode:
		return &ptnUnit{}
	case *BoolPtnNode:
		return &ptnBool{n.Value}
	case *IntPtnNode:
		i, _ := StrToInt(n.Value.Text)
		return &ptnInt{i}
	case *StrPtnNode:
		return &ptnStr{n.Value.Text}
	case *VarPtnNode:
		return &ptnVar{n.Name.Text}
	default:
		panic("notimpl")
		return nil
	}
}

func (p *Pattern) Eval(env *Env, v Value) bool {
	return p.Comp.Eval(env, v)
}

func (p *Pattern) Type() int {
	return ValueTypePattern
}

func (p *Pattern) Desc() string {
	return fmt.Sprintf("<pattern %s>", p.Comp.Desc())
}

type ptnUnit struct {
}

func (p *ptnUnit) Eval(env *Env, v Value) bool {
	return v == SharedUnit
}

func (p *ptnUnit) Desc() string {
	return "()"
}

type ptnBool struct {
	v bool
}

func (p *ptnBool) Eval(env *Env, v Value) bool {
	if b, ok := ValueToBool(v); ok {
		return p.v == b.Value
	} else {
		return false
	}
}

func (p *ptnBool) Desc() string {
	if p.v {
		return "true"
	} else {
		return "false"
	}
}

type ptnInt struct {
	v int
}

func (p *ptnInt) Eval(env *Env, v Value) bool {
	if i, ok := ValueToInt(v); ok {
		return p.v == i.Value
	} else {
		return false
	}
}

func (p *ptnInt) Desc() string {
	return fmt.Sprintf("%d", p.v)
}

/*
type ptnFloat struct {
	v float64
}

func (p *ptnFloat) Eval(env *Env, v Value) bool {
	if f, ok := ValueToFloat(v); ok {
		return p.v == f.Value
	} else {
		return false
	}
}

func (p *ptnFloat) Desc() string {
	return fmt.Sprintf("%f", p.v)
}
*/

type ptnStr struct {
	v string
}

func (p *ptnStr) Eval(env *Env, v Value) bool {
	if s, ok := ValueToString(v); ok {
		return p.v == s.Value
	} else {
		return false
	}
}

func (p *ptnStr) Desc() string {
	return fmt.Sprintf("\"%s\"", p.v)
}

type ptnList struct {
	comps []ptnComp
}

func (p *ptnList) Eval(env *Env, v Value) bool {
	if l, ok := ValueToList(v); ok {
		if len(p.comps) == l.Len() {
			iter := l
			for _, comp := range p.comps {
				if !comp.Eval(env, iter.Value) {
					return false
				}
			}
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (p *ptnList) Desc() string {
	desc := "["
	for i, e := range p.comps {
		desc += e.Desc()
		if i+1 < len(p.comps) {
			desc += ", "
		}
	}
	desc += "]"
	return desc
}

type ptnTuple struct {
	comps []ptnComp
}

func newPtnTuple(c ...ptnComp) *ptnTuple {
	return &ptnTuple{c}
}

func (p *ptnTuple) Eval(env *Env, v Value) bool {
	if t, ok := ValueToTuple(v); ok {
		if len(p.comps) != t.Len() {
			return false
		}

		for i := 0; i < len(p.comps); i++ {
			e1 := p.comps[i]
			e2 := t.Values[i]
			if !e1.Eval(env, e2) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (p *ptnTuple) Desc() string {
	desc := "("
	for i, e := range p.comps {
		desc += e.Desc()
		if i+1 < len(p.comps) {
			desc += ", "
		}
	}
	desc += ")"
	return desc
}

type ptnOpt struct {
	v Pattern
}

type ptnVar struct {
	Name string
}

func (p *ptnVar) Eval(env *Env, v Value) bool {
	if !strings.HasPrefix(p.Name, "_") {
		env.Set(p.Name, v)
	}
	return true
}

func (p *ptnVar) Desc() string {
	return fmt.Sprintf("$%s", p.Name)
}

type ptnPin struct {
	Name string
}

func (p *ptnPin) Eval(env *Env, v Value) bool {
	// TODO
	return true
}

func (p *ptnPin) Desc() string {
	return fmt.Sprintf("^%s", p.Name)
}
