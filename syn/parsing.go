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

func parseExpr(toks lex.Tokens) (IExpr, *Error) {
	var lastexpr IExpr
	for len(toks) > 0 {
		var expr IExpr

		if expr == nil { // LIT or IDENT or OP?
			if lit := parseLit(toks[0]); lit != nil {
				expr, toks = lit, toks[1:]
			} else if tid, _ := toks[0].(*lex.TokenIdent); tid != nil {
				expr, toks = Id(tid.Token), toks[1:]
			} else if toth, _ := toks[0].(*lex.TokenOther); toth != nil {
				expr, toks = IdO(toth.Token), toks[1:]
			}
		}

		if expr == nil { // LAMBDA?
			if tlam, _ := toks[0].(*lex.TokenOther); tlam != nil && tlam.Token == "\\" {
				lam := Ab(nil, nil)
				lamargs, lambody := toks.BreakOnOther("->")
				if len(lambody) == 0 {
					return nil, errPos(toks[0], "missing body for lambda expression", 0)
				}
				for _, lamarg := range lamargs {
					if tid, _ := lamarg.(*lex.TokenIdent); tid != nil {
						lam.Args = append(lam.Args, tid.Token)
					} else {
						return nil, errPos(lamarg, "expected identifier for lambda argument instead of `"+lamarg.String()+"`", len(lamarg.String()))
					}
				}
				lamexpr, lamerr := parseExpr(lambody)
				if lam.Body = lamexpr; lamerr != nil {
					return nil, lamerr
				}
				expr, toks = lam, nil
			}
		}

		if expr == nil { // PARENS SUB-EXPR?
			sub, subtail, numunclosed := toks.SubTokens("(", ")")
			if numunclosed > 0 {
				return nil, errPos(toks[0], "unclosed parens in current indent level", 1)
			} else if len(sub) == 0 {
				return nil, errPos(toks[0], "unrecognized syntax", 0)
			} else if subexpr, suberr := parseExpr(sub); suberr == nil {
				expr, toks = subexpr, subtail
			} else {
				return nil, suberr
			}
		}

		// NEXT
		if lastexpr == nil {
			lastexpr = expr
		} else {
			lastexpr = Ap(lastexpr, expr)
		}
	}

	return lastexpr, nil
}

func parseLit(token lex.IToken) IExpr {
	switch t := token.(type) {
	case *lex.TokenFloat:
		return Lf(t.Token)
	case *lex.TokenUint:
		return Lu(t.Token, t.Base)
	case *lex.TokenRune:
		return Lr(t.Token)
	case *lex.TokenStr:
		return Lt(t.Token)
	}
	return nil
}
