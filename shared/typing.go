package trompe

import (
	"fmt"
	"strconv"
)

type inferer struct {
	state      *State
	env        *Env
	block      *TypedBlockNode
	unifyCount int // for debug
}

func (state *State) Typing(env *Env, pnode *Node) (*TypedNode, error) {
	Debugf("==> BEGIN typing")
	var inf inferer
	inf.state = state
	inf.env = env
	tnode, err := inf.infer(pnode)
	Debugf("<== END typing")
	return tnode, err
}

func (inf *inferer) beginEnv() *Env {
	inf.env = NewEnv(inf.env)
	return inf.env
}

func (inf *inferer) endEnv() {
	inf.env = inf.env.Parent
}

func (inf *inferer) infer(pnode *Node) (*TypedNode, error) {
	if pnode == nil {
		return nil, nil
	}

	LogTypingf("==> BEGIN infer type at %s", pnode.Loc.StartString())
	defer LogTypingf("==> END infer type at %s", pnode.Loc.StartString())

	switch desc := pnode.Desc.(type) {
	case *UnitNode:
		return NewTypedNode(pnode.Loc, TUnit, &TypedUnitNode{}), nil

	case *BoolNode:
		return NewTypedNode(pnode.Loc, TBool, &TypedBoolNode{desc.Value}), nil

	case *IntNode:
		n := &TypedIntNode{Value: desc.Value}
		v, err := strconv.ParseInt(desc.Value, 10, 64)
		if err != nil {
			n.Decimal = true
		} else {
			n.SmallInt = v
		}
		return NewTypedNode(pnode.Loc, TInt, n), nil

	case *FloatNode:
		return NewTypedNode(pnode.Loc, TFloat, &TypedFloatNode{desc.Value}), nil

	case *CharNode:
		return NewTypedNode(pnode.Loc, TChar, &TypedCharNode{desc.Value}), nil

	case *StringNode:
		return NewTypedNode(pnode.Loc, TString, &TypedStringNode{desc.Value}), nil

	case *ProgramNode:
		titems, err := inf.inferNodes(desc.Items)
		if err == nil {
			return NewTypedNode(pnode.Loc, TUnit,
				&TypedProgramNode{titems}), nil
		} else {
			return nil, err
		}

		/*
			case *ImportNode:
				path := desc.Path.Copy()
				name, ok := path.Unshift()
				if !ok {
					// TODO
				}
				sig, ok := inf.state.LoadInterface(name)
				if !ok {
					// TODO
				}
				inf.env.AddSig(sig)
				// TODO: モジュールパスを渡した場合
		*/

	case *LetNode:
		if desc.Rec {
			for _, n := range desc.Bindings {
				switch bind := n.Desc.(type) {
				case *LetBindingNode:
					switch ptn := bind.Ptn.Desc.(type) {
					case *PtnIdentNode:
						LogTypingf("let rec pattern %s", ptn.Name)
						inf.env.AddVarType(ptn.Name, inf.env.NewMeta())
					}
				case *BlockNode:
					f := bind.Name.Desc.(*IdentNode)
					LogTypingf("let rec %s", f.Name)
					inf.env.AddVarType(f.Name, inf.env.NewMeta())
				}
			}
		}

		tbs, err := inf.inferNodes(desc.Bindings)
		if err != nil {
			return nil, err
		}

		inf.beginEnv()
		tbody, err := inf.infer(desc.Body)
		if err != nil {
			return nil, err
		}
		inf.endEnv()

		var ty Type
		if tbody != nil {
			ty = tbody.Type
		} else {
			ty = TUnit
		}
		return NewTypedNode(pnode.Loc, ty,
			&TypedLetNode{Public: desc.Public, Rec: desc.Rec,
				Bindings: tbs, Body: tbody}), nil

	case *LetBindingNode:
		tptn, err := inf.inferPattern(desc.Ptn, inf.env.NewMeta())
		if err != nil {
			return nil, err
		}
		tbody, err := inf.infer(desc.Body)
		if err != nil {
			return nil, err
		}
		if err := inf.unify(desc.Ptn.Loc, tptn.Type, tbody.Type); err != nil {
			return nil, err
		}

		return NewTypedNode(pnode.Loc, tptn.Type,
			&TypedLetBindingNode{Ptn: tptn, Body: tbody}), nil

	case *BlockNode:
		name := desc.Name.Desc.(*IdentNode).Name

		if desc.Rec {
			inf.beginEnv()
		} else {
			inf.beginEnv()
		}

		// params
		tparams, err := inf.inferParams(desc.Params)
		if err != nil {
			return nil, err
		}

		// check labeled arguments
		var labels []string
		for _, tparam := range tparams {
			if _, ok := tparam.Desc.(*TypedLabelParamNode); ok {
				labels = make([]string, len(tparams))
				for i, tparam := range tparams {
					var name string
					if desc, ok := tparam.Desc.(*TypedLabelParamNode); ok {
						name = desc.Name
					}
					labels[i] = name
				}
				break
			}
		}

		// body
		tbody, err := inf.infer(desc.Body)
		if err != nil {
			return nil, err
		}
		args := TypesOfTypedNodes(tparams)
		args = append(args, tbody.Type)
		gargs := make([]Type, len(args))
		tvenv := NewTyvarEnv()
		for i, arg := range args {
			g, err := inf.generalize(tvenv, arg)
			if err != nil {
				panic(err)
			}
			gargs[i] = g
		}

		// type constructor
		var tycon Tycon = TcArrow
		if labels != nil {
			tycon = &TyconKeyArrow{Keywords: labels}
		}

		// type
		ty := TApp(tycon, gargs)
		if len(tvenv.Tyvars) > 0 {
			ty = TPoly(tvenv.Tyvars, ty)
		}
		Debugf("before generalize: %s", ReprOfType(ty))
		gty, err := inf.generalize(NewTyvarEnv(), ty)
		if err != nil {
			Panicf("generalize `%s' failed: %s", name, ReprOfType(ty))
		}
		ty = gty
		Debugf("block %s : %s", name, ReprOfType(ty))

		// node
		tname := NewTypedNode(desc.Name.Loc, ty,
			&TypedIdentNode{Name: name})
		tnode := NewTypedNode(pnode.Loc, ty,
			&TypedBlockNode{Name: tname, Params: tparams, Body: tbody})
		inf.endEnv()
		inf.env.AddVarType(desc.Name.Desc.(*IdentNode).Name, ty)
		return tnode, nil

	case *IfNode:
		tcond, err := inf.infer(desc.Cond)
		if err != nil {
			return nil, err
		}

		if err := inf.unify(desc.Cond.Loc, TBool, tcond.Type); err != nil {
			return nil, err
		}

		ttrue, err := inf.infer(desc.True)
		if err != nil {
			return nil, err
		}

		var tfalse *TypedNode
		if desc.False == nil {
			if err := inf.unify(desc.True.Loc, TUnit, ttrue.Type); err != nil {
				return nil, err
			} else {
				tfalse = NewTypedNode(pnode.Loc, TUnit, &TypedUnitNode{})
			}
		} else {
			tfalse, err = inf.infer(desc.False)
			if err != nil {
				return nil, err
			}
			if err := inf.unify(desc.False.Loc, ttrue.Type, tfalse.Type); err != nil {
				return nil, err
			}
		}
		return NewTypedNode(pnode.Loc, ttrue.Type,
			&TypedIfNode{Cond: tcond, True: ttrue, False: tfalse}), nil

	case *ForNode:
		inf.beginEnv()
		tinit, err := inf.infer(desc.Init)
		if err != nil {
			return nil, err
		}
		if err := inf.unify(desc.Init.Loc, TInt, tinit.Type); err != nil {
			return nil, err
		}
		name := desc.Name.Desc.(*IdentNode).Name
		inf.env.AddVarType(name, TInt)
		tname := NewTypedNode(desc.Name.Loc, TInt, &TypedIdentNode{Name: name})
		tlimit, err := inf.infer(desc.Limit)
		if err != nil {
			return nil, err
		}
		if err := inf.unify(desc.Limit.Loc, TInt, tlimit.Type); err != nil {
			return nil, err
		}
		tbody, err := inf.infer(desc.Body)
		if err != nil {
			return nil, err
		}
		if err := inf.unify(desc.Body.Loc, TUnit, tbody.Type); err != nil {
			return nil, err
		}
		inf.endEnv()
		return NewTypedNode(pnode.Loc, TUnit,
			&TypedForNode{Name: tname, Init: tinit, Dir: desc.Dir,
				Limit: tlimit, Body: tbody}), nil

	case *SeqExpNode:
		texps := make([]*TypedNode, len(desc.Exps))
		for i, exp := range desc.Exps {
			texp, err := inf.infer(exp)
			if err != nil {
				return nil, err
			}
			if i+1 < len(desc.Exps) {
				if err := inf.unify(texp.Loc, TUnit, texp.Type); err != nil {
					return nil, err
				}
			}
			texps[i] = texp
		}
		ty := texps[len(texps)-1].Type
		return NewTypedNode(pnode.Loc, ty, &TypedSeqExpNode{Exps: texps}), nil

	case *CaseNode:
		texp, err := inf.infer(desc.Exp)
		if err != nil {
			return nil, err
		}
		ty := texp.Type
		ret := inf.env.NewMeta()
		tms, err := inf.inferMatchList(desc.Match, ty, ret)
		if err != nil {
			return nil, err
		}
		return NewTypedNode(pnode.Loc, ret,
			&TypedCaseNode{Exp: texp, Match: tms}), nil

	case *FunctionNode:
		inf.beginEnv()
		cond := inf.env.NewMeta()
		ret := inf.env.NewMeta()
		tms, err := inf.inferMatchList(desc.Match, cond, ret)
		if err != nil {
			return nil, err
		}
		gcond, err := inf.generalize(NewTyvarEnv(), cond)
		if err != nil {
			return nil, err
		}
		gret, err := inf.generalize(NewTyvarEnv(), ret)
		if err != nil {
			return nil, err
		}
		inf.endEnv()
		ty := TApp(TcArrow, TArgs(gcond, gret))
		Debugf("function : %s", ReprOfType(ty))
		return NewTypedNode(pnode.Loc, ty,
			&TypedFunctionNode{Match: tms}), nil

	case *FunNode:
		// TODO: keyword
		env := inf.beginEnv()
		match := desc.MultiMatch.Desc.(*MultiMatchNode)
		tparams := make([]*TypedNode, len(match.Params))
		args := make([]Type, len(match.Params))
		for i, param := range match.Params {
			tparam, err := inf.inferPattern(param, env.NewMeta())
			if err != nil {
				return nil, err
			}
			tparams[i] = tparam
			args[i] = tparam.Type
		}

		var tcond *TypedNode
		if match.Cond != nil {
			t, err := inf.infer(match.Cond)
			if err != nil {
				return nil, err
			}
			if err := inf.unify(match.Cond.Loc, TBool, t.Type); err != nil {
				return nil, err
			}
			tcond = t
		}

		tbody, err := inf.infer(match.Body)
		if err != nil {
			return nil, err
		}
		inf.endEnv()

		args = append(args, tbody.Type)
		arrow := TApp(TcArrow, args)
		tvenv := NewTyvarEnv()
		arrow, err = inf.generalize(tvenv, arrow)
		if err != nil {
			return nil, err
		}
		Debugf("generalize : %s", ReprOfType(arrow))
		tmatch := NewTypedNode(pnode.Loc, tbody.Type,
			&TypedMultiMatchNode{Params: tparams, Cond: tcond, Body: tbody})
		return NewTypedNode(pnode.Loc, arrow,
			&TypedFunNode{MultiMatch: tmatch}), nil

	case *AddNode:
		return inf.inferBinexp(pnode, &TypedAddNode{}, desc.Left, desc.Right)

	case *SubNode:
		return inf.inferBinexp(pnode, &TypedSubNode{}, desc.Left, desc.Right)

	case *MulNode:
		return inf.inferBinexp(pnode, &TypedMulNode{}, desc.Left, desc.Right)

	case *DivNode:
		return inf.inferBinexp(pnode, &TypedDivNode{}, desc.Left, desc.Right)

	case *ModNode:
		return inf.inferBinexp(pnode, &TypedModNode{}, desc.Left, desc.Right)

	case *FAddNode:
		return inf.inferBinexp(pnode, &TypedFAddNode{}, desc.Left, desc.Right)

	case *FSubNode:
		return inf.inferBinexp(pnode, &TypedFSubNode{}, desc.Left, desc.Right)

	case *FMulNode:
		return inf.inferBinexp(pnode, &TypedFMulNode{}, desc.Left, desc.Right)

	case *FDivNode:
		return inf.inferBinexp(pnode, &TypedFDivNode{}, desc.Left, desc.Right)

	case *EqNode:
		return inf.inferBinexp(pnode, &TypedEqNode{}, desc.Left, desc.Right)

	case *NeNode:
		return inf.inferBinexp(pnode, &TypedNeNode{}, desc.Left, desc.Right)

	case *LtNode:
		return inf.inferBinexp(pnode, &TypedLtNode{}, desc.Left, desc.Right)

	case *LeNode:
		return inf.inferBinexp(pnode, &TypedLeNode{}, desc.Left, desc.Right)

	case *GtNode:
		return inf.inferBinexp(pnode, &TypedGtNode{}, desc.Left, desc.Right)

	case *GeNode:
		return inf.inferBinexp(pnode, &TypedGeNode{}, desc.Left, desc.Right)

	case *ValuePathNode:
		path := desc.NamePath()
		mod, err := inf.env.FindModuleOfPath(path)
		if err != nil {
			return nil, NewRuntimeError(pnode.Loc, err)
		}

		if ty, ok := mod.FindFieldType(path.Name); ok {
			return NewTypedNode(pnode.Loc, ty,
				&TypedValuePathNode{Path: path}), nil
		} else {
			return nil, RuntimeErrorf(pnode.Loc,
				"%s: Unbound value path `%s'",
				pnode.Loc.StartString(), path.String())
		}

	case *IdentNode:
		Debugf("find var %s", desc.Name)
		if ty, ok := inf.env.FindVarType(desc.Name); ok {
			Debugf("found var %s : %s", desc.Name, ReprOfType(ty))
			return NewTypedNode(pnode.Loc, ty,
				&TypedIdentNode{Name: desc.Name}), nil
		} else {
			return nil, RuntimeErrorf(pnode.Loc, "%s: Unbound value `%s'",
				pnode.Loc.StartString(), desc.Name)
		}

	case *KeywordNode:
		texp, err := inf.infer(desc.Exp)
		if err != nil {
			return nil, err
		}
		return NewTypedNode(pnode.Loc, texp.Type,
			&TypedKeywordNode{Keyword: desc.Keyword.Value, Exp: texp}), nil

	case *AppNode:
		texp, err := inf.infer(desc.Exp)
		if err != nil {
			return nil, err
		}

		expTy, err := inf.instantiate(texp.Type)
		if err != nil {
			return nil, err
		}

		// infer the types of the arguments
		targs := make([]*TypedNode, len(desc.Args))
		for i := 0; i < len(desc.Args); i++ {
			targ, err := inf.infer(desc.Args[i])
			if err != nil {
				return nil, err
			}
			argTy, err := inf.instantiate(targ.Type)
			if err != nil {
				return nil, err
			}
			targ.Type = argTy
			targs[i] = targ
		}

		// if the type of the function is not yet undefined
		if _, ok := expTy.(*TypeMeta); ok {
			fargs := make([]Type, len(targs))
			for i, targ := range targs {
				fargs[i] = targ.Type
			}
			ret := inf.env.NewMeta()
			fargs = append(fargs, ret)
			arrow := TApp(TcArrow, fargs)
			err := inf.unify(texp.Loc, arrow, expTy)
			if err != nil {
				return nil, err
			}
			return NewTypedNode(pnode.Loc, ret,
				&TypedAppNode{Exp: texp, Args: targs}), nil
		}

		funApp, ok := TypeArrowOfType(expTy)
		if !ok {
			return nil, NotFunError(pnode.Loc, expTy)

		}

		// check format type
		for i := 0; i < len(desc.Args); i++ {
			if param, ok := funApp.Args[i].(*TypeApp); ok {
				if arg, ok := desc.Args[i].Desc.(*StringNode); ok &&
					param.Tycon.TyconTag() == TyconTagFormat {
					ty, err := translateFormat(funApp, i, arg.Value)
					if err != nil {
						return nil, err
					}
					funApp = ty
					break
				}
			}
		}

		var ret Type
		var funTy Type = funApp
		for i := 0; i < len(desc.Args); i++ {
			if funTy == nil {
				return nil, TooManyArgsError(texp.Loc, expTy)
			}
			LogTypingf("apply: %s", ReprOfType(funTy))

			targ := targs[i]
			funTy1, err := inf.instantiate(funTy)
			if err != nil {
				return nil, err
			}
			if _, ok := TypeArrowOfType(funTy1); !ok {
				return nil, NotFunError(targ.Loc, funTy1)
			}

			var head, tail Type
			switch funApp.Tycon.TyconTag() {
			case TyconTagArrow:
				if key, ok := targ.Desc.(*TypedKeywordNode); ok {
					return nil, RuntimeErrorf(targ.Loc,
						"The function applied to this argument has type `%s'. This argument cannot be applied with keyword `%s'",
						StringOfType(funTy1), key.Keyword)
				}
				head1, tail1, ok := PartialArrow(funTy1)
				if !ok {
					return nil, TooManyArgsError(texp.Loc, expTy)
				}
				head = head1
				tail = tail1
			case TyconTagKeyArrow:
				key, ok := targ.Desc.(*TypedKeywordNode)
				if !ok {
					return nil, RuntimeErrorf(targ.Loc,
						"labels were omitted in the application of this function:\n       %s",
						StringOfType(funTy1))
				}
				Debugf("key arraw = %s, %s", ReprOfType(funTy1), key.Keyword)
				head1, tail1, ok := PartialKeyArrow(key.Keyword, funTy1)
				if !ok {
					return nil, TooManyArgsError(texp.Loc, expTy)
				}
				head = head1
				tail = tail1
			default:
				panic("not arrow tycon")
			}

			head1, err := inf.instantiate(head)
			if err != nil {
				return nil, err
			}

			err = inf.unify(targ.Loc, head1, targ.Type)
			if err != nil {
				return nil, err
			}

			funTy = tail
			if funTy == nil {
				LogTypingf("return %d: %s", i, ReprOfType(head1))
			} else {
				LogTypingf("arg %d: %s => %s", i, ReprOfType(head1), ReprOfType(tail))
			}
			ret = funTy
			LogTypingf("apply return: %s", ret)
		}

		return NewTypedNode(pnode.Loc, ret,
			&TypedAppNode{Exp: texp, Args: targs}), nil

	case *ConstrAppNode:
		tconstr, err := inf.infer(desc.Constr)
		if err != nil {
			return nil, err
		}

		Debugf("constr type %s", ReprOfType(tconstr.Type))
		constrDesc, ok := tconstr.Desc.(*TypedConstrNode)
		if !ok {
			panic("error")
		}

		path := constrDesc.Path
		Debugf("find %s", path.String())
		var conTy Type
		if path.HasBase() {
			mod, err := inf.env.FindModuleOfPath(path)
			if err != nil {
				return nil, NewRuntimeError(pnode.Loc, err)
			}

			ty, ok := mod.FindType(path.Name)
			if !ok {
				return nil, RuntimeErrorf(pnode.Loc, "Unbound module value %s", path.String())
			}
			conTy = ty
		} else {
			// constructor or exception
			if ty, ok := inf.env.FindType(path.Name); ok {
				conTy = ty
			} else if ty, ok := inf.env.FindExnType(path.Name); ok {
				conTy = ty
			} else {
				return nil, RuntimeErrorf(pnode.Loc, "Unbound constructor value %s", path)
			}
		}
		Debugf("constructor type = %s", ReprOfType(conTy))

		var texp *TypedNode
		if desc.Exp != nil {
			tnode, err := inf.infer(desc.Exp)
			if err != nil {
				return nil, err
			}
			texp = tnode
			Debugf("texp = %s", ReprOfType(texp.Type))
		}

		app1, ok := PickTypeApp(conTy)
		if !ok {
			panic("error")
		}
		numArgs1 := len(app1.Args)
		if texp == nil {
			if numArgs1 == 0 {
				return NewTypedNode(pnode.Loc, tconstr.Type,
					&TypedConstrAppNode{Constr: tconstr, Exp: texp}), nil
			} else {
				return nil, ConstrArityError(pnode.Loc, path.String(), numArgs1, 0)
			}
		} else {
			app2, ok := PickTypeApp(texp.Type)
			if !ok {
				panic("error")
			}
			var args2 []Type
			if _, ok := app2.Tycon.(*TyconTuple); ok {
				args2 = app2.Args
			} else {
				if t, ok := NewTypeTuple(app2); ok {
					args2 = t.(*TypeApp).Args
				} else {
					panic("creating tuple failed")
				}
			}
			numArgs2 := len(args2)
			if numArgs1 != numArgs2 {
				return nil, ConstrArityError(pnode.Loc,
					path.String(), numArgs1, numArgs2)
			} else {
				for i := 0; i < numArgs1; i++ {
					if err := inf.unify(pnode.Loc,
						app1.Args[i], args2[i]); err != nil {
						return nil, err
					}
				}
			}
			return NewTypedNode(pnode.Loc, tconstr.Type,
				&TypedConstrAppNode{Constr: tconstr, Exp: texp}), nil
		}

	case *ConstrNode:
		// TODO: Path
		name := desc.Name.Value
		path := NewNamePath()
		path.AddName(name)
		Debugf("constr path %s", path.String())
		ty, ok := inf.env.FindType(name)
		if !ok {
			ty, ok = inf.env.FindExnType(name)
			if !ok {
				return nil, RuntimeErrorf(pnode.Loc, "%s is not found", name)
			}
		}
		ty, err := inf.instantiate(ty)
		if err != nil {
			return nil, err
		}
		return NewTypedNode(pnode.Loc, ty,
			&TypedConstrNode{Path: path}), nil

	case *TupleNode:
		tcomps := make([]*TypedNode, len(desc.Comps))
		appArgs := make([]Type, len(desc.Comps))
		for i, elt := range desc.Comps {
			telt, err := inf.infer(elt)
			if err != nil {
				return nil, err
			}
			tcomps[i] = telt
			appArgs[i] = telt.Type
			Debugf("tuple %d: %s", i, ReprOfType(telt.Type))
		}
		ty := &TypeApp{Tycon: &TyconTuple{}, Args: appArgs}
		return NewTypedNode(pnode.Loc, ty,
			&TypedTupleNode{Comps: tcomps}), nil

	case *ListNode:
		telts := make([]*TypedNode, len(desc.Elts))
		var arg Type
		if len(desc.Elts) == 0 {
			// empty list => 'a list
			arg = inf.env.NewMeta()
		} else {
			for i, elt := range desc.Elts {
				telt, err := inf.infer(elt)
				if err != nil {
					return nil, err
				}
				if i == 0 {
					arg = telt.Type
				} else {
					err = inf.unify(elt.Loc, arg, telt.Type)
					if err != nil {
						return nil, err
					}
				}
				telts[i] = telt
			}
		}
		ty := &TypeApp{Tycon: &TyconList{}, Args: []Type{arg}}
		return NewTypedNode(pnode.Loc, ty, &TypedListNode{Elts: telts}), nil

	case *ListConsNode:
		thead, err := inf.infer(desc.Head)
		if err != nil {
			return nil, err
		}
		ttail, err := inf.infer(desc.Tail)
		if err != nil {
			return nil, err
		}
		theadL := TApp(TcList, TArgs(thead.Type))
		err = inf.unify(thead.Loc, theadL, ttail.Type)
		if err != nil {
			return nil, err
		}
		return NewTypedNode(pnode.Loc, ttail.Type,
			&TypedListConsNode{Head: thead, Tail: ttail}), nil

	case *ArrayNode:
		telts := make([]*TypedNode, len(desc.Elts))
		var arg Type
		if len(desc.Elts) == 0 {
			// empty array => 'a arry
			arg = inf.env.NewMeta()
		} else {
			for i, elt := range desc.Elts {
				telt, err := inf.infer(elt)
				if err != nil {
					return nil, err
				}
				if i == 0 {
					arg = telt.Type
				} else {
					err = inf.unify(elt.Loc, arg, telt.Type)
					if err != nil {
						return nil, err
					}
				}
				telts[i] = telt
			}
		}
		ty := &TypeApp{Tycon: &TyconArray{}, Args: []Type{arg}}
		return NewTypedNode(pnode.Loc, ty, &TypedArrayNode{Elts: telts}), nil

	case *ArrayAccessNode:
		tary, err := inf.infer(desc.Array)
		if err != nil {
			return nil, err
		}
		eltTy := inf.env.NewMeta()
		aryTy := &TypeApp{Tycon: &TyconArray{}, Args: []Type{eltTy}}
		err = inf.unify(tary.Loc, aryTy, tary.Type)
		if err != nil {
			return nil, err
		}

		tidx, err := inf.infer(desc.Index)
		if err != nil {
			return nil, err
		}
		err = inf.unify(tidx.Loc, TInt, tidx.Type)
		if err != nil {
			return nil, err
		}

		var ty Type
		var tset *TypedNode
		if desc.Set != nil {
			ty = TUnit
			tset, err = inf.infer(desc.Set)
			if err != nil {
				return nil, err
			}

			err = inf.unify(tset.Loc, eltTy, tset.Type)
			if err != nil {
				return nil, err
			}
		} else {
			ty = eltTy
		}
		return NewTypedNode(pnode.Loc, ty,
			&TypedArrayAccessNode{Array: tary, Index: tidx, Set: tset}), nil

	case *OptionNode:
		ty := TNone
		var tval *TypedNode
		if desc.Value != nil {
			tval, err := inf.infer(desc.Value)
			if err != nil {
				return nil, err
			}
			ty = NewTypeSome(tval.Type)
		}
		return NewTypedNode(pnode.Loc, ty,
			&TypedOptionNode{Value: tval}), nil

	default:
		panic(fmt.Sprintf("not impl at %s, %s", pnode.Loc.StartString(), pnode))
	}
	return nil, RuntimeErrorf(pnode.Loc, "internal error")
}

func (inf *inferer) inferNodes(pnodes []*Node) ([]*TypedNode, error) {
	tnodes := make([]*TypedNode, len(pnodes))
	for i, pnode := range pnodes {
		tnode, err := inf.infer(pnode)
		if err == nil {
			tnodes[i] = tnode
		} else {
			return nil, err
		}
	}
	return tnodes, nil
}

func (inf *inferer) inferMatch(match *Node, cond Type, ret Type) (*TypedNode, error) {
	m := match.Desc.(*MatchNode)
	var tptn, tcond *TypedNode
	var err error
	if !m.Ptn.isWildcard() {
		tptn, err = inf.inferPattern(m.Ptn, cond)
		if err != nil {
			return nil, err
		}
	}

	if m.Cond != nil {
		tcond, err := inf.infer(m.Cond)
		if err != nil {
			return nil, err
		}
		if err := inf.unify(tcond.Loc, TBool, tcond.Type); err != nil {
			return nil, err
		}
	}

	tbody, err := inf.infer(m.Body)
	if err != nil {
		return nil, err
	}
	if err := inf.unify(tbody.Loc, ret, tbody.Type); err != nil {
		return nil, err
	}

	return NewTypedNode(match.Loc, ret,
		&TypedMatchNode{Ptn: tptn, Cond: tcond, Body: tbody}), nil
}

func (inf *inferer) inferMatchList(match []*Node, cond Type, ret Type) ([]*TypedNode, error) {
	tms := make([]*TypedNode, len(match))
	for i, e := range match {
		tm, err := inf.inferMatch(e, cond, ret)
		if err != nil {
			return nil, err
		}
		tms[i] = tm
	}
	return tms, nil
}

func (inf *inferer) inferPattern(ptn *Node, expTy Type) (*TypedNode, error) {
	Debugf("ptn expected %s", ReprOfType(expTy))
	var ptnTy Type
	var tnode interface{}
	switch desc := ptn.Desc.(type) {
	case *WildcardNode:
		ptnTy = expTy
		tnode = &TypedWildcardNode{}

	case *PtnConstNode:
		tval, err := inf.infer(desc.Value)
		if err != nil {
			return nil, err
		}
		ptnTy = tval.Type
		tnode = &TypedPtnConstNode{Value: tval}

	case *PtnIdentNode:
		if _, ok := expTy.(*TypeMeta); ok {
			ptnTy = expTy
		} else {
			ptnTy = inf.env.NewMeta()
		}
		tnode = &TypedPtnIdentNode{Name: desc.Name}
		inf.env.AddVarType(desc.Name, ptnTy)

	case *PtnTupleNode:
		tcomps := make([]*TypedNode, len(desc.Comps))
		compTys := make([]Type, len(desc.Comps))
		for i, e := range desc.Comps {
			ty := inf.env.NewMeta()
			te, err := inf.inferPattern(e, ty)
			if err != nil {
				return nil, err
			}
			tcomps[i] = te
			compTys[i] = ty
			Debugf("ptn tuple %d: %s", i, ReprOfType(ty))
		}
		ptnTy = &TypeApp{Tycon: &TyconTuple{}, Args: compTys}
		tnode = &TypedPtnTupleNode{Comps: tcomps}

	case *PtnListNode:
		telts := make([]*TypedNode, len(desc.Elts))
		var arg Type
		if len(desc.Elts) == 0 {
			// empty list => 'a list
			arg = inf.env.NewMeta()
		} else {
			for i, e := range desc.Elts {
				ty := inf.env.NewMeta()
				telt, err := inf.inferPattern(e, ty)
				if err != nil {
					return nil, err
				}
				if i == 0 {
					arg = ty
				} else {
					err = inf.unify(e.Loc, arg, telt.Type)
					if err != nil {
						return nil, err
					}
				}
				telts[i] = telt
			}
		}
		ptnTy = &TypeApp{Tycon: &TyconList{}, Args: []Type{arg}}
		tnode = &TypedPtnListNode{Elts: telts}

	case *PtnListConsNode:
		headTy := inf.env.NewMeta()
		thead, err := inf.inferPattern(desc.Head, headTy)
		if err != nil {
			return nil, err
		}

		ttail, err := inf.inferPattern(desc.Tail, expTy)
		if err != nil {
			return nil, err
		}

		theadL := TApp(TcList, TArgs(thead.Type))
		if err := inf.unify(ptn.Loc, theadL, ttail.Type); err != nil {
			return nil, err
		}

		ptnTy = ttail.Type
		tnode = &TypedPtnListConsNode{Head: thead, Tail: ttail}

	case *PtnSomeNode:
		tptn, err := inf.inferPattern(desc.Ptn, inf.env.NewMeta())
		if err != nil {
			return nil, err
		}
		ptnTy = NewTypeSome(tptn.Type)

	default:
		Panicf("notimpl: %s", ptn)
	}

	if ptnTy == nil {
		panic("type of pattern is not found")
	}
	if err := inf.unify(ptn.Loc, expTy, ptnTy); err != nil {
		return nil, err
	}
	return NewTypedNode(ptn.Loc, expTy, tnode), nil
}

func (inf *inferer) inferParams(params []*Node) ([]*TypedNode, error) {
	tparams := make([]*TypedNode, len(params))
	for i, param := range params {
		tparam, err := inf.inferParam(param)
		if err != nil {
			return nil, err
		}
		tparams[i] = tparam
	}
	return tparams, nil
}

func (inf *inferer) inferParam(param *Node) (*TypedNode, error) {
	var ty Type
	var tparam interface{}
	switch desc := param.Desc.(type) {
	case *PtnIdentNode:
		ty = inf.env.NewMeta()
		inf.env.AddVarType(desc.Name, ty)
		tparam = &TypedPtnIdentNode{Name: desc.Name}
	case *LabelParamNode:
		tptn, err := inf.inferParam(desc.Ptn)
		if err != nil {
			return nil, err
		}
		ty = tptn.Type
		tparam = &TypedLabelParamNode{Name: desc.Name.Value, Ptn: tptn}
	default:
		panic(fmt.Errorf("not impl %s", desc))
	}
	return NewTypedNode(param.Loc, ty, tparam), nil
}

func (inf *inferer) inferBinexp(pnode *Node, desc interface{}, left *Node, right *Node) (*TypedNode, error) {
	tleft, err := inf.infer(left)
	if err != nil {
		return nil, err
	}

	var ty Type
	switch desc.(type) {
	case *TypedAddNode, *TypedSubNode, *TypedMulNode, *TypedDivNode, *TypedModNode:
		ty = TInt
	case *TypedFAddNode, *TypedFSubNode, *TypedFMulNode, *TypedFDivNode:
		ty = TFloat
	case *TypedEqNode, *TypedNeNode, *TypedLtNode, *TypedLeNode,
		*TypedGtNode, *TypedGeNode:
		ty = tleft.Type
	default:
		panic("notimpl")
	}

	tright, err := inf.infer(right)
	if err != nil {
		return nil, err
	}
	if err := inf.unify(left.Loc, ty, tleft.Type); err != nil {
		return nil, err
	}
	if err := inf.unify(right.Loc, ty, tright.Type); err != nil {
		return nil, err
	}
	Debugf("binexp = %s, %s", ReprOfType(tleft.Type), ReprOfType(tright.Type))
	switch t := desc.(type) {
	case *TypedAddNode:
		t.Left = tleft
		t.Right = tright
	case *TypedSubNode:
		t.Left = tleft
		t.Right = tright
	case *TypedMulNode:
		t.Left = tleft
		t.Right = tright
	case *TypedDivNode:
		t.Left = tleft
		t.Right = tright
	case *TypedModNode:
		t.Left = tleft
		t.Right = tright
	case *TypedFAddNode:
		t.Left = tleft
		t.Right = tright
	case *TypedFSubNode:
		t.Left = tleft
		t.Right = tright
	case *TypedFMulNode:
		t.Left = tleft
		t.Right = tright
	case *TypedFDivNode:
		t.Left = tleft
		t.Right = tright
	case *TypedEqNode:
		t.Left = tleft
		t.Right = tright
		ty = TBool
	case *TypedNeNode:
		t.Left = tleft
		t.Right = tright
		ty = TBool
	case *TypedLtNode:
		t.Left = tleft
		t.Right = tright
		ty = TBool
	case *TypedLeNode:
		t.Left = tleft
		t.Right = tright
		ty = TBool
	case *TypedGtNode:
		t.Left = tleft
		t.Right = tright
		ty = TBool
	case *TypedGeNode:
		t.Left = tleft
		t.Right = tright
		ty = TBool
	default:
		panic("notimpl")
	}
	Debugf("binexp finish %s", ty)
	return NewTypedNode(pnode.Loc, ty, desc), nil
}

func (inf *inferer) parseFormat(ty Type, arg *TypedNode) (*Format, error) {
	sdesc, ok := arg.Desc.(*TypedStringNode)
	if !ok {
		return nil, RuntimeErrorf(arg.Loc,
			"This expression has `%s' type, but must be is expected to a format string as string literal",
			StringOfType(arg.Type))
	}
	f, err := NewFormat(sdesc.Value)
	if err != nil {
		return nil, err
	}
	err = f.SetType(ty)
	if err != nil {
		return nil, err
	} else {
		return f, nil
	}
}
