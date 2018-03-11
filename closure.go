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
		case OpLoadZero:
			s += "load 0"
		case OpLoadOne:
			s += "load 1"
		case OpLoadNegOne:
			s += "load -1"
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
		case OpLoadArg:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load arg %d", i)
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
		case OpReturnUnit:
			s += "return ()"
		case OpLabel:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("label L%d", i)
		case OpJump:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("jump L%d", i)
		case OpBranchTrue:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("branch true L%d", i)
		case OpBranchFalse:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("branch false L%d", i)
		case OpBegin:
			s += "begin block"
		case OpEnd:
			s += "end block"
		case OpCall:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("call %d", i)
		case OpEq:
			s += "=="
		case OpNe:
			s += "=="
		case OpLt:
			s += "<"
		case OpLe:
			s += "<="
		case OpGt:
			s += ">"
		case OpGe:
			s += ">="
		case OpMatch:
			s += "match"
		case OpAdd:
			s += "+"
		case OpSub:
			s += "-"
		case OpMul:
			s += "*"
		case OpDiv:
			s += "/"
		case OpMod:
			s += "%"
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
			panic(fmt.Sprintf("unknown opcode %d", code.Ops[pc]))
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

/*
	show "Hello, world!"
*/
func TestCompiledCodeHelloWorld() {
	code := CompiledCode{
		Lits: []Value{
			CreateValStr("show"),
			CreateValStr("Hello, world!"),
		},
		Ops: []int{
			OpLoadPrim, 0, // "show"
			OpLoadLit, 1, // "Hello, world!"
			OpCall, 1,
			OpReturnUnit,
		},
	}
	fmt.Println(code.Inspect())

	/*
		prog := Program{}
		ctx := CreateContext(nil, &code, nil, 0)
		prog.Eval(&ctx)
	*/
}

/*
	match (arg1 % 3, arg1 % 5) with
	case (0, 0) then show("FizzBuzz")
	case (0, _) then show("Fizz")
	case (_, 0) then show("Buzz")
	else show(arg1)
	end
*/

func TestCompiledCodeFizzBuzzMatch() {
	code := CompiledCode{
		Lits: []Value{
			CreateValStr("show"),
			CreateValStr("Fizz"),
			CreateValStr("Buzz"),
			CreateValStr("FizzBuzz"),
			CreateValPtn(CreatePtnTuple(&PtnInt{0}, &PtnInt{0})),
			CreateValPtn(CreatePtnTuple(&PtnInt{0}, PtnWildcard)),
			CreateValPtn(CreatePtnTuple(PtnWildcard, &PtnInt{0})),
		},
		Ops: []int{
			OpLoadArg, 0,
			OpLoadInt, 3,
			OpMod,
			OpLoadArg, 0,
			OpLoadInt, 5,
			OpMod,
			OpTuple, 2,

			OpBegin,
			OpLoadLit, 4, // pattern 1
			OpMatch,
			OpBranchFalse, 0, // label 0
			OpLoadPrim, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz",
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 0,

			OpBegin,
			OpLoadLit, 5, // pattern 2
			OpMatch,
			OpBranchFalse, 1, // label 1
			OpLoadPrim, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz",
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 1,

			OpBegin,
			OpLoadLit, 6, // pattern 3
			OpMatch,
			OpBranchFalse, 2, // label 2
			OpLoadPrim, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz",
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 2,

			OpBegin,
			OpLoadLit, 6, // pattern 3
			OpMatch,
			OpBranchFalse, 2, // label 3
			OpLoadPrim, 0, // "show"
			OpLoadArg, 0, // arg 1
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 3,

			OpLabel, 4, // match end
			OpReturnUnit,
		},
	}
	fmt.Println(code.Inspect())

	/*
		prog := Program{}
		ctx := CreateContext(nil, &code, nil, 0)
		prog.Eval(&ctx)
	*/
}

/*
	if arg1 == 3 then
		show("Fizz")
	else if arg1 == 5 then
		show("Buzz")
	else if arg1 == 15 then
		show("FizzBuzz")
	else
		show(arg1)
	end
*/
func TestCompiledCodeFizzBuzzCompare() {
	code := CompiledCode{
		Lits: []Value{
			CreateValStr("show"),
			CreateValStr("Fizz"),
			CreateValStr("Buzz"),
			CreateValStr("FizzBuzz"),
		},
		Ops: []int{
			OpLoadArg, 0, // arg 1
			OpLoadInt, 3, // 3
			OpEq,             // ==
			OpBranchFalse, 0, // label 0
			OpLoadPrim, 0, // "show"
			OpLoadLit, 1, // "Fizz"
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 0,

			OpLoadArg, 0, // arg 1
			OpLoadInt, 5, // 5
			OpEq,             // ==
			OpBranchFalse, 1, // label 1
			OpLoadPrim, 0, // "show"
			OpLoadLit, 2, // "Buzz"
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 1,

			OpLoadArg, 0, // arg 1
			OpLoadInt, 15, // 15
			OpEq,             // ==
			OpBranchFalse, 2, // label 2
			OpLoadPrim, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz"
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 2,

			OpLoadPrim, 0, // "show"
			OpLoadArg, 0, // arg 1
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 3,

			OpLabel, 4,
			OpReturnUnit,
		},
	}
	fmt.Println(code.Inspect())

	/*
		prog := Program{}
		ctx := CreateContext(nil, &code, nil, 0)
		prog.Eval(&ctx)
	*/
}
