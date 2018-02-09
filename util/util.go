package clutil

import (
	"strconv"
)

type IMachine interface {
	Eval(string) (interface{}, Stats, error)
}

type INode interface {
}

type Addr int

func (me Addr) String() string { return strconv.Itoa(int(me)) }

type Heap map[Addr]INode

func (me Heap) Alloc(obj INode) (addr Addr) {
	addr = me.NextAddr()
	me[addr] = obj
	return
}

func (me Heap) NextAddr() Addr {
	return Addr(1 + len(me))
}

type Env map[string]Addr

func (me Env) LookupOrPanic(name string) (addr Addr) {
	if addr = me[name]; addr == 0 {
		panic("undefined: " + name)
	}
	return
}

type Stack []Addr

func (me Stack) Dropped(n int) Stack {
	return me[:len(me)-n]
}

func (me Stack) Pos(i int) int {
	return len(me) - (1 + i)
}

func (me *Stack) Push(addr Addr) {
	*me = append(*me, addr)
}

func (me Stack) Pushed(addr Addr) Stack {
	return append(me, addr)
}

func (me Stack) Top(i int) Addr {
	return me[len(me)-(1+i)]
}
