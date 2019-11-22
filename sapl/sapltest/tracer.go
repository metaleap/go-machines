package main

import (
	"os"
	"strconv"
	"strings"

	. "github.com/metaleap/go-machines/sapl"
)

type tracer struct {
	root *traceStep
}

type traceStep struct {
	expr     Expr
	stack    []Expr
	subSteps []*traceStep
	ret      Expr
}

func (me *tracer) onEvalStep(expr Expr, stack []Expr) (end func(Expr) Expr) {
	parent, curstep := me.root, &traceStep{expr, make([]Expr, len(stack)), nil, nil}
	copy(curstep.stack, stack)
	me.root, parent.subSteps = curstep, append(parent.subSteps, curstep)
	os.Stdin.Read([]byte{10})
	return func(expr Expr) Expr {
		me.root, curstep.ret = parent, expr
		return expr
	}
}

func (me *traceStep) str(level int) string {
	ret := strings.Repeat("  ", level) + me.expr.String() + "\t=>\t" + me.ret.String() + "\t["
	ret += strconv.Itoa(len(me.stack))
	// for _, expr := range me.stack {
	// 	ret += expr.String() + "] ["
	// }
	ret += "]\n"
	for _, sub := range me.subSteps {
		ret += sub.str(level + 1)
	}
	return ret
}
