package climpl

import (
	"errors"
	"strconv"
	"time"

	// "github.com/metaleap/go-corelang"
	. "github.com/metaleap/go-corelang/syn"
	util "github.com/metaleap/go-corelang/util"
)

type compilation func(IExpr, env) code

type compilationN func(int, IExpr, env) code

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

var primsMark7Globals = map[string]*SynDef{
	"negate": {TopLevel: true, Name: "neg", Args: []string{"x"}, Body: Ap(Id("neg"), Id("x"))},
	"if":     {TopLevel: true, Name: "if", Args: []string{"c", "t", "f"}, Body: Ap(Ap(Ap(Id("if"), Id("c")), Id("t")), Id("f"))},
	"True":   {TopLevel: true, Name: "True", Body: Ct(2, 0)},
	"False":  {TopLevel: true, Name: "False", Body: Ct(1, 0)},
}

func init() {
	for _, opname := range []string{"+", "-", "*", "/", "==", "!=", "<", "<=", ">", ">="} {
		primsMark7Globals[opname] = &SynDef{TopLevel: true, Name: opname, Args: []string{"x", "y"}, Body: Ap(Ap(Id(opname), Id("x")), Id("y"))}
	}
}

func CompileToMachine(mod *SynMod) (util.IMachine, []error) {
	errs, me := []error{}, gMachine{
		Heap:    make(util.HeapA, 1, 1024*1024),
		Globals: make(util.Env, len(mod.Defs)),
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

func (me *gMachine) compileGlobal_SchemeSC(global *SynDef) (node nodeGlobal, err error) {
	defer util.Catch(&err)
	argsenv := make(env, len(global.Args))
	for i, arg := range global.Args {
		argsenv[arg] = i
	}
	node = nodeGlobal{len(global.Args), me.compileGlobalBody_SchemeR(global.Body, argsenv)}
	if global.Name == "fac" {
		println(node.Code.String())
	}
	return
}

func (me *gMachine) compileGlobalBody_SchemeR(bodyexpr IExpr, argsEnv env) code {
	return append(me.compileExprStrict_SchemeE(bodyexpr, argsEnv),
		instr{Op: INSTR_UPDATE, Int: len(argsEnv)},
		instr{Op: INSTR_POP, Int: len(argsEnv)},
		instr{Op: INSTR_UNWIND},
	)
}

func (me *gMachine) compileExprStrict_SchemeE(expression IExpr, argsEnv env) code {
	switch expr := expression.(type) {
	case *ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *ExprLetIn:
		return me.compileLet(me.compileExprStrict_SchemeE, expr, argsEnv, INSTR_SLIDE)
	case *ExprCtor:
		comp := me.compileExprLazy_SchemeC
		if !MARK7 {
			comp = me.compileExprStrict_SchemeE
		}
		return me.compileCtorAppl(comp, expr, nil, argsEnv, MARK7)
	case *ExprCaseOf:
		return append(me.compileExprStrict_SchemeE(expr.Scrut, argsEnv), instr{Op: INSTR_CASE_JUMP,
			CaseJump: me.compileCaseAlts_SchemeD(me.compileExprStrictSplitSlide_SchemeA, expr.Alts, argsEnv)})
	case *ExprCall:
		if ctor, ctorrevargs := expr.FlattenedIfEffectivelyCtor(); ctor != nil {
			comp := me.compileExprLazy_SchemeC
			if !MARK7 {
				comp = me.compileExprStrict_SchemeE
			}
			return me.compileCtorAppl(comp, ctor, ctorrevargs, argsEnv, MARK7)
		}

		if instrs := me.compilePrimsMaybe(me.compileExprStrict_SchemeE, expr, argsEnv, 1, MARK7); len(instrs) > 0 {
			return instrs
		}
	}
	return append(me.compileExprLazy_SchemeC(expression, argsEnv), instr{Op: INSTR_EVAL})
}

func (me *gMachine) compileExprStrictSplitSlide_SchemeA(offset int, expression IExpr, argsEnv env) code {
	// does not `offset` the given `argsEnv`, expected to be passed already-`envOffsetBy`
	return append(append(code{{Op: INSTR_CASE_SPLIT, Int: offset}},
		me.compileExprStrict_SchemeE(expression, argsEnv)...,
	), instr{Op: INSTR_SLIDE, Int: offset})
}

func (me *gMachine) compileExprStrictMark7_SchemeB(expression IExpr, argsEnv env) code {
	switch expr := expression.(type) {
	case *ExprLitUInt:
		return code{{Op: INSTR_MARK7_PUSHINTVAL, Int: int(expr.Lit)}}
	case *ExprLetIn:
		return me.compileLet(me.compileExprStrictMark7_SchemeB, expr, argsEnv, INSTR_POP)
	case *ExprCall:
		if instrs := me.compilePrimsMaybe(me.compileExprStrictMark7_SchemeB, expr, argsEnv, 0, false); len(instrs) > 0 {
			return instrs
		}
	}
	return append(me.compileExprStrict_SchemeE(expression, argsEnv), instr{Op: INSTR_MARK7_PUSHNODEINT})
}

func (me *gMachine) compileExprLazy_SchemeC(expression IExpr, argsEnv env) code {
	switch expr := expression.(type) {
	case *ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *ExprIdent:
		if i, islocal := argsEnv[expr.Name]; islocal {
			return code{{Op: INSTR_PUSHARG, Int: i}}
		}
		return code{{Op: INSTR_PUSHGLOBAL, Name: expr.Name}}
	case *ExprCall:
		if ctor, ctorrevargs := expr.FlattenedIfEffectivelyCtor(); ctor != nil {
			return me.compileCtorAppl(me.compileExprLazy_SchemeC, ctor, ctorrevargs, argsEnv, false)
		}
		return append(append(
			me.compileExprLazy_SchemeC(expr.Arg, argsEnv),
			me.compileExprLazy_SchemeC(expr.Callee, me.envOffsetBy(argsEnv, 1))...,
		), instr{Op: INSTR_MAKEAPPL})
	case *ExprLetIn:
		return me.compileLet(me.compileExprLazy_SchemeC, expr, argsEnv, INSTR_SLIDE)
	case *ExprCtor:
		return me.compileCtorAppl(me.compileExprLazy_SchemeC, expr, nil, argsEnv, false)
	case *ExprCaseOf:
		dynname := "#case" + strconv.FormatInt(time.Now().UnixNano(), 16)
		dynglobaldef, dynglobalcall := expr.ExtractIntoDef(dynname, true, NewLookupEnv(nil, me.Globals, nil, nil))
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

func (me *gMachine) compileCtorAppl(comp compilation, ctor *ExprCtor, reverseArgs []IExpr, argsEnv env, fromMark7E bool) code {
	if (!fromMark7E) && len(reverseArgs) != ctor.Arity {
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
		dynglobalcall := Call(Id(dynglobalname), reverseArgs...)
		// println("\n\nCALL:\n\n")
		// src, _ := (&corelang.SyntaxTreePrinter{}).Expr(dynglobalcall)
		// println(src)
		return comp(dynglobalcall, argsEnv)
	}
	instrs := make(code, 0, len(reverseArgs)*2)
	for i, arg := range reverseArgs {
		argsenv := argsEnv
		if !fromMark7E {
			argsenv = me.envOffsetBy(argsEnv, i)
		}
		instrs = append(instrs, comp(arg, argsenv)...)
	}
	return append(instrs, instr{Op: INSTR_CTOR_PACK, Int: ctor.Tag, CtorArity: ctor.Arity})
}

func (me *gMachine) compileLet(compbody compilation, let *ExprLetIn, argsEnv env, finalOp instruction) (instrs code) {
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
	instrs = append(instrs, instr{Op: finalOp, Int: n})
	return
}

func (me *gMachine) compileCaseAlts_SchemeD(compn compilationN, caseAlts []*SynCaseAlt, argsEnv env) (jumpblocks []code) {
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

func (me *gMachine) compilePrimsMaybe(comp compilation, expr *ExprCall, argsEnv env, dyadicOffset int, mark7WrapInB bool) code {
	switch callee := expr.Callee.(type) {
	case *ExprIdent:
		if callee.Name == "neg" {
			if mark7WrapInB {
				return append(me.compileExprStrictMark7_SchemeB(expr, argsEnv), instr{Op: INSTR_MARK7_MAKENODEINT})
			}
			return append(comp(expr.Arg, argsEnv), instr{Op: INSTR_PRIM_AR_NEG})
		}
	case *ExprCall:
		if maybeop, _ := callee.Callee.(*ExprIdent); maybeop != nil {
			if primdyadic := primDyadicsForStrict[maybeop.Name]; primdyadic != 0 {
				if mark7WrapInB {
					finalop := INSTR_MARK7_MAKENODEBOOL
					if primdyadic == INSTR_PRIM_AR_ADD || primdyadic == INSTR_PRIM_AR_SUB || primdyadic == INSTR_PRIM_AR_MUL || primdyadic == INSTR_PRIM_AR_DIV {
						finalop = INSTR_MARK7_MAKENODEINT
					}
					return append(me.compileExprStrictMark7_SchemeB(expr, argsEnv), instr{Op: finalop})
				}
				return append(append(
					comp(expr.Arg, argsEnv),
					comp(callee.Arg, me.envOffsetBy(argsEnv, dyadicOffset))...,
				), instr{Op: primdyadic})
			}
		} else if maybeif, _ := callee.Callee.(*ExprCall); maybeif != nil {
			if ifname, _ := maybeif.Callee.(*ExprIdent); ifname != nil && ifname.Name == "if" {
				cond, condthen, condelse := maybeif.Arg, callee.Arg, expr.Arg
				compcond := comp
				if mark7WrapInB {
					compcond = me.compileExprStrictMark7_SchemeB
				}
				return append(compcond(cond, argsEnv), instr{
					Op:       INSTR_PRIM_COND,
					CondThen: comp(condthen, argsEnv),
					CondElse: comp(condelse, argsEnv),
				})
			}
		}
	}
	return nil
}

func (*gMachine) envOffsetBy(argsEnv env, offsetBy int) (envOffset env) {
	envOffset = make(env, len(argsEnv))
	for k, v := range argsEnv {
		envOffset[k] = v + offsetBy
	}
	return
}
