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

func main() {
	mul, readln, write := &multiplicator{}, bufio.NewScanner(os.Stdin), os.Stdout.WriteString
	for {
		if write("Keep entering 2 ints (separated by 1 space) to have them (horribly inefficiently) multiplied by a state transition machine via mere incr1/decr1 operations:\n"); !readln.Scan() {
			return
		} else if operands := strings.Split(strings.TrimSpace(readln.Text()), " "); len(operands) != 2 {
			write("try again\n")
		} else {
			operand1, _ := strconv.ParseInt(operands[0], 0, 64)
			operand2, _ := strconv.ParseInt(operands[1], 0, 64)
			result := mul.eval(operand1, operand2)
			write(strconv.FormatInt(result, 10) + "\n")
		}
	}
}

func (this *multiplicator) eval(op1 int64, op2 int64) int64 {
	for this.init(op1, op2); !this.finalState(); this.step() {
		// nothing else to do here while we step
	}
	return this.RunningTotal
}

func (this *multiplicator) init(op1 int64, op2 int64) {
	this.Operand1, this.Operand2, this.Jobber, this.RunningTotal = op1, op2, 0, 0
}

func (this *multiplicator) finalState() bool {
	return this.Operand2 == 0 && this.Jobber == 0
}

func (this *multiplicator) step() {
	if this.Jobber == 0 {
		this.Jobber, this.Operand2 = this.Operand1, this.Operand2-1
	} else {
		this.Jobber, this.RunningTotal = this.Jobber-1, this.RunningTotal+1
	}
}
