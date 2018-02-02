package clsyn

func Alt(tag int, binds []string, body IExpr) *CaseAlt {
	return &CaseAlt{Tag: tag, Binds: binds, Body: body}
}

type CaseAlt struct {
	syn
	Tag   int
	Binds []string
	Body  IExpr
}
