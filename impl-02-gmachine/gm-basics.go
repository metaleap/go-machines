package climpl

import (
	"github.com/metaleap/go-corelang/util"
)

type gMachine struct {
	Stack         []clutil.Addr // push-to and pop-from end
	Heap          clutil.Heap   // no GC here, forever growing
	Globals       map[string]clutil.Addr
	Code          code // sequential ordering
	NumStepsTaken int
}

func (me *gMachine) Eval(name string) (val interface{}, numSteps int, err error) {
	// defer clutil.Catch(&err)
	me.Code = code{{Op: INSTR_PUSHGLOBAL, Name: name}, {Op: INSTR_UNWIND}}
	// println(me.Heap[me.Globals[name]].(nodeGlobal).Code.String())
	me.eval()
	numSteps, val = me.NumStepsTaken, me.Heap[me.Stack[len(me.Stack)-1]]
	return
}

func (me *gMachine) eval() {
	for me.NumStepsTaken = 0; len(me.Code) != 0 && me.NumStepsTaken <= 99999; me.NumStepsTaken++ {
		me.step()
	}
	if me.NumStepsTaken >= 99999 {
		panic("infinite loop")
	}
}

func (me *gMachine) step() {
	me.Code = me.dispatch(me.Code[0], me.Code[1:])
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
