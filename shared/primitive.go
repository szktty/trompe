package trompe

import ()

type Primitive func(*State, *Context, []Value) (Value, error)

type PrimRef struct {
	Key string
}

var PrimMap = map[string]Primitive{
	"create_module": Prim_create_module,
	"array_get":     Prim_array_get,
	"array_set":     Prim_array_set,
}

var PrimNameList = make([]string, 0)

func AddPrim(name string, prim Primitive) {
	PrimMap[name] = prim
	PrimNameList = append(PrimNameList, name)
}

func PrimIndex(name string) (int, bool) {
	for i, k := range PrimNameList {
		if k == name {
			return i, true
		}
	}
	return -1, false
}

func NewPrimRef(key string) *PrimRef {
	return &PrimRef{Key: key}
}

func Prim_create_module(state *State, parent *Context, args []Value) (Value, error) {
	ary := args[0].([]Value)
	name := ary[0].(string)
	Debugf("create module %s", name)
	return NewModule(name), nil
}

func Prim_array_get(state *State, parent *Context, args []Value) (Value, error) {
	ary := args[0].([]Value)
	idx := args[1].(int64)
	return ary[idx], nil
}

func Prim_array_set(state *State, parent *Context, args []Value) (Value, error) {
	ary := args[0].([]Value)
	idx := args[1].(int64)
	ary[idx] = args[2]
	return UnitValue, nil
}
