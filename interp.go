package corelang

type iInterpreter interface {
	prog(*aProgram, ...interface{}) (interface{}, error)
	def(*aDef, ...interface{}) (interface{}, error)
	expr(iExpr, ...interface{}) (interface{}, error)
}
