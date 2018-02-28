package climpl

import (
	"strconv"

	corelang "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

var (
	identCounter = 'a' - 1
	identSuff    = "a"
)

func CompileToMachine(mod *corelang.SynMod) (_ util.IMachine, errs []error) {
	me, modenv := &stgMachine{}, corelang.NewLookupEnv(mod.Defs, nil, nil, nil)
	me.mod.Binds = make(synBindings, 0, len(mod.Defs))
	for _, global := range mod.Defs {
		if bind, err := compileCoreGlobalToStg(modenv, global); err != nil {
			errs = append(errs, err)
		} else {
			me.mod.Binds = append(me.mod.Binds, bind)
		}
	}
	me.mod.setUpd()
	// println(me.mod.String())
	return me, errs
}

func nextIdent() string {
	if identCounter++; identCounter >= 'z' {
		if identCounter = 'a'; identSuff[0] == 'z' {
			identSuff = "a" + identSuff
		} else {
			identSuff = string(identSuff[0]+1) + identSuff[1:]
		}
	}
	return string(identCounter) + identSuff + "’"
}

func compileCoreGlobalToStg(modEnv map[string]bool, global *corelang.SynDef) (bind *synBinding, err error) {
	defer util.Catch(&err)
	for _, argname := range global.Args {
		if _, argnameshadowsglobal := modEnv[argname]; argnameshadowsglobal {
			modEnv[argname] = false // hacky hiding of shadowed globals: for compileCtorApplMaybeLambda()
		}
	}
	bind = compileCoreDefToStgBind(modEnv, global)
	for globalname := range modEnv {
		modEnv[globalname] = true // undo above's hack..
	}
	return
}

func compileCoreDefToStgBind(modEnv map[string]bool, clDef *corelang.SynDef) (bind *synBinding) {
	bind = &synBinding{Name: clDef.Name}
	bind.LamForm.Args = make([]*synExprAtomIdent, len(clDef.Args))
	for i, argname := range clDef.Args {
		bind.LamForm.Args[i] = &synExprAtomIdent{Name: argname}
	}

	freevars := map[string]bool{}
	clDef.FreeVars(freevars, modEnv)
	bind.LamForm.Free = make([]*synExprAtomIdent, len(freevars))
	var i int
	for freevarname := range freevars {
		i, bind.LamForm.Free[i] = i+1, &synExprAtomIdent{Name: freevarname}
	}

	bind.LamForm.Body = compileCoreExprToStgExpr(modEnv, clDef.Body)
	return
}

func compileCoreExprToStgExpr(modEnv map[string]bool, clExpr corelang.IExpr) iSynExpr {
	switch x := clExpr.(type) {
	case *corelang.ExprLitFloat:
		return &synExprAtomLitFloat{Lit: x.Lit}
	case *corelang.ExprLitUInt:
		return &synExprAtomLitUInt{Lit: x.Lit}
	case *corelang.ExprLitText:
		return &synExprAtomLitText{Lit: x.Lit}
	case *corelang.ExprLitRune:
		return &synExprAtomLitRune{Lit: x.Lit}
	case *corelang.ExprIdent:
		if letbinds, letbody := compileCoreCallToStgPrimOpMaybe(modEnv, x, nil); letbody != nil {
			return &synExprLet{Binds: letbinds, Body: letbody}
		}
		return &synExprAtomIdent{Name: x.Name}
	case *corelang.ExprLetIn:
		let := &synExprLet{Rec: x.Rec, Body: compileCoreExprToStgExpr(modEnv, x.Body), Binds: make(synBindings, len(x.Defs))}
		for i, def := range x.Defs {
			let.Binds[i] = compileCoreDefToStgBind(modEnv, def)
		}
		return let
	case *corelang.ExprCtor: // not already captured by outer call (see below), so nilary application
		var let synExprLet
		me := synExprCtor{Tag: &synExprAtomIdent{Name: x.Tag}}
		compileCtorApplMaybeLambda(modEnv, x.Arity, x, &me, &let)
		if len(let.Binds) == 0 {
			return &me
		}
		return &let
	case *corelang.ExprCall:
		var let synExprLet
		call, revargs := x.Flattened()
		if ctor, ok := call.(*corelang.ExprCtor); ok {
			argscap := ctor.Arity
			if argscap < len(revargs) {
				panic("fully-saturated ctor applied like a function")
			}
			me := synExprCtor{Tag: &synExprAtomIdent{Name: ctor.Tag}, Args: make([]iSynExprAtom, len(revargs), argscap)}
			for i, ctorarg := range revargs {
				if _i := len(revargs) - (1 + i); ctorarg.IsAtomic() {
					me.Args[_i] = compileCoreExprToStgExpr(modEnv, ctorarg).(iSynExprAtom)
				} else {
					name := nextIdent()
					let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, &corelang.SynDef{Name: name, Body: ctorarg}))
					me.Args[_i] = &synExprAtomIdent{Name: name}
				}
			}
			let.Body = &me
			compileCtorApplMaybeLambda(modEnv, ctor.Arity-len(revargs), ctor, &me, &let)
		} else {
			me := synExprCall{Args: make([]iSynExprAtom, len(revargs))}
			switch callee := call.(type) {
			case *corelang.ExprIdent:
				if let.Binds, let.Body = compileCoreCallToStgPrimOpMaybe(modEnv, callee, revargs); let.Body != nil {
					goto retLet
				}
				me.Callee = &synExprAtomIdent{Name: callee.Name}
			default:
				me.Callee = &synExprAtomIdent{Name: "CALL" + nextIdent()}
				let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, &corelang.SynDef{Name: me.Callee.Name, Body: callee}))
			}
			for i, callarg := range revargs {
				if _i := len(me.Args) - (1 + i); callarg.IsAtomic() {
					me.Args[_i] = compileCoreExprToStgExpr(modEnv, callarg).(iSynExprAtom)
				} else {
					name := nextIdent()
					let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, &corelang.SynDef{Name: name, Body: callarg}))
					me.Args[_i] = &synExprAtomIdent{Name: name}
				}
			}
			let.Body = &me
		}
	retLet:
		if len(let.Binds) == 0 {
			return let.Body
		}
		return &let
	case *corelang.ExprLambda:
		bind := compileCoreDefToStgBind(modEnv, &corelang.SynDef{Body: x.Body, Args: x.Args, Name: nextIdent()})
		return &synExprLet{Body: &synExprAtomIdent{Name: bind.Name}, Binds: synBindings{bind}}
	case *corelang.ExprCaseOf:
		caseof := synExprCaseOf{Scrut: compileCoreExprToStgExpr(modEnv, x.Scrut), Alts: make([]*synCaseAlt, len(x.Alts))}
		for i, alt := range x.Alts {
			caseof.Alts[i] = &synCaseAlt{Body: compileCoreExprToStgExpr(modEnv, alt.Body)}
			if alt.Tag != "_" {
				caseof.Alts[i].Ctor.Tag = &synExprAtomIdent{Name: alt.Tag}
				caseof.Alts[i].Ctor.Vars = make([]*synExprAtomIdent, len(alt.Binds))
				for j, altbind := range alt.Binds {
					caseof.Alts[i].Ctor.Vars[j] = &synExprAtomIdent{Name: altbind}
				}
			}
		}
		return &caseof
	}
	return nil
}

func compileCtorApplMaybeLambda(modEnv map[string]bool, diff int, expr *corelang.ExprCtor, me *synExprCtor, let *synExprLet) {
	if diff > 0 {
		i, fv, lamdef := 0, map[string]bool{}, synBinding{Name: nextIdent()}
		expr.FreeVars(fv, modEnv)
		lamdef.LamForm.Free = make([]*synExprAtomIdent, len(fv))
		for k := range fv {
			i, lamdef.LamForm.Free[i].Name = i+1, k
		}
		lamdef.LamForm.Args = make([]*synExprAtomIdent, diff)
		for i = 0; i < diff; i++ {
			lamdef.LamForm.Args[i] = &synExprAtomIdent{Name: nextIdent()}
			me.Args = append(me.Args, lamdef.LamForm.Args[i])
		}
		lamdef.LamForm.Body = me
		let.Binds = append(let.Binds, &lamdef)
		let.Body = &synExprAtomIdent{Name: lamdef.Name}
	}
}

func compileCoreCallToStgPrimOpMaybe(modEnv map[string]bool, callee *corelang.ExprIdent, revArgs []corelang.IExpr) (binds synBindings, expr iSynExpr) {
	switch callee.Name {
	case "+", "-", "*", "/", "==", "!=", "<=", ">=", ">", "<":
		num, op := len(revArgs), &synExprPrimOp{PrimOp: callee.Name}
		if expr = op; num > 2 {
			panic("prim-op `" + op.PrimOp + "` over-saturated: expected 2 operands, not " + strconv.Itoa(num))
		}
		if num > 0 {
			left := compileCoreExprToStgExpr(modEnv, revArgs[num-1])
			if op.Left, _ = left.(iSynExprAtom); op.Left == nil {
				bind := synBinding{Name: nextIdent()}
				bind.LamForm.Body = left
				binds, op.Left = append(binds, &bind), &synExprAtomIdent{Name: bind.Name}
			}
		}
		if num > 1 {
			right := compileCoreExprToStgExpr(modEnv, revArgs[0])
			if op.Right, _ = right.(iSynExprAtom); op.Right == nil {
				bind := synBinding{Name: nextIdent()}
				bind.LamForm.Body = right
				binds, op.Right = append(binds, &bind), &synExprAtomIdent{Name: bind.Name}
			}
		}
		if op.Left == nil || op.Right == nil { // num<2
			lamdef := synBinding{Name: nextIdent()}
			lamdef.LamForm.Body = op
			if op.Left == nil { // num==0
				name := &synExprAtomIdent{Name: nextIdent()}
				op.Left, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			if op.Right == nil { // num<=1
				name := &synExprAtomIdent{Name: nextIdent()}
				op.Right, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			binds = append(binds, &lamdef)
			expr = &synExprAtomIdent{Name: lamdef.Name}
		}
	}
	return
}
