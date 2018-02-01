package corelang

import (
	"bytes"
	"fmt"
	"strconv"
)

type InterpPrettyPrint struct {
	curIndent int
}

func (me *InterpPrettyPrint) Prog(prog *aProgram, args ...interface{}) (interface{}, error) {
	var buf bytes.Buffer
	for _, def := range prog.Defs {
		me.curIndent = 0
		me.def(&buf, def)
		buf.WriteString("\n\n")
	}
	return buf.String(), nil
}

func (me *InterpPrettyPrint) def(w *bytes.Buffer, def *aDef, _ ...interface{}) {
	w.WriteString(def.Name)
	for _, defarg := range def.Args {
		w.WriteRune(' ')
		w.WriteString(defarg)
	}
	w.WriteString(" = ")
	me.expr(w, def.Body)
}

func (me *InterpPrettyPrint) Def(def *aDef, args ...interface{}) (interface{}, error) {
	var buf bytes.Buffer
	me.def(&buf, def, args...)
	return buf.String(), nil
}

func (me *InterpPrettyPrint) expr(w *bytes.Buffer, expression iExpr) {
	switch expr := expression.(type) {
	case *aExprSym:
		w.WriteString(expr.Name)
	case *aExprNum:
		w.WriteString(strconv.Itoa(expr.Lit))
	case *aExprLambda:
		w.WriteString("(\\")
		for _, lamarg := range expr.Args {
			w.WriteString(lamarg)
			w.WriteRune(' ')
		}
		w.WriteString("-> ")
		me.expr(w, expr.Body)
		w.WriteRune(')')
	case *aExprCall:
		w.WriteRune('(')
		me.expr(w, expr.Callee)
		w.WriteRune(' ')
		me.expr(w, expr.Arg)
		w.WriteRune(')')
	case *aExprLet:
		w.WriteString("let ")
		for i, letdef := range expr.Defs {
			me.def(w, letdef)
			if i < (len(expr.Defs) - 1) {
				w.WriteString("; ")
			}
		}
		w.WriteString(" in ")
		me.expr(w, expr.Body)
	case *aExprCtor:
		w.WriteString("Pack{")
		w.WriteString(strconv.Itoa(expr.Tag))
		w.WriteRune(',')
		w.WriteString(strconv.Itoa(expr.Arity))
		w.WriteRune('}')
	case *aExprCaseAlt:
		w.WriteString(strconv.Itoa(expr.Tag))
		w.WriteString(" -> ")
		me.expr(w, expr.Body)
	case *aExprCase:
		w.WriteString("(case ")
		me.expr(w, expr.Scrut)
		w.WriteString(" of ")
		for i, alt := range expr.Alts {
			me.expr(w, alt)
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

func (me *InterpPrettyPrint) Expr(expr iExpr) (interface{}, error) {
	var buf bytes.Buffer
	me.expr(&buf, expr)
	return buf.String(), nil
}
