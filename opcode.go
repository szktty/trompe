package trompe

type Opcode = int

const (
	OpNop = iota
	OpLoadUnit
	OpLoadTrue
	OpLoadFalse
	OpLoadInt // int
	OpLoadSome
	OpLoadNone
	OpLoadLit   // index of value in literal list *)
	OpLoadLocal // index
	OpLoadAttr  // index of literal string
	OpStore     // index of local
	OpStoreAttr // index of literal string
	OpPop
	OpReturn
	OpLabel       // name
	OpJump        // index
	OpBranchTrue  // index
	OpBranchFalse // index
	OpCall        // length
	OpPrimitive   // index of literal string
	OpList        // length
	OpTuple       // length
)
