package clutil

import (
	"strconv"
)

type IMachine interface {
	Eval(string) (interface{}, int, int, error)
}

type INode interface {
}

type Addr int

func (me Addr) String() string { return strconv.Itoa(int(me)) }

type Heap map[Addr]INode
