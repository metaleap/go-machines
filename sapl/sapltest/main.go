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
	ctx := &CtxEval{Tracer: trace.onEvalStep}
	if !tracing {
		ctx.Tracer = nil
	}

	_, outbytes := prog.Eval(ctx, ExprAppl{Callee: ExprAppl{
		Callee: ExprFnRef(len(prog) - 1), Arg: ListsFrom(os.Args[2:]),
	}, Arg: ListsFrom(os.Environ())})

	os.Stdout.Write(append(outbytes, 10))
	println(ctx.String())
	if tracing {
		println(trace.root.subSteps[0].str(0))
	}
}
