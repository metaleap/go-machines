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
	JumpTable []stepInstr
}

type stepInstr func(*code)

type dumpedState struct {
	Code  code
	Stack clutil.StackA
}

func (me *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	me.StackA, me.StackDump, me.StackInts = make(clutil.StackA, 0, 64), make([]dumpedState, 0, 16), make(clutil.StackI, 0, 64)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_EVAL}}
	// println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
	me.eval()
	stats, val = me.Stats, me.Heap[me.StackA.Top(0)]
	return
}

func (me *gMachine) eval() {
	for me.Stats.NumSteps, me.Stats.NumAppls, me.Stats.MaxStack = 0, 0, 0; len(me.Code) != 0; me.Stats.HeapSize = len(me.Heap) {
		me.Stats.NumSteps++
		next := me.Code[1:]

		switch me.Code[0].Op {
		case INSTR_EVAL:
			me.step_INSTR_EVAL(&next)
		case INSTR_UNWIND:
			me.step_INSTR_UNWIND(&next)
		case INSTR_PUSHGLOBAL:
			me.step_INSTR_PUSHGLOBAL(nil)
		case INSTR_PUSHINT:
			me.step_INSTR_PUSHINT(nil)
		case INSTR_PUSHARG:
			me.step_INSTR_PUSHARG(nil)
		case INSTR_MAKEAPPL:
			me.step_INSTR_MAKEAPPL(nil)
		case INSTR_UPDATE:
			me.step_INSTR_UPDATE(nil)
		case INSTR_POP:
			me.step_INSTR_POP(nil)
		case INSTR_SLIDE:
			me.step_INSTR_SLIDE(nil)
		case INSTR_ALLOC:
			me.step_INSTR_ALLOC(nil)
		case INSTR_PRIM_CMP_EQ, INSTR_PRIM_CMP_NEQ, INSTR_PRIM_CMP_LT, INSTR_PRIM_CMP_LEQ, INSTR_PRIM_CMP_GT, INSTR_PRIM_CMP_GEQ:
			me.step_INSTR_PRIM_CMP(nil)
		case INSTR_PRIM_AR_ADD, INSTR_PRIM_AR_SUB, INSTR_PRIM_AR_MUL, INSTR_PRIM_AR_DIV:
			me.step_INSTR_PRIM_AR(nil)
		case INSTR_PRIM_AR_NEG:
			me.step_INSTR_PRIM_AR_NEG(nil)
		case INSTR_PRIM_COND:
			me.step_INSTR_PRIM_COND(&next)
		case INSTR_CTOR_PACK:
			me.step_INSTR_CTOR_PACK(nil)
		case INSTR_CASE_JUMP:
			me.step_INSTR_CASE_JUMP(&next)
		case INSTR_CASE_SPLIT:
			me.step_INSTR_CASE_SPLIT(nil)
		case INSTR_MARK7_PUSHINTVAL:
			me.step_INSTR_MARK7_PUSHINTVAL(nil)
		case INSTR_MARK7_MAKENODEBOOL:
			me.step_INSTR_MARK7_MAKENODEBOOL(nil)
		case INSTR_MARK7_MAKENODEINT:
			me.step_INSTR_MARK7_MAKENODEINT(nil)
		case INSTR_MARK7_PUSHNODEINT:
			me.step_INSTR_MARK7_PUSHNODEINT(nil)
		default:
			panic(me.Code[0].Op)
		}
		me.Code = next

		if me.Stats.MaxStack < len(me.StackA) {
			me.Stats.MaxStack = len(me.StackA)
		}
	}
}

func (me *gMachine) _step_old_oddlySlowerThanSwitchBlockSoUnusedButHadToTry() {
	// reminder to self for future proto VMs, the earlier big-uber-switch with direct-code per case (instead of method dispatch) was simply fastest, contrary to commonly held wisdom about a method-jumptable being preferable to dozens of cmp ops in the switch!
	me.Stats.NumSteps++
	next := me.Code[1:]
	me.JumpTable[me.Code[0].Op](&next)
	me.Code = next
}

func (me *gMachine) step_INSTR_PUSHGLOBAL(*code) {
	addr := me.Globals.LookupOrPanic(me.Code[0].Name)
	me.StackA.Push(addr)
}

func (me *gMachine) step_INSTR_PUSHINT(*code) {
	addr := me.Heap.Alloc(nodeInt(me.Code[0].Int))
	me.StackA.Push(addr)
}

func (me *gMachine) step_INSTR_PUSHARG(*code) {
	me.StackA.Push(me.StackA.Top(me.Code[0].Int))
}

func (me *gMachine) step_INSTR_MAKEAPPL(*code) {
	addrcallee := me.StackA.Top(0)
	addrarg := me.StackA.Top(1)
	addr := me.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
	me.StackA[me.StackA.Pos(1)] = addr
	me.StackA = me.StackA.Dropped(1)
}

func (me *gMachine) step_INSTR_SLIDE(*code) {
	keep := me.StackA.Top(0)
	me.StackA = me.StackA.Dropped(me.Code[0].Int)
	me.StackA[me.StackA.Pos(0)] = keep
}

func (me *gMachine) step_INSTR_UPDATE(*code) {
	pointee := me.StackA.Top(0)
	addrptr := me.Heap.Alloc(nodeIndirection{Addr: pointee})
	me.StackA = me.StackA.Dropped(1)
	me.StackA[me.StackA.Pos(me.Code[0].Int)] = addrptr
}

func (me *gMachine) step_INSTR_POP(*code) {
	me.StackA = me.StackA.Dropped(me.Code[0].Int)
}

func (me *gMachine) step_INSTR_ALLOC(*code) {
	for i := 0; i < me.Code[0].Int; i++ {
		me.StackA.Push(me.Heap.Alloc(nodeIndirection{}))
	}
}

func (me *gMachine) step_INSTR_EVAL(next *code) {
	pos := me.StackA.Pos(0)
	me.StackDump = append(me.StackDump, dumpedState{Code: *next, Stack: me.StackA[:pos]})
	me.StackA = me.StackA[pos:]
	*next = code{{Op: INSTR_UNWIND}}
}

func (me *gMachine) step_INSTR_CTOR_PACK(*code) {
	arity := me.Code[0].CtorArity
	node := nodeCtor{Tag: me.Code[0].Int, Items: make([]clutil.Addr, arity)}
	for i := 0; i < arity; i++ {
		node.Items[i] = me.StackA.Top(i)
	}
	me.StackA = me.StackA.Dropped(arity).Pushed(me.Heap.Alloc(node))
}

func (me *gMachine) step_INSTR_CASE_JUMP(next *code) {
	node := me.Heap[me.StackA.Top(0)].(nodeCtor)
	*next = append(me.Code[0].CaseJump[node.Tag], *next...)
}

func (me *gMachine) step_INSTR_CASE_SPLIT(*code) {
	node := me.Heap[me.StackA.Top(0)].(nodeCtor)
	me.StackA = me.StackA.Dropped(1)
	for i := /*len(node.Items)*/ me.Code[0].Int - 1; i > -1; i-- {
		me.StackA.Push(node.Items[i])
	}
}

func (me *gMachine) step_INSTR_MARK7_PUSHINTVAL(*code) {
	me.StackInts.Push(me.Code[0].Int)
}

func (me *gMachine) step_INSTR_MARK7_MAKENODEBOOL(*code) {
	me.StackA.Push(me.Heap.Alloc(nodeCtor{Tag: me.StackInts.Top(0)}))
	me.StackInts = me.StackInts.Dropped(1)
}

func (me *gMachine) step_INSTR_MARK7_MAKENODEINT(*code) {
	me.StackA.Push(me.Heap.Alloc(nodeInt(me.StackInts.Top(0))))
	me.StackInts = me.StackInts.Dropped(1)
}

func (me *gMachine) step_INSTR_MARK7_PUSHNODEINT(*code) {
	addr := me.StackA.Top(0)
	me.StackA = me.StackA.Dropped(1)
	switch node := me.Heap[addr].(type) {
	case nodeCtor:
		me.StackInts.Push(node.Tag)
	case nodeInt:
		me.StackInts.Push(int(node))
	}
}

func (me *gMachine) step_INSTR_PRIM_COND(next *code) {
	if MARK7 {
		bnum := me.StackInts.Top(0)
		me.StackInts = me.StackInts.Dropped(1)
		if bnum == 2 {
			*next = append(me.Code[0].CondThen, *next...)
		} else if bnum == 1 {
			*next = append(me.Code[0].CondElse, *next...)
		} else {
			panic(bnum)
		}
	} else {
		if node := me.Heap[me.StackA.Top(0)].(nodeCtor); node.Tag == 2 {
			*next = append(me.Code[0].CondThen, *next...)
		} else if node.Tag == 1 {
			*next = append(me.Code[0].CondElse, *next...)
		} else {
			panic(node.Tag)
		}
		me.StackA = me.StackA.Dropped(1)
	}
}

func (me *gMachine) step_INSTR_PRIM_AR_NEG(*code) {
	if MARK7 {
		me.StackInts[me.StackInts.Pos(0)] = -me.StackInts[me.StackInts.Pos(0)]
	} else {
		node := me.Heap[me.StackA.Top(0)].(nodeInt)
		addr := me.Heap.Alloc(-node)
		me.StackA[me.StackA.Pos(0)] = addr
	}
}

func (me *gMachine) step_INSTR_PRIM_AR(*code) {
	if MARK7 {
		num1, num2 := me.StackInts.Top(0), me.StackInts.Top(1)
		var result int
		switch me.Code[0].Op {
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
		switch me.Code[0].Op {
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
}

func (me *gMachine) step_INSTR_PRIM_CMP(*code) {
	if MARK7 {
		num1, num2 := me.StackInts.Top(0), me.StackInts.Top(1)
		var istrue bool
		switch me.Code[0].Op {
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
		switch me.Code[0].Op {
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
}

func (me *gMachine) step_INSTR_UNWIND(next *code) {
	addr := me.StackA.Top(0)
	node := me.Heap[addr]
	switch n := node.(type) {
	case nodeInt, nodeCtor:
		if len(me.StackDump) == 0 {
			*next = code{}
		} else {
			restore := me.StackDump[len(me.StackDump)-1]
			*next, me.StackDump, me.StackA =
				restore.Code, me.StackDump[:len(me.StackDump)-1], append(restore.Stack, addr)
		}
	case nodeIndirection:
		me.StackA[me.StackA.Pos(0)] = n.Addr
		*next = code{instr{Op: INSTR_UNWIND}} // unwind again
	case nodeAppl:
		me.Stats.NumAppls++
		me.StackA.Push(n.Callee)
		*next = code{instr{Op: INSTR_UNWIND}} // unwind again
	case nodeGlobal:
		if (len(me.StackA) - 1) < n.NumArgs {
			if len(me.StackDump) == 0 {
				panic("unwinding with too few arguments")
			}
			restore := me.StackDump[len(me.StackDump)-1]
			me.StackDump = me.StackDump[:len(me.StackDump)-1]
			*next = restore.Code
			me.StackA = restore.Stack.Pushed(me.StackA[0])
		} else {
			nustack := make(clutil.StackA, 0, n.NumArgs)
			for i := n.NumArgs; i > 0; i-- {
				nustack.Push(me.Heap[me.StackA.Top(i)].(nodeAppl).Arg)
			}
			me.StackA = append(me.StackA.Dropped(n.NumArgs), nustack...)
			*next = n.Code
		}
	default:
		panic(n)
	}
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
