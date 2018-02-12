package clutil

import (
	"strconv"
)

type IMachine interface {
	Eval(string) (interface{}, Stats, error)
	String(interface{}) string
}

type INode interface {
}

type Addr int

func (me Addr) String() string { return strconv.Itoa(int(me)) }

type Env map[string]Addr

func (me Env) LookupOrPanic(name string) (addr Addr) {
	if addr = me[name]; addr == 0 {
		panic("undefined: " + name)
	}
	return
}

type StackA []Addr

func (me StackA) Dropped(n int) StackA {
	return me[:len(me)-n]
}

func (me StackA) Pos(i int) int {
	return len(me) - (1 + i)
}

func (me *StackA) Push(addr Addr) {
	*me = append(*me, addr)
}

func (me StackA) Pushed(addr Addr) StackA {
	return append(me, addr)
}

func (me StackA) Top(i int) Addr {
	return me[len(me)-(1+i)]
}

type StackI []int

func (me StackI) Dropped(n int) StackI {
	return me[:len(me)-n]
}

func (me StackI) Pos(i int) int {
	return len(me) - (1 + i)
}

func (me *StackI) Push(i int) {
	*me = append(*me, i)
}

func (me StackI) Pushed(i int) StackI {
	return append(me, i)
}

func (me StackI) Top(i int) int {
	return me[len(me)-(1+i)]
}

type HeapM map[Addr]INode

func (me HeapM) Alloc(obj INode) (addr Addr) {
	addr = me.NextAddr()
	me[addr] = obj
	return
}

func (me HeapM) NextAddr() Addr {
	return Addr(1 + len(me))
}

type HeapA []INode

func (me *HeapA) Alloc(obj INode) (addr Addr) {
	addr = me.NextAddr()
	*me = append(*me, obj)
	return
}

func (me HeapA) NextAddr() Addr {
	return Addr(len(me))
}
