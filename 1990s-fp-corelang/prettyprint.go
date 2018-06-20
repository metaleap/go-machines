package corelang

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	. "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
)

type SyntaxTreePrinter struct {
	curIndent int
}

func (this *SyntaxTreePrinter) Mod(mod *SynMod) string {
	var buf bytes.Buffer
	for _, def := range mod.Defs {
		this.curIndent = 0
		this.def(&buf, def)
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func (this *SyntaxTreePrinter) def(w *bytes.Buffer, def *SynDef) {
	w.WriteString(def.Name)
	for _, defarg := range def.Args {
		w.WriteRune(' ')
		w.WriteString(defarg)
	}
	w.WriteString(" =\n")
	this.curIndent++
	w.WriteString(strings.Repeat("  ", this.curIndent))
	this.expr(w, def.Body, false)
	this.curIndent--
}

func (this *SyntaxTreePrinter) Def(def *SynDef) string {
	var buf bytes.Buffer
	this.def(&buf, def)
	return buf.String()
}

func (this *SyntaxTreePrinter) expr(w *bytes.Buffer, expression IExpr, parensUnlessAtomic bool) {
	if parensUnlessAtomic && !expression.IsAtomic() {
		w.WriteRune('(')
	}
	switch expr := expression.(type) {
	case *ExprIdent:
		if expr.OpLike && expr.OpLone {
			w.WriteRune('(')
		}
		w.WriteString(expr.Name)
		if expr.OpLike && expr.OpLone {
			w.WriteRune(')')
		}
	case *ExprLitFloat:
		w.WriteString(strconv.FormatFloat(expr.Lit, 'g', -1, 64))
	case *ExprLitUInt:
		if expr.Base == 16 {
			w.WriteString("0x")
		} else if expr.Base == 8 {
			w.WriteRune('0')
		}
		w.WriteString(strconv.FormatUint(expr.Lit, expr.Base))
	case *ExprLitText:
		w.WriteString(strconv.Quote(expr.Lit))
	case *ExprLitRune:
		w.WriteString(strconv.QuoteRune(expr.Lit))
	case *ExprLambda:
		w.WriteString("\\")
		for _, lamarg := range expr.Args {
			w.WriteString(lamarg)
			w.WriteRune(' ')
		}
		w.WriteString("-> ")
		this.expr(w, expr.Body, true)
	case *ExprCall:
		this.expr(w, expr.Callee, true)
		w.WriteRune(' ')
		this.expr(w, expr.Arg, true)
	case *ExprLetIn:
		w.WriteString("LET\n")
		this.curIndent++
		for _, letdef := range expr.Defs {
			w.WriteString(strings.Repeat("  ", this.curIndent))
			this.def(w, letdef)
			w.WriteRune('\n')
		}
		this.curIndent--
		w.WriteString(strings.Repeat("  ", this.curIndent))
		w.WriteString("IN\n")
		this.curIndent++
		w.WriteString(strings.Repeat("  ", this.curIndent))
		this.expr(w, expr.Body, false)
		this.curIndent--
	case *ExprCtor:
		w.WriteRune('(')
		w.WriteString(expr.Tag)
		w.WriteRune(' ')
		w.WriteString(strconv.Itoa(expr.Arity))
		w.WriteRune(')')
	case *ExprCaseOf:
		w.WriteString("CASE ")
		this.expr(w, expr.Scrut, false)
		w.WriteString(" OF\n")
		this.curIndent++
		w.WriteString(strings.Repeat("  ", this.curIndent))
		for _, alt := range expr.Alts {
			w.WriteString(alt.Tag)
			for _, bind := range alt.Binds {
				w.WriteRune(' ')
				w.WriteString(bind)
			}
			w.WriteString(" ->\n")
			this.curIndent++
			w.WriteString(strings.Repeat("  ", this.curIndent))
			this.expr(w, alt.Body, false)
			this.curIndent--
			w.WriteRune('\n')
			w.WriteString(strings.Repeat("  ", this.curIndent))
		}
		this.curIndent--
	default:
		panic(fmt.Errorf("unknown expression type %T — %#v", expr, expr))
	}
	if parensUnlessAtomic && !expression.IsAtomic() {
		w.WriteRune(')')
	}
	return
}

func (this *SyntaxTreePrinter) Expr(expr IExpr) string {
	var buf bytes.Buffer
	this.expr(&buf, expr, false)
	return buf.String()
}
