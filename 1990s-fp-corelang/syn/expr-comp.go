package clsyn

func Ap(callee IExpr, arg IExpr) *ExprCall     { return &ExprCall{Callee: callee, Arg: arg} }
func Ab(args []string, body IExpr) *ExprLambda { return &ExprLambda{Args: args, Body: body} }
func Ct(tag string, arity uint64) *ExprCtor    { return &ExprCtor{Tag: tag, Arity: int(arity)} }

func Call(callee IExpr /*argsReversed bool,*/, args ...IExpr) (call *ExprCall) {
	var i int
	// if argsReversed {
	i = len(args) - 1
	// }
	call = Ap(callee, args[i])
	// if argsReversed {
	for i = i - 1; i >= 0; i-- {
		call = Ap(call, args[i])
	}
	// } else {
	// 	for i = 1; i < len(args); i++ {
	// 		call = Ap(call, args[i])
	// 	}
	// }
	return
}

type ExprCtor struct {
	exprComp
	Tag   string
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
	Defs []*SynDef
	Body IExpr
}

type ExprCaseOf struct {
	exprComp
	Scrut IExpr
	Alts  []*SynCaseAlt
}
