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

func instantiateNodeFromExpr(body clsyn.IExpr, heap clutil.Heap, env map[string]clutil.Addr) (nuHeap clutil.Heap, resultAddr clutil.Addr) {
	switch expr := body.(type) {
	case *clsyn.ExprLitFloat:
		nuHeap, resultAddr = heap.Alloc(nodeNumFloat(expr.Lit))
	case *clsyn.ExprLitUInt:
		nuHeap, resultAddr = heap.Alloc(nodeNumUint(expr.Lit))
	case *clsyn.ExprCall:
		heap1, a1 := instantiateNodeFromExpr(expr.Callee, heap, env)
		heap2, a2 := instantiateNodeFromExpr(expr.Arg, heap1, env)
		nuHeap, resultAddr = heap2.Alloc(&nodeAp{Callee: a1, Arg: a2})
	case *clsyn.ExprIdent:
		if nuHeap, resultAddr = heap, env[expr.Name]; resultAddr == 0 {
			panic(expr.Name + ": undefined")
		}
	case *clsyn.ExprCaseOf, *clsyn.ExprLetIn, *clsyn.ExprCtor:
		panic("instantiateNodeFromExpr: expr type coming soon")
	default:
		panic("instantiateNodeFromExpr: expr type not yet implemented")
	}
	return
}
