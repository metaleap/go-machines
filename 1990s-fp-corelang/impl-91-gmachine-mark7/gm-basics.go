package climpl

import (
	"strconv"

	"github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

const _MARK7 = false // not a big gain in practice for this unoptimized prototype and its toy examples, still intrinsically a sane (and for real-world likely crucial) approach to have separate val stacks (in addition to addr stack)

type gMachine struct {
	Heap      clutil.HeapA // no GC here, forever growing
	Globals   clutil.Env
	Code      code          // evaluated l2r
	StackA    clutil.StackA // push-to and pop-from its end
	StackDump []dumpedState
	StackInts clutil.StackI // used if _MARK7
	StackStrs clutil.StackS // used if _MARK7 ever since we swapped ctor int tags for str tags
	Stats     clutil.Stats
}

type dumpedState struct {
	Code  code
	Stack clutil.StackA
}

func (me *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	me.StackA, me.StackDump, me.StackInts, me.StackStrs = make(clutil.StackA, 0, 64), make([]dumpedState, 0, 16), make(clutil.StackI, 0, 64), make(clutil.StackS, 0, 16)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_EVAL}}
	// println(me.Heap[me.Globals["times"]].(nodeGlobal).Code.String())
	me.eval()
	stats, val = me.Stats, me.Heap[me.StackA.Top0()]
	return
}

func (me *gMachine) eval() {
	for me.Stats.NumSteps, me.Stats.NumAppls, me.Stats.MaxStack = 0, 0, 0; len(me.Code) > 0; me.Stats.NumSteps++ {
		next := me.Code[1:]

		switch me.Code[0].Op {
		case INSTR_PUSHGLOBAL:
			addr := me.Globals.LookupOrPanic(me.Code[0].Name)
			me.StackA.Push(addr)
		case INSTR_PUSHINT:
			addr := me.Heap.Alloc(nodeInt(me.Code[0].Int))
			me.StackA.Push(addr)
		case INSTR_PUSHARG:
			me.StackA.Push(me.StackA.Top(me.Code[0].Int))
		case INSTR_MAKEAPPL:
			addrcallee := me.StackA.Top0()
			addrarg := me.StackA.Top1()
			addr := me.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
			me.StackA[me.StackA.Pos1()] = addr
			me.StackA = me.StackA.Dropped(1)
		case INSTR_UPDATE:
			pointee := me.StackA.Top0()
			addrptr := me.Heap.Alloc(nodeIndirection{Addr: pointee})
			me.StackA = me.StackA.Dropped(1)
			me.StackA[me.StackA.Pos(me.Code[0].Int)] = addrptr
		case INSTR_POP:
			me.StackA = me.StackA.Dropped(me.Code[0].Int)
		case INSTR_SLIDE:
			keep := me.StackA.Top0()
			me.StackA = me.StackA.Dropped(me.Code[0].Int)
			me.StackA[me.StackA.Pos0()] = keep
		case INSTR_ALLOC:
			for i := 0; i < me.Code[0].Int; i++ {
				me.StackA.Push(me.Heap.Alloc(nodeIndirection{}))
			}
		case INSTR_EVAL:
			pos := me.StackA.Pos0()
			me.StackDump = append(me.StackDump, dumpedState{Code: next, Stack: me.StackA[:pos]})
			me.StackA = me.StackA[pos:]
			next = code{{Op: INSTR_UNWIND}}
		case INSTR_UNWIND:
			addr := me.StackA.Top0()
			node := me.Heap[addr]
			switch n := node.(type) {
			case nodeInt, nodeCtor:
				if len(me.StackDump) == 0 {
					next = code{}
				} else {
					restore := me.StackDump[len(me.StackDump)-1]
					next, me.StackDump, me.StackA =
						restore.Code, me.StackDump[:len(me.StackDump)-1], append(restore.Stack, addr)
				}
			case nodeIndirection:
				me.StackA[me.StackA.Pos0()] = n.Addr
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
		case INSTR_PRIM_CMP_EQ, INSTR_PRIM_CMP_NEQ, INSTR_PRIM_CMP_LT, INSTR_PRIM_CMP_LEQ, INSTR_PRIM_CMP_GT, INSTR_PRIM_CMP_GEQ:
			if _MARK7 {
				num1, num2 := me.StackInts.Top0(), me.StackInts.Top1()
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
				var result string
				if istrue {
					result = "True"
				} else {
					result = "False"
				}
				me.StackInts = me.StackInts.Dropped(2)
				me.StackStrs.Push(result)
			} else {
				node1, node2 := me.Heap[me.StackA.Top0()].(nodeInt), me.Heap[me.StackA.Top1()].(nodeInt)
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
					result.Tag = "True"
				} else {
					result.Tag = "False"
				}
				addr := me.Heap.Alloc(result)
				me.StackA = me.StackA.Dropped(1)
				me.StackA[me.StackA.Pos0()] = addr
			}
		case INSTR_PRIM_AR_ADD, INSTR_PRIM_AR_SUB, INSTR_PRIM_AR_MUL, INSTR_PRIM_AR_DIV:
			if _MARK7 {
				num1, num2 := me.StackInts.Top0(), me.StackInts.Top1()
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
				me.StackInts[me.StackInts.Pos0()] = result
			} else {
				node1, node2 := me.Heap[me.StackA.Top0()].(nodeInt), me.Heap[me.StackA.Top1()].(nodeInt)
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
				me.StackA[me.StackA.Pos0()] = addr
			}
		case INSTR_PRIM_AR_NEG:
			if _MARK7 {
				me.StackInts[me.StackInts.Pos0()] = -me.StackInts[me.StackInts.Pos0()]
			} else {
				node := me.Heap[me.StackA.Top0()].(nodeInt)
				addr := me.Heap.Alloc(-node)
				me.StackA[me.StackA.Pos0()] = addr
			}
		case INSTR_PRIM_COND:
			if _MARK7 {
				ctortag := me.StackStrs.Top0()
				me.StackStrs = me.StackStrs.Dropped(1)
				if ctortag == "True" {
					next = append(me.Code[0].CondThen, next...)
				} else if ctortag == "False" {
					next = append(me.Code[0].CondElse, next...)
				} else {
					panic(ctortag)
				}
			} else {
				if node := me.Heap[me.StackA.Top0()].(nodeCtor); node.Tag == "True" {
					next = append(me.Code[0].CondThen, next...)
				} else if node.Tag == "False" {
					next = append(me.Code[0].CondElse, next...)
				} else {
					panic(node.Tag)
				}
				me.StackA = me.StackA.Dropped(1)
			}
		case INSTR_CTOR_PACK:
			arity := me.Code[0].CtorArity
			node := nodeCtor{Tag: me.Code[0].Name, Items: make([]clutil.Addr, arity)}
			for i := 0; i < arity; i++ {
				node.Items[i] = me.StackA.Top(i)
			}
			me.StackA = me.StackA.Dropped(arity).Pushed(me.Heap.Alloc(node))
		case INSTR_CASE_JUMP:
			node := me.Heap[me.StackA.Top0()].(nodeCtor)
			if code := me.Code[0].CaseJump[node.Tag]; len(code) > 0 {
				next = append(code, next...)
			} else if code = me.Code[0].CaseJump["_"]; len(code) > 0 { // jump to default case
				next = append(code, next...)
			} else {
				panic("no matching alternative in CASE OF for ‹" + node.Tag + "," + strconv.Itoa(len(node.Items)) + "› and no default (tag 0) alternative either")
			}
		case INSTR_CASE_SPLIT:
			node := me.Heap[me.StackA.Top0()].(nodeCtor)
			me.StackA = me.StackA.Dropped(1)
			for i := /*len(node.Items)*/ me.Code[0].Int - 1; i > -1; i-- {
				me.StackA.Push(node.Items[i])
			}
		case INSTR_MARK7_PUSHINTVAL:
			me.StackInts.Push(me.Code[0].Int)
		case INSTR_MARK7_MAKENODEBOOL:
			me.StackA.Push(me.Heap.Alloc(nodeCtor{Tag: me.StackStrs.Top0()}))
			me.StackStrs = me.StackStrs.Dropped(1)
		case INSTR_MARK7_MAKENODEINT:
			me.StackA.Push(me.Heap.Alloc(nodeInt(me.StackInts.Top0())))
			me.StackInts = me.StackInts.Dropped(1)
		case INSTR_MARK7_PUSHNODEINT:
			addr := me.StackA.Top0()
			me.StackA = me.StackA.Dropped(1)
			switch node := me.Heap[addr].(type) {
			case nodeCtor:
				me.StackStrs.Push(node.Tag)
			case nodeInt:
				me.StackInts.Push(int(node))
			}
		default:
			panic(me.Code[0].Op)
		}

		if me.Code = next; me.Stats.MaxStack < len(me.StackA) {
			me.Stats.MaxStack = len(me.StackA)
		}
		if me.Stats.NumSteps > 999999 {
			panic("exceeded 1 million steps: probable infinite loop, stopping evaluation")
		}
	}
	me.Stats.HeapSize = len(me.Heap)
}

func (me *gMachine) String(result interface{}) string {
	switch res := result.(type) {
	case nodeInt:
		return "#" + strconv.Itoa(int(res))
	case nodeCtor:
		s := "‹" + res.Tag
		for _, addr := range res.Items {
			s += " " + me.String(me.Heap[addr])
		}
		return s + "›"
	case nodeIndirection:
		return "@" + res.Addr.String()
	case nodeGlobal:
		return strconv.Itoa(res.NumArgs) + "@" + res.Code.String()
	case nodeAppl:
		return "(" + res.Callee.String() + " " + res.Arg.String() + ")"
	}
	panic(result)
}
