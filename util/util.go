package clutil

import (
	"strconv"
)

type INode interface {
}

type Addr int

func (me Addr) String() string { return strconv.Itoa(int(me)) }

type Heap map[Addr]INode

func (me Heap) copy() (nu Heap) {
	nu = make(Heap, len(me))
	for k, v := range me {
		nu[k] = v
	}
	return
}

func (me Heap) Alloc(obj INode) (nu Heap, addr Addr) {
	nu = me.copy()
	addr = nu.NextAddr()
	nu[addr] = obj
	return
}

func (me Heap) Free(addr Addr) (nu Heap) {
	nu = me.copy()
	delete(nu, addr)
	return
}

func (me Heap) NextAddr() Addr {
	return Addr(len(me) + 1)
}

func (me Heap) Update(addr Addr, obj INode) (nu Heap) {
	nu = me.copy()
	nu[addr] = obj
	return
}
