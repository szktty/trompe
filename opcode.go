package trompe

type Opcode = int

const (
	OpNop = iota
	OpLoadUnit
	OpLoadTrue
	OpLoadFalse
	OpLoadZero   // load 0
	OpLoadOne    // load 1
	OpLoadNegOne // load -1
	OpLoadInt    // int
	OpLoadNone
	OpLoadRef
	OpLoadLit   // index of value in literal list *)
	OpLoadLocal // index
	OpLoadAttr  // index of literal string
	OpLoadPrim  // index of literal string
	OpLoadArg   // index
	OpStore     // index of literal string
	OpStoreRef
	OpStoreAttr // index of literal string
	OpPop
	OpReturn
	OpReturnUnit
	OpLabel // name
	OpLoopHead
	OpJump        // index
	OpBranchTrue  // index
	OpBranchFalse // index
	OpBegin
	OpEnd
	OpCall // length
	OpEq
	OpNe
	OpLt
	OpLe
	OpGt
	OpGe
	OpMatch
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpSome
	OpList  // length
	OpTuple // length
)
