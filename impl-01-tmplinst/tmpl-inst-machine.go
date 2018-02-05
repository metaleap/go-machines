package climpl

import (
	"errors"
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

func CompileToMachine(mod *clsyn.SynMod) (initialMachineState *TiState) {
	initialMachineState = &TiState{
		Globals: map[string]clutil.Addr{},
		Heap:    clutil.Heap{},
	}
	for _, def := range mod.Defs {
		ndef := nodeDef(*def)
		initialMachineState.Heap, initialMachineState.Globals[def.Name] = initialMachineState.Heap.Alloc(&ndef)
	}
	return
}

func (me *TiState) Eval(name string) (allSteps []*TiState, err error) {
	defer clutil.Catch(&err)
	addr := me.Globals[name]
	if addr == 0 {
		return nil, errors.New("undefined: " + name)
	}
	me.Stack = []clutil.Addr{addr}
	allSteps = me.eval()
	return
}

func (me *TiState) eval() (allSteps []*TiState) {
	if allSteps = []*TiState{me}; !me.isFinalState() {
		allSteps = append(allSteps, me.step().eval()...)
	}
	return
}

func (me *TiState) isFinalState() bool {
	if len(me.Stack) == 0 {
		panic("isFinalState: empty stack")
	}
	return len(me.Stack) == 1 && isDataNode(me.Heap[me.Stack[0]])
}

func (me *TiState) step() *TiState {
	headaddr := me.Stack[0]
	nu, obj := *me, me.Heap[headaddr]
	switch node := obj.(type) {
	case nodeNumFloat:
		panic("float applied as a function")
	case nodeNumUint:
		panic("uint applied as a function")
	case *nodeAp:
		nu.Stack = append([]clutil.Addr{node.Callee}, me.Stack...)
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

		nuheap, resultaddr := instantiateNodeFromExpr(node.Body, me.Heap, env)
		nu.Heap, nu.Stack = nuheap, append([]clutil.Addr{resultaddr}, me.Stack[1+len(node.Args):]...)
	default:
		panic("step: node type not yet implemented")
	}
	nu.Stats.NumberOfStepsTaken = me.Stats.NumberOfStepsTaken + 1
	return &nu
}

func (me *TiState) getArgs(num int) (argsaddrs []clutil.Addr) {
	if num >= len(me.Stack) {
		panic("not enough arguments given")
	}
	for i := 1; i <= num; i++ {
		addr := me.Stack[i]
		nap, _ := me.Heap[addr].(*nodeAp)
		argsaddrs = append(argsaddrs, nap.Arg)
	}
	return
}
