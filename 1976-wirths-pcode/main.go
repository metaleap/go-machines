package main

import (
	"bufio"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var opExit = instr{Op: OP_JUMP, A: 0}

func main() {
	debug.SetGCPercent(-1)
	runtime.GOMAXPROCS(1)
	runDemo("(123×456)÷789", "71", []instr{
		{Op: OP_LIT, A: 123},
		{Op: OP_LIT, A: 456},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		{Op: OP_LIT, A: 789},
		{Op: OP_EXEC, A: EXEC_AR_DIV},
		opExit,
	})

	runDemo("987×(654-321)", "328671", []instr{
		{Op: OP_LIT, A: 987},
		{Op: OP_LIT, A: 654},
		{Op: OP_LIT, A: 321},
		{Op: OP_EXEC, A: EXEC_AR_SUB},
		{Op: OP_EXEC, A: EXEC_AR_MUL},
		opExit,
	})

	var result int
	var timetaken time.Duration
	readln, write := bufio.NewScanner(os.Stdin), os.Stdout.WriteString
	write("\n\nEnter one of the following function names,\nfollowed by 1 space and 1 integral number:\n")
	write("· negodd ‹num›\n  — negates if odd\n")
	write("· negeven ‹num›\n  — negates if even\n")
	write("· fac ‹max 20›\n  — factorial\n\n")
	for readln.Scan() {
		if ln := strings.TrimSpace(readln.Text()); ln != "" {
			isnodd, isnev, isfac := strings.HasPrefix(ln, "negodd"), strings.HasPrefix(ln, "negeven"), strings.HasPrefix(ln, "fac")
			if i := strings.IndexRune(ln, ' '); i > 0 && (isnodd || isnev || isfac) {
				num, err := strconv.ParseInt(ln[i+1:], 0, 64)
				if err != nil {
					println(err.Error())
				} else {
					if arg := int(num); isnodd {
						result, timetaken = runNegIf(1, arg)
					} else if isnev {
						result, timetaken = runNegIf(0, arg)
					} else {
						result, timetaken = runFac(arg)
					}
					write(timetaken.String() + "\n")
					println(result)
				}
			}
		}
	}
}

func runDemo(descr string, expectedResult string, code []instr) {
	println("Calcing " + descr + ".. — should be: " + expectedResult)
	result, _ := interp(code)
	println(result)
}

func runNegIf(isNegIfOdd int, num int) (int, time.Duration) {
	code := codeNegIfEven
	if isNegIfOdd != 0 {
		code = codeNegIfOdd
	}
	code[0].A, code[3].A, code[5+isNegIfOdd].A = num, num, num
	return interp(code)
}

var codeNegIfEven = []instr{
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_ODD},
	{Op: OP_JUMPCOND, A: 5},
	{Op: OP_LIT},
	opExit,
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_NEG},
	opExit,
}

var codeNegIfOdd = []instr{
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_ODD},
	{Op: OP_JUMPCOND, A: 6},
	{Op: OP_LIT},
	{Op: OP_EXEC, A: EXEC_NEG},
	opExit,
	{Op: OP_LIT},
	opExit,
}

func runFac(num int) (int, time.Duration) {
	code := codeFactorialLoop
	code[1].A, code[2].A = num, len(code)-1
	return interp(code)
}

var codeFactorialLoop = []instr{ // r := 1; for n>0 { r=r*n ; n = n-1 }
	{Op: OP_LIT, A: 1}, //r=1, t1
	{Op: OP_LIT},       //n, t2

	{Op: OP_JUMPCOND_K},            // n>0, t2
	{Op: OP_STORE_K, A: 11 - 10},   // stow away n, t2 (addr '11-10' is just for our own later readability because `1` at first reads ambiguous)
	{Op: OP_EXEC, A: EXEC_AR_MUL},  // r=r*n, t1
	{Op: OP_LOAD, A: 11 - 10},      // restore n, t2 (addr '11-10' see note above)
	{Op: OP_EXEC, A: EXEC_AR_SUB1}, // n=n-1, t2
	{Op: OP_JUMP, A: 2},            // t2

	opExit,
}
