package climpl

import (
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

type stgMachine struct {
	mod synMod
}

func (me *stgMachine) Eval(argLessDefName string) (result interface{}, stats util.Stats, err error) {
	for i := range me.mod.Binds {
		if me.mod.Binds[i].Name == argLessDefName {
			result = me.mod.Binds[i]
			return
		}
	}
	return
}

func (me *stgMachine) String(result interface{}) string {
	if syn, _ := result.(iSyn); syn != nil {
		return syn.String()
	}
	return "?"
}
