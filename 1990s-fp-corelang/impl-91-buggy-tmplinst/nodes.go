package climpl

import (
	"github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	"github.com/metaleap/go-machines/1990s-fp-corelang/util"
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

func (this *tiMachine) instantiate(expression clsyn.IExpr) (resultAddr clutil.Addr) {
	switch expr := expression.(type) {
	case *clsyn.ExprLitFloat:
		resultAddr = this.Heap.Alloc(nodeNumFloat(expr.Lit))
	case *clsyn.ExprLitUInt:
		resultAddr = this.Heap.Alloc(nodeNumUint(expr.Lit))
	case *clsyn.ExprCall:
		resultAddr = this.Heap.Alloc(&nodeAp{this.instantiate(expr.Callee), this.instantiate(expr.Arg)})
	case *clsyn.ExprIdent:
		resultAddr = this.Env.LookupOrPanic(expr.Name)
	case *clsyn.ExprLetIn:
		for _, def := range expr.Defs {
			ndef := nodeDef(*def)
			this.Env[def.Name] = this.Heap.Alloc(&ndef)
		}
		resultAddr = this.instantiate(expr.Body)
	case *clsyn.ExprCaseOf, *clsyn.ExprCtor:
		panic("instantiate: expr type coming soon")
	default:
		panic("instantiate: expr type not yet implemented")
	}
	return
}
