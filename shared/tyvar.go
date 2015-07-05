package trompe

type TyvarEnv struct {
	Map    map[string]Type
	Tyvars []string
}

func TyvarsOfType(ty Type) []string {
	return collectTyvarsOfType(ty, make([]string, 0))
}

func collectTyvarsOfType(ty Type, accu []string) []string {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return collectTyvarsOfType(desc.Type, accu)
		} else {
			return accu
		}
	case *TypePoly:
		for _, tyvar := range desc.Tyvars {
			if !ContainsString(accu, tyvar) {
				accu = append(accu, tyvar)
			}
		}
		return accu
	default:
		return accu
	}
}

func NewTyvarEnv() *TyvarEnv {
	return &TyvarEnv{Map: make(map[string]Type), Tyvars: make([]string, 0)}
}

func TyvarEnvOfLists(names []string, tys []Type) *TyvarEnv {
	env := NewTyvarEnv()
	env.SetLists(names, tys)
	return env
}

func (self *TyvarEnv) Length() int {
	return len(self.Map)
}

func (self *TyvarEnv) Copy() *TyvarEnv {
	cpy := NewTyvarEnv()
	for k, v := range self.Map {
		cpy.Set(k, v)
	}
	return cpy
}

func (self *TyvarEnv) Get(name string) (Type, bool) {
	v, ok := self.Map[name]
	return v, ok
}

func (self *TyvarEnv) Set(name string, ty Type) {
	self.Map[name] = ty
	for _, e := range self.Tyvars {
		if e == name {
			return
		}
	}
	self.Tyvars = append(self.Tyvars, name)
}

func (self *TyvarEnv) SetVar(k string, v string) {
	self.Set(k, TVar(v))
}

func (self *TyvarEnv) SetLists(names []string, tys []Type) {
	if len(names) != len(tys) {
		panic("error")
	}
	for i, name := range names {
		self.Set(name, tys[i])
	}
}

func (self *TyvarEnv) NewTyvar() string {
	return TyvarNames[len(self.Map)]
}
