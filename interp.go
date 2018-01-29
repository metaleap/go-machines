package corelang

type IInterpreter interface {
	Prog(*aProgram, ...interface{}) (interface{}, error)
	Def(*aDef, ...interface{}) (interface{}, error)
	Expr(iExpr) (interface{}, error)
}
