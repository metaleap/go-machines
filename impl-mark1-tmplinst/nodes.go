package climpl

import (
	"github.com/metaleap/go-corelang/syn"
)

type nodeAp clsyn.ExprCall

func (nodeAp) IsValue() bool { return false }

type nodeDef clsyn.SynDef

func (nodeDef) IsValue() bool { return false }

type nodeNumFloat clsyn.ExprLitFloat

func (nodeNumFloat) IsValue() bool { return true }

type nodeNumUint clsyn.ExprLitUInt

func (nodeNumUint) IsValue() bool { return true }
