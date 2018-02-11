package climpl

import (
	"errors"

	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type compilation func(clsyn.IExpr, env) code

type compilationN func(int, clsyn.IExpr, env) code

type env map[string]int

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
		argsenv := make(env, len(global.Args))
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

func (me *gMachine) compileGlobal(bodyexpr clsyn.IExpr, argsEnv env) (bodycode code, err error) {
	defer clutil.Catch(&err)
	bodycode = append(me.compileExprStrict(bodyexpr, argsEnv),
		instr{Op: INSTR_UPDATE, Int: len(argsEnv)},
		instr{Op: INSTR_POP, Int: len(argsEnv)},
		instr{Op: INSTR_UNWIND},
	)
	return
}

func (me *gMachine) compileExprStrict(expression clsyn.IExpr, argsEnv env) code {
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
			if op, _ := callee.Callee.(*clsyn.ExprIdent); op != nil && op.OpLike {
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

func (me *gMachine) compileExprStrictSplitSlide(offset int, expression clsyn.IExpr, argsEnv env) code {
	// does not `offset` the given `argsEnv` though, which is assumed to already be properly `envOffsetBy`
	return append(append(code{{Op: INSTR_CASE_SPLIT, Int: offset}},
		me.compileExprStrict(expression, argsEnv)...,
	), instr{Op: INSTR_SLIDE, Int: offset})
}

func (me *gMachine) compileExprLazy(expression clsyn.IExpr, argsEnv env) code {
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

func (me *gMachine) compileLet(compbody compilation, let *clsyn.ExprLetIn, argsEnv env) (instrs code) {
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

func (me *gMachine) compileCaseAlts(compn compilationN, caseAlts []*clsyn.SynCaseAlt, argsEnv env) (jumpblocks []code) {
	jumpblocks = make([]code, len(caseAlts))
	for _, alt := range caseAlts {
		n := len(alt.Binds)
		bodyargsenv := me.envOffsetBy(argsEnv, n)
		for i, name := range alt.Binds {
			bodyargsenv[name] = i // = n - (i + 1)
		}
		jumpblocks[alt.Tag-1] = compn(n, alt.Body, bodyargsenv)
	}
	return
}

func (*gMachine) envOffsetBy(argsEnv env, offsetBy int) (envOffset env) {
	envOffset = make(env, len(argsEnv))
	for k, v := range argsEnv {
		envOffset[k] = v + offsetBy
	}
	return
}
