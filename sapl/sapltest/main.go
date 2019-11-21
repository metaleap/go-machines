package main

import (
	"io/ioutil"
	"os"

	. "github.com/metaleap/go-machines/sapl"
)

const tracing = true

func main() {
	src, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	prog, trace := LoadFromJson(src), &tracer{root: &traceStep{}}
	tracestep := trace.onEvalStep
	if !tracing {
		tracestep = nil
	}
	result, timetaken := prog.Eval(prog[len(prog)-1].Expr, tracestep)
	println(timetaken.String(), result.String())
}

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
	oldstep, curstep := me.root, &traceStep{expr, stack, nil, nil}
	me.root, oldstep.subSteps = curstep, append(oldstep.subSteps, curstep)
	return func(expr Expr) Expr {
		me.root = oldstep
		return expr
	}
}

func (me *traceStep) String() string {
	return "foo"
}
