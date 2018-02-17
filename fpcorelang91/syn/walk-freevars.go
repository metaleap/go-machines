package clsyn

func (me *syn) FreeVars(map[string]bool, ...map[string]bool) {}

func (me *ExprIdent) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	for _, lookupenv := range lookupEnvs {
		if lookupenv[me.Name] {
			return
		}
	}
	freeVarNames[me.Name] = true
}

func (me *ExprCall) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Callee.FreeVars(freeVarNames, lookupEnvs...)
	me.Arg.FreeVars(freeVarNames, lookupEnvs...)
}

func (me *ExprLambda) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Body.FreeVars(freeVarNames, append(lookupEnvs, NewLookupEnv(nil, nil, nil, me.Args))...)
}

func (me *ExprLetIn) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	defsenv := NewLookupEnv(me.Defs, nil, nil, nil)
	combined := append(lookupEnvs, defsenv)
	for _, def := range me.Defs {
		if me.Rec {
			def.Body.FreeVars(freeVarNames, combined...)
		} else {
			def.Body.FreeVars(freeVarNames, lookupEnvs...)
		}
	}
	me.Body.FreeVars(freeVarNames, combined...)
}

func (me *ExprCaseOf) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Scrut.FreeVars(freeVarNames, lookupEnvs...)
	for _, alt := range me.Alts {
		alt.Body.FreeVars(freeVarNames, append(lookupEnvs, NewLookupEnv(nil, nil, nil, alt.Binds))...)
	}
}
