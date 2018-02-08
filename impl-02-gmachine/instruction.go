package climpl

import (
	"strconv"
)

type instruction int

const (
	_ instruction = iota
	INSTR_UNWIND
	INSTR_PUSHGLOBAL
	INSTR_PUSHINT
	INSTR_PUSH
	INSTR_MAKEAPPL
	INSTR_SLIDE
)

type instr struct {
	Op   instruction
	Int  int
	Name string
}

func (me instr) String() string {
	switch me.Op {
	case INSTR_UNWIND:
		return "Unwind"
	case INSTR_PUSHGLOBAL:
		return "PushG:" + me.Name
	case INSTR_PUSHINT:
		return "PushI:" + strconv.Itoa(me.Int)
	case INSTR_PUSH:
		return "Push:" + strconv.Itoa(me.Int)
	case INSTR_SLIDE:
		return "Slide:" + strconv.Itoa(me.Int)
	case INSTR_MAKEAPPL:
		return "MkAp"
	}
	return strconv.Itoa(int(me.Op))
}

type code []instr

func (me code) String() (s string) {
	s = "["
	for i, instr := range me {
		if i > 0 {
			s += " · "
		}
		s += instr.String()
	}
	return s + "]"
}

func (me *gMachine) dispatch(cur instr, nuCode code) code {
	stackpos := len(me.Stack) - 1
	switch cur.Op {
	case INSTR_PUSHGLOBAL:
		addr := me.lookup(cur.Name)
		me.Stack = append(me.Stack, addr)
	case INSTR_PUSHINT:
		addr := me.alloc(nodeNumUint(cur.Int))
		me.Stack = append(me.Stack, addr)
	case INSTR_MAKEAPPL:
		addrcallee := me.Stack[stackpos]
		addrarg := me.Stack[stackpos-1]
		addrappl := me.alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
		me.Stack[stackpos-1] = addrappl
		me.Stack = me.Stack[:len(me.Stack)-1]
	case INSTR_PUSH:
		applarg := me.Heap[me.Stack[stackpos-(1+cur.Int)]].(nodeAppl).Arg
		me.Stack = append(me.Stack, applarg)
	case INSTR_SLIDE:
		keep := me.Stack[stackpos]
		less := me.Stack[:len(me.Stack)-(1+cur.Int)]
		me.Stack = append(less, keep)
	case INSTR_UNWIND:
		addr := me.Stack[stackpos]
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodeNumUint:
			if len(nuCode) > 0 {
				panic("unexpected? or not?")
			}
			// nuCode = nil
		case nodeAppl:
			me.Stack = append(me.Stack, n.Callee)
			nuCode = code{{Op: INSTR_UNWIND}}
		case nodeGlobal:
			if (len(me.Stack) - 1) < n.NumArgs {
				panic("unwinding with too few arguments")
			}
			nuCode = n.Code
		default:
			panic(n)
		}
	default:
		panic(cur.Op)
	}
	return nuCode
}
