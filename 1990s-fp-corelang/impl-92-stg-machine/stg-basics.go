package climpl

import (
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

type stgMachine struct {
	mod synMod
}

func (me *stgMachine) Eval(argLessDefName string) (result interface{}, stats util.Stats, err error) {
	return
}

func (me *stgMachine) String(result interface{}) string {
	return "?"
}
