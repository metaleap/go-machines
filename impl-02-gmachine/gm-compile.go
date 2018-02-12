package climpl

import (
	"errors"
	"strconv"
	"time"

	// "github.com/metaleap/go-corelang"
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type compilation func(clsyn.IExpr, env) code

type compilationN func(int, clsyn.IExpr, env) code

type env map[string]int

var primDyadicsForStrict = map[string]instruction{
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

var primsPrecompiledForLazy = map[string]nodeGlobal{
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

	for primname, primnode := range primsPrecompiledForLazy {
		me.Globals[primname] = me.Heap.Alloc(primnode)
	}
	for _, global := range mod.Defs {
		if node, err := me.compileGlobal_SchemeSC(global); err != nil {
			errs = append(errs, errors.New(global.Name+": "+err.Error()))
		} else {
			me.Globals[global.Name] = me.Heap.Alloc(node)
		}
	}
	return &me, errs
}

func (me *gMachine) compileGlobal_SchemeSC(global *clsyn.SynDef) (node nodeGlobal, err error) {
	defer clutil.Catch(&err)
	argsenv := make(env, len(global.Args))
	for i, arg := range global.Args {
		argsenv[arg] = i
	}
	node = nodeGlobal{len(global.Args), me.compileGlobalBody_SchemeR(global.Body, argsenv)}
	return
}

func (me *gMachine) compileGlobalBody_SchemeR(bodyexpr clsyn.IExpr, argsEnv env) code {
	return append(me.compileExprStrict_SchemeE(bodyexpr, argsEnv),
		instr{Op: INSTR_UPDATE, Int: len(argsEnv)},
		instr{Op: INSTR_POP, Int: len(argsEnv)},
		instr{Op: INSTR_UNWIND},
	)
}

func (me *gMachine) compileExprStrict_SchemeE(expression clsyn.IExpr, argsEnv env) code {
	switch expr := expression.(type) {
	case *clsyn.ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *clsyn.ExprLetIn:
		return me.compileLet(me.compileExprStrict_SchemeE, expr, argsEnv)
	case *clsyn.ExprCtor:
		return me.compileCtorAppl(me.compileExprStrict_SchemeE, expr, nil, argsEnv)
	case *clsyn.ExprCaseOf:
		return append(me.compileExprStrict_SchemeE(expr.Scrut, argsEnv), instr{Op: INSTR_CASE_JUMP,
			CaseJump: me.compileCaseAlts_SchemeD(me.compileExprStrictSplitSlide_SchemeA, expr.Alts, argsEnv)})
	case *clsyn.ExprCall:
		if ctor, ctorrevargs := expr.FlattenedIfEffectivelyCtor(); ctor != nil {
			return me.compileCtorAppl(me.compileExprStrict_SchemeE, ctor, ctorrevargs, argsEnv)
		}
		switch callee := expr.Callee.(type) {
		case *clsyn.ExprIdent:
			if callee.Name == "neg" {
				return append(me.compileExprStrict_SchemeE(expr.Arg, argsEnv), instr{Op: INSTR_PRIM_AR_NEG})
			}
		case *clsyn.ExprCall:
			if op, _ := callee.Callee.(*clsyn.ExprIdent); op != nil {
				if primdyadic := primDyadicsForStrict[op.Name]; primdyadic != 0 {
					return append(append(
						me.compileExprStrict_SchemeE(expr.Arg, argsEnv),
						me.compileExprStrict_SchemeE(callee.Arg, me.envOffsetBy(argsEnv, 1))...,
					), instr{Op: primdyadic})
				}
			} else if maybeif, _ := callee.Callee.(*clsyn.ExprCall); maybeif != nil {
				if ifname, _ := maybeif.Callee.(*clsyn.ExprIdent); ifname != nil && ifname.Name == "if" {
					cond, condthen, condelse := maybeif.Arg, callee.Arg, expr.Arg
					return append(me.compileExprStrict_SchemeE(cond, argsEnv), instr{
						Op:       INSTR_PRIM_COND,
						CondThen: me.compileExprStrict_SchemeE(condthen, argsEnv),
						CondElse: me.compileExprStrict_SchemeE(condelse, argsEnv),
					})
				}
			}
		}
	}
	return append(me.compileExprLazy_SchemeC(expression, argsEnv), instr{Op: INSTR_EVAL})
}

func (me *gMachine) compileExprStrictSplitSlide_SchemeA(offset int, expression clsyn.IExpr, argsEnv env) code {
	// does not `offset` the given `argsEnv`, expected to be passed already-`envOffsetBy`
	return append(append(code{{Op: INSTR_CASE_SPLIT, Int: offset}},
		me.compileExprStrict_SchemeE(expression, argsEnv)...,
	), instr{Op: INSTR_SLIDE, Int: offset})
}

func (me *gMachine) compileExprLazy_SchemeC(expression clsyn.IExpr, argsEnv env) code {
	switch expr := expression.(type) {
	case *clsyn.ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *clsyn.ExprIdent:
		if i, islocal := argsEnv[expr.Name]; islocal {
			return code{{Op: INSTR_PUSHARG, Int: i}}
		}
		return code{{Op: INSTR_PUSHGLOBAL, Name: expr.Name}}
	case *clsyn.ExprCall:
		if ctor, ctorrevargs := expr.FlattenedIfEffectivelyCtor(); ctor != nil {
			return me.compileCtorAppl(me.compileExprLazy_SchemeC, ctor, ctorrevargs, argsEnv)
		}
		return append(append(
			me.compileExprLazy_SchemeC(expr.Arg, argsEnv),
			me.compileExprLazy_SchemeC(expr.Callee, me.envOffsetBy(argsEnv, 1))...,
		), instr{Op: INSTR_MAKEAPPL})
	case *clsyn.ExprLetIn:
		return me.compileLet(me.compileExprLazy_SchemeC, expr, argsEnv)
	case *clsyn.ExprCtor:
		return me.compileCtorAppl(me.compileExprLazy_SchemeC, expr, nil, argsEnv)
	case *clsyn.ExprCaseOf:
		dynname := "#case" + strconv.FormatInt(time.Now().UnixNano(), 16)
		dynglobaldef, dynglobalcall := expr.ExtractIntoDef(dynname, true, clsyn.LookupEnvFrom(nil, me.Globals, nil, nil))
		// println("DEF:")
		// src, _ := (&corelang.SyntaxTreePrinter{}).Def(dynglobaldef)
		// println(src)
		// println("CALL:")
		// src, _ = (&corelang.SyntaxTreePrinter{}).Expr(dynglobalcall)
		// println(src)
		if dynglobalnode, err := me.compileGlobal_SchemeSC(dynglobaldef); err != nil {
			panic(err)
		} else {
			me.Globals[dynglobaldef.Name] = me.Heap.Alloc(dynglobalnode)
		}
		return me.compileExprLazy_SchemeC(dynglobalcall, argsEnv)
	default:
		panic(expr)
	}
}

func (me *gMachine) compileCtorAppl(comp compilation, ctor *clsyn.ExprCtor, reverseArgs []clsyn.IExpr, argsEnv env) code {
	if len(reverseArgs) != ctor.Arity {
		dynglobalname := "#ctor#" + strconv.Itoa(ctor.Tag) + "#" + strconv.Itoa(ctor.Arity)
		dynglobaladdr := me.Globals[dynglobalname]
		if dynglobaladdr == 0 {
			dynglobaldef := ctor.ExtractIntoDef(dynglobalname, true)
			if dynglobalnode, err := me.compileGlobal_SchemeSC(dynglobaldef); err != nil {
				panic(err)
			} else {
				me.Globals[dynglobalname] = me.Heap.Alloc(dynglobalnode)
			}
			// println("\n\nDEF:\n\n")
			// src, _ := (&corelang.SyntaxTreePrinter{}).Def(dynglobaldef)
			// println(src)
		}
		dynglobalcall := clsyn.Call(clsyn.Id(dynglobalname), reverseArgs...)
		// println("\n\nCALL:\n\n")
		// src, _ := (&corelang.SyntaxTreePrinter{}).Expr(dynglobalcall)
		// println(src)
		return comp(dynglobalcall, argsEnv)
	}
	instrs := make(code, 0, len(reverseArgs)+len(reverseArgs))
	for i, arg := range reverseArgs {
		instrs = append(instrs, comp(arg, me.envOffsetBy(argsEnv, i))...)
	}
	return append(instrs, instr{Op: INSTR_CTOR_PACK, Int: ctor.Tag, CtorArity: ctor.Arity})
}

func (me *gMachine) compileLet(compbody compilation, let *clsyn.ExprLetIn, argsEnv env) (instrs code) {
	n := len(let.Defs)
	if let.Rec {
		instrs = code{{Op: INSTR_ALLOC, Int: n}}
	}

	bodyargsenv := me.envOffsetBy(argsEnv, n)
	for i, def := range let.Defs {
		if bodyargsenv[def.Name] = n - (i + 1); !let.Rec {
			instrs = append(instrs, me.compileExprLazy_SchemeC(def.Body, me.envOffsetBy(argsEnv, i))...)
		}
	}

	if let.Rec {
		for i, def := range let.Defs {
			instrs = append(instrs, me.compileExprLazy_SchemeC(def.Body, bodyargsenv)...)
			instrs = append(instrs, instr{Op: INSTR_UPDATE, Int: n - (i + 1)})
		}
	}

	instrs = append(instrs, compbody(let.Body, bodyargsenv)...)
	instrs = append(instrs, instr{Op: INSTR_SLIDE, Int: n})
	return
}

func (me *gMachine) compileCaseAlts_SchemeD(compn compilationN, caseAlts []*clsyn.SynCaseAlt, argsEnv env) (jumpblocks []code) {
	jumpblocks = make([]code, len(caseAlts))
	for _, alt := range caseAlts {
		n := len(alt.Binds)
		bodyargsenv := me.envOffsetBy(argsEnv, n)
		for i, name := range alt.Binds {
			bodyargsenv[name] = i // = n - (i + 1)
		}
		jumpblocks[alt.Tag] = compn(n, alt.Body, bodyargsenv)
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
