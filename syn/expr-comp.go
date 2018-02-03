package clsyn

func Ap(callee IExpr, arg IExpr) *ExprCall            { return &ExprCall{Callee: callee, Arg: arg} }
func Ab(args []string, body IExpr) *ExprLambda        { return &ExprLambda{Args: args, Body: body} }
func Ct(tag uint64, arity uint64) *ExprCtor           { return &ExprCtor{Tag: tag, Arity: arity} }
func Co(scrut IExpr, alts ...*SynCaseAlt) *ExprCaseOf { return &ExprCaseOf{Scrut: scrut, Alts: alts} }

type ExprCtor struct {
	exprComp
	Tag   uint64
	Arity uint64
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
	Defs []*SynDef
	Body IExpr
}
type ExprCaseOf struct {
	exprComp
	Scrut IExpr
	Alts  []*SynCaseAlt
}
