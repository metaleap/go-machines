// SAPL interpreter implementation following: **"Efficient Interpretation by Transforming Data Types and Patterns to Functions"** (Jan Martin Jansen, Pieter Koopman, Rinus Plasmeijer)
//
// Divergence from the paper: NumArgs is not carried around with the Func Ref but stored in the top-level-funcs array together with that func's expression.
//
// "Non"-Parser loads from a JSON format: no need to expressly spec it out here, it's under 40 LoC in `prog.go`'s `LoadFromJson([]byte)`.
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
			default:
				if op := OpCode(it); (op == OpEq && lhs == rhs) || (op == OpLt && lhs < rhs) || (op == OpGt && lhs > rhs) {
					it, numargs = 1, 2
				} else {
					it, numargs = 2, 2
				}
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
