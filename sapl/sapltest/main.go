package main

import (
	"io/ioutil"
	"os"

	. "github.com/metaleap/go-machines/sapl"
)

const tracing = false

func main() {
	src, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	prog, trace := LoadFromJson(src), &tracer{root: &traceStep{}}
	tracestep := trace.onEvalStep
	if !tracing {
		tracestep = nil
	}
	result, stats := prog.Eval(ExprAppl{Callee: ExprFnRef(len(prog) - 1), Arg: ExprNum(7)}, tracestep)
	println(result.String())
	println(stats.String())
	if tracing {
		println(trace.root.subSteps[0].str(0))
	}
}
