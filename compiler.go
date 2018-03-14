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
	}
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

func (c *codeComp) code() *CompiledCode {
	// TODO
	return nil
}

func (c *codeComp) compile(node Node) {
	switch node := node.(type) {
	case *Chunk:
		c.compile(node.Block)
	case *Block:
		for _, stat := range node.Stats {
			c.compile(stat)
		}
	case *LetStat:
		c.compile(node.Ptn)
		c.compile(node.Exp)
		c.addOp(OpMatch)
		c.addOp(OpPop)
	case *IfStat:
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
	case *CaseStat:
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
	case *RetStat:
		if node.Value == nil {
			c.addOp(OpReturnUnit)
		} else {
			c.compile(node.Value)
			c.addOp(OpReturn)
		}
	case *FunCallStat:
		c.compile(node.Exp)
		c.addOp(OpPop)
	case *FunCallExp:
		c.compile(node.Prefix)
		for _, arg := range node.Args.Elts {
			c.compile(arg)
		}
		c.addOp(OpCall)
		c.addOp(len(node.Args.Elts))
	case *VarExp:
		i := c.addStr(node.Name)
		c.addOp(OpLoadLocal)
		c.addOp(i)
	case *UnitExp:
		c.addOp(OpLoadUnit)
	case *BoolExp:
		if node.Value {
			c.addOp(OpLoadTrue)
		} else {
			c.addOp(OpLoadFalse)
		}
	case *IntExp:
		val, err := strconv.Atoi(node.Value)
		if err != nil {
			panic(fmt.Sprintf("atoi failed: %s", err.Error()))
		}
		c.addOp(OpLoadInt)
		c.addOp(val)
	case *StrExp:
		i := c.addStr(node.Value)
		c.addOp(OpLoadLit)
		c.addOp(i)
	case *ListExp:
		ln := len(node.Elts.Elts)
		for _, elt := range node.Elts.Elts {
			c.compile(elt)
		}
		c.addOp(OpList)
		c.addOp(ln)
	case *TupleExp:
		ln := len(node.Elts.Elts)
		for _, elt := range node.Elts.Elts {
			c.compile(elt)
		}
		c.addOp(OpTuple)
		c.addOp(ln)
	case *SomeExp:
		c.compile(node.Value)
		c.addOp(OpSome)
	case *NoneExp:
		c.addOp(OpLoadNone)
	case *AnonFunExp:
		anonComp := createCodeComp(c.comp)
		for _, arg := range node.Args {
			if nameExp, ok := arg.(*StrExp); ok {
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
