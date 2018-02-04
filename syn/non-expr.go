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
	Tag   uint64
	Binds []string
	Body  IExpr
}

func Alt(tag uint64, binds []string, body IExpr) *SynCaseAlt {
	return &SynCaseAlt{Tag: tag, Binds: binds, Body: body}
}
