package climpl

import (
	"github.com/metaleap/go-corelang/util"
)

type gMachine struct {
	Stack   clutil.Stack // push-to and pop-from its end
	Heap    clutil.Heap  // no GC here, forever growing
	Globals clutil.Env
	Code    code // evaluated l2r
	Dump    []gDumpItem
	Stats   clutil.Stats
}

type gDumpItem struct {
	Code  code
	Stack clutil.Stack
}

func (me *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	// defer clutil.Catch(&err)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_UNWIND}}
	println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
	me.eval()
	stats, val = me.Stats, me.Heap[me.Stack.Top(0)]
	return
}

func (me *gMachine) eval() {
	for me.Stats.NumSteps, me.Stats.NumAppls = 0, 0; len(me.Code) != 0; me.step() {
	}
}

func (me *gMachine) step() {
	me.Stats.NumSteps++
	next := me.Code[1:]

	const cur = 0
	switch me.Code[cur].Op {
	case INSTR_PUSHGLOBAL:
		addr := me.Globals.LookupOrPanic(me.Code[cur].Name)
		me.Stack.Push(addr)
	case INSTR_PUSHINT:
		addr := me.Heap.Alloc(nodeLitUint(me.Code[cur].Int))
		me.Stack.Push(addr)
	case INSTR_PUSHARG:
		if MARK3_REARRANGESTACK {
			me.Stack.Push(me.Stack.Top(me.Code[cur].Int))
		} else {
			addrarg := me.Heap[me.Stack.Top(1+me.Code[cur].Int)].(nodeAppl).Arg
			me.Stack.Push(addrarg)
		}
	case INSTR_MAKEAPPL:
		addrcallee := me.Stack.Top(0)
		addrarg := me.Stack.Top(1)
		addr := me.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
		me.Stack[me.Stack.Pos(1)] = addr
		me.Stack = me.Stack.Dropped(1)
	case INSTR_SLIDE:
		keep := me.Stack.Top(0)
		me.Stack = me.Stack.Dropped(me.Code[cur].Int)
		me.Stack[me.Stack.Pos(0)] = keep
	case INSTR_UPDATE:
		pointee := me.Stack.Top(0)
		addrptr := me.Heap.Alloc(nodeIndirection{Addr: pointee})
		me.Stack = me.Stack.Dropped(1)
		me.Stack[me.Stack.Pos(me.Code[cur].Int)] = addrptr
	case INSTR_POP:
		me.Stack = me.Stack.Dropped(me.Code[cur].Int)
	case INSTR_ALLOC:
		for i := 0; i < me.Code[cur].Int; i++ {
			me.Stack.Push(me.Heap.Alloc(nodeIndirection{}))
		}
	case INSTR_UNWIND:
		addr := me.Stack.Top(0)
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodeLitUint:
			next = nil
		case nodeIndirection:
			me.Stack[me.Stack.Pos(0)] = n.Addr
			next = code{instr{Op: INSTR_UNWIND}} // unwind again
		case nodeAppl:
			me.Stats.NumAppls++
			me.Stack.Push(n.Callee)
			next = code{instr{Op: INSTR_UNWIND}} // unwind again
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
		panic(me.Code[cur].Op)
	}
	me.Code = next
}
