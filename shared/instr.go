package trompe

import (
	"math"
)

type BytecodeBuilder struct {
	Instrs []*Instr
	Consts []Value
	Labels []*Instr
}

type Instr struct {
	Loc   *Loc
	Pc    int
	Size  int
	Bytes [3]uint8
	Desc  InstrDesc
}

type InstrDesc interface {
	Generate(*BytecodeBuilder, *Instr)
	StackSize() int
}

type LoadLocalInstr struct {
	Index int
}

type LoadLocalIndirectInstr struct {
	Array int
	Elt   int
}

type LoadIndirectInstr struct {
	Pop   bool
	Index int
}

type LoadGlobalInstr struct {
	Index int
}

type LoadConstInstr struct {
	Index int
}

type LoadValueInstr struct {
	Const Value
	Index int
}

type PopLoadFieldInstr struct {
	Index int
}

type LoadUnitInstr struct{}
type LoadNilInstr struct{}
type LoadHeadInstr struct{}
type PopLoadTailInstr struct{}

type LoadBoolInstr struct {
	Flag bool
}

type LoadIntInstr struct {
	Value int
}

type StorePopLocalInstr struct {
	Index int
}

type StorePopLocalIndirectInstr struct {
	Array int
	Elt   int
}

type StorePopGlobalInstr struct {
	Index int
}

type StorePopFieldInstr struct {
	Index int
}

type ReturnInstr struct{}
type ReturnUnitInstr struct{}

type ReturnBoolInstr struct {
	Value bool
}

type LoopHeadInstr struct{}

type LabelInstr struct {
	Name string
}

type JumpInstr struct {
	Dest *Instr
}

const (
	_BranchType = iota
	BranchFalse
	BranchTrue
	BranchNe
	BranchNeSizes
)

type BranchInstr struct {
	Type int
	Dest *Instr
}

type ApplyInstr struct {
	Arity int
}

type ApplyDirectInstr struct {
	Local int
	Arity int
}

type PopInstr struct{}
type SwapPopInstr struct{}
type DupInstr struct{}
type AddInstr struct{}
type Add1Instr struct{}
type SubInstr struct{}
type Sub1Instr struct{}
type MulInstr struct{}
type DivInstr struct{}
type PowInstr struct{}
type ModInstr struct{}
type EqIntsInstr struct{}
type NeIntsInstr struct{}
type LtIntsInstr struct{}
type LeIntsInstr struct{}
type GtIntsInstr struct{}
type GeIntsInstr struct{}
type EqInstr struct{}
type NeInstr struct{}
type LtInstr struct{}
type LeInstr struct{}
type GtInstr struct{}
type GeInstr struct{}
type LshiftInstr struct{}
type RshiftInstr struct{}
type NegIntInstr struct{}
type BnotInstr struct{}
type BandInstr struct{}
type BorInstr struct{}
type BxorInstr struct{}
type NotInstr struct{}

type CreateArrayInstr struct {
	Size int
}

type ConsArrayInstr struct {
	Size int
}

type ConsListInstr struct {
	Size int
}

type FullBlockInstr struct {
	Index int
}

type CopyingBlockInstr struct {
	Index   int
	NumCopy int
	Full    bool
}

type CopyValuesInstr struct {
	Size int
}

type CountValuesInstr struct{}

type PrimInstr struct {
	Index int
	Arity int
}

func NewBytecodeBuilder() *BytecodeBuilder {
	return &BytecodeBuilder{Consts: make([]Value, 0), Instrs: make([]*Instr, 0)}
}

func (bld *BytecodeBuilder) AddInstr(instr *Instr) *Instr {
	bld.Instrs = append(bld.Instrs, instr)
	return instr
}

func (bld *BytecodeBuilder) IndexOfInstr(instr *Instr) (int, bool) {
	for i, instr1 := range bld.Instrs {
		if instr == instr1 {
			return i, true
		}
	}
	return -1, false
}

func CompareConsts(v1 Value, v2 Value) bool {
	switch desc1 := v1.(type) {
	case []Value:
		if desc2, ok := v2.([]Value); ok && len(desc1) == len(desc2) {
			for i := 0; i < len(desc1); i++ {
				if !CompareConsts(desc1[i], desc2[i]) {
					return false
				}
			}
			return true
		}
		return false
	case *NamePath:
		if desc2, ok := v2.(*NamePath); ok {
			return desc1.String() == desc2.String()
		}
		return false
	default:
		return v1 == v2
	}
}

func (bld *BytecodeBuilder) AddConst(v Value) int {
	for i, cst := range bld.Consts {
		if CompareConsts(cst, v) {
			return int(i)
		}
	}
	bld.Consts = append(bld.Consts, v)
	return int(len(bld.Consts) - 1)
}

func (bld *BytecodeBuilder) Generate() ([]uint8, int) {
	bld.Optimize()
	bld.GenerateBasic()
	bld.AssignPcToInstrs()
	for i := 0; i < 3; i++ {
		bld.GenerateJump()
		bld.AssignPcToInstrs()
	}

	bytes := make([]uint8, 0)
	ptr := 0
	maxPtr := 0
	for _, ins := range bld.Instrs {
		for i := 0; i < ins.Size; i++ {
			bytes = append(bytes, ins.Bytes[i])
			ptr += ins.Desc.StackSize()
			if ptr >= maxPtr {
				maxPtr = ptr
			}
		}
	}
	return bytes, maxPtr
}

func (bld *BytecodeBuilder) Optimize() {
	rms := make([]*Instr, 0)
	filter := make([]*Instr, 0)

	// optimize
	var pre *Instr
	for _, instr := range bld.Instrs {
		if pre != nil {
			switch pre.Desc.(type) {
			case *EqInstr:
				if desc, ok := instr.Desc.(*BranchInstr); ok {
					if desc.Type == BranchFalse {
						rms = append(rms, pre)
						desc.Type = BranchNe
					}
				}
			}
		}
		pre = instr
	}

	// remove unused instructions
	for _, instr := range bld.Instrs {
		for _, rm := range rms {
			if rm == instr {
				goto next
			}
		}
		filter = append(filter, instr)
	next:
	}
	bld.Instrs = filter
}

func (bld *BytecodeBuilder) GenerateBasic() {
	for _, instr := range bld.Instrs {
		instr.Generate(bld)
	}
}

func (bld *BytecodeBuilder) AssignPcToInstrs() {
	pc := 0
	for _, instr := range bld.Instrs {
		instr.Pc = pc
		pc += instr.Size
	}
}

func (bld *BytecodeBuilder) Distance(src *Instr, dest *Instr) int {
	if _, ok := dest.Desc.(*LabelInstr); !ok {
		panic("error")
	}
	return dest.Pc - src.NextPc()
}

func (bld *BytecodeBuilder) GenerateJump() {
	for _, instr := range bld.Instrs {
		if srcDesc, ok := instr.Desc.(*JumpInstr); ok {
			dist := bld.Distance(instr, srcDesc.Dest)
			if dist >= 0 {
				if dist <= OpMaxShortJump-OpShortJump0 {
					instr.SetBytes(uint8(OpShortJump0 + dist - 1))
				} else {
					instr.SetBytes(uint8(OpLongJumpBase+dist/256), uint8(dist%256))
				}
			} else {
				dist *= -1
				instr.SetBytes(uint8(OpLongJumpBase-1-dist/256), uint8(256-dist%256))
			}
		} else if srcDesc, ok := instr.Desc.(*BranchInstr); ok {
			dist := bld.Distance(instr, srcDesc.Dest)
			if dist == 0 {
				Panicf("distance must not be 0")
			}
			switch srcDesc.Type {
			case BranchTrue:
				instr.SetBytes(uint8(OpLongBranchTrue0+dist/256), uint8(dist%256))
			case BranchFalse:
				if 0 < dist && dist <= OpMaxShortBranchFalse-OpShortBranchFalse0 {
					//Debugf("dist %d, %d", dist, OpMaxShortBranchFalse-OpShortBranchFalse0)
					instr.SetBytes(uint8(OpShortBranchFalse0 + dist))
				} else {
					if dist < 0 {
						dist *= -1
					}
					instr.SetBytes(uint8(OpLongBranchFalse0+dist/256),
						uint8(dist%256))
				}
			case BranchNe:
				instr.SetBytes(uint8(OpBranchNe), uint8(dist-1))
			case BranchNeSizes:
				instr.SetBytes(uint8(OpBranchNeSizes), uint8(dist-1))
			default:
				panic("unknown branch type")
			}
		}
	}
}

func NewInstr(loc *Loc, desc InstrDesc) *Instr {
	return &Instr{Loc: loc, Desc: desc}
}

func (instr *Instr) NextPc() int {
	return instr.Pc + instr.Size
}

func (instr *Instr) SetBytes(bs ...uint8) {
	instr.Size = len(bs)
	for i, b := range bs {
		instr.Bytes[i] = b
	}
}

func (instr *Instr) SetShortCode(i int,
	maxShort uint8, shortOp uint8, ops ...uint8) {
	if i <= int(maxShort-shortOp) {
		instr.SetBytes(shortOp + uint8(i))
	} else if len(ops) > 0 && i <= math.MaxUint8 {
		instr.SetBytes(ops[0], uint8(i))
	} else if len(ops) > 1 {
		instr.Size = 3
		instr.SetBytes(ops[1], uint8(i/256), uint8(i%256))
	} else {
		panic("error")
	}
}

func (instr *Instr) SetTwoByteCode(i int, shortOp uint8, twoOp uint8) {
	if i <= math.MaxUint8 {
		instr.SetBytes(shortOp + uint8(i))
	} else {
		instr.SetBytes(twoOp, uint8(i))
	}
}

func (instr *Instr) SetIndexedCode(i1 int, i2 int, twoOp uint8, threeOp uint8) {
	if i1 < 16 && i2 < 16 {
		instr.SetBytes(twoOp, uint8(i1*16)+uint8(i2))
	} else {
		instr.SetBytes(threeOp, uint8(i1), uint8(i2))
	}
}

func (instr *Instr) SetThreeByteCode(i1 uint8, i2 uint8,
	twoOp uint8, threeOp uint8) {
	if i1 < 32 && i2 < 32 {
		instr.Size = 2
		instr.Bytes[0] = twoOp
		instr.Bytes[1] = i1/32 + i2%32
	} else {
		instr.Size = 3
		instr.Bytes[0] = threeOp
		instr.Bytes[1] = i1
		instr.Bytes[2] = i2
	}
}

func (instr *Instr) Finished() bool {
	return instr.Size > 0
}

func (instr *Instr) Generate(bld *BytecodeBuilder) {
	instr.Desc.Generate(bld, instr)
}

func (desc *LoadLocalInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Index, OpMaxLoadLocal, OpLoadLocal0, OpXLoadLocal)
}

func (desc *LoadLocalInstr) StackSize() int {
	return 1
}

func (desc *LoadLocalIndirectInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetIndexedCode(desc.Array, desc.Elt, OpLoadLocalIndirect,
		OpXLoadLocalIndirect)
}

func (desc *LoadLocalIndirectInstr) StackSize() int {
	return 1
}

func (desc *LoadIndirectInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	if desc.Pop {
		if desc.Index > OpMaxPopLoadIndirect {
			panic("error")
		}
		instr.SetBytes(uint8(OpPopLoadIndirect0 + desc.Index))
	} else {
		instr.SetShortCode(desc.Index, OpMaxLoadIndirect, OpLoadIndirect0,
			OpXLoadIndirect)
	}
}

func (desc *LoadIndirectInstr) StackSize() int {
	if desc.Pop {
		return 0
	} else {
		return 1
	}
}

func (desc *LoadGlobalInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Index, OpMaxLoadGlobal, OpLoadGlobal0, OpXLoadGlobal)
}

func (desc *LoadGlobalInstr) StackSize() int {
	return 1
}

func (desc *LoadConstInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Index, OpMaxLoadConst, OpLoadConst0, OpXLoadConst)
}

func (desc *LoadConstInstr) StackSize() int {
	return 1
}

func (desc *LoadValueInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Index, OpMaxLoadValue, OpLoadValue0, OpXLoadValue)
}

func (desc *LoadValueInstr) StackSize() int {
	return 1
}

func (desc *LoadUnitInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLoadUnit)
}

func (desc *LoadUnitInstr) StackSize() int {
	return 1
}

func (desc *LoadNilInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLoadNil)
}

func (desc *LoadNilInstr) StackSize() int {
	return 1
}

func (desc *LoadHeadInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLoadHead)
}

func (desc *LoadHeadInstr) StackSize() int {
	return 1
}

func (desc *PopLoadTailInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpPopLoadTail)
}

func (desc *PopLoadTailInstr) StackSize() int {
	return 0
}

func NewLoadBoolInstr(loc *Loc, flag bool) *Instr {
	return NewInstr(loc, &LoadBoolInstr{Flag: flag})
}

func (desc *LoadBoolInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	if desc.Flag {
		instr.SetBytes(OpLoadTrue)
	} else {
		instr.SetBytes(OpLoadFalse)
	}
}

func (desc *LoadBoolInstr) StackSize() int {
	return 1
}

func NewLoadIntInstr(loc *Loc, v int) *Instr {
	return NewInstr(loc, &LoadIntInstr{Value: v})
}

func (desc *LoadIntInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	v := desc.Value
	if v >= 0 {
		if v <= OpMaxLoadInt-OpLoadInt0 {
			instr.SetBytes(OpLoadInt0 + uint8(v))
		} else if v <= math.MaxUint8 {
			instr.SetBytes(OpXLoadInt, uint8(v))
		} else if v <= math.MaxInt16 {
			instr.SetBytes(OpXXLoadInt, uint8(v/256), uint8(v%256))
		} else {
			// TODO: const
			panic("notimpl")
		}
	} else {
		if math.MinInt16 <= v {
			b1, b2 := TwoBytesOfNegInt(v)
			instr.SetBytes(OpXXLoadInt, b1, b2)
		} else {
			// TODO
			panic("notimpl")
		}
	}
}

func (desc *LoadIntInstr) StackSize() int {
	return 1
}

func (desc *StorePopLocalInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Index, OpMaxStorePopLocal, OpStorePopLocal0, OpXStorePopLocal)
}

func (desc *StorePopLocalInstr) StackSize() int {
	return -1
}

func (desc *StorePopLocalIndirectInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetIndexedCode(desc.Array, desc.Elt, OpStorePopLocalIndirect,
		OpXStorePopLocalIndirect)
}

func (desc *StorePopLocalIndirectInstr) StackSize() int {
	return -1
}

func (desc *StorePopGlobalInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	Debugf("store pop global %d", desc.Index)
	instr.SetShortCode(desc.Index, OpMaxStorePopGlobal, OpStorePopGlobal0, OpXStorePopGlobal)
}

func (desc *StorePopGlobalInstr) StackSize() int {
	return -1
}

func (desc *StorePopFieldInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Index, OpMaxStorePopField, OpStorePopField0, OpXStorePopField)
}

func (desc *StorePopFieldInstr) StackSize() int {
	return -1
}

func NewReturnInstr(loc *Loc) *Instr {
	return NewInstr(loc, &ReturnInstr{})
}

func (desc *ReturnInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpReturn)
}

func (desc *ReturnInstr) StackSize() int {
	return -1
}

func NewReturnUnitInstr(loc *Loc) *Instr {
	return NewInstr(loc, &ReturnUnitInstr{})
}

func (desc *ReturnUnitInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpReturnUnit)
}

func (desc *ReturnUnitInstr) StackSize() int {
	return 0
}

func NewReturnBoolInstr(loc *Loc, cond bool) *Instr {
	return NewInstr(loc, &ReturnBoolInstr{Value: cond})
}

func (desc *ReturnBoolInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	if desc.Value {
		instr.SetBytes(OpReturnTrue)
	} else {
		instr.SetBytes(OpReturnFalse)
	}
}

func (desc *ReturnBoolInstr) StackSize() int {
	return 0
}

func (desc *LoopHeadInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLoopHead)
}

func (desc *LoopHeadInstr) StackSize() int {
	return 0
}

func NewLabelInstr(loc *Loc, name string) *Instr {
	return NewInstr(loc, &LabelInstr{Name: name})
}

func (desc *LabelInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.Size = 0
}

func (desc *LabelInstr) StackSize() int {
	return 0
}

func (desc *JumpInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	// assume byte size to be 2
	instr.Size = 2
}

func (desc *JumpInstr) StackSize() int {
	return 0
}

func (desc *BranchInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	// assume byte size to be 2
	instr.Size = 2
}

func (desc *BranchInstr) StackSize() int {
	switch desc.Type {
	case BranchNe, BranchNeSizes:
		return -2
	default:
		return -1
	}
}

func (desc *ApplyInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetShortCode(desc.Arity-1, OpMaxApply, OpApply1, OpXApply)
}

func (desc *ApplyInstr) StackSize() int {
	return desc.Arity
}

func (desc *ApplyDirectInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	switch desc.Arity {
	case 1:
		instr.SetShortCode(desc.Local, OpMaxApplyDirect1, OpApplyDirect1_0, OpXApplyDirect1)
	case 2:
		instr.SetShortCode(desc.Local, OpMaxApplyDirect2, OpApplyDirect2_0, OpXApplyDirect2)
	case 3:
		instr.SetShortCode(desc.Local, OpMaxApplyDirect3, OpApplyDirect3_0, OpXApplyDirect3)
	default:
		instr.SetBytes(OpXXApplyDirect, uint8(desc.Local), uint8(desc.Arity))
	}
}

func (desc *ApplyDirectInstr) StackSize() int {
	return desc.Arity
}

func (desc *PopInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpPop)
}

func (desc *PopInstr) StackSize() int {
	return -1
}

func (desc *SwapPopInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpSwapPop)
}

func (desc *SwapPopInstr) StackSize() int {
	return -1
}

func (desc *DupInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpDup)
}

func (desc *DupInstr) StackSize() int {
	return 1
}

func (desc *AddInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpAdd)
}

func (desc *AddInstr) StackSize() int {
	return -1
}

func (desc *Add1Instr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpAdd1)
}

func (desc *Add1Instr) StackSize() int {
	return 0
}

func (desc *SubInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpSub)
}

func (desc *SubInstr) StackSize() int {
	return -1
}

func (desc *Sub1Instr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpSub1)
}

func (desc *Sub1Instr) StackSize() int {
	return 0
}

func (desc *MulInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpMul)
}

func (desc *MulInstr) StackSize() int {
	return -1
}

func (desc *DivInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpDiv)
}

func (desc *DivInstr) StackSize() int {
	return -1
}

func (desc *PowInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpPow)
}

func (desc *PowInstr) StackSize() int {
	return -1
}

func (desc *ModInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpMod)
}

func (desc *ModInstr) StackSize() int {
	return -1
}

func (desc *EqIntsInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpEqInts)
}

func (desc *EqIntsInstr) StackSize() int {
	return -1
}

func (desc *NeIntsInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpNeInts)
}

func (desc *NeIntsInstr) StackSize() int {
	return -1
}

func (desc *LtIntsInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLtInts)
}

func (desc *LtIntsInstr) StackSize() int {
	return -1
}

func (desc *LeIntsInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLeInts)
}

func (desc *LeIntsInstr) StackSize() int {
	return -1
}

func (desc *GtIntsInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpGtInts)
}

func (desc *GtIntsInstr) StackSize() int {
	return -1
}

func (desc *GeIntsInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpGeInts)
}

func (desc *GeIntsInstr) StackSize() int {
	return -1
}

func (desc *EqInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpEq)
}

func (desc *EqInstr) StackSize() int {
	return -1
}

func (desc *NeInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpNe)
}

func (desc *NeInstr) StackSize() int {
	return -1
}

func (desc *LtInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLt)
}

func (desc *LtInstr) StackSize() int {
	return -1
}

func (desc *LeInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLe)
}

func (desc *LeInstr) StackSize() int {
	return -1
}

func (desc *GtInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpGt)
}

func (desc *GtInstr) StackSize() int {
	return -1
}

func (desc *GeInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpGe)
}

func (desc *GeInstr) StackSize() int {
	return -1
}

func (desc *LshiftInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpLshift)
}

func (desc *LshiftInstr) StackSize() int {
	return -1
}

func (desc *RshiftInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpRshift)
}

func (desc *RshiftInstr) StackSize() int {
	return -1
}

func (desc *BnotInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpBnot)
}

func (desc *BnotInstr) StackSize() int {
	return -1
}

func (desc *BandInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpBand)
}

func (desc *BandInstr) StackSize() int {
	return -1
}

func (desc *BorInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpBor)
}

func (desc *BorInstr) StackSize() int {
	return -1
}

func (desc *BxorInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpBxor)
}

func (desc *BxorInstr) StackSize() int {
	return -1
}

func (desc *NegIntInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpNegInt)
}

func (desc *NegIntInstr) StackSize() int {
	return 0
}

func (desc *NotInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpNot)
}

func (desc *NotInstr) StackSize() int {
	return 0
}

func (desc *CreateArrayInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpCreateArray, uint8(desc.Size)-1)
}

func (desc *CreateArrayInstr) StackSize() int {
	return 1
}

func (desc *ConsArrayInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpConsArray, uint8(desc.Size)-1)
}

func (desc *ConsArrayInstr) StackSize() int {
	return 1
}

func (desc *ConsListInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpConsList, uint8(desc.Size)-1)
}

func (desc *ConsListInstr) StackSize() int {
	return 1
}

func (desc *FullBlockInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpFullBlock, uint8(desc.Index))
}

func (desc *FullBlockInstr) StackSize() int {
	return 1
}

func (desc *CopyingBlockInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	if desc.Full {
		instr.SetBytes(OpFullCopyingBlock, uint8(desc.Index), uint8(desc.NumCopy))
	} else {
		instr.SetBytes(OpCopyingBlock, uint8(desc.Index), uint8(desc.NumCopy))
	}
}

func (desc *CopyingBlockInstr) StackSize() int {
	return -(desc.NumCopy) + 1
}

func (desc *CopyValuesInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpCopyValues, uint8(desc.Size))
}

func (desc *CopyValuesInstr) StackSize() int {
	return desc.Size
}

func (desc *CountValuesInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpCountValues)
}

func (desc *CountValuesInstr) StackSize() int {
	return 1
}

func NewPrimInstr(loc *Loc, idx int, arity int) *Instr {
	return NewInstr(loc, &PrimInstr{Index: idx, Arity: arity})
}

func (desc *PrimInstr) Generate(bld *BytecodeBuilder, instr *Instr) {
	instr.SetBytes(OpPrimitive, uint8(desc.Index), uint8(desc.Arity))
}

func (desc *PrimInstr) StackSize() int {
	return desc.Arity - 1
}
