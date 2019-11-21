package sapl

import (
	"time"
)

type OpCode int

const (
	OpPanic OpCode = -1234567890
	OpAdd   OpCode = -1
	OpSub   OpCode = -2
	OpMul   OpCode = -3
	OpDiv   OpCode = -4
	OpMod   OpCode = -5
	OpEq    OpCode = -6
	OpLt    OpCode = -7
	OpGt    OpCode = -8
)

func (me Prog) Eval(expr Expr, maybeTracer func(Expr, []Expr) func(Expr) Expr) (ret Expr, timeTaken time.Duration) {
	stack := make([]Expr, 0, 128)
	tstart := time.Now().UnixNano()
	if maybeTracer == nil {
		ret := func(it Expr) Expr { return it }
		maybeTracer = func(Expr, []Expr) func(Expr) Expr { return ret }
	}
	ret = me.eval(expr, stack, maybeTracer)
	timeTaken = time.Duration(time.Now().UnixNano() - tstart)
	return
}

func (me Prog) eval(expr Expr, stack []Expr, tracer func(Expr, []Expr) func(Expr) Expr) Expr {
	ret := tracer(expr, stack)
	switch it := expr.(type) {
	case ExprAppl:
		return ret(me.eval(it.Callee, append(stack, it.Arg), tracer))
	case ExprFnRef:
		numargs, isopcode := 2, (it < 0)
		if !isopcode {
			numargs = me[it].NumArgs
		}
		if len(stack) < numargs { // not enough args on stack: a partial-application aka closure
			for i := len(stack) - 1; i >= 0; i-- {
				expr = ExprAppl{expr, stack[i]}
			}
			return ret(expr)
		} else if isopcode {
			lhs, rhs := me.eval(stack[len(stack)-1], stack, tracer).(ExprNum), me.eval(stack[len(stack)-2], stack, tracer).(ExprNum)
			stack = stack[:len(stack)-2]
			switch OpCode(it) {
			case OpAdd:
				return ret(lhs + rhs)
			case OpSub:
				return ret(lhs - rhs)
			case OpMul:
				return ret(lhs * rhs)
			case OpDiv:
				return ret(lhs / rhs)
			case OpMod:
				return ret(lhs % rhs)
			case OpEq, OpGt, OpLt:
				if op := OpCode(it); (op == OpEq && lhs == rhs) || (op == OpLt && lhs < rhs) || (op == OpGt && lhs > rhs) {
					it, numargs = 1, 2
				} else {
					it, numargs = 2, 2
				}
			default:
				panic(stack)
			}
		}
		return ret(me.eval(argRefsResolvedToCurrentStackEntries(me[it].Expr, stack), stack[:len(stack)-numargs], tracer))
	}
	return ret(expr)
}

func argRefsResolvedToCurrentStackEntries(expr Expr, stack []Expr) Expr {
	switch it := expr.(type) {
	case ExprAppl:
		return ExprAppl{argRefsResolvedToCurrentStackEntries(it.Callee, stack), argRefsResolvedToCurrentStackEntries(it.Arg, stack)}
	case ExprArgRef:
		return stack[len(stack)+int(it)]
	}
	return expr
}
