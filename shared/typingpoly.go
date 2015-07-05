package trompe

func (state *State) subst(ty Type, env *TyvarEnv) (Type, error) {
	LogTypingf("==> BEGIN subst: %s", ReprOfType(ty))
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return state.subst(desc.Type, env)
		} else {
			return desc, nil
		}
	case *TypeVar:
		Debugf("type var")
		if ty1, ok := env.Get(desc.Name); ok {
			Debugf("found %s -> %s, %s", desc.Name, ReprOfType(ty1), env)
			return ty1, nil
		} else {
			Debugf("not found %s", desc.Name)
			return ty, nil
		}
	case *TypePoly:
		newEnv := NewTyvarEnv()
		vlen := len(desc.Tyvars)
		newTyvars := make([]string, vlen)
		for i, _ := range desc.Tyvars {
			//newName := TyvarNames[i+vlen]
			newName := TyvarNames[i+len(newEnv.Map)]
			newTyvars[i] = newName
			newEnv.SetVar(newName, newName)
		}
		ty2, err := state.subst(desc.Type, newEnv)
		if err != nil {
			return nil, err
		}
		ty3, err := state.subst(ty2, env)
		if err != nil {
			return nil, err
		}
		return TPoly(newTyvars, ty3), nil
	case *TypeApp:
		if tyfun, ok := desc.Tycon.(*TyconTyFun); ok {
			Debugf("type app subst %d %d ", len(tyfun.Tyvars), len(desc.Args))
			Debugf("subst tyfun.Type = %s", ReprOfType(tyfun.Type))
			funEnv := NewTyvarEnv()
			funEnv.SetLists(tyfun.Tyvars, desc.Args)
			funTy, err := state.subst(tyfun.Type, funEnv)
			if err != nil {
				return nil, err
			}
			retTy, err := state.subst(funTy, env)
			if err != nil {
				return nil, err
			}
			Debugf("<== subst tyfun: %s => %s",
				ReprOfType(ty), ReprOfType(retTy))
			return retTy, nil
		} else {
			Debugf("not tyfun = %s", ty)
			env = env.Copy()
			args := make([]Type, len(desc.Args))
			for i, arg := range desc.Args {
				Debugf("subst arg %d: %s", i, ReprOfType(arg))
				arg, err := state.subst(arg, env)
				if err != nil {
					return nil, err
				}
				args[i] = arg
			}
			return TApp(desc.Tycon, args), nil
		}
	default:
		Panicf("notimpl %s", ty)
		return nil, nil
	}
}

func (inf *inferer) generalize(env *TyvarEnv, ty Type) (Type, error) {
	LogTypingf("==> BEGIN generalize: %s", ReprOfType(ty))
	env = env.Copy()
	mvars := FreeMetavarsOfType(ty)
	tyvars := make([]string, len(mvars))
	for i, mvar := range mvars {
		tyvar := env.NewTyvar()
		tyvars[i] = tyvar
		env.SetVar(mvar, tyvar)
	}
	ty = inf.generalizeFreeMetas(env, ty)
	ret := TPoly(tyvars, ty)
	LogTypingf("<== END generalize: %s => %s", ReprOfType(ty), ReprOfType(ret))
	return ret, nil
}

func (inf *inferer) generalizeFreeMetas(env *TyvarEnv, ty Type) Type {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return inf.generalizeFreeMetas(env, desc.Type)
		} else {
			if va, ok := env.Get(desc.Name); ok {
				return va
			} else {
				Panicf("cannot find Var(%s) for Meta(%s)", desc.Name, desc.Name)
				return nil
			}
		}
	case *TypePoly:
		return TPoly(desc.Tyvars, inf.generalizeFreeMetas(env, desc.Type))
	case *TypeApp:
		args := make([]Type, len(desc.Args))
		for i, arg := range desc.Args {
			args[i] = inf.generalizeFreeMetas(env, arg)
		}
		return TApp(desc.Tycon, args)
	default:
		return ty
	}
}

func (inf *inferer) instantiate(ty Type) (Type, error) {
	Debugf("==> BEGIN instantiate: %s", ReprOfType(ty))
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return inf.instantiate(desc.Type)
		} else {
			Debugf("<== END instantiate")
			return ty, nil
		}
	case *TypePoly:
		env := NewTyvarEnv()
		Debugf("instantiate tyvars: %s", desc.Tyvars)
		for _, name := range desc.Tyvars {
			env.Set(name, inf.env.NewMeta())
		}
		Debugf("instantiate with env: %s", env)
		ty2, err := inf.state.subst(desc.Type, env)
		if err != nil {
			return nil, err
		}
		Debugf("<== END instantiate: %s => %s", ReprOfType(ty), ReprOfType(ty2))
		return ty2, nil
	default:
		Debugf("<== END instantiate: nothing changed")
		return ty, nil
	}
}
