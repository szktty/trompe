package trompe

import "fmt"

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

type Value interface {
	Type() int
	Desc() string
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

func (val *ValUnit) Type() int {
	return ValUnitType
}

func (val *ValUnit) Desc() string {
	return "()"
}

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

func (val *ValBool) Type() int {
	return ValBoolType
}

func (val *ValBool) Desc() string {
	if val.Value {
		return "true"
	} else {
		return "false"
	}
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

func (val *ValInt) Type() int {
	return ValIntType
}

func (val *ValInt) Desc() string {
	return fmt.Sprintf("%d", val.Value)
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

func (val *ValStr) Type() int {
	return ValStrType
}

func (val *ValStr) Desc() string {
	return fmt.Sprintf("\"%s\"", val.Value)
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

func CreateValStr(value string) *ValStr {
	return &ValStr{value}
}

func (val *ValStr) Closure() Closure {
	panic("not closure")
}

type ValList struct {
	Value *List
}

func (val *ValList) Type() int {
	return ValListType
}

func (val *ValList) Desc() string {
	// TODO
	return fmt.Sprintf("list")
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

func (val *ValClos) Type() int {
	return ValClosType
}

func (val *ValClos) Desc() string {
	return fmt.Sprintf("closure %p", val.Value)
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
