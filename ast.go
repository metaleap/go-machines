package corelang

type aProgram struct {
	Defs []*aDef
}

type aDef struct {
	Name string
	Args []string
	Body iExpr
}

type iExpr interface {
	isAtomic() bool
}

func aSym(name string) iExpr              { return &aExprSym{Name: name} }
func aNum(lit int) iExpr                  { return &aExprNum{Lit: lit} }
func aCall(callee iExpr, arg iExpr) iExpr { return &aExprCall{Callee: callee, Arg: arg} }

type aExpr struct {
}

func (me *aExpr) isAtomic() bool { return false }

type aExprSym struct {
	aExpr
	Name string
}

func (me *aExprSym) isAtomic() bool { return true }

type aExprNum struct {
	aExpr
	Lit int
}

func (me *aExprNum) isAtomic() bool { return true }

type aExprCtor struct {
	aExpr
	Tag   uint8
	Arity uint8
}

type aExprCall struct {
	aExpr
	Callee iExpr
	Arg    iExpr
}

type aExprLet struct {
	aExpr
	Rec bool
	Let map[string]iExpr
	In  iExpr
}

type aExprCase struct {
	aExpr
	Scrut iExpr
	Alts  []*aExprCaseAlt
}

type aExprCaseAlt struct {
	aExpr
	Tag   int
	Binds []string
	Body  iExpr
}

type aExprLambda struct {
	aExpr
	Args []string
	Body iExpr
}
