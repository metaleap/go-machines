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
	// we get more consistent timings (less runtime background work) and can better micro-experiment with the run() loop with:
	debug.SetGCPercent(-1) // GC off
	runtime.GOMAXPROCS(1)  // no thread scheduling

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

	readln, write := bufio.NewScanner(os.Stdin), os.Stdout.WriteString
	write("\n\nEnter one of the following function names,\nfollowed by 1 space and 1 integral number:\n")
	write("· negodd ‹num›\n  — negates if odd\n")
	write("· negeven ‹num›\n  — negates if even\n")
	write("· fac ‹max 20›\n  — factorial\n\n")
	var result int
	var timetaken time.Duration
	const numrums = 999999
	for readln.Scan() {
		if ln := strings.TrimSpace(readln.Text()); ln != "" {
			isnodd, isnev, isfac := strings.HasPrefix(ln, "negodd"), strings.HasPrefix(ln, "negeven"), strings.HasPrefix(ln, "fac")
			if i := strings.IndexRune(ln, ' '); i > 0 && (isnodd || isnev || isfac) {
				num, err := strconv.ParseInt(ln[i+1:], 0, 64)
				if err != nil {
					println(err.Error())
				} else {
					if arg := int(num); isnodd {
						result, timetaken = runNegIf(1, arg, numrums)
					} else if isnev {
						result, timetaken = runNegIf(0, arg, numrums)
					} else {
						result, timetaken = runFac(arg, numrums)
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

func runNegIf(isNegIfOdd int, num int, runs int) (result int, timeTaken time.Duration) {
	code := codeNegIfEven
	if isNegIfOdd != 0 {
		code = codeNegIfOdd
	}
	code[0].A, code[3].A, code[5+isNegIfOdd].A = num, num, num

	var totalnanos int64
	for i := 0; i < runs; i++ {
		result, timeTaken = interp(code)
		totalnanos += int64(timeTaken)
	}
	timeTaken = time.Duration(int64(float64(totalnanos) / float64(runs)))
	return
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

func runFac(num int, runs int) (result int, timeTaken time.Duration) {
	code := codeFactorial_Opt // codeFactorial_Orig
	code[1].A, code[2].A = num, len(code)-1
	var totalnanos int64
	for i := 0; i < runs; i++ {
		result, timeTaken = interp(code)
		totalnanos += int64(timeTaken)
	}
	timeTaken = time.Duration(int64(float64(totalnanos) / float64(runs)))
	if runs > 1 {
		var cr, ctr int
		var ctn int64
		var ctt time.Duration
		for i := 0; i < runs; i++ {
			cr, ctt = facCompiled(num)
			ctr, ctn = ctr+cr, ctn+int64(ctt)
		}
		if ctr != num+num { // always true but ensure ctr is computed not optimized away by a future too-clever compiler
			println("(avg. over ", runs, " runs, vs. ", time.Duration(int64(float64(ctn)/float64(runs))).String(), ")")
		}
	}
	return
}

// uses only original-pcode opcodes — longer=slower vs codeFactorial_Opt
var codeFactorial_Orig = []instr{ // r := 1; for n>0 { r=r*n ; n = n-1 }
	{Op: OP_LIT, A: 1}, //r=1, t1
	{Op: OP_LIT},       //n, t2

	{Op: OP_JUMPCOND},             // n>0, t1
	{Op: OP_INCR, A: 1},           // restore n, t2
	{Op: OP_STORE, A: 11 - 10},    // stow away n, t1 (addr '11-10' is just for our own later readability because there are `1`s everywhere in here)
	{Op: OP_INCR, A: 1},           // but keep it on-stack still, t2
	{Op: OP_EXEC, A: EXEC_AR_MUL}, // r=r*n, t1
	{Op: OP_LOAD, A: 11 - 10},     // restore n, t2 (addr '11-10' see note above)
	{Op: OP_LIT, A: 1},            // 1, t3
	{Op: OP_EXEC, A: EXEC_AR_SUB}, // n=n-1, t2
	{Op: OP_JUMP, A: 2},           // t2

	opExit,
}

// uses additional/invented/custom opcodes — shorter=faster
var codeFactorial_Opt = []instr{ // r := 1; for n>0 { r=r*n ; n = n-1 }
	{Op: OP_LIT, A: 1}, //r=1, t1
	{Op: OP_LIT},       //n, t2

	{Op: OP_JUMPCOND_KEEP},             // n>0, t2
	{Op: OP_EXEC, A: EXEC_AR_MUL_KEEP}, // r=r*n, t1
	{Op: OP_EXEC, A: EXEC_AR_SUB1},     // n=n-1, t2
	{Op: OP_JUMP, A: 2},                // t2

	opExit,
}

func facCompiled(n int) (result int, timeTaken time.Duration) {
	z, _ := strconv.ParseInt(os.Getenv("_DOESNT_REALLY_EVER_EXIST"), 10, 64)
	zero := int(z)
	timestarted := time.Now()
	for result = 1; n > zero; n-- {
		result = result * n
	}
	timeTaken = time.Now().Sub(timestarted)
	return
}
