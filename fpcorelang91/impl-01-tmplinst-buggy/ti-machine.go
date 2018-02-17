package climpl

import (
	"fmt"

	"github.com/metaleap/go-machines/fpcorelang91/syn"
	"github.com/metaleap/go-machines/fpcorelang91/util"
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

func (me *tiMachine) String(result interface{}) string { return fmt.Sprintf("%v", result) }

func (me *tiMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	defer clutil.Catch(&err)
	addr := me.Env[name]
	if addr == 0 {
		panic("undefined: " + name)
	} else {
		me.Stack = clutil.StackA{addr}
		me.eval()
		val, stats = me.Heap[me.Stack[0]], me.Stats
	}
	return
}

func (me *tiMachine) eval() {
	for me.Stats.NumSteps, me.Stats.NumAppls = 0, 0; !me.isFinalState(); me.step() {
	}
}

func (me *tiMachine) isFinalState() bool {
	return len(me.Stack) == 1 && isDataNode(me.Heap[me.Stack[0]])
}

func (me *tiMachine) step() {
	if me.Stats.NumSteps++; me.Stats.NumSteps > 9999 {
		panic("infinite loop")
	}
	addr := me.Stack.Top(0)
	obj := me.Heap[addr]
	switch n := obj.(type) {
	case nodeNumFloat, nodeNumUint:
		panic("number applied as a function")
	case *nodeAp:
		me.Stats.NumAppls++
		me.Stack.Push(n.Callee)
	case *nodeDef:
		oldenv := me.Env
		me.Env = make(map[string]clutil.Addr, len(n.Args)+len(oldenv))
		for k, v := range oldenv {
			me.Env[k] = v
		}

		argsaddrs := me.getArgs(n.Name, len(n.Args))
		for i, argname := range n.Args {
			me.Env[argname] = argsaddrs[i]
		}

		resultaddr := me.instantiate(n.Body)
		me.Stack = me.Stack.Dropped(1 + len(n.Args)).Pushed(resultaddr)

		// me.Env = oldenv
	}
}

func (me *tiMachine) getArgs(name string, count int) (argsaddrs []clutil.Addr) {
	pos := me.Stack.Pos(count)
	if pos < 0 {
		panic(name + ": not enough arguments given")
	}
	argsaddrs = make([]clutil.Addr, count)
	for i := 0; i < count; i++ {
		addr := me.Stack[pos+i]
		nap, _ := me.Heap[addr].(*nodeAp)
		argsaddrs[count-1-i] = nap.Arg
	}
	return
}
