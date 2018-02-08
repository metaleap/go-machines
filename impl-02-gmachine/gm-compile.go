package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

func CompileToMachine(mod *clsyn.SynMod) clutil.IMachine {
	me := gMachine{Heap: clutil.Heap{}, Globals: make(map[string]clutil.Addr, len(mod.Defs))}
	for _, def := range mod.Defs {
		args := make(map[string]int, len(def.Args))
		for i, arg := range def.Args {
			args[arg] = i
		}
		me.Globals[def.Name] = me.alloc(nodeGlobal{
			NumArgs: len(args),
			Code:    me.compileExpr(def.Body, args),
		})
	}
	return &me
}

func (me *gMachine) compileExpr(expr clsyn.IExpr, env map[string]int) code {
	return append(
		me.compileAppl(expr, env),
		instr{Op: INSTR_SLIDE, Int: 1 + len(env)},
		instr{Op: INSTR_UNWIND},
	)
}

func (me *gMachine) compileAppl(expression clsyn.IExpr, env map[string]int) code {
	switch expr := expression.(type) {
	case *clsyn.ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *clsyn.ExprIdent:
		if i, islocal := env[expr.Name]; islocal {
			return code{{Op: INSTR_PUSH, Int: i}}
		}
		return code{{Op: INSTR_PUSHGLOBAL, Name: expr.Name}}
	case *clsyn.ExprCall:
		return append(append(
			me.compileAppl(expr.Arg, env),
			me.compileAppl(expr.Callee, me.envOffsetBy(env, 1))...,
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
