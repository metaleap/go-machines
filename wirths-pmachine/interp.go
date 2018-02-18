package main

type execOpCode = int

const (
	EXEC_RET execOpCode = iota
	EXEC_NEG
	EXEC_AR_ADD
	EXEC_AR_SUB
	EXEC_AR_MUL
	EXEC_AR_DIV
	EXEC_ODD
	_
	EXEC_CMP_EQ
	EXEC_CMP_NEQ
	EXEC_CMP_LT
	EXEC_CMP_GEQ
	EXEC_CMP_GT
	EXEC_CMP_LEQ
)

type interp struct {
	Prog     int
	Base     int
	TopStack int
	Instr    int
	Stack    [1024]int
	Code     []instr
}

func (me *interp) base(l int) (b1 int) {
	b1 = me.Base
	for l > 0 {
		b1 = me.Stack[b1]
		l--
	}
	return
}

func (me *interp) run() {
	me.TopStack, me.Base, me.Prog = 0, 1, 0
	me.Stack[1], me.Stack[2], me.Stack[3] = 0, 0, 0
	var i int
	for {
		i = me.Prog
		me.Prog++
		switch me.Code[i].Op {
		case OP_LIT:
			me.TopStack++
			me.Stack[me.TopStack] = me.Code[i].A
		case OP_LOAD:
			me.TopStack++
			me.Stack[me.TopStack] = me.Stack[me.base(me.Code[i].L)+me.Code[i].A]
		case OP_STORE:
			me.Stack[me.base(me.Code[i].L)+me.Code[i].A] = me.Stack[me.TopStack]
			me.TopStack--
		case OP_CALL:
			me.Stack[me.TopStack+1] = me.base(me.Code[i].L)
			me.Stack[me.TopStack+2] = me.Base
			me.Stack[me.TopStack+3] = me.Prog
			me.Base = me.TopStack + 1
			me.Prog = me.Code[i].A
		case OP_INCR:
			me.TopStack = me.TopStack + me.Code[i].A
		case OP_JUMP:
			me.Prog = me.Code[i].A
		case OP_JUMPCOND:
			if me.Stack[me.TopStack] == 0 {
				me.Prog = me.Code[i].A
			}
			me.TopStack--
		case OP_EXEC:
			switch me.Code[i].A {
			case EXEC_RET:
				me.TopStack = me.Base - 1
				me.Prog = me.Stack[me.TopStack+3]
				me.Base = me.Stack[me.TopStack+2]
			case EXEC_NEG:
				me.Stack[me.TopStack] = -me.Stack[me.TopStack]
			case EXEC_AR_ADD:
				me.TopStack--
				me.Stack[me.TopStack] = me.Stack[me.TopStack] + me.Stack[me.TopStack+1]
			case EXEC_AR_SUB:
				me.TopStack--
				me.Stack[me.TopStack] = me.Stack[me.TopStack] - me.Stack[me.TopStack+1]
			case EXEC_AR_MUL:
				me.TopStack--
				me.Stack[me.TopStack] = me.Stack[me.TopStack] * me.Stack[me.TopStack+1]
			case EXEC_AR_DIV:
				me.TopStack--
				me.Stack[me.TopStack] = me.Stack[me.TopStack] / me.Stack[me.TopStack+1]
			case EXEC_ODD:
				me.Stack[me.TopStack] = me.Stack[me.TopStack] & 1
			case EXEC_CMP_EQ:
				me.TopStack--
				if me.Stack[me.TopStack] == me.Stack[me.TopStack+1] {
					me.Stack[me.TopStack] = 1
				} else {
					me.Stack[me.TopStack] = 0
				}
			case EXEC_CMP_NEQ:
				me.TopStack--
				if me.Stack[me.TopStack] != me.Stack[me.TopStack+1] {
					me.Stack[me.TopStack] = 1
				} else {
					me.Stack[me.TopStack] = 0
				}
			case EXEC_CMP_LT:
				me.TopStack--
				if me.Stack[me.TopStack] < me.Stack[me.TopStack+1] {
					me.Stack[me.TopStack] = 1
				} else {
					me.Stack[me.TopStack] = 0
				}
			case EXEC_CMP_GEQ:
				me.TopStack--
				if me.Stack[me.TopStack] >= me.Stack[me.TopStack+1] {
					me.Stack[me.TopStack] = 1
				} else {
					me.Stack[me.TopStack] = 0
				}
			case EXEC_CMP_GT:
				me.TopStack--
				if me.Stack[me.TopStack] > me.Stack[me.TopStack+1] {
					me.Stack[me.TopStack] = 1
				} else {
					me.Stack[me.TopStack] = 0
				}
			case EXEC_CMP_LEQ:
				me.TopStack--
				if me.Stack[me.TopStack] <= me.Stack[me.TopStack+1] {
					me.Stack[me.TopStack] = 1
				} else {
					me.Stack[me.TopStack] = 0
				}
			}
		}
		if me.Prog == 0 {
			break
		}
	}
}
