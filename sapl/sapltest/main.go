package main

import (
	"bufio"
	"bytes"
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

	outexpr, outlist := prog.Eval(ctx, ExprAppl{Callee: ExprAppl{
		Callee: ExprFnRef(len(prog) - 1), Arg: ListsFrom(os.Args[2:]),
	}, Arg: ListsFrom(os.Environ())})

	println(ctx.String())
	if outbytes := prog.ToBytes(outlist); outbytes != nil {
		os.Stdout.Write(append(outbytes, 10))
	} else if outlist == nil {
		os.Stdout.WriteString("EXPR:\t" + outexpr.String() + "\n")
	} else if !maybeRepl(prog, ctx, outlist) {
		os.Stdout.WriteString("[ ")
		for _, expr := range outlist {
			os.Stdout.WriteString(expr.String() + " , ")
		}
		os.Stdout.WriteString("]\n")
	}
	if tracing {
		println(trace.root.subSteps[0].str(0))
	}
}

func maybeRepl(prog Prog, ctx *CtxEval, outList []Expr) bool {
	if len(outList) == 3 {
		if fnhandler, okf := outList[0].(ExprFnRef); okf {
			if sepchar, oks := outList[1].(ExprNum); oks {
				if appl, oka := outList[2].(ExprAppl); oka {
					if strintro := prog.ToBytes(prog.List(ctx, appl)); strintro != nil {
						handleinput := func(input []byte, state Expr) Expr {
							stateless := (state == nil)
							if stateless {
								state = ExprFnRef(3)
							}
							if retexpr, retlist := prog.Eval(ctx, ExprAppl{Callee: ExprAppl{Callee: fnhandler, Arg: state}, Arg: ListFrom(input)}); retlist == nil {
								panic(retexpr.String())
							} else if stateless {
								os.Stdout.Write(prog.ToBytes(retlist))
							} else if len(retlist) == 2 {
								state = retlist[0]
								os.Stdout.Write(prog.ToBytes(prog.List(ctx, retlist[1])))
							} else {
								panic(retexpr.String())
							}
							return state
						}
						if os.Stdout.Write(strintro); sepchar == 0 {
							if allinputatonce, err := ioutil.ReadAll(os.Stdin); err != nil {
								panic(err)
							} else {
								handleinput(allinputatonce, nil)
							}
						} else {
							reader := bufio.NewScanner(os.Stdin)
							if sepchar != '\n' {
								reader.Split(readSplitter(byte(sepchar)))
							}
							for state := Expr(ExprFnRef(3)); reader.Scan(); {
								state = handleinput(reader.Bytes(), state)
							}
							if err := reader.Err(); err != nil {
								panic(err)
							}
						}
						return true
					}
				}
			}
		}
	}
	return false
}

func readSplitter(sep byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.IndexByte(data, sep); i >= 0 {
			advance, token = i+1, data[0:i]
		} else if atEOF {
			advance, token = len(data), data
		}
		return
	}
}
