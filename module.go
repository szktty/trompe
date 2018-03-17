package trompe

import (
	"strings"
)

type Module struct {
	Parent  *Module
	Subs    map[string]*Module
	Name    string
	File    string
	Attrs   map[string]Value
	Imports []*Module
}

var root = CreateModule(nil, "")

var sep = "."

func GetModule(path string) *Module {
	comps := strings.Split(path, sep)
	owner := root
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
	root.AddSub(m)
}

func CreateModule(parent *Module, name string) *Module {
	return &Module{
		Parent:  parent,
		Subs:    make(map[string]*Module, 8),
		Name:    name,
		Attrs:   make(map[string]Value, 16),
		Imports: make([]*Module, 4)}
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
