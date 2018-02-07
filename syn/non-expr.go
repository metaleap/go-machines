package clsyn

type SynMod struct {
	Defs map[string]*SynDef
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
	Tag   uint64
	Binds []string
	Body  IExpr
}
