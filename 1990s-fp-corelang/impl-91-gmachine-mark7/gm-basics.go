package climpl

import (
	"strconv"

	"github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

const _MARK7 = true // not a big gain in practice for this unoptimized prototype and its toy examples, still intrinsically a sane (and for real-world likely crucial) approach to have separate val stacks (in addition to addr stack)

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

func (this *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	this.StackA, this.StackDump, this.StackInts, this.StackStrs = make(clutil.StackA, 0, 64), make([]dumpedState, 0, 16), make(clutil.StackI, 0, 64), make(clutil.StackS, 0, 16)
	this.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_EVAL}}
	// println(me.Heap[me.Globals["times"]].(nodeGlobal).Code.String())
	this.eval()
	stats, val = this.Stats, this.Heap[this.StackA.Top0()]
	return
}

func (this *gMachine) eval() {
	for this.Stats.NumSteps, this.Stats.NumAppls, this.Stats.MaxStack = 0, 0, 0; len(this.Code) > 0; this.Stats.NumSteps++ {
		next := this.Code[1:]

		switch this.Code[0].Op {
		case INSTR_PUSHGLOBAL:
			addr := this.Globals.LookupOrPanic(this.Code[0].Name)
			this.StackA.Push(addr)
		case INSTR_PUSHINT:
			addr := this.Heap.Alloc(nodeInt(this.Code[0].Int))
			this.StackA.Push(addr)
		case INSTR_PUSHARG:
			this.StackA.Push(this.StackA.Top(this.Code[0].Int))
		case INSTR_MAKEAPPL:
			addrcallee := this.StackA.Top0()
			addrarg := this.StackA.Top1()
			addr := this.Heap.Alloc(nodeAppl{Callee: addrcallee, Arg: addrarg})
			this.StackA[this.StackA.Pos1()] = addr
			this.StackA = this.StackA.Dropped(1)
		case INSTR_UPDATE:
			pointee := this.StackA.Top0()
			addrptr := this.Heap.Alloc(nodeIndirection{Addr: pointee})
			this.StackA = this.StackA.Dropped(1)
			this.StackA[this.StackA.Pos(this.Code[0].Int)] = addrptr
		case INSTR_POP:
			this.StackA = this.StackA.Dropped(this.Code[0].Int)
		case INSTR_SLIDE:
			keep := this.StackA.Top0()
			this.StackA = this.StackA.Dropped(this.Code[0].Int)
			this.StackA[this.StackA.Pos0()] = keep
		case INSTR_ALLOC:
			for i := 0; i < this.Code[0].Int; i++ {
				this.StackA.Push(this.Heap.Alloc(nodeIndirection{}))
			}
		case INSTR_EVAL:
			pos := this.StackA.Pos0()
			this.StackDump = append(this.StackDump, dumpedState{Code: next, Stack: this.StackA[:pos]})
			this.StackA = this.StackA[pos:]
			next = code{{Op: INSTR_UNWIND}}
		case INSTR_UNWIND:
			addr := this.StackA.Top0()
			node := this.Heap[addr]
			switch n := node.(type) {
			case nodeInt, nodeCtor:
				if len(this.StackDump) == 0 {
					next = code{}
				} else {
					restore := this.StackDump[len(this.StackDump)-1]
					next, this.StackDump, this.StackA =
						restore.Code, this.StackDump[:len(this.StackDump)-1], append(restore.Stack, addr)
				}
			case nodeIndirection:
				this.StackA[this.StackA.Pos0()] = n.Addr
				next = code{instr{Op: INSTR_UNWIND}} // unwind again
			case nodeAppl:
				this.Stats.NumAppls++
				this.StackA.Push(n.Callee)
				next = code{instr{Op: INSTR_UNWIND}} // unwind again
			case nodeGlobal:
				if (len(this.StackA) - 1) < n.NumArgs {
					if len(this.StackDump) == 0 {
						panic("unwinding with too few arguments")
					}
					restore := this.StackDump[len(this.StackDump)-1]
					this.StackDump = this.StackDump[:len(this.StackDump)-1]
					next = restore.Code
					this.StackA = restore.Stack.Pushed(this.StackA[0])
				} else {
					nustack := make(clutil.StackA, 0, n.NumArgs)
					for i := n.NumArgs; i > 0; i-- {
						nustack.Push(this.Heap[this.StackA.Top(i)].(nodeAppl).Arg)
					}
					this.StackA = append(this.StackA.Dropped(n.NumArgs), nustack...)
					next = n.Code
				}
			default:
				panic(n)
			}
		case INSTR_PRIM_CMP_EQ, INSTR_PRIM_CMP_NEQ, INSTR_PRIM_CMP_LT, INSTR_PRIM_CMP_LEQ, INSTR_PRIM_CMP_GT, INSTR_PRIM_CMP_GEQ:
			if _MARK7 {
				num1, num2 := this.StackInts.Top0(), this.StackInts.Top1()
				var istrue bool
				switch this.Code[0].Op {
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
				this.StackInts = this.StackInts.Dropped(2)
				this.StackStrs.Push(result)
			} else {
				node1, node2 := this.Heap[this.StackA.Top0()].(nodeInt), this.Heap[this.StackA.Top1()].(nodeInt)
				var istrue bool
				switch this.Code[0].Op {
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
				addr := this.Heap.Alloc(result)
				this.StackA = this.StackA.Dropped(1)
				this.StackA[this.StackA.Pos0()] = addr
			}
		case INSTR_PRIM_AR_ADD, INSTR_PRIM_AR_SUB, INSTR_PRIM_AR_MUL, INSTR_PRIM_AR_DIV:
			if _MARK7 {
				num1, num2 := this.StackInts.Top0(), this.StackInts.Top1()
				var result int
				switch this.Code[0].Op {
				case INSTR_PRIM_AR_ADD:
					result = num1 + num2
				case INSTR_PRIM_AR_SUB:
					result = num1 - num2
				case INSTR_PRIM_AR_MUL:
					result = num1 * num2
				case INSTR_PRIM_AR_DIV:
					result = num1 / num2
				}
				this.StackInts = this.StackInts.Dropped(1)
				this.StackInts[this.StackInts.Pos0()] = result
			} else {
				node1, node2 := this.Heap[this.StackA.Top0()].(nodeInt), this.Heap[this.StackA.Top1()].(nodeInt)
				var result nodeInt
				switch this.Code[0].Op {
				case INSTR_PRIM_AR_ADD:
					result = node1 + node2
				case INSTR_PRIM_AR_SUB:
					result = node1 - node2
				case INSTR_PRIM_AR_MUL:
					result = node1 * node2
				case INSTR_PRIM_AR_DIV:
					result = node1 / node2
				}
				addr := this.Heap.Alloc(result)
				this.StackA = this.StackA.Dropped(1)
				this.StackA[this.StackA.Pos0()] = addr
			}
		case INSTR_PRIM_AR_NEG:
			if _MARK7 {
				this.StackInts[this.StackInts.Pos0()] = -this.StackInts[this.StackInts.Pos0()]
			} else {
				node := this.Heap[this.StackA.Top0()].(nodeInt)
				addr := this.Heap.Alloc(-node)
				this.StackA[this.StackA.Pos0()] = addr
			}
		case INSTR_PRIM_COND:
			if _MARK7 {
				ctortag := this.StackStrs.Top0()
				this.StackStrs = this.StackStrs.Dropped(1)
				if ctortag == "True" {
					next = append(this.Code[0].CondThen, next...)
				} else if ctortag == "False" {
					next = append(this.Code[0].CondElse, next...)
				} else {
					panic(ctortag)
				}
			} else {
				if node := this.Heap[this.StackA.Top0()].(nodeCtor); node.Tag == "True" {
					next = append(this.Code[0].CondThen, next...)
				} else if node.Tag == "False" {
					next = append(this.Code[0].CondElse, next...)
				} else {
					panic(node.Tag)
				}
				this.StackA = this.StackA.Dropped(1)
			}
		case INSTR_CTOR_PACK:
			arity := this.Code[0].CtorArity
			node := nodeCtor{Tag: this.Code[0].Name, Items: make([]clutil.Addr, arity)}
			for i := 0; i < arity; i++ {
				node.Items[i] = this.StackA.Top(i)
			}
			this.StackA = this.StackA.Dropped(arity).Pushed(this.Heap.Alloc(node))
		case INSTR_CASE_JUMP:
			node := this.Heap[this.StackA.Top0()].(nodeCtor)
			if code := this.Code[0].CaseJump[node.Tag]; len(code) > 0 {
				next = append(code, next...)
			} else if code = this.Code[0].CaseJump["_"]; len(code) > 0 { // jump to default case
				next = append(code, next...)
			} else {
				panic("no matching alternative in CASE OF for ‹" + node.Tag + "," + strconv.Itoa(len(node.Items)) + "› and no default (tag 0) alternative either")
			}
		case INSTR_CASE_SPLIT:
			node := this.Heap[this.StackA.Top0()].(nodeCtor)
			this.StackA = this.StackA.Dropped(1)
			for i := /*len(node.Items)*/ this.Code[0].Int - 1; i > -1; i-- {
				this.StackA.Push(node.Items[i])
			}
		case INSTR_MARK7_PUSHINTVAL:
			this.StackInts.Push(this.Code[0].Int)
		case INSTR_MARK7_MAKENODEBOOL:
			this.StackA.Push(this.Heap.Alloc(nodeCtor{Tag: this.StackStrs.Top0()}))
			this.StackStrs = this.StackStrs.Dropped(1)
		case INSTR_MARK7_MAKENODEINT:
			this.StackA.Push(this.Heap.Alloc(nodeInt(this.StackInts.Top0())))
			this.StackInts = this.StackInts.Dropped(1)
		case INSTR_MARK7_PUSHNODEINT:
			addr := this.StackA.Top0()
			this.StackA = this.StackA.Dropped(1)
			switch node := this.Heap[addr].(type) {
			case nodeCtor:
				this.StackStrs.Push(node.Tag)
			case nodeInt:
				this.StackInts.Push(int(node))
			}
		default:
			panic(this.Code[0].Op)
		}

		if this.Code = next; this.Stats.MaxStack < len(this.StackA) {
			this.Stats.MaxStack = len(this.StackA)
		}
		if this.Stats.NumSteps > 9999999 {
			panic("exceeded 10 million steps: probable infinite loop, stopping evaluation")
		}
	}
	this.Stats.HeapSize = len(this.Heap)
}

func (this *gMachine) String(result interface{}) string {
	switch res := result.(type) {
	case nodeInt:
		return "#" + strconv.Itoa(int(res))
	case nodeCtor:
		s := "‹" + res.Tag
		for _, addr := range res.Items {
			s += " " + this.String(this.Heap[addr])
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
