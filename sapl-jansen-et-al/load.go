package sapl

import (
	"encoding/json"
	"strconv"
)

type Prog []Expr

func LoadFromJson(src []byte) Prog {
	toplevel := make([]interface{}, 0, 128)
	if e := json.Unmarshal(src, &toplevel); e != nil {
		panic(e)
	}
	me := make(Prog, 0, len(toplevel))
	for _, it := range toplevel {
		me = append(me, exprFromJson(it))
	}
	return me
}

func exprFromJson(from interface{}) Expr {
	switch it := from.(type) {
	case float64:
		return ExprNum(int(it))
	case string:
		if n, e := strconv.ParseInt(it, 10, 0); e != nil {
			panic(e)
		} else {
			return ExprArgRef(int(n))
		}
	case []interface{}:
		if len(it) > 1 {
			expr := exprFromJson(it[0])
			for i := 1; i < len(it); i++ {
				expr = ExprAppl{expr, exprFromJson(it[i])}
			}
			return expr
		} else {
			arr := it[0].([]interface{})
			return ExprFnRef{int(arr[0].(float64)), int(arr[1].(float64))}
		}
	}
	panic(from)
}
