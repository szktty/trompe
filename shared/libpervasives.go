package trompe

import (
	"bytes"
	"fmt"
	"math"
)

func ModulePervasives() *Module {
	m := NewModule("Pervasives")

	m.SetTycon("list", TcList)
	m.SetTycon("option", TcOption)
	m.SetTycon("result", TcResult)

	m.SetFieldType("=",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(TVar("a"), TVar("a"), TBool))),
				TArgs(TVar("a")))))
	m.SetPrim("=", Pervasives_eq)

	m.SetFieldType("min",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(TVar("a"), TVar("a"), TVar("a")))),
				TArgs(TVar("a")))))
	m.SetPrim("min", Pervasives_min)

	m.SetFieldType("max",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(TVar("a"), TVar("a"), TVar("a")))),
				TArgs(TVar("a")))))
	m.SetPrim("max", Pervasives_max)

	m.SetFieldType("show",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(TVar("a"), TUnit))),
				TArgs(TVar("a")))))
	m.SetPrim("show", Pervasives_show)

	m.SetFieldType("print_string", TApp(TcArrow, TArgs(TString, TUnit)))
	m.SetPrim("print_string", Pervasives_print_string)

	m.SetFieldType("print_char", TApp(TcArrow, TArgs(TChar, TUnit)))
	m.SetPrim("print_char", Pervasives_print_char)

	m.SetFieldType("print_int", TApp(TcArrow, TArgs(TInt, TUnit)))
	m.SetPrim("print_int", Pervasives_print_int)

	m.SetFieldType("print_float", TApp(TcArrow, TArgs(TFloat, TUnit)))
	m.SetPrim("print_float", Pervasives_print_float)

	m.SetFieldType("print_nl", TApp(TcArrow, TArgs(TUnit, TUnit)))
	m.SetPrim("print_nl", Pervasives_print_nl)

	m.SetFieldType("printf",
		TApp(TcArrow, TArgs(
			TApp(TcFormat, TArgs(TUnit, TUnit)),
			TApp(TcFormatter, TArgs(TUnit, TUnit)))))
	m.SetPrim("printf", Pervasives_printf)

	m.SetFieldType("sprintf",
		TApp(TcArrow, TArgs(
			TApp(TcFormat, TArgs(TUnit, TUnit)),
			TApp(TcFormatter, TArgs(TUnit, TString)))))
	m.SetPrim("sprintf", Pervasives_sprintf)

	m.SetFieldType("raise",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(TExn, TVar("a")))),
				TArgs(TVar("a")))))

	return m
}

/*
external printf : (unit, unit) format -> (unit, unit) formatter = "printf"
external sprintf : (unit, string) format -> (unit, string) formatter = "sprintf"
external raise : exn -> 'a = "raise"
*/

// val (=) : 'a -> 'a -> bool
func Pervasives_eq(state *State, parent *Context, args []Value) (Value, error) {
	// TODO
	x := args[0]
	y := args[1]
	return x == y, nil
}

// val min : 'a -> 'a -> 'a
func Pervasives_min(state *State, parent *Context, args []Value) (Value, error) {
	// TODO
	x := args[0]
	y := args[1]
	if xv, ok := x.(int64); ok {
		if yv, ok := y.(int64); ok {
			if xv < yv {
				return xv, nil
			} else {
				return yv, nil
			}
		}
	} else if xv, ok := x.(float64); ok {
		if yv, ok := y.(float64); ok {
			return math.Min(xv, yv), nil
		}
	}
	return x, nil
}

// val max : 'a -> 'a -> 'a
func Pervasives_max(state *State, parent *Context, args []Value) (Value, error) {
	// TODO
	x := args[0]
	y := args[1]
	if xv, ok := x.(int64); ok {
		if yv, ok := y.(int64); ok {
			if xv > yv {
				return xv, nil
			} else {
				return yv, nil
			}
		}
	} else if xv, ok := x.(float64); ok {
		if yv, ok := y.(float64); ok {
			return math.Max(xv, yv), nil
		}
	}
	return x, nil
}

func Pervasives_show(state *State, parent *Context, args []Value) (Value, error) {
	switch v := args[0].(type) {
	case int64:
		fmt.Printf("%d", v)
	case float64:
		fmt.Printf("%f", v)
	case string:
		fmt.Printf("%s", v)
	default:
		fmt.Printf("%s", v)
	}
	fmt.Printf("\n")
	return UnitValue, nil
}

func Pervasives_print_char(state *State, parent *Context, args []Value) (Value, error) {
	fmt.Printf("%s", args[0].(string))
	return UnitValue, nil
}

func Pervasives_print_string(state *State, parent *Context, args []Value) (Value, error) {
	fmt.Printf("%s", args[0].(string))
	return UnitValue, nil
}

func Pervasives_print_int(state *State, parent *Context, args []Value) (Value, error) {
	fmt.Printf("%d", args[0].(int64))
	return UnitValue, nil
}

func Pervasives_print_float(state *State, parent *Context, args []Value) (Value, error) {
	fmt.Printf("%f", args[0].(float64))
	return UnitValue, nil
}

func Pervasives_print_nl(state *State, parent *Context, args []Value) (Value, error) {
	fmt.Printf("\n")
	return UnitValue, nil
}

func Pervasives_printf(state *State, parent *Context, args []Value) (Value, error) {
	s, err := Pervasives_sprintf(state, parent, args)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", s)
	return UnitValue, nil
}

func Pervasives_sprintf(state *State, parent *Context, args []Value) (Value, error) {
	return FormatString(args)
}

func FormatString(args []Value) (string, error) {
	s := args[0].(string)
	f, err := NewFormat(s)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	i := 1
	for _, comp := range f.Comps {
		switch comp.Type {
		case FormatTypeNone:
			buf.WriteString(comp.String)
			continue
		case FormatTypeBool:
			if args[i].(bool) {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		case FormatTypeInt:
			buf.WriteString(fmt.Sprintf("%d", args[i].(int64)))
		case FormatTypeFloat:
			buf.WriteString(fmt.Sprintf("%f", args[i].(float64)))
		case FormatTypeString:
			buf.WriteString(args[i].(string))
		case FormatTypeChar:
			buf.WriteString(fmt.Sprintf("%c", args[i].(string)))
		default:
			Panicf("unsupported format type %d", comp.Type)
		}
		i++
	}
	return buf.String(), nil
}
