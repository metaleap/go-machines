package main

type opCode int

const (
	_ opCode = iota
	OP_LIT
	OP_EXEC
	OP_LOAD
	OP_STORE
	OP_CALL
	OP_INCR
	OP_JUMP
	OP_JUMPCOND
)

type instr struct {
	Op opCode
	A  int
	L  int
}

type execOpCode = int

const (
	EXEC_RET execOpCode = iota
	EXEC_NEG
	EXEC_AR_ADD
	EXEC_AR_SUB
	EXEC_AR_MUL
	EXEC_AR_DIV
	EXEC_ODD
	EXEC_DBG
	EXEC_CMP_EQ
	EXEC_CMP_NEQ
	EXEC_CMP_LT
	EXEC_CMP_GEQ
	EXEC_CMP_GT
	EXEC_CMP_LEQ
)

type interp struct {
	p    int      // 'program register' (aka (next-)instruction pointer)
	b    int      // 'base register'
	t    int      // 'top-stack register' (aka stack pointer)
	st   [512]int // stack
	code []instr
}

func (me *interp) base(l int) (b int) {
	for b = me.b; l > 0; l-- {
		b = me.st[b]
	}
	return
}

func (me *interp) run() int {
	me.t, me.b, me.p = 0, 1, 0
	me.st[1], me.st[2], me.st[3] = 0, 0, 0

	for i, running := 0, true; running; running = me.p != 0 {
		i = me.p
		me.p++

		switch me.code[i].Op {
		case OP_LIT:
			me.t++
			me.st[me.t] = me.code[i].A

		case OP_INCR:
			me.t = me.t + me.code[i].A

		case OP_JUMP:
			me.p = me.code[i].A

		case OP_JUMPCOND:
			if me.st[me.t] == 0 {
				me.p = me.code[i].A
			}
			me.t--

		case OP_LOAD:
			me.t++
			me.st[me.t] = me.st[me.base(me.code[i].L)+me.code[i].A]

		case OP_STORE:
			me.st[me.base(me.code[i].L)+me.code[i].A] = me.st[me.t]
			me.t--

		case OP_CALL:
			me.st[me.t+1] = me.base(me.code[i].L)
			me.st[me.t+2] = me.b
			me.st[me.t+3] = me.p
			me.b = me.t + 1
			me.p = me.code[i].A

		case OP_EXEC:
			switch me.code[i].A {
			case EXEC_RET:
				me.t = me.b - 1
				me.p = me.st[me.t+3]
				me.b = me.st[me.t+2]

			case EXEC_NEG:
				me.st[me.t] = -me.st[me.t]

			case EXEC_AR_ADD:
				me.t--
				me.st[me.t] = me.st[me.t] + me.st[me.t+1]

			case EXEC_AR_SUB:
				me.t--
				me.st[me.t] = me.st[me.t] - me.st[me.t+1]

			case EXEC_AR_MUL:
				me.t--
				me.st[me.t] = me.st[me.t] * me.st[me.t+1]

			case EXEC_AR_DIV:
				me.t--
				me.st[me.t] = me.st[me.t] / me.st[me.t+1]

			case EXEC_ODD:
				me.st[me.t] = me.st[me.t] & 1

			case EXEC_CMP_EQ:
				me.t--
				if me.st[me.t] == me.st[me.t+1] {
					me.st[me.t] = 1
				} else {
					me.st[me.t] = 0
				}

			case EXEC_CMP_NEQ:
				me.t--
				if me.st[me.t] != me.st[me.t+1] {
					me.st[me.t] = 1
				} else {
					me.st[me.t] = 0
				}

			case EXEC_CMP_LT:
				me.t--
				if me.st[me.t] < me.st[me.t+1] {
					me.st[me.t] = 1
				} else {
					me.st[me.t] = 0
				}

			case EXEC_CMP_GEQ:
				me.t--
				if me.st[me.t] >= me.st[me.t+1] {
					me.st[me.t] = 1
				} else {
					me.st[me.t] = 0
				}

			case EXEC_CMP_GT:
				me.t--
				if me.st[me.t] > me.st[me.t+1] {
					me.st[me.t] = 1
				} else {
					me.st[me.t] = 0
				}

			case EXEC_CMP_LEQ:
				me.t--
				if me.st[me.t] <= me.st[me.t+1] {
					me.st[me.t] = 1
				} else {
					me.st[me.t] = 0
				}

			case EXEC_DBG:
				print("t")
				print(me.t)
				print("=")
				println(me.st[me.t])
			}
		}
	}
	return me.st[me.t]
}
