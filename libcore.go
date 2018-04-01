package trompe

import "fmt"

func LibCoreId(ctx *Context, args []Value, nargs int) (Value, error) {
	return args[0], nil
}

func LibCoreShow(ctx *Context, args []Value, nargs int) (Value, error) {
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
