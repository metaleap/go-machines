package climpl

import (
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

type gMachine struct {
}

func CompileToMachine(mod *clsyn.SynMod) clutil.IMachine {
	return &gMachine{}
}

func (me *gMachine) Eval(name string) (val interface{}, numSteps int, err error) {
	defer clutil.Catch(&err)
	return
}
