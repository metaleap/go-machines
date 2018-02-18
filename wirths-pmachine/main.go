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
	println("https://en.wikipedia.org/wiki/P-code_machine#Example_machine")
}
