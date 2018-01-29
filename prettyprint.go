package corelang

import (
	"fmt"
	"strings"
)

type interpPrettyPrint struct {
	curIndent int
}

func (me *interpPrettyPrint) prog(prog *aProgram, args ...interface{}) (result interface{}, err error) {
	var s string
	for _, def := range prog.Defs {
		if result, err = me.def(def); err != nil {
			return
		}
		s += fmt.Sprintf("%s\n\n", result)
	}
	result = s
	return
}

func (me *interpPrettyPrint) def(def *aDef, args ...interface{}) (result interface{}, err error) {
	s := strings.Join(append([]string{def.Name}, def.Args...), " ")
	if result, err = me.expr(def.Body); err == nil {
		result = fmt.Sprintf("%s = %s", s, result)
	}
	return
}

func (me *interpPrettyPrint) expr(expression iExpr) (result interface{}, err error) {
	switch expr := expression.(type) {
	case *aExprSym:
		result = expr.Name
	case *aExprNum:
		result = fmt.Sprint(expr.Lit)
	case *aExprLambda:
		if result, err = me.expr(expr.Body); err == nil {
			result = fmt.Sprintf("(\\%s -> %s)", strings.Join(expr.Args, " "), result)
		}
	case *aExprCall:
		var callee, arg interface{}
		if callee, err = me.expr(expr.Callee); err == nil {
			if arg, err = me.expr(expr.Arg); err == nil {
				result = fmt.Sprintf("(%s %s)", callee, arg)
			}
		}
	case *aExprLet:
		s := "let"
		for letname, letexpr := range expr.Let {
			if result, err = me.expr(letexpr); err != nil {
				return
			}
			s += fmt.Sprintf(" %s = %s;", letname, result)
		}
		if result, err = me.expr(expr.In); err == nil {
			result = fmt.Sprintf("(%s in %s)", s[:len(s)-1], result)
		}
	case *aExprCtor:
		result = fmt.Sprintf("(Ctor %d %d)", expr.Tag, expr.Arity)
	case *aExprCaseAlt:
		if result, err = me.expr(expr.Body); err == nil {
			result = fmt.Sprintf("%d -> %s", expr.Tag, result)
		}
	case *aExprCase:
		if result, err = me.expr(expr.Scrut); err == nil {
			s := fmt.Sprintf("case %s of ", result)
			for i, alt := range expr.Alts {
				if result, err = me.expr(alt); err != nil {
					return
				}
				if s += fmt.Sprint(result); i < len(expr.Alts)-1 {
					s += "; "
				}
			}
			result = "(" + s + ")"
		}
	default:
		err = fmt.Errorf("unknown expression type %T — %#v", expr, expr)
	}
	return
}
