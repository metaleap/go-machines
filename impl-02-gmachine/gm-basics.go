package climpl

import (
	"github.com/metaleap/go-corelang/util"
)

type gMachine struct {
	Stack           []clutil.Addr // push-to and pop-from its end
	Heap            clutil.Heap   // no GC here, forever growing
	Globals         map[string]clutil.Addr
	Code            code // evaluated l2r
	NumApplications int
	NumStepsTaken   int
}

func (me *gMachine) Eval(name string) (val interface{}, numAppl int, numSteps int, err error) {
	// defer clutil.Catch(&err)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_UNWIND}}
	// println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
	me.eval()
	numAppl, numSteps, val = me.NumApplications, me.NumStepsTaken, me.Heap[me.Stack[len(me.Stack)-1]]
	return
}

func (me *gMachine) eval() {
	for me.NumStepsTaken, me.NumApplications = 0, 0; len(me.Code) != 0; me.step() {
	}
}

func (me *gMachine) step() {
	me.NumStepsTaken, me.Code = me.NumStepsTaken+1, me.dispatch(me.Code[0], me.Code[1:])
}

func (me *gMachine) alloc(obj clutil.INode) (addr clutil.Addr) {
	addr = me.nextAddr()
	me.Heap[addr] = obj
	return
}

func (me *gMachine) lookup(name string) (addr clutil.Addr) {
	if addr = me.Globals[name]; addr == 0 {
		panic("undefined: " + name)
	}
	return
}

func (me *gMachine) nextAddr() clutil.Addr {
	return clutil.Addr(len(me.Heap) + 1)
}
