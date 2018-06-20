package climpl

import (
	"fmt"

	"github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	"github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

type tiMachine struct {
	Heap  clutil.HeapM
	Stack clutil.StackA
	Env   clutil.Env
	Stats clutil.Stats
}

func CompileToMachine(mod *clsyn.SynMod) (clutil.IMachine, []error) {
	me := &tiMachine{
		Env:  make(clutil.Env, len(mod.Defs)),
		Heap: clutil.HeapM{},
	}
	for _, def := range mod.Defs {
		addr, ndef := me.Heap.NextAddr(), nodeDef(*def)
		me.Env[def.Name], me.Heap[addr] = addr, &ndef
	}
	return me, nil
}

func (this *tiMachine) String(result interface{}) string { return fmt.Sprintf("%v", result) }

func (this *tiMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	addr := this.Env[name]
	if addr == 0 {
		panic("undefined: " + name)
	} else {
		this.Stack = clutil.StackA{addr}
		this.eval()
		val, stats = this.Heap[this.Stack[0]], this.Stats
	}
	return
}

func (this *tiMachine) eval() {
	for this.Stats.NumSteps, this.Stats.NumAppls = 0, 0; !this.isFinalState(); this.step() {
	}
}

func (this *tiMachine) isFinalState() bool {
	return len(this.Stack) == 1 && isDataNode(this.Heap[this.Stack[0]])
}

func (this *tiMachine) step() {
	if this.Stats.NumSteps++; this.Stats.NumSteps > 9999 {
		panic("infinite loop")
	}
	addr := this.Stack.Top(0)
	obj := this.Heap[addr]
	switch n := obj.(type) {
	case nodeNumFloat, nodeNumUint:
		panic("number applied as a function")
	case *nodeAp:
		this.Stats.NumAppls++
		this.Stack.Push(n.Callee)
	case *nodeDef:
		oldenv := this.Env
		this.Env = make(map[string]clutil.Addr, len(n.Args)+len(oldenv))
		for k, v := range oldenv {
			this.Env[k] = v
		}

		argsaddrs := this.getArgs(n.Name, len(n.Args))
		for i, argname := range n.Args {
			this.Env[argname] = argsaddrs[i]
		}

		resultaddr := this.instantiate(n.Body)
		this.Stack = this.Stack.Dropped(1 + len(n.Args)).Pushed(resultaddr)

		// me.Env = oldenv
	}
}

func (this *tiMachine) getArgs(name string, count int) (argsaddrs []clutil.Addr) {
	pos := this.Stack.Pos(count)
	if pos < 0 {
		panic(name + ": not enough arguments given")
	}
	argsaddrs = make([]clutil.Addr, count)
	for i := 0; i < count; i++ {
		addr := this.Stack[pos+i]
		nap, _ := this.Heap[addr].(*nodeAp)
		argsaddrs[count-1-i] = nap.Arg
	}
	return
}
