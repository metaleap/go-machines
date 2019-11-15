package tinylam

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Prog struct {
	LazyEval     bool
	TopDefs      map[string]Expr
	OnInstrMSG   func(string, Value)
	NumEvalSteps int

	exprBoolTrue         *ExprFunc
	exprBoolFalse        *ExprFunc
	exprListNil          *ExprFunc
	exprListConsCtorBody Expr
}

func (me *Prog) RunAsMain(mainFuncExpr Expr, osProcArgs []string) (ret Value) {
	loc, expr2eval := mainFuncExpr.locInfo(), mainFuncExpr
	fillarg := func(argval Expr) Expr { return &ExprCall{loc, expr2eval, argval} }
	for fn, _ := mainFuncExpr.(*ExprFunc); fn != nil; fn, _ = fn.Body.(*ExprFunc) {

		if fn.ArgName == "args" {
			expr2eval = fillarg(me.newListOfStrs(loc, osProcArgs))

		} else if fn.ArgName == "env" {
			expr2eval = fillarg(me.newListOfStrs(loc, os.Environ()))

		} else if fn.ArgName == "stdin" {
			if stdin, err := ioutil.ReadAll(os.Stdin); err != nil {
				panic(err)
			} else {
				expr2eval = fillarg(me.newStr(loc, string(stdin)))
			}

		} else {
			panic("no idea how to fill in `" + fn.ArgName + "` arg: for your entry-point-ish top-level defs to be `Prog.RunAsMain`, use any combination of these supported arg names: `stdin`, `env`, `args`.")
		}
	}
	return me.Value(me.Eval(expr2eval, nil))
}

func (me *Prog) newCons(loc *nodeLocInfo, head Expr, tail Expr) Expr {
	cons := me.TopDefs[StdRequiredDefs_listCons]
	if cons == nil {
		cons = &ExprName{loc, StdRequiredDefs_listCons, 0}
	}
	return &ExprCall{loc, &ExprCall{loc, cons, head}, tail}
}

func (me *Prog) newStr(loc *nodeLocInfo, str string) Expr {
	return me.newList(loc, len(str), func(i int) Expr { return &ExprLitNum{loc, int(str[i])} })
}

func (me *Prog) newList(loc *nodeLocInfo, length int, next func(int) Expr) (list Expr) {
	if list = me.TopDefs[StdRequiredDefs_listNil]; list == nil {
		list = &ExprName{loc, StdRequiredDefs_listNil, 0}
	}
	for i := length - 1; i >= 0; i-- {
		list = me.newCons(loc, next(i), list)
	}
	return
}

func (me *Prog) newListOfStrs(loc *nodeLocInfo, vals []string) Expr {
	return me.newList(loc, len(vals), func(i int) Expr { return me.newStr(loc, vals[i]) })
}

func (me *Prog) value(it Value) Value {
	if cl := it.isClosure(); cl != nil {
		if isfalse, istrue := (cl.body == me.exprBoolFalse.Body), (cl.body == me.exprBoolTrue.Body); isfalse || istrue {
			it = valFinalBool(istrue)
		} else if cl.body == me.exprListNil.Body {
			it = valFinalList(nil)
		} else if cl.body == me.exprListConsCtorBody {
			it = &valTempCons{me.value(cl.env[len(cl.env)-2]), me.value(cl.env[len(cl.env)-1])}
		} else if fn, _ := cl.body.(*ExprFunc); fn != nil && strings.HasPrefix(fn.ArgName, "__") && strings.Contains(fn.ArgName, "Of") {
			println(fn.ArgName + fmt.Sprintf("\t%T", cl.body) + "=\t" + cl.body.String() + "\nENV\t" + fmt.Sprintf("%v", cl.env))
			for fnsub, _ := fn.Body.(*ExprFunc); fnsub != nil; fnsub, _ = fn.Body.(*ExprFunc) {
				fn = fnsub
			}
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
	}
	return
}

func ValueBool(it Value) (bool, bool) {
	v, ok := it.(valFinalBool)
	return bool(v), ok
}

func ValueBytes(it Value) ([]byte, bool) {
	v, ok := it.(valFinalBytes)
	return []byte(v), ok
}

func ValueNum(it Value) (int, bool) {
	v, ok := it.(valNum)
	return int(v), ok
}

func ValueSlice(it Value) (Values, bool) {
	v, ok := it.(valFinalList)
	return Values(v), ok
}

type valFinalBool bool

func (me valFinalBool) eq(cmp Value) bool      { it, ok := cmp.(valFinalBool); return ok && me == it }
func (me valFinalBool) force() Value           { return me }
func (me valFinalBool) isClosure() *valClosure { return nil }
func (me valFinalBool) isNum() *valNum         { return nil }
func (me valFinalBool) String() string         { return strconv.FormatBool(bool(me)) }

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
	for i, pref := 0, ""; i < len(me); i, pref = i+1, " " {
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
