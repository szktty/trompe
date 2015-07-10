package trompe

import (
	"fmt"
	"github.com/jteeuwen/go-pkg-optarg"
	"os"
	"path/filepath"
	"strings"
)

type LoadingContext struct {
	File      string
	Name      string
	NoWarn    bool
	WarnError bool
}

const (
	LoadingLogType = iota
	LoadingLogWarn
	LoadingLogError
)

const (
	LoadingPhase = iota
	LoadingPhaseSyntax
	LoadingPhaseTyping
	LoadingPhaseCompilation
	LoadingPhaseAll
)

func ModuleNameOfPath(path string) string {
	_, base := filepath.Split(path)
	name := strings.ToLower(strings.Split(base, ".")[0])
	return fmt.Sprintf("%c%s", strings.ToUpper(name)[0], name[1:])
}

func (state *State) LoadFile(fpath string, phase int, opts []*optarg.Option) (*Module, error) {
	m, code, err := state.CompileFile(fpath, phase, opts)
	if err != nil {
		return nil, err
	}
	if phase == LoadingPhaseCompilation {
		return m, nil
	}

	if code != nil {
		state.Root.AddModule(m)
		block := NewBlockClosure(code)
		ctx := state.NewContext(m, nil, block, nil)
		res, err := state.Exec(ctx)
		if err != nil {
			return nil, err
		}
		Debugf("res = %s", res)
	}
	return m, nil
}

func (state *State) CompileFile(fpath string, phase int, opts []*optarg.Option) (*Module, *CompiledCode, error) {
	// TODO: ファイルパスがある場合のモジュールの扱い
	name := ModuleNameOfPath(fpath)
	if m, ok := state.FindModule(name); ok {
		return m, nil, nil
	}

	ctx := &LoadingContext{File: fpath, Name: name}
	printsIntf := false
	for _, opt := range opts {
		switch opt.Name {
		case "nowarn":
			ctx.NoWarn = true
		case "warn-error":
			ctx.WarnError = true
		case "interface":
			printsIntf = true
		}
	}

	scn := NewLexerFromFile(fpath)
	node, err := Parse(scn)
	if err != nil {
		return nil, nil, err
	}
	if phase == LoadingPhaseSyntax {
		return nil, nil, nil
	}

	env, _ := state.NewTopEnv()
	tnode, err := state.Typing(env, node)
	if err != nil {
		return nil, nil, err
	}
	if phase == LoadingPhaseTyping {
		return nil, nil, nil
	}

	mod := NewModule(name)
	mod.MergeEnv(env)
	if printsIntf {
		mod.PrintIntf()
		return mod, nil, nil
	}

	code, err := state.Compile(ctx, tnode)
	if err != nil {
		return nil, nil, err
	}
	return mod, code, nil
}

func (ctx *LoadingContext) Warn(err error) {
	if ctx.NoWarn {
		return
	}
	if ctx.WarnError {
		ctx.Error(err)
	}
	fmt.Printf("Warning: %s\n", err.Error())
}

func (ctx *LoadingContext) Warnf(f string, v ...interface{}) {
	ctx.Warn(fmt.Errorf(f, v))
}

func (ctx *LoadingContext) Error(err error) {
	fmt.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

func (ctx *LoadingContext) Errorf(f string, v ...interface{}) {
	ctx.Error(fmt.Errorf(f, v))
}
