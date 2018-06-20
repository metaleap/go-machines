package climpl

import (
	util "github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

type stgMachine struct {
	mod synMod
}

func (this *stgMachine) Eval(argLessDefName string) (result interface{}, stats util.Stats, err error) {
	for i := range this.mod.Binds {
		if this.mod.Binds[i].Name == argLessDefName {
			result = this.mod.Binds[i]
			return
		}
	}
	return
}

func (this *stgMachine) String(result interface{}) string {
	if syn, _ := result.(iSyn); syn != nil {
		return syn.String()
	}
	return "?"
}
