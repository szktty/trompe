package trompe

import (
	"fmt"
	"strings"
)

type State struct {
	Root       *Module
	SearchPath []string
	Stack      *Stack
}

func Init() {
	// check bytecode size
	if OpLastCode > 255 {
		Panicf("number of opcode must be less than 255 (%d)", OpLastCode)
	}
	Debugf("number of opcode: %d", OpLastCode)
}

func NewState() *State {
	root := ModuleTrompe()
	return &State{Root: root, SearchPath: SearchPath(),
		Stack: NewStack(NormalFrameSize)}
}

func (state *State) NewTopEnv() (*Env, bool) {
	env := NewEnv(nil)
	env.AddFullImport(state.Root)
	return env, true
}

func (state *State) LoadModule(name string) (*Module, bool) {
	if mod, ok := state.Root.FindModule(name); ok {
		return mod, true
	} else {
		mod, err := state.LoadIntf(name)
		if err != nil {
			return nil, false
		}
		return mod, true
	}
}

func (state *State) FindModule(name string) (*Module, bool) {
	return state.Root.FindModule(name)
}

func (state *State) FindModuleOfPath(path *NamePath) (*Module, error) {
	if path == nil {
		return state.Root, nil
	}

	mod := state.Root
	accu := make([]string, 0)
	for _, name := range path.Base {
		if name == state.Root.Name && mod == state.Root {
			continue
		}
		accu = append(accu, name)
		sub, ok := mod.FindModule(name)
		if !ok {
			return nil, fmt.Errorf("module %s is not found",
				strings.Join(accu, "."))
		}
		mod = sub
	}
	return mod, nil
}

func (state *State) FindFieldValueOfPath(path *NamePath) (Value, error) {
	mod, err := state.FindModuleOfPath(path)
	if err != nil {
		return nil, err
	}
	if f, ok := mod.FindFieldValue(path.Name); ok {
		return f, nil
	} else {
		return nil, fmt.Errorf("field %s is not found", path.Name)
	}
}

func (state *State) Apply(parent *Context, f Value, args []Value) (Value, error) {
	switch desc := f.(type) {
	case *BlockClosure:
		ctx := state.NewContext(parent.Module, parent, desc, args)
		return state.Exec(ctx)
	case Primitive:
		return desc(state, parent, args)
	default:
		return nil, fmt.Errorf("value %s is not applicable", StringOfValue(f))
	}
}
