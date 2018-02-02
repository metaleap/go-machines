package clsyn

func Alt(tag int, binds []string, body IExpr) *SynCaseAlt {
	return &SynCaseAlt{Tag: tag, Binds: binds, Body: body}
}

type SynCaseAlt struct {
	syn
	Tag   int
	Binds []string
	Body  IExpr
}
