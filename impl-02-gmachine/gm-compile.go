package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

func CompileToMachine(mod *clsyn.SynMod) clutil.IMachine {
	me := gMachine{
		Heap:    clutil.Heap{},
		Globals: make(map[string]clutil.Addr, len(mod.Defs)),
	}
	for _, def := range mod.Defs {
		argsenv := make(map[string]int, len(def.Args))
		for i, arg := range def.Args {
			argsenv[arg] = i
		}

		me.Globals[def.Name] = me.alloc(nodeGlobal{
			NumArgs: len(argsenv),
			Code:    me.compileR(def.Body, argsenv),
		})
	}
	return &me
}

func (me *gMachine) compileR(expr clsyn.IExpr, argsEnv map[string]int) code {
	return append(
		me.compileC(expr, argsEnv),
		instr{Op: INSTR_SLIDE, Int: 1 + len(argsEnv)},
		instr{Op: INSTR_UNWIND},
	)
}

func (me *gMachine) compileC(expression clsyn.IExpr, argsEnv map[string]int) code {
	switch expr := expression.(type) {
	case *clsyn.ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *clsyn.ExprIdent:
		if i, islocal := argsEnv[expr.Name]; islocal {
			return code{{Op: INSTR_PUSHARG, Int: i}}
		}
		return code{{Op: INSTR_PUSHGLOBAL, Name: expr.Name}}
	case *clsyn.ExprCall:
		return append(append(
			me.compileC(expr.Arg, argsEnv),
			me.compileC(expr.Callee, me.envOffsetBy(argsEnv, 1))...,
		), instr{Op: INSTR_MAKEAPPL})
	default:
		panic(expr)
	}
}

func (*gMachine) envOffsetBy(env map[string]int, offsetBy int) (envOffset map[string]int) { // TODO: ditch this and replace with a new `offset` arg to compileAppl()
	envOffset = make(map[string]int, len(env))
	for k, v := range env {
		envOffset[k] = v + offsetBy
	}
	return
}
