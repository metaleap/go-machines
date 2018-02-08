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
	INSTR_PUSHARG
	INSTR_MAKEAPPL
	INSTR_SLIDE
	INSTR_UPDATE
	INSTR_POP
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
		return "Push`" + me.Name
	case INSTR_PUSHINT:
		return "Push=" + strconv.Itoa(me.Int)
	case INSTR_PUSHARG:
		return "Push@" + strconv.Itoa(me.Int)
	case INSTR_SLIDE:
		return "Slide:" + strconv.Itoa(me.Int)
	case INSTR_MAKEAPPL:
		return "MkAp"
	case INSTR_UPDATE:
		return "Upd@" + strconv.Itoa(me.Int)
	case INSTR_POP:
		return "Pop@" + strconv.Itoa(me.Int)
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
		addr := me.alloc(nodeLitUint(cur.Int))
		me.Stack = append(me.Stack, addr)
	case INSTR_MAKEAPPL:
		addrcallee := me.Stack[stackpos]
		addrarg := me.Stack[stackpos-1]
		addr := me.alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
		me.Stack[stackpos-1] = addr
		me.Stack = me.Stack[:len(me.Stack)-1]
	case INSTR_PUSHARG:
		addrarg := me.Heap[me.Stack[stackpos-(1+cur.Int)]].(nodeAppl).Arg
		me.Stack = append(me.Stack, addrarg)
	case INSTR_SLIDE:
		keep := me.Stack[stackpos]
		// less := me.Stack[:len(me.Stack)-(1+cur.Int)]
		// me.Stack = append(less, keep)
		me.Stack = me.Stack[:len(me.Stack)-cur.Int]
		me.Stack[len(me.Stack)-1] = keep
	case INSTR_UPDATE:
		pointee := me.Stack[stackpos]
		addrptr := me.alloc(nodePointTo{Addr: pointee})
		me.Stack = me.Stack[:len(me.Stack)-1]
		me.Stack[len(me.Stack)-(1+cur.Int)] = addrptr
	case INSTR_POP:
		me.Stack = me.Stack[:len(me.Stack)-cur.Int]
	case INSTR_UNWIND:
		addr := me.Stack[stackpos]
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodePointTo:
			me.Stack[stackpos] = n.Addr
			nuCode = code{{Op: INSTR_UNWIND}} // if bugs, revert to (cur:nuCode)
		case nodeLitUint:
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
