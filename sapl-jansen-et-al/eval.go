package sapl

func (me Prog) Eval(expr Expr) Expr {
	return me.eval(expr, make([]Expr, 0, 32))
}

func (me Prog) eval(expr Expr, stack []Expr) Expr {
	switch it := expr.(type) {
	case ExprAppl:
		return me.eval(it.Callee, append(stack, it.Arg))
	case ExprFnRef:
		numargs, isop := 2, (it < 0)
		if !isop {
			numargs = me[it].NumArgs
		}
		if len(stack) < numargs {
			return rebuildAppl(it, stack)
		} else if isop {
			lhs, rhs := me.eval(stack[len(stack)-1], stack).(ExprNum), me.eval(stack[len(stack)-2], stack).(ExprNum)
			stack = stack[:len(stack)-2]
			switch OpCode(it) {
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
				if numargs, it = 2, 2; lhs == rhs {
					it = 1
				}
			case OpLt:
				if numargs, it = 2, 2; lhs < rhs {
					it = 1
				}
			case OpGt:
				if numargs, it = 2, 2; lhs > rhs {
					it = 1
				}
			default:
				panic(stack)
			}
		}
		return me.eval(inst(me[it].Expr, stack), stack[:len(stack)-numargs])
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
