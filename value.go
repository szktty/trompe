package trompe

import (
	"fmt"
	"strconv"
)

const (
	ValueTypeUnit = iota
	ValueTypeBool
	ValueTypeInt
	ValueTypeString
	ValueTypeList
	ValueTypeTuple
	ValueTypeClos
	ValueTypeOption
	ValueTypeRange
	ValueTypeIter
	ValueTypePattern
	ValueTypeRef
)

type Value interface {
	Type() int
	Desc() string
}

type Unit struct{}

type Bool struct {
	Value bool
}

type Int struct {
	Value int
}

type String struct {
	Value string
}

type Tuple struct {
	Values []Value
}

var SharedUnit = &Unit{}
var SharedTrue = &Bool{true}
var SharedFalse = &Bool{false}
var SharedNone = &Option{nil}

func NewBool(b bool) *Bool {
	if b {
		return SharedTrue
	} else {
		return SharedFalse
	}
}

func ValueToBool(v Value) (*Bool, bool) {
	switch v := v.(type) {
	case *Bool:
		return v, true
	default:
		return nil, false
	}
}

func ValueToInt(v Value) (*Int, bool) {
	switch v := v.(type) {
	case *Int:
		return v, true
	default:
		return nil, false
	}
}

func ValueToString(v Value) (*String, bool) {
	switch v := v.(type) {
	case *String:
		return v, true
	default:
		return nil, false
	}
}

func ValueToList(v Value) (*List, bool) {
	switch v := v.(type) {
	case *List:
		return v, true
	default:
		return nil, false
	}
}

func ValueToTuple(v Value) (*Tuple, bool) {
	switch v := v.(type) {
	case *Tuple:
		return v, true
	default:
		return nil, false
	}
}

func ValueToClos(v Value) (Closure, bool) {
	switch v := v.(type) {
	case *CompiledCode:
		return v, true
	case *Primitive:
		return v, true
	default:
		return nil, false
	}
}

func ValueToOption(v Value) (*Option, bool) {
	switch v := v.(type) {
	case *Option:
		return v, true
	default:
		return nil, false
	}
}

func ValueToIter(v Value) (Iter, bool) {
	switch v := v.(type) {
	case Iter:
		return v, true
	default:
		return nil, false
	}
}

func ValueToRef(v Value) (*Ref, bool) {
	switch v := v.(type) {
	case *Ref:
		return v, true
	default:
		return nil, false
	}
}

func (u *Unit) Type() int {
	return ValueTypeUnit
}

func (u *Unit) Desc() string {
	return "()"
}

func (b *Bool) Type() int {
	return ValueTypeBool
}

func (b *Bool) Desc() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

func NewInt(i int) *Int {
	return &Int{i}
}

func (i *Int) Type() int {
	return ValueTypeInt
}

func (i *Int) Desc() string {
	return fmt.Sprintf("%d", i.Value)
}

func NewString(s string) *String {
	return &String{s}
}

func (s *String) Type() int {
	return ValueTypeString
}

func (s *String) Desc() string {
	return s.Value
}

func NewTuple(v ...Value) *Tuple {
	return &Tuple{v}
}

func (t *Tuple) Type() int {
	return ValueTypeTuple
}

func (t *Tuple) Desc() string {
	// TODO
	return fmt.Sprintf("tuple")
}

func (t *Tuple) Len() int {
	return len(t.Values)
}

type Ref struct {
	Path   string
	module *Module
}

func NewRef(path string, mod *Module) *Ref {
	return &Ref{Path: path, module: mod}
}

func (r *Ref) Type() int {
	return ValueTypeRef
}

func (r *Ref) Desc() string {
	return fmt.Sprintf("{%s}", r.Path)
}

func (r *Ref) Module() *Module {
	if r.module == nil {
		r.module = GetModule(r.Path)
	}
	return r.module
}

func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
