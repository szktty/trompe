package trompe

import (
	"fmt"
)

type Env struct {
	Parent      *Env
	Venv        map[string]Type
	Tenv        map[string]Type
	Eenv        map[string]Type // exception
	Tycons      map[string]Tycon
	Prims       map[string]string
	ModVars     map[string]*Module
	Imports     []*Module
	FullImports []*Module
	Metas       []*TypeMeta
}

func NewEnv(parent *Env) *Env {
	return &Env{Parent: parent, Venv: make(map[string]Type),
		Tenv: make(map[string]Type), Eenv: make(map[string]Type),
		Tycons: make(map[string]Tycon), Prims: make(map[string]string),
		ModVars: make(map[string]*Module),
		Metas:   make([]*TypeMeta, 0), Imports: make([]*Module, 0)}
}

func (env *Env) NewMeta() *TypeMeta {
	i := 0
	e := env
	for e != nil {
		i += len(e.Metas)
		e = e.Parent
	}
	name := fmt.Sprintf("_%d", i)
	meta := &TypeMeta{Name: name}
	env.Metas = append(env.Metas, meta)
	return meta
}

func (env *Env) AddType(name string, ty Type) {
	env.Tenv[name] = ty
}

func (env *Env) AddVarType(name string, ty Type) {
	env.Venv[name] = ty
}

func (env *Env) AddExnType(name string, ty ...Type) {
	env.Eenv[name] = NewTypeExn(name, ty...)
}

func (env *Env) FindVarType(name string) (Type, bool) {
	for env != nil {
		if ty, ok := env.Venv[name]; ok {
			return ty, true
		} else {
			for _, m := range env.FullImports {
				if ty, ok := m.FindFieldType(name); ok {
					return ty, true
				}
			}
		}
		env = env.Parent
	}
	return nil, false
}

func (env *Env) FindExnType(name string) (Type, bool) {
	for env != nil {
		if ty, ok := env.Eenv[name]; ok {
			return ty, true
		} else {
			for _, m := range env.FullImports {
				if ty, ok := m.FindExnType(name); ok {
					return ty, true
				}
			}
		}
		env = env.Parent
	}
	return nil, false
}

func (env *Env) FindModuleOfVar(name string) (*Module, bool) {
	for env != nil {
		if mod, ok := env.ModVars[name]; ok {
			return mod, true
		}
		env = env.Parent
	}
	return nil, false
}

func (env *Env) FindModuleOfPath(path *NamePath) (*Module, error) {
	var mod *Module
	var mname string
	if len(path.Base) == 0 {
		mname = path.Name
	} else {
		// TODO
		mname = path.Base[0]
	}

	for env != nil {
		if ty, ok := env.FindVarType(mname); ok {
			if deref, ok := DerefType(ty); ok {
				if app, ok := deref.(*TypeApp); ok {
					if tycMod, ok := app.Tycon.(*TyconModule); ok {
						mod = tycMod.Module
						break
					}
				}
			}
		}
		for _, m := range env.Imports {
			if m.Name == mname {
				mod = m
				break
			}
		}
		for _, m := range env.FullImports {
			for _, m := range m.Submods {
				if m.Name == mname {
					mod = m
					break
				}
			}
		}
		env = env.Parent
	}
	if mod == nil {
		return nil, fmt.Errorf("Unbound module `%s'", mname)
	}

	for i := 1; i < len(path.Base); i++ {
		if sub, ok := mod.FindModule(path.Base[i]); ok {
			mod = sub
		} else {
			return nil, fmt.Errorf("Unbound module value `%s'", path.StringUpto(i))
		}
	}
	return mod, nil
}

func (env *Env) FindType(name string) (Type, bool) {
	for env != nil {
		if ty, ok := env.Tenv[name]; ok {
			return ty, true
		} else {
			env = env.Parent
		}
	}
	return nil, false
}

func (env *Env) AddFullImport(m *Module) {
	for _, m1 := range env.FullImports {
		if m == m1 {
			return
		}
	}
	env.FullImports = append(env.FullImports, m)
}
