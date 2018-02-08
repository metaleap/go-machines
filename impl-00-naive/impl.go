package climpl

import (
	"fmt"

	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type naiveMachine struct {
	globals       map[string]clsyn.ISyn
	locals        map[string]clsyn.ISyn
	args          []clsyn.ISyn
	numStepsTaken int
	printSteps    bool
}

func CompileToMachine(mod *clsyn.SynMod) clutil.IMachine {
	globals := make(map[string]clsyn.ISyn, len(mod.Defs))
	for _, def := range mod.Defs {
		globals[def.Name] = def
	}
	return &naiveMachine{globals: globals}
}

func (me *naiveMachine) Eval(name string) (val interface{}, numSteps int, err error) {
	defer clutil.Catch(&err)
	def := me.resolveIdent(name)
	me.numStepsTaken, me.locals = 0, map[string]clsyn.ISyn{}
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
	numSteps = me.numStepsTaken
	return
}

func (me *naiveMachine) reduce(syn clsyn.ISyn) clsyn.ISyn {
	if me.printSteps {
		fmt.Printf("\n\n%d — %T\n\t%v\n\t%v\n", me.numStepsTaken, syn, me.args, me.locals)
	}
	me.numStepsTaken++
	switch n := syn.(type) {
	case *clsyn.ExprLitFloat, *clsyn.ExprLitRune, *clsyn.ExprLitText, *clsyn.ExprLitUInt:
		return syn
	case *clsyn.ExprIdent:
		if me.printSteps {
			fmt.Printf("\t%s\n", n.Name)
		}
		return me.reduce(me.resolveIdent(n.Name))
	case *clsyn.SynDef:
		if me.printSteps {
			fmt.Printf("\t%s\n", n.Name)
		}
		if len(me.args) < len(n.Args) {
			for i, arg := range me.args {
				me.locals[n.Args[i]] = arg
			}
			me.args = []clsyn.ISyn{}
		} else {
			for i, arg := range n.Args {
				me.locals[arg] = me.args[i]
			}
			me.args = me.args[len(n.Args):]
		}
		val := me.reduce(n.Body)
		return val
	case *clsyn.ExprCall:
		me.args = append([]clsyn.ISyn{n.Arg}, me.args...)
		return me.reduce(n.Callee)
	case *clsyn.ExprLetIn:
		for _, def := range n.Defs {
			me.locals[def.Name] = def
		}
		return me.reduce(n.Body)
	}
	panic(syn)
}

func (me *naiveMachine) resolveIdent(name string) (syn clsyn.ISyn) {
	if syn = me.locals[name]; syn == nil {
		if syn = me.globals[name]; syn == nil {
			panic("undefined: " + name)
		}
	}
	return
}
