package climpl

import (
	"strconv"

	corelang "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

func CompileToMachine(mod *corelang.SynMod) (_ util.IMachine, errs []error) {
	me, modenv := &stgMachine{}, corelang.NewLookupEnv(mod.Defs_(), nil, nil, nil)
	me.mod.Binds = make([]synBinding, 0, len(mod.Defs))
	for _, global := range mod.Defs {
		if bind, err := compileCoreGlobalToStg(modenv, global); err != nil {
			errs = append(errs, err)
		} else {
			me.mod.Binds = append(me.mod.Binds, bind)
		}
	}
	println(me.mod.String())
	return me, errs
}

func compileCoreGlobalToStg(modEnv map[string]bool, global *corelang.SynDef) (bind synBinding, err error) {
	defer util.Catch(&err)
	bind = compileCoreDefToStgBind(modEnv, "", global)
	return
}

func compileCoreDefToStgBind(modEnv map[string]bool, prefix string, clDef *corelang.SynDef) (bind synBinding) {
	bind.Name = clDef.Name

	bind.LamForm.Args = make([]synExprAtomIdent, len(clDef.Args))
	for i, argname := range clDef.Args {
		bind.LamForm.Args[i].Name = argname
	}

	freevars := map[string]bool{}
	clDef.FreeVars(freevars, modEnv)
	bind.LamForm.Free = make([]synExprAtomIdent, len(freevars))
	var i int
	for freevarname := range freevars {
		i, bind.LamForm.Free[i].Name = i+1, freevarname
	}

	bind.LamForm.Body = compileCoreExprToStgExpr(modEnv, prefix+clDef.Name+"·", clDef.Body)
	return
}

func compileCoreExprToStgExpr(modEnv map[string]bool, prefix string, clExpr corelang.IExpr) iSynExpr {
	switch x := clExpr.(type) {
	case *corelang.ExprLitFloat:
		return synExprAtomLitFloat{Lit: x.Lit}
	case *corelang.ExprLitUInt:
		return synExprAtomLitUInt{Lit: x.Lit}
	case *corelang.ExprLitText:
		return synExprAtomLitText{Lit: x.Lit}
	case *corelang.ExprLitRune:
		return synExprAtomLitRune{Lit: x.Lit}
	case *corelang.ExprIdent:
		if letbinds, letbody := compileCoreCallToStgPrimOpMaybe(modEnv, prefix, x, nil); letbody != nil {
			return synExprLet{Binds: letbinds, Body: letbody}
		}
		return synExprAtomIdent{Name: x.Name}
	case *corelang.ExprLetIn:
		let := synExprLet{Rec: x.Rec, Body: compileCoreExprToStgExpr(modEnv, prefix, x.Body), Binds: make([]synBinding, len(x.Defs))}
		for i, def := range x.Defs {
			let.Binds[i] = compileCoreDefToStgBind(modEnv, prefix, def)
		}
		return let
	case *corelang.ExprCtor: // not already captured by outer call (see below), so nilary
		return synExprCtor{Tag: synExprAtomIdent{Name: x.Tag}}
	case *corelang.ExprCall:
		var let synExprLet
		call, revargs := x.Flattened()
		if ctor, ok := call.(*corelang.ExprCtor); ok {
			argscap := ctor.Arity
			if argscap < len(revargs) {
				panic("fully-saturated ctor applied like a function")
			}
			me := synExprCtor{Tag: synExprAtomIdent{Name: ctor.Tag}, Args: make([]iSynExprAtom, len(revargs), argscap)}
			prefix += me.Tag.Name + "·"
			for i, ctorarg := range revargs {
				if _i := len(revargs) - (1 + i); ctorarg.IsAtomic() {
					me.Args[_i] = compileCoreExprToStgExpr(modEnv, prefix, ctorarg).(iSynExprAtom)
				} else {
					name := prefix + strconv.Itoa(i)
					let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, "", &corelang.SynDef{Name: name, Body: ctorarg}))
					me.Args[_i] = synExprAtomIdent{Name: name}
				}
			}
			if diff := ctor.Arity - len(revargs); diff > 0 {
				lamdef := synBinding{Name: prefix + "CLAM"}
				lamdef.LamForm.Args = make([]synExprAtomIdent, diff)
				for i := 0; i < diff; i++ {
					lamdef.LamForm.Args[i].Name = lamdef.Name + "·A" + strconv.Itoa(i)
					me.Args = append(me.Args, lamdef.LamForm.Args[i])
				}
				lamdef.LamForm.Body = me
				let.Binds = append(let.Binds, lamdef)
				let.Body = synExprAtomIdent{Name: lamdef.Name}
			} else {
				let.Body = me
			}
		} else {
			me := synExprCall{Args: make([]iSynExprAtom, len(revargs))}
			switch callee := call.(type) {
			case *corelang.ExprIdent:
				if let.Binds, let.Body = compileCoreCallToStgPrimOpMaybe(modEnv, prefix, callee, revargs); let.Body != nil {
					goto retLet
				}
				me.Callee = synExprAtomIdent{Name: callee.Name}
			default:
				me.Callee = synExprAtomIdent{Name: prefix + "CALL"}
				let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, "", &corelang.SynDef{Name: me.Callee.Name, Body: callee}))
			}
			prefix += me.Callee.Name + "·"
			for i, callarg := range revargs {
				if _i := len(me.Args) - (1 + i); callarg.IsAtomic() {
					me.Args[_i] = compileCoreExprToStgExpr(modEnv, prefix, callarg).(iSynExprAtom)
				} else {
					name := prefix + strconv.Itoa(i)
					let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, "", &corelang.SynDef{Name: name, Body: callarg}))
					me.Args[_i] = synExprAtomIdent{Name: name}
				}
			}
			let.Body = me
		}
	retLet:
		if len(let.Binds) == 0 {
			return let.Body
		}
		return let
	case *corelang.ExprLambda:
		bind := compileCoreDefToStgBind(modEnv, "", &corelang.SynDef{Body: x.Body, Args: x.Args, Name: prefix + "LAM"})
		return synExprLet{Binds: []synBinding{bind}, Body: synExprAtomIdent{Name: bind.Name}}
	case *corelang.ExprCaseOf:
		caseof := synExprCaseOf{Scrut: compileCoreExprToStgExpr(modEnv, prefix, x.Scrut), Alts: make([]synCaseAlt, len(x.Alts))}
		for i, alt := range x.Alts {
			caseof.Alts[i].Body = compileCoreExprToStgExpr(modEnv, prefix, alt.Body)
			if alt.Tag != "_" {
				caseof.Alts[i].Ctor.Tag = synExprAtomIdent{Name: alt.Tag}
				caseof.Alts[i].Ctor.Vars = make([]synExprAtomIdent, len(alt.Binds))
				for j, altbind := range alt.Binds {
					caseof.Alts[i].Ctor.Vars[j] = synExprAtomIdent{Name: altbind}
				}
			}
		}
		return caseof
	}
	return nil
}

func compileCoreCallToStgPrimOpMaybe(modEnv map[string]bool, prefix string, callee *corelang.ExprIdent, revArgs []corelang.IExpr) (binds []synBinding, expr iSynExpr) {
	switch callee.Name {
	case "+", "-", "*", "/", "==", "!=", "<=", ">=", ">", "<":
		num, op := len(revArgs), &synExprPrimOp{PrimOp: callee.Name}
		if expr, prefix = op, prefix+callee.Name+"·"; num > 2 {
			panic("prim-op `" + op.PrimOp + "` over-saturated: expected 2 operands, not " + strconv.Itoa(num))
		}
		if num > 0 {
			left := compileCoreExprToStgExpr(modEnv, prefix, revArgs[num-1])
			if op.Left, _ = left.(iSynExprAtom); op.Left == nil {
				bind := synBinding{Name: prefix + "¯L"}
				bind.LamForm.Body = left
				binds, op.Left = append(binds, bind), synExprAtomIdent{Name: bind.Name}
			}
		}
		if num > 1 {
			right := compileCoreExprToStgExpr(modEnv, prefix, revArgs[0])
			if op.Right, _ = right.(iSynExprAtom); op.Right == nil {
				bind := synBinding{Name: prefix + "¯R"}
				bind.LamForm.Body = right
				binds, op.Right = append(binds, bind), synExprAtomIdent{Name: bind.Name}
			}
		}
		if op.Left == nil || op.Right == nil { // num<2
			lamdef := synBinding{Name: prefix + "PLAM"}
			lamdef.LamForm.Body = op
			if op.Left == nil { // num==0
				name := synExprAtomIdent{Name: lamdef.Name + "¨L"}
				op.Left, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			if op.Right == nil { // num<=1
				name := synExprAtomIdent{Name: lamdef.Name + "¨R"}
				op.Right, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			binds = append(binds, lamdef)
			expr = synExprAtomIdent{Name: lamdef.Name}
		}
	}
	return
}
