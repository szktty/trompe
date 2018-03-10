package trompe

type Opcode = int

const (
	OpNop = iota
	OpLoadUnit
	OpLoadTrue
	OpLoadFalse
	OpLoadInt // int
	OpLoadNone
	OpLoadRef
	OpLoadLit   // index of value in literal list *)
	OpLoadLocal // index
	OpLoadAttr  // index of literal string
	OpLoadPrim  // index of literal string
	OpStore     // index of literal string
	OpStoreRef
	OpStoreAttr // index of literal string
	OpPop
	OpReturn
	OpLabel       // name
	OpJump        // index
	OpBranchTrue  // index
	OpBranchFalse // index
	OpCall        // length
	OpSome
	OpList  // length
	OpTuple // length
)
