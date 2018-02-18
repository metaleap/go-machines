package climpl

import (
	"github.com/metaleap/go-machines/1991-fpcorelang/syn"
	"github.com/metaleap/go-machines/1991-fpcorelang/util"
)

type nodeAp struct {
	Callee clutil.Addr
	Arg    clutil.Addr
}

type nodeDef clsyn.SynDef

type nodeNumFloat float64

type nodeNumUint uint64

func isDataNode(node clutil.INode) (isvalue bool) {
	switch node.(type) {
	case nodeNumFloat, nodeNumUint:
		isvalue = true
	}
	return
}

func (me *tiMachine) instantiate(expression clsyn.IExpr) (resultAddr clutil.Addr) {
	switch expr := expression.(type) {
	case *clsyn.ExprLitFloat:
		resultAddr = me.Heap.Alloc(nodeNumFloat(expr.Lit))
	case *clsyn.ExprLitUInt:
		resultAddr = me.Heap.Alloc(nodeNumUint(expr.Lit))
	case *clsyn.ExprCall:
		resultAddr = me.Heap.Alloc(&nodeAp{me.instantiate(expr.Callee), me.instantiate(expr.Arg)})
	case *clsyn.ExprIdent:
		resultAddr = me.Env.LookupOrPanic(expr.Name)
	case *clsyn.ExprLetIn:
		for _, def := range expr.Defs {
			ndef := nodeDef(*def)
			me.Env[def.Name] = me.Heap.Alloc(&ndef)
		}
		resultAddr = me.instantiate(expr.Body)
	case *clsyn.ExprCaseOf, *clsyn.ExprCtor:
		panic("instantiate: expr type coming soon")
	default:
		panic("instantiate: expr type not yet implemented")
	}
	return
}
