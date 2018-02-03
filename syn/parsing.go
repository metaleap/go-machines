package clsyn

import (
	"text/scanner"

	lex "github.com/go-leap/dev/lex"
)

func errPos(pos lex.IPos, msg string, rangeLen int) *Error {
	return &Error{Pos: pos.Pos().Position, msg: msg, RangeLen: rangeLen}
}

type Error struct {
	msg      string
	Pos      scanner.Position
	RangeLen int
}

func (me *Error) Error() string { return me.msg }

func parseExpr(parent ISyn, tokens lex.Tokens) (IExpr, lex.Tokens, *Error) {
	toks, tail := tokens.BreakOnIndent()
	if len(toks) == 0 {
		return nil, nil, errPos(lex.Pos(nil, parent, ""), "not enough tokens to form an expression under ", 0)
	}

	sub, subtail, numunclosed := toks.SubTokens("(", ")")
	if numunclosed > 0 {
		return nil, nil, errPos(toks[0], "unclosed parens in current indent level", 1)
	} else if len(sub) > 0 {
		subexpr, unlikelytail, suberr := parseExpr(parent, sub)
		if suberr == nil && len(unlikelytail) > 0 {
			return nil, nil, errPos(unlikelytail[0], "dangling tokens in parenthesized sub-expression", 0)
		}
		return subexpr, append(subtail, tail...), suberr
	}

	if tlam, _ := toks[0].(*lex.TokenOther); tlam != nil && tlam.Token == "\\" {
		lam := Ab(nil, nil)
		lamargs, lambody := toks.BreakOnOther("->")
		if lam.parent = parent; len(lambody) == 0 {
			return nil, nil, errPos(toks[0], "missing body for lambda expression", 0)
		}
		for _, lamarg := range lamargs {
			if tid, _ := lamarg.(*lex.TokenIdent); tid != nil {
				lam.Args = append(lam.Args, tid.Token)
			} else {
				return nil, nil, errPos(lamarg, "expected identifier for lambda argument instead of `"+lamarg.String()+"`", len(lamarg.String()))
			}
		}
		lamexpr, lamtail, lamerr := parseExpr(lam, lambody)
		if lam.Body = lamexpr; lamerr != nil {
			return nil, nil, lamerr
		}
		return lam, append(lamtail, tail...), nil
	}

	if len(toks) == 1 {
		if lit := parseLit(parent, toks[0]); lit != nil {
			return lit, append(toks[1:], tail...), nil
		} else if tid, _ := toks[0].(*lex.TokenIdent); tid != nil {
			id := Id(tid.Token)
			id.parent = parent
			return id, nil, nil
		} else if toth, _ := toks[0].(*lex.TokenOther); toth != nil {
			id := IdO(toth.Token)
			id.parent = parent
			return id, nil, nil
		}
		return nil, nil, errPos(toks[0], "expected identifier or literal instead of `"+toks[0].String()+"`", len(toks[0].String()))
	}

	ap := Ap(nil, nil)
	ap.parent = parent
	if apcallee, _, errcallee := parseExpr(ap, toks[:1]); errcallee == nil {
		ap.Callee = apcallee
	} else {
		return nil, nil, errcallee
	}
	aparg, aptail, errarg := parseExpr(ap, toks[1:])
	if errarg != nil {
		return nil, nil, errarg
	}
	ap.Arg = aparg
	return ap, append(aptail, tail...), nil
}

func parseLit(parent ISyn, token lex.IToken) IExpr {
	switch t := token.(type) {
	case *lex.TokenFloat:
		lit := Lf(t.Token)
		lit.parent = parent
		return lit
	case *lex.TokenUint:
		lit := Lu(t.Token, t.Base)
		lit.parent = parent
		return lit
	case *lex.TokenRune:
		lit := Lr(t.Token)
		lit.parent = parent
		return lit
	case *lex.TokenStr:
		lit := Lt(t.Token)
		lit.parent = parent
		return lit
	}
	return nil
}
