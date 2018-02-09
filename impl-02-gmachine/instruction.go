package climpl

import (
	"strconv"

	"github.com/metaleap/go-corelang/util"
)

const MARK3_REARRANGESTACK = true

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
	INSTR_ALLOC
)

type instr struct {
	Op   instruction
	Int  int
	Name string
}

func (me instr) String() string {
	switch me.Op {
	case INSTR_UNWIND:
		return "Unwd"
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
	case INSTR_ALLOC:
		return "Alloc=" + strconv.Itoa(me.Int)
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
	switch cur.Op {
	case INSTR_PUSHGLOBAL:
		addr := me.Globals.LookupOrPanic(cur.Name)
		me.Stack.Push(addr)
	case INSTR_PUSHINT:
		addr := me.Heap.Alloc(nodeLitUint(cur.Int))
		me.Stack.Push(addr)
	case INSTR_MAKEAPPL:
		addrcallee := me.Stack.Top(0)
		addrarg := me.Stack.Top(1)
		addr := me.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
		me.Stack[me.Stack.Pos(1)] = addr
		me.Stack = me.Stack.Dropped(1)
	case INSTR_PUSHARG:
		if MARK3_REARRANGESTACK {
			me.Stack.Push(me.Stack.Top(cur.Int))
		} else {
			addrarg := me.Heap[me.Stack.Top(1+cur.Int)].(nodeAppl).Arg
			me.Stack.Push(addrarg)
		}
	case INSTR_SLIDE:
		keep := me.Stack.Top(0)
		me.Stack = me.Stack.Dropped(cur.Int)
		me.Stack[me.Stack.Pos(0)] = keep
	case INSTR_UPDATE:
		pointee := me.Stack.Top(0)
		addrptr := me.Heap.Alloc(nodeIndirection{Addr: pointee})
		me.Stack = me.Stack.Dropped(1)
		me.Stack[me.Stack.Pos(cur.Int)] = addrptr
	case INSTR_POP:
		me.Stack = me.Stack.Dropped(cur.Int)
	case INSTR_ALLOC:
		for i := 0; i < cur.Int; i++ {
			me.Stack.Push(me.Heap.Alloc(nodeIndirection{}))
		}
	case INSTR_UNWIND:
		addr := me.Stack.Top(0)
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodeLitUint:
			if len(next) > 0 { // temporarily to observe
				panic("nodeLitUint: code remaining")
				// next =nil
			}
		case nodeIndirection:
			me.Stack[me.Stack.Pos(0)] = n.Addr
			if len(next) > 0 { // temporarily to observe
				panic("nodeIndirection: code remaining")
				// next = append(code{cur}, next...)
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
			if MARK3_REARRANGESTACK {
				nustack := make(clutil.Stack, 0, n.NumArgs)
				for i := n.NumArgs; i > 0; i-- {
					nustack.Push(me.Heap[me.Stack.Top(i)].(nodeAppl).Arg)
				}
				me.Stack = append(me.Stack.Dropped(n.NumArgs), nustack...)
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
