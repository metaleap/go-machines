package corelang

import (
	"fmt"
	"strings"
)

type interpPrettyPrint struct {
	curIndent int
}

func (interpPrettyPrint) prog(me *aProgram, args ...interface{}) (result interface{}, err error) {
	return
}

func (interpPrettyPrint) def(me *aDef, args ...interface{}) (result interface{}, err error) {
	return
}

func (me *interpPrettyPrint) expr(expression iExpr, args ...interface{}) (result interface{}, err error) {
	switch expr := expression.(type) {
	case *aExprSym:
		result = expr.Name
	case *aExprNum:
		result = fmt.Sprint(expr.Lit)
	case *aExprLambda:
		if result, err = me.expr(expr.Body); err != nil {
			result = nil
		} else {
			result = fmt.Sprintf("(\\%s -> %s)", strings.Join(expr.Args, " "), result)
		}
	case *aExprCall:
		var callee, arg interface{}
		if callee, err = me.expr(expr.Callee); err == nil {
			if arg, err = me.expr(expr.Arg); err == nil {
				result = fmt.Sprintf("(%s %s)", callee, arg)
			}
		}
	default:
		err = fmt.Errorf("unknown expression type %T — %#v", expr, expr)
	}
	return
}
