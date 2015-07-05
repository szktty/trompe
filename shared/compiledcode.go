package trompe

import (
	"bytes"
	"fmt"
)

type CompiledCode struct {
	File      string
	Name      string
	Lines     []LineInfo
	Bytes     []uint8
	NumArgs   int
	NumTemps  int
	NumCopied int
	FrameSize int
	Consts    []Value
}

type LineInfo struct {
	Line int
	Pc   int
}

type PartialCompiledCode struct {
	Code *CompiledCode
	Args []Value
}

func (code *CompiledCode) LineAt(pc int) int {
	line := -1
	for _, info := range code.Lines {
		if pc >= info.Pc {
			line = info.Line
		} else {
			break
		}
	}
	return line
}

func (code *CompiledCode) NumLocals() int {
	return code.NumArgs + code.NumTemps
}

func (code *CompiledCode) String() string {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("File: %s\n", code.File))
	if code.Name != "" {
		buf.WriteString(fmt.Sprintf("Name: %s\n", code.Name))
	}
	buf.WriteString(fmt.Sprintf("NumArgs: %d\n", code.NumArgs))
	buf.WriteString(fmt.Sprintf("NumTemps: %d\n", code.NumTemps))
	buf.WriteString(fmt.Sprintf("FrameSize: %d\n", code.FrameSize))
	buf.WriteString(fmt.Sprintf("Consts: %d\n", len(code.Consts)))
	if len(code.Consts) > 0 {
		for i, lit := range code.Consts {
			buf.WriteString(fmt.Sprintf("    %d %s\n", i, StringOfValue(lit)))
		}
	}

	numBytes := len(code.Bytes)
	if numBytes > 0 {
		buf.WriteString("Bytes:\n")
		buf.WriteString(code.StringOfBytes())
		/*
			pre := -1
			var line int
		*/
		/*
			line = code.LineAt(pc)
			if (pre < 0 || pre != line) && line >= 0 {
				buf.WriteString(fmt.Sprintf("    Line %d:\n", line+1))
				pre = line
			}
		*/
	}
	return buf.String()
}

func (code *CompiledCode) StringOfConst(i int) string {
	return StringOfValue(code.Consts[i])
}

func (code *CompiledCode) StringOfBytes() string {
	buf := bytes.NewBufferString("")
	pcs, descs := code.StringsOfBytes(true)
	for i := 0; i < len(pcs); i++ {
		pc := pcs[i]
		buf.WriteString(fmt.Sprintf("    % 4d %s\n", pc, descs[pc]))
	}
	return buf.String()
}

func (code *CompiledCode) StringsOfBytes(pretty bool) ([]int, map[int]string) {
	pcs := make([]int, 0)
	descs := make(map[int]string)
	for pc := 0; pc < len(code.Bytes); {
		basePc := pc
		pcs = append(pcs, pc)
		buf := bytes.NewBufferString("")
		bc := int(code.Bytes[pc])
		pc++
		switch {
		case bc <= OpLastShortCode:
			buf.WriteString(fmt.Sprintf("<%02X> ", bc))
			if pretty {
				buf.WriteString("      ")
			}
			switch {
			case OpLoadLocal0 <= bc && bc <= OpMaxLoadLocal:
				buf.WriteString(fmt.Sprintf("push local %d", bc-OpLoadLocal0))
			case OpLoadGlobal0 <= bc && bc <= OpMaxLoadGlobal:
				buf.WriteString(fmt.Sprintf("push global %s",
					code.StringOfConst(bc-OpLoadGlobal0)))
			case OpLoadConst0 <= bc && bc <= OpMaxLoadConst:
				buf.WriteString(fmt.Sprintf("push %s",
					code.StringOfConst(bc-OpLoadConst0)))
			case OpLoadValue0 <= bc && bc <= OpMaxLoadValue:
				buf.WriteString(fmt.Sprintf("push %s",
					code.Consts[bc-OpLoadValue0].(*NamePath)))
			case OpLoadInt0 <= bc && bc <= OpMaxLoadInt:
				buf.WriteString(fmt.Sprintf("push %d", bc-OpLoadInt0))
			case OpLoadIndirect0 <= bc && bc <= OpMaxLoadIndirect:
				buf.WriteString(fmt.Sprintf("push indirect %d",
					bc-OpLoadIndirect0))
			case OpPopLoadIndirect0 <= bc && bc <= OpMaxPopLoadIndirect:
				buf.WriteString(fmt.Sprintf("pop; push indirect %d",
					bc-OpPopLoadIndirect0))
			case OpStorePopLocal0 <= bc && bc <= OpMaxStorePopLocal:
				buf.WriteString(fmt.Sprintf("store local %d; pop", bc-OpStorePopLocal0))
			case OpStorePopGlobal0 <= bc && bc <= OpMaxStorePopGlobal:
				buf.WriteString(fmt.Sprintf("store global %s; pop",
					code.StringOfConst(bc-OpStorePopGlobal0)))
			case OpStorePopField0 <= bc && bc <= OpMaxStorePopField:
				buf.WriteString(fmt.Sprintf("store field %s; pop",
					code.StringOfConst(bc-OpStorePopField0)))
			case OpApply1 <= bc && bc <= OpMaxApply:
				buf.WriteString(fmt.Sprintf("apply (%d)", bc-OpApply1+1))
			case OpApplyDirect1_0 <= bc && bc <= OpMaxApplyDirect1:
				buf.WriteString(fmt.Sprintf("apply local %d (1)",
					bc-OpApplyDirect1_0))
			case OpApplyDirect2_0 <= bc && bc <= OpMaxApplyDirect2:
				buf.WriteString(fmt.Sprintf("apply local %d (2)",
					bc-OpApplyDirect2_0))
			case OpApplyDirect3_0 <= bc && bc <= OpMaxApplyDirect3:
				buf.WriteString(fmt.Sprintf("apply local %d (3)",
					bc-OpApplyDirect3_0))
			case OpShortJump0 <= bc && bc <= OpMaxShortJump:
				buf.WriteString(fmt.Sprintf("jump %d", pc+int(bc-OpShortJump0)+1))
			case OpShortBranchFalse0 <= bc && bc <= OpMaxShortBranchFalse:
				buf.WriteString(fmt.Sprintf("jump false %d",
					pc+int(bc-OpShortBranchFalse0)))
			default:
				switch bc {
				case OpNoOp:
					buf.WriteString("no op")
				case OpLoadUnit:
					buf.WriteString("push ()")
				case OpLoadTrue:
					buf.WriteString("push true")
				case OpLoadFalse:
					buf.WriteString("push false")
				case OpLoadNil:
					buf.WriteString("push []")
				case OpLoadHead:
					buf.WriteString("push head")
				case OpPopLoadTail:
					buf.WriteString("pop; push tail")
				case OpLoadNone:
					buf.WriteString("push None")
				case OpLoadSome:
					buf.WriteString("pop; push Some")
				case OpPop:
					buf.WriteString("pop")
				case OpSwapPop:
					buf.WriteString("swap; pop")
				case OpDup:
					buf.WriteString("dup")
				case OpCountValues:
					buf.WriteString("count values")
				case OpReturn:
					buf.WriteString("return")
				case OpReturnUnit:
					buf.WriteString("return ()")
				case OpReturnTrue:
					buf.WriteString("return true")
				case OpReturnFalse:
					buf.WriteString("return false")
				case OpNot:
					buf.WriteString("not")
				case OpAdd:
					buf.WriteString("+")
				case OpSub:
					buf.WriteString("-")
				case OpMul:
					buf.WriteString("*")
				case OpDiv:
					buf.WriteString("/")
				case OpMod:
					buf.WriteString("mod")
				case OpEqInts:
					buf.WriteString("=")
				case OpNeInts:
					buf.WriteString("<>")
				case OpLtInts:
					buf.WriteString("<")
				case OpLeInts:
					buf.WriteString("<=")
				case OpGtInts:
					buf.WriteString(">")
				case OpGeInts:
					buf.WriteString(">=")
				case OpEq:
					buf.WriteString("=")
				case OpNe:
					buf.WriteString("<>")
				case OpLt:
					buf.WriteString("<")
				case OpLe:
					buf.WriteString("<=")
				case OpGt:
					buf.WriteString(">")
				case OpGe:
					buf.WriteString(">=")
				case OpBand:
					buf.WriteString("land")
				case OpBor:
					buf.WriteString("lor")
				case OpBxor:
					buf.WriteString("lxor")
				case OpLshift:
					buf.WriteString("lsl")
				case OpRshift:
					buf.WriteString("lsr")
				case OpNegInt:
					buf.WriteString("* -1")
				case OpLoopHead:
					buf.WriteString("loop head")
				case OpAdd1:
					buf.WriteString("+ 1")
				case OpSub1:
					buf.WriteString("- 1")
				default:
					Panicf("unknown short code %02X", bc)
				}
			}
		case bc <= OpLastTwoByteCode:
			b1 := int(code.Bytes[pc])
			pc++
			buf.WriteString(fmt.Sprintf("<%02X %02X> ", bc, b1))
			if pretty {
				buf.WriteString("   ")
			}
			switch {
			case OpXApplyDirect1 <= bc && bc <= OpMaxXApplyDirect:
				buf.WriteString(fmt.Sprintf("apply %s (%d)",
					code.StringOfConst(b1), bc-OpXApplyDirect1+1))
			case OpLongBranchFalse0 <= bc && bc <= OpMaxLongBranchFalse:
				buf.WriteString(fmt.Sprintf("jump false %d",
					pc+(int(bc)-OpLongBranchFalse0)*256+int(b1)))
			case OpLongBranchTrue0 <= bc && bc <= OpMaxLongBranchTrue:
				buf.WriteString(fmt.Sprintf("jump true %d",
					pc+(int(bc)-OpLongBranchTrue0)*256+int(b1)))
			case OpLongJump0 <= bc && bc <= OpMaxLongJump:
				buf.WriteString(fmt.Sprintf("jump %d",
					pc+(int(bc)-OpLongJumpBase)*256+int(b1)))
			default:
				switch bc {
				case OpXLoadLocal:
					buf.WriteString(fmt.Sprintf("push local %d", b1))
				case OpXLoadGlobal:
					buf.WriteString(fmt.Sprintf("push %s", code.StringOfConst(b1)))
				case OpXLoadConst:
					buf.WriteString(fmt.Sprintf("push %s", code.StringOfConst(b1)))
				case OpXLoadInt:
					buf.WriteString(fmt.Sprintf("push %d", b1))
				case OpXLoadValue:
					buf.WriteString(fmt.Sprintf("push %s",
						code.Consts[b1].(*NamePath)))
				case OpLoadLocalIndirect:
					buf.WriteString(fmt.Sprintf("push local %d at %d", b1/16, b1%16))
				case OpXLoadIndirect:
					buf.WriteString(fmt.Sprintf("push indirect %d", b1))
				case OpXStorePopLocal:
					buf.WriteString(fmt.Sprintf("store local %d; pop", b1))
				case OpStorePopLocalIndirect:
					buf.WriteString(fmt.Sprintf("store local %d at %d; pop",
						b1/16, b1%16))
				case OpXStorePopGlobal:
					buf.WriteString(fmt.Sprintf("store global %s; pop",
						code.StringOfConst(b1)))
				case OpBranchNe:
					buf.WriteString(fmt.Sprintf("jump <> %d", pc+b1+1))
				case OpBranchNeSizes:
					buf.WriteString(fmt.Sprintf("jump <> sizes %d", pc+b1+1))
				case OpXApply:
					buf.WriteString(fmt.Sprintf("apply (%d)", b1+1))
				case OpCreateArray:
					buf.WriteString(fmt.Sprintf("create array %d", b1+1))
				case OpConsArray:
					buf.WriteString(fmt.Sprintf("construct array %d", b1+1))
				case OpConsList:
					buf.WriteString(fmt.Sprintf("construct list %d", b1+1))
				case OpCopyValues:
					buf.WriteString(fmt.Sprintf("push %d copied values", b1))
				case OpFullBlock:
					buf.WriteString("make full block")
				default:
					Panicf("unknown two-byte code %02X", bc)
				}
			}
		case bc <= OpLastThreeByteCode:
			b1 := code.Bytes[pc]
			b2 := code.Bytes[pc+1]
			pc += 2
			buf.WriteString(fmt.Sprintf("<%02X %02X %02X> ", bc, b1, b2))
			switch bc {
			case OpXXLoadLocal:
				buf.WriteString(fmt.Sprintf("push local %d", IntOfTwoBytes(b1, b2)))
				/*
					case OpXPopLoadField:
						buf.WriteString(fmt.Sprintf("pop; push {%s.%s}",
							code.Consts[b1].(*BindingRef)).Key,
							code.Consts[b2].(string))
					case OpXStorePopField:
						buf.WriteString(fmt.Sprintf("store {%s.%s}; pop",
							code.Consts[b1].(*BindingRef)).Key,
							code.Consts[b2].(string))
				*/
			case OpXXLoadInt:
				buf.WriteString(fmt.Sprintf("push %d", IntOfTwoBytes(b1, b2)))
			case OpCopyingBlock:
				buf.WriteString(fmt.Sprintf("make copying block (%d)", b2))
			case OpFullCopyingBlock:
				buf.WriteString(fmt.Sprintf("make full copying block (%d)", b2))
			case OpPrimitive:
				buf.WriteString(fmt.Sprintf("primitive %s (%d)",
					code.StringOfConst(int(b1)), b2))
			default:
				Panicf("unknown three-byte code %02X", bc)
			}
		default:
			Panicf("unknown bytecode %02X", bc)
		}
		descs[basePc] = buf.String()
	}
	return pcs, descs
}
