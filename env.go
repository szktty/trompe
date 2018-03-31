package trompe

type Env struct {
	Attrs   map[string]Value
	Imports []*Module
}

func NewEnv(src *Env) *Env {
	var newMap map[string]Value
	var imports []*Module
	if src != nil {
		newMap = make(map[string]Value, len(src.Attrs))
		for k, v := range src.Attrs {
			newMap[k] = v
		}
		imports = make([]*Module, len(src.Imports))
		for i, m := range src.Imports {
			imports[i] = m
		}
	} else {
		newMap = make(map[string]Value, 16)
		imports = []*Module{}
	}
	return &Env{newMap, imports}
}

func (env *Env) AddImport(m *Module) {
	env.Imports = append(env.Imports, m)
}

func (env *Env) Get(name string) Value {
	if value := env.Attrs[name]; value != nil {
		return value
	}
	return GetModuleAttr(env.Imports, name)
}

// TODO: Set -> Add
func (env *Env) Set(name string, value Value) {
	env.Attrs[name] = value
}
