package climpl

import (
	"errors"
	"strconv"
	"time"

	. "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
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

	if _MARK7 {
		for name, def := range primsMark7Globals {
			if node, err := me.compileGlobal_SchemeSC(def); err != nil {
				errs = append(errs, err)
			} else {
				primsPrecompiledForLazy[name] = node
			}
		}
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
	if _MARK7 {
		node = nodeGlobal{len(global.Args), me.compileExprMark7_SchemeR(len(global.Args))(global.Body, argsenv)}
	} else {
		node = nodeGlobal{len(global.Args), me.compileGlobalBody_SchemeR(global.Body, argsenv)}
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

func (me *gMachine) compileExprMark7_SchemeR(d int) compilation {
	return func(expr IExpr, argsenv env) code {
		return me.compileExprMark7_SchemeR_(expr, argsenv, d)
	}
}

func (me *gMachine) compileExprMark7_SchemeR_(expression IExpr, argsEnv env, d int) code {
	switch expr := expression.(type) {
	case *ExprLetIn:
		return me.compileLet(me.compileExprMark7_SchemeR(d+len(expr.Defs)), expr, argsEnv, 0)
	case *ExprCall:
		if callee, _ := expr.Callee.(*ExprCall); callee != nil {
			if ifcode := me.compilePrimIfMaybe(me.compileExprMark7_SchemeR(d), expr, callee, argsEnv, true); len(ifcode) > 0 {
				return ifcode
			}
		}
	case *ExprCaseOf:
		return append(me.compileExprStrict_SchemeE(expr.Scrut, argsEnv),
			instr{Op: INSTR_CASE_JUMP, CaseJump: me.compileCaseAlts_SchemeD(me.compileExprStrictSplit_SchemeR(d), expr.Alts, argsEnv)})
	}
	return append(me.compileExprStrict_SchemeE(expression, argsEnv),
		instr{Op: INSTR_UPDATE, Int: d}, instr{Op: INSTR_POP, Int: d}, instr{Op: INSTR_UNWIND})
}

func (me *gMachine) compileExprStrict_SchemeE(expression IExpr, argsEnv env) code {
	switch expr := expression.(type) {
	case *ExprLitUInt:
		return code{{Op: INSTR_PUSHINT, Int: int(expr.Lit)}}
	case *ExprLetIn:
		return me.compileLet(me.compileExprStrict_SchemeE, expr, argsEnv, INSTR_SLIDE)
	case *ExprCtor:
		comp := me.compileExprStrict_SchemeE
		if _MARK7 {
			comp = me.compileExprLazy_SchemeC
		}
		return me.compileCtorAppl(comp, expr, nil, argsEnv, _MARK7)
	case *ExprCaseOf:
		return append(me.compileExprStrict_SchemeE(expr.Scrut, argsEnv), instr{Op: INSTR_CASE_JUMP,
			CaseJump: me.compileCaseAlts_SchemeD(me.compileExprStrictSplitSlide_SchemeA, expr.Alts, argsEnv)})
	case *ExprCall:
		if ctor, ctorrevargs := expr.FlattenedIfCtor(); ctor != nil {
			comp := me.compileExprStrict_SchemeE
			if _MARK7 {
				comp = me.compileExprLazy_SchemeC
			}
			return me.compileCtorAppl(comp, ctor, ctorrevargs, argsEnv, _MARK7)
		}

		if instrs := me.compilePrimsMaybe(me.compileExprStrict_SchemeE, expr, argsEnv, 1, _MARK7); len(instrs) > 0 {
			return instrs
		}
	}
	return append(me.compileExprLazy_SchemeC(expression, argsEnv), instr{Op: INSTR_EVAL})
}

func (me *gMachine) compileExprStrictSplitSlide_SchemeA(offset int, expr IExpr, argsEnv env) code {
	return append(me.compileExprStrictSplit(me.compileExprStrict_SchemeE, expr, argsEnv, offset),
		instr{Op: INSTR_SLIDE, Int: offset})
}

func (me *gMachine) compileExprStrictSplit_SchemeR(d int) compilationN {
	return func(offset int, expr IExpr, argsEnv env) code {
		return me.compileExprStrictSplit(me.compileExprMark7_SchemeR(d), expr, argsEnv, offset)
	}
}

func (me *gMachine) compileExprStrictSplit(comp compilation, expr IExpr, argsEnv env, offset int) code {
	return append(code{{Op: INSTR_CASE_SPLIT, Int: offset}}, comp(expr, argsEnv)...)
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
		if ctor, ctorrevargs := expr.FlattenedIfCtor(); ctor != nil {
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
		}
		var dynglobalcall IExpr = Id(dynglobalname)
		if len(reverseArgs) > 0 {
			dynglobalcall = Call(dynglobalcall, reverseArgs...)
		}
		return comp(dynglobalcall, argsEnv)
	}
	instrs := make(code, 0, len(reverseArgs)*3) // arbitrary extra cap, exact need not known
	for i, arg := range reverseArgs {
		instrs = append(instrs, comp(arg, me.envOffsetBy(argsEnv, i))...)
		if fromMark7E {
			instrs = append(instrs, instr{Op: INSTR_EVAL})
		}
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

	if instrs = append(instrs, compbody(let.Body, bodyargsenv)...); finalOp > 0 {
		instrs = append(instrs, instr{Op: finalOp, Int: n})
	}
	return
}

func (me *gMachine) compileCaseAlts_SchemeD(compn compilationN, caseAlts []*SynCaseAlt, argsEnv env) (jumpblocks []code) {
	var tagmax int
	for i := 0; i < len(caseAlts); i++ {
		if caseAlts[i].Tag > tagmax {
			tagmax = caseAlts[i].Tag
		}
	}
	jumpblocks = make([]code, tagmax+1)
	for _, alt := range caseAlts {
		jumpblocks[alt.Tag] = me.compileCaseAlt_SchemeA(compn, alt, argsEnv)
	}
	return
}

func (me *gMachine) compileCaseAlt_SchemeA(compn compilationN, alt *SynCaseAlt, argsEnv env) code {
	n := len(alt.Binds)
	bodyargsenv := me.envOffsetBy(argsEnv, n)
	for i, name := range alt.Binds {
		bodyargsenv[name] = i // = n - (i + 1)
	}
	return compn(n, alt.Body, bodyargsenv)
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
				compcond, cond, condthen, condelse := comp, maybeif.Arg, callee.Arg, expr.Arg
				if mark7WrapInB {
					compcond = me.compileExprStrictMark7_SchemeB
				}
				return me.compilePrimIf(comp, compcond, cond, condthen, condelse, argsEnv)
			}
		} else if ifcode := me.compilePrimIfMaybe(comp, expr, callee, argsEnv, mark7WrapInB); len(ifcode) > 0 {
			return ifcode
		}
	}
	return nil
}

func (me *gMachine) compilePrimIfMaybe(comp compilation, expr *ExprCall, exprCallee *ExprCall, argsEnv env, mark7WrapInB bool) code {
	if maybeif, _ := exprCallee.Callee.(*ExprCall); maybeif != nil {
		if ifname, _ := maybeif.Callee.(*ExprIdent); ifname != nil && ifname.Name == "if" {
			compcond, cond, condthen, condelse := comp, maybeif.Arg, exprCallee.Arg, expr.Arg
			if mark7WrapInB {
				compcond = me.compileExprStrictMark7_SchemeB
			}
			return me.compilePrimIf(comp, compcond, cond, condthen, condelse, argsEnv)
		}
	}
	return nil
}

func (me *gMachine) compilePrimIf(comp compilation, compcond compilation, cond IExpr, condthen IExpr, condelse IExpr, argsEnv env) code {
	return append(compcond(cond, argsEnv), instr{
		Op:       INSTR_PRIM_COND,
		CondThen: comp(condthen, argsEnv),
		CondElse: comp(condelse, argsEnv),
	})
}

func (*gMachine) envOffsetBy(argsEnv env, offsetBy int) (envOffset env) {
	envOffset = make(env, len(argsEnv))
	for k, v := range argsEnv {
		envOffset[k] = v + offsetBy
	}
	return
}
