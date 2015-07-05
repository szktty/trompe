package trompe

type Type interface {
	// this method is to avoid using interface{} to prevent type error
	TypeTag() int
}

const (
	_TypeTag = iota
	TypeTagMeta
	TypeTagVar
	TypeTagApp
	TypeTagPoly
)

type TypeVar struct {
	Name string
}

type TypePoly struct {
	Tyvars []string
	Type   Type
}

type TypeApp struct {
	Tycon Tycon
	Args  []Type
}

// type constructor
type Tycon interface {
	TyconTag() int
}

const (
	_TyconTag = iota
	TyconTagUnit
	TyconTagBool
	TyconTagInt
	TyconTagFloat
	TyconTagChar
	TyconTagString
	TyconTagExn
	TyconTagFormat
	TyconTagFormatter
	TyconTagRegexp
	TyconTagList
	TyconTagArray
	TyconTagTuple
	TyconTagModule
	TyconTagArrow
	TyconTagKeyArrow
	TyconTagVariant
	TyconTagVariantTag
	TyconTagUnique
	TyconTagTyFun
	TyconTagSig // TODO:derecated
)

type TypeMeta struct {
	Name string
	Type Type
}

type TyconUnit struct{}
type TyconBool struct{}
type TyconInt struct{}
type TyconFloat struct{}
type TyconChar struct{}
type TyconString struct{}
type TyconFormat struct{}
type TyconFormatter struct{}
type TyconRegexp struct{}
type TyconList struct{}
type TyconArray struct{}
type TyconArrow struct{}
type TyconTuple struct{}

type TyconModule struct {
	Module *Module
}

type TyconTyFun struct {
	Tyvars []string
	Type   Type
}

type TyconExn struct {
	Name string
}

type TyconRecord struct {
	Module *Module
	Name   string
	Fields []string
}

type TyconVariant struct {
	Module  *Module
	Name    string
	Constrs []string
}

type TyconVariantTag struct {
	Tag     string
	Variant Tycon
}

type TyconUnique struct {
	Tycon Tycon
	ID    int
}

type TyconKeyArrow struct {
	Keywords []string
}

var TyvarNames = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

func (ty *TypeVar) TypeTag() int {
	return TypeTagVar
}

func (ty *TypeApp) TypeTag() int {
	return TypeTagApp
}

func (ty *TypePoly) TypeTag() int {
	return TypeTagPoly
}

func (ty *TypeMeta) TypeTag() int {
	return TypeTagMeta
}

func (tycon *TyconUnit) TyconTag() int {
	return TyconTagUnit
}

func (tycon *TyconBool) TyconTag() int {
	return TyconTagBool
}

func (tycon *TyconString) TyconTag() int {
	return TyconTagString
}

func (tycon *TyconChar) TyconTag() int {
	return TyconTagChar
}

func (tycon *TyconInt) TyconTag() int {
	return TyconTagInt
}

func (tycon *TyconFloat) TyconTag() int {
	return TyconTagFloat
}

func (tycon *TyconExn) TyconTag() int {
	return TyconTagExn
}

func (tycon *TyconFormat) TyconTag() int {
	return TyconTagFormat
}

func (tycon *TyconFormatter) TyconTag() int {
	return TyconTagFormatter
}

func (tycon *TyconRegexp) TyconTag() int {
	return TyconTagRegexp
}

func (tycon *TyconList) TyconTag() int {
	return TyconTagList
}

func (tycon *TyconArray) TyconTag() int {
	return TyconTagArray
}

func (tycon *TyconTuple) TyconTag() int {
	return TyconTagTuple
}

func (tycon *TyconModule) TyconTag() int {
	return TyconTagModule
}

func (tycon *TyconArrow) TyconTag() int {
	return TyconTagArrow
}

func (tycon *TyconKeyArrow) TyconTag() int {
	return TyconTagKeyArrow
}

func (tycon *TyconTyFun) TyconTag() int {
	return TyconTagTyFun
}

func (tycon *TyconVariant) TyconTag() int {
	return TyconTagVariant
}

func (tycon *TyconVariantTag) TyconTag() int {
	return TyconTagVariantTag
}

func (tycon *TyconUnique) TyconTag() int {
	return TyconTagUnique
}

func ListOfTypes(v ...interface{}) []Type {
	ts := make([]Type, len(v))
	for i, t := range v {
		ts[i] = t.(Type)
	}
	return ts
}

func TyconTagOfTypeApp(ty Type) (int, bool) {
	if deref, ok := DerefType(ty); ok {
		if app, ok := deref.(*TypeApp); ok {
			return app.Tycon.TyconTag(), true
		}
	}
	return -1, false
}

// NOTE: Do not use this function before typing
func EqualTypes(ty1 Type, ty2 Type) bool {
	if ty1 == ty2 {
		return true
	}

	tag1 := ty1.TypeTag()
	tag2 := ty2.TypeTag()
	if tag1 != tag2 {
		return false
	} else if tag1 == TypeTagMeta && tag2 == TypeTagMeta {
		return EqualTypes(ty1.(*TypeMeta).Type, ty2.(*TypeMeta).Type)
	} else if tag1 == TypeTagMeta {
		panic("cannot compare the types before typing")
	} else if tag2 == TypeTagMeta {
		return EqualTypes(ty2.(*TypeMeta).Type, ty1)
	}

	// TODO
	return false
}

func DerefType(ty Type) (Type, bool) {
	if meta, ok := ty.(*TypeMeta); ok {
		if meta.Type != nil {
			return DerefType(meta.Type)
		} else {
			return nil, false
		}
	} else {
		return ty, true
	}
}

func FreeMetavarsOfType(ty Type) []string {
	return scanFreeMetavars(ty, make([]string, 0))
}

func scanFreeMetavars(ty Type, accu []string) []string {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return scanFreeMetavars(desc.Type, accu)
		} else {
			return AppendStringIfAbsent(accu, desc.Name)
		}
	case *TypeVar:
		return AppendStringIfAbsent(accu, desc.Name)
	case *TypePoly:
		return scanFreeMetavars(desc.Type, accu)
	case *TypeApp:
		for _, arg := range desc.Args {
			accu = scanFreeMetavars(arg, accu)
		}
		return accu
	default:
		return accu
	}
}

func ContainsMeta(ty1 *TypeMeta, ty2 Type) bool {
	switch desc := ty2.(type) {
	case *TypeMeta:
		if ty1.Name == desc.Name {
			return true
		} else if desc.Type != nil {
			return ContainsMeta(ty1, desc.Type)
		} else {
			return false
		}
	case *TypePoly:
		return ContainsMeta(ty1, desc.Type)
	case *TypeApp:
		for _, arg := range desc.Args {
			if ContainsMeta(ty1, arg) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func FirstArgOfType(ty Type) (Type, bool) {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return FirstArgOfType(desc.Type)
		} else {
			return nil, false
		}
	case *TypePoly:
		if ty1, ok := FirstArgOfType(desc.Type); ok {
			return TPoly(desc.Tyvars, ty1), true
		} else {
			return nil, false
		}
	case *TypeApp:
		switch desc.Tycon.(type) {
		case *TyconTyFun:
			panic("error")
		default:
			if len(desc.Args) == 0 {
				return nil, false
			} else {
				return desc.Args[0], true
			}
		}
	default:
		return nil, false
	}
}

func PartialArrow(ty Type) (head Type, tail Type, ok bool) {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return PartialArrow(desc.Type)
		} else {
			return nil, nil, false
		}
	case *TypePoly:
		if head, tail, ok := PartialArrow(desc.Type); ok {
			return head, reassignTyvars(tail), true
		} else {
			return nil, nil, false
		}
	case *TypeApp:
		switch tycon := desc.Tycon.(type) {
		case *TyconTyFun:
			panic("error")
		case *TyconArrow:
			switch len(desc.Args) {
			case 0:
				panic("invalid no arguments arrow")
			case 1:
				panic("invalid 1 argument arrow")
			default:
				argLen := len(desc.Args)
				args := make([]Type, argLen-1)
				for i := 1; i < argLen; i++ {
					args[i-1] = desc.Args[i]
				}
				if argLen == 2 {
					Debugf("args = %s", desc.Args)
					return desc.Args[0], desc.Args[1], true
				} else {
					return desc.Args[0], TApp(tycon, args), true
				}
			}
		default:
			return nil, nil, false
		}
	default:
		return nil, nil, false
	}
}

func PartialKeyArrow(kw string, ty Type) (head Type, tail Type, ok bool) {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return PartialKeyArrow(kw, desc.Type)
		} else {
			return nil, nil, false
		}
	case *TypePoly:
		if head, tail, ok := PartialKeyArrow(kw, desc.Type); ok {
			return head, reassignTyvars(tail), true
		} else {
			return nil, nil, false
		}
	case *TypeApp:
		switch tycon := desc.Tycon.(type) {
		case *TyconTyFun:
			panic("error")
		case *TyconKeyArrow:
			switch len(desc.Args) {
			case 0:
				panic("invalid no arguments keyarrow")
			case 1:
				panic("invalid 1 argument keyarrow")
			default:
				i := -1
				for j, kw1 := range tycon.Keywords {
					if kw == kw1 {
						i = j
						break
					}
				}
				if i < 0 {
					return nil, nil, false
				}

				argLen := len(desc.Args)
				tailKws := make([]string, 0)
				tailArgs := make([]Type, 0)
				for j := 0; j < len(tycon.Keywords); j++ {
					if i != j {
						tailKws = append(tailKws, tycon.Keywords[j])
						tailArgs = append(tailArgs, desc.Args[j])
					}
				}
				tailArgs = append(tailArgs, desc.Args[argLen-1])

				if argLen == 2 {
					return desc.Args[0], desc.Args[1], true
				} else {
					return desc.Args[i], TApp(TcKeyArrow(tailKws), tailArgs), true
				}
			}
		default:
			return nil, nil, false
		}
	default:
		return nil, nil, false
	}
}

func reassignTyvars(ty Type) Type {
	tyvars := TyvarsOfType(ty)
	tyvars1 := make([]string, len(tyvars))
	env := make(map[string]string)
	for i, tyvar := range tyvars {
		tyvar1 := TyvarNames[i]
		env[tyvar] = tyvar1
		tyvars1[i] = tyvar1
	}
	return TPoly(tyvars1, reassignTyvars1(env, ty))
}

func reassignTyvars1(env map[string]string, ty Type) Type {
	switch desc := ty.(type) {
	case *TypeVar:
		return TVar(env[desc.Name])
	case *TypePoly:
		tyvars := make([]string, len(desc.Tyvars))
		for i, tyvar := range desc.Tyvars {
			tyvars[i] = env[tyvar]
		}
		return TPoly(tyvars, reassignTyvars1(env, desc.Type))
	case *TypeApp:
		switch tycon := desc.Tycon.(type) {
		case *TyconTyFun:
			tyvars := make([]string, len(tycon.Tyvars))
			for i, tyvar := range tycon.Tyvars {
				tyvars[i] = env[tyvar]
			}
			ty1 := reassignTyvars1(env, tycon.Type)
			args := make([]Type, len(desc.Args))
			for i, arg := range desc.Args {
				args[i] = reassignTyvars1(env, arg)
			}
			return TApp(TcTyFun(tyvars, ty1), args)
		default:
			args := make([]Type, len(desc.Args))
			for i, arg := range desc.Args {
				args[i] = reassignTyvars1(env, arg)
			}
			return TApp(tycon, args)
		}
	default:
		return ty
	}
}

func EqualsTyconTags(x Tycon, y Tycon) bool {
	return x.TyconTag() == y.TyconTag()
}

func RootTyconOfType(ty Type) (Tycon, bool) {
	if tycon, ok := TyconOfType(ty); ok {
		switch desc := tycon.(type) {
		case *TyconTyFun:
			return RootTyconOfType(desc.Type)
		default:
			return desc, true
		}
	} else {
		return nil, false
	}
}

func TyconOfType(ty Type) (Tycon, bool) {
	if poly, ok := ty.(*TypePoly); ok {
		return TyconOfType(poly.Type)
	} else if app, ok := ty.(*TypeApp); ok {
		return app.Tycon, true
	} else {
		return nil, false
	}
}

func TypeArrowOfType(ty Type) (*TypeApp, bool) {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type != nil {
			return TypeArrowOfType(desc.Type)
		} else {
			return nil, false
		}
	case *TypePoly:
		Panicf("type must be instantiated: %s", ReprOfType(ty))
		return nil, false
	case *TypeApp:
		switch desc.Tycon.(type) {
		case *TyconArrow, *TyconKeyArrow:
			return desc, true
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}

func (app *TypeApp) HeadArg() (Type, bool) {
	if len(app.Args) == 0 {
		return nil, false
	} else {
		return app.Args[0], true
	}
}

func (app *TypeApp) Tail() (*TypeApp, bool) {
	if len(app.Args) <= 1 {
		return nil, false
	} else {
		return &TypeApp{Tycon: app.Tycon, Args: app.Args[1:]}, true
	}
}

func TypeAppOfType(ty Type) (*TypeApp, bool) {
	switch desc := ty.(type) {
	case *TypeMeta:
		if desc.Type == nil {
			return nil, false
		} else {
			return TypeAppOfType(desc.Type)
		}
	case *TypePoly:
		return TypeAppOfType(desc.Type)
	case *TypeApp:
		return desc, true
	default:
		return nil, false
	}
}

func PickTypeApp(ty Type) (*TypeApp, bool) {
	ty, ok := DerefType(ty)
	if !ok {
		return nil, false
	}
	if poly, ok := ty.(*TypePoly); ok {
		app, ok := poly.Type.(*TypeApp)
		return app, ok
	} else {
		app, ok := ty.(*TypeApp)
		return app, ok
	}
}
