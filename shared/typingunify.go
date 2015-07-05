package trompe

func (inf *inferer) unify(loc *Loc, ty1 Type, ty2 Type) error {
	LogTypingf("try unify at %s", loc.StartString())
	if inf.basicUnify(ty1, ty2) {
		return nil
	} else {
		LogTypingf("type mismatch: %s and %s", ReprOfType(ty1), ReprOfType(ty2))
		gty1, err := inf.generalize(NewTyvarEnv(), ty1)
		if err != nil {
			Panicf("generalize failed: %s", ReprOfType(ty1))
		}
		gty2, err := inf.generalize(NewTyvarEnv(), ty2)
		if err != nil {
			Panicf("generalize failed: %s", ReprOfType(ty2))
		}
		return RuntimeErrorf(loc, "Type mismatch:\n    This expression has type: %s\n    but the expression was expected of type: %s",
			StringOfType(gty2), StringOfType(gty1))
	}
}

var retMsg = map[bool]string{true: "success", false: "failed"}

func (inf *inferer) basicUnify(ty1 Type, ty2 Type) (res bool) {
	ity1, err := inf.instantiate(ty1)
	if err != nil {
		panic("instantiate failed")
	}
	ity2, err := inf.instantiate(ty2)
	if err != nil {
		panic("instantiate failed")
	}
	ty1 = ity1
	ty2 = ity2

	id := inf.unifyCount
	inf.unifyCount++
	LogTypingf("==> BEGIN unify %d: %s and %s",
		id, ReprOfType(ty1), ReprOfType(ty2))
	defer LogTypingf("<== END unify %d", id)

	if ty1 == ty2 {
		return true
	} else if y, ok := ty2.(*TypeMeta); ok {
		if y.Type == nil {
			y.Type = ty1
			return true
		} else if x, ok := ty1.(*TypeMeta); ok {
			if x.Name == y.Name {
				return true
			} else if x.Type == nil {
				x.Type = y
			} else {
				return inf.basicUnify(x.Type, y.Type)
			}
		} else {
			return inf.basicUnify(ty1, y.Type)
		}
	}

	switch x := ty1.(type) {
	case *TypeMeta:
		if x.Type != nil {
			return inf.basicUnify(x.Type, ty2)
		} else {
			switch y := ty2.(type) {
			/*
				case *TypeMeta:
					if y.Type != nil {
						x.Type = y.Type
						return true
					} else if x.Name == y.Name {
						return true
					} else {
						x.Type = y
						return true
					}
			*/
			default:
				if ContainsMeta(x, ty2) {
					Panicf("same metavar must not be included in other type: %s and %s",
						ReprOfType(ty1), ReprOfType(ty2))
					return false
				}
				x.Type = y
				return true
			}
		}

	case *TypeVar:
		switch y := ty2.(type) {
		case *TypeVar:
			return x.Name == y.Name
		default:
			return false
		}

	case *TypePoly:
		switch y := ty2.(type) {
		case *TypePoly:
			return inf.basicUnify(x.Type, y.Type)
		default:
			return inf.basicUnify(x.Type, y)
		}

	case *TypeApp:
		switch y := ty2.(type) {
		case *TypePoly:
			return inf.basicUnify(x, y.Type)

		case *TypeApp:
			tag1 := x.Tycon.TyconTag()
			tag2 := y.Tycon.TyconTag()
			switch {
			case tag1 == TyconTagExn && tag2 == TyconTagExn:
				return true
			case len(x.Args) == 0 && len(y.Args) == 0:
				switch {
				case tag1 == TyconTagUnit && tag2 == TyconTagUnit,
					tag1 == TyconTagBool && tag2 == TyconTagBool,
					tag1 == TyconTagInt && tag2 == TyconTagInt,
					tag1 == TyconTagFloat && tag2 == TyconTagFloat,
					tag1 == TyconTagString && tag2 == TyconTagString,
					tag1 == TyconTagChar && tag2 == TyconTagChar,
					tag1 == TyconTagList && tag2 == TyconTagList,
					tag1 == TyconTagTuple && tag2 == TyconTagTuple,
					tag1 == TyconTagRegexp && tag2 == TyconTagRegexp:
					return true
				default:
					return false
				}
			case len(x.Args) == len(y.Args):
				switch {
				case tag1 == TyconTagVariantTag && tag2 == TyconTagVariantTag:
					tycon1 := x.Tycon.(*TyconVariantTag)
					tycon2 := y.Tycon.(*TyconVariantTag)
					return tycon1.Variant == tycon2.Variant
				case tag1 == TyconTagFormat && (tag2 == TyconTagFormat ||
					tag2 == TyconTagString):
					return true
				default:
					n := len(x.Args)
					if n != len(y.Args) {
						return false
					}
					for i := 0; i < n; i++ {
						if !inf.basicUnify(x.Args[i], y.Args[i]) {
							return false
						}
					}
					return true
				}
			default:
				return false
			}
		default:
			return false
		}

	default:
		LogTypingf("ty1: %s", ReprOfType(ty1))
		LogTypingf("ty2: %s", ReprOfType(ty2))
		panic("notimpl")
	}
}

func TooManyArgsError(loc *Loc, ty Type) error {
	return RuntimeErrorf(loc, "This function has type `%s'. It is applied to too many arguments; maybe you forgot a `;'.", StringOfType(ty))
}

func NotFunError(loc *Loc, ty Type) error {
	return RuntimeErrorf(loc, "This function has type `%s'. This is not a function; it cannot be applied.", StringOfType(ty))
}

func ConstrArityError(loc *Loc, name string, expected int, actual int) error {
	return RuntimeErrorf(loc, "The constructor `%s' expects %d argument(s), but is applied here to %d argument(s)",
		name, expected, actual)
}
