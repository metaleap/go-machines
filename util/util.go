package clutil

import (
	"strconv"
)

type INode interface {
}

type Addr int

func (me Addr) String() string { return strconv.Itoa(int(me)) }

type Heap map[Addr]INode
