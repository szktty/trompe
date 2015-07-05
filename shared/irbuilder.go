package trompe

import (
	"strconv"
)

// intermediate representation builder
type IRBuilder struct {
	Parent *IRBuilder
	Bytes  *BytecodeBuilder
	Scope  *LocalScope
	Name   string
}

func NewIRBuilder(parent *IRBuilder, scope *LocalScope) *IRBuilder {
	return &IRBuilder{Parent: parent, Bytes: NewBytecodeBuilder(),
		Scope: scope}
}

func (bld *IRBuilder) CompiledCode() *CompiledCode {
	bytes, frameSize := bld.Bytes.Generate()
	return &CompiledCode{File: "<none>", Name: bld.Name,
		Lines: make([]LineInfo, 0),
		Bytes: bytes, NumArgs: bld.Scope.NumArgs(),
		NumTemps:  bld.Scope.NumTempSlots(),
		NumCopied: bld.Scope.NumCopied(),
		FrameSize: frameSize, Consts: bld.Bytes.Consts}
}

func (bld *IRBuilder) IsTop() bool {
	return bld.Parent == nil
}

func (bld *IRBuilder) AddConst(v Value) int {
	return bld.Bytes.AddConst(v)
}

/*
func (bld *IRBuilder) AddArray(size int) int {
	bld.Arrays = append(bld.Arrays, size)
	return len(bld.Temps) + len(bld.Outers) + len(bld.Arrays) - 1
}
*/

func (bld *IRBuilder) LocalIndex(name string) (int, bool) {
	if va, ok := bld.Scope.FindTemp(name); ok {
		return va.Index, true
	} else {
		return -1, false
	}
}

func (bld *IRBuilder) LocalIndirectIndex(name string) (int, int, bool) {
	if grp, elt, ok := bld.Scope.FindMember(name); ok {
		return grp.Index, elt.Index, true
	} else {
		return -1, -1, false
	}
}

/*
func (bld *IRBuilder) RecArrayIndex() int {
	return len(bld.Temps) + len(bld.Outers)
}
*/

func (bld *IRBuilder) AddInstr(instr *Instr) {
	bld.Bytes.AddInstr(instr)
}

func (bld *IRBuilder) PushLocal(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadLocalInstr{Index: idx}))
}

func (bld *IRBuilder) PushLocalIndirect(loc *Loc, array int, elt int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadLocalIndirectInstr{Array: array, Elt: elt}))
}

func (bld *IRBuilder) PushGlobalValue(loc *Loc, name string) {
	i := bld.AddConst(name)
	bld.Bytes.AddInstr(NewInstr(loc, &LoadGlobalInstr{Index: i}))
}

func (bld *IRBuilder) PushGlobal(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadGlobalInstr{Index: idx}))
}

func (bld *IRBuilder) PushIndirect(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadIndirectInstr{Index: idx}))
}

func (bld *IRBuilder) PopPushIndirect(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadIndirectInstr{Pop: true, Index: idx}))
}

func (bld *IRBuilder) PushVar(loc *Loc, name string) bool {
	if _, ok := bld.Scope.FindGlobal(name); ok {
		bld.PushGlobalValue(loc, name)
		return true
	} else if i, ok := bld.LocalIndex(name); ok {
		Debugf("push %s at %d", name, i)
		bld.PushLocal(loc, i)
		return true
	} else if aryIdx, eltIdx, ok := bld.LocalIndirectIndex(name); ok {
		bld.PushLocalIndirect(loc, aryIdx, eltIdx)
		return true
	} else {
		return false
	}
}

func (bld *IRBuilder) PushConst(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadConstInstr{Index: idx}))
}

func (bld *IRBuilder) PushValue(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadValueInstr{Index: idx}))
}

func (bld *IRBuilder) PushConstValue(loc *Loc, v Value) {
	i := bld.AddConst(v)
	bld.PushConst(loc, i)
}

func (bld *IRBuilder) PushUnit(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadUnitInstr{}))
}

func (bld *IRBuilder) PushNil(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadNilInstr{}))
}

func (bld *IRBuilder) PushHead(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadHeadInstr{}))
}

func (bld *IRBuilder) PopPushTail(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &PopLoadTailInstr{}))
}

func (bld *IRBuilder) PushBool(loc *Loc, flag bool) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadBoolInstr{Flag: flag}))
}

func (bld *IRBuilder) PushInt(loc *Loc, v int) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoadIntInstr{Value: v}))
}

func (bld *IRBuilder) PushIntOfString(loc *Loc, s string) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	bld.PushInt(loc, int(v))
}

func (bld *IRBuilder) StorePopLocal(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &StorePopLocalInstr{Index: idx}))
}

func (bld *IRBuilder) StorePopLocalIndirect(loc *Loc, ary int, elt int) {
	bld.Bytes.AddInstr(NewInstr(loc,
		&StorePopLocalIndirectInstr{Array: ary, Elt: elt}))
}

func (bld *IRBuilder) StorePopGlobal(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &StorePopGlobalInstr{Index: idx}))
}

func (bld *IRBuilder) StorePopVar(loc *Loc, name string) bool {
	if _, ok := bld.Scope.FindGlobal(name); ok {
		i := bld.AddConst(name)
		bld.StorePopGlobal(loc, i)
		return true
	} else if i, ok := bld.LocalIndex(name); ok {
		bld.StorePopLocal(loc, i)
		return true
	} else if aryIdx, eltIdx, ok := bld.LocalIndirectIndex(name); ok {
		bld.StorePopLocalIndirect(loc, aryIdx, eltIdx)
		return true
	} else {
		return false
	}
}

func (bld *IRBuilder) PutReturn(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &ReturnInstr{}))
}

func (bld *IRBuilder) PutReturnUnit(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &ReturnUnitInstr{}))
}

func (bld *IRBuilder) Pop(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &PopInstr{}))
}

func (bld *IRBuilder) SwapPop(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &SwapPopInstr{}))
}

func (bld *IRBuilder) Dup(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &DupInstr{}))
}

func (bld *IRBuilder) PutLoopHead(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LoopHeadInstr{}))
}

func (bld *IRBuilder) NewLabel(loc *Loc) *Instr {
	return NewLabelInstr(loc, "")
}

func (bld *IRBuilder) PutLabel(loc *Loc, name string) *Instr {
	desc := NewLabelInstr(loc, name)
	bld.Bytes.AddInstr(desc)
	return desc
}

func (bld *IRBuilder) PutJump(loc *Loc, dest *Instr) {
	bld.Bytes.AddInstr(NewInstr(loc, &JumpInstr{Dest: dest}))
}

func (bld *IRBuilder) PutBranchTrue(loc *Loc, dest *Instr) {
	bld.Bytes.AddInstr(NewInstr(loc, &BranchInstr{Type: BranchTrue, Dest: dest}))
}

func (bld *IRBuilder) PutBranchFalse(loc *Loc, dest *Instr) {
	bld.Bytes.AddInstr(NewInstr(loc, &BranchInstr{Type: BranchFalse, Dest: dest}))
}

func (bld *IRBuilder) PutBranchNe(loc *Loc, dest *Instr) {
	bld.Bytes.AddInstr(NewInstr(loc, &BranchInstr{Type: BranchNe, Dest: dest}))
}

func (bld *IRBuilder) PutBranchNeSizes(loc *Loc, dest *Instr) {
	bld.Bytes.AddInstr(NewInstr(loc, &BranchInstr{Type: BranchNeSizes, Dest: dest}))
}

func (bld *IRBuilder) Apply(loc *Loc, arity int) {
	bld.Bytes.AddInstr(NewInstr(loc, &ApplyInstr{Arity: arity}))
}

func (bld *IRBuilder) ApplyDirect(loc *Loc, local int, arity int) {
	bld.Bytes.AddInstr(NewInstr(loc, &ApplyDirectInstr{Local: local, Arity: arity}))
}

func (bld *IRBuilder) PutAdd(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &AddInstr{}))
}

func (bld *IRBuilder) PutAdd1(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &Add1Instr{}))
}

func (bld *IRBuilder) PutSub(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &SubInstr{}))
}

func (bld *IRBuilder) PutSub1(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &Sub1Instr{}))
}

func (bld *IRBuilder) PutMul(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &MulInstr{}))
}

func (bld *IRBuilder) PutDiv(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &DivInstr{}))
}

func (bld *IRBuilder) PutPow(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &PowInstr{}))
}

func (bld *IRBuilder) PutMod(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &ModInstr{}))
}

func (bld *IRBuilder) PutEq(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &EqInstr{}))
}

func (bld *IRBuilder) PutNe(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &NeInstr{}))
}

func (bld *IRBuilder) PutLt(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LtInstr{}))
}

func (bld *IRBuilder) PutLe(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LeInstr{}))
}

func (bld *IRBuilder) PutGt(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &GtInstr{}))
}

func (bld *IRBuilder) PutGe(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &GeInstr{}))
}

func (bld *IRBuilder) PutEqInts(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &EqIntsInstr{}))
}

func (bld *IRBuilder) PutNeInts(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &NeIntsInstr{}))
}

func (bld *IRBuilder) PutLtInts(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LtIntsInstr{}))
}

func (bld *IRBuilder) PutLeInts(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LeIntsInstr{}))
}

func (bld *IRBuilder) PutGtInts(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &GtIntsInstr{}))
}

func (bld *IRBuilder) PutGeInts(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &GeIntsInstr{}))
}

func (bld *IRBuilder) PutLshift(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &LshiftInstr{}))
}

func (bld *IRBuilder) PutRshift(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &RshiftInstr{}))
}

func (bld *IRBuilder) PutNegInt(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &NegIntInstr{}))
}

func (bld *IRBuilder) PutBnot(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &BnotInstr{}))
}

func (bld *IRBuilder) PutBand(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &BandInstr{}))
}

func (bld *IRBuilder) PutBor(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &BorInstr{}))
}

func (bld *IRBuilder) PutBxor(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &BxorInstr{}))
}

func (bld *IRBuilder) PutNot(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &NotInstr{}))
}

func (bld *IRBuilder) PutCreateArray(loc *Loc, size int) {
	bld.Bytes.AddInstr(NewInstr(loc, &CreateArrayInstr{Size: size}))
}

func (bld *IRBuilder) PutConsArray(loc *Loc, size int) {
	bld.Bytes.AddInstr(NewInstr(loc, &ConsArrayInstr{Size: size}))
}

func (bld *IRBuilder) PutConsList(loc *Loc, size int) {
	bld.Bytes.AddInstr(NewInstr(loc, &ConsListInstr{Size: size}))
}

func (bld *IRBuilder) MakeFullBlock(loc *Loc, idx int) {
	bld.Bytes.AddInstr(NewInstr(loc, &FullBlockInstr{Index: idx}))
}

func (bld *IRBuilder) MakeCopyingBlock(loc *Loc, idx int, numCopy int) {
	bld.Bytes.AddInstr(NewInstr(loc,
		&CopyingBlockInstr{Index: idx, NumCopy: numCopy}))
}

func (bld *IRBuilder) MakeFullCopyingBlock(loc *Loc, idx int, numCopy int) {
	bld.Bytes.AddInstr(NewInstr(loc,
		&CopyingBlockInstr{Index: idx, NumCopy: numCopy, Full: true}))
}

func (bld *IRBuilder) PushCopied(loc *Loc, n int) {
	bld.Bytes.AddInstr(NewInstr(loc, &CopyValuesInstr{Size: n}))
}

func (bld *IRBuilder) PutCountValues(loc *Loc) {
	bld.Bytes.AddInstr(NewInstr(loc, &CountValuesInstr{}))
}

func (bld *IRBuilder) PutPrimitive(loc *Loc, name string, arity int) {
	i := bld.AddConst(name)
	bld.Bytes.AddInstr(NewPrimInstr(loc, i, arity))
}
