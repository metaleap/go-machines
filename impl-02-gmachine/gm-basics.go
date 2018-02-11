package climpl

import (
	"github.com/metaleap/go-corelang/util"
)

type gMachine struct {
	Stack   clutil.Stack // push-to and pop-from its end
	Heap    clutil.Heap  // no GC here, forever growing
	Globals clutil.Env
	Code    code // evaluated l2r
	Dump    []dumpItem
	Stats   clutil.Stats
}

type dumpItem struct {
	Code  code
	Stack clutil.Stack
}

func (me *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_EVAL}}
	// println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
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
		addr := me.Heap.Alloc(nodeInt(me.Code[cur].Int))
		me.Stack.Push(addr)
	case INSTR_PUSHARG:
		me.Stack.Push(me.Stack.Top(me.Code[cur].Int))
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
	case INSTR_EVAL:
		pos := me.Stack.Pos(0)
		me.Dump = append(me.Dump, dumpItem{Code: next, Stack: me.Stack[:pos]})
		me.Stack = me.Stack[pos:]
		next = code{{Op: INSTR_UNWIND}}
	case INSTR_PRIM_CMP_EQ, INSTR_PRIM_CMP_NEQ, INSTR_PRIM_CMP_LT, INSTR_PRIM_CMP_LEQ, INSTR_PRIM_CMP_GT, INSTR_PRIM_CMP_GEQ:
		node1, node2 := me.Heap[me.Stack.Top(0)].(nodeInt), me.Heap[me.Stack.Top(1)].(nodeInt)
		var istrue bool
		switch me.Code[cur].Op {
		case INSTR_PRIM_CMP_EQ:
			istrue = (node1 == node2)
		case INSTR_PRIM_CMP_NEQ:
			istrue = (node1 != node2)
		case INSTR_PRIM_CMP_LT:
			istrue = (node1 < node2)
		case INSTR_PRIM_CMP_LEQ:
			istrue = (node1 <= node2)
		case INSTR_PRIM_CMP_GT:
			istrue = (node1 > node2)
		case INSTR_PRIM_CMP_GEQ:
			istrue = (node1 >= node2)
		}
		var result nodeInt
		if istrue {
			result = 1
		}
		addr := me.Heap.Alloc(result)
		me.Stack = me.Stack.Dropped(1)
		me.Stack[me.Stack.Pos(0)] = addr
	case INSTR_PRIM_AR_ADD, INSTR_PRIM_AR_SUB, INSTR_PRIM_AR_MUL, INSTR_PRIM_AR_DIV:
		node1, node2 := me.Heap[me.Stack.Top(0)].(nodeInt), me.Heap[me.Stack.Top(1)].(nodeInt)
		var result nodeInt
		switch me.Code[cur].Op {
		case INSTR_PRIM_AR_ADD:
			result = node1 + node2
		case INSTR_PRIM_AR_SUB:
			result = node1 - node2
		case INSTR_PRIM_AR_MUL:
			result = node1 * node2
		case INSTR_PRIM_AR_DIV:
			result = node1 / node2
		}
		addr := me.Heap.Alloc(result)
		me.Stack = me.Stack.Dropped(1)
		me.Stack[me.Stack.Pos(0)] = addr
	case INSTR_PRIM_AR_NEG:
		node := me.Heap[me.Stack.Top(0)].(nodeInt)
		addr := me.Heap.Alloc(-node)
		me.Stack[me.Stack.Pos(0)] = addr
	case INSTR_PRIM_COND:
		if node := me.Heap[me.Stack.Top(0)].(nodeInt); node == 1 {
			next = append(me.Code[0].CondThen, next...)
		} else if node == 0 {
			next = append(me.Code[0].CondElse, next...)
		} else {
			panic("boolean bug")
		}
		me.Stack = me.Stack.Dropped(1)
	case INSTR_CTOR_PACK:
		arity := me.Code[cur].CtorArity
		node := nodeCtor{Tag: me.Code[cur].Int, Items: make([]clutil.Addr, arity)}
		for i := 0; i < arity; i++ {
			node.Items[i] = me.Stack.Top(i)
		}
		me.Stack = me.Stack.Dropped(arity).Pushed(me.Heap.Alloc(node))
	case INSTR_CASE_JUMP:
		node := me.Heap[me.Stack.Top(0)].(nodeCtor)
		next = append(me.Code[cur].CaseJump[node.Tag], next...)
	case INSTR_CASE_SPLIT:
		node := me.Heap[me.Stack.Top(0)].(nodeCtor)
		me.Stack = me.Stack.Dropped(1)
		for i := len(node.Items) - 1; i > -1; i-- {
			me.Stack.Push(node.Items[i])
		}
	case INSTR_UNWIND:
		addr := me.Stack.Top(0)
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodeInt, nodeCtor:
			if len(me.Dump) == 0 {
				next = nil
			} else {
				restore := me.Dump[len(me.Dump)-1]
				next, me.Dump, me.Stack =
					restore.Code, me.Dump[:len(me.Dump)-1], append(restore.Stack, addr)
			}
		case nodeIndirection:
			me.Stack[me.Stack.Pos(0)] = n.Addr
			next = code{instr{Op: INSTR_UNWIND}} // unwind again
		case nodeAppl:
			me.Stats.NumAppls++
			me.Stack.Push(n.Callee)
			next = code{instr{Op: INSTR_UNWIND}} // unwind again
		case nodeGlobal:
			if (len(me.Stack) - 1) < n.NumArgs {
				if len(me.Dump) == 0 {
					panic("unwinding with too few arguments")
				}
				restore := me.Dump[len(me.Dump)-1]
				me.Dump = me.Dump[:len(me.Dump)-1]
				next = restore.Code
				me.Stack = restore.Stack.Pushed(me.Stack[0])
			} else {
				nustack := make(clutil.Stack, 0, n.NumArgs)
				for i := n.NumArgs; i > 0; i-- {
					nustack.Push(me.Heap[me.Stack.Top(i)].(nodeAppl).Arg)
				}
				me.Stack = append(me.Stack.Dropped(n.NumArgs), nustack...)
				next = n.Code
			}
		default:
			panic(n)
		}
	default:
		panic(me.Code[cur].Op)
	}
	me.Code = next
}
