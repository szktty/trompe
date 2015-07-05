package trompe

import (
	"bytes"
	"fmt"
)

type Format struct {
	Type   *TypeApp // format type. ex: (unit, string) format
	String string
	Comps  []*FormatComp
}

type FormatComp struct {
	Type   int
	String string
}

const (
	_FormatType = iota
	FormatTypeNone
	FormatTypeRepr // @
	FormatTypeInt
	FormatTypeUint
	FormatTypeLowerHexUint
	FormatTypeUpperHexUint
	FormatTypeOctalUint
	FormatTypeFloat
	FormatTypeString
	FormatTypeChar
	FormatTypeBool
)

func NewFormat(s string) (*Format, error) {
	comps := make([]*FormatComp, 0)
	flag := false
	buf := bytes.NewBuffer(nil)
	var comp *FormatComp
	for _, c := range s {
		if flag {
			switch c {
			case '%':
				comp = &FormatComp{Type: FormatTypeString, String: "%"}
			case '@':
				comp = &FormatComp{Type: FormatTypeRepr}
			case 'B':
				comp = &FormatComp{Type: FormatTypeBool}
			case 'd':
				comp = &FormatComp{Type: FormatTypeInt}
			case 'u':
				comp = &FormatComp{Type: FormatTypeUint}
			case 's':
				comp = &FormatComp{Type: FormatTypeString}
			case 'x':
				comp = &FormatComp{Type: FormatTypeLowerHexUint}
			case 'X':
				comp = &FormatComp{Type: FormatTypeUpperHexUint}
			default:
				return nil, fmt.Errorf("invalid flag %%%c", c)
			}
			comps = append(comps, comp)
			flag = false
		} else {
			switch c {
			case '%':
				flag = true
				comps = append(comps,
					&FormatComp{Type: FormatTypeNone, String: buf.String()})
				buf.Reset()
			default:
				buf.WriteRune(rune(c))
			}
		}
	}
	if buf.Len() > 0 {
		comps = append(comps,
			&FormatComp{Type: FormatTypeNone, String: buf.String()})
	}
	return &Format{String: s, Comps: comps}, nil
}

func (f *Format) SetType(ty Type) error {
	ty1, ok := DerefType(ty)
	if !ok {
		return fmt.Errorf("invalid formatter")
	}
	app, ok := ty1.(*TypeApp)
	if !ok {
		return fmt.Errorf("invalid formatter")
	}
	f.Type = app
	return nil
}

func StringOfFormat(f *Format) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("\"")
	for _, comp := range f.Comps {
		switch comp.Type {
		case FormatTypeNone:
			buf.WriteString(comp.String)
		case FormatTypeInt:
			// TODO: escape special characters
			buf.WriteString("%d")
		default:
			panic("unsupported format type")
		}
	}
	buf.WriteString("\"")
	return buf.String()
}

func (f *Format) Formatter() (*TypeApp, error) {
	// pick return type
	if len(f.Type.Args) != 2 {
		return nil, fmt.Errorf("formatter parameter must be 2 (given %d)",
			len(f.Type.Args))
	}
	ret := f.Type.Args[1]

	args := make([]Type, 0)
	for _, comp := range f.Comps {
		switch comp.Type {
		case FormatTypeNone:
			break
		case FormatTypeBool:
			args = append(args, TBool)
		case FormatTypeInt:
			args = append(args, TInt)
		case FormatTypeFloat:
			args = append(args, TFloat)
		case FormatTypeString:
			args = append(args, TString)
		case FormatTypeChar:
			args = append(args, TChar)
		default:
			Panicf("unsupported format type %d", comp.Type)
		}
	}
	args = append(args, ret)
	return &TypeApp{Tycon: &TyconArrow{}, Args: args}, nil
}

func translateFormat(fun *TypeApp, i int, f string) (*TypeApp, error) {
	if fun.Tycon.TyconTag() != TyconTagArrow {
		return nil, fmt.Errorf("type including format and formatter must be a function type: %s", StringOfType(fun))
	}

	fter := fun.Args[i]
	format, err := NewFormat(f)
	if err != nil {
		return nil, err
	}
	err = format.SetType(fter)
	if err != nil {
		return nil, err
	}

	fter1, err := format.Formatter()
	if err != nil {
		return nil, err
	}

	// replace format and formatter types with expansions
	args := make([]Type, len(fun.Args))
	for i := 0; i < len(fun.Args); i++ {
		ty, _ := DerefType(fun.Args[i])
		args[i] = ty
		if app, ok := ty.(*TypeApp); ok {
			switch app.Tycon.(type) {
			case *TyconFormat:
				args[i] = TString
			case *TyconFormatter:
				args[i] = fter1
			}
		}
	}
	ty := &TypeApp{Tycon: TcArrow, Args: args}
	LogTypingf("translate format => %s", ReprOfType(ty))
	return ty, nil
}
