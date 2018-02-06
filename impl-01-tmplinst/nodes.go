package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
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

func (me *TiMachine) instantiate(expression clsyn.IExpr) (resultAddr clutil.Addr) {
	switch expr := expression.(type) {
	case *clsyn.ExprLitFloat:
		resultAddr = me.alloc(nodeNumFloat(expr.Lit))
	case *clsyn.ExprLitUInt:
		resultAddr = me.alloc(nodeNumUint(expr.Lit))
	case *clsyn.ExprCall:
		resultAddr = me.alloc(&nodeAp{me.instantiate(expr.Callee), me.instantiate(expr.Arg)})
	case *clsyn.ExprIdent:
		if resultAddr = me.Env[expr.Name]; resultAddr == 0 {
			panic(expr.Name + ": undefined")
		}
	case *clsyn.ExprLetIn:
		for _, def := range expr.Defs {
			ndef := nodeDef(*def)
			me.Env[def.Name] = me.alloc(&ndef)
		}
		resultAddr = me.instantiate(expr.Body)
	case *clsyn.ExprCaseOf, *clsyn.ExprCtor:
		panic("instantiate: expr type coming soon")
	default:
		panic("instantiate: expr type not yet implemented")
	}
	return
}
