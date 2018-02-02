package clsyn

type SynDef struct {
	syn
	Name string
	Args []string
	Body IExpr
}
