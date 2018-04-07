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
	ValRangeType
	ValPtnType
	ValModRefType
)

type Value interface {
	Type() int
	Desc() string
	Bool() bool
	Int() int
	String() string
	List() *List
	Tuple() []Value
	Closure() Closure
	// Pattern() Pattern
	Iter() ValIter
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

type ValTuple struct {
	Values []Value
}

type ValIter interface {
	Value
	Next() Value
}

func NewValIter(val Value) ValIter {
	switch val := val.(type) {
	case *Range:
		return val.NewValIter()
	default:
		panic("unsupported")
	}
}

var LangUnit = &ValUnit{}
var LangTrue = &ValBool{true}
var LangFalse = &ValBool{false}
var LangNone = &ValOpt{nil}

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

func (val *ValUnit) List() *List {
	panic("unit")
}

func (val *ValUnit) Tuple() []Value {
	panic("unit")
}

func (val *ValUnit) Iter() ValIter {
	return nil
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

func (val *ValBool) List() *List {
	panic("bool")
}

func (val *ValBool) Tuple() []Value {
	panic("bool")
}

func (val *ValBool) Iter() ValIter {
	return nil
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

func (val *ValInt) List() *List {
	panic("int")
}

func (val *ValInt) Tuple() []Value {
	panic("int")
}

func (val *ValInt) Iter() ValIter {
	return nil
}

func (val *ValStr) Type() int {
	return ValStrType
}

func (val *ValStr) Desc() string {
	return val.Value
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

func NewValStr(value string) *ValStr {
	return &ValStr{value}
}

func (val *ValStr) Closure() Closure {
	panic("not closure")
}

func (val *ValStr) List() *List {
	panic("string")
}

func (val *ValStr) Tuple() []Value {
	panic("string")
}

func (val *ValStr) Iter() ValIter {
	return nil
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

func (val *ValList) List() *List {
	return val.Value
}

func (val *ValList) Tuple() []Value {
	panic("list")
}

func (val *ValList) Iter() ValIter {
	// TODO
	return nil
}

func NewValList(value *List) *ValList {
	return &ValList{value}
}

func NewValTuple(value ...Value) *ValTuple {
	return &ValTuple{value}
}

func (val *ValTuple) Type() int {
	return ValTupleType
}

func (val *ValTuple) Desc() string {
	// TODO
	return fmt.Sprintf("list")
}

func (val *ValTuple) Bool() bool {
	panic("tuple")
}

func (val *ValTuple) Int() int {
	panic("tuple")
}

func (val *ValTuple) String() string {
	panic("tuple")
}

func (val *ValTuple) Closure() Closure {
	panic("not closure")
}

func (val *ValTuple) List() *List {
	panic("tuple")
}

func (val *ValTuple) Tuple() []Value {
	return val.Values
}

func (val *ValTuple) Iter() ValIter {
	return nil
}

type ValOpt struct {
	Value Value // nullable
}

func NewValOpt(value Value) *ValOpt {
	return &ValOpt{value}
}

func (val *ValOpt) Type() int {
	return ValOptType
}

func (val *ValOpt) Desc() string {
	return fmt.Sprintf("closure %p", val.Value)
}

func (val *ValOpt) Bool() bool {
	panic("Opt")
}

func (val *ValOpt) Int() int {
	panic("Opt")
}

func (val *ValOpt) String() string {
	panic("Opt")
}

func (val *ValOpt) Closure() Closure {
	panic("Opt")
}

func (val *ValOpt) List() *List {
	panic("Opt")
}

func (val *ValOpt) Tuple() []Value {
	panic("Opt")
}

func (val *ValOpt) Iter() ValIter {
	return nil
}

type ValPtn struct {
	Value Pattern
}

func NewValPtn(ptn Pattern) *ValPtn {
	return &ValPtn{ptn}
}

func (val *ValPtn) Type() int {
	return ValPtnType
}

func (val *ValPtn) Desc() string {
	return val.Value.Desc()
}

func (val *ValPtn) Bool() bool {
	panic("Pattern")
}

func (val *ValPtn) Int() int {
	panic("Pattern")
}

func (val *ValPtn) String() string {
	panic("Pattern")
}

func (val *ValPtn) Closure() Closure {
	panic("Pattern")
}

func (val *ValPtn) List() *List {
	panic("Pattern")
}

func (val *ValPtn) Tuple() []Value {
	panic("Pattern")
}

func (val *ValPtn) Iter() ValIter {
	return nil
}

type ValModRef struct {
	Path  string
	Cache *Module
}

func NewValModRef(path string) *ValModRef {
	return &ValModRef{Path: path}
}

func NewValModRefWithModule(m *Module) *ValModRef {
	return &ValModRef{Cache: m}
}

func (val *ValModRef) Type() int {
	return ValModRefType
}

func (val *ValModRef) Desc() string {
	return fmt.Sprintf("{%s}", val.Path)
}

func (val *ValModRef) Bool() bool {
	panic("ModRef")
}

func (val *ValModRef) Int() int {
	panic("ModRef")
}

func (val *ValModRef) String() string {
	panic("ModRef")
}

func (val *ValModRef) Closure() Closure {
	panic("ModRef")
}

func (val *ValModRef) List() *List {
	panic("ModRef")
}

func (val *ValModRef) Tuple() []Value {
	panic("ModRef")
}

func (val *ValModRef) Module() *Module {
	if val.Cache == nil {
		val.Cache = GetModule(val.Path)
	}
	return val.Cache
}

func (val *ValModRef) Iter() ValIter {
	return nil
}
