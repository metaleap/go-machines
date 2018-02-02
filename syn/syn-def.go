package clsyn

type Def struct {
	syn
	Name string
	Args []string
	Body IExpr
}
