package clsyn

func Ap(callee IExpr, arg IExpr) *ExprCall     { return &ExprCall{Callee: callee, Arg: arg} }
func Ab(args []string, body IExpr) *ExprLambda { return &ExprLambda{Args: args, Body: body} }
func Ct(tag uint64, arity uint64) *ExprCtor    { return &ExprCtor{Tag: int(tag), Arity: int(arity)} }

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

func (me *ExprCall) FlattenedIfEffectivelyCtor() (ctor *ExprCtor, reverseArgs []IExpr) {
	reverseArgs = []IExpr{me.Arg}
	for callee := me.Callee; callee != nil; {
		switch c := callee.(type) {
		case *ExprCtor:
			ctor = c
			return
		case *ExprCall:
			callee, reverseArgs = c.Callee, append(reverseArgs, c.Arg)
		default:
			return nil, nil
		}
	}
	return
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
