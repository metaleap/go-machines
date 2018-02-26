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
	me, modenv := &stgMachine{}, corelang.NewLookupEnv(mod.Defs_(), nil, nil, nil)
	me.mod.Binds = make([]synBinding, 0, len(mod.Defs))
	for _, global := range mod.Defs {
		if bind, err := compileCoreGlobalToStg(modenv, global); err != nil {
			errs = append(errs, err)
		} else {
			me.mod.Binds = append(me.mod.Binds, bind)
		}
	}
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

func compileCoreGlobalToStg(modEnv map[string]bool, global *corelang.SynDef) (bind synBinding, err error) {
	defer util.Catch(&err)

	for _, argname := range global.Args {
		modEnv[argname] = false // hacky hiding of shadowed globals: for compileCoreExprToStgExpr/ExprCall/ExprCtor/arityDiff case way below
	}
	bind = compileCoreDefToStgBind(modEnv, global)
	for globalname := range modEnv {
		modEnv[globalname] = true // undo above's hack..
	}
	return
}

func compileCoreDefToStgBind(modEnv map[string]bool, clDef *corelang.SynDef) (bind synBinding) {
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

	bind.LamForm.Body = compileCoreExprToStgExpr(modEnv, clDef.Body)
	return
}

func compileCoreExprToStgExpr(modEnv map[string]bool, clExpr corelang.IExpr) iSynExpr {
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
		if letbinds, letbody := compileCoreCallToStgPrimOpMaybe(modEnv, x, nil); letbody != nil {
			return synExprLet{Binds: letbinds, Body: letbody}
		}
		return synExprAtomIdent{Name: x.Name}
	case *corelang.ExprLetIn:
		let := synExprLet{Rec: x.Rec, Body: compileCoreExprToStgExpr(modEnv, x.Body), Binds: make([]synBinding, len(x.Defs))}
		for i, def := range x.Defs {
			let.Binds[i] = compileCoreDefToStgBind(modEnv, def)
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
			for i, ctorarg := range revargs {
				if _i := len(revargs) - (1 + i); ctorarg.IsAtomic() {
					me.Args[_i] = compileCoreExprToStgExpr(modEnv, ctorarg).(iSynExprAtom)
				} else {
					name := nextIdent()
					let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, &corelang.SynDef{Name: name, Body: ctorarg}))
					me.Args[_i] = synExprAtomIdent{Name: name}
				}
			}
			if diff := ctor.Arity - len(revargs); diff > 0 {
				i, fv, lamdef := 0, map[string]bool{}, synBinding{Name: nextIdent()}
				x.FreeVars(fv, modEnv)
				lamdef.LamForm.Free = make([]synExprAtomIdent, len(fv))
				for k := range fv {
					i, lamdef.LamForm.Free[i].Name = i+1, k
				}
				lamdef.LamForm.Args = make([]synExprAtomIdent, diff)
				for i = 0; i < diff; i++ {
					lamdef.LamForm.Args[i].Name = nextIdent()
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
				if let.Binds, let.Body = compileCoreCallToStgPrimOpMaybe(modEnv, callee, revargs); let.Body != nil {
					goto retLet
				}
				me.Callee = synExprAtomIdent{Name: callee.Name}
			default:
				me.Callee = synExprAtomIdent{Name: "CALL" + nextIdent()}
				let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, &corelang.SynDef{Name: me.Callee.Name, Body: callee}))
			}
			for i, callarg := range revargs {
				if _i := len(me.Args) - (1 + i); callarg.IsAtomic() {
					me.Args[_i] = compileCoreExprToStgExpr(modEnv, callarg).(iSynExprAtom)
				} else {
					name := nextIdent()
					let.Binds = append(let.Binds, compileCoreDefToStgBind(modEnv, &corelang.SynDef{Name: name, Body: callarg}))
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
		bind := compileCoreDefToStgBind(modEnv, &corelang.SynDef{Body: x.Body, Args: x.Args, Name: nextIdent()})
		return synExprLet{Binds: []synBinding{bind}, Body: synExprAtomIdent{Name: bind.Name}}
	case *corelang.ExprCaseOf:
		caseof := synExprCaseOf{Scrut: compileCoreExprToStgExpr(modEnv, x.Scrut), Alts: make([]synCaseAlt, len(x.Alts))}
		for i, alt := range x.Alts {
			caseof.Alts[i].Body = compileCoreExprToStgExpr(modEnv, alt.Body)
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

func compileCoreCallToStgPrimOpMaybe(modEnv map[string]bool, callee *corelang.ExprIdent, revArgs []corelang.IExpr) (binds []synBinding, expr iSynExpr) {
	switch callee.Name {
	case "+", "-", "*", "/", "==", "!=", "<=", ">=", ">", "<":
		num, op := len(revArgs), &synExprPrimOp{PrimOp: callee.Name}
		if expr = op; num > 2 {
			panic("prim-op `" + op.PrimOp + "` over-saturated: expected 2 operands, not " + strconv.Itoa(num))
		}
		if num > 0 {
			left := compileCoreExprToStgExpr(modEnv, revArgs[num-1])
			if op.Left, _ = left.(iSynExprAtom); op.Left == nil {
				bind := synBinding{Name: "L¯" + nextIdent()}
				bind.LamForm.Body = left
				binds, op.Left = append(binds, bind), synExprAtomIdent{Name: bind.Name}
			}
		}
		if num > 1 {
			right := compileCoreExprToStgExpr(modEnv, revArgs[0])
			if op.Right, _ = right.(iSynExprAtom); op.Right == nil {
				bind := synBinding{Name: "R¯" + nextIdent()}
				bind.LamForm.Body = right
				binds, op.Right = append(binds, bind), synExprAtomIdent{Name: bind.Name}
			}
		}
		if op.Left == nil || op.Right == nil { // num<2
			lamdef := synBinding{Name: nextIdent()}
			lamdef.LamForm.Body = op
			if op.Left == nil { // num==0
				name := synExprAtomIdent{Name: "L¨" + nextIdent()}
				op.Left, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			if op.Right == nil { // num<=1
				name := synExprAtomIdent{Name: "R¨" + nextIdent()}
				op.Right, lamdef.LamForm.Args = name, append(lamdef.LamForm.Args, name)
			}
			binds = append(binds, lamdef)
			expr = synExprAtomIdent{Name: lamdef.Name}
		}
	}
	return
}
