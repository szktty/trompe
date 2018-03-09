package trompe

const (
	ValUnitType = iota
	ValBoolType
	ValIntType
	ValStrType
	ValListType
	ValTupleType
	ValClosType
	ValOptType
)

/*
type Value interface {
	Type() int
	Bool() bool
	Int() int
	String() string
	Closure() *Closure
}
*/
type Value interface {
	Bool() bool
	Int() int
	String() string
}

type ValUnit struct{}

type ValBool struct {
	Value bool
}

type ValInt struct {
	Value int
}

type ValStr struct {
	Value string
}

type Closure struct {
	Lits []Value
	Ops  []Opcode
}

type ValOpt struct {
	Value *Value // nullable
}

var SharedUnit = &ValUnit{}
var SharedTrue = &ValBool{true}
var SharedFalse = &ValBool{false}

func (val *ValUnit) Bool() bool {
	panic("unit")
}

func (val *ValUnit) Int() int {
	panic("unit")
}

func (val *ValUnit) String() string {
	panic("unit")
}

func (val *ValBool) Bool() bool {
	return val.Value
}

func (val *ValBool) Int() int {
	panic("bool")
}

func (val *ValBool) String() string {
	panic("bool")
}

func (val *ValInt) Bool() bool {
	panic("int")
}

func (val *ValInt) Int() int {
	return val.Value
}

func (val *ValInt) String() string {
	panic("int")
}

func (val *ValStr) Bool() bool {
	panic("string")
}

func (val *ValStr) Int() int {
	panic("string")
}

func (val *ValStr) String() string {
	return val.Value
}

type ValList struct {
	Value *List
}

func (val *ValList) Bool() bool {
	panic("list")
}

func (val *ValList) Int() int {
	panic("list")
}

func (val *ValList) String() string {
	panic("list")
}

func CreateValList(value *List) *ValList {
	return &ValList{value}
}
