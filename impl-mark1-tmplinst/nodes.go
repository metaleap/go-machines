package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type nodeAp struct {
	Callee clutil.Addr
	Arg    clutil.Addr
}

func (*nodeAp) IsValue() bool { return false }

type nodeDef clsyn.SynDef

func (*nodeDef) IsValue() bool { return false }

type nodeNumFloat float64

func (nodeNumFloat) IsValue() bool { return true }

type nodeNumUint uint64

func (nodeNumUint) IsValue() bool { return true }
