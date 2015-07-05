package trompe

import (
	"fmt"
)

type Context struct {
	State  *State
	Parent *Context
	Module *Module
	Block  *BlockClosure
	Pc     int
	Ptr    int
	Stack  []Value
}

func (ctx *Context) Const(i int) Value {
	return ctx.Block.Code.Consts[i]
}

func (ctx *Context) Top() Value {
	return ctx.Stack[ctx.Ptr]
}

func (ctx *Context) TopArray() []Value {
	if v, ok := ctx.Top().([]Value); ok {
		return v
	} else {
		panic("not array")
	}
}

func (ctx *Context) Push(obj Value) {
	if v, ok := obj.(int); ok {
		obj = int64(v)
	}
	ctx.Ptr++
	ctx.Stack[ctx.Ptr] = obj
}

func (ctx *Context) PushLocal(i int) {
	ctx.Push(ctx.Stack[i])
}

func (ctx *Context) PushLocalIndirect(b int) {
	ary := ctx.Stack[b/16].([]Value)
	ctx.Push(ary[b%16])
}

func (ctx *Context) PushConst(i int) {
	ctx.Push(ctx.Block.Code.Consts[i])
}

func (ctx *Context) PushUnit() {
	ctx.Push(UnitValue)
}

func (ctx *Context) PushNil() {
	ctx.Push(NilValue)
}

func (ctx *Context) PushBool(v bool) {
	ctx.Push(v)
}

func (ctx *Context) Pop() Value {
	ctx.Ptr--
	return ctx.Stack[ctx.Ptr+1]
}

func (ctx *Context) SwapPop() {
	ctx.Stack[ctx.Ptr-1] = ctx.Stack[ctx.Ptr]
	ctx.Ptr--
}

func (ctx *Context) PopBool() bool {
	return ctx.Pop().(bool)
}

func (ctx *Context) PopInt() int64 {
	if v, ok := ctx.Pop().(int64); ok {
		return v
	} else {
		panic("not int")
	}
}

func (ctx *Context) PopArray() []Value {
	if v, ok := ctx.Pop().([]Value); ok {
		return v
	} else {
		panic("not array")
	}
}

func (ctx *Context) PopValues(n int) []Value {
	vs := make([]Value, n)
	for i := 0; i < n; i++ {
		vs[n-1-i] = ctx.Pop()
	}
	return vs
}

func (ctx *Context) StorePopLocal(i int) {
	ctx.Stack[i] = ctx.Pop()
}

func (ctx *Context) StorePopLocalIndirect(b int) {
	ary := ctx.Stack[b/16].([]Value)
	ary[b%16] = ctx.Pop()
}

func (ctx *Context) ValidateJumpBackward(i int) {
	if ctx.Pc > int(i) && ctx.Block.Code.Bytes[i] != OpLoopHead {
		panic("target of a backwards branch must be OpLoopHead")
	}
}

func (ctx *Context) CompareValues() (ComparisonResult, error) {
	v2 := ctx.Pop()
	v1 := ctx.Pop()
	if res, ok := CompareValues(v1, v2); ok {
		return res, nil
	} else {
		Debugf("compare %s %s", v1, v2)
		return 0, fmt.Errorf("both values must be same types")
	}
}

// TODO
func (ctx *Context) Eval(blk *BlockClosure, args []Value) (Value, error) {
	return nil, fmt.Errorf("error")
}

func (state *State) Exec(mod *Module, parent *Context, block *BlockClosure, args []Value) (Value, error) {
	code := block.Code
	var ctx Context
	ctx.State = state
	ctx.Parent = parent
	ctx.Module = mod
	ctx.Block = block
	ctx.Ptr = code.NumArgs + code.NumTemps - 1
	// +2 will be for frameless functions (not yet implemented)
	ctx.Stack = make([]Value, code.NumLocals()+code.FrameSize+2)
	Debugf("eval start %d on %d", ctx.Ptr, len(ctx.Stack))

	logExec := LogGroupEnabled(LogGroupExec)
	var bcDescs map[int]string
	if logExec {
		_, descs := code.StringsOfBytes(false)
		bcDescs = descs
		Logf(LogGroupExec, "exec %s", code.Name)
	}

	if len(args) == code.NumArgs {
		for i := 0; i < code.NumArgs; i++ {
			ctx.Stack[i] = args[i]
			if logExec {
				Logf(LogGroupExec, "arg %d: %s", i, StringOfValue(args[i]))
			}
		}
	} else {
		return nil, fmt.Errorf("arity not equal")
	}

loophead:
	if ctx.Pc+1 > len(code.Bytes) {
		panic("return instruction must be needed")
	}
	bc := int(code.Bytes[ctx.Pc])
	if logExec {
		Debugf("%d %s", ctx.Pc, bcDescs[ctx.Pc])
		if ctx.Ptr >= 0 {
			for i := ctx.Ptr; i >= 0; i-- {
				Debugf("    at %d: %s", i, StringOfValue(ctx.Stack[i]))
			}
		}
	}

	ctx.Pc++
	switch {
	case OpLoadLocal0 <= bc && bc <= OpMaxLoadLocal:
		ctx.PushLocal(bc - OpLoadLocal0)
	case OpLoadGlobal0 <= bc && bc <= OpMaxLoadGlobal:
		name := ctx.Const(bc - OpLoadGlobal0).(string)
		v, ok := ctx.Module.FindFieldValue(name)
		if !ok {
			Panicf("field %s is not found", name)
		}
		ctx.Push(v)
	case OpLoadConst0 <= bc && bc <= OpMaxLoadConst:
		ctx.PushConst(bc - OpLoadConst0)
	case OpLoadValue0 <= bc && bc <= OpMaxLoadValue:
		path := ctx.Const(bc - OpLoadValue0).(*NamePath)
		f, err := state.FindFieldValueOfPath(path)
		if err != nil {
			return nil, err
		}
		ctx.Push(f)
	case OpLoadInt0 <= bc && bc <= OpMaxLoadInt:
		ctx.Push(bc - OpLoadInt0)
	case OpLoadIndirect0 <= bc && bc <= OpMaxLoadIndirect:
		ary := ctx.TopArray()
		ctx.Push(ary[bc-OpLoadIndirect0])
	case OpPopLoadIndirect0 <= bc && bc <= OpMaxPopLoadIndirect:
		ary := ctx.PopArray()
		ctx.Push(ary[bc-OpPopLoadIndirect0])
	case OpStorePopLocal0 <= bc && bc <= OpMaxStorePopLocal:
		ctx.StorePopLocal(bc - OpStorePopLocal0)
	case OpStorePopGlobal0 <= bc && bc <= OpMaxStorePopGlobal:
		name := ctx.Const(bc - OpStorePopGlobal0).(string)
		ctx.Module.SetFieldValue(name, ctx.Pop())
	case OpShortBranchFalse0 <= bc && bc <= OpMaxShortBranchFalse:
		if ctx.Pop() == false {
			ctx.Pc += bc - OpShortBranchFalse0
		}
	case OpShortJump0 <= bc && bc <= OpMaxShortJump:
		ctx.Pc += bc - OpShortJump0 + 1
	case OpApply1 <= bc && bc <= OpMaxApply:
		args := ctx.PopValues(bc - OpApply1 + 1)
		f := ctx.Pop()
		Debugf("apply %s", StringOfValue(f))
		ret, err := state.Apply(&ctx, f, args)
		if err != nil {
			return nil, err
		} else {
			Debugf("return")
			ctx.Push(ret)
		}
	case OpApplyDirect1_0 <= bc && bc <= OpMaxApplyDirect1:
		args := ctx.PopValues(1)
		f := ctx.Stack[bc-OpApplyDirect1_0]
		Debugf("apply %s", StringOfValue(f))
		ret, err := state.Apply(&ctx, f, args)
		if err != nil {
			return nil, err
		} else {
			Debugf("return")
			ctx.Push(ret)
		}
	default:
		switch bc {
		case OpLoadUnit:
			ctx.PushUnit()
		case OpLoadNil:
			ctx.PushNil()
		case OpLoadHead:
			list := ctx.Top().(*List)
			ctx.Push(list.Head)
		case OpPopLoadTail:
			list := ctx.Pop().(*List)
			ctx.Push(list.Tail)
		case OpLoadTrue:
			ctx.PushBool(true)
		case OpLoadFalse:
			ctx.PushBool(false)
		case OpPop:
			ctx.Pop()
		case OpSwapPop:
			ctx.SwapPop()
		case OpDup:
			ctx.Push(ctx.Top())
		case OpReturn:
			return ctx.Top(), nil
		case OpReturnUnit:
			return UnitValue, nil
		case OpReturnTrue:
			return true, nil
		case OpReturnFalse:
			return false, nil
		case OpAdd:
			ctx.Push(ctx.PopInt() + ctx.PopInt())
		case OpSub:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.Push(v1 - v2)
		case OpMul:
			ctx.Push(ctx.PopInt() * ctx.PopInt())
		case OpDiv:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.Push(v1 / v2)
		case OpMod:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.Push(v1 % v2)
		case OpEq:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			ctx.PushBool(res == OrderedSame)
		case OpNe:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			ctx.PushBool(res != OrderedSame)
		case OpLe:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			ctx.PushBool(res != OrderedAscending)
		case OpLt:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			ctx.PushBool(res == OrderedDescending)
		case OpGe:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			ctx.PushBool(res != OrderedDescending)
		case OpGt:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			ctx.PushBool(res == OrderedAscending)
		case OpEqInts:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.PushBool(v1 == v2)
		case OpNeInts:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.PushBool(v1 != v2)
		case OpLtInts:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.PushBool(v1 < v2)
		case OpLeInts:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.PushBool(v1 <= v2)
		case OpGtInts:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.PushBool(v1 > v2)
		case OpGeInts:
			v2 := ctx.PopInt()
			v1 := ctx.PopInt()
			ctx.PushBool(v1 >= v2)
		case OpAdd1:
			v := ctx.PopInt()
			ctx.Push(v + 1)
		case OpSub1:
			v := ctx.PopInt()
			ctx.Push(v - 1)
		case OpCountValues:
			switch vs := ctx.Top().(type) {
			case []Value:
				ctx.Push(len(vs))
			case *List:
				ctx.Push(vs.Length())
			default:
				panic("the top value of the stack must be an array or list")
			}
		case OpLoopHead:
			goto loophead
		default:
			goto twoByteCode
		}
	}
	goto loophead

twoByteCode:
	b1 := int(code.Bytes[ctx.Pc])
	ctx.Pc++
	switch {
	case OpLongJump0 <= bc && bc <= OpMaxLongJump:
		ctx.Pc += (bc-OpLongJump4)*256 + b1
	case OpLongBranchTrue0 <= bc && bc <= OpMaxLongBranchTrue:
		if ctx.PopBool() {
			ctx.Pc += (bc-OpLongBranchTrue0)*256 + b1
		}
	case OpLongBranchFalse0 <= bc && bc <= OpMaxLongBranchFalse:
		if !ctx.PopBool() {
			ctx.Pc += (bc-OpLongBranchFalse0)*256 + b1
		}
	default:
		switch bc {
		case OpLoadLocalIndirect:
			ctx.PushLocalIndirect(b1)
		case OpXLoadIndirect:
			ary := ctx.TopArray()
			ctx.Push(ary[b1])
		case OpXStorePopLocal:
			ctx.StorePopLocal(b1)
		case OpStorePopLocalIndirect:
			ctx.StorePopLocalIndirect(b1)
		case OpXStorePopGlobal:
			name := ctx.Const(b1).(string)
			ctx.Module.SetFieldValue(name, ctx.Pop())
		case OpXLoadInt:
			ctx.Push(b1)
		case OpBranchNe:
			v2 := ctx.Pop()
			v1 := ctx.Pop()
			res, ok := CompareValues(v1, v2)
			if !ok || (ok && res != OrderedSame) {
				ctx.Pc += b1 + 1
			}
		case OpBranchNeSizes:
			v2 := len(ctx.PopArray())
			v1 := len(ctx.PopArray())
			if v1 != v2 {
				ctx.Pc += b1 + 1
			}
		case OpCreateArray:
			ary := make([]Value, b1+1)
			ctx.Push(ary)
		case OpConsArray:
			ary := ctx.PopValues(b1 + 1)
			ctx.Push(ary)
		case OpConsList:
			list := NewListFromArray(ctx.PopValues(b1 + 1))
			ctx.Push(list)
		case OpCopyValues:
			if b1 != len(block.Copied) {
				Panicf("number of copied values %d of the block is not equal to bytecode %d", len(block.Copied), b1)
			}
			for i := 0; i < b1; i++ {
				ctx.Push(block.Copied[i])
			}
		case OpFullBlock:
			code := ctx.Const(b1).(*CompiledCode)
			blk := NewBlockClosure(code)
			blk.Context = &ctx
			ctx.Push(blk)
		default:
			goto threeByteCode
		}
	}
	goto loophead

threeByteCode:
	b2 := int(code.Bytes[ctx.Pc])
	ctx.Pc++
	switch bc {
	case OpXXLoadInt:
		ctx.Push(b1*256 + b2)
	case OpCopyingBlock, OpFullCopyingBlock:
		code := ctx.Const(b1).(*CompiledCode)
		blk := NewBlockClosure(code)
		blk.Copied = ctx.PopValues(b2)
		ctx.Push(blk)
		if bc == OpFullCopyingBlock {
			blk.Context = &ctx
		}
	case OpPrimitive:
		key := ctx.Const(b1).(string)
		args := ctx.PopValues(b2)
		if f, ok := PrimMap[key]; ok {
			ret, err := f(state, parent, args)
			if err != nil {
				return nil, err
			}
			ctx.Push(ret)
		} else {
			return nil, fmt.Errorf("primitive `%s' is not found", key)
		}
	default:
		panic("notimpl op")
	}
	goto loophead
}
