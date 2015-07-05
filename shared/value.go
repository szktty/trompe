package trompe

import (
	"bytes"
	"fmt"
	//"strings"
)

type Value interface{}

type Unit struct{}

type Tuple struct {
	Comps []Value
}

type List struct {
	Head Value
	Tail *List
}

type BlockClosure struct {
	Context *Context
	Code    *CompiledCode
	Copied  []Value
}

var UnitValue = &Unit{}
var NilValue = &List{}

func StringOfValue(v Value) string {
	switch desc := v.(type) {
	case nil:
		return "nil"
	case *Unit:
		return "()"
	case bool:
		if desc {
			return "true"
		} else {
			return "false"
		}
	case string:
		return fmt.Sprintf("\"%s\"", desc)
	case int:
		panic("must use int64")
	case int32: // character
		return fmt.Sprintf("%c", desc)
	case int64:
		return fmt.Sprintf("%d", desc)
	case float64:
		return fmt.Sprintf("%f", desc)
	case *Module:
		return fmt.Sprintf("<Module \"%s\">", desc.Name)
	case *CompiledCode:
		return fmt.Sprintf("<CompiledCode %p>", v)
	case *BlockClosure:
		return fmt.Sprintf("<BlockClosure %p>", v)
	case *PrimRef:
		return fmt.Sprintf("<primitive \"%s\">", desc.Key)
	case *NamePath:
		return fmt.Sprintf("{%s}", desc.String())
	case *Tuple:
		buf := bytes.NewBufferString("(")
		for i, comp := range desc.Comps {
			buf.WriteString(StringOfValue(comp))
			if i+1 < len(desc.Comps) {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(")")
		return buf.String()
	case *List:
		buf := bytes.NewBufferString("[")
		l := desc
		for l != NilValue {
			buf.WriteString(StringOfValue(l.Head))
			if l.Tail != NilValue {
				buf.WriteString("; ")
			}
			l = l.Tail
		}
		buf.WriteString("]")
		return buf.String()
	case []Value:
		buf := bytes.NewBufferString("#[")
		for i, v := range desc {
			buf.WriteString(StringOfValue(v))
			if i+1 < len(desc) {
				buf.WriteString("; ")
			}
		}
		buf.WriteString("]")
		return buf.String()
	case Primitive:
		return fmt.Sprintf("<primitive %p>", desc)
	default:
		panic(fmt.Errorf("unknown value %s", v))
	}
}

func NewTuple(comps []Value) *Tuple {
	return &Tuple{Comps: comps}
}

func (t *Tuple) Get(i int) Value {
	return t.Comps[i]
}

func NewList(v Value) *List {
	if v == nil {
		panic("cannot allow nil")
	}
	return &List{Head: v, Tail: NilValue}
}

func NewListFromArray(es []Value) *List {
	l := NilValue
	for _, e := range es {
		l = l.Cons(e)
	}
	return l.Rev()
}

func (l *List) IsNilValue() bool {
	return l == NilValue || l.Length() == 0
}

func (l *List) Length() int {
	v := 0
	for l != NilValue {
		v++
		l = l.Tail
	}
	return v
}

func (l *List) Cons(v Value) *List {
	return &List{Head: v, Tail: l}
}

func (l *List) Rev() *List {
	rev := NilValue
	for l != NilValue {
		rev = rev.Cons(l.Head)
		l = l.Tail
	}
	return rev
}

func NewBlockClosure(code *CompiledCode) *BlockClosure {
	return &BlockClosure{Code: code}
}

type ComparisonResult int

const (
	_ComparisonResult = iota
	OrderedSame
	OrderedAscending
	OrderedDescending
)

func CompareValues(left Value, right Value) (ComparisonResult, bool) {
	switch ldesc := left.(type) {
	case bool:
		if rdesc, ok := right.(bool); ok {
			if ldesc == rdesc {
				return OrderedSame, true
			} else if ldesc {
				return OrderedAscending, true
			} else {
				return OrderedDescending, true
			}
		}
	case string:
		if rdesc, ok := right.(string); ok {
			if ldesc == rdesc {
				return OrderedSame, true
			} else if ldesc > rdesc {
				return OrderedAscending, true
			} else {
				return OrderedDescending, true
			}
		}
	case int32:
		if rdesc, ok := right.(int32); ok {
			if ldesc == rdesc {
				return OrderedSame, true
			} else if ldesc > rdesc {
				return OrderedAscending, true
			} else {
				return OrderedDescending, true
			}
		}
	case int64:
		if rdesc, ok := right.(int64); ok {
			if ldesc == rdesc {
				return OrderedSame, true
			} else if ldesc > rdesc {
				return OrderedAscending, true
			} else {
				return OrderedDescending, true
			}
		}
	case float64:
		if rdesc, ok := right.(float64); ok {
			if ldesc == rdesc {
				return OrderedSame, true
			} else if ldesc > rdesc {
				return OrderedAscending, true
			} else {
				return OrderedDescending, true
			}
		}
	}
	Panicf("unknown types: %s and %s", left, right)
	return 0, false
}
