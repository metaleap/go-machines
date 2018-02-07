package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type TiMachine struct {
	Heap  clutil.Heap
	Stack []clutil.Addr
	Env   map[string]clutil.Addr
	Stats struct {
		NumStepsTaken int
	}
}

func CompileToMachine(mod *clsyn.SynMod) (me *TiMachine) {
	me = &TiMachine{
		Env:  make(map[string]clutil.Addr, len(mod.Defs)),
		Heap: clutil.Heap{},
	}
	for _, def := range mod.Defs {
		addr, ndef := me.nextAddr(), nodeDef(*def)
		me.Env[def.Name], me.Heap[addr] = addr, &ndef
	}
	return
}

func (me *TiMachine) Eval(name string) (val interface{}, numsteps int, err error) {
	defer clutil.Catch(&err)
	addr := me.Env[name]
	if me.Stats.NumStepsTaken = 0; addr == 0 {
		panic("undefined: " + name)
	} else {
		me.Stack = []clutil.Addr{addr}
		me.eval()
		val, numsteps = me.Heap[me.Stack[0]], me.Stats.NumStepsTaken
	}
	return
}

func (me *TiMachine) eval() {
	for !me.isFinalState() {
		me.step()
	}
}

func (me *TiMachine) isFinalState() bool {
	if len(me.Stack) == 0 {
		panic("isFinalState: empty stack")
	}
	return len(me.Stack) == 1 && isDataNode(me.Heap[me.Stack[0]])
}

func (me *TiMachine) step() {
	if me.Stats.NumStepsTaken++; me.Stats.NumStepsTaken > 9999 {
		panic("infinite loop")
	}
	addr := me.Stack[len(me.Stack)-1]
	obj := me.Heap[addr]
	switch n := obj.(type) {
	case nodeNumFloat, nodeNumUint:
		panic("number applied as a function")
	case *nodeAp:
		me.Stack = append(me.Stack, n.Callee)
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
		me.Stack = append(me.Stack[:len(me.Stack)-(1+len(n.Args))], resultaddr)

		// me.Env = oldenv
	}
}

func (me *TiMachine) alloc(obj clutil.INode) (addr clutil.Addr) {
	addr = me.nextAddr()
	me.Heap[addr] = obj
	return
}

func (me *TiMachine) getArgs(name string, count int) (argsaddrs []clutil.Addr) {
	stackzero := len(me.Stack) - (1 + count)
	if stackzero < 0 {
		panic(name + ": not enough arguments given")
	}
	argsaddrs = make([]clutil.Addr, count)
	for i := 0; i < count; i++ {
		addr := me.Stack[stackzero+i]
		nap, _ := me.Heap[addr].(*nodeAp)
		argsaddrs[count-1-i] = nap.Arg
	}
	return
}

func (me *TiMachine) nextAddr() clutil.Addr {
	return clutil.Addr(len(me.Heap) + 1)
}
