package climpl

import (
	"errors"

	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

const MARK2_LAZY = true

type compilation func(clsyn.IExpr, map[string]int) code

func CompileToMachine(mod *clsyn.SynMod) (clutil.IMachine, []error) {
	errs, me := []error{}, gMachine{
		Heap:    clutil.Heap{},
		Globals: make(clutil.Env, len(mod.Defs)),
		Stack:   make(clutil.Stack, 0, 128),
	}

	for _, global := range mod.Defs {
		argsenv := make(map[string]int, len(global.Args))
		for i, arg := range global.Args {
			argsenv[arg] = i
		}

		if bodycode, err := me.compileBody(global.Body, argsenv); err != nil {
			errs = append(errs, errors.New(global.Name+": "+err.Error()))
		} else {
			me.Globals[global.Name] = me.Heap.Alloc(nodeGlobal{len(argsenv), bodycode})
		}
	}
	return &me, errs
}

func (me *gMachine) compileBody(bodyexpr clsyn.IExpr, argsEnv map[string]int) (bodycode code, err error) {
	defer clutil.Catch(&err)
	numargs, codeexpr := len(argsEnv), me.compileExpr(bodyexpr, argsEnv)

	if MARK2_LAZY {
		bodycode = append(codeexpr,
			instr{Op: INSTR_UPDATE, Int: numargs},
			instr{Op: INSTR_POP, Int: numargs},
			instr{Op: INSTR_UNWIND},
		)
	}

	bodycode = append(codeexpr,
		instr{Op: INSTR_SLIDE, Int: 1 + numargs},
		instr{Op: INSTR_UNWIND},
	)
	return
}

func (me *gMachine) compileExpr(expression clsyn.IExpr, argsEnv map[string]int) code {
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
			me.compileExpr(expr.Arg, argsEnv),
			me.compileExpr(expr.Callee, me.envOffsetBy(argsEnv, 1))...,
		), instr{Op: INSTR_MAKEAPPL})
	case *clsyn.ExprLetIn:
		if expr.Rec {
			return me.compileLetRec(me.compileExpr)
		}
		return me.compileLet(me.compileExpr)
	default:
		panic(expr)
	}
}

func (me *gMachine) compileLet(comp compilation) code {
	return nil
}

func (me *gMachine) compileLetRec(comp compilation) code {
	return nil
}

func (*gMachine) envOffsetBy(env map[string]int, offsetBy int) (envOffset map[string]int) { // TODO: ditch this and replace with a new `offset` arg to compileAppl()
	envOffset = make(map[string]int, len(env))
	for k, v := range env {
		envOffset[k] = v + offsetBy
	}
	return
}
