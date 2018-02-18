package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

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

var opExit = instr{Op: OP_JUMP}

func main() {
	machine := interp{}

	machine.simpleDemo("(123×456)÷789", "71", []instr{
		{Op: OP_LIT, A: 123},
		{Op: OP_LIT, A: 456},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		{Op: OP_LIT, A: 789},
		{Op: OP_EXEC, A: EXEC_AR_DIV},
		opExit,
	})

	machine.simpleDemo("987×(654+321)", "962325", []instr{
		{Op: OP_LIT, A: 987},
		{Op: OP_LIT, A: 654},
		{Op: OP_LIT, A: 321},
		{Op: OP_EXEC, A: EXEC_AR_ADD},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		opExit,
	})

	readln, write := bufio.NewScanner(os.Stdin), os.Stdout.WriteString
	write("\n\nEnter one of the following function names, followed by 1 space and 1 number:\n")
	write("· negodd — negates the number if it is odd\n")
	write("· negeven — negates the number if it is even\n")
	write("· fac — computes the number's factorial\n")
	write("· q — quits\n\n")
	for readln.Scan() {
		if ln := strings.TrimSpace(readln.Text()); ln != "" {
			if ln == "q" {
				return
			}
			isnodd, isnev, isfac := strings.HasPrefix(ln, "negodd"), strings.HasPrefix(ln, "negeven"), strings.HasPrefix(ln, "fac")
			if i := strings.IndexRune(ln, ' '); i > 0 && (isnodd || isnev || isfac) {
				num, err := strconv.ParseInt(ln[i+1:], 0, 64)
				if err != nil {
					println(err.Error())
				}
				arg := int(num)

				var result int
				var timetaken time.Duration
				if isnodd {
					result, timetaken = machine.simpleNegIf(negIfOdd, 1, arg)
				} else if isnev {
					result, timetaken = machine.simpleNegIf(negIfEven, 0, arg)
				}
				println(timetaken.String())
				println(result)
			}
		}
	}
}

func (me *interp) simpleDemo(descr string, expectedResult string, programCode []instr) {
	println("Calcing " + descr + ".. — should be: " + expectedResult)
	me.code = programCode
	println(me.run())
}

func (me *interp) simpleNegIf(negIfCode []instr, off int, num int) (result int, timeTaken time.Duration) {
	me.code = negIfCode
	me.code[0].A, me.code[3].A, me.code[5+off].A = num, num, num

	timestarted := time.Now()
	result = me.run()
	timeTaken = time.Now().Sub(timestarted)
	return
}

var negIfEven = []instr{
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_ODD},
	{Op: OP_JUMPCOND, A: 5},
	{Op: OP_LIT},
	opExit,
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_NEG},
	opExit,
}

var negIfOdd = []instr{
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_ODD},
	{Op: OP_JUMPCOND, A: 6},
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_NEG},
	opExit,
	{Op: OP_LIT},
	opExit,
}
