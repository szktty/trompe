package trompe

import (
	"fmt"
	"strconv"
)

type codeComp struct {
	comp     *compiler
	params   []string
	syms     []string
	lits     []Value
	ops      []int
	labels   int
	labelMap map[int]int
	funComps map[string]*codeComp
}

type compiler struct {
	path string
}

func newCodeComp(comp *compiler) *codeComp {
	return &codeComp{
		comp:     comp,
		syms:     make([]string, 0),
		lits:     make([]Value, 0),
		ops:      make([]int, 0),
		labels:   -1,
		labelMap: make(map[int]int, 16),
		funComps: make(map[string]*codeComp, 16),
	}
}

func (c *codeComp) newCodeComp() *codeComp {
	new := newCodeComp(c.comp)
	return new
}

func (c *codeComp) addParam(name string) {
	c.params = append(c.params, name)
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

func (c *codeComp) addOpPop() {
	c.addOp(OpPop)
}

func (c *codeComp) addOpJump(label int) {
	c.addOp(OpJump)
	c.addOp(label)
}

func (c *codeComp) addOpBranch(flag bool, label int) {
	if flag {
		c.addOp(OpBranchTrue)
	} else {
		c.addOp(OpBranchFalse)
	}
	c.addOp(label)
}

func (c *codeComp) addSym(name string) int {
	for i, name1 := range c.syms {
		if name1 == name {
			return i
		}
	}
	c.syms = append(c.syms, name)
	return len(c.syms) - 1
}

func (c *codeComp) addLit(val Value) int {
	c.lits = append(c.lits, val)
	return len(c.lits) - 1
}

func (c *codeComp) addStr(s string) int {
	for i, lit := range c.lits {
		if v, ok := ValueToString(lit); ok {
			if v.Value == s {
				return i
			}
		}
	}
	return c.addLit(NewString(s))
}

func (c *codeComp) addFun(name string, comp *codeComp) {
	c.funComps[name] = comp
}

func (c *codeComp) addMatch(n PtnNode) {
	ptn := NewPatternFromNode(n)
	i := c.addLit(ptn)
	c.addOp(OpLoadLit)
	c.addOp(i)
	c.addOp(OpMatch)
}

func (c *codeComp) addOpPanic(kind int) {
	c.addOp(OpPanic)
	c.addOp(kind)
}

func (c *codeComp) code() *CompiledCode {
	code := NewCompiledCode()
	code.Syms = c.syms
	code.Lits = c.lits
	code.Ops = c.ops
	code.Labels = c.labelMap
	return code
}

func (c *codeComp) compile(node Node) {
	switch node := node.(type) {
	case *ChunkNode:
		c.compile(node.Block)
		c.addOpPop()
	case *BlockNode:
		c.addOp(OpBegin)
		l := len(node.Stats)
		for i, stat := range node.Stats {
			c.compile(stat)
			if i+1 < l {
				c.addOpPop()
			}
		}
		c.addOp(OpEnd)
	case *LetStatNode:
		c.compile(node.Ptn)
		c.compile(node.Exp)
		c.addOp(OpMatch)
		c.addOp(OpPop)
	case *DefStatNode:
		defComp := c.newCodeComp()
		defComp.params = node.Params.NameStrs()
		defComp.compile(&node.Block)
		c.addFun(node.Name.Text, defComp)
	case *ShortDefStatNode:
		defComp := c.newCodeComp()
		defComp.params = node.Params.NameStrs()
		defComp.compile(node.Exp)
		defComp.addOp(OpReturn)
		c.addFun(node.Name.Text, defComp)
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
			c.compile(node.ElseAction)
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
			c.compile(node.ElseAction)
		}
		c.addLabel(endL)
	case *ForStatNode:
		beginL := c.newLabel()
		panicL := c.newLabel()
		endL := c.newLabel()
		c.addOp(OpBegin)
		c.compile(node.Exp)
		c.addOp(OpIter)

		// loop ahead
		c.addLabel(beginL)
		c.addOp(OpBranchNext)
		c.addOp(endL)
		c.addMatch(node.Ptn)
		c.addOpBranch(false, panicL)
		c.compile(&node.Block)
		c.addOpPop()
		c.addOpJump(beginL)

		// pattern matching error
		c.addLabel(panicL)
		c.addOpPanic(OpPanicMatch)

		c.addLabel(endL)
		c.addOp(OpEnd)
		c.addOpPop()
	case *RetStatNode:
		if node.Exp == nil {
			c.addOp(OpReturnUnit)
		} else {
			c.compile(node.Exp)
			c.addOp(OpReturn)
		}
	case *FunCallExpNode:
		c.compile(node.Callable)
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
		i := c.addSym(node.Name.Text)
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
		val, err := strconv.Atoi(node.Value.Text)
		if err != nil {
			panic(fmt.Sprintf("atoi failed: %s", err.Error()))
		}
		c.addOp(OpLoadInt)
		c.addOp(val)
	case *StrExpNode:
		i := c.addStr(node.Value.Text)
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
		anonComp := newCodeComp(c.comp)
		for _, name := range node.Params.Names {
			anonComp.addParam(name.Text)
		}
		for _, stat := range node.Stats {
			anonComp.compile(stat)
		}
		anonComp.addOpPop()
		anonComp.compile(node.Exp)
		anonComp.addOp(OpReturn)
		code := anonComp.code()
		c.addLit(code)
	case *RangeExpNode:
		c.compile(node.Left)
		c.compile(node.Right)
		if node.Close {
			c.addOp(OpClosedRange)
		} else {
			c.addOp(OpHalfOpenRange)
		}
	default:
		panic(fmt.Sprintf("unsupported node %s", NodeDesc(node)))
	}
}

func Compile(path string, node Node) *CompiledCode {
	comp := &compiler{path: path}
	codeComp := newCodeComp(comp)
	codeComp.compile(node)
	return codeComp.code()
}
