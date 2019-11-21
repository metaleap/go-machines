package sapl

import (
	"strconv"
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

type CtxEval struct {
	tracer func(Expr, []Expr) func(Expr) Expr

	MaxStack    int
	NumSteps    int
	NumRebuilds int
	NumCalls    int
	TimeTaken   time.Duration
}

func (me *CtxEval) String() string {
	return me.TimeTaken.String() + "\tMaxStack=" + strconv.Itoa(me.MaxStack) + "\tNumSteps=" + strconv.Itoa(me.NumSteps) + "\tNumRebuilds=" + strconv.Itoa(me.NumRebuilds) + "\tNumCalls=" + strconv.Itoa(me.NumCalls)
}

func (me Prog) Eval(expr Expr, maybeTracer func(Expr, []Expr) func(Expr) Expr) (ret Expr, stats CtxEval) {
	if stats.tracer = maybeTracer; stats.tracer == nil {
		ret := func(it Expr) Expr { return it }
		stats.tracer = func(Expr, []Expr) func(Expr) Expr { return ret }
	}
	stack := make([]Expr, 0, 128)
	tstart := time.Now().UnixNano()
	ret = me.eval(expr, stack, &stats)
	stats.TimeTaken = time.Duration(time.Now().UnixNano() - tstart)
	return
}

func (me Prog) eval(expr Expr, stack []Expr, ctx *CtxEval) Expr {
	if ctx.NumSteps++; len(stack) > ctx.MaxStack {
		ctx.MaxStack = len(stack)
	}
	ret := ctx.tracer(expr, stack)
	switch it := expr.(type) {
	case ExprAppl:
		return ret(me.eval(it.Callee, append(stack, it.Arg), ctx))
	case ExprFnRef:
		numargs, isopcode := 2, (it < 0)
		if !isopcode {
			numargs = me[it].NumArgs
		}
		if len(stack) < numargs { // not enough args on stack: a partial-application aka closure
			ctx.NumRebuilds++
			for i := len(stack) - 1; i >= 0; i-- {
				expr = ExprAppl{expr, stack[i]}
			}
			return ret(expr)
		} else if isopcode {
			lhs, rhs := me.eval(stack[len(stack)-1], stack, ctx).(ExprNum), me.eval(stack[len(stack)-2], stack, ctx).(ExprNum)
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
		ctx.NumCalls++
		return ret(me.eval(argRefsResolvedToCurrentStackEntries(me[it].Expr, stack), stack[:len(stack)-numargs], ctx))
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
