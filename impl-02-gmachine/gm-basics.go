package climpl

import (
	"github.com/metaleap/go-corelang/util"
)

type gMachine struct {
	Stack   clutil.Stack // push-to and pop-from its end
	Heap    clutil.Heap  // no GC here, forever growing
	Globals clutil.Env
	Code    code // evaluated l2r
	Stats   clutil.Stats
}

func (me *gMachine) Eval(name string) (val interface{}, stats clutil.Stats, err error) {
	// defer clutil.Catch(&err)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_UNWIND}}
	println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
	me.eval()
	stats, val = me.Stats, me.Heap[me.Stack.Top(0)]
	return
}

func (me *gMachine) eval() {
	for me.Stats.NumSteps, me.Stats.NumAppls = 0, 0; len(me.Code) != 0; me.step() {
	}
}

func (me *gMachine) step() {
	me.Stats.NumSteps, me.Code = me.Stats.NumSteps+1, me.dispatch(me.Code[0], me.Code[1:])
}
