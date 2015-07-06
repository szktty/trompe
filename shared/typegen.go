package trompe

func NewTypeExn(name string, ty ...Type) Type {
	// TODO: poly
	return &TypeApp{Tycon: &TyconExn{Name: name}, Args: ty}
}

func NewTypeTuple(args ...Type) (Type, bool) {
	if len(args) == 0 {
		return nil, false
	} else {
		return &TypeApp{Tycon: &TyconTuple{}, Args: args}, true
	}
}

func Tyvars(names ...string) []string {
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = name
	}
	return ret
}

func TArgs(tys ...Type) []Type {
	ret := make([]Type, len(tys))
	for i, ty := range tys {
		ret[i] = ty
	}
	return ret
}

func TPoly(tyvars []string, ty Type) Type {
	return &TypePoly{Tyvars: tyvars, Type: ty}
}

func TVar(name string) Type {
	return &TypeVar{Name: name}
}

func TApp(tycon Tycon, args []Type) Type {
	return &TypeApp{Tycon: tycon, Args: args}
}

func TcTyFun(tyvars []string, ty Type) Tycon {
	return &TyconTyFun{Tyvars: tyvars, Type: ty}
}

func TcVariant(mod *Module, name string, constrs []string) Tycon {
	return &TyconVariant{Module: mod, Name: name, Constrs: constrs}
}

func TcVariantv(mod *Module, name string, constrs ...string) Tycon {
	return TcVariant(mod, name, constrs)
}

func NewTypeSome(ty Type) Type {
	return TPoly(Tyvars("a"),
		TApp(&TyconVariantTag{Tag: "Some", Variant: TcOption}, TArgs(ty)))
}

var uniqueID = 0

func TcUnique(tycon Tycon) Tycon {
	u := &TyconUnique{Tycon: tycon, ID: uniqueID}
	uniqueID++
	return u
}

func TcKeyArrow(kws []string) Tycon {
	return &TyconKeyArrow{Keywords: kws}
}

func TcKeyArrowv(kws ...string) Tycon {
	return TcKeyArrow(kws)
}

var TUnit = &TypeApp{Tycon: &TyconUnit{}}
var TBool = &TypeApp{Tycon: &TyconBool{}}
var TInt = &TypeApp{Tycon: &TyconInt{}}
var TFloat = &TypeApp{Tycon: &TyconFloat{}}
var TChar = &TypeApp{Tycon: &TyconChar{}}
var TString = &TypeApp{Tycon: &TyconString{}}
var TExn = &TypeApp{Tycon: &TyconExn{}}

var TList = TPoly(Tyvars("a"), TApp(&TyconList{}, TArgs(TVar("a"))))

var TFormat = TPoly(Tyvars("a"),
	TApp(
		TcTyFun(Tyvars("a", "b"),
			TApp(TcFormat, TArgs(TVar("a"), TVar("b")))),
		TArgs(TVar("a"), TVar("b"))))

var TFormatter = TPoly(Tyvars("a"),
	TApp(
		TcTyFun(Tyvars("a", "b"),
			TApp(&TyconFormatter{}, TArgs(TVar("a"), TVar("b")))),
		TArgs(TVar("a"), TVar("b"))))

var TNone = TPoly(Tyvars("a"),
	TApp(&TyconVariantTag{Tag: "None", Variant: TcOption}, nil))

var TcArrow = &TyconArrow{}
var TcFormat = &TyconFormat{}
var TcFormatter = &TyconFormatter{}

var TcList = TcTyFun(Tyvars("a"), TApp(&TyconList{}, TArgs(TVar("a"))))

var TcOption = TcTyFun(
	Tyvars("a"),
	TApp(
		TcVariantv(nil, "option", "Some", "None"),
		TArgs(TVar("a"), TUnit)))

var TcResult = TcTyFun(
	Tyvars("a", "b"),
	TApp(
		TcVariantv(nil, "result", "Ok", "Error"),
		TArgs(TVar("a"), TVar("b"))))

func WrapInTypePoly1(ty Type) Type {
	return TPoly(Tyvars("a"),
		TApp(
			TcTyFun(Tyvars("a"), ty),
			TArgs(TVar("a"))))
}

func WrapInTypePoly2(ty Type) Type {
	return TPoly(Tyvars("a", "b"),
		TApp(
			TcTyFun(Tyvars("a", "b"), ty),
			TArgs(TVar("a"), TVar("b"))))
}
