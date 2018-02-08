package climpl

import (
	"fmt"

	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

const PrintSteps = false

type naiveMachine struct {
	Globals         map[string]clsyn.ISyn
	Locals          map[string]clsyn.ISyn
	Args            []clsyn.ISyn
	NumStepsTaken   int
	NumApplications int
}

func CompileToMachine(mod *clsyn.SynMod) clutil.IMachine {
	globals := make(map[string]clsyn.ISyn, len(mod.Defs))
	for _, def := range mod.Defs {
		globals[def.Name] = def
	}
	return &naiveMachine{Globals: globals}
}

func (me *naiveMachine) Eval(name string) (val interface{}, numAppl int, numSteps int, err error) {
	defer clutil.Catch(&err)
	def := me.resolveIdent(name)
	me.NumStepsTaken, me.NumApplications, me.Locals = 0, 0, map[string]clsyn.ISyn{}
	syn := me.reduce(def)
	switch n := syn.(type) {
	case *clsyn.ExprLitFloat:
		val = n.Lit
	case *clsyn.ExprLitRune:
		val = n.Lit
	case *clsyn.ExprLitText:
		val = n.Lit
	case *clsyn.ExprLitUInt:
		val = n.Lit
	default:
		panic("no atomic result")
	}
	numAppl, numSteps = me.NumApplications, me.NumStepsTaken
	return
}

func (me *naiveMachine) reduce(syn clsyn.ISyn) clsyn.ISyn {
	if PrintSteps {
		fmt.Printf("\n\n%d — %T\n\t%v\n\t%v\n", me.NumStepsTaken, syn, me.Args, me.Locals)
	}
	me.NumStepsTaken++
	switch n := syn.(type) {
	case *clsyn.ExprLitFloat, *clsyn.ExprLitRune, *clsyn.ExprLitText, *clsyn.ExprLitUInt:
		return syn
	case *clsyn.ExprIdent:
		if PrintSteps {
			fmt.Printf("\t%s\n", n.Name)
		}
		return me.reduce(me.resolveIdent(n.Name))
	case *clsyn.SynDef:
		if PrintSteps {
			fmt.Printf("\t%s\n", n.Name)
		}
		if len(me.Args) < len(n.Args) {
			for i, arg := range me.Args {
				me.Locals[n.Args[i]] = arg
			}
			me.Args = []clsyn.ISyn{}
		} else {
			for i, arg := range n.Args {
				me.Locals[arg] = me.Args[i]
			}
			me.Args = me.Args[len(n.Args):]
		}
		val := me.reduce(n.Body)
		return val
	case *clsyn.ExprCall:
		me.Args = append([]clsyn.ISyn{n.Arg}, me.Args...)
		me.NumApplications++
		return me.reduce(n.Callee)
	case *clsyn.ExprLetIn:
		for _, def := range n.Defs {
			me.Locals[def.Name] = def
		}
		return me.reduce(n.Body)
	}
	panic(syn)
}

func (me *naiveMachine) resolveIdent(name string) (syn clsyn.ISyn) {
	if syn = me.Locals[name]; syn == nil {
		if syn = me.Globals[name]; syn == nil {
			panic("undefined: " + name)
		}
	}
	return
}
