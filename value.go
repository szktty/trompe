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
	ValPrimType
	ValOptType
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
	//Prim() string
	// Pattern() Pattern
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

func CreateValStr(value string) *ValStr {
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

func CreateValList(value *List) *ValList {
	return &ValList{value}
}

func CreateValTuple(value ...Value) *ValTuple {
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

func (val *ValClos) List() *List {
	panic("Clos")
}

func (val *ValClos) Tuple() []Value {
	panic("closure")
}

func CreateValClos(value Closure) *ValClos {
	return &ValClos{value}
}

type ValPrim struct {
	Name string
}

func NewValPrim(name string) *ValPrim {
	return &ValPrim{name}
}

func (val *ValPrim) Type() int {
	return ValPrimType
}

func (val *ValPrim) Desc() string {
	return fmt.Sprintf("prim %s", val.Name)
}

func (val *ValPrim) Bool() bool {
	panic("Prim")
}

func (val *ValPrim) Int() int {
	panic("Prim")
}

func (val *ValPrim) String() string {
	panic("Prim")
}

func (val *ValPrim) Closure() Closure {
	return GetPrim(val.Name)
}

func (val *ValPrim) List() *List {
	panic("Prim")
}

func (val *ValPrim) Tuple() []Value {
	panic("Prim")
}

type ValOpt struct {
	Value Value // nullable
}

func CreateValOpt(value Value) *ValOpt {
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

type ValPtn struct {
	Value Pattern
}

func CreateValPtn(ptn Pattern) *ValPtn {
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
