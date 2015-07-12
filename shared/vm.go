package trompe

import (
	"fmt"
)

type Context struct {
	State   *State
	Parent  *Context
	Module  *Module
	Block   *BlockClosure
	Args    []Value
	Pc      int
	BasePtr int
}

func (state *State) NewContext(mod *Module, parent *Context,
	block *BlockClosure, args []Value) *Context {
	code := block.Code
	if len(args) != code.NumArgs {
		Panicf("arity must be %d", code.NumArgs)
	}
	return &Context{State: state, Module: mod, Parent: parent,
		Block: block, Args: args, BasePtr: state.Stack.Ptr}
}

func (ctx *Context) Const(i int) Value {
	return ctx.Block.Code.Consts[i]
}

func (ctx *Context) PushConst(i int) {
	ctx.State.Stack.Push(ctx.Block.Code.Consts[i])
}

func (ctx *Context) Return(v Value) {
	if ctx.Parent != nil {
		ctx.State.Stack.BasePtr = ctx.Parent.BasePtr
	} else {
		ctx.State.Stack.BasePtr = -1
	}
	ctx.State.Stack.Ptr = ctx.BasePtr
	ctx.State.Stack.Push(v)
}

func (ctx *Context) ValidateJumpBackward(i int) {
	if ctx.Pc > int(i) && ctx.Block.Code.Bytes[i] != OpLoopHead {
		panic("target of a backwards branch must be OpLoopHead")
	}
}

func (ctx *Context) CompareValues() (ComparisonResult, error) {
	v2 := ctx.State.Stack.Pop()
	v1 := ctx.State.Stack.Pop()
	if res, ok := CompareValues(v1, v2); ok {
		return res, nil
	} else {
		Debugf("compare %s %s", v1, v2)
		return 0, fmt.Errorf("both values must be same types")
	}
}

func (state *State) Exec(ctx *Context) (Value, error) {
	logExec := LogGroupEnabled(LogGroupExec)
	var bcDescs map[int]string
	var funToApply interface{}
	var argsToApply []Value
	var numSlots int
	stack := state.Stack

apply:
	if funToApply == nil {
		goto init
	}
	Debugf("apply %s", StringOfValue(funToApply))
	switch desc := funToApply.(type) {
	case *BlockClosure:
		ctx = state.NewContext(ctx.Module, ctx, desc, argsToApply)
		stack.BasePtr = stack.Ptr
		goto init
	case Primitive:
		ret, err := desc(state, ctx, argsToApply)
		if err != nil {
			return nil, err
		}
		state.Stack.Push(ret)
		goto loophead
	default:
		Panicf("not function %s", funToApply)
	}

init:
	numSlots = ctx.Block.Code.NumLocals() + ctx.Block.Code.FrameSize + ExtraNumSlots
	Debugf("stack grow if %d+%d > %d", stack.Ptr+1, numSlots, stack.Capa)
	if stack.Ptr+1+numSlots >= stack.Capa {
		stack.Increase(numSlots)
	}
	Debugf("eval start %d on %d", stack.Ptr, stack.Capa)

	if logExec {
		_, descs := ctx.Block.Code.StringsOfBytes(false)
		bcDescs = descs
		Logf(LogGroupExec, "==> exec %s", ctx.Block.Code.Name)
	}

	for i := 0; i < ctx.Block.Code.NumArgs; i++ {
		stack.Push(ctx.Args[i])
		if logExec {
			Logf(LogGroupExec, "    arg %d: %s", i, StringOfValue(ctx.Args[i]))
		}
	}
	stack.Ptr += ctx.Block.Code.NumArgs + ctx.Block.Code.NumTemps

loophead:
	for ctx.Pc+1 > len(ctx.Block.Code.Bytes) {
		Debugf("<== return")
		ctx = ctx.Parent
		if ctx == nil {
			return stack.Top(), nil
		}
		if logExec {
			_, descs := ctx.Block.Code.StringsOfBytes(false)
			bcDescs = descs
		}
	}
	bc := int(ctx.Block.Code.Bytes[ctx.Pc])
	if logExec {
		Debugf("%d %s", ctx.Pc, bcDescs[ctx.Pc])
		for i := stack.Ptr; i >= 0; i-- {
			v := StringOfValue(stack.Slots[i])
			if i == ctx.BasePtr+1 {
				Debugf(" -- at %d: %s", i, v)
			} else {
				Debugf("    at %d: %s", i, v)
			}
		}
	}

	ctx.Pc++
	switch {
	case OpLoadLocal0 <= bc && bc <= OpMaxLoadLocal:
		stack.PushLocal(bc - OpLoadLocal0)
	case OpLoadGlobal0 <= bc && bc <= OpMaxLoadGlobal:
		name := ctx.Const(bc - OpLoadGlobal0).(string)
		v, ok := ctx.Module.FindFieldValue(name)
		if !ok {
			Panicf("field %s is not found", name)
		}
		stack.Push(v)
	case OpLoadConst0 <= bc && bc <= OpMaxLoadConst:
		ctx.PushConst(bc - OpLoadConst0)
	case OpLoadValue0 <= bc && bc <= OpMaxLoadValue:
		path := ctx.Const(bc - OpLoadValue0).(*NamePath)
		f, err := state.FindFieldValueOfPath(path)
		if err != nil {
			return nil, err
		}
		stack.Push(f)
	case OpLoadInt0 <= bc && bc <= OpMaxLoadInt:
		stack.Push(bc - OpLoadInt0)
	case OpLoadIndirect0 <= bc && bc <= OpMaxLoadIndirect:
		ary := stack.TopArray()
		stack.Push(ary[bc-OpLoadIndirect0])
	case OpPopLoadIndirect0 <= bc && bc <= OpMaxPopLoadIndirect:
		ary := stack.PopArray()
		stack.Push(ary[bc-OpPopLoadIndirect0])
	case OpStorePopLocal0 <= bc && bc <= OpMaxStorePopLocal:
		stack.StorePopLocal(bc - OpStorePopLocal0)
	case OpStorePopGlobal0 <= bc && bc <= OpMaxStorePopGlobal:
		name := ctx.Const(bc - OpStorePopGlobal0).(string)
		ctx.Module.SetFieldValue(name, stack.Pop())
	case OpShortBranchFalse0 <= bc && bc <= OpMaxShortBranchFalse:
		if stack.Pop() == false {
			ctx.Pc += bc - OpShortBranchFalse0
		}
	case OpShortJump0 <= bc && bc <= OpMaxShortJump:
		ctx.Pc += bc - OpShortJump0 + 1
	case OpApply1 <= bc && bc <= OpMaxApply:
		argsToApply = stack.PopValues(bc - OpApply1 + 1)
		funToApply = stack.Pop()
		goto apply
	case OpApplyDirect1_0 <= bc && bc <= OpMaxApplyDirect1:
		argsToApply = stack.PopValues(1)
		funToApply = stack.Slots[ctx.BasePtr+bc-OpApplyDirect1_0]
		goto apply
	default:
		switch bc {
		case OpLoadUnit:
			stack.PushUnit()
		case OpLoadNil:
			stack.PushNil()
		case OpLoadHead:
			list := stack.Top().(*List)
			stack.Push(list.Head)
		case OpPopLoadTail:
			list := stack.Pop().(*List)
			stack.Push(list.Tail)
		case OpLoadTrue:
			stack.PushBool(true)
		case OpLoadFalse:
			stack.PushBool(false)
		case OpPop:
			stack.Pop()
		case OpSwapPop:
			stack.SwapPop()
		case OpDup:
			stack.Push(stack.Top())
		case OpReturn:
			ctx.Return(stack.Top())
			goto loophead
		case OpReturnUnit:
			ctx.Return(UnitValue)
			goto loophead
		case OpReturnTrue:
			ctx.Return(true)
			goto loophead
		case OpReturnFalse:
			ctx.Return(false)
			goto loophead
		case OpAdd:
			stack.Push(stack.PopInt() + stack.PopInt())
		case OpSub:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.Push(v1 - v2)
		case OpMul:
			stack.Push(stack.PopInt() * stack.PopInt())
		case OpDiv:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.Push(v1 / v2)
		case OpMod:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.Push(v1 % v2)
		case OpEq:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			stack.PushBool(res == OrderedSame)
		case OpNe:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			stack.PushBool(res != OrderedSame)
		case OpLe:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			stack.PushBool(res != OrderedAscending)
		case OpLt:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			stack.PushBool(res == OrderedDescending)
		case OpGe:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			stack.PushBool(res != OrderedDescending)
		case OpGt:
			res, err := ctx.CompareValues()
			if err != nil {
				return nil, err
			}
			stack.PushBool(res == OrderedAscending)
		case OpEqInts:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.PushBool(v1 == v2)
		case OpNeInts:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.PushBool(v1 != v2)
		case OpLtInts:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.PushBool(v1 < v2)
		case OpLeInts:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.PushBool(v1 <= v2)
		case OpGtInts:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.PushBool(v1 > v2)
		case OpGeInts:
			v2 := stack.PopInt()
			v1 := stack.PopInt()
			stack.PushBool(v1 >= v2)
		case OpAdd1:
			v := stack.PopInt()
			stack.Push(v + 1)
		case OpSub1:
			v := stack.PopInt()
			stack.Push(v - 1)
		case OpCountValues:
			switch vs := stack.Top().(type) {
			case []Value:
				stack.Push(len(vs))
			case *List:
				stack.Push(vs.Length())
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
	b1 := int(ctx.Block.Code.Bytes[ctx.Pc])
	ctx.Pc++
	switch {
	case OpLongJump0 <= bc && bc <= OpMaxLongJump:
		ctx.Pc += (bc-OpLongJump4)*256 + b1
	case OpLongBranchTrue0 <= bc && bc <= OpMaxLongBranchTrue:
		if stack.PopBool() {
			ctx.Pc += (bc-OpLongBranchTrue0)*256 + b1
		}
	case OpLongBranchFalse0 <= bc && bc <= OpMaxLongBranchFalse:
		if !stack.PopBool() {
			ctx.Pc += (bc-OpLongBranchFalse0)*256 + b1
		}
	default:
		switch bc {
		case OpLoadLocalIndirect:
			stack.PushLocalIndirect(b1)
		case OpXLoadIndirect:
			ary := stack.TopArray()
			stack.Push(ary[b1])
		case OpXStorePopLocal:
			stack.StorePopLocal(b1)
		case OpStorePopLocalIndirect:
			stack.StorePopLocalIndirect(b1)
		case OpXStorePopGlobal:
			name := ctx.Const(b1).(string)
			ctx.Module.SetFieldValue(name, stack.Pop())
		case OpXLoadInt:
			stack.Push(b1)
		case OpBranchNe:
			v2 := stack.Pop()
			v1 := stack.Pop()
			res, ok := CompareValues(v1, v2)
			if !ok || (ok && res != OrderedSame) {
				ctx.Pc += b1 + 1
			}
		case OpBranchNeSizes:
			v2 := len(stack.PopArray())
			v1 := len(stack.PopArray())
			if v1 != v2 {
				ctx.Pc += b1 + 1
			}
		case OpCreateArray:
			ary := make([]Value, b1+1)
			stack.Push(ary)
		case OpConsArray:
			ary := stack.PopValues(b1 + 1)
			stack.Push(ary)
		case OpConsList:
			list := NewListFromArray(stack.PopValues(b1 + 1))
			stack.Push(list)
		case OpCopyValues:
			if b1 != len(ctx.Block.Copied) {
				Panicf("number of copied values %d of the block is not equal to bytecode %d", len(ctx.Block.Copied), b1)
			}
			for i := 0; i < b1; i++ {
				stack.Push(ctx.Block.Copied[i])
			}
		case OpFullBlock:
			code := ctx.Const(b1).(*CompiledCode)
			blk := NewBlockClosure(code)
			blk.Context = ctx
			stack.Push(blk)
		default:
			goto threeByteCode
		}
	}
	goto loophead

threeByteCode:
	b2 := int(ctx.Block.Code.Bytes[ctx.Pc])
	ctx.Pc++
	switch bc {
	case OpXXLoadInt:
		stack.Push(b1*256 + b2)
	case OpCopyingBlock, OpFullCopyingBlock:
		code := ctx.Const(b1).(*CompiledCode)
		blk := NewBlockClosure(code)
		blk.Copied = stack.PopValues(b2)
		stack.Push(blk)
		if bc == OpFullCopyingBlock {
			blk.Context = ctx
		}
	case OpPrimitive:
		key := ctx.Const(b1).(string)
		args := stack.PopValues(b2)
		if f, ok := PrimMap[key]; ok {
			ret, err := f(state, ctx, args)
			if err != nil {
				return nil, err
			}
			stack.Push(ret)
		} else {
			return nil, fmt.Errorf("primitive `%s' is not found", key)
		}
	default:
		panic("notimpl op")
	}
	goto loophead
}
