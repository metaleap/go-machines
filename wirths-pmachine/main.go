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

func main() {
	machine := interp{}
	machine.Code = []instr{
		{Op: OP_LIT, A: 1234},
		{Op: OP_LIT, A: 5678},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		{Op: OP_JUMP, A: 0},
	}
	println("Calcing 1234×5678..")
	machine.Run()
	println(machine.st[machine.t])
}
