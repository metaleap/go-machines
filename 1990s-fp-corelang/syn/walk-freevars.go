package clsyn

func (me *syn) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {}

func (me *SynDef) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Body.FreeVars(freeVarNames, append(append([]map[string]bool{}, lookupEnvs...), NewLookupEnv(nil, nil, nil, me.Args))...)
}

func (me *ExprIdent) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	if !freeVarNames[me.Name] {
		for _, lookupenv := range lookupEnvs {
			if lookupenv[me.Name] {
				return
			}
		}
		freeVarNames[me.Name] = true
	}
}

func (me *ExprCall) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Callee.FreeVars(freeVarNames, lookupEnvs...)
	me.Arg.FreeVars(freeVarNames, lookupEnvs...)
}

func (me *ExprLambda) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Body.FreeVars(freeVarNames, append(append([]map[string]bool{}, lookupEnvs...), NewLookupEnv(nil, nil, nil, me.Args))...)
}

func (me *ExprLetIn) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	defsenv := NewLookupEnv(me.Defs, nil, nil, nil)
	combined := append(append([]map[string]bool{}, lookupEnvs...), defsenv)
	for _, def := range me.Defs {
		if me.Rec {
			def.FreeVars(freeVarNames, combined...)
		} else {
			def.FreeVars(freeVarNames, lookupEnvs...)
		}
	}
	me.Body.FreeVars(freeVarNames, combined...)
}

func (me *ExprCaseOf) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	me.Scrut.FreeVars(freeVarNames, lookupEnvs...)
	for _, alt := range me.Alts {
		alt.Body.FreeVars(freeVarNames, append(append([]map[string]bool{}, lookupEnvs...), NewLookupEnv(nil, nil, nil, alt.Binds))...)
	}
}
