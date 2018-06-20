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

func (this Addr) String() string { return strconv.Itoa(int(this)) }

type Env map[string]Addr

func (this Env) LookupOrPanic(name string) (addr Addr) {
	if addr = this[name]; addr == 0 {
		panic("undefined: " + name)
	}
	return
}

type StackA []Addr

func (this StackA) Dropped(n int) StackA {
	return this[:len(this)-n]
}

func (this StackA) Pos0() int {
	return len(this) - 1
}

func (this StackA) Pos1() int {
	return len(this) - 2
}

func (this StackA) Pos(i int) int {
	return len(this) - (1 + i)
}

func (this *StackA) Push(addr Addr) {
	*this = append(*this, addr)
}

func (this StackA) Pushed(addr Addr) StackA {
	return append(this, addr)
}

func (this StackA) Top0() Addr {
	return this[len(this)-1]
}

func (this StackA) Top1() Addr {
	return this[len(this)-2]
}

func (this StackA) Top(i int) Addr {
	return this[len(this)-(1+i)]
}

type StackI []int

func (this StackI) Dropped(n int) StackI {
	return this[:len(this)-n]
}

func (this StackI) Pos0() int {
	return len(this) - 1
}

func (this StackI) Pos1() int {
	return len(this) - 2
}

func (this StackI) Pos(i int) int {
	return len(this) - (1 + i)
}

func (this *StackI) Push(i int) {
	*this = append(*this, i)
}

func (this StackI) Pushed(i int) StackI {
	return append(this, i)
}

func (this StackI) Top0() int {
	return this[len(this)-1]
}

func (this StackI) Top1() int {
	return this[len(this)-2]
}

func (this StackI) Top(i int) int {
	return this[len(this)-(1+i)]
}

type StackS []string

func (this StackS) Dropped(n int) StackS {
	return this[:len(this)-n]
}

func (this StackS) Pos0() int {
	return len(this) - 1
}

func (this StackS) Pos1() int {
	return len(this) - 2
}

func (this StackS) Pos(i int) int {
	return len(this) - (1 + i)
}

func (this *StackS) Push(i string) {
	*this = append(*this, i)
}

func (this StackS) Pushed(i string) StackS {
	return append(this, i)
}

func (this StackS) Top0() string {
	return this[len(this)-1]
}

func (this StackS) Top1() string {
	return this[len(this)-2]
}

func (this StackS) Top(i int) string {
	return this[len(this)-(1+i)]
}

type HeapM map[Addr]INode

func (this HeapM) Alloc(obj INode) (addr Addr) {
	addr = this.NextAddr()
	this[addr] = obj
	return
}

func (this HeapM) NextAddr() Addr {
	return Addr(1 + len(this))
}

type HeapA []INode

func (this *HeapA) Alloc(obj INode) (addr Addr) {
	addr = this.NextAddr()
	*this = append(*this, obj)
	return
}

func (this HeapA) NextAddr() Addr {
	return Addr(len(this))
}
