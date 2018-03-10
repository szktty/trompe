package trompe

import (
	"fmt"
	"strings"
)

type Closure interface {
	Apply(*Program, *Context) Value
}

type CompiledCode struct {
	Lits []Value
	Ops  []Opcode
}

func (code *CompiledCode) Apply(prog *Program, ctx *Context) Value {
	return prog.Eval(ctx)
}

func (code *CompiledCode) LiteralDesc(i int) string {
	return code.Lits[i].Desc()
}

func (code *CompiledCode) Inspect() string {
	var b strings.Builder
	b.WriteString("literals:\n")
	for i, value := range code.Lits {
		s := fmt.Sprintf("    %d: %s\n", i, value.Desc())
		b.WriteString(s)
	}

	b.WriteString("\nopcodes:\n")
	pc := 0
	for ; pc < len(code.Ops); pc++ {
		s := fmt.Sprintf("    %d: ", pc)
		switch code.Ops[pc] {
		case OpNop:
			s += "nop"
		case OpLoadUnit:
			s += "load ()"
		case OpLoadTrue:
			s += "load true"
		case OpLoadFalse:
			s += "load false"
		case OpLoadInt:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load %d", i)
		case OpLoadNone:
			s += "load none"
		case OpLoadRef:
			s += "load ref"
		case OpLoadLit:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load literal %s", code.LiteralDesc(i))
		case OpLoadLocal:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load local %s", code.LiteralDesc(i))
		case OpLoadAttr:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load attr %s", code.LiteralDesc(i))
		case OpLoadPrim:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load primitive %s", code.LiteralDesc(i))
		case OpStore:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("store %s", code.LiteralDesc(i))
		case OpStoreRef:
			s += fmt.Sprintf("store ref")
		case OpStoreAttr:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("store attr %s", code.LiteralDesc(i))
		case OpPop:
			s += "pop"
		case OpReturn:
			s += "return"
		case OpLabel:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("label %s", code.LiteralDesc(i))
		case OpJump:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("jump %s", code.LiteralDesc(i))
		case OpBranchTrue:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("branch true %s", code.LiteralDesc(i))
		case OpBranchFalse:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("branch false %s", code.LiteralDesc(i))
		case OpCall:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("call %d", i)
		case OpSome:
			s += "create some"
		case OpList:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("create list %d", i)
		case OpTuple:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("create tuple %d", i)
		default:
			panic("unknown opcode")
		}
		b.WriteString(s + "\n")
	}
	return b.String()
}

type Primitive struct {
	Func func(*Program, []Value, int) Value
}

func (prim *Primitive) Apply(prog *Program, ctx *Context) Value {
	return prim.Func(prog, ctx.Args, ctx.Len)
}

func CompiledCodeTest() {
	code := CompiledCode{
		Lits: []Value{
			CreateValStr("hello"),
			CreateValStr("loophead"),
		},
		Ops: []int{
			OpNop,
			OpLoadUnit,
			OpLoadTrue,
			OpLoadFalse,
			OpLoadInt, 12345,
			OpLoadLit, 0,
			OpPop,
			OpLabel, 1,
			OpReturn,
		},
	}
	fmt.Println(code.Inspect())

	prog := Program{}
	ctx := CreateContext(nil, &code, nil, 0)
	prog.Eval(&ctx)
}
