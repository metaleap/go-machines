package clsyn

func Ap(callee IExpr, arg IExpr) *ExprCall         { return &ExprCall{Callee: callee, Arg: arg} }
func Ab(args []string, body IExpr) *ExprLambda     { return &ExprLambda{Args: args, Body: body} }
func Ct(tag int, arity int) *ExprCtor              { return &ExprCtor{Tag: tag, Arity: arity} }
func Co(scrut IExpr, alts ...*CaseAlt) *ExprCaseOf { return &ExprCaseOf{Scrut: scrut, Alts: alts} }

type ExprCtor struct {
	exprComp
	Tag   int
	Arity int
}

type ExprCall struct {
	exprComp
	Callee IExpr
	Arg    IExpr
}

type ExprLambda struct {
	exprComp
	Args []string
	Body IExpr
}

type ExprLetIn struct {
	exprComp
	Rec  bool
	Defs []*Def
	Body IExpr
}

type ExprCaseOf struct {
	exprComp
	Scrut IExpr
	Alts  []*CaseAlt
}
