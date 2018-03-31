package trompe

import (
	"strings"
)

type Module struct {
	Parent *Module
	Subs   map[string]*Module
	Name   string
	File   string
	Env    *Env
}

var RootModule *Module

var OpenedModules []*Module

var sep = "."

func GetModule(path string) *Module {
	comps := strings.Split(path, sep)
	owner := RootModule
	for _, name := range comps {
		if m, ok := owner.Subs[name]; ok {
			owner = m
		} else {
			return nil
		}
	}
	return owner
}

func AddTopModule(m *Module) {
	RootModule.AddSub(m)
}

func AddOpenedModule(m *Module) {
	OpenedModules = append(OpenedModules, m)
}

func GetModuleAttr(imports []*Module, name string) Value {
	for _, m := range imports {
		if value := m.GetAttr(name); value != nil {
			return value
		}
	}
	for _, m := range OpenedModules {
		if value := m.GetAttr(name); value != nil {
			return value
		}
	}
	return nil
}

func NewModule(parent *Module, name string) *Module {
	return &Module{
		Parent: parent,
		Subs:   make(map[string]*Module, 8),
		Name:   name,
		Env:    NewEnv(nil),
	}
}

func (m *Module) Path() string {
	path := m.Name
	cur := m.Parent
	for cur != nil {
		path += sep
		path += cur.Name
		cur = cur.Parent
	}
	return path
}

func (m *Module) AddSub(sub *Module) {
	m.Subs[sub.Name] = sub
}

func (m *Module) GetAttr(name string) Value {
	return m.Env.Get(name)
}

func (m *Module) AddAttr(name string, value Value) {
	m.Env.Set(name, value)
}

func (m *Module) AddPrim(name string,
	f func(*Context, []Value, int) (Value, error),
	arity int) {
	m.AddAttr(name, NewPrim(f, arity))
}
