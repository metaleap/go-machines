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

func (me *gMachine) dispatch(cur instr, next code) code {
	stackpos := me.Stack.Pos(0)
	switch cur.Op {
	case INSTR_PUSHGLOBAL:
		addr := me.Globals.Lookup(cur.Name)
		me.Stack.Push(addr)
	case INSTR_PUSHINT:
		addr := me.Heap.Alloc(nodeLitUint(cur.Int))
		me.Stack.Push(addr)
	case INSTR_MAKEAPPL:
		addrcallee := me.Stack.Top(0)
		addrarg := me.Stack.Top(1)
		addr := me.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
		me.Stack[stackpos-1] = addr
		me.Stack = me.Stack.Dropped(1)
	case INSTR_PUSHARG:
		addrarg := me.Heap[me.Stack.Top(1+cur.Int)].(nodeAppl).Arg
		me.Stack.Push(addrarg)
	case INSTR_SLIDE:
		keep := me.Stack.Top(0)
		// less := me.Stack[:len(me.Stack)-(1+cur.Int)]
		// me.Stack = append(less, keep)
		me.Stack = me.Stack.Dropped(cur.Int)
		me.Stack[len(me.Stack)-1] = keep
	case INSTR_UPDATE:
		pointee := me.Stack.Top(0)
		addrptr := me.Heap.Alloc(nodeIndirection{Addr: pointee})
		me.Stack = me.Stack.Dropped(1)
		me.Stack[len(me.Stack)-(1+cur.Int)] = addrptr
	case INSTR_POP:
		me.Stack = me.Stack.Dropped(cur.Int)
	case INSTR_UNWIND:
		addr := me.Stack.Top(0)
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodeLitUint:
			// nothing to do
		case nodeIndirection:
			me.Stack[stackpos] = n.Addr
			if len(next) > 0 { // temporarily to observe
				panic("does dis ever happen?")
				// nuCode = append(code{cur}, nuCode...)
			}
			next = code{cur}
		case nodeAppl:
			me.Stats.NumAppls++
			me.Stack.Push(n.Callee)
			next = code{cur}
		case nodeGlobal:
			if (len(me.Stack) - 1) < n.NumArgs {
				panic("unwinding with too few arguments")
			}
			next = n.Code
		default:
			panic(n)
		}
	default:
		panic(cur.Op)
	}
	return next
}
