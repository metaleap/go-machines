package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type TiState struct {
	Heap    clutil.Heap
	Stack   []clutil.Addr
	Globals map[string]clutil.Addr
	Dump    struct{}
	Stats   struct {
		NumberOfStepsTaken int
	}
}

func CompileToMachine(mod *clsyn.SynMod, mainName string) (initialMachineState *TiState) {
	initialMachineState = &TiState{
		Globals: map[string]clutil.Addr{},
		Heap:    clutil.Heap{},
	}
	for _, def := range mod.Defs {
		ndef := nodeDef(*def)
		initialMachineState.Heap, initialMachineState.Globals[def.Name] = initialMachineState.Heap.Alloc(&ndef)
	}
	initialMachineState.Stack = []clutil.Addr{initialMachineState.Globals[mainName]}
	return
}

func (me *TiState) Eval() (allSteps []*TiState, err error) {
	defer clutil.Catch(&err)
	allSteps = me.eval()
	return
}

func (me *TiState) eval() (allSteps []*TiState) {
	if allSteps = []*TiState{me}; !me.isFinalState() {
		allSteps = append(allSteps, me.step().stats().eval()...)
	}
	return
}

func (me *TiState) isFinalState() bool {
	if len(me.Stack) == 0 {
		panic("isFinalState: empty stack")
	}
	return len(me.Stack) == 1 && me.Heap[me.Stack[0]].IsValue()
}

func (me *TiState) stats() *TiState {
	me.Stats.NumberOfStepsTaken++
	return me
}

func (me *TiState) step() *TiState {
	headaddr := me.Stack[len(me.Stack)-1]
	nu, obj := *me, me.Heap[headaddr]
	switch node := obj.(type) {
	case *nodeNumFloat:
		panic("float applied as a function")
	case *nodeNumUint:
		panic("uint applied as a function")
	case *nodeAp:
		nu.Stack = append(nu.Stack, node.Callee)
	case *nodeDef:
		argaddrs, argbinds := me.getArgs(len(node.Args)), make(map[string]clutil.Addr, len(node.Args))
		for i, argname := range node.Args {
			argbinds[argname] = argaddrs[i]
		}

		env := make(map[string]clutil.Addr, len(argbinds)+len(me.Globals))
		for k, v := range me.Globals {
			env[k] = v
		}
		for k, v := range argbinds {
			env[k] = v
		}

		nuheap, resultaddr := me.instantiate(node.Body, env)
		pos := len(me.Stack) - (1 + len(node.Args))
		nu.Heap, nu.Stack = nuheap, append(me.Stack[:pos], resultaddr)
	default:
		panic("step: node type not yet implemented")
	}
	return &nu
}

func (me *TiState) getArgs(num int) (argsaddrs []clutil.Addr) {
	for _, addr := range me.Stack {
		if nap, _ := me.Heap[addr].(*nodeAp); nap != nil {
			if argsaddrs = append(argsaddrs, nap.Arg); len(argsaddrs) == num {
				return
			}
		}
	}
	return
}

func (me *TiState) instantiate(body clsyn.IExpr, env map[string]clutil.Addr) (nuHeap clutil.Heap, resultAddr clutil.Addr) {
	switch expr := body.(type) {
	case *clsyn.ExprLitFloat:
		nuHeap, resultAddr = me.Heap.Alloc(nodeNumFloat(expr.Lit))
	case *clsyn.ExprLitUInt:
		nuHeap, resultAddr = me.Heap.Alloc(nodeNumUint(expr.Lit))
	default:
		panic("instantiate: expr type not yet implemented")
	}
	return
}
