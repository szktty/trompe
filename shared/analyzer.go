package trompe

// analyze variable scope in blocks

import (
	"fmt"
	"strings"
)

type Analyzer struct {
	State   *State
	Blocks  []*AnalysisBlock
	Current *AnalysisBlock
	Scope   map[string]string
	Nodes   map[*TypedNode]*AnalysisBlock
}

const (
	_AnalysisVarType    = iota
	AnalysisVarTypeCopy = 1
	AnalysisVarTypeFull = 1 << 1
	AnalysisVarTypePerm = 1 << 2 // global
	AnalysisVarTypeArg  = 1 << 3
)

type AnalysisVar struct {
	Name  string
	Home  *AnalysisBlock
	Types int
}

type AnalysisBlock struct {
	Node     *TypedNode
	Parent   *AnalysisBlock
	Name     string
	Import   []*AnalysisModule  // imported modules
	Perms    []*AnalysisVar     // current module variables
	Refs     []*AnalysisVar     // variables from outer scope
	Shared   []*AnalysisVar     // variable shared with inner blocks
	Temps    []*AnalysisVar     // temporary variables
	Bindings map[string]*Module // outer module variable
	Scope    *LocalScope

	// TODO: deprecated?
	Rec []*AnalysisVar // recursive functions
	Top bool
}

type AnalysisModule struct {
	Module *Module
	Fully  bool
}

func (state *State) AnalyzeBlock(node *TypedNode) (*Analyzer, error) {
	Debugf("==> BEGIN analyze block")
	azer := &Analyzer{State: state,
		Blocks: make([]*AnalysisBlock, 0),
		Scope:  make(map[string]string),
		Nodes:  make(map[*TypedNode]*AnalysisBlock)}
	// TODO: analyze error
	azer.Analyze(node)
	if LogGroupEnabled(LogGroupDebug) {
		for i, asis := range azer.Blocks {
			Debugf("%d: %s", i, asis)
		}
	}
	for _, blk := range azer.Blocks {
		blk.Finish()
	}
	Debugf("<== END analyze block")
	return azer, nil
}

func (asis *AnalysisVar) IsArg() bool {
	return asis.Types&AnalysisVarTypeArg > 0
}

func (asis *AnalysisVar) IsFullyShared() bool {
	return asis.Types&AnalysisVarTypeFull > 0
}

func (asis *AnalysisVar) IsPermFull() bool {
	return asis.Types&AnalysisVarTypePerm > 0
}

func (asis *AnalysisVar) DoCopy() {
	asis.Types |= AnalysisVarTypeCopy
}

func (asis *AnalysisVar) DoFullCopy() {
	asis.Types |= AnalysisVarTypeCopy | AnalysisVarTypeFull
	asis.Home.AddShared(asis)
}

func (asis *AnalysisVar) DoPermFull() {
	asis.Types |= AnalysisVarTypePerm | AnalysisVarTypeFull
	asis.Home.AddShared(asis)
}

func NewAnalysisBlock(parent *AnalysisBlock) *AnalysisBlock {
	return &AnalysisBlock{Parent: parent,
		Import:   make([]*AnalysisModule, 0),
		Bindings: make(map[string]*Module),
		Perms:    make([]*AnalysisVar, 0),
		Refs:     make([]*AnalysisVar, 0),
		Shared:   make([]*AnalysisVar, 0),
		Temps:    make([]*AnalysisVar, 0),
		Rec:      make([]*AnalysisVar, 0)}
}

func (asis *AnalysisBlock) TopBlock() *AnalysisBlock {
	if asis.Parent == nil {
		return asis
	} else {
		return asis.Parent.TopBlock()
	}
}

func (asis *AnalysisBlock) IsTopBlock() bool {
	return asis == asis.TopBlock()
}

func (asis *AnalysisBlock) AddImport(m *Module, fully bool) {
	for _, other := range asis.Import {
		if other.Module == m && other.Fully == fully {
			return
		}
	}
	va := &AnalysisModule{Module: m, Fully: fully}
	asis.Import = append(asis.Import, va)
}

func (asis *AnalysisBlock) NumArgs() int {
	i := 0
	for _, va := range asis.Temps {
		if va.IsArg() {
			i++
		}
	}
	return i
}

func (asis *AnalysisBlock) NumTemps() int {
	return len(asis.Temps)
}

func (asis *AnalysisBlock) NumTempSlots() int {
	temps := len(asis.Temps)
	shared := len(asis.Shared)
	if shared == 0 {
		return temps
	} else {
		return temps - shared + 1
	}
}

func (asis *AnalysisBlock) AddNewTemp(name string) *AnalysisVar {
	if va, ok := asis.ContainsTempOrPerm(name); ok {
		return va
	}
	va := &AnalysisVar{Name: name, Home: asis}
	asis.Temps = append(asis.Temps, va)
	return va
}

func (asis *AnalysisBlock) AddNewArg(name string) *AnalysisVar {
	temp := asis.AddNewTemp(name)
	temp.Types |= AnalysisVarTypeArg
	return temp
}

func (asis *AnalysisBlock) FindTemp(name string) (*AnalysisVar, bool) {
	for asis != nil {
		for _, temp := range asis.Temps {
			if temp.Name == name {
				return temp, true
			}
		}
		asis = asis.Parent
	}
	return nil, false
}

// TODO: deprecated
func (asis *AnalysisBlock) FindTempRec(name string) (*AnalysisVar, bool) {
	for asis != nil {
		if va, ok := asis.FindTemp(name); ok {
			return va, true
		}
		asis = asis.Parent
	}
	return nil, false
}

func (asis *AnalysisBlock) FindPerm(name string) (*AnalysisVar, bool) {
	for _, va := range asis.Perms {
		if va.Name == name {
			return va, true
		}
	}
	return nil, false
}

func (asis *AnalysisBlock) FindPermRec(name string) (*AnalysisVar, bool) {
	for asis != nil {
		if va, ok := asis.FindPerm(name); ok {
			return va, true
		}
		asis = asis.Parent
	}
	return nil, false
}

func (asis *AnalysisBlock) ContainsTempOrPerm(name string) (*AnalysisVar, bool) {
	for _, temp := range asis.Temps {
		if temp.Name == name {
			return temp, true
		}
	}
	for _, temp := range asis.Perms {
		if temp.Name == name {
			return temp, true
		}
	}
	return nil, false
}

func (asis *AnalysisBlock) AddNewPerm(name string) *AnalysisVar {
	if va, ok := asis.ContainsTempOrPerm(name); ok {
		return va
	}
	va := &AnalysisVar{Name: name, Home: asis}
	va.Types |= AnalysisVarTypePerm
	asis.Perms = append(asis.Perms, va)
	return va
}

func (asis *AnalysisBlock) AddPerm(temp *AnalysisVar) {
	if asis.HasVar(asis.Perms, temp) {
		return
	}
	temp.Types |= AnalysisVarTypePerm
	asis.Perms = append(asis.Perms, temp)
}

func (asis *AnalysisBlock) HasVar(list []*AnalysisVar,
	temp *AnalysisVar) bool {
	for _, va := range list {
		if va == temp {
			return true
		}
	}
	return false
}

func (asis *AnalysisBlock) AddShared(temp *AnalysisVar) {
	if asis.HasVar(asis.Shared, temp) {
		return
	}
	asis.Shared = append(asis.Shared, temp)
	if temp.IsPermFull() {
		asis.Perms = append(asis.Perms, temp)
	}
	if temp.Home != asis && asis.Parent != nil {
		asis.AddRef(temp)
		asis.Parent.AddShared(temp)
	}
}

func (asis *AnalysisBlock) AddRef(temp *AnalysisVar) {
	for _, va := range asis.Refs {
		if va == temp {
			return
		}
	}
	temp.Types |= AnalysisVarTypeCopy
	asis.Refs = append(asis.Refs, temp)
}

func (asis *AnalysisBlock) AddRec(temp *AnalysisVar) {
	for _, rec := range asis.Rec {
		if rec == temp {
			return
		}
	}

	home := asis.Parent
findHome:
	for home != nil {
		for _, rec := range home.Rec {
			if rec == temp {
				Debugf("found home %p for rec %s", home, temp.Name)
				break findHome
			}
		}
	}
	if home == nil {
		home = asis
	}

	//temp.Types |= TRec
	temp.Home = home
	asis.Rec = append(asis.Rec, temp)
}

func (asis *AnalysisBlock) AddNewRec(name string) {
	va := &AnalysisVar{Name: name, Home: asis,
		Types: AnalysisVarTypeFull}
	asis.AddRec(va)
	asis.AddShared(va)
}

func (asis *AnalysisBlock) FindRec(name string) (*AnalysisVar, bool) {
	for asis != nil {
		for _, rec := range asis.Rec {
			if rec.Name == name {
				return rec, true
			}
		}
		asis = asis.Parent
	}
	return nil, false
}

func (asis *AnalysisBlock) FindBinding(name string) (*Module, bool) {
	for asis != nil {
		for _, m := range asis.Import {
			if m.Fully {
				if _, ok := m.Module.FindFieldValue(name); ok {
					return m.Module, true
				}
			}
		}
		asis = asis.Parent
	}
	return nil, false
}

func (asis *AnalysisBlock) AddNewBinding(m *Module, name string) {
	for k, _ := range asis.Bindings {
		if k == name {
			return
		}
	}
	asis.Bindings[name] = m
}

func (asis *AnalysisBlock) VarNames(vars []*AnalysisVar) []string {
	names := make([]string, len(vars))
	for i, va := range vars {
		names[i] = va.Name
	}
	return names
}

func (asis *AnalysisBlock) String() string {
	temps := asis.VarNames(asis.Temps)
	shared := asis.VarNames(asis.Shared)
	refs := asis.VarNames(asis.Refs)
	rec := asis.VarNames(asis.Rec)
	globals := asis.VarNames(asis.Perms)
	binds := make([]string, 0)
	for name, mod := range asis.Bindings {
		binds = append(binds, fmt.Sprintf("%s:%s", name, mod.Path().String()))
	}
	return fmt.Sprintf("<AnalysisBlock %p Parent=%p Refs=[%s] Shared=[%s] Temps=[%s] Rec=[%s] Perms=[%s] Bindings=[%s]>",
		asis, asis.Parent,
		strings.Join(refs, " "),
		strings.Join(shared, " "),
		strings.Join(temps, " "),
		strings.Join(rec, " "),
		strings.Join(globals, " "),
		strings.Join(binds, " "))
}

func (azer *Analyzer) BeginBlock(node *TypedNode) {
	azer.Current = NewAnalysisBlock(azer.Current)
	azer.Current.Node = node
	azer.Nodes[node] = azer.Current
	azer.Blocks = append(azer.Blocks, azer.Current)
}

func (azer *Analyzer) EndBlock() {
	azer.Current = azer.Current.Parent
}

func (azer *Analyzer) ScopeOfNode(node *TypedNode) *LocalScope {
	return azer.Nodes[node].Scope
}

func (azer *Analyzer) Analyze(node *TypedNode) {
	switch desc := node.Desc.(type) {
	case *TypedProgramNode:
		azer.BeginBlock(node)
		azer.Current.Top = true
		azer.Current.AddImport(azer.State.Root, true)
		for _, e := range desc.Items {
			azer.Analyze(e)
		}
		azer.EndBlock()

	case *TypedLetNode:
		if desc.Rec {
			recs := make([]string, 0)
			for _, e := range desc.Bindings {
				switch bind := e.Desc.(type) {
				case *TypedLetBindingNode:
					switch ptn := bind.Ptn.Desc.(type) {
					case *PtnIdentNode:
						// azer.Current.AddNewRec(ptn.Name)
						recs = append(recs, ptn.Name)
					}
				case *TypedBlockNode:
					// azer.Current.AddNewRec(name)
					recs = append(recs, bind.Name.NameExn())
				}
			}
			isTop := azer.Current.IsTopBlock()
			for _, rec := range recs {
				temp := azer.Current.AddNewTemp(rec)
				if isTop {
					temp.DoPermFull()
				} else {
					temp.DoFullCopy()
				}
			}
		}
		for _, e := range desc.Bindings {
			azer.Analyze(e)
		}
		if desc.Body != nil {
			azer.Analyze(desc.Body)
		}

	case *TypedLetBindingNode:
		azer.Analyze(desc.Ptn)
		azer.Analyze(desc.Body)

	case *TypedBlockNode:
		name := desc.Name.NameExn()
		// TODO: deprecated
		/*
			parent := azer.Current
			if _, ok := parent.Node.Desc.(*TypedProgramNode); ok {
				parent.AddNewPerm(name)
			} else if _, ok := parent.FindRec(name); !ok {
				parent.AddNewTemp(name)
			}
		*/
		temp := azer.Current.AddNewTemp(name)
		if azer.Current.IsTopBlock() {
			temp.DoPermFull()
		}

		azer.BeginBlock(node)
		azer.AnalyzeParamList(desc.Params)
		azer.Analyze(desc.Body)
		azer.EndBlock()

	case *TypedIdentNode:
		if temp, ok := azer.Current.FindTemp(desc.Name); ok {
			if temp.Home != azer.Current {
				azer.Current.AddRef(temp)
				azer.Current.Parent.AddShared(temp)
				if temp.IsPermFull() {
					azer.Current.AddPerm(temp)
				}
			}
		} else if m, ok := azer.Current.FindBinding(desc.Name); ok {
			azer.Current.Bindings[desc.Name] = m
		} else {
			Panicf("unbound %s", desc.Name)
		}
		// TODO: deprecated
		/*
			if _, ok := azer.Current.FindTemp(desc.Name); ok {
				// do nothing
			} else if va, ok := azer.Current.FindTempRec(desc.Name); ok {
				azer.Current.AddRef(va)
				azer.Current.Parent.AddShared(va)
			} else if va, ok := azer.Current.FindRec(desc.Name); ok {
				azer.Current.AddRec(va)
				azer.Current.AddRef(va)
				azer.Current.Parent.AddShared(va)
			} else if va, ok := azer.Current.FindPermRec(desc.Name); ok {
				azer.Current.AddPerm(va)
			} else if m, ok := azer.Current.FindBinding(desc.Name); ok {
				azer.Current.AddNewBinding(m, desc.Name)
			} else {
				Panicf("unbound %s", desc.Name)
			}
		*/

	case *TypedKeywordNode:
		azer.Analyze(desc.Exp)

	case *TypedSeqExpNode:
		for _, e := range desc.Exps {
			azer.Analyze(e)
		}

	case *TypedForNode:
		azer.Current.AddNewTemp(desc.Name.NameExn())
		azer.Analyze(desc.Init)
		azer.Analyze(desc.Limit)
		azer.Analyze(desc.Body)

	case *TypedIfNode:
		azer.Analyze(desc.Cond)
		azer.Analyze(desc.True)
		if desc.False != nil {
			azer.Analyze(desc.False)
		}

	case *TypedAppNode:
		azer.Analyze(desc.Exp)
		for _, arg := range desc.Args {
			azer.Analyze(arg)
		}

	case *TypedCaseNode:
		azer.Analyze(desc.Exp)
		for _, m := range desc.Match {
			azer.Analyze(m)
		}

	case *TypedMatchNode:
		if desc.Cond != nil {
			azer.Analyze(desc.Cond)
		}
		if desc.Ptn != nil {
			azer.Analyze(desc.Ptn)
		}
		azer.Analyze(desc.Body)

	case *TypedPtnIdentNode:
		azer.Current.AddNewTemp(desc.Name)

	case *TypedPtnTupleNode:
		for _, comp := range desc.Comps {
			azer.Analyze(comp)
		}

	case *TypedPtnListNode:
		for _, e := range desc.Elts {
			azer.Analyze(e)
		}

	case *TypedPtnListConsNode:
		azer.Analyze(desc.Head)
		azer.Analyze(desc.Tail)

	case *TypedPtnArrayNode:
		for _, e := range desc.Elts {
			azer.Analyze(e)
		}

	case *TypedWildcardNode:
	case *TypedPtnConstNode:
	case *TypedValuePathNode:
		// do nothing

	case *TypedTupleNode:
		for _, comp := range desc.Comps {
			azer.Analyze(comp)
		}

	case *TypedListNode:
		for _, e := range desc.Elts {
			azer.Analyze(e)
		}

	case *TypedArrayNode:
		for _, e := range desc.Elts {
			azer.Analyze(e)
		}

	case *TypedArrayAccessNode:
		azer.Analyze(desc.Array)
		azer.Analyze(desc.Index)
		if desc.Set != nil {
			azer.Analyze(desc.Set)
		}

	case *TypedFunNode:
		azer.Analyze(desc.MultiMatch)

	case *TypedMultiMatchNode:
		azer.BeginBlock(node)
		azer.AnalyzeParamList(desc.Params)
		if desc.Cond != nil {
			azer.Analyze(desc.Cond)
		}
		azer.Analyze(desc.Body)
		azer.EndBlock()

	case *TypedAddNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedSubNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedMulNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedDivNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedModNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedEqNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedNeNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedLtNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedLeNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedGtNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedGeNode:
		azer.Analyze(desc.Left)
		azer.Analyze(desc.Right)

	case *TypedUnitNode:
	case *TypedBoolNode:
	case *TypedIntNode:
	case *TypedFloatNode:
	case *TypedStringNode:
	case *TypedCharNode:
		// do nothing

	default:
		asis := azer.Current
		for asis != nil {
			Debugf("%s", asis)
			asis = asis.Parent
		}
		Panicf("block analyzer notimpl %s\n", node)
	}
}

func (azer *Analyzer) AnalyzeParamList(params []*TypedNode) {
	// TODO: other patterns
	for _, param := range params {
		switch ptn := param.Desc.(type) {
		case *TypedPtnIdentNode:
			azer.Current.AddNewArg(ptn.Name)
		default:
			panic(fmt.Errorf("not impl %s", ptn))
		}
	}
}

func (asis *AnalysisBlock) Finish() {
	if asis.Parent != nil {
		asis.Scope = NewLocalScope(asis.Parent.Scope)
	} else {
		asis.Scope = NewLocalScope(nil)
	}
	sc := asis.Scope
	sc.Asis = asis

	// arguments and temporary variables
	for _, temp := range asis.Temps {
		if temp.IsArg() && !temp.IsFullyShared() {
			sc.AddNewArg(temp.Name)
		}
	}
	for _, temp := range asis.Temps {
		if !temp.IsArg() && !temp.IsFullyShared() {
			sc.AddNewTemp(temp.Name)
		}
	}

	// shared variables
	if len(asis.Shared) > 0 {
		names := make([]string, 0)
		for _, temp := range asis.Shared {
			if temp.IsFullyShared() && !temp.IsPermFull() {
				names = append(names, temp.Name)
			}
		}
		if len(names) > 0 {
			sc.SetNewShared(names)
		}
	}

	// copied immutable values (outer variables)
	if len(asis.Refs) > 0 && asis.Parent != nil {
		for _, temp := range asis.Refs {
			if temp.IsPermFull() {
				sc.AddNewGlobal(temp.Name)
			} else if temp.IsFullyShared() {
				shared, ok := sc.Outer.FindShared(temp.Name)
				if !ok {
					Panicf("fully shared %s is not found", temp.Name)
				}
				sc.AddCopied(shared)
			} else if !temp.IsFullyShared() {
				sc.AddNewCopied(temp.Name)
			}
		}
	}

	// permanently fully shared  variables (module global variables)
	for _, va := range asis.Perms {
		sc.AddNewGlobal(va.Name)
	}

	// bindings
	for k, v := range asis.Bindings {
		sc.AddBinding(k, v.Path().AddName(k))
	}

	sc.Finish()
	Debugf("local scope: %s", sc)
}

func (asis *AnalysisBlock) IsRec(name string) bool {
	for _, rec := range asis.Rec {
		if rec.Name == name {
			return true
		}
	}
	return false
}
