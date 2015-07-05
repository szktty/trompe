package trompe

import (
	"fmt"
	"os"
	"path/filepath"
)

func IntfFileOfModule(name string) string {
	return fmt.Sprintf("%s.tmi", name)
}

func (state *State) LoadIntf(mname string) (*Module, error) {
	Verbosef("search interface file of %s", mname)
	fname := IntfFileOfModule(mname)
	fpath, err := state.findIntfFile(fname)
	if err != nil {
		Verbosef("interface file not found %s", fpath)
		return nil, err
	}
	Verbosef("interface file found %s", fpath)
	return state.parseIntfFile(fpath, mname)
}

func (state *State) findIntfFile(fname string) (string, error) {
	for _, base := range state.SearchPath {
		fpath := filepath.Join(base, fname)
		if _, err := os.Stat(fpath); err == nil {
			return fpath, nil
		}
	}
	return "", fmt.Errorf("interface file of %s is not found", fname)
}

func (state *State) parseIntfFile(fpath string, mname string) (*Module, error) {
	scn := NewLexerFromFile(fpath)
	Debugf("begin parse")
	node, err := Parse(scn)
	if err != nil {
		return nil, err
	}
	Debugf("end parse")
	env, _ := state.NewTopEnv()
	env = NewEnv(env)
	if err := state.parseIntf(env, node); err != nil {
		return nil, err
	} else {
		m := NewModule(mname)
		m.MergeEnv(env)
		return m, nil
	}
}

func (state *State) parseIntf(env *Env, node *Node) error {
	switch desc := node.Desc.(type) {
	case *ProgramNode:
		for _, item := range desc.Items {
			if err := state.parseIntf(env, item); err != nil {
				return err
			}
		}
	case *ExtNode:
		// TODO: 多相型の場合は TyFun と Poly でラップする
		ty, err := state.parseType(env, desc.TypeExp)
		if err != nil {
			return err
		}
		env.Venv[desc.Name] = ty
		env.Prims[desc.Name] = desc.Prim
		Debugf("external %s : %s = \"%s\"", desc.Name, StringOfType(ty), desc.Prim)
		Debugf("ty = %s", ty)
	default:
		panic(fmt.Errorf("notimpl %s", desc))
	}
	return nil
}

func (state *State) parseType(env *Env, node *Node) (Type, error) {
	switch desc := node.Desc.(type) {
	case *TypeVarNode:
		return TVar(desc.Name.Value), nil
	case *TypePolyNode:
		ty, err := state.parseType(env, desc.App)
		if err != nil {
			return nil, err
		}
		return TPoly(StringsOfWords(desc.Vars), ty), nil
	case *TypeArrowNode:
		ty, err := state.parseType(env, desc.Left)
		if err != nil {
			return nil, err
		}
		args := make([]Type, 1)
		args[0] = ty
		args, err = state.parseTypeArrow(env, desc.Right, args)
		if err != nil {
			return nil, err
		}
		return &TypeApp{Tycon: &TyconArrow{}, Args: args}, nil
	case *TypeConstrAppNode:
		// TODO: Exps
		constr := desc.Constr.Desc.(*TypeConstrNode)
		Debugf("constr exps = %s", constr)
		path := constr.NamePath()
		mod, err := state.FindModuleOfPath(path.ModulePath())
		if err != nil {
			return nil, err
		}
		if ty, ok := mod.FindType(path.Name); ok {
			return ty, nil
		} else {
			return nil, fmt.Errorf("type constructor `%s' is not found", path.String())
		}
	case *TypeConstrNode:
		// TODO: Path
		name := desc.Name.Value
		if ty, ok := env.FindType(name); ok {
			return ty, nil
		} else {
			return nil, RuntimeErrorf(node.Loc,
				"type constructor `%s' is not found", name)
		}
	case *TypeParamConstrNode:
		args, err := state.parseTypes(env, desc.Exps)
		if err != nil {
			return nil, err
		}

		// TODO
		ty, err := state.FindTypeOfConstr(env, desc.Constr)
		if err != nil {
			return nil, err
		}
		vars := ty.(*TypePoly).Tyvars
		Debugf("vars = %s", vars)
		Debugf("args = %s", args)
		if len(vars) != len(args) {
			panic(fmt.Errorf("number of variables (%d) is not equal to one of arguments (%d)", len(vars), len(args)))
		}
		ty, err = state.subst(ty, TyvarEnvOfLists(vars, args))
		if err != nil {
			return nil, err
		}

		Debugf("subst result = %s, %s", ReprOfType(ty), ty)
		return ty, nil
	default:
		panic(fmt.Sprintf("notimpl %s", node))
	}
	return nil, nil
}

func (state *State) FindTypeOfConstr(env *Env, node *Node) (Type, error) {
	desc := node.Desc.(*TypeConstrNode)
	name := desc.Name.Value
	if ty, ok := env.FindType(name); ok {
		return ty, nil
	} else {
		return nil, RuntimeErrorf(node.Loc,
			"type constructor `%s' is not found", name)
	}
}

func (state *State) parseTypes(env *Env, nodes []*Node) ([]Type, error) {
	tys := make([]Type, 0)
	for _, node := range nodes {
		ty, err := state.parseType(env, node)
		if err != nil {
			return nil, err
		}
		tys = append(tys, ty)
	}
	return tys, nil
}

func (state *State) parseTypeArrow(env *Env, node *Node, args []Type) ([]Type, error) {
	switch desc := node.Desc.(type) {
	case *TypeArrowNode:
		ty, err := state.parseType(env, desc.Left)
		if err != nil {
			return nil, err
		}
		args, err = state.parseTypeArrow(env, desc.Right, append(args, ty))
		if err != nil {
			return nil, err
		}
	default:
		ty, err := state.parseType(env, node)
		if err != nil {
			return nil, err
		}
		args = append(args, ty)
	}
	return args, nil
}
