package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type multiplicator struct {
	Operand1     int64
	Operand2     int64
	Jobber       int64
	RunningTotal int64
}

func (me *multiplicator) finalState() bool {
	return me.Operand2 == 0 && me.Jobber == 0
}

func (me *multiplicator) init(op1 int64, op2 int64) {
	me.Operand1, me.Operand2, me.Jobber, me.RunningTotal = op1, op2, 0, 0
}

func (me *multiplicator) step() {
	if me.Jobber > 0 {
		me.Jobber, me.RunningTotal = me.Jobber-1, me.RunningTotal+1
	} else {
		me.Operand2, me.Jobber = me.Operand2-1, me.Operand1
	}
}

func (me *multiplicator) eval(op1 int64, op2 int64) int64 {
	for me.init(op1, op2); !me.finalState(); me.step() {
		// nothing else to do here while we step
	}
	return me.RunningTotal
}

func main() {
	write, readln := os.Stdout.WriteString, bufio.NewScanner(os.Stdin)
	for {
		write("Keep entering 2 ints (separated by 1 space) to have them multiplied horribly inefficiently by a mere state transition machine using +/- prim-ops:\n")
		if !readln.Scan() {
			return
		}
		if snums := strings.Split(strings.TrimSpace(readln.Text()), " "); len(snums) != 2 {
			write("try again\n")
		} else {
			i1, _ := strconv.ParseInt(snums[0], 0, 64)
			i2, _ := strconv.ParseInt(snums[1], 0, 64)
			result := (&multiplicator{}).eval(i1, i2)
			write(strconv.FormatInt(result, 10) + "\n")
		}
	}
}
