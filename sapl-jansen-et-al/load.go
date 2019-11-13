package sapl

import (
	"strconv"
	"strings"
)

type Prog []Expr

func parseInt(s string) int {
	if n, e := strconv.ParseInt(s, 0, 0); e != nil {
		panic(e)
	} else {
		return int(n)
	}
}

func LoadFrom(src string) Prog {
	me := make(Prog, 0, 128)
	for _, ln := range strings.Split(src, "\n") {
		var fn Expr
		var stash []ExprAppl
		for _, tok := range strings.Fields(ln) {
			var cur Expr
			if len(tok) > 2 && tok[0] == '0' && tok[1] == 'x' {
				cur = ExprFunc{2, -parseInt(tok)}
			} else if tok[0] == '0' && len(tok) > 1 {
				cur = ExprVar(parseInt(tok))
			} else if tok[0] >= '0' && tok[0] <= '9' {
				cur = ExprNum(parseInt(tok))
			} else if pos := strings.IndexByte(tok, '@'); pos > 0 {
				cur = ExprFunc{parseInt(tok[:pos]), parseInt(tok[pos+1:])}
			} else if tok == "(" {
				stash = append(stash, ExprAppl{})
			} else if tok == ")" {
				cur, stash = stash[len(stash)-1], stash[:len(stash)-1]
			}
			if fn == nil {
				fn = cur
			} else if len(stash) == 0 {
				fn = ExprAppl{fn, cur}
			} else if appl := &stash[len(stash)-1]; appl.lhs == nil {
				appl.lhs = cur
			} else if appl.rhs == nil {
				appl.rhs = cur
			} else {
				appl.lhs, appl.rhs = *appl, cur
			}
		}
		me = append(me, fn)
	}
	return me
}

func (me ExprNum) String() string { return strconv.Itoa(int(me)) }

func (me ExprVar) String() string { return "@" + strconv.Itoa(int(me)) }

func (me ExprFunc) String() string { return strconv.Itoa(me.numArgs) + "@" + strconv.Itoa(me.idx) }

func (me ExprAppl) String() string { return "(" + me.lhs.String() + " " + me.rhs.String() + ")" }
