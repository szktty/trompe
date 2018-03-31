package trompe

import "fmt"

func LibCoreShow(ctx *Context, args []Value, nargs int) (Value, error) {
	if err := ValidateArity(nil, 1, nargs); err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", args[0].Desc())
	return LangUnit, nil
}

func InstallLibCore() {
	SetPrim("show", LibCoreShow, 1)
	m := NewModule(nil, "core")
	m.AddPrim("show", "show")
	AddTopModule(m)
	AddOpenedModule(m)

}
