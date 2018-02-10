package climpl

import (
	"errors"

	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

const MARK2_LAZY = true

type compilation func(clsyn.IExpr, map[string]int) code

var preCompiledPrims = map[string]nodeGlobal{
	"+":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_ADD}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"-":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_SUB}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"*":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_MUL}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"/":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_DIV}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"neg": {1, code{{Op: INSTR_PUSHARG}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_NEG}, {Op: INSTR_UPDATE, Int: 1}, {Op: INSTR_POP, Int: 1}, {Op: INSTR_UNWIND}}},

	"==": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_EQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"!=": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_NEQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"<":  {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_LT}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	"<=": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_LEQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	">":  {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_GT}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
	">=": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_GEQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},

	"if": {3, code{{Op: INSTR_PUSHARG}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_COND, CondThen: code{{Op: INSTR_PUSHARG, Int: 1}}, CondElse: code{{Op: INSTR_PUSHARG, Int: 2}}}, {Op: INSTR_UPDATE, Int: 3}, {Op: INSTR_POP, Int: 3}, {Op: INSTR_UNWIND}}},
}

func CompileToMachine(mod *clsyn.SynMod) (clutil.IMachine, []error) {
	errs, me := []error{}, gMachine{
		Heap:    clutil.Heap{},
		Globals: make(clutil.Env, len(mod.Defs)),
		Stack:   make(clutil.Stack, 0, 128),
	}

	for primname, primnode := range preCompiledPrims {
		me.Globals[primname] = me.Heap.Alloc(primnode)
	}

	for _, global := range mod.Defs {
		argsenv := make(map[string]int, len(global.Args))
		for i, arg := range global.Args {
			argsenv[arg] = i
		}

		if bodycode, err := me.compileTopLevelDefBody(global.Body, argsenv); err != nil {
			errs = append(errs, errors.New(global.Name+": "+err.Error()))
		} else {
			me.Globals[global.Name] = me.Heap.Alloc(nodeGlobal{len(argsenv), bodycode})
		}
	}
	return &me, errs
}

func (me *gMachine) compileTopLevelDefBody(bodyexpr clsyn.IExpr, argsEnv map[string]int) (bodycode code, err error) {
	defer clutil.Catch(&err)
	numargs, codeexpr := len(argsEnv), me.compileExpr(bodyexpr, argsEnv)
	if MARK2_LAZY {
		bodycode = append(codeexpr,
			instr{Op: INSTR_UPDATE, Int: numargs},
			instr{Op: INSTR_POP, Int: numargs},
			instr{Op: INSTR_UNWIND},
		)
	} else {
		bodycode = append(codeexpr,
			instr{Op: INSTR_SLIDE, Int: 1 + numargs},
			instr{Op: INSTR_UNWIND},
		)
	}
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
			return me.compileLetRec(me.compileExpr, expr, argsEnv)
		}
		return me.compileLet(me.compileExpr, expr, argsEnv)
	default:
		panic(expr)
	}
}

func (me *gMachine) compileLet(compbody compilation, let *clsyn.ExprLetIn, argsEnv map[string]int) code {
	n, instrs := len(let.Defs), code{}

	bodyargsenv := me.envOffsetBy(argsEnv, n)
	for i, def := range let.Defs {
		bodyargsenv[def.Name] = n - (i + 1)
		instrs = append(instrs, me.compileExpr(def.Body, me.envOffsetBy(argsEnv, i))...)
	}

	instrs = append(instrs, compbody(let.Body, bodyargsenv)...)
	return append(instrs, instr{Op: INSTR_SLIDE, Int: n})
}

func (me *gMachine) compileLetRec(compbody compilation, let *clsyn.ExprLetIn, argsEnv map[string]int) code {
	n := len(let.Defs)
	instrs := code{{Op: INSTR_ALLOC, Int: n}}

	bodyargsenv := me.envOffsetBy(argsEnv, n)
	for i, def := range let.Defs {
		bodyargsenv[def.Name] = n - (i + 1)
	}

	for i, def := range let.Defs {
		instrs = append(instrs, me.compileExpr(def.Body, bodyargsenv)...)
		instrs = append(instrs, instr{Op: INSTR_UPDATE, Int: n - (i + 1)})
	}

	instrs = append(instrs, compbody(let.Body, bodyargsenv)...)
	return append(instrs, instr{Op: INSTR_SLIDE, Int: n})
}

func (*gMachine) envOffsetBy(env map[string]int, offsetBy int) (envOffset map[string]int) {
	envOffset = make(map[string]int, len(env))
	for k, v := range env {
		envOffset[k] = v + offsetBy
	}
	return
}
