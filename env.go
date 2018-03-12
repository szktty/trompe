package trompe

type Env struct {
	Parent *Env
	Attrs  map[string]Value
}

func CreateEnv(parent *Env) *Env {
	return &Env{parent, make(map[string]Value, 0)}
}

func (env *Env) GetAttr(name string) Value {
	for env != nil {
		if value := env.Attrs[name]; value != nil {
			return value
		}
		env = env.Parent
	}
	return nil
}

func (env *Env) SetAttr(name string, value Value) {
	env.Attrs[name] = value
}
