package corelang

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	. "github.com/metaleap/go-corelang/syn"
)

type InterpPrettyPrint struct {
	curIndent int
}

func (me *InterpPrettyPrint) Mod(mod *Module, args ...interface{}) (interface{}, error) {
	var buf bytes.Buffer
	for _, def := range mod.Defs {
		me.curIndent = 0
		me.def(&buf, def)
		buf.WriteString("\n\n")
	}
	return buf.String(), nil
}

func (me *InterpPrettyPrint) def(w *bytes.Buffer, def *Def, _ ...interface{}) {
	w.WriteString(def.Name)
	for _, defarg := range def.Args {
		w.WriteRune(' ')
		w.WriteString(defarg)
	}
	w.WriteString(" =\n")
	me.curIndent++
	w.WriteString(strings.Repeat("  ", me.curIndent))
	me.expr(w, def.Body)
	me.curIndent--
}

func (me *InterpPrettyPrint) Def(def *Def, args ...interface{}) (interface{}, error) {
	var buf bytes.Buffer
	me.def(&buf, def, args...)
	return buf.String(), nil
}

func (me *InterpPrettyPrint) expr(w *bytes.Buffer, expression IExpr) {
	switch expr := expression.(type) {
	case *ExprIdent:
		w.WriteString(expr.Val)
	case *ExprLitFloat:
		w.WriteString(strconv.FormatFloat(expr.Val, 'g', -1, 64))
	case *ExprLambda:
		w.WriteString("(\\")
		for _, lamarg := range expr.Args {
			w.WriteString(lamarg)
			w.WriteRune(' ')
		}
		w.WriteString("-> ")
		me.expr(w, expr.Body)
		w.WriteRune(')')
	case *ExprCall:
		w.WriteRune('(')
		me.expr(w, expr.Callee)
		w.WriteRune(' ')
		me.expr(w, expr.Arg)
		w.WriteRune(')')
	case *ExprLetIn:
		w.WriteString("let ")
		for i, letdef := range expr.Defs {
			me.def(w, letdef)
			if i < (len(expr.Defs) - 1) {
				w.WriteString("; ")
			}
		}
		w.WriteString(" in ")
		me.expr(w, expr.Body)
	case *ExprCtor:
		w.WriteString("Pack{")
		w.WriteString(strconv.Itoa(expr.Tag))
		w.WriteRune(',')
		w.WriteString(strconv.Itoa(expr.Arity))
		w.WriteRune('}')
	case *ExprCaseOf:
		w.WriteString("(case ")
		me.expr(w, expr.Scrut)
		w.WriteString(" of ")
		for i, alt := range expr.Alts {
			w.WriteString(strconv.Itoa(alt.Tag))
			w.WriteString(" -> ")
			me.expr(w, alt.Body)

			if i < (len(expr.Alts) - 1) {
				w.WriteString("; ")
			}
		}
		w.WriteRune(')')
	default:
		panic(fmt.Errorf("unknown expression type %T — %#v", expr, expr))
	}
	return
}

func (me *InterpPrettyPrint) Expr(expr IExpr) (interface{}, error) {
	var buf bytes.Buffer
	me.expr(&buf, expr)
	return buf.String(), nil
}
