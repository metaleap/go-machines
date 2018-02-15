package climpl

import (
	"fmt"

	"github.com/metaleap/go-corelang/util"
)

const MARK7 = false // dont set true yet. compilation part of "mark 7" section (p143ff) still missing the updates to the R scheme (p147)

type gMachine struct {
	Heap      clutil.HeapA // no GC here, forever growing
	Globals   clutil.Env
	Code      code          // evaluated l2r
	StackA    clutil.StackA // push-to and pop-from its end
	StackDump []dumpedState
	StackInts clutil.StackI
	Stats     clutil.Stats
}

type dumpedState struct {
	Code  code
	Stack clutil.StackA
}

func (me *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	me.StackA, me.StackDump, me.StackInts = make(clutil.StackA, 0, 64), make([]dumpedState, 0, 16), make(clutil.StackI, 0, 64)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_EVAL}}
	println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
	me.eval()
	stats, val = me.Stats, me.Heap[me.StackA.Top(0)]
	return
}

func (me *gMachine) eval() {
	for me.Stats.NumSteps, me.Stats.NumAppls, me.Stats.MaxStack = 0, 0, 0; len(me.Code) != 0; me.step() {
		if me.Stats.HeapSize = len(me.Heap); me.Stats.MaxStack < len(me.StackA) {
			me.Stats.MaxStack = len(me.StackA)
		}
	}
}

func (me *gMachine) step() {
	me.Stats.NumSteps++
	next := me.Code[1:]

	const cur = 0
	switch me.Code[cur].Op {
	case INSTR_PUSHGLOBAL:
		addr := me.Globals.LookupOrPanic(me.Code[cur].Name)
		me.StackA.Push(addr)
	case INSTR_PUSHINT:
		addr := me.Heap.Alloc(nodeInt(me.Code[cur].Int))
		me.StackA.Push(addr)
	case INSTR_PUSHARG:
		me.StackA.Push(me.StackA.Top(me.Code[cur].Int))
	case INSTR_MAKEAPPL:
		addrcallee := me.StackA.Top(0)
		addrarg := me.StackA.Top(1)
		addr := me.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
		me.StackA[me.StackA.Pos(1)] = addr
		me.StackA = me.StackA.Dropped(1)
	case INSTR_SLIDE:
		keep := me.StackA.Top(0)
		me.StackA = me.StackA.Dropped(me.Code[cur].Int)
		me.StackA[me.StackA.Pos(0)] = keep
	case INSTR_UPDATE:
		pointee := me.StackA.Top(0)
		addrptr := me.Heap.Alloc(nodeIndirection{Addr: pointee})
		me.StackA = me.StackA.Dropped(1)
		me.StackA[me.StackA.Pos(me.Code[cur].Int)] = addrptr
	case INSTR_POP:
		me.StackA = me.StackA.Dropped(me.Code[cur].Int)
	case INSTR_ALLOC:
		for i := 0; i < me.Code[cur].Int; i++ {
			me.StackA.Push(me.Heap.Alloc(nodeIndirection{}))
		}
	case INSTR_EVAL:
		pos := me.StackA.Pos(0)
		me.StackDump = append(me.StackDump, dumpedState{Code: next, Stack: me.StackA[:pos]})
		me.StackA = me.StackA[pos:]
		next = code{{Op: INSTR_UNWIND}}
	case INSTR_PRIM_CMP_EQ, INSTR_PRIM_CMP_NEQ, INSTR_PRIM_CMP_LT, INSTR_PRIM_CMP_LEQ, INSTR_PRIM_CMP_GT, INSTR_PRIM_CMP_GEQ:
		if MARK7 {
			num1, num2 := me.StackInts.Top(0), me.StackInts.Top(1)
			var istrue bool
			switch me.Code[cur].Op {
			case INSTR_PRIM_CMP_EQ:
				istrue = (num1 == num2)
			case INSTR_PRIM_CMP_NEQ:
				istrue = (num1 != num2)
			case INSTR_PRIM_CMP_LT:
				istrue = (num1 < num2)
			case INSTR_PRIM_CMP_LEQ:
				istrue = (num1 <= num2)
			case INSTR_PRIM_CMP_GT:
				istrue = (num1 > num2)
			case INSTR_PRIM_CMP_GEQ:
				istrue = (num1 >= num2)
			}
			var result int
			if istrue {
				result = 2
			} else {
				result = 1
			}
			me.StackInts = me.StackInts.Dropped(1)
			me.StackInts[me.StackInts.Pos(0)] = result
		} else {
			node1, node2 := me.Heap[me.StackA.Top(0)].(nodeInt), me.Heap[me.StackA.Top(1)].(nodeInt)
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
			var result nodeCtor
			if istrue {
				result.Tag = 2
			} else {
				result.Tag = 1
			}
			addr := me.Heap.Alloc(result)
			me.StackA = me.StackA.Dropped(1)
			me.StackA[me.StackA.Pos(0)] = addr
		}
	case INSTR_PRIM_AR_ADD, INSTR_PRIM_AR_SUB, INSTR_PRIM_AR_MUL, INSTR_PRIM_AR_DIV:
		if MARK7 {
			num1, num2 := me.StackInts.Top(0), me.StackInts.Top(1)
			var result int
			switch me.Code[cur].Op {
			case INSTR_PRIM_AR_ADD:
				result = num1 + num2
			case INSTR_PRIM_AR_SUB:
				result = num1 - num2
			case INSTR_PRIM_AR_MUL:
				result = num1 * num2
			case INSTR_PRIM_AR_DIV:
				result = num1 / num2
			}
			me.StackInts = me.StackInts.Dropped(1)
			me.StackInts[me.StackInts.Pos(0)] = result
		} else {
			node1, node2 := me.Heap[me.StackA.Top(0)].(nodeInt), me.Heap[me.StackA.Top(1)].(nodeInt)
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
			me.StackA = me.StackA.Dropped(1)
			me.StackA[me.StackA.Pos(0)] = addr
		}
	case INSTR_PRIM_AR_NEG:
		if MARK7 {
			me.StackInts[me.StackInts.Pos(0)] = -me.StackInts[me.StackInts.Pos(0)]
		} else {
			node := me.Heap[me.StackA.Top(0)].(nodeInt)
			addr := me.Heap.Alloc(-node)
			me.StackA[me.StackA.Pos(0)] = addr
		}
	case INSTR_PRIM_COND:
		if MARK7 {
			bnum := me.StackInts.Top(0)
			me.StackInts = me.StackInts.Dropped(1)
			if bnum == 2 {
				next = append(me.Code[cur].CondThen, next...)
			} else if bnum == 1 {
				next = append(me.Code[cur].CondElse, next...)
			} else {
				panic(bnum)
			}
		} else {
			if node := me.Heap[me.StackA.Top(0)].(nodeCtor); node.Tag == 2 {
				next = append(me.Code[cur].CondThen, next...)
			} else if node.Tag == 1 {
				next = append(me.Code[cur].CondElse, next...)
			} else {
				panic(node.Tag)
			}
			me.StackA = me.StackA.Dropped(1)
		}
	case INSTR_CTOR_PACK:
		arity := me.Code[cur].CtorArity
		node := nodeCtor{Tag: me.Code[cur].Int, Items: make([]clutil.Addr, arity)}
		for i := 0; i < arity; i++ {
			node.Items[i] = me.StackA.Top(i)
		}
		me.StackA = me.StackA.Dropped(arity).Pushed(me.Heap.Alloc(node))
	case INSTR_CASE_JUMP:
		node := me.Heap[me.StackA.Top(0)].(nodeCtor)
		next = append(me.Code[cur].CaseJump[node.Tag], next...)
	case INSTR_CASE_SPLIT:
		node := me.Heap[me.StackA.Top(0)].(nodeCtor)
		me.StackA = me.StackA.Dropped(1)
		for i := /*len(node.Items)*/ me.Code[cur].Int - 1; i > -1; i-- {
			me.StackA.Push(node.Items[i])
		}
	case INSTR_MARK7_PUSHINTVAL:
		me.StackInts.Push(me.Code[cur].Int)
	case INSTR_MARK7_MAKENODEBOOL:
		me.StackA.Push(me.Heap.Alloc(nodeCtor{Tag: me.StackInts.Top(0)}))
		me.StackInts = me.StackInts.Dropped(1)
	case INSTR_MARK7_MAKENODEINT:
		me.StackA.Push(me.Heap.Alloc(nodeInt(me.StackInts.Top(0))))
		me.StackInts = me.StackInts.Dropped(1)
	case INSTR_MARK7_PUSHNODEINT:
		addr := me.StackA.Top(0)
		me.StackA = me.StackA.Dropped(1)
		switch node := me.Heap[addr].(type) {
		case nodeCtor:
			me.StackInts.Push(node.Tag)
		case nodeInt:
			me.StackInts.Push(int(node))
		}
	case INSTR_UNWIND:
		addr := me.StackA.Top(0)
		node := me.Heap[addr]
		switch n := node.(type) {
		case nodeInt, nodeCtor:
			if len(me.StackDump) == 0 {
				next = nil
			} else {
				restore := me.StackDump[len(me.StackDump)-1]
				next, me.StackDump, me.StackA =
					restore.Code, me.StackDump[:len(me.StackDump)-1], append(restore.Stack, addr)
			}
		case nodeIndirection:
			me.StackA[me.StackA.Pos(0)] = n.Addr
			next = code{instr{Op: INSTR_UNWIND}} // unwind again
		case nodeAppl:
			me.Stats.NumAppls++
			me.StackA.Push(n.Callee)
			next = code{instr{Op: INSTR_UNWIND}} // unwind again
		case nodeGlobal:
			if (len(me.StackA) - 1) < n.NumArgs {
				if len(me.StackDump) == 0 {
					panic("unwinding with too few arguments")
				}
				restore := me.StackDump[len(me.StackDump)-1]
				me.StackDump = me.StackDump[:len(me.StackDump)-1]
				next = restore.Code
				me.StackA = restore.Stack.Pushed(me.StackA[0])
			} else {
				nustack := make(clutil.StackA, 0, n.NumArgs)
				for i := n.NumArgs; i > 0; i-- {
					nustack.Push(me.Heap[me.StackA.Top(i)].(nodeAppl).Arg)
				}
				me.StackA = append(me.StackA.Dropped(n.NumArgs), nustack...)
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

func (me *gMachine) String(result interface{}) string {
	if ctor, ok := result.(nodeCtor); ok {
		s := fmt.Sprintf("‹%d", ctor.Tag)
		for _, addr := range ctor.Items {
			s += " " + me.String(me.Heap[addr])
		}
		return s + "›"
	}
	return fmt.Sprintf("%#v", result)
}
