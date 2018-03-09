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
	OpStore     // index of local
	OpPop
	OpReturn
	OpLabel       // name
	OpJump        // index
	OpBranchTrue  // index
	OpBranchFalse // index
	OpCall        // length
	OpPrimitive   // index of string as literal
	OpList        // length
	OpTuple       // length
)
