package trompe

import (
	"fmt"
	"math/rand"
	"strings"
)

type Closure interface {
	Arity() int
	Apply(*Interp, *Context, *Env) (Value, error)
}

type CompiledCode struct {
	Id     int
	Syms   []string
	Lits   []Value
	Ops    []Opcode
	Labels map[int]int
}

func NewCompiledCode() *CompiledCode {
	return &CompiledCode{
		Id:     int(rand.Int31()),
		Syms:   []string{},
		Lits:   []Value{},
		Ops:    []Opcode{},
		Labels: make(map[int]int, 0),
	}
}

func (code *CompiledCode) AddLit(value Value) {
	code.Lits = append(code.Lits, value)
}

func (code *CompiledCode) LiteralDesc(i int) string {
	return code.Lits[i].Desc()
}

func (code *CompiledCode) Inspect() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("id: %d\n", code.Id))

	b.WriteString("symbols:\n")
	for i, name := range code.Syms {
		s := fmt.Sprintf("    %d: \"%s\"\n", i, name)
		b.WriteString(s)
	}

	b.WriteString("literals:\n")
	for i, value := range code.Lits {
		s := fmt.Sprintf("    %d: %s\n", i, value.Desc())
		b.WriteString(s)
	}

	b.WriteString("opcodes:\n")
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
			s += fmt.Sprintf("load local \"%s\"", code.Syms[i])
		case OpLoadAttr:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load attr \"%s\"", code.Syms[i])
		case OpLoadArg:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("load arg %d", i)
		case OpLoadModule:
			s += "load module"
		case OpStoreLocal:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("store local \"%s\"", code.LiteralDesc(i))
		case OpStoreRef:
			s += fmt.Sprintf("store ref")
		case OpStoreAttr:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("store attr \"%s\"", code.LiteralDesc(i))
		case OpPop:
			s += "pop"
		case OpDup:
			s += "dup"
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
		case OpBranchNext:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("branch next L%d", i)
		case OpBegin:
			s += "begin block"
		case OpEnd:
			s += "end block"
		case OpCall:
			i := code.Ops[pc+1]
			pc++
			s += fmt.Sprintf("call with %d args", i)
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
		case OpClosedRange:
			s += fmt.Sprintf("create closed range")
		case OpHalfOpenRange:
			s += fmt.Sprintf("create half-open range")
		default:
			panic(fmt.Sprintf("unknown opcode %d", code.Ops[pc]))
		}
		b.WriteString(s + "\n")
	}
	return b.String()
}

func (code *CompiledCode) Type() int {
	return ValClosType
}

func (code *CompiledCode) Desc() string {
	return fmt.Sprintf("CompiledCode %p", code)
}

func (code *CompiledCode) Bool() bool {
	panic("CompiledCode")
}

func (code *CompiledCode) Int() int {
	panic("CompiledCode")
}

func (code *CompiledCode) String() string {
	panic("CompiledCode")
}

func (code *CompiledCode) Closure() Closure {
	return code
}

func (code *CompiledCode) List() *List {
	panic("CompiledCode")
}

func (code *CompiledCode) Tuple() []Value {
	panic("CompiledCode")
}

func (code *CompiledCode) Arity() int {
	// TODO
	return 0
}

func (code *CompiledCode) Apply(ip *Interp, ctx *Context, env *Env) (Value, error) {
	return ip.Eval(ctx, env, code)
}

/*
	show "Hello, world!"
*/
func TestCompiledCodeHelloWorld() {
	code := CompiledCode{
		Syms: []string{
			"show",
		},
		Lits: []Value{
			NewValStr("show"),
			NewValStr("Hello, world!"),
		},
		Ops: []int{
			OpLoadLocal, 0, // "show"
			OpLoadLit, 1, // "Hello, world!"
			OpCall, 1,
			OpReturnUnit,
		},
	}
	fmt.Println(code.Inspect())

	ip := NewInterp()
	ctx := NewContext(nil, nil, &code, nil, 0)
	env := NewEnv(nil)
	ip.Eval(&ctx, env, &code)
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
		Syms: []string{
			"show",
		},
		Lits: []Value{
			NewValStr("Fizz"),
			NewValStr("Buzz"),
			NewValStr("FizzBuzz"),
			NewValPtn(NewPtnTuple(&PtnInt{0}, &PtnInt{0})),
			NewValPtn(NewPtnTuple(&PtnInt{0}, PtnWildcard)),
			NewValPtn(NewPtnTuple(PtnWildcard, &PtnInt{0})),
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
			OpLoadLocal, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz",
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 0,

			OpBegin,
			OpLoadLit, 5, // pattern 2
			OpMatch,
			OpBranchFalse, 1, // label 1
			OpLoadLocal, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz",
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 1,

			OpBegin,
			OpLoadLit, 6, // pattern 3
			OpMatch,
			OpBranchFalse, 2, // label 2
			OpLoadLocal, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz",
			OpCall, 1,
			OpEnd,
			OpJump, 4, // label 4
			OpLabel, 2,

			OpBegin,
			OpLoadLit, 6, // pattern 3
			OpMatch,
			OpBranchFalse, 2, // label 3
			OpLoadLocal, 0, // "show"
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
		ctx := NewContext(nil, &code, nil, 0)
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
		Syms: []string{
			"show",
		},
		Lits: []Value{
			NewValStr("show"),
			NewValStr("Fizz"),
			NewValStr("Buzz"),
			NewValStr("FizzBuzz"),
		},
		Ops: []int{
			OpLoadArg, 0, // arg 1
			OpLoadInt, 3, // 3
			OpEq,             // ==
			OpBranchFalse, 0, // label 0
			OpLoadLocal, 0, // "show"
			OpLoadLit, 1, // "Fizz"
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 0,

			OpLoadArg, 0, // arg 1
			OpLoadInt, 5, // 5
			OpEq,             // ==
			OpBranchFalse, 1, // label 1
			OpLoadLocal, 0, // "show"
			OpLoadLit, 2, // "Buzz"
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 1,

			OpLoadArg, 0, // arg 1
			OpLoadInt, 15, // 15
			OpEq,             // ==
			OpBranchFalse, 2, // label 2
			OpLoadLocal, 0, // "show"
			OpLoadLit, 3, // "FizzBuzz"
			OpCall, 1,
			OpJump, 4, // label 4
			OpLabel, 2,

			OpLoadLocal, 0, // "show"
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
		ctx := NewContext(nil, &code, nil, 0)
		prog.Eval(&ctx)
	*/
}
