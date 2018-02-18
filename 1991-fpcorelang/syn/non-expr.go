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
	Tag   int
	Binds []string
	Body  IExpr
}
