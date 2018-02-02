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

func (me *InterpPrettyPrint) Mod(mod *SynMod, args ...interface{}) (interface{}, error) {
	var buf bytes.Buffer
	for _, def := range mod.Defs {
		me.curIndent = 0
		me.def(&buf, def)
		buf.WriteString("\n\n")
	}
	return buf.String(), nil
}

func (me *InterpPrettyPrint) def(w *bytes.Buffer, def *SynDef, _ ...interface{}) {
	w.WriteString(def.Name)
	for _, defarg := range def.Args {
		w.WriteRune(' ')
		w.WriteString(defarg)
	}
	w.WriteString(" =\n")
	me.curIndent++
	w.WriteString(strings.Repeat("  ", me.curIndent))
	me.expr(w, def.Body, false)
	me.curIndent--
}

func (me *InterpPrettyPrint) Def(def *SynDef, args ...interface{}) (interface{}, error) {
	var buf bytes.Buffer
	me.def(&buf, def, args...)
	return buf.String(), nil
}

func (me *InterpPrettyPrint) expr(w *bytes.Buffer, expression IExpr, couldBeParensed bool) {
	if couldBeParensed && !expression.IsAtomic() {
		w.WriteRune('(')
	}
	switch expr := expression.(type) {
	case *ExprIdent:
		w.WriteString(expr.Val)
	case *ExprLitFloat:
		w.WriteString(strconv.FormatFloat(expr.Val, 'g', -1, 64))
	case *ExprLambda:
		w.WriteString("\\")
		for _, lamarg := range expr.Args {
			w.WriteString(lamarg)
			w.WriteRune(' ')
		}
		w.WriteString("-> ")
		me.expr(w, expr.Body, true)
	case *ExprCall:
		me.expr(w, expr.Callee, true)
		w.WriteRune(' ')
		me.expr(w, expr.Arg, true)
	case *ExprLetIn:
		w.WriteString("let\n")
		me.curIndent++
		w.WriteString(strings.Repeat("  ", me.curIndent))
		for _, letdef := range expr.Defs {
			me.def(w, letdef)
			w.WriteRune('\n')
			w.WriteString(strings.Repeat("  ", me.curIndent))
		}
		me.curIndent--
		w.WriteString("in ")
		me.expr(w, expr.Body, false)
	case *ExprCtor:
		w.WriteString("Pack{")
		w.WriteString(strconv.Itoa(expr.Tag))
		w.WriteRune(',')
		w.WriteString(strconv.Itoa(expr.Arity))
		w.WriteRune('}')
	case *ExprCaseOf:
		w.WriteString("case ")
		me.expr(w, expr.Scrut, false)
		w.WriteString(" of\n")
		me.curIndent++
		w.WriteString(strings.Repeat("  ", me.curIndent))
		for _, alt := range expr.Alts {
			w.WriteString(strconv.Itoa(alt.Tag))
			w.WriteString(" ->\n")
			me.curIndent++
			w.WriteString(strings.Repeat("  ", me.curIndent))
			me.expr(w, alt.Body, false)
			me.curIndent--
		}
		me.curIndent--
	default:
		panic(fmt.Errorf("unknown expression type %T — %#v", expr, expr))
	}
	if couldBeParensed && !expression.IsAtomic() {
		w.WriteRune(')')
	}
	return
}

func (me *InterpPrettyPrint) Expr(expr IExpr) (interface{}, error) {
	var buf bytes.Buffer
	me.expr(&buf, expr, false)
	return buf.String(), nil
}
