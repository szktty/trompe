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
}
*/
type Value interface {
	Bool() bool
	Int() int
	String() string
	Closure() Closure
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

type ValOpt struct {
	Value *Value // nullable
}

var SharedValUnit = &ValUnit{}
var SharedValTrue = &ValBool{true}
var SharedValFalse = &ValBool{false}

func (val *ValUnit) Bool() bool {
	panic("unit")
}

func (val *ValUnit) Int() int {
	panic("unit")
}

func (val *ValUnit) String() string {
	panic("unit")
}

func (val *ValUnit) Closure() Closure {
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

func (val *ValBool) Closure() Closure {
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

func (val *ValInt) Closure() Closure {
	panic("not closure")
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

func (val *ValStr) Closure() Closure {
	panic("not closure")
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

func (val *ValList) Closure() Closure {
	panic("not closure")
}

func CreateValList(value *List) *ValList {
	return &ValList{value}
}

type ValClos struct {
	Value Closure
}

func (val *ValClos) Bool() bool {
	panic("Clos")
}

func (val *ValClos) Int() int {
	panic("Clos")
}

func (val *ValClos) String() string {
	panic("Clos")
}

func (val *ValClos) Closure() Closure {
	return val.Value
}

func CreateValClos(value Closure) *ValClos {
	return &ValClos{value}
}
