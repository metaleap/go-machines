package corelang

type iExpr interface{}

type exprSymbol struct {
	Name string
}

type exprNumber struct {
	Lit int
}

type exprConstructor struct {
	Tag   uint8
	Arity uint8
}

type exprCall struct {
	Callee iExpr
	Arg    iExpr
}

type exprLetIn struct {
	Rec bool
	Let map[string]iExpr
	In  iExpr
}

type exprCaseOf struct {
	Scrut iExpr
	Alts  map[iExpr]iExpr
}

type exprLambda struct {
	Args []string
	Body iExpr
}
