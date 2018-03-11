package trompe

import (
	"fmt"
	"strings"
)

var PtnWildcard = &PtnVar{"_"}

type Pattern interface {
	Eval(*Env, Value) bool
	Desc() string
}

type PtnBool struct {
	Value bool
}

func (ptn *PtnBool) Eval(env *Env, value Value) bool {
	return value.Type() == ValUnitType && ptn.Value == value.Bool()
}

func (ptn *PtnBool) Desc() string {
	if ptn.Value {
		return "true"
	} else {
		return "false"
	}
}

type PtnInt struct {
	Value int
}

func (ptn *PtnInt) Eval(env *Env, value Value) bool {
	return value.Type() == ValIntType && ptn.Value == value.Int()
}

func (ptn *PtnInt) Desc() string {
	return fmt.Sprintf("%d", ptn.Value)
}

type PtnFloat struct {
	Value float64
}

func (ptn *PtnFloat) Desc() string {
	return fmt.Sprintf("%f", ptn.Value)
}

type PtnStr struct {
	Value string
}

func (ptn *PtnStr) Eval(env *Env, value Value) bool {
	return value.Type() == ValStrType && ptn.Value == value.String()
}

func (ptn *PtnStr) Desc() string {
	return fmt.Sprintf("\"%s\"", ptn.Value)
}

type PtnList struct {
	Ptns []Pattern
}

func (ptn *PtnList) Eval(env *Env, value Value) bool {
	if value.Type() == ValListType {
		list := value.List()
		if len(ptn.Ptns) == list.Len() {
			iter := list
			for _, ptn := range ptn.Ptns {
				if !ptn.Eval(env, iter.Value) {
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

func (ptn *PtnList) Desc() string {
	desc := "["
	for i, e := range ptn.Ptns {
		desc += e.Desc()
		if i+1 < len(ptn.Ptns) {
			desc += ", "
		}
	}
	desc += "]"
	return desc
}

type PtnTuple struct {
	Ptns []Pattern
}

func CreatePtnTuple(ptn ...Pattern) *PtnTuple {
	return &PtnTuple{ptn}
}

func (ptn *PtnTuple) Eval(env *Env, value Value) bool {
	if value.Type() != ValTupleType {
		return false
	}

	tuple := value.Tuple()
	if len(ptn.Ptns) != len(tuple) {
		return false
	}

	for i := 0; i < len(ptn.Ptns); i++ {
		e1 := ptn.Ptns[i]
		e2 := tuple[i]
		if !e1.Eval(env, e2) {
			return false
		}
	}
	return true
}

func (ptn *PtnTuple) Desc() string {
	desc := "("
	for i, e := range ptn.Ptns {
		desc += e.Desc()
		if i+1 < len(ptn.Ptns) {
			desc += ", "
		}
	}
	desc += ")"
	return desc
}

type PtnOpt struct {
	Value Pattern
}

type PtnVar struct {
	Name string
}

func (ptn *PtnVar) Eval(env *Env, value Value) bool {
	if !strings.HasPrefix(ptn.Name, "_") {
		// TODO: set attr
	}
	return true
}

func (ptn *PtnVar) Desc() string {
	return fmt.Sprintf("$%s", ptn.Name)
}

type PtnPin struct {
	Name string
}

func (ptn *PtnPin) Eval(env *Env, value Value) bool {
	// TODO
	return true
}

func (ptn *PtnPin) Desc() string {
	return fmt.Sprintf("^%s", ptn.Name)
}
