package trompe

import (
	"fmt"
)

type Module struct {
	Parent     *Module
	Submods    []*Module
	Name       string
	Includes   []*Module
	FieldTypes map[string]Type
	FieldVals  map[string]Value
	Types      map[string]Type
	Exns       map[string]Type
	Tycons     map[string]Tycon
}

func NewModule(name string) *Module {
	return &Module{Name: name, Includes: make([]*Module, 0),
		FieldTypes: make(map[string]Type),
		Types:      make(map[string]Type), Exns: make(map[string]Type),
		Tycons: make(map[string]Tycon), FieldVals: make(map[string]Value),
		Submods: make([]*Module, 0)}
}

func (m *Module) FindType(name string) (Type, bool) {
	fv, ok := m.FieldTypes[name]
	return fv, ok
}

func (m *Module) FindExnType(name string) (Type, bool) {
	fv, ok := m.Exns[name]
	return fv, ok
}

func (m *Module) SetType(name string, ty Type) {
	m.Types[name] = ty
}

func (m *Module) SetTycon(name string, tyc Tycon) {
	m.Tycons[name] = tyc
}

func (m *Module) SetExn(name string, ty ...Type) {
	m.Exns[name] = NewTypeExn(name, ty...)
}

func (m *Module) FindFieldType(name string) (Type, bool) {
	ty, ok := m.FieldTypes[name]
	return ty, ok
}

func (m *Module) SetFieldType(name string, ty Type) {
	m.FieldTypes[name] = ty
}

func (m *Module) FindFieldValue(name string) (Value, bool) {
	Debugf("find %s.%s", m.Name, name)
	fv, ok := m.FieldVals[name]
	return fv, ok
}

func (m *Module) SetFieldValue(name string, v Value) {
	m.FieldVals[name] = v
}

func (m *Module) SetPrim(name string,
	f func(*State, *Context, []Value) (Value, error)) {
	m.FieldVals[name] = Primitive(f)
}

func (m *Module) MergeEnv(env *Env) {
	for k, v := range env.Venv {
		m.SetFieldType(k, v)
	}
	for k, v := range env.Tenv {
		m.SetType(k, v)
	}
	for k, v := range env.Eenv {
		m.SetExn(k, v)
	}
	for k, v := range env.Tycons {
		m.SetTycon(k, v)
	}
}

func (m *Module) FindModule(name string) (*Module, bool) {
	for _, sub := range m.Submods {
		if sub.Name == name {
			return sub, true
		}
	}
	return nil, false
}

func (m *Module) AddModule(other *Module) {
	m.Submods = append(m.Submods, other)
	other.Parent = m
}

func (m *Module) AddInclude(other *Module) {
	for k, v := range other.FieldVals {
		m.FieldVals[k] = v
	}
	for k, v := range other.FieldTypes {
		m.FieldTypes[k] = v
	}
	for k, v := range other.Types {
		m.Types[k] = v
	}
	for k, v := range other.Tycons {
		m.Tycons[k] = v
	}
}

func (m *Module) Path() *NamePath {
	comps := make([]string, 0)
	if m != nil {
		comps = append(comps, m.Name)
		m = m.Parent
	}
	return NamePathOfComps(RevString(comps))
}

func (m *Module) PrintIntf() {
	for name, ty := range m.Types {
		fmt.Printf("type %s = %s\n", name, StringOfType(ty))
	}
	for name, ty := range m.FieldTypes {
		fmt.Printf("val %s : %s\n", name, StringOfType(ty))
	}
}
