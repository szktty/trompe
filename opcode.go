package trompe

type Opcode = int

const (
	OpNop      = iota
	OpPosition // position
	OpLoadUnit
	OpLoadTrue
	OpLoadFalse
	OpLoadInt // int
	OpLoadSome
	OpLoadNone
	OpLoadLiteral // index of value in literal list *)
	OpLoadLocal   // index
	OpStore       // index of local
	OpPop
	OpReturn
	OpLabel // name
	OpLoopHead
	OpJump        // index
	OpBranchTrue  // index
	OpBranchFalse // index
	OpCall        // length
	OpPrimitive   // index of string as literal
	OpList        // length
	OpTuple       // length
)
