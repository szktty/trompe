package trompe

import (
	"bytes"
	"fmt"
)

type LocalScope struct {
	Asis     *AnalysisBlock
	Outer    *LocalScope
	Stack    []*VarDef
	VarMap   map[string]*VarDef
	Copied   []*VarDef
	Shared   *VarDef
	Globals  []*VarDef
	Rec      []*VarDef
	RecMap   map[*VarDef]string
	Bindings map[string]*NamePath

	// TODO: deprecated: Copied と意味が同じ
	Refs []*VarDef
}

const (
	_VarDefType = iota
	VarDefTypeArgVar
	VarDefTypeTempVar
	VarDefTypeRefVar
	VarDefTypeGroupVar
	VarDefTypeGlobalVar
	VarDefTypeRecVar
	VarDefTypeDynamicVar
)

type VarDef struct {
	Scope *LocalScope
	Name  string
	Type  int
	Index int

	Ref   *VarDef   // TTempVar (outer variable)
	Group []*VarDef // TGroupVar
	RecId string    // TRecVar
	Owner *VarDef   // TDynamicVar
}

func NewVarDef(owner *LocalScope, ty int, name string) *VarDef {
	return &VarDef{Scope: owner, Type: ty, Name: name}
}

func NewLocalScope(outer *LocalScope) *LocalScope {
	return &LocalScope{Outer: outer, Stack: make([]*VarDef, 0),
		VarMap:   make(map[string]*VarDef),
		Rec:      make([]*VarDef, 0),
		RecMap:   make(map[*VarDef]string),
		Bindings: make(map[string]*NamePath)}
}

func (sc *LocalScope) NumArgs() int {
	n := 0
	for _, va := range sc.Stack {
		if va.Type == VarDefTypeArgVar {
			n++
		}
	}
	return n
}

func (sc *LocalScope) NumTempsOnly() int {
	n := 0
	for _, va := range sc.Stack {
		if va.Type == VarDefTypeTempVar {
			n++
		}
	}
	return n
}

func (sc *LocalScope) NumTempSlots() int {
	n := sc.NumTempsOnly()
	if sc.Shared != nil {
		return n + 1
	} else {
		return n
	}
}

func (sc *LocalScope) NumCopied() int {
	return len(sc.Copied)
}

func (sc *LocalScope) NumSharedSlots() int {
	if sc.Outer == nil || sc.Shared == nil {
		return 0
	} else {
		n := 0
		for _, temp := range sc.Shared.Group {
			shared := true
			for _, global := range sc.Globals {
				if temp == global {
					shared = false
					break
				}
			}
			if shared {
				n++
			}
		}
		return n
	}
}

func (sc *LocalScope) IsClean() bool {
	return !sc.IsFull() && !sc.IsCopying()
}

func (sc *LocalScope) IsCopying() bool {
	return len(sc.Copied) > 0
}

func (sc *LocalScope) IsFull() bool {
	return len(sc.Globals) > 0 || sc.Outer.Shared != nil
}

func (sc *LocalScope) IsFullCopying() bool {
	return sc.IsFull() && sc.IsCopying()
}

func (sc *LocalScope) Add(va *VarDef) {
	if va.Name != "" {
		sc.VarMap[va.Name] = va
	}
	sc.Stack = append(sc.Stack, va)
}

func (sc *LocalScope) AddNewArg(name string) {
	va := NewVarDef(sc, VarDefTypeArgVar, name)
	sc.Add(va)
}

func (sc *LocalScope) AddNewTemp(name string) {
	va := NewVarDef(sc, VarDefTypeTempVar, name)
	sc.Add(va)
}

func (sc *LocalScope) AddNewGlobal(name string) {
	for _, va := range sc.Globals {
		if va.Name == name {
			return
		}
	}
	va := NewVarDef(sc, VarDefTypeGlobalVar, name)
	sc.Globals = append(sc.Globals, va)
}

func (sc *LocalScope) AddNewCopied(name string) {
	va := sc.Outer.VarMap[name]
	sc.AddCopied(va)
	Debugf("add new copied %s", name)
}

func (sc *LocalScope) AddCopied(va *VarDef) {
	if ref, isNew := sc.AddRef(va); isNew {
		sc.Copied = append(sc.Copied, ref)
	}
}

func (sc *LocalScope) AddRef(va *VarDef) (*VarDef, bool) {
	for _, ref := range sc.Refs {
		deref := sc.Deref(ref)
		if deref.Type == VarDefTypeGroupVar &&
			va.Type == VarDefTypeGroupVar &&
			len(deref.Group) == len(va.Group) {
			all := true
			for i := 0; i < len(deref.Group); i++ {
				if deref.Group[i].Name != va.Group[i].Name {
					all = false
					break
				}
			}
			if all {
				return ref, false
			}
		}
	}
	ref := NewVarDef(sc, VarDefTypeRefVar, va.Name)
	ref.Ref = va
	sc.Add(ref)
	sc.Refs = append(sc.Refs, ref)
	return ref, true
}

// TODO: deprecated
/*
func (sc *LocalScope) AddNewShared(name string, indirect bool) {
	va := NewVarDef(sc, VarDefTypeTempVar, name)
	if !indirect {
		if sc.Copied == nil {
			sc.Copied = make([]*VarDef, 0)
		}
		sc.Copied = append(sc.Copied, va)
	} else {
		if sc.Shared == nil {
			sc.Shared = NewVarDef(sc, VarDefTypeGroupVar, "")
		}
		sc.Shared.Group = append(sc.Shared.Group, va)
	}
}
*/

func (sc *LocalScope) SetNewShared(names []string) {
	if len(names) == 0 {
		panic("shared variables must not be 0")
	}
	grp := NewVarDef(sc, VarDefTypeGroupVar, "")
	grp.Group = make([]*VarDef, len(names))
	for i, name := range names {
		va := NewVarDef(sc, VarDefTypeDynamicVar, name)
		va.Index = i
		va.Owner = grp
		grp.Group[i] = va
		sc.VarMap[name] = grp
	}
	sc.Add(grp)
	sc.Shared = grp
}

func (sc *LocalScope) AddBinding(name string, path *NamePath) {
	sc.Bindings[name] = path
}

func (sc *LocalScope) AddNewRec(name string) {
	if _, ok := sc.FindRec(name); ok {
		return
	}
	va := NewVarDef(sc, VarDefTypeRecVar, name)
	va.RecId = sc.NewRecId(va)
	sc.Rec = append(sc.Rec, va)
	sc.RecMap[va] = va.RecId
}

func (sc *LocalScope) NewRecId(va *VarDef) string {
	nest := 0
	e := sc
	for e != va.Scope {
		nest++
		e = e.Outer
	}
	count := len(sc.Rec)
	return fmt.Sprintf("#rec_s%d_f%d", nest, count)
}

// TODO: deprecated?
func (sc *LocalScope) FindRec(name string) (*VarDef, bool) {
	for sc != nil {
		for _, va := range sc.Rec {
			if va.Name == name {
				return va, true
			}
		}
		sc = sc.Outer
	}
	return nil, false
}

func (sc *LocalScope) FindShared(name string) (*VarDef, bool) {
	if sc.Shared == nil {
		return nil, false
	}

	for _, temp := range sc.Shared.Group {
		if temp.Name == name {
			if temp.Type == VarDefTypeDynamicVar {
				return temp.Owner, true
			} else {
				return temp, true
			}
		}
	}
	if sc.Outer == nil {
		return nil, false
	} else {
		return sc.Outer.FindShared(name)
	}
}

func (sc *LocalScope) String() string {
	buf := bytes.NewBufferString("<LocalScope ")
	buf.WriteString(fmt.Sprintf("Outer:%p VarDef:[", sc.Outer))
	for i, va := range sc.Stack {
		buf.WriteString(fmt.Sprintf("%d:", i))
		if va == sc.Shared {
			buf.WriteString("<Shared>")
		} else {
			write := true
			for _, copied := range sc.Copied {
				if copied == va {
					buf.WriteString("<")
					sc.WriteStringOfVar(buf, va)
					buf.WriteString(">")
					write = false
					break
				}
			}
			if write {
				sc.WriteStringOfVar(buf, va)
			}
		}
		if i+1 < len(sc.Stack) {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("]")

	if len(sc.Refs) > 0 {
		buf.WriteString(" Refs:[")
		for i, va := range sc.Refs {
			sc.WriteStringOfVar(buf, va)
			if i+1 < len(sc.Refs) {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("]")
	}

	if sc.Copied != nil {
		buf.WriteString(" Copied:[")
		for i, va := range sc.Copied {
			buf.WriteString(fmt.Sprintf("%d:", i))
			sc.WriteStringOfVar(buf, va)
			if i+1 < len(sc.Copied) {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("]")
	}

	if sc.Shared != nil {
		buf.WriteString(" Shared:")
		sc.WriteStringOfVar(buf, sc.Shared)
	}

	if sc.Globals != nil {
		buf.WriteString(" Globals:[")
		for i, va := range sc.Globals {
			sc.WriteStringOfVar(buf, va)
			if i+1 < len(sc.Globals) {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("]")
	}

	if len(sc.Bindings) > 0 {
		buf.WriteString(" Bindings:[")
		i := 0
		for name, bin := range sc.Bindings {
			buf.WriteString(fmt.Sprintf("%s:%s", name, bin.String()))
			if i+1 < len(sc.Bindings) {
				buf.WriteString(" ")
			}
			i++
		}
		buf.WriteString("]")
	}

	buf.WriteString(">")
	return buf.String()
}

func (sc *LocalScope) WriteStringOfVar(buf *bytes.Buffer, va *VarDef) {
	switch va.Type {
	case VarDefTypeGroupVar:
		buf.WriteString("[")
		for i, e := range va.Group {
			buf.WriteString(fmt.Sprintf("%d:%s", i, e.Name))
			if i+1 < len(va.Group) {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("]")
	case VarDefTypeRefVar:
		sc.WriteStringOfVar(buf, va.Ref)
	default:
		buf.WriteString(va.Name)
	}
}

func (sc *LocalScope) Finish() {
	sc.Refs = make([]*VarDef, 0)
	for i, va := range sc.Stack {
		va.Index = i
		if va.Type == VarDefTypeRefVar {
			sc.Refs = append(sc.Refs, va)
		}
	}
}

func (sc *LocalScope) Deref(va *VarDef) *VarDef {
	if va.Type == VarDefTypeRefVar {
		return sc.Deref(va.Ref)
	} else {
		return va
	}
}

func (sc *LocalScope) FindTemp(name string) (*VarDef, bool) {
	for _, va := range sc.Stack {
		deref := sc.Deref(va)
		if deref.Type != VarDefTypeGroupVar && deref.Name == name {
			return va, true
		}
	}
	return nil, false
}

func (sc *LocalScope) FindGlobal(name string) (*VarDef, bool) {
	for _, va := range sc.Globals {
		if va.Name == name {
			return va, true
		}
	}
	return nil, false
}

func (sc *LocalScope) FindMember(name string) (*VarDef, *VarDef, bool) {
	for _, va := range sc.Stack {
		deref := sc.Deref(va)
		if deref.Type == VarDefTypeGroupVar {
			for _, mem := range deref.Group {
				if mem.Name == name {
					return va, mem, true
				}
			}
		}
	}
	return nil, nil, false
}
