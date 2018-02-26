package main

import (
	"time"
)

type opCode int

const (
	_           opCode = iota
	OP_LIT             // lit
	OP_EXEC            // opr
	OP_LOAD            // lod
	OP_STORE           // sto
	OP_CALL            // cal
	OP_INCR            // int
	OP_JUMP            // jmp
	OP_JUMPCOND        // jpc

	OP_INCR1         // custom, not pcode orig
	OP_JUMPCOND_KEEP // dito
	OP_STORE_KEEP    // dito
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
	EXEC_AR_SUB1
	EXEC_AR_MUL
	EXEC_AR_MUL_KEEP
	EXEC_AR_DIV
	EXEC_ODD
	EXEC_DBG // custom, not original
	EXEC_CMP_EQ
	EXEC_CMP_NEQ
	EXEC_CMP_LT
	EXEC_CMP_GEQ
	EXEC_CMP_GT
	EXEC_CMP_LEQ
)

func interp(code []instr) (int, time.Duration) {
	var (
		p  int    // 'program register' (aka (next-)instruction pointer)
		i  int    // curr-instruction pointer
		b  = 1    // 'base register'
		t  int    // 'top-stack register' (aka stack pointer)
		st [4]int // stack — make it bigger as needed, 4 is the minimum for current built-ins like fac

		tmpb int
		tmpl int
	)

	timestarted := time.Now()
	for running := true; running; running = (p != 0) {
		i = p
		p++
		switch code[i].Op {
		case OP_LIT:
			t++
			st[t] = code[i].A

		case OP_JUMP:
			p = code[i].A

		case OP_JUMPCOND_KEEP:
			if st[t] == 0 {
				p = code[i].A
				t--
			}

		case OP_STORE_KEEP:
			for tmpb, tmpl = b, code[i].L; tmpl > 0; tmpl-- {
				tmpb = st[tmpb]
			}
			st[tmpb+code[i].A] = st[t]

		case OP_LOAD:
			for tmpb, tmpl = b, code[i].L; tmpl > 0; tmpl-- {
				tmpb = st[tmpb]
			}
			t++
			st[t] = st[tmpb+code[i].A]

		case OP_STORE:
			for tmpb, tmpl = b, code[i].L; tmpl > 0; tmpl-- {
				tmpb = st[tmpb]
			}
			st[tmpb+code[i].A] = st[t]
			t--

		case OP_INCR:
			t = t + code[i].A

		case OP_INCR1:
			t++

		case OP_JUMPCOND:
			if st[t] == 0 {
				p = code[i].A
			}
			t--

		case OP_CALL:
			for tmpb, tmpl = b, code[i].L; tmpl > 0; tmpl-- {
				tmpb = st[tmpb]
			}
			st[t+1] = tmpb
			st[t+2] = b
			st[t+3] = p
			b = t + 1
			p = code[i].A

		default: // case OP_EXEC:
			switch code[i].A {
			case EXEC_AR_MUL_KEEP:
				st[t-1] = st[t-1] * st[t]

			case EXEC_AR_SUB1:
				st[t]--

			case EXEC_AR_MUL:
				t--
				st[t] = st[t] * st[t+1]

			case EXEC_AR_SUB:
				t--
				st[t] = st[t] - st[t+1]

			case EXEC_AR_ADD:
				t--
				st[t] = st[t] + st[t+1]

			case EXEC_AR_DIV:
				t--
				st[t] = st[t] / st[t+1]

			case EXEC_ODD:
				st[t] = st[t] & 1

			case EXEC_CMP_EQ:
				t--
				if st[t] == st[t+1] {
					st[t] = 1
				} else {
					st[t] = 0
				}

			case EXEC_CMP_NEQ:
				t--
				if st[t] != st[t+1] {
					st[t] = 1
				} else {
					st[t] = 0
				}

			case EXEC_CMP_LT:
				t--
				if st[t] < st[t+1] {
					st[t] = 1
				} else {
					st[t] = 0
				}

			case EXEC_CMP_GEQ:
				t--
				if st[t] >= st[t+1] {
					st[t] = 1
				} else {
					st[t] = 0
				}

			case EXEC_CMP_GT:
				t--
				if st[t] > st[t+1] {
					st[t] = 1
				} else {
					st[t] = 0
				}

			case EXEC_CMP_LEQ:
				t--
				if st[t] <= st[t+1] {
					st[t] = 1
				} else {
					st[t] = 0
				}

			case EXEC_NEG:
				st[t] = -st[t]

			case EXEC_RET:
				t = b - 1
				p = st[t+3]
				b = st[t+2]

			case EXEC_DBG:
				print("t")
				print(t)
				print("=")
				println(st[t])
			}
		}
	}
	timetaken := time.Now().Sub(timestarted)

	return st[t], timetaken
}
