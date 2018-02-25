package clsyn

type SynMod struct {
	Defs map[string]*SynDef
}

func (me *SynMod) Defs_() (defs []*SynDef) {
	defs = make([]*SynDef, len(me.Defs))
	var i int
	for _, def := range me.Defs {
		defs[i] = def
		i++
	}
	return
}

type SynDef struct {
	syn
	Name     string
	Args     []string
	Body     IExpr
	TopLevel bool
}

type SynCaseAlt struct {
	syn
	Tag   string
	Binds []string
	Body  IExpr
}
