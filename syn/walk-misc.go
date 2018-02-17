package clsyn

import (
	"strconv"
)

func (me *ExprCaseOf) ExtractIntoDef(name string, topLevel bool, lookupEnv map[string]bool) (*SynDef, IExpr) {
	return me.extractIntoDef(me, name, topLevel, lookupEnv)
}

func (me *exprComp) extractIntoDef(this IExpr, name string, topLevel bool, lookupEnv map[string]bool) (*SynDef, IExpr) {
	i, free, def := 0, make(map[string]bool, 8), SynDef{Name: name, TopLevel: topLevel, Body: this}
	this.FreeVars(free, lookupEnv)
	def.toks, def.Args = me.toks, make([]string, len(free))
	for name := range free {
		def.Args[i], i = name, i+1
	}

	if len(def.Args) == 0 {
		return &def, Id(def.Name)
	}
	call := ExprCall{Callee: Id(def.Name), Arg: Id(def.Args[0])}
	for i = 1; i < len(def.Args); i++ {
		call = ExprCall{Callee: &call, Arg: Id(def.Args[i])}
	}
	call.toks = me.toks
	return &def, &call
}

func (me *ExprCtor) ExtractIntoDef(name string, topLevel bool) *SynDef {
	def := SynDef{Name: name, TopLevel: topLevel, Body: me, Args: make([]string, me.Arity)}
	for i := 0; i < me.Arity; i++ {
		def.Args[i] = "v" + strconv.Itoa(i)
		def.Body = Ap(def.Body, Id(def.Args[i]))
	}
	return &def
}

func (me *ExprCall) FlattenedIfEffectivelyCtor() (ctor *ExprCtor, reverseArgs []IExpr) {
	reverseArgs = []IExpr{me.Arg}
	for callee := me.Callee; callee != nil; {
		switch c := callee.(type) {
		case *ExprCtor:
			ctor = c
			return
		case *ExprCall:
			callee, reverseArgs = c.Callee, append(reverseArgs, c.Arg)
		default:
			return nil, nil
		}
	}
	return
}
