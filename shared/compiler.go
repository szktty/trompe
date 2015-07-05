package trompe

import (
	"fmt"
	"math"
	"strconv"
)

type Compiler struct {
	State   *State
	Context *LoadingContext
	Azer    *Analyzer
	Code    *CompiledCode
}

func (state *State) Compile(ctx *LoadingContext, node *TypedNode) (*CompiledCode, error) {
	if _, ok := node.Desc.(*TypedProgramNode); !ok {
		return nil, fmt.Errorf("node %p must be TypedProgramNode", node)
	}
	c := NewCompiler(state, ctx)
	c.Compile(nil, node)
	return c.Code, nil
}

func NewCompiler(state *State, ctx *LoadingContext) *Compiler {
	c := &Compiler{State: state, Context: ctx}
	/*
		var c compiler
		c.init(state, ctx)
		bld := newBuilder(&c, nil)
		moduleVar := "*module*"
		bld.addLocal(moduleVar)
		bld.createStructure(node.Loc)
		bld.storePopLocal(node.Loc, moduleVar)
		c.compile(bld, node)
		bld.pushLocal(node.Loc, moduleVar)
		bld.return_()
		code := c.finish(bld)
		return code, nil
	*/
	return c
}

func (c *Compiler) CompiledCode(bld *IRBuilder) *CompiledCode {
	code := bld.CompiledCode()
	code.File = c.Context.File
	return code
}

func (c *Compiler) Finish(bld *IRBuilder) *CompiledCode {
	// TODO: もしかしてこの関数は不要か？
	/*
		if bld.Parent == nil {
			bld.ReturnUnit()
		} else {
			bld.return_()
		}
	*/

	code := c.CompiledCode(bld)
	Debugf("compiled code:\n%s", code.String())
	return code
}

func (c *Compiler) Compile(bld *IRBuilder, node *TypedNode) {
	switch desc := node.Desc.(type) {
	case *TypedProgramNode:
		// TODO: builder で full import すべき？
		//c.FullImport(state.Pervasives)

		azer, err := c.State.AnalyzeBlock(node)
		if err != nil {
			panic(err)
		}
		Debugf("azer = %s", azer)

		c.Azer = azer
		bld := NewIRBuilder(nil, azer.ScopeOfNode(node))
		Debugf("prog = %s, %s", bld.Scope.Asis, bld.Scope)

		// shared array on the top level is needed?
		numSharedSlots := bld.Scope.NumSharedSlots()
		if numSharedSlots > 0 {
			bld.PutCreateArray(nil, numSharedSlots)
			bld.StorePopLocal(nil, bld.Scope.Shared.Index)
		}

		for _, item := range desc.Items {
			c.Compile(bld, item)
		}
		bld.PutReturnUnit(nil)

		c.Code = c.CompiledCode(bld)
		LogCompilingf("program code:\n%s", c.Code)

	case *TypedLetNode:
		/*
			if desc.Rec {
				bld.rec = make(map[string]*TypedNode)
				for _, n := range desc.Bindings {
					switch bind := n.Desc.(type) {
					case *TypedLetBindingNode:
						switch ptn := bind.Ptn.Desc.(type) {
						case *PtnIdentNode:
							bld.rec[ptn.Name] = n
						}
					case *TypedBlockNode:
						f := bind.Name.Desc.(*TypedIdentNode)
						bld.rec[f.Name] = n
					}
				}
				for name, node := range bld.rec {
					bld.addLocalRef(name)
					bld.createRef(node.Loc)
					bld.storePopLocal(node.Loc, name)
				}
			}
		*/

		for _, bind := range desc.Bindings {
			c.Compile(bld, bind)
		}
		if desc.Body != nil {
			c.Compile(bld, desc.Body)
		}

	case *TypedLetBindingNode:
		c.Compile(bld, desc.Body)
		c.CompileMatchPattern(bld, desc.Ptn)

	case *TypedBlockNode:
		name, _ := desc.Name.Name()
		inner := NewIRBuilder(bld, c.Azer.ScopeOfNode(node))
		inner.Name = name

		// outer variables
		Debugf("block = %s", inner.Scope.Asis)
		numRefs := len(inner.Scope.Refs)
		if numRefs > 0 {
			Debugf("push copied %d", numRefs)
			inner.PushCopied(nil, numRefs)
		}

		// shared variables (array, recursive functions)
		numSharedSlots := inner.Scope.NumSharedSlots()
		if numSharedSlots > 0 {
			inner.PutCreateArray(nil, numSharedSlots)
			inner.StorePopLocal(nil, inner.Scope.Shared.Index)
		}

		c.Compile(inner, desc.Body)
		inner.PutReturn(nil)
		code := c.CompiledCode(inner)
		LogCompilingf("block code:\n%s", code)
		bld = inner.Parent
		if inner.Scope.IsClean() {
			// clean block
			i := bld.AddConst(NewBlockClosure(code))
			bld.PushConst(node.Loc, i)

		} else if inner.Scope.IsFull() && !inner.Scope.IsCopying() {
			codeIdx := bld.AddConst(code)
			bld.MakeFullBlock(nil, codeIdx)

		} else if inner.Scope.IsCopying() || inner.Scope.IsFullCopying() {
			// push variables to copy
			codeIdx := bld.AddConst(code)
			for _, ref := range inner.Scope.Refs {
				bld.PushLocal(nil, ref.Ref.Index)
			}

			// make block
			numCopied := inner.Scope.NumCopied()
			Debugf("make block %d for %s", numCopied, inner.Scope)
			if inner.Scope.IsFullCopying() {
				//bld.PushLocal(nil, bld.Scope.Shared.Index)
				bld.MakeFullCopyingBlock(nil, codeIdx, numCopied)
			} else {
				bld.MakeCopyingBlock(nil, codeIdx, numCopied)
			}
		} else {
			panic("block must be one of clean, copying, full or full copying")
		}

		if bld.Scope.Outer == nil {
			i := bld.AddConst(name)
			bld.StorePopGlobal(nil, i)
		} else {
			// TODO: 再帰の場合、 array にセット
			if !bld.StorePopVar(node.Loc, name) {
				Panicf("local or global %s is not found", name)
			}
		}

	case *TypedIdentNode:
		if BeginsWithLowerCase(desc.Name) {
			if bld.PushVar(node.Loc, desc.Name) {
				/*
					} else if path, ok := bld.findExtVal(desc.Name); ok {
						Debugf("add static ref %s", path)
						bld.pushStatic(node.Loc, path)
				*/
			} else if path, ok := bld.Scope.Bindings[desc.Name]; ok {
				i := bld.AddConst(path)
				bld.PushValue(nil, i)
			} else {
				// TODO
				Panicf("notimpl ident %s", desc.Name)
				/*
					if mod, ok := node.Scope.ModuleOfVar(desc.Name); ok {
						// module variable
						path := mod.Path().AddName(desc.Name)
						i := bld.AddConst(path)
						bld.PushValue(node.Loc, i)
					} else {
						Panicf("variable `%s' is not in the variable scope", desc.Name)
					}
				*/
			}
		} else {
			// TODO: withuppercase
			panic("notimpl titlecase")
		}

	case *TypedSeqExpNode:
		for i, exp := range desc.Exps {
			c.Compile(bld, exp)
			if i+1 < len(desc.Exps) {
				bld.Pop(exp.Loc)
			}
		}
		if !bld.IsTop() && len(desc.Exps) > 1 {
			bld.PushUnit(nil)
		}

	case *TypedIfNode:
		c.Compile(bld, desc.Cond)
		condFalse := bld.NewLabel(nil)
		bld.PutBranchFalse(nil, condFalse)
		c.Compile(bld, desc.True)
		if desc.False != nil {
			trueTo := bld.NewLabel(nil)
			bld.PutJump(nil, trueTo)
			bld.AddInstr(condFalse)
			c.Compile(bld, desc.False)
			bld.AddInstr(trueTo)
		} else {
			bld.AddInstr(condFalse)
			bld.PushUnit(nil)
		}

	case *TypedForNode:
		name := desc.Name.NameExn()
		c.Compile(bld, desc.Init)
		bld.StorePopVar(nil, name)

		lh := bld.NewLabel(nil)
		end := bld.NewLabel(nil)
		bld.AddInstr(lh)
		bld.PutLoopHead(nil)
		bld.PushVar(nil, name)
		c.Compile(bld, desc.Limit)
		switch desc.Dir {
		case ForDirTo:
			bld.PutLe(nil)
		case ForDirDownTo:
			bld.PutGe(nil)
		default:
			panic("error")
		}
		bld.PutBranchFalse(nil, end)

		c.Compile(bld, desc.Body)
		bld.Pop(nil)
		bld.PushVar(nil, name)
		switch desc.Dir {
		case ForDirTo:
			bld.PutAdd1(nil)
		case ForDirDownTo:
			bld.PutSub1(nil)
		default:
			panic("error")
		}
		bld.StorePopVar(nil, name)
		bld.PutJump(nil, lh)

		bld.AddInstr(end)

	case *TypedCaseNode:
		c.Compile(bld, desc.Exp)
		c.CompileMatch(bld, desc.Match)

	case *TypedAppNode:
		direct := -1
		if id, ok := desc.Exp.Desc.(*TypedIdentNode); ok {
			if i, ok := bld.LocalIndex(id.Name); ok && i < math.MaxUint8 {
				direct = i
				Debugf("direct = %d, %s", i, id.Name)
			}
		}

		if direct < 0 {
			c.Compile(bld, desc.Exp)
		}
		for _, arg := range desc.Args {
			c.Compile(bld, arg)
		}
		if direct < 0 {
			bld.Apply(node.Loc, len(desc.Args))
		} else {
			bld.ApplyDirect(node.Loc, direct, len(desc.Args))
		}

	case *TypedAddNode:
		add1 := false
		if exp, ok := desc.Left.Desc.(*TypedIntNode); ok {
			v, err := strconv.ParseInt(exp.Value, 10, 64)
			if err == nil && v == 1 {
				add1 = true
				c.Compile(bld, desc.Right)
				bld.PutAdd1(node.Loc)
			}
		} else if exp, ok := desc.Right.Desc.(*TypedIntNode); ok {
			v, err := strconv.ParseInt(exp.Value, 10, 64)
			if err == nil && v == 1 {
				add1 = true
				c.Compile(bld, desc.Left)
				bld.PutAdd1(node.Loc)
			}
		}

		if !add1 {
			c.Compile(bld, desc.Left)
			c.Compile(bld, desc.Right)
			bld.PutAdd(node.Loc)
		}

	case *TypedSubNode:
		sub1 := false
		if exp, ok := desc.Left.Desc.(*TypedIntNode); ok {
			v, err := strconv.ParseInt(exp.Value, 10, 64)
			if err == nil && v == 1 {
				sub1 = true
				c.Compile(bld, desc.Right)
				bld.PutSub1(node.Loc)
			}
		} else if exp, ok := desc.Right.Desc.(*TypedIntNode); ok {
			v, err := strconv.ParseInt(exp.Value, 10, 64)
			if err == nil && v == 1 {
				sub1 = true
				c.Compile(bld, desc.Left)
				bld.PutSub1(node.Loc)
			}
		}

		if !sub1 {
			c.Compile(bld, desc.Left)
			c.Compile(bld, desc.Right)
			bld.PutSub(node.Loc)
		}

	case *TypedMulNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		bld.PutMul(node.Loc)

	case *TypedDivNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		bld.PutDiv(node.Loc)

		/*
			case *TypedPowNode:
				c.Compile(bld, desc.Left)
				c.Compile(bld, desc.Right)
				bld.PutPow(node.Loc)
		*/

	case *TypedModNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		bld.PutMod(node.Loc)

	case *TypedEqNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		if IsTypeIntOfTypes(desc.Left.Type, desc.Right.Type) {
			bld.PutEqInts(node.Loc)
		} else {
			bld.PutEq(node.Loc)
		}

	case *TypedNeNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		if IsTypeIntOfTypes(desc.Left.Type, desc.Right.Type) {
			bld.PutNeInts(node.Loc)
		} else {
			bld.PutNe(node.Loc)
		}

	case *TypedLtNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		if IsTypeIntOfTypes(desc.Left.Type, desc.Right.Type) {
			bld.PutLtInts(node.Loc)
		} else {
			bld.PutLt(node.Loc)
		}

	case *TypedLeNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		if IsTypeIntOfTypes(desc.Left.Type, desc.Right.Type) {
			bld.PutLeInts(node.Loc)
		} else {
			bld.PutLe(node.Loc)
		}

	case *TypedGtNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		if IsTypeIntOfTypes(desc.Left.Type, desc.Right.Type) {
			bld.PutGtInts(node.Loc)
		} else {
			bld.PutGt(node.Loc)
		}

	case *TypedGeNode:
		c.Compile(bld, desc.Left)
		c.Compile(bld, desc.Right)
		if IsTypeIntOfTypes(desc.Left.Type, desc.Right.Type) {
			bld.PutGeInts(node.Loc)
		} else {
			bld.PutGe(node.Loc)
		}

	case *TypedUnitNode:
		bld.PushUnit(node.Loc)

	case *TypedBoolNode:
		bld.PushBool(node.Loc, desc.Value)

	case *TypedIntNode:
		bld.PushIntOfString(node.Loc, desc.Value)

	case *TypedStringNode:
		bld.PushConstValue(node.Loc, desc.Value)

	case *TypedListNode:
		if len(desc.Elts) == 0 {
			bld.PushNil(node.Loc)
		} else if ary, ok := ConstValueArray(desc.Elts); ok {
			bld.PushConstValue(nil, NewListFromArray(ary))
		} else {
			for _, e := range desc.Elts {
				c.Compile(bld, e)
			}
			bld.PutConsList(nil, len(desc.Elts))
		}

	case *TypedArrayNode:
		if len(desc.Elts) == 0 {
			bld.PushConstValue(nil, make([]Value, 0))
		} else if ary, ok := ConstValueArray(desc.Elts); ok {
			bld.PushConstValue(nil, ary)
		} else {
			for _, e := range desc.Elts {
				c.Compile(bld, e)
			}
			bld.PutConsArray(nil, len(desc.Elts))
		}

	case *TypedArrayAccessNode:
		c.Compile(bld, desc.Array)
		c.Compile(bld, desc.Index)
		if desc.Set != nil {
			c.Compile(bld, desc.Set)
			bld.PutPrimitive(nil, "array_set", 3)
		} else {
			bld.PutPrimitive(nil, "array_get", 2)
		}

	case *TypedTupleNode:
		for _, comp := range desc.Comps {
			c.Compile(bld, comp)
		}
		bld.PutConsArray(nil, len(desc.Comps))

	case *TypedValuePathNode:
		i := bld.AddConst(desc.Path)
		bld.PushValue(node.Loc, i)

	default:
		Panicf("compiler notimpl %s", node)
	}
}

func IsTypeIntOfTypes(ty1 Type, ty2 Type) bool {
	tycon1, ok1 := TyconTagOfTypeApp(ty1)
	tycon2, ok2 := TyconTagOfTypeApp(ty2)
	return ok1 && ok2 && tycon1 == TyconTagInt && tycon2 == TyconTagInt
}

func ConstValueArray(nodes []*TypedNode) ([]Value, bool) {
	var ary []Value
	for _, node := range nodes {
		var v Value
		switch desc := node.Desc.(type) {
		case *TypedUnitNode:
			v = UnitValue
		case *TypedBoolNode:
			v = desc.Value
		case *TypedIntNode:
			if desc.Decimal {
				return nil, false
			} else {
				v = desc.SmallInt
			}
		case *TypedCharNode:
			v = desc.Value
		case *TypedStringNode:
			v = desc.Value
		case *TypedTupleNode:
			if ary, ok := ConstValueArray(desc.Comps); ok {
				v = ary
			} else {
				return nil, false
			}
		case *TypedListNode:
			if len(desc.Elts) == 0 {
				v = NilValue
			} else {
				return nil, false
			}
		default:
			return nil, false
		}
		ary = append(ary, v)
	}
	return ary, true
}

func (c *Compiler) CompileMatch(bld *IRBuilder, matches []*TypedNode) {
	toEnd := bld.NewLabel(nil)
	for i, match := range matches {
		last := i+1 == len(matches)
		desc := match.Desc.(*TypedMatchNode)
		toBody := bld.NewLabel(desc.Ptn.Loc)
		toNext := bld.NewLabel(desc.Ptn.Loc)

		if desc.Cond != nil {
			c.Compile(bld, desc.Cond)
			bld.PutBranchFalse(nil, toNext)
		}
		if desc.Ptn != nil {
			bld.Dup(desc.Ptn.Loc)
			c.CompileMatchPattern(bld, desc.Ptn)
			if !last {
				bld.PutBranchFalse(nil, toNext)
			} else {
				bld.Pop(nil)
			}
		}

		bld.AddInstr(toBody)
		c.Compile(bld, desc.Body)
		if !last {
			bld.PutJump(nil, toEnd)
			bld.AddInstr(toNext)
		}
	}
	bld.AddInstr(toEnd)
	bld.SwapPop(nil)
}

func (c *Compiler) CompileMatchPattern(bld *IRBuilder, ptn *TypedNode) {
	// matching pattern process must be push a bool as to whether
	// matching success or failure onto the stack
	switch desc := ptn.Desc.(type) {
	case *TypedWildcardNode:
		bld.Pop(ptn.Loc)
		bld.PushBool(ptn.Loc, true)

	case *TypedPtnConstNode:
		c.Compile(bld, desc.Value)
		bld.PutEq(nil)

	case *TypedPtnIdentNode:
		if BeginsWithLowerCase(desc.Name) {
			bld.StorePopVar(ptn.Loc, desc.Name)
		} else {
			// title-case
			panic("notimpl")
		}
		bld.PushBool(ptn.Loc, true)

	case *TypedPtnTupleNode:
		// compare length of the tuples
		fail := bld.NewLabel(nil)
		succ := bld.NewLabel(nil)
		bld.PutCountValues(nil)
		bld.PushInt(nil, len(desc.Comps))
		bld.PutBranchNe(nil, fail)

		// compare each components
		for i, e := range desc.Comps {
			bld.PushIndirect(e.Loc, i)
			c.CompileMatchPattern(bld, e)
			bld.PutBranchFalse(nil, fail)
		}

		bld.PushBool(ptn.Loc, true)
		bld.PutJump(nil, succ)
		bld.AddInstr(fail)
		bld.PushBool(ptn.Loc, false)
		bld.AddInstr(succ)
		bld.SwapPop(nil)

	case *TypedPtnListNode:
		// compare length of the lists
		fail := bld.NewLabel(nil)
		succ := bld.NewLabel(nil)
		bld.PutCountValues(nil)
		bld.PushInt(nil, len(desc.Elts))
		bld.PutBranchNe(nil, fail)

		// compare each elements
		if len(desc.Elts) > 0 {
			bld.Dup(nil) // first element
			for _, e := range desc.Elts {
				// car
				bld.PushHead(e.Loc)
				c.CompileMatchPattern(bld, e)
				bld.PutBranchFalse(nil, fail)

				// cdr
				bld.PopPushTail(e.Loc)
			}
			bld.Pop(nil)
		}
		bld.PushBool(ptn.Loc, true)
		bld.PutJump(nil, succ)
		bld.AddInstr(fail)
		bld.PushBool(ptn.Loc, false)
		bld.AddInstr(succ)
		bld.SwapPop(nil)

	case *TypedPtnListConsNode:
		// check length of the list
		fail := bld.NewLabel(nil)
		succ := bld.NewLabel(nil)
		bld.PutCountValues(nil)
		bld.PushInt(nil, 1)
		bld.PutGe(nil)
		bld.PutBranchFalse(nil, fail)

		// head
		bld.PushHead(ptn.Loc)
		c.CompileMatchPattern(bld, desc.Head)
		bld.PutBranchFalse(nil, fail)

		// tail
		bld.PopPushTail(nil)
		c.CompileMatchPattern(bld, desc.Tail)
		bld.PutJump(nil, succ)

		bld.AddInstr(fail)
		bld.Pop(nil)
		bld.PushBool(ptn.Loc, false)
		bld.AddInstr(succ)
		bld.SwapPop(nil)

	default:
		Panicf("compiler match pattern notimpl %s", ptn)
	}
}

/*

func (c *compiler) open(mname string) {
	if sig, ok := c.state.FindSig(mname); ok {
		c.imports = append(c.imports, sig)
	} else {
		panic(fmt.Sprintf("module signature of %s not found", mname))
	}
}

func (c *compiler) findSignature(name string) (*Signature, bool) {
	for _, sig := range c.imports {
		if sig, ok := sig.Subsig(name); ok {
			return sig, true
		}
	}
	return nil, false
}

func (c *compiler) compile(bld *builder, node *TypedNode) {
	//Debugf("compile %s", node)
	switch desc := node.Desc.(type) {
	case *TypedProgramNode:
		for _, item := range desc.Items {
			c.compile(bld, item)
		}

	case *TypedLetNode:
		if desc.Rec {
			bld.rec = make(map[string]*TypedNode)
			for _, n := range desc.Bindings {
				switch bind := n.Desc.(type) {
				case *TypedLetBindingNode:
					switch ptn := bind.Ptn.Desc.(type) {
					case *PtnIdentNode:
						bld.rec[ptn.Name] = n
					}
				case *TypedBlockNode:
					f := bind.Name.Desc.(*TypedIdentNode)
					bld.rec[f.Name] = n
				}
			}
			for name, node := range bld.rec {
				bld.addLocalRef(name)
				bld.createRef(node.Loc)
				bld.storePopLocal(node.Loc, name)
			}
		}

		for _, bnd := range desc.Bindings {
			c.compile(bld, bnd)
		}
		if desc.Body != nil {
			c.compile(bld, desc.Body)
		}
		bld.rec = nil

	case *TypedLetBindingNode:
		switch ptn := desc.Ptn.Desc.(type) {
		case *TypedPtnIdentNode:
			if _, ok := bld.isRecRef(ptn.Name); ok {
				Debugf("rec ref ok ok")
				// TODO: 呼ばれてない？
				bld.addLocalRef(ptn.Name)
				c.compile(bld, desc.Body)
				if !bld.storePopLocalIndirect(desc.Body.Loc, ptn.Name) {
					panic("failed")
				}
			} else {
				c.compile(bld, desc.Body)
				if ptn.Name == "_" {
					bld.pop(desc.Ptn.Loc)
				} else {
					bld.addLocal(ptn.Name)
					bld.storePopLocal(desc.Body.Loc, ptn.Name)
				}
			}

		case *TypedPtnConstNode:
			c.compile(bld, desc.Body)

		default:
			panic("not impl")
		}

	case *TypedBlockNode:
		var name string
		ref := false
		if desc.Name != nil {
			name = desc.Name.Desc.(*TypedIdentNode).Name
			if _, ok := bld.isRecRef(name); ok {
				bld.addLocalRef(name)
				ref = true
			} else {
				bld.addLocal(name)
			}
		}
		bld = newBuilder(c, bld)

		for _, param := range desc.Params {
			switch ptn := param.Desc.(type) {
			case *TypedPtnIdentNode:
				Debugf("add ptn param %s", ptn.Name)
				bld.addArg(ptn.Name)
			default:
				panic(fmt.Errorf("notimpl %s", param))
			}
		}

		c.compile(bld, desc.Body)
		code := c.finish(bld)
		inner := bld
		bld = bld.parent
		if code.Outers > 0 {
			for _, l := range inner.locals {
				if inner.isOuter(l) {
					bld.pushLocal(node.Loc, l.name)
				}
			}
			i := bld.addConst(code)
			bld.copyingBlock(node.Loc, uint8(i), uint8(code.Outers))
			if ref {
				bld.storePopLocalIndirect(node.Loc, name)
			} else {
				bld.storePopLocal(node.Loc, name)
			}
		} else {
			bld.pushConst(node.Loc, NewBlockClosure(code))
			bld.storePopLocal(node.Loc, name)
		}

	case *TypedIfNode:
		c.compile(bld, desc.Cond)
		br := bld.beginBranchForward(desc.Cond.Loc)
		c.compile(bld, desc.True)
		j := bld.beginJumpForward(desc.True.Loc)
		bld.endBranchForward(br, false)
		c.compile(bld, desc.False)
		bld.endJumpForward(j)

	case *TypedForNode:
		name := desc.Name.Desc.(*TypedIdentNode).Name
		bld.addLocal(name)
		c.compile(bld, desc.Init)
		bld.storePopLocal(desc.Init.Loc, name)

		fst := bld.beginJumpForward(desc.Limit.Loc)
		lh := bld.beginBranchBackward(desc.Body.Loc)
		bld.pushLocal(desc.Init.Loc, name)
		bld.pushAdd1(desc.Init.Loc)
		bld.storePopLocal(desc.Init.Loc, name)
		bld.endJumpForward(fst)
		c.compile(bld, desc.Body)
		bld.pop(desc.Body.Loc) // discard value of the expression
		bld.pushLocal(desc.Limit.Loc, name)
		c.compile(bld, desc.Limit)
		bld.binop(desc.Limit.Loc, OpLt)
		bld.endBranchBackward(desc.Limit.Loc, lh, true)

	case *TypedCaseNode:
		c.compile(bld, desc.Exp)
		c.compileMatch(bld, desc.Match)

	case *TypedAppNode:
		c.compile(bld, desc.Exp)
		for _, arg := range desc.Args {
			c.compile(bld, arg)
		}
		bld.apply(node.Loc, uint8(len(desc.Args)))

	case *TypedSeqExpNode:
		for i, exp := range desc.Exps {
			c.compile(bld, exp)
			if i+1 < len(desc.Exps) {
				bld.pop(exp.Loc)
			}
		}

	case *TypedAddNode:
		if exp, ok := desc.Left.Desc.(*TypedIntNode); ok && exp.Value == "1" {
			c.compile(bld, desc.Right)
			bld.pushAdd1(node.Loc)
		} else if exp, ok := desc.Right.Desc.(*TypedIntNode); ok && exp.Value == "1" {
			c.compile(bld, desc.Left)
			bld.pushAdd1(node.Loc)
		} else {
			c.compile(bld, desc.Left)
			c.compile(bld, desc.Right)
			bld.binop(node.Loc, OpAdd)
		}

	case *TypedSubNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpSub)

	case *TypedMulNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpMul)

	case *TypedDivNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpDiv)

	case *TypedModNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpMod)

	case *TypedFAddNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpFAdd)

	case *TypedFSubNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpFSub)

	case *TypedFMulNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpFMul)

	case *TypedFDivNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpFDiv)

	case *TypedEqNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpEq)

	case *TypedNeNode:
		c.compile(bld, desc.Left)
		c.compile(bld, desc.Right)
		bld.binop(node.Loc, OpNe)

	case *TypedIdentNode:
		if BeginsWithLowerCase(desc.Name) {
			if _, l, ok := bld.localIndex(desc.Name); ok {
				Debugf("registered ident node = %s, %b", desc.Name, l.ref)
				if l.ref {
					bld.pushLocalIndirect(node.Loc, desc.Name)
				} else {
					Debugf("push local %s", desc.Name)
					bld.pushLocal(node.Loc, desc.Name)
				}
			} else if _, ok := bld.isRecRef(desc.Name); ok {
				bld.addLocalRef(desc.Name)
				bld.pushLocalIndirect(node.Loc, desc.Name)
			} else if path, ok := bld.findExtVal(desc.Name); ok {
				Debugf("add static ref %s", path)
				bld.pushStatic(node.Loc, path)
			} else {
				Debugf("add local %s", desc.Name)
				bld.addLocal(desc.Name)
				bld.pushLocal(node.Loc, desc.Name)
			}
		} else {
			// TODO: withuppercase
			panic("notimpl titlecase")
		}

	case *TypedUnitNode:
		bld.pushUnit(node.Loc)

	case *TypedBoolNode:
		bld.pushBool(node.Loc, desc.Value)

	case *TypedIntNode:
		bld.pushIntOfString(node.Loc, desc.Value)

	case *TypedStringNode:
		bld.pushConst(node.Loc, desc.Value)

	case *TypedTupleNode:
		for _, e := range desc.Comps {
			c.compile(bld, e)
		}
		bld.createTuple(node.Loc, uint8(len(desc.Comps)))

	case *TypedListNode:
		for _, e := range desc.Elts {
			c.compile(bld, e)
		}
		bld.createList(node.Loc, uint8(len(desc.Elts)))

	default:
		panic("not impl")
	}
}

func (c *compiler) compileMatch(bld *builder, match []*TypedNode) {
	jmps := make([]uint16, 0)
	for _, mnode := range match {
		m := mnode.Desc.(*TypedMatchNode)
		Debugf("match = %s\n", m)

		var condBr, ptnBr uint16
		if m.Cond != nil {
			Debugf("cond = %s\n", m.Cond)
			bld.dup(mnode.Loc)
			c.compile(bld, m.Cond)
			condBr = bld.beginBranchForward(m.Cond.Loc)
		}

		if m.Ptn != nil {
			bld.dup(mnode.Loc)
			c.compileMatchPattern(bld, m.Ptn)
			ptnBr = bld.beginBranchForward(m.Ptn.Loc)
		}

		c.compile(bld, m.Body)
		jmps = append(jmps, bld.beginJumpForward(m.Body.Loc))
		if m.Ptn != nil {
			bld.endBranchForward(ptnBr, false)
		}
		if m.Cond != nil {
			bld.endBranchForward(condBr, false)
		}

		for _, jmp := range jmps {
			bld.endJumpForward(jmp)
		}
	}
}

func (c *compiler) compileMatchPattern(bld *builder, ptn *TypedNode) {
	Debugf("ptn = %s\n", ptn)
	switch desc := ptn.Desc.(type) {
	case *TypedPtnConstNode:
		c.compile(bld, desc.Value)
		bld.binop(ptn.Loc, OpEq)

	case *TypedPtnIdentNode:
		if desc.Name == "_" {
			// do nothing
		} else if BeginsWithLowerCase(desc.Name) {
			//bld.dup(ptn.Loc)
			bld.addLocal(desc.Name)
			bld.storePopLocal(ptn.Loc, desc.Name)
		} else {
			// title-case
			panic("notimpl")
		}
		bld.pushBool(ptn.Loc, true)

	case *TypedPtnTupleNode:
		brs := make([]uint16, len(desc.Comps))
		for i, e := range desc.Comps {
			bld.dup(e.Loc)
			bld.pushComp(e.Loc, uint8(i))
			c.compileMatchPattern(bld, e)
			brs[i] = bld.beginBranchForward(e.Loc)
		}
		bld.pushBool(ptn.Loc, true)
		jmp := bld.beginJumpForward(ptn.Loc)
		for _, br := range brs {
			bld.endBranchForward(br, false)
		}
		bld.pushBool(ptn.Loc, false)
		bld.endJumpForward(jmp)

	case *TypedPtnListNode:
		if len(desc.Elts) == 0 {
			bld.matchNil(ptn.Loc)
		} else {
			// compare length of the lists
			bld.dup(ptn.Loc)
			bld.pushLen(ptn.Loc)
			bld.pushInt(ptn.Loc, int64(len(desc.Elts)))
			lenBr := bld.beginBranchNe(ptn.Loc)

			brs := make([]uint16, len(desc.Elts))
			for i, e := range desc.Elts {
				bld.dup(e.Loc)
				bld.pushHead(e.Loc)
				c.compileMatchPattern(bld, e)
				brs[i] = bld.beginBranchForward(e.Loc)
				bld.pushTail(e.Loc)
			}
			bld.pop(ptn.Loc)
			bld.pushBool(ptn.Loc, true)
			jmp := bld.beginJumpForward(ptn.Loc)
			for _, br := range brs {
				bld.endBranchForward(br, false)
			}
			bld.pop(ptn.Loc)
			bld.endBranchNe(lenBr)
			bld.pushBool(ptn.Loc, false)
			bld.endJumpForward(jmp)
		}

	case *TypedPtnListConsNode:
		// compare length of the lists
		bld.dup(ptn.Loc)
		bld.pushLen(ptn.Loc)
		bld.pushInt(ptn.Loc, 0)
		lenBr := bld.beginBranchEq(ptn.Loc)

		// compare the head of the list
		bld.dup(ptn.Loc)
		bld.pushHead(ptn.Loc)
		c.compileMatchPattern(bld, desc.Head)
		headBr := bld.beginBranchForward(ptn.Loc)

		// compare the rest elements
		bld.pushTail(ptn.Loc)
		c.compileMatchPattern(bld, desc.Tail)
		jmp := bld.beginJumpForward(ptn.Loc)
		bld.endBranchEq(lenBr)
		bld.endBranchForward(headBr, false)
		bld.pushBool(ptn.Loc, false)
		bld.endJumpForward(jmp)

	default:
		panic("notimpl")
	}
}
*/
