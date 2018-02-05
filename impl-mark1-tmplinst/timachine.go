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
		initialMachineState.Heap, initialMachineState.Globals[def.Name] = initialMachineState.Heap.Alloc(nodeDef(*def))
	}
	initialMachineState.Stack = []clutil.Addr{initialMachineState.Globals[mainName]}
	return
}

func (me *TiState) Eval() (allSteps []*TiState) {
	if allSteps = []*TiState{me}; !me.isFinalState() {
		allSteps = append(allSteps, me.step().stats().Eval()...)
	}
	return
}

func (me *TiState) stats() *TiState {
	me.Stats.NumberOfStepsTaken++
	return me
}

func (me *TiState) step() (nu *TiState) {
	return
}

func (me *TiState) isFinalState() bool {
	if len(me.Stack) == 0 {
		panic("empty stack")
	}
	return len(me.Stack) == 1 && me.Heap[me.Stack[0]].IsValue()
}
