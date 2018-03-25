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
	OpLoadLit   // index of value in literal list
	OpLoadLocal // index
	OpLoadAttr  // index of literal string
	OpLoadPrim  // index of literal string
	OpLoadArg   // index
	OpLoadModule
	OpStoreLocal // index of literal string
	OpStoreRef
	OpStoreAttr // index of literal string
	OpPop
	OpDup
	OpReturn
	OpReturnUnit
	OpLabel       // label number
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

func GetOpName(op int) string {
	switch op {
	case OpNop:
		return "OpNop"
	case OpLoadUnit:
		return "OpLoadUnit"
	case OpLoadTrue:
		return "OpLoadTrue"
	case OpLoadFalse:
		return "OpLoadFalse"
	case OpLoadZero:
		return "OpLoadZero"
	case OpLoadOne:
		return "OpLoadOne"
	case OpLoadNegOne:
		return "OpLoadNegOne"
	case OpLoadInt:
		return "OpLoadInt"
	case OpLoadNone:
		return "OpLoadNone"
	case OpLoadRef:
		return "OpLoadRef"
	case OpLoadLit:
		return "OpLoadLit"
	case OpLoadLocal:
		return "OpLoadLocal"
	case OpLoadAttr:
		return "OpLoadAttr"
	case OpLoadPrim:
		return "OpLoadPrim"
	case OpLoadArg:
		return "OpLoadArg"
	case OpStoreLocal:
		return "OpStoreLocal"
	case OpStoreRef:
		return "OpStoreRef"
	case OpStoreAttr:
		return "OpStoreAttr"
	case OpPop:
		return "OpPop"
	case OpReturn:
		return "OpReturn"
	case OpReturnUnit:
		return "OpReturnUnit"
	case OpLabel:
		return "OpLabel"
	case OpJump:
		return "OpJump"
	case OpBranchTrue:
		return "OpBranchTrue"
	case OpBranchFalse:
		return "OpBranchFalse"
	case OpBegin:
		return "OpBegin"
	case OpEnd:
		return "OpEnd"
	case OpCall:
		return "OpCall"
	case OpEq:
		return "OpEq"
	case OpNe:
		return "OpNe"
	case OpLt:
		return "OpLt"
	case OpLe:
		return "OpLe"
	case OpGt:
		return "OpGt"
	case OpGe:
		return "OpGe"
	case OpMatch:
		return "OpMatch"
	case OpAdd:
		return "OpAdd"
	case OpSub:
		return "OpSub"
	case OpMul:
		return "OpMul"
	case OpDiv:
		return "OpDiv"
	case OpMod:
		return "OpMod"
	case OpSome:
		return "OpSome"
	case OpList:
		return "OpList"
	case OpTuple:
		return "OpTuple"
	default:
		panic("unknown opcode")
	}
}
