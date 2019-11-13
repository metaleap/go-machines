package sapl

import (
	"strconv"
)

type OpCode int

const (
	OpAdd OpCode = -1
	OpSub OpCode = -2
	OpMul OpCode = -3
	OpDiv OpCode = -4
	OpMod OpCode = -5
	OpEq  OpCode = -6
	OpLt  OpCode = -7
	OpGt  OpCode = -8
)

type Expr interface{ String() string }

func (me ExprNum) String() string    { return strconv.Itoa(int(me)) }
func (me ExprArgRef) String() string { return "#" + strconv.Itoa(int(me)) }
func (me ExprFnRef) String() string  { return strconv.Itoa(me.NumArgs) + "@" + strconv.Itoa(me.Idx) }
func (me ExprAppl) String() string   { return "(" + me.Callee.String() + " " + me.Arg.String() + ")" }

type ExprNum int

type ExprArgRef int

type ExprAppl struct {
	Callee Expr
	Arg    Expr
}

type ExprFnRef struct {
	NumArgs int
	Idx     int
}

func (me *ExprFnRef) boolish(isTrue bool) {
	if me.NumArgs, me.Idx = 2, 2; isTrue {
		me.Idx = 1
	}
}

func (me Prog) Eval(expr Expr) Expr {
	return me.eval(expr, make([]Expr, 0, 32))
}

func (me Prog) eval(expr Expr, stack []Expr) Expr {
	switch it := expr.(type) {
	case ExprAppl:
		return me.eval(it.Callee, append(stack, it.Arg))
	case ExprFnRef:
		if len(stack) < it.NumArgs {
			return rebuildAppl(it, stack)
		}
		if it.Idx < 0 {
			lhs, rhs := stack[len(stack)-1].(ExprNum), stack[len(stack)-2].(ExprNum)
			stack = stack[:len(stack)-2]
			switch OpCode(it.Idx) {
			case OpAdd:
				return (lhs + rhs)
			case OpSub:
				return (lhs - rhs)
			case OpMul:
				return (lhs * rhs)
			case OpDiv:
				return (lhs / rhs)
			case OpMod:
				return (lhs % rhs)
			case OpEq:
				it.boolish(lhs == rhs)
			case OpLt:
				it.boolish(lhs < rhs)
			case OpGt:
				it.boolish(lhs > rhs)
			default:
				panic(stack)
			}
		}
		return me.eval(inst(me[it.Idx], stack), stack[:len(stack)-it.NumArgs])
	}
	return expr
}

func rebuildAppl(expr Expr, stack []Expr) Expr {
	for len(stack) > 0 {
		expr, stack = ExprAppl{expr, stack[len(stack)-1]}, stack[:len(stack)-1]
	}
	return expr
}

func inst(expr Expr, stack []Expr) Expr {
	switch it := expr.(type) {
	case ExprAppl:
		return ExprAppl{inst(it.Callee, stack), inst(it.Arg, stack)}
	case ExprArgRef:
		return stack[(len(stack)-1)-int(it)]
	}
	return expr
}
