package trompe

import (
	"strings"
)

type NamePath struct {
	Base []string
	Name string
}

func NewNamePath() *NamePath {
	return &NamePath{Base: make([]string, 0)}
}

func NamePathOfComps(comps []string) *NamePath {
	Debugf("path = %s", comps)
	path := &NamePath{Base: make([]string, len(comps)-1)}
	for i := 0; i < (len(comps) - 1); i++ {
		path.Base[i] = comps[i]
	}
	path.Name = comps[len(comps)-1]
	return path
}

func NamePathOfWordList(words []*Word) *NamePath {
	path := NewNamePath()
	for _, w := range words {
		path.AddName(w.Value)
	}
	return path
}

func (path *NamePath) HasBase() bool {
	return len(path.Base) > 0
}

func (path *NamePath) AddName(name string) *NamePath {
	if path.Name != "" {
		path.Base = append(path.Base, path.Name)
	}
	path.Name = name
	return path
}

func (path *NamePath) BaseString() string {
	return strings.Join(path.Base, ".")
}

func (path *NamePath) String() string {
	if len(path.Base) == 0 {
		return path.Name
	} else {
		return strings.Join(path.Base, ".") + "." + path.Name
	}
}

func (path *NamePath) StringUpto(i int) string {
	return strings.Join(path.Base[0:i], ".")
}

func (path *NamePath) ModulePath() *NamePath {
	if len(path.Base) == 1 {
		return nil
	} else {
		return &NamePath{Base: path.Base[0 : len(path.Base)-3],
			Name: path.Base[len(path.Base)-1]}
	}
}

func (path *NamePath) List() []string {
	return append(path.Base, path.Name)
}
