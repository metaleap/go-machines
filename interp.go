package corelang

import (
	. "github.com/metaleap/go-corelang/syn"
)

type IInterpreter interface {
	Mod(*SynMod, ...interface{}) (interface{}, error)
	Def(*SynDef, ...interface{}) (interface{}, error)
	Expr(IExpr) (interface{}, error)
}
