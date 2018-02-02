package corelang

import (
	. "github.com/metaleap/go-corelang/syn"
)

type IInterpreter interface {
	Mod(*Module, ...interface{}) (interface{}, error)
	Def(*Def, ...interface{}) (interface{}, error)
	Expr(IExpr) (interface{}, error)
}
