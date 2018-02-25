package clutil

import (
	"strconv"
)

type IMachine interface {
	Eval(argLessDefName string) (result interface{}, stats Stats, err error)
	String(result interface{}) string
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

func (me StackA) Pos0() int {
	return len(me) - 1
}

func (me StackA) Pos1() int {
	return len(me) - 2
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

func (me StackA) Top0() Addr {
	return me[len(me)-1]
}

func (me StackA) Top1() Addr {
	return me[len(me)-2]
}

func (me StackA) Top(i int) Addr {
	return me[len(me)-(1+i)]
}

type StackI []int

func (me StackI) Dropped(n int) StackI {
	return me[:len(me)-n]
}

func (me StackI) Pos0() int {
	return len(me) - 1
}

func (me StackI) Pos1() int {
	return len(me) - 2
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

func (me StackI) Top0() int {
	return me[len(me)-1]
}

func (me StackI) Top1() int {
	return me[len(me)-2]
}

func (me StackI) Top(i int) int {
	return me[len(me)-(1+i)]
}

type StackS []string

func (me StackS) Dropped(n int) StackS {
	return me[:len(me)-n]
}

func (me StackS) Pos0() int {
	return len(me) - 1
}

func (me StackS) Pos1() int {
	return len(me) - 2
}

func (me StackS) Pos(i int) int {
	return len(me) - (1 + i)
}

func (me *StackS) Push(i string) {
	*me = append(*me, i)
}

func (me StackS) Pushed(i string) StackS {
	return append(me, i)
}

func (me StackS) Top0() string {
	return me[len(me)-1]
}

func (me StackS) Top1() string {
	return me[len(me)-2]
}

func (me StackS) Top(i int) string {
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
