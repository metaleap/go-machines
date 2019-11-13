package sapl

type OpCode int

const (
	_ OpCode = -iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpEq
	OpLt
	OpGt
)

type Expr interface{ String() string }

type ExprNum int

type ExprVar int

type ExprAppl struct {
	lhs Expr
	rhs Expr
}

type ExprFunc struct {
	numArgs int
	idx     int
}

func (me *ExprFunc) boolish(isTrue bool) {
	if me.numArgs, me.idx = 2, 2; isTrue {
		me.idx = 1
	}
}

func (me Prog) Eval(expr Expr) Expr {
	return me.eval(expr, make([]Expr, 0, 32))
}

func (me Prog) eval(expr Expr, stack []Expr) Expr {
	switch it := expr.(type) {
	case ExprAppl:
		return me.eval(it.lhs, append(stack, it.rhs))
	case ExprFunc:
		if len(stack) < it.numArgs {
			return rebuildAppl(it, stack)
		}
		if it.idx < 0 {
			lhs, rhs := stack[len(stack)-2].(ExprNum), stack[len(stack)-1].(ExprNum)
			stack = stack[:len(stack)-2]
			switch OpCode(it.idx) {
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
		return me.eval(inst(me[it.idx], stack), stack[:len(stack)-it.numArgs])
	}
	return expr
}

func rebuildAppl(expr Expr, stack []Expr) Expr {
	for len(stack) > 0 {
		expr, stack = ExprAppl{lhs: expr, rhs: stack[len(stack)-1]}, stack[:len(stack)-1]
	}
	return expr
}

func inst(expr Expr, stack []Expr) Expr {
	switch it := expr.(type) {
	case ExprAppl:
		return ExprAppl{lhs: inst(it.lhs, stack), rhs: inst(it.rhs, stack)}
	case ExprVar:
		return stack[(len(stack)-1)-int(it)]
	}
	return expr
}
