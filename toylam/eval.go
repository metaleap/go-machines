package toylam

import (
	"strconv"
)

type Instr int

const ( // if wanting to re-arrange: `Eval` requires that all arith-instrs do precede all compare-instrs (of which EQ must be the first)
	_ Instr = iota
	InstrADD
	InstrMUL
	InstrSUB
	InstrDIV
	InstrMOD
	InstrMSG
	InstrERR
	InstrEQ
	InstrGT
	InstrLT
)

var (
	instrs     = map[string]Instr{"ADD": InstrADD, "MUL": InstrMUL, "SUB": InstrSUB, "DIV": InstrDIV, "MOD": InstrMOD, "MSG": InstrMSG, "ERR": InstrERR, "EQ": InstrEQ, "GT": InstrGT, "LT": InstrLT}
	instrNames = []string{"?!?", "ADD", "MUL", "SUB", "DIV", "MOD", "MSG", "ERR", "EQ", "GT", "LT"}
)

type Expr interface {
	LocInfo() *Loc
	NamesDeclared() []string
	ReplaceName(string, string) int
	RewriteName(string, Expr) Expr
	String() string
}

type ExprLitNum struct {
	*Loc
	NumVal int
}

func (me *ExprLitNum) NamesDeclared() []string        { return nil }
func (me *ExprLitNum) RewriteName(string, Expr) Expr  { return me }
func (me *ExprLitNum) ReplaceName(string, string) int { return 0 }
func (me *ExprLitNum) String() string                 { return strconv.FormatInt(int64(me.NumVal), 10) }

type ExprName struct {
	*Loc
	NameVal    string
	IdxOrInstr int // if <0 then De Bruijn index, if >0 then instrCode
}

func (me *ExprName) NamesDeclared() []string { return nil }
func (me *ExprName) RewriteName(name string, with Expr) Expr {
	if me.NameVal == name {
		return with
	} else if me.IdxOrInstr < 0 {
		me.IdxOrInstr = 0
	}
	return me
}
func (me *ExprName) ReplaceName(nameOld string, nameNew string) (didReplace int) {
	if me.NameVal == nameOld { // even if nameOld==nameNew, by design: as we use it also to check "refersTo" by doing `ReplaceName("foo", "foo")`
		didReplace, me.NameVal = 1, nameNew
	}
	return
}
func (me *ExprName) String() string { return me.NameVal + ":" + strconv.Itoa(me.IdxOrInstr) }

type ExprCall struct {
	*Loc
	Callee  Expr
	CallArg Expr
}

func (me *ExprCall) NamesDeclared() []string {
	return append(me.Callee.NamesDeclared(), me.CallArg.NamesDeclared()...)
}
func (me *ExprCall) RewriteName(name string, with Expr) Expr {
	me.Callee, me.CallArg = me.Callee.RewriteName(name, with), me.CallArg.RewriteName(name, with)
	return me
}
func (me *ExprCall) ReplaceName(nameOld string, nameNew string) int {
	return me.Callee.ReplaceName(nameOld, nameNew) + me.CallArg.ReplaceName(nameOld, nameNew)
}
func (me *ExprCall) String() string {
	return "(" + me.Callee.String() + " " + me.CallArg.String() + ")"
}

type ExprFunc struct {
	*Loc
	ArgName string
	Body    Expr

	numArgUses int
}

func (me *ExprFunc) NamesDeclared() []string { return append(me.Body.NamesDeclared(), me.ArgName) }
func (me *ExprFunc) RewriteName(name string, with Expr) Expr {
	me.Body = me.Body.RewriteName(name, with)
	return me
}
func (me *ExprFunc) ReplaceName(old string, new string) int { return me.Body.ReplaceName(old, new) }
func (me *ExprFunc) String() string                         { return "{ " + me.ArgName + " -> " + me.Body.String() + " }" }

type Value interface {
	isClosure() *valClosure
	isNum() *valNum
	force() Value
	String() string
}

type valEq interface{ eq(Value) bool }

type valNum int

func (me valNum) eq(cmp Value) bool      { it := cmp.isNum(); return it != nil && me == *it }
func (me valNum) force() Value           { return me }
func (me valNum) isClosure() *valClosure { return nil }
func (me valNum) isNum() *valNum         { return &me }
func (me valNum) String() string         { return strconv.FormatInt(int64(me), 10) }

type valClosure struct {
	env     Values
	body    Expr
	instr   Instr
	argDrop bool
}

func (me *valClosure) force() Value           { return me }
func (me *valClosure) isClosure() *valClosure { return me }
func (me *valClosure) isNum() *valNum         { return nil }
func (me *valClosure) String() (r string) {
	if r = "closureEnv#" + strconv.Itoa(len(me.env)) + "#"; me.body != nil {
		r += me.body.String()
	} else if instr := int(me.instr); instr != 0 && instr < len(instrNames) {
		if instr < 0 {
			instr = 0 - instr
		}
		r += instrNames[instr]
	} else {
		r += "?!NEWBUG!?"
	}
	return
}

type valThunk struct{ val interface{} }

func (me *valThunk) eq(cmp Value) bool {
	self, _ := me.force().(valEq)
	return self != nil && self.eq(cmp.force())
}
func (me *valThunk) isClosure() *valClosure { return me.force().isClosure() }
func (me *valThunk) isNum() *valNum         { return me.force().isNum() }
func (me *valThunk) String() string         { return me.force().String() }
func (me *valThunk) force() (r Value) {
	if r, _ = me.val.(Value); r == nil {
		r = (me.val.(func(*valThunk) Value)(me)).force()
	}
	return
}

type Values []Value

func (me Values) shallowCopy() Values { return append(make(Values, 0, len(me)), me...) }

func (me *Prog) Eval(expr Expr, env Values) Value {
	me.NumEvalSteps++
	switch it := expr.(type) {
	case *ExprLitNum:
		return valNum(it.NumVal)
	case *ExprFunc:
		return &valClosure{body: it.Body, env: env.shallowCopy(), argDrop: it.numArgUses == 0}
	case *ExprName:
		if it.IdxOrInstr > 0 { // it's never 0 thanks to prior & completed `Prog.preResolveNames`
			return &valClosure{instr: Instr(-it.IdxOrInstr), env: env.shallowCopy()}
		} else if it.IdxOrInstr == 0 {
			panic(it.LocStr() + "NEWLY INTRODUCED INTERPRETER BUG: " + it.String())
		}
		return env[len(env)+it.IdxOrInstr]
	case *ExprCall:
		callee := me.Eval(it.Callee, env)
		closure := callee.isClosure()
		if closure == nil {
			panic(it.LocStr() + "not callable: `" + it.Callee.String() + "`, which equates to `" + callee.String() + "`, in: " + it.String())
		}
		var argval Value
		if !closure.argDrop {
			if me.LazyEval {
				argval = &valThunk{val: func(set *valThunk) Value { ret := me.Eval(it.CallArg, env); set.val = ret; return ret }}
			} else {
				argval = me.Eval(it.CallArg, env)
			}
		}
		if closure.instr < 0 {
			closure.instr, closure.env = -closure.instr, append(closure.env, argval) // return &valClosure{instr: -closure.instr, env: append(closure.env.Copy(), argval)}
			return closure
		} else if closure.instr > 0 {
			lhs, rhs := closure.env[len(closure.env)-1], argval
			lnum, rnum := lhs.isNum(), rhs.isNum()
			if iserr := closure.instr == InstrERR; iserr || (closure.instr == InstrMSG) {
				if strmsg := me.Value(lhs).(valFinalBytes); iserr {
					if r := me.Value(rhs).isClosure(); r == nil || r.body == nil || r.body != me.exprId.Body {
						panic(valFinalList{strmsg, me.Value(rhs)})
					}
					panic(strmsg)
				} else if me.OnInstrMSG != nil {
					me.OnInstrMSG(string(strmsg), rhs)
				}
				return argval
			}
			if closure.instr < InstrEQ && lnum != nil && rnum != nil {
				return closure.instr.callCalc(it.LocInfo(), *lnum, *rnum)
			}
			var retbool *ExprFunc
			if closure.instr >= InstrEQ && lnum != nil && rnum != nil {
				retbool = me.newBool(closure.instr.callCmp(it.LocInfo(), *lnum, *rnum))
			} else if closure.instr == InstrEQ {
				if eq, _ := me.value(lhs).(valEq); eq != nil {
					retbool = me.newBool(eq.eq(me.value(rhs)))
				}
			}
			if retbool == nil {
				panic(it.LocStr() + "invalid operands for '" + instrNames[closure.instr] + "' instruction: `" + lhs.String() + "` and `" + rhs.String() + "`, in: `" + it.String() + "`")
			}
			return me.Eval(retbool, env)
		}
		return me.Eval(closure.body, append(closure.env, argval))
	}
	panic(expr)
}

func (me *Prog) newBool(b bool) (exprTrueOrFalse *ExprFunc) {
	if exprTrueOrFalse = me.exprBoolFalse; b {
		exprTrueOrFalse = me.exprBoolTrue
	}
	return
}

func (me Instr) callCalc(loc *Loc, lhs valNum, rhs valNum) valNum {
	switch me {
	case InstrADD:
		return lhs + rhs
	case InstrSUB:
		return lhs - rhs
	case InstrMUL:
		return lhs * rhs
	case InstrDIV:
		return lhs / rhs
	case InstrMOD:
		return lhs % rhs
	}
	panic(loc.LocStr() + "unknown calc-instruction code: " + strconv.Itoa(int(me)))
}

func (me Instr) callCmp(loc *Loc, lhs valNum, rhs valNum) bool {
	switch me {
	case InstrEQ:
		return (lhs == rhs)
	case InstrGT:
		return (lhs > rhs)
	case InstrLT:
		return (lhs < rhs)
	}
	panic(loc.LocStr() + "unknown compare-instruction code: " + strconv.Itoa(int(me)))
}

func Walk(expr Expr, visitor func(Expr)) {
	visitor(expr)
	switch it := expr.(type) {
	case *ExprFunc:
		Walk(it.Body, visitor)
	case *ExprCall:
		Walk(it.Callee, visitor)
		Walk(it.CallArg, visitor)
	default:
		panic(it)
	}
}
