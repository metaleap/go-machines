package main

import (
	"io"
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
	ctx := &CtxEval{Tracer: trace.onEvalStep, Outputs: []io.Writer{os.Stdout, os.Stderr}}
	if !tracing {
		ctx.Tracer = nil
	}
	_ = prog.Eval(ctx, ExprAppl{Callee: ExprFnRef(len(prog) - 1), Arg: ExprNum(88)})
	// println(result.String())
	println(ctx.String())
	if tracing {
		println(trace.root.subSteps[0].str(0))
	}
}
