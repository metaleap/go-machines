package climpl

import (
	"errors"

	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type compilation func(clsyn.IExpr, map[string]int) code

var primDyadic = map[string]instruction{
	"+":  INSTR_PRIM_AR_ADD,
	"-":  INSTR_PRIM_AR_SUB,
	"*":  INSTR_PRIM_AR_MUL,
	"/":  INSTR_PRIM_AR_DIV,
	"==": INSTR_PRIM_CMP_EQ,
	"!=": INSTR_PRIM_CMP_NEQ,
	"<":  INSTR_PRIM_CMP_LT,
	"<=": INSTR_PRIM_CMP_LEQ,
	">":  INSTR_PRIM_CMP_GT,
	">=": INSTR_PRIM_CMP_GEQ,
}

// var preCompiledPrims = map[string]nodeGlobal{
// 	"+":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_ADD}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"-":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_SUB}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"*":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_MUL}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"/":   {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_DIV}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"neg": {1, code{{Op: INSTR_PUSHARG}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_AR_NEG}, {Op: INSTR_UPDATE, Int: 1}, {Op: INSTR_POP, Int: 1}, {Op: INSTR_UNWIND}}},

// 	"==": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_EQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"!=": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_NEQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"<":  {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_LT}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	"<=": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_LEQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	">":  {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_GT}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},
// 	">=": {2, code{{Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PUSHARG, Int: 1}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_CMP_GEQ}, {Op: INSTR_UPDATE, Int: 2}, {Op: INSTR_POP, Int: 2}, {Op: INSTR_UNWIND}}},

// 	"if": {3, code{{Op: INSTR_PUSHARG}, {Op: INSTR_EVAL}, {Op: INSTR_PRIM_COND, CondThen: code{{Op: INSTR_PUSHARG, Int: 1}}, CondElse: code{{Op: INSTR_PUSHARG, Int: 2}}}, {Op: INSTR_UPDATE, Int: 3}, {Op: INSTR_POP, Int: 3}, {Op: INSTR_UNWIND}}},
// }

func CompileToMachine(mod *clsyn.SynMod) (clutil.IMachine, []error) {
	errs, me := []error{}, gMachine{
		Heap:    clutil.Heap{},
		Globals: make(clutil.Env, len(mod.Defs)),
		Stack:   make(clutil.Stack, 0, 128),
	}

	// for primname, primnode := range preCompiledPrims {
	// 	me.Globals[primname] = me.Heap.Alloc(primnode)
	// }

	for _, global := range mod.Defs {
		argsenv := make(map[string]int, len(global.Args))
		for i, arg := range global.Args {
			argsenv[arg] = i
		}

		if bodycode, err := me.compileGlobal(global.Body, argsenv); err != nil {
			errs = append(errs, errors.New(global.Name+": "+err.Error()))
		} else {
			me.Globals[global.Name] = me.Heap.Alloc(nodeGlobal{len(global.Args), bodycode})
		}
	}
	return &me, errs
}

func (me *gMachine) compileGlobal(bodyexpr clsyn.IExpr, argsEnv map[string]int) (bodycode code, err error) {
	defer clutil.Catch(&err)
	bodycode = append(me.compileExprStrict(bodyexpr, argsEnv),
		instr{Op: INSTR_UPDATE, Int: len(argsEnv)},
		instr{Op: INSTR_POP, Int: len(argsEnv)},
		instr{Op: INSTR_UNWIND},
	)
	return
}

func (me *gMachine) compileExprStrict(expression clsyn.IExpr, argsEnv map[string]int) code {
	switch expr := expression.(type) {
	case *clsyn.ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *clsyn.ExprLetIn:
		return me.compileLet(me.compileExprStrict, expr, argsEnv)
	case *clsyn.ExprCall:
		switch callee := expr.Callee.(type) {
		case *clsyn.ExprIdent:
			if callee.Name == "neg" {
				return append(me.compileExprStrict(expr.Arg, argsEnv), instr{Op: INSTR_PRIM_AR_NEG})
			}
		case *clsyn.ExprCall:
			if op, _ := callee.Callee.(*clsyn.ExprIdent); op != nil {
				if primdyadic := primDyadic[op.Name]; primdyadic != 0 {
					return append(append(
						me.compileExprStrict(expr.Arg, argsEnv),
						me.compileExprStrict(callee.Arg, me.envOffsetBy(argsEnv, 1))...,
					), instr{Op: primdyadic})
				}
			} else if maybeif, _ := callee.Callee.(*clsyn.ExprCall); maybeif != nil {
				if ifname, _ := maybeif.Callee.(*clsyn.ExprIdent); ifname != nil && ifname.Name == "if" {
					cond, condthen, condelse := maybeif.Arg, callee.Arg, expr.Arg
					return append(me.compileExprStrict(cond, argsEnv), instr{
						Op:       INSTR_PRIM_COND,
						CondThen: me.compileExprStrict(condthen, argsEnv),
						CondElse: me.compileExprStrict(condelse, argsEnv),
					})
				}
			}
		}
	}
	return append(me.compileExprLazy(expression, argsEnv), instr{Op: INSTR_EVAL})
}

func (me *gMachine) compileExprLazy(expression clsyn.IExpr, argsEnv map[string]int) code {
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
			me.compileExprLazy(expr.Arg, argsEnv),
			me.compileExprLazy(expr.Callee, me.envOffsetBy(argsEnv, 1))...,
		), instr{Op: INSTR_MAKEAPPL})
	case *clsyn.ExprLetIn:
		return me.compileLet(me.compileExprLazy, expr, argsEnv)
	default:
		panic(expr)
	}
}

func (me *gMachine) compileLet(compbody compilation, let *clsyn.ExprLetIn, argsEnv map[string]int) (instrs code) {
	n := len(let.Defs)
	if let.Rec {
		instrs = code{{Op: INSTR_ALLOC, Int: n}}
	}

	bodyargsenv := me.envOffsetBy(argsEnv, n)
	for i, def := range let.Defs {
		if bodyargsenv[def.Name] = n - (i + 1); !let.Rec {
			instrs = append(instrs, me.compileExprLazy(def.Body, me.envOffsetBy(argsEnv, i))...)
		}
	}

	if let.Rec {
		for i, def := range let.Defs {
			instrs = append(instrs, me.compileExprLazy(def.Body, bodyargsenv)...)
			instrs = append(instrs, instr{Op: INSTR_UPDATE, Int: n - (i + 1)})
		}
	}

	instrs = append(instrs, compbody(let.Body, bodyargsenv)...)
	instrs = append(instrs, instr{Op: INSTR_SLIDE, Int: n})
	return
}

func (*gMachine) envOffsetBy(env map[string]int, offsetBy int) (envOffset map[string]int) {
	envOffset = make(map[string]int, len(env))
	for k, v := range env {
		envOffset[k] = v + offsetBy
	}
	return
}
