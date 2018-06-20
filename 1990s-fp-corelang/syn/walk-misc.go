package clsyn

import (
	"strconv"
)

func (this *ExprCaseOf) ExtractIntoDef(name string, topLevel bool, lookupEnv map[string]bool) (*SynDef, IExpr) {
	return this.extractIntoDef(this, name, topLevel, lookupEnv)
}

func (this *exprComp) extractIntoDef(self IExpr, name string, topLevel bool, lookupEnv map[string]bool) (*SynDef, IExpr) {
	i, free, def := 0, make(map[string]bool, 8), SynDef{Name: name, TopLevel: topLevel, Body: self}
	self.FreeVars(free, lookupEnv)
	def.toks, def.Args = this.toks, make([]string, len(free))
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
	call.toks = this.toks
	return &def, &call
}
func (this *ExprCtor) ExtractIntoDef(name string, topLevel bool) *SynDef {
	def := SynDef{Name: name, TopLevel: topLevel, Body: this, Args: make([]string, this.Arity)}
	for i := 0; i < this.Arity; i++ {
		def.Args[i] = "v" + strconv.Itoa(i)
		def.Body = Ap(def.Body, Id(def.Args[i]))
	}
	return &def
}

func (this *ExprCall) Flattened() (callee IExpr, reverseArgs []IExpr) {
	for ; this != nil; this, _ = callee.(*ExprCall) {
		callee, reverseArgs = this.Callee, append(reverseArgs, this.Arg)
	}
	return
}

func (this *ExprCall) FlattenedIfCtor() (ctor *ExprCtor, reverseArgs []IExpr) {
	callee, revargs := this.Flattened()
	if ctormaybe, ok := callee.(*ExprCtor); ok {
		ctor, reverseArgs = ctormaybe, revargs
	}
	return
}
