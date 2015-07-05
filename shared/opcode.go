package trompe

import (
//"fmt"
)

const (
	//
	// Short codes
	//

	OpLoadLocal0 = iota
	OpLoadLocal1
	OpLoadLocal2
	OpLoadLocal3
	OpLoadLocal4
	OpLoadLocal5
	OpLoadLocal6
	OpLoadLocal7
	OpLoadLocal8
	OpLoadLocal9
	OpLoadLocal10
	OpMaxLoadLocal

	OpLoadGlobal0
	OpLoadGlobal1
	OpLoadGlobal2
	OpLoadGlobal3
	OpLoadGlobal4
	OpLoadGlobal5
	OpLoadGlobal6
	OpLoadGlobal7
	OpLoadGlobal8
	OpLoadGlobal9
	OpLoadGlobal10
	OpMaxLoadGlobal

	// Push the constant at N onto the stack.
	OpLoadConst0
	OpLoadConst1
	OpLoadConst2
	OpLoadConst3
	OpLoadConst4
	OpLoadConst5
	OpLoadConst6
	OpLoadConst7
	OpLoadConst8
	OpLoadConst9
	OpLoadConst10
	OpLoadConst11
	OpLoadConst12
	OpLoadConst13
	OpLoadConst14
	OpMaxLoadConst

	// The top of the stack must be a binding reference.
	// Push the constant at N onto the stack, and resolve the reference.
	OpLoadValue0
	OpLoadValue1
	OpLoadValue2
	OpLoadValue3
	OpLoadValue4
	OpLoadValue5
	OpLoadValue6
	OpLoadValue7
	OpLoadValue8
	OpLoadValue9
	OpLoadValue10
	OpMaxLoadValue

	// Pop the top value off the stack and push the field value of it onto the stack.
	// The field value is found at constant index N.
	OpPopLoadField0
	OpPopLoadField1
	OpPopLoadField2
	OpPopLoadField3
	OpPopLoadField4
	OpPopLoadField5
	OpPopLoadField6
	OpPopLoadField7
	OpPopLoadField8
	OpPopLoadField9
	OpPopLoadField10
	OpPopLoadField11
	OpPopLoadField12
	OpPopLoadField13
	OpPopLoadField14
	OpMaxPopLoadField

	OpPopLoadIndirect0
	OpMaxPopLoadIndirect

	// Do nothing. This bytecode can be used as a placeholder by bytecode optimizer.
	OpNoOp

	// Push unit onto the stack.
	OpLoadUnit

	// Push true onto the stack.
	OpLoadTrue

	// Push false onto the stack.
	OpLoadFalse

	// Push an integer N onto the stack.
	OpLoadInt0
	OpLoadInt1
	OpMaxLoadInt

	// Push an empty list onto the stack.
	OpLoadNil

	// The top of the stack must be a list.
	// Pop the top value off the stack, and push head of the list onto the stack.
	OpLoadHead

	// The top of the stack must be a list.
	// Pop the top value off the stack, and push tail of the list onto the stack.
	OpPopLoadTail

	// Push "None" of "option" type onto the stack.
	OpLoadNone

	// Pop the top value off the stack, and push "Some" of
	// "option" type with the popped value onto the stack.
	OpLoadSome

	// Pop the top value off the stack and store it in the local at N.
	OpStorePopLocal0
	OpStorePopLocal1
	OpStorePopLocal2
	OpStorePopLocal3
	OpStorePopLocal4
	OpStorePopLocal5
	OpStorePopLocal6
	OpMaxStorePopLocal

	OpStorePopGlobal0
	OpStorePopGlobal1
	OpStorePopGlobal2
	OpStorePopGlobal3
	OpStorePopGlobal4
	OpStorePopGlobal5
	OpStorePopGlobal6
	OpMaxStorePopGlobal

	// Pop the top value off the stack and store it in the field value of the top of the stack.
	// The field value is found at constant index N.
	OpStorePopField0
	OpStorePopField1
	OpStorePopField2
	OpStorePopField3
	OpStorePopField4
	OpStorePopField5
	OpStorePopField6
	OpMaxStorePopField

	OpReturn
	OpReturnUnit
	OpReturnTrue
	OpReturnFalse

	// Target of a backwards (looping) branch.
	// It is an error to have a backward branch whose target is not this instruction.
	OpLoopHead

	// Jump forward N+1 bytes.
	OpShortJump0
	OpShortJump1
	OpShortJump2
	OpShortJump3
	OpShortJump4
	OpMaxShortJump

	// The top of the stack must be boolean. Pop it off, and if false,
	// jump forward N bytes.
	OpShortBranchFalse0
	OpShortBranchFalse1
	OpShortBranchFalse2
	OpShortBranchFalse3
	OpShortBranchFalse4
	OpShortBranchFalse5
	OpShortBranchFalse6
	OpMaxShortBranchFalse

	// Apply a function with N arguments.
	// The function and arguments are all on the stack.
	OpApply1
	OpApply2
	OpApply3
	OpMaxApply

	// Apply a local N function with 1 argument.
	OpApplyDirect1_0
	OpApplyDirect1_1
	OpApplyDirect1_2
	OpApplyDirect1_3
	OpApplyDirect1_4
	OpApplyDirect1_5
	OpApplyDirect1_6
	OpApplyDirect1_7
	OpApplyDirect1_8
	OpApplyDirect1_10
	OpApplyDirect1_11
	OpApplyDirect1_12
	OpApplyDirect1_13
	OpApplyDirect1_14
	OpMaxApplyDirect1

	// Apply a local N function with 2 argument.
	OpApplyDirect2_0
	OpApplyDirect2_1
	OpApplyDirect2_2
	OpApplyDirect2_3
	OpApplyDirect2_4
	OpApplyDirect2_5
	OpApplyDirect2_6
	OpApplyDirect2_7
	OpApplyDirect2_8
	OpApplyDirect2_10
	OpApplyDirect2_11
	OpApplyDirect2_12
	OpApplyDirect2_13
	OpApplyDirect2_14
	OpMaxApplyDirect2

	// Apply a local N function with 3 argument.
	OpApplyDirect3_0
	OpApplyDirect3_1
	OpApplyDirect3_2
	OpApplyDirect3_3
	OpApplyDirect3_4
	OpApplyDirect3_5
	OpApplyDirect3_6
	OpMaxApplyDirect3

	OpEq
	OpNe
	OpLt
	OpLe
	OpGt
	OpGe
	OpAdd
	OpAdd1
	OpSub
	OpSub1
	OpMul
	OpDiv
	OpPow
	OpMod
	OpEqInts
	OpNeInts
	OpLtInts
	OpLeInts
	OpGtInts
	OpGeInts
	OpLshift
	OpRshift
	OpNegInt
	OpBnot
	OpBand
	OpBor
	OpBxor
	OpNot

	// The top of the stack must be an array or a list.
	// Push length of the top value onto the stack.
	OpCountValues

	OpLoadIndirect0
	OpLoadIndirect1
	OpLoadIndirect2
	OpMaxLoadIndirect

	OpPop
	OpSwapPop
	OpDup

	//
	// Two-byte codes
	//

	// Push onto the stack the value of the element B1%16 of
	// the local variable at B1/16.
	// The element of arrays are created by OpCreateArray.
	OpLoadLocalIndirect

	// Push the local variable found at B1 onto the stack.
	OpXLoadLocal

	OpXLoadGlobal

	// Push the constant at B1 onto the stack.
	OpXLoadConst

	// The top of the stack must be a binding reference.
	// Push the constant at B1 onto the stack, and resolve the reference.
	OpXLoadValue

	// Pop the top value off the stack and push the field value of it onto the stack.
	// The field value is found at constant index B1.
	OpXPopLoadField

	// Push a signed integer whose value is B1 onto the stack.
	OpXLoadInt

	OpXLoadIndirect

	// Pop the top value off the stack and store it in the local at B1.
	OpXStorePopLocal

	// Pop the top value off the stack and store it in the array B1%16 of the local at B1/16.
	OpStorePopLocalIndirect

	OpXStorePopGlobal

	// Pop the top value off the stack and store it in the field value of the top of the stack.
	// The field value is found at constant index B1.
	OpXStorePopField

	// Jump forward (or backward if the value is negative) N*256+B1 bytes.
	OpLongJump0
	OpLongJump1
	OpLongJump2
	OpLongJump3
	OpLongJump4
	OpLongJump5
	OpLongJump6
	OpMaxLongJump

	// The top of the stack must be boolean. Pop it off, and if false,
	// jump forward N*256+B1 bytes.
	OpLongBranchFalse0
	OpLongBranchFalse1
	OpLongBranchFalse2
	OpMaxLongBranchFalse

	// The top of the stack must be boolean. Pop it off, and if true,
	// jump forward N*256+B1 bytes.
	OpLongBranchTrue0
	OpLongBranchTrue1
	OpLongBranchTrue2
	OpMaxLongBranchTrue

	OpBranchNe
	OpBranchNeSizes

	// Apply a function with B1 arguments.
	// The function and arguments are all on the stack.
	OpXApply

	// Apply a local B1 function with N arguments on the stack.
	OpXApplyDirect1
	OpXApplyDirect2
	OpXApplyDirect3
	OpMaxXApplyDirect

	// Create a new array of size B1+1 and push it onto the stack.
	// The array represents:
	//
	// - a reference (size 1)
	// - a tuple (size is equal to the size of it)
	// - a list (size 2, head and tail)
	OpCreateArray

	// Replace the top B1+1 values on the stack with a new array
	// containg those values. Their order within the array is the same as
	// the order in which they were pushed on the stack.
	OpConsArray

	OpConsList

	OpCopyValues

	OpFullBlock

	//
	// Three-byte codes
	//

	OpXXLoadLocal
	OpXXLoadGlobal

	// Push onto the stack the contents of the reference of the local variable
	// at B1*256+B2.
	OpXLoadLocalIndirect

	// Push a signed integer whose value is B1*256+B2 onto the stack.
	OpXXLoadInt

	// Pop the top value off the stack and store it in the array B2 of the local at B1.
	OpXStorePopLocalIndirect

	// Apply a local B1 function with B2 arguments on the stack.
	OpXXApplyDirect

	OpCopyingBlock
	OpFullCopyingBlock

	// Call a primitive whose name is local B2. The number of arguments is B1.
	// If the call is success, push the return value onto the stack.
	OpPrimitive
)

var OpLongJumpBase = OpLongJump4
var OpLastShortCode = OpDup
var OpLastTwoByteCode = OpFullBlock
var OpLastThreeByteCode = OpPrimitive
var OpLastCode = OpLastThreeByteCode

func IntOfTwoBytes(b1 uint8, b2 uint8) int {
	return int(b1)*256 + int(b2)
}

func TwoBytesOfInt(v int) (uint8, uint8) {
	return uint8(v / 256), uint8(v % 256)
}

func TwoBytesOfNegInt(v int) (uint8, uint8) {
	return uint8(-(v/256 - 256)), uint8(v % 256)
}
