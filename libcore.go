package trompe

import "fmt"

func LibCoreId(ctx *Context, args []Value, nargs int) (Value, error) {
	if err := ValidateArity(nil, 1, nargs); err != nil {
		return nil, err
	}
	return args[0], nil
}

func LibCoreShow(ctx *Context, args []Value, nargs int) (Value, error) {
	if err := ValidateArity(nil, 1, nargs); err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", args[0].Desc())
	return LangUnit, nil
}

func InstallLibCore() {
	m := NewModule(nil, "core")
	m.AddPrim("id", LibCoreId, 1)
	m.AddPrim("show", LibCoreShow, 1)
	AddTopModule(m)
	AddOpenedModule(m)

}
