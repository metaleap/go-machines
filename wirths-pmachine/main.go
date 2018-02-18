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

	machine.simpleDemo("(123×456)÷789", "71", []instr{
		{Op: OP_LIT, A: 123},
		{Op: OP_LIT, A: 456},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		{Op: OP_LIT, A: 789},
		{Op: OP_EXEC, A: EXEC_AR_DIV},
		{Op: OP_JUMP, A: 0},
	})

	machine.simpleDemo("987×(654+321)", "962325", []instr{
		{Op: OP_LIT, A: 987},
		{Op: OP_LIT, A: 654},
		{Op: OP_LIT, A: 321},
		{Op: OP_EXEC, A: EXEC_AR_ADD},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		{Op: OP_JUMP, A: 0},
	})

	println("entering REPL..")
}

func (me *interp) simpleDemo(descr string, expectedResult string, programCode []instr) {
	println("Calcing " + descr + ".. — should be: " + expectedResult)
	me.Code = programCode
	me.Run()
	println(me.st[me.t])
}
