package clsyn

func (this *syn) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {}

func (this *SynDef) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	this.Body.FreeVars(freeVarNames, append(lookupEnvs, NewLookupEnv(nil, nil, nil, this.Args))...)
}

func (this *ExprIdent) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	if !freeVarNames[this.Name] {
		for _, lookupenv := range lookupEnvs {
			if lookupenv[this.Name] {
				return
			}
		}
		freeVarNames[this.Name] = true
	}
}

func (this *ExprCall) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	this.Callee.FreeVars(freeVarNames, lookupEnvs...)
	this.Arg.FreeVars(freeVarNames, lookupEnvs...)
}

func (this *ExprLambda) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	this.Body.FreeVars(freeVarNames, append(lookupEnvs, NewLookupEnv(nil, nil, nil, this.Args))...)
}

func (this *ExprLetIn) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	defsenv := NewLookupEnv(this.Defs, nil, nil, nil)
	combined := append(lookupEnvs, defsenv)
	for _, def := range this.Defs {
		if this.Rec {
			def.FreeVars(freeVarNames, combined...)
		} else {
			def.FreeVars(freeVarNames, lookupEnvs...)
		}
	}
	this.Body.FreeVars(freeVarNames, combined...)
}

func (this *ExprCaseOf) FreeVars(freeVarNames map[string]bool, lookupEnvs ...map[string]bool) {
	this.Scrut.FreeVars(freeVarNames, lookupEnvs...)
	for _, alt := range this.Alts {
		alt.Body.FreeVars(freeVarNames, append(lookupEnvs, NewLookupEnv(nil, nil, nil, alt.Binds))...)
	}
}
