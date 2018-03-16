package trompe

import (
	"fmt"
	"strconv"
)

type codeComp struct {
	comp     *compiler
	args     []string
	lits     []Value
	ops      []int
	labels   int
	labelMap map[int]int
	funComps map[string]*codeComp
}

type compiler struct {
	path string
}

func createCodeComp(comp *compiler) *codeComp {
	return &codeComp{
		comp:     comp,
		lits:     make([]Value, 16),
		ops:      make([]int, 64),
		labels:   -1,
		labelMap: make(map[int]int, 16),
		funComps: make(map[string]*codeComp, 16),
	}
}

func (c *codeComp) newCodeComp() *codeComp {
	new := createCodeComp(c.comp)
	return new
}

func (c *codeComp) addArg(name string) {
	c.args = append(c.args, name)
}

func (c *codeComp) newLabel() int {
	c.labels += 1
	return c.labels
}

func (c *codeComp) addLabel(label int) {
	c.labelMap[label] = len(c.ops)
	c.addOp(OpLabel)
	c.addOp(label)
}

func (c *codeComp) addOp(op int) {
	c.ops = append(c.ops, op)
}

func (c *codeComp) addLit(val Value) int {
	c.lits = append(c.lits, val)
	return len(c.lits) - 1
}

func (c *codeComp) addStr(s string) int {
	for i, lit := range c.lits {
		if val, ok := lit.(*ValStr); ok {
			if val.Value == s {
				return i
			}
		}
	}
	return c.addLit(CreateValStr(s))
}

func (c *codeComp) addFun(name string, comp *codeComp) {
	c.funComps[name] = comp
}

func (c *codeComp) code() *CompiledCode {
	// TODO
	return nil
}

func (c *codeComp) compile(node Node) {
	switch node := node.(type) {
	case *ChunkNode:
		c.compile(&node.Block)
	case *BlockNode:
		for _, stat := range node.Stats {
			c.compile(stat)
		}
	case *LetStatNode:
		c.compile(node.Ptn)
		c.compile(node.Exp)
		c.addOp(OpMatch)
		c.addOp(OpPop)
	case *DefStatNode:
		defComp := c.newCodeComp()
		defComp.args = node.Args.NameStrs()
		defComp.compile(&node.Block)
		c.addFun(node.Name.Value, defComp)
	case *ShortDefStatNode:
		defComp := c.newCodeComp()
		defComp.args = node.Args.NameStrs()
		defComp.compile(node.Exp)
		defComp.addOp(OpReturn)
		c.addFun(node.Name.Value, defComp)
	case *IfStatNode:
		endL := c.newLabel()
		for _, cond := range node.Cond {
			nextL := c.newLabel()
			c.compile(cond.Cond)
			c.addOp(OpBranchFalse)
			c.addOp(nextL)
			c.compile(&cond.Action)
			c.addOp(OpJump)
			c.addOp(endL)
			c.addLabel(nextL)
		}
		if node.Else != nil {
			c.compile(&node.Else.Action)
		}
		c.addOp(OpLabel)
		c.addLabel(endL)
	case *CaseStatNode:
		endL := c.newLabel()
		c.compile(node.Cond)
		for _, clau := range node.Claus {
			nextL := c.newLabel()
			c.addOp(OpDup)
			c.compile(clau.Ptn)
			c.addOp(OpMatch)
			c.addOp(OpBranchFalse)
			c.addOp(nextL)
			c.compile(clau.Action)
			c.addOp(OpJump)
			c.addOp(endL)
			c.addLabel(nextL)
		}
		c.addOp(OpPop) // Cond
		if node.Else != nil {
			c.compile(&node.Else.Action)
		}
		c.addLabel(endL)
	case *RetStatNode:
		if node.Value == nil {
			c.addOp(OpReturnUnit)
		} else {
			c.compile(node.Value)
			c.addOp(OpReturn)
		}
	case *FunCallStatNode:
		c.compile(node.Exp)
		c.addOp(OpPop)
	case *FunCallExpNode:
		c.compile(node.Prefix)
		for _, arg := range node.Args.Elts {
			c.compile(arg)
		}
		c.addOp(OpCall)
		c.addOp(len(node.Args.Elts))
	case *CondOpExpNode:
		falseL := c.newLabel()
		endL := c.newLabel()
		c.compile(node.Cond)
		c.addOp(OpBranchFalse)
		c.addOp(falseL)
		c.compile(node.True)
		c.addOp(OpJump)
		c.addOp(endL)
		c.addLabel(falseL)
		c.compile(node.False)
		c.addLabel(endL)
	case *VarExpNode:
		i := c.addStr(node.Name)
		c.addOp(OpLoadLocal)
		c.addOp(i)
	case *UnitExpNode:
		c.addOp(OpLoadUnit)
	case *BoolExpNode:
		if node.Value {
			c.addOp(OpLoadTrue)
		} else {
			c.addOp(OpLoadFalse)
		}
	case *IntExpNode:
		val, err := strconv.Atoi(node.Value)
		if err != nil {
			panic(fmt.Sprintf("atoi failed: %s", err.Error()))
		}
		c.addOp(OpLoadInt)
		c.addOp(val)
	case *StrExpNode:
		i := c.addStr(node.Value)
		c.addOp(OpLoadLit)
		c.addOp(i)
	case *ListExpNode:
		ln := len(node.Elts.Elts)
		for _, elt := range node.Elts.Elts {
			c.compile(elt)
		}
		c.addOp(OpList)
		c.addOp(ln)
	case *TupleExpNode:
		ln := len(node.Elts.Elts)
		for _, elt := range node.Elts.Elts {
			c.compile(elt)
		}
		c.addOp(OpTuple)
		c.addOp(ln)
	case *SomeExpNode:
		c.compile(node.Value)
		c.addOp(OpSome)
	case *NoneExpNode:
		c.addOp(OpLoadNone)
	case *AnonFunExpNode:
		anonComp := createCodeComp(c.comp)
		for _, arg := range node.Args {
			if nameExp, ok := arg.(*StrExpNode); ok {
				anonComp.addArg(nameExp.Value)
			} else {
				panic("not StrExp")
			}
		}
		anonComp.compile(&node.Block)
		code := anonComp.code()
		val := CreateValClos(code)
		c.addLit(val)
	default:
		panic(fmt.Sprintf("unsupported node %p", node))
	}
}

type CompiledProg struct {
	Path string
	Code *CompiledCode
}

func Compile(path string, node Node) *CompiledProg {
	comp := &compiler{path: path}
	codeComp := &codeComp{comp: comp}
	codeComp.compile(node)
	code := codeComp.code()
	return &CompiledProg{Path: path, Code: code}
}
