package trompe

import (
	"bytes"
	"fmt"
)

func StringOfType(ty Type) string {
	buf := bytes.NewBufferString("")
	occur := make([]interface{}, 0)
	writeType(nil, ty, buf, occur)
	return buf.String()
}

// TODO: poly 対応
func writeType(tyvars []string, ty Type, buf *bytes.Buffer, occur []interface{}) []interface{} {
	for _, exist := range occur {
		if ty == exist {
			return occur
		}
	}

	switch desc := ty.(type) {
	case nil:
		panic("type must not be nil")

	case *TypeMeta:
		occur = append(occur, ty)
		if desc.Type != nil {
			return writeType(tyvars, desc.Type, buf, occur)
		} else {
			buf.WriteString("'" + desc.Name)
		}

	case *TypeVar:
		buf.WriteString("'" + desc.Name)
	case *TypePoly:
		occur = writeType(desc.Tyvars, desc.Type, buf, append(occur, ty))
	case *TypeApp:
		occur = writeTycon(tyvars, desc, desc.Tycon, buf, occur)
	default:
		Debugf("type = %s", desc)
		panic("unknown type")
	}
	return occur
}

func writeTycon(tyvars []string, ty *TypeApp, tycon Tycon, buf *bytes.Buffer,
	occur []interface{}) []interface{} {
	for _, exist := range occur {
		if tycon == exist {
			buf.WriteString("...")
			return occur
		}
	}

	switch tycon := tycon.(type) {
	case *TyconUnit:
		buf.WriteString("unit")
	case *TyconBool:
		buf.WriteString("bool")
	case *TyconInt:
		buf.WriteString("int")
	case *TyconFloat:
		buf.WriteString("float")
	case *TyconChar:
		buf.WriteString("char")
	case *TyconString:
		buf.WriteString("string")
	case *TyconExn:
		buf.WriteString("exn")
	case *TyconFormat:
		buf.WriteString("(")
		occur = writeVarsOfType(tyvars, ty.Args, ", ", buf, occur)
		buf.WriteString(") ")
		buf.WriteString("format")
	case *TyconFormatter:
		buf.WriteString("(")
		occur = writeVarsOfType(tyvars, ty.Args, ", ", buf, occur)
		buf.WriteString(") ")
		buf.WriteString("formatter")
	case *TyconRegexp:
		buf.WriteString("regexp")
	case *TyconList:
		occur = append(occur, tycon)
		if len(ty.Args) == 0 {
			buf.WriteString("'a")
		} else {
			occur = writeType(tyvars, ty.Args[0], buf, occur)
		}
		buf.WriteString(" list")
	case *TyconArray:
		occur = append(occur, tycon)
		if len(ty.Args) == 0 {
			buf.WriteString("'a")
		} else {
			occur = writeType(tyvars, ty.Args[0], buf, occur)
		}
		buf.WriteString(" array")
	case *TyconTuple:
		occur = append(occur, tycon)
		buf.WriteString("(")
		occur = writeVarsOfType(tyvars, ty.Args, " * ", buf, occur)
		buf.WriteString(")")
	case *TyconArrow:
		occur = append(occur, tycon)
		outer := buf.Len() == 0
		if !outer {
			buf.WriteString("(")
		}
		occur = writeVarsOfType(tyvars, ty.Args, " -> ", buf, occur)
		if !outer {
			buf.WriteString(")")
		}
	case *TyconLabeledArrow:
		occur = append(occur, tycon)
		outer := buf.Len() == 0
		if !outer {
			buf.WriteString("(")
		}
		for i := 0; i < len(ty.Args); i++ {
			if i < len(tycon.Names) {
				label := tycon.Names[i]
				if label != TyconUnlabeledName {
					buf.WriteString(label + ":")
				}
				occur = writeType(tyvars, ty.Args[i], buf, occur)
				buf.WriteString(" -> ")
			} else {
				occur = writeType(tyvars, ty.Args[i], buf, occur)
			}
		}
		if !outer {
			buf.WriteString(")")
		}
	case *TyconTyFun:
		occur = append(occur, tycon)
		occur = writeType(tyvars, tycon.Type, buf, occur)
	case *TyconModule:
		buf.WriteString("module " + tycon.Module.Name)
	case *TyconUnique:
		occur = append(occur, tycon)
		occur = writeTycon(tyvars, ty, tycon.Tycon, buf, occur)
	case *TyconVariantTag:
		writeTycon(tyvars, ty, tycon.Variant, buf, occur)
	case *TyconVariant:
		switch len(tyvars) {
		case 0:
			break
		case 1:
			buf.WriteString("'" + tyvars[0] + " ")
		default:
			for _, tyvar := range tyvars {
				buf.WriteString("'" + tyvar + " ")
			}
		}
		if tycon.Module != nil {
			buf.WriteString(tycon.Module.Name + ".")
		}
		buf.WriteString(tycon.Name)
	default:
		panic(fmt.Errorf("unknown tycon %s", ty.Tycon))
	}
	return occur
}

func writeVarsOfType(tyvars []string, args []Type, sep string, buf *bytes.Buffer, occur []interface{}) []interface{} {
	for i, arg := range args {
		occur = writeType(tyvars, arg, buf, occur)
		if i+1 < len(args) {
			buf.WriteString(sep)
		}
	}
	return occur
}

func ReprOfType(ty Type) string {
	buf := bytes.NewBufferString("")
	occur := make([]interface{}, 0)
	writeReprOfType(ty, buf, occur)
	return buf.String()
}

func writeReprOfType(ty Type, buf *bytes.Buffer, occur []interface{}) []interface{} {
	for _, exist := range occur {
		if ty == exist {
			buf.WriteString("...")
			return occur
		}
	}

	switch desc := ty.(type) {
	case nil:
		panic("type must not be nil")

	case *TypeMeta:
		occur = append(occur, ty)
		buf.WriteString("Meta(" + desc.Name)
		if desc.Type != nil {
			buf.WriteString(", ")
			occur = writeReprOfType(desc.Type, buf, occur)
		}
		buf.WriteString(")")

	case *TypeVar:
		buf.WriteString("Var(" + desc.Name + ")")

	case *TypePoly:
		occur = append(occur, ty)
		buf.WriteString("Poly([")
		for i, name := range desc.Tyvars {
			buf.WriteString(name)
			if i+1 < len(desc.Tyvars) {
				buf.WriteString(", ")
			}
		}
		buf.WriteString("], ")
		writeReprOfType(desc.Type, buf, occur)
		buf.WriteString(")")

	case *TypeApp:
		switch desc.Tycon.(type) {
		case *TyconUnit, *TyconBool, *TyconInt, *TyconFloat, *TyconString,
			*TyconChar:
			break
		default:
			occur = append(occur, ty)
		}
		buf.WriteString("App(")
		occur = writeReprOfTycon(desc.Tycon, buf, occur)
		buf.WriteString(", [")
		occur = writeReprVarsOfType(desc.Args, ", ", buf, occur)
		buf.WriteString("])")
	default:
		Panicf("unknown type %s", desc)
	}
	return occur
}

func writeReprOfTycon(desc Tycon, buf *bytes.Buffer, occur []interface{}) []interface{} {
	for _, exist := range occur {
		if desc == exist {
			buf.WriteString("...")
			return occur
		}
	}

	switch tycon := desc.(type) {
	case *TyconUnit:
		buf.WriteString("Unit")
	case *TyconBool:
		buf.WriteString("Bool")
	case *TyconInt:
		buf.WriteString("Int")
	case *TyconFloat:
		buf.WriteString("Float")
	case *TyconChar:
		buf.WriteString("Char")
	case *TyconString:
		buf.WriteString("String")
	case *TyconExn:
		buf.WriteString("Exn")
	case *TyconFormat:
		buf.WriteString("Format")
	case *TyconFormatter:
		buf.WriteString("Formatter")
	case *TyconRegexp:
		buf.WriteString("Regexp")
	case *TyconList:
		buf.WriteString("List")
	case *TyconArray:
		buf.WriteString("Array")
	case *TyconTuple:
		buf.WriteString("Tuple")
	case *TyconArrow:
		buf.WriteString("Arrow")
	case *TyconLabeledArrow:
		buf.WriteString("LabeledArrow([")
		for i, label := range tycon.Names {
			buf.WriteString(label)
			if i+1 < len(tycon.Names) {
				buf.WriteString(", ")
			}
		}
		buf.WriteString("])")
	case *TyconTyFun:
		buf.WriteString("TyFun([")
		for i, name := range tycon.Tyvars {
			buf.WriteString(name)
			if i+1 < len(tycon.Tyvars) {
				buf.WriteString(", ")
			}
		}
		buf.WriteString("], ")
		writeReprOfType(tycon.Type, buf, occur)
		buf.WriteString(")")
	case *TyconVariantTag:
		buf.WriteString(fmt.Sprintf("VariantTag(%d, %s)",
			tycon.Variant.(*TyconUnique).ID, tycon.Tag))
	case *TyconModule:
		buf.WriteString("Module(\"" + tycon.Module.Name + "\")")
	case *TyconUnique:
		buf.WriteString(fmt.Sprintf("Unique(%d, ", tycon.ID))
		occur = writeReprOfTycon(tycon.Tycon, buf, occur)
		buf.WriteString(")")
	default:
		panic(fmt.Errorf("unknown tycon %s", desc))
	}
	return occur
}

func writeReprVarsOfType(args []Type, sep string, buf *bytes.Buffer,
	occur []interface{}) []interface{} {
	for i, arg := range args {
		occur = writeReprOfType(arg, buf, occur)
		if i+1 < len(args) {
			buf.WriteString(sep)
		}
	}
	return occur
}
