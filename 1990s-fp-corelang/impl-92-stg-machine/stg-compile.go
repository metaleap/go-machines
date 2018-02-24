package climpl

import (
	"strconv"

	corelang "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

func CompileToMachine(mod *corelang.SynMod) (util.IMachine, []error) {
	me, modenv := &stgMachine{}, corelang.NewLookupEnv(mod.Defs_(), nil, nil, nil)
	for _, global := range mod.Defs {
		me.mod.Binds = append(me.mod.Binds, compileBind(modenv, "", global))
	}
	return me, nil
}

func compileBind(modEnv map[string]bool, prefix string, clDef *corelang.SynDef) (bind synBinding) {
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

	bind.LamForm.Body = compileExpr(modEnv, prefix+clDef.Name+"_", clDef.Body)
	return
}

func compileExpr(modEnv map[string]bool, prefix string, clExpr corelang.IExpr) iSynExpr {
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
		if letbinds, letbody := compilePrimOpMaybe(modEnv, prefix, x, nil); letbody != nil {
			return synExprLet{Binds: letbinds, Body: letbody}
		}
		return synExprAtomIdent{Name: x.Name}
	case *corelang.ExprLetIn:
		let := synExprLet{Rec: x.Rec, Body: compileExpr(modEnv, prefix, x.Body), Binds: make([]synBinding, len(x.Defs))}
		for i, def := range x.Defs {
			let.Binds[i] = compileBind(modEnv, prefix, def)
		}
		return let
	case *corelang.ExprCtor: // not already captured by outer call (see below), so nilary
		return synExprCtor{Tag: synExprAtomIdent{Name: strconv.Itoa(x.Tag)}}
	case *corelang.ExprCall:
		var let synExprLet
		call, revargs := x.Flattened()
		if ctor, ok := call.(*corelang.ExprCtor); ok {
			me := synExprCtor{Tag: synExprAtomIdent{Name: strconv.Itoa(ctor.Tag)}, Args: make([]iSynExprAtom, ctor.Arity)}
			prefix += me.Tag.Name + "_"
			for i, ctorarg := range revargs {
				if _i := len(revargs) - (1 + i); ctorarg.IsAtomic() {
					me.Args[_i] = compileExpr(modEnv, prefix, ctorarg).(iSynExprAtom)
				} else {
					name := prefix + strconv.Itoa(i)
					let.Binds = append(let.Binds, compileBind(modEnv, "", &corelang.SynDef{Name: name, Body: ctorarg}))
					me.Args[_i] = synExprAtomIdent{Name: name}
				}
			}
			let.Body = me
			if diff := ctor.Arity - len(revargs); diff < 0 {
				panic("fully-saturated ctor applied like a function")
			} else if diff > 0 {
				lamdef := synBinding{Name: prefix + "lam"}
				lamdef.LamForm.Body, lamdef.LamForm.Args = me, make([]synExprAtomIdent, diff)
				for i := 0; i < diff; i++ {
					lamdef.LamForm.Args[i].Name = lamdef.Name + "_a_" + strconv.Itoa(i)
					me.Args = append(me.Args, lamdef.LamForm.Args[i])
				}
				let.Binds = append(let.Binds, lamdef)
				let.Body = synExprAtomIdent{Name: lamdef.Name}
			}
		} else {
			me := synExprCall{Args: make([]iSynExprAtom, len(revargs))}
			switch callee := call.(type) {
			case *corelang.ExprIdent:
				if let.Binds, let.Body = compilePrimOpMaybe(modEnv, prefix, callee, revargs); let.Body != nil {
					goto retLet
				}
				me.Callee = synExprAtomIdent{Name: callee.Name}
			default:
				me.Callee = synExprAtomIdent{Name: prefix + "callee"}
				let.Binds = append(let.Binds, compileBind(modEnv, "", &corelang.SynDef{Name: me.Callee.Name, Body: callee}))
			}
			prefix += me.Callee.Name + "_"
			for i, callarg := range revargs {
				if _i := len(me.Args) - (1 + i); callarg.IsAtomic() {
					me.Args[_i] = compileExpr(modEnv, prefix, callarg).(iSynExprAtom)
				} else {
					name := prefix + strconv.Itoa(i)
					let.Binds = append(let.Binds, compileBind(modEnv, "", &corelang.SynDef{Name: name, Body: callarg}))
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
		bind := compileBind(modEnv, "", &corelang.SynDef{Body: x.Body, Args: x.Args, Name: prefix + "lam"})
		return synExprLet{Binds: []synBinding{bind}, Body: synExprAtomIdent{Name: bind.Name}}
	case *corelang.ExprCaseOf:
		caseof := synExprCaseOf{Scrut: compileExpr(modEnv, prefix, x.Scrut), Alts: make([]synCaseAlt, len(x.Alts))}
		for i, alt := range x.Alts {
			caseof.Alts[i].Body = compileExpr(modEnv, prefix, alt.Body)
			if alt.Tag > 0 {
				caseof.Alts[i].Ctor.Tag = synExprAtomIdent{Name: strconv.Itoa(alt.Tag)}
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

func compilePrimOpMaybe(modEnv map[string]bool, prefix string, callee *corelang.ExprIdent, revArgs []corelang.IExpr) (binds []synBinding, expr iSynExpr) {
	switch callee.Name {
	case "+", "-", "*", "/", "==", "!=", "<=", ">=", ">", "<":
		num, op := len(revArgs), &synExprPrimOp{PrimOp: callee.Name}
		if expr, prefix = op, prefix+callee.Name+"_"; num > 2 {
			panic("prim-op `" + op.PrimOp + "` over-saturated: expecting 2 operands, not " + strconv.Itoa(num))
		}
		if num > 0 {
			left := compileExpr(modEnv, prefix, revArgs[num-1])
			if op.Left, _ = left.(iSynExprAtom); op.Left == nil {
				bind := synBinding{Name: prefix + "l"}
				bind.LamForm.Body = left
				binds, op.Left = append(binds, bind), synExprAtomIdent{Name: bind.Name}
			}
		}
		if num > 1 {
			right := compileExpr(modEnv, prefix, revArgs[0])
			if op.Right, _ = right.(iSynExprAtom); op.Right == nil {
				bind := synBinding{Name: prefix + "r"}
				bind.LamForm.Body = right
				binds, op.Right = append(binds, bind), synExprAtomIdent{Name: bind.Name}
			}
		}
		if op.Left == nil || op.Right == nil { // num<2
			lamdef := synBinding{Name: prefix + "lam"}
			lamdef.LamForm.Body = op
			if op.Left == nil { // num==0
				name := synExprAtomIdent{Name: lamdef.Name + "_l"}
				op.Left, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			if op.Right == nil { // num<=1
				name := synExprAtomIdent{Name: lamdef.Name + "_r"}
				op.Right, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			binds = append(binds, lamdef)
			expr = synExprAtomIdent{Name: lamdef.Name}
		}
	}
	return
}
