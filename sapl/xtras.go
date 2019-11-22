package sapl

func (me Prog) List(ctx *CtxEval, expr Expr) (ret []Expr) {
	ret = make([]Expr, 0, 1024)
	for again, next := true, expr; again; {
		again = false
		if fouter, ok0 := next.(ExprFnRef); ok0 && fouter == 3 { // clean end-of-list
			break
		} else if aouter, ok1 := next.(ExprAppl); ok1 {
			if ainner, ok2 := aouter.Callee.(ExprAppl); ok2 {
				if finner, ok3 := ainner.Callee.(ExprFnRef); ok3 && finner == 4 {
					elem := me.eval(ainner.Arg, nil, ctx)
					again, next, ret = true, me.eval(aouter.Arg, nil, ctx), append(ret, elem)
				}
			}
		}
		if !again {
			ret = nil
		}
	}
	return
}

func (me Prog) ToBytes(maybeNumList []Expr) (retNumListAsBytes []byte) {
	if maybeNumList != nil {
		retNumListAsBytes = make([]byte, 0, len(maybeNumList))
		for _, expr := range maybeNumList {
			if num, ok := expr.(ExprNum); ok {
				retNumListAsBytes = append(retNumListAsBytes, byte(num))
			} else {
				retNumListAsBytes = nil
				break
			}
		}
	}
	return
}

func ListsFrom(strs []string) (ret Expr) {
	ret = ExprFnRef(3)
	for i := len(strs) - 1; i > -1; i-- {
		ret = ExprAppl{ExprAppl{ExprFnRef(4), ListFrom([]byte(strs[i]))}, ret}
	}
	return
}

func ListFrom(str []byte) (ret Expr) {
	ret = ExprFnRef(3)
	for i := len(str) - 1; i > -1; i-- {
		ret = ExprAppl{ExprAppl{ExprFnRef(4), ExprNum(str[i])}, ret}
	}
	return
}
