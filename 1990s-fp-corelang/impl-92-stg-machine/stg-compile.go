package climpl

import (
	"strconv"

	corelang "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

func CompileToMachine(mod *corelang.SynMod) (util.IMachine, []error) {
	me := &stgMachine{}

	modenv := corelang.NewLookupEnv(mod.Defs_(), nil, nil, nil)

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
		return synExprAtomIdent{Name: x.Name}
	case *corelang.ExprLetIn:
		let := synExprLet{Rec: x.Rec, Body: compileExpr(modEnv, prefix, x.Body), Binds: make([]synBinding, len(x.Defs))}
		for i, def := range x.Defs {
			let.Binds[i] = compileBind(modEnv, prefix, def)
		}
		return let
	case *corelang.ExprCtor:
		return synExprCtor{Tag: synExprAtomIdent{Name: strconv.Itoa(x.Tag)}}
	case *corelang.ExprCall:
		var let synExprLet
		call, revargs := x.Flattened()
		if ctor, _ := call.(*corelang.ExprCtor); ctor != nil {
			me := synExprCtor{Tag: synExprAtomIdent{Name: strconv.Itoa(ctor.Tag)}, Args: make([]iSynExprAtom, len(revargs))}
			prefix += me.Tag.Name + "_"
			for i, ctorarg := range revargs {
				if _i := len(me.Args) - (1 + i); ctorarg.IsAtomic() {
					me.Args[_i] = compileExpr(modEnv, prefix, ctorarg).(iSynExprAtom)
				} else {
					name := prefix + strconv.Itoa(i)
					let.Binds = append(let.Binds, compileBind(modEnv, "", &corelang.SynDef{Name: name, Body: ctorarg}))
					me.Args[_i] = synExprAtomIdent{Name: name}
				}
			}
			let.Body = me
		} else {
			me := synExprCall{Args: make([]iSynExprAtom, len(revargs))}
			switch callee := call.(type) {
			case *corelang.ExprIdent:
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
		if len(let.Binds) == 0 {
			return let.Body
		}
		return let
	case *corelang.ExprLambda:
		bind := compileBind(modEnv, "", &corelang.SynDef{Body: x.Body, Args: x.Args, Name: prefix + "lam"})
		return synExprLet{Binds: []synBinding{bind}, Body: synExprAtomIdent{Name: bind.Name}}
	case *corelang.ExprCaseOf:
	}
	return nil
}
