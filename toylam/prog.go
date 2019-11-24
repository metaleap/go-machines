package toylam

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type localDef = struct {
	name string
	Expr
}

type Prog struct {
	LazyEval        bool
	TopDefs         map[string]Expr
	TopDefSepLocals map[string][]localDef
	OnInstrMSG      func(string, Value)
	NumEvalSteps    int

	exprBoolTrue         *ExprFunc
	exprBoolFalse        *ExprFunc
	exprListNil          *ExprFunc
	exprListConsCtorBody Expr
	pseudoSumTypes       map[string][]pseudoSumTypeCtor
}

type pseudoSumTypeCtor = struct {
	name  string
	arity int
}

func (me *Prog) RunAsMain(mainFuncExpr Expr, osProcArgs []string) (ret Value) {
	loc, expr2eval := mainFuncExpr.LocInfo(), mainFuncExpr
	fillarg := func(argval Expr) Expr { return &ExprCall{loc, expr2eval, argval} }
	for fn, _ := mainFuncExpr.(*ExprFunc); fn != nil; fn, _ = fn.Body.(*ExprFunc) {
		if fn.numArgUses == 0 {
			fn.numArgUses = -1
		}
		if fn.ArgName == "args" {
			expr2eval = fillarg(me.newListOfStrs(false, loc, osProcArgs))
		} else if fn.ArgName == "env" {
			expr2eval = fillarg(me.newListOfStrs(false, loc, os.Environ()))
		} else if fn.ArgName == "stdin" {
			if stdin, err := ioutil.ReadAll(os.Stdin); err != nil {
				panic(err)
			} else {
				expr2eval = fillarg(me.newStr(false, loc, string(stdin)))
			}
		} else {
			panic("no idea how to fill in `" + fn.ArgName + "` arg: for your entry-point-ish top-level defs to be `Prog.RunAsMain`, use any combination of these supported arg names: `stdin`, `env`, `args`.")
		}
	}
	ret = me.Value(me.Eval(expr2eval, nil))
	return
}

func (me *Prog) newStr(forceNames bool, loc *Loc, str string) Expr {
	return me.newList(forceNames, loc, len(str), func(i int) Expr { return &ExprLitNum{loc, int(str[i])} })
}

func (me *Prog) newList(forceNames bool, loc *Loc, length int, next func(int) Expr) (list Expr) {
	cons := me.TopDefs[StdRequiredDefs_listCons]
	if cons == nil || forceNames {
		cons = &ExprName{loc, StdRequiredDefs_listCons, 0}
	}
	if list = me.TopDefs[StdRequiredDefs_listNil]; list == nil || forceNames {
		list = &ExprName{loc, StdRequiredDefs_listNil, 0}
	}
	for i := length - 1; i >= 0; i-- {
		list = &ExprCall{loc, &ExprCall{loc, cons, next(i)}, list}
	}
	return
}

func (me *Prog) newListOfStrs(forceNames bool, loc *Loc, vals []string) Expr {
	return me.newList(forceNames, loc, len(vals), func(i int) Expr { return me.newStr(forceNames, loc, vals[i]) })
}
func (me *Prog) value(it Value) Value {
	if cl := it.isClosure(); cl != nil {
		if cl.body == me.exprListNil.Body {
			it = valFinalList(nil)
		} else if cl.body == me.exprListConsCtorBody {
			it = &valTempCons{me.value(cl.env[len(cl.env)-2]), me.value(cl.env[len(cl.env)-1])}
		}
	}
	return it
}

func (me *Prog) Value(it Value) (retVal Value) {
	retVal = me.value(it)
	if cons, _ := retVal.(*valTempCons); cons != nil {
		list := make(valFinalList, 0, 128)
		for cons != nil {
			list = append(list, me.Value(cons.head))
			if eol, ok := cons.tail.(valFinalList); ok && eol == nil {
				cons = nil
			} else if next, _ := cons.tail.(*valTempCons); next != nil {
				cons = next
			} else {
				cons = &valTempCons{me.Value(cons.tail), valFinalList(nil)}
			}
		}
		retVal = list

		if list[0].isNum() != nil && list[len(list)-1].isNum() != nil {
			allbytes := make(valFinalBytes, len(list))
			for i := 0; i < len(allbytes); i++ {
				if n := list[i].isNum(); n != nil && *n > 0 && *n < 256 {
					allbytes[i] = byte(*n)
				} else {
					allbytes = nil
				}
			}
			if allbytes != nil {
				retVal = allbytes
			}
		}
	} else if cl := retVal.isClosure(); cl != nil {
		var name *ExprName
		var args []*ExprName
		for name, _ = cl.body.(*ExprName); name == nil && cl.body != nil; name, _ = cl.body.(*ExprName) {
			if call, _ := cl.body.(*ExprCall); call != nil {
				cl.body = call.Callee
				if argname, _ := call.CallArg.(*ExprName); argname == nil {
					break
				} else {
					args = append(args, argname)
				}
			} else if fn, _ := cl.body.(*ExprFunc); fn != nil {
				cl.body = fn.Body
			}
		}
		if str := ""; name != nil && strings.HasPrefix(name.NameVal, "__") && strings.Contains(name.NameVal, "_Of_") {
			if str = strings.TrimPrefix(name.NameVal[strings.Index(name.NameVal, "_Of_")+len("_Of_"):], "__"); len(args) > 0 {
				str = "(" + str
				for i := len(args) - 1; i >= 0; i-- {
					if idx := (len(cl.env) - 1) - i; idx >= 0 && idx < len(cl.env) && cl.env[idx] != nil {
						str += " " + me.Value(cl.env[idx]).String()
					} else {
						str += " ‹?›"
					}
				}
				str += ")"
			}
			retVal = valFinalOther(str)
		}
	}
	return
}

func ValueBool(it Value) (bool, bool) {
	v, _ := it.(valFinalOther)
	return v == "True", v == "True" || v == "False"
}

func ValueBytes(it Value) ([]byte, bool) {
	v, ok := it.(valFinalBytes)
	return []byte(v), ok
}

func ValueOther(it Value) (string, bool) {
	v, ok := it.(valFinalOther)
	return string(v), ok
}

func ValueNum(it Value) (int, bool) {
	v, ok := it.(valNum)
	return int(v), ok
}

func ValueSlice(it Value) (Values, bool) {
	v, ok := it.(valFinalList)
	return Values(v), ok
}

type valFinalOther string

func (me valFinalOther) force() Value           { return me }
func (me valFinalOther) isClosure() *valClosure { return nil }
func (me valFinalOther) isNum() *valNum         { return nil }
func (me valFinalOther) String() string         { return string(me) }

type valFinalBytes []byte

func (me valFinalBytes) force() Value           { return me }
func (me valFinalBytes) isClosure() *valClosure { return nil }
func (me valFinalBytes) isNum() *valNum         { return nil }
func (me valFinalBytes) String() string         { return strconv.Quote(string(me)) }

type valFinalList Values

func (me valFinalList) eq(cmp Value) bool {
	it, ok := cmp.(valFinalList)
	return ok && len(me) == 0 && len(it) == 0
}
func (me valFinalList) force() Value {
	for i, v := range me {
		me[i] = v.force()
	}
	return me
}
func (me valFinalList) isClosure() *valClosure { return nil }
func (me valFinalList) isNum() *valNum         { return nil }
func (me valFinalList) String() string {
	str := "["
	for i, pref := 0, ""; i < len(me); i, pref = i+1, ", " {
		str = str + pref + me[i].String()
	}
	return str + "]"
}

type valTempCons struct {
	head Value
	tail Value
}

func (me *valTempCons) eq(cmp Value) bool {
	if it, ok := cmp.(*valTempCons); ok {
		if head, okh := me.head.(valEq); okh {
			if tail, okt := me.tail.(valEq); okt {
				return head.eq(it.head) && tail.eq(it.tail)
			}
		}
	}
	return false
}
func (me *valTempCons) force() Value           { me.head.force(); me.tail.force(); return me }
func (me *valTempCons) isClosure() *valClosure { return nil }
func (me *valTempCons) isNum() *valNum         { return nil }
func (me *valTempCons) String() string {
	return "(+> " + me.head.String() + " " + me.tail.String() + ")"
}
