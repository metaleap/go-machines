// SAPL interpreter implementation following: **"Efficient Interpretation by Transforming Data Types and Patterns to Functions"** (Jan Martin Jansen, Pieter Koopman, Rinus Plasmeijer)
//
// More specifically, implementation in Go of the elegantly slim spec on pages 8-9 (chapter 1.4), excluding all the optimizations starting from section 1.4.1 (p.9ff). No GC / heap / dump etc, stack-only. Go does GC anyway.
//
// Divergence from the paper: NumArgs is not carried around with the Func Ref but stored in the top-level-funcs array together with that func's expression.
//
// "Non"-Parser loads from a simple JSON arrays-of-arrays format: no need to expressly spec out the details here, it's under 40 LoC in the `LoadFromJson` and `exprFromJson` funcs. See the `sapltest/foo.json` that can be fed into the `sapltest/main.go` program via `stdin`.
package sapl

import (
	"encoding/json"
	"strconv"
)

type TopDef = struct {
	NumArgs int
	Expr    Expr
}

type Prog []TopDef

type Expr interface{ String() string }

func (me ExprNum) String() string    { return strconv.Itoa(int(me)) }
func (me ExprArgRef) String() string { return "@" + strconv.Itoa(int(me)) }
func (me ExprFnRef) String() string  { return "^" + strconv.Itoa(int(me)) }
func (me ExprAppl) String() string   { return "(" + me.Callee.String() + " " + me.Arg.String() + ")" }

type ExprNum int

type ExprArgRef int

type ExprFnRef int

type ExprAppl struct {
	Callee Expr
	Arg    Expr
}

type any = interface{}

func LoadFromJson(src []byte) Prog {
	arr := make([][]any, 0, 128)
	if e := json.Unmarshal(src, &arr); e != nil {
		panic(e)
	}
	me := make(Prog, 0, len(arr))
	for _, it := range arr {
		me = append(me, TopDef{int(it[0].(float64)), exprFromJson(it[1])})
	}
	return me
}

func exprFromJson(from any) Expr {
	switch it := from.(type) {
	case float64:
		return ExprNum(int(it))
	case string:
		if n, e := strconv.ParseInt(it, 10, 0); e != nil {
			panic(e)
		} else { // rewrite arg-refs for later stack-access-from-tail-end: 0 -> -1, 1 -> -2, 2 -> -3
			return ExprArgRef(int(-(n + 1)))
		}
	case []any:
		if len(it) == 0 {
			return ExprFnRef(OpPanic) // only `forced panic` mechanism we have in the JSON syntax for now
		} else if len(it) == 1 {
			return ExprFnRef(int(it[0].(float64)))
		}
		expr := exprFromJson(it[0])
		for i := 1; i < len(it); i++ {
			expr = ExprAppl{expr, exprFromJson(it[i])}
		}
		return expr
	case map[string]any: // allows for free-form annotations / comments / meta-data like orig-source-file/line-number mappings...
		return exprFromJson(it[""]) // ... by digging into this field and ignoring all others
	}
	panic(from)
}
