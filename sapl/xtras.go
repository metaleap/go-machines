package sapl

func (me Prog) BytesFromList(ctx *CtxEval, ret Expr, preAlloc []byte) (retIntListAsBytes []byte) {
	retIntListAsBytes = preAlloc
	for again, next := true, ret; again; { // if the ret is an int-list, force it into `retIntListAsBytes`
		again = false
		if fouter, ok0 := next.(ExprFnRef); ok0 && fouter == 3 { // clean end-of-list
			break
		} else if aouter, ok1 := next.(ExprAppl); ok1 {
			if ainner, ok2 := aouter.Callee.(ExprAppl); ok2 {
				if finner, ok3 := ainner.Callee.(ExprFnRef); ok3 && finner == 4 {
					if hd, ok4 := me.eval(ainner.Arg, nil, ctx).(ExprNum); ok4 {
						again, next, retIntListAsBytes = true, me.eval(aouter.Arg, nil, ctx), append(retIntListAsBytes, byte(hd))
					}
				}
			}
		}
		if !again {
			retIntListAsBytes = nil
		}
	}
	return
}

func ListsFrom(strs []string) (ret Expr) {
	ret = ExprFnRef(3)
	for i := len(strs) - 1; i > -1; i-- {
		ret = ExprAppl{ExprAppl{ExprFnRef(4), ListFrom(strs[i])}, ret}
	}
	return
}

func ListFrom(str string) (ret Expr) {
	ret = ExprFnRef(3)
	for i := len(str) - 1; i > -1; i-- {
		ret = ExprAppl{ExprAppl{ExprFnRef(4), ExprNum(str[i])}, ret}
	}
	return
}
