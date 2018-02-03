package clsyn

type SynMod struct {
	Defs map[string]*SynDef
}

type SynDef struct {
	syn
	Name string
	Args []string
	Body IExpr
}

type SynCaseAlt struct {
	syn
	Tag   int
	Binds []string
	Body  IExpr
}

func Alt(tag int, binds []string, body IExpr) *SynCaseAlt {
	return &SynCaseAlt{Tag: tag, Binds: binds, Body: body}
}
