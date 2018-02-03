package clsyn

import (
	"strings"

	lex "github.com/go-leap/dev/lex"
)

type Keyword func(lex.Tokens) (IExpr, lex.Tokens, *Error)

var keywords = map[string]Keyword{}

func init() {
	RegisterKeyword("LET", parseKeywordLet)
	RegisterKeyword("CASE", parseKeywordCase)
}

func RegisterKeyword(triggerWord string, keyword Keyword) string {
	if triggerWord = strings.ToUpper(strings.TrimSpace(triggerWord)); triggerWord != "" && keyword != nil && keywords[triggerWord] == nil {
		keywords[triggerWord] = keyword
		return triggerWord
	}
	return ""
}

func Lex(srcFilePath string, src string) (lex.Tokens, []*lex.Error) {
	return lex.Lex(srcFilePath, src, true, "(", ")", "\\")
}

func ParseDefs(srcFilePath string, tokens lex.Tokens) (defs []*SynDef, errs []*Error) {
	defs, errs = parseDefs(tokens)
	for _, e := range errs {
		e.Pos.Filename = srcFilePath
	}
	return
}

func parseDefs(tokens lex.Tokens) (defs []*SynDef, errs []*Error) {
	for len(tokens) > 0 {
		def, tail, deferr := parseDef(tokens)
		if tokens = tail; deferr != nil {
			errs = append(errs, deferr)
		} else {
			defs = append(defs, def)
		}
	}
	return
}

func parseDef(tokens lex.Tokens) (*SynDef, lex.Tokens, *Error) {
	if len(tokens) < 3 {
		return nil, nil, errPos(lex.Pos(tokens, nil, ""), "not enough tokens to form a definition", 0)
	}

	tid, _ := tokens[0].(*lex.TokenIdent)
	if tid == nil {
		return nil, nil, errPos(tokens[0], "expected identifier instead of `"+tokens[0].String()+"`", len(tokens[0].String()))
	}

	toks, tail := tokens[1:].BreakOnIndent(tid.LineIndent)
	if len(toks) < 2 {
		return nil, nil, errPos(tid, "not enough tokens to form a definition", len(tid.Token))
	}

	i, def := 0, &SynDef{Name: tid.Token}
	def.syn.pos = tid.TokenMeta

	// args up until `=`
	for ; i < len(toks); i++ {
		if t, _ := toks[i].(*lex.TokenOther); t != nil && t.Token == "=" {
			i++
			break
		} else if t, _ := toks[i].(*lex.TokenIdent); t != nil {
			def.Args = append(def.Args, t.Token)
		} else {
			return nil, nil, errPos(toks[i], def.Name+": expected argument name or `=` instead of `"+toks[i].String()+"`", len(toks[i].String()))
		}
	}

	// body of definition after `=`
	bodytoks := toks[i:]
	if len(bodytoks) == 0 {
		return nil, nil, errPos(toks[len(toks)-1], def.Name+": missing body of definition", 0)
	}
	expr, exprerr := parseExpr(toks[i:])
	if def.Body = expr; exprerr != nil {
		exprerr.msg = def.Name + ": " + exprerr.msg
	}
	return def, tail, exprerr
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

func parseExpr(toks lex.Tokens) (IExpr, *Error) {
	var lastexpr IExpr
	for len(toks) > 0 {
		var expr IExpr

		if expr == nil { // LAMBDA?
			if tlam, _ := toks[0].(*lex.TokenSep); tlam != nil && tlam.Token == "\\" {
				if toks = toks[1:]; len(toks) == 0 {
					return nil, errPos(tlam, "expected complete lambda abstraction", 0)
				}
				lamargs, lambody := toks.BreakOnOther("->")
				if len(lamargs) == 0 {
					return nil, errPos(toks[0], "missing argument(s) for lambda expression", 0)
				}
				lam := Ab(nil, nil)
				for _, lamarg := range lamargs {
					if tid, _ := lamarg.(*lex.TokenIdent); tid != nil {
						lam.Args = append(lam.Args, tid.Token)
					} else {
						return nil, errPos(lamarg, "expected `->` or identifier for lambda argument instead of `"+lamarg.String()+"`", len(lamarg.String()))
					}
				}
				if len(lambody) == 0 {
					return nil, errPos(toks[0], "missing body for lambda expression", 0)
				}
				lamexpr, lamerr := parseExpr(lambody)
				if lam.Body = lamexpr; lamerr != nil {
					return nil, lamerr
				}
				expr, toks = lam, nil
			}
		}

		if expr == nil { // LIT or IDENT or OP or KEYWORD?
			if lit := parseLit(toks[0]); lit != nil {
				expr, toks = lit, toks[1:]
			} else if toth, _ := toks[0].(*lex.TokenOther); toth != nil {
				expr, toks = IdO(toth.Token), toks[1:]
			} else if tid, _ := toks[0].(*lex.TokenIdent); tid != nil {
				if keyword := keywords[tid.Token]; keyword == nil || len(toks) == 1 {
					expr, toks = Id(tid.Token), toks[1:]
				} else if kx, kt, ke := keyword(toks); ke != nil {
					return nil, ke
				} else {
					expr, toks = kx, kt
				}
			}
		}

		if expr == nil { // PARENS SUB-EXPR?
			if tsep, _ := toks[0].(*lex.TokenSep); tsep != nil && tsep.Token == "(" {
				sub, subtail, numunclosed := toks.SubTokens("(", ")")
				if numunclosed > 0 {
					return nil, errPos(toks[0], "unclosed parentheses in current indent level", 1)
				} else if len(sub) == 0 {
					return nil, errPos(toks[0], "empty or mis-matched parentheses", 0)
				} else if subexpr, suberr := parseExpr(sub); suberr == nil {
					expr, toks = subexpr, subtail
				} else {
					return nil, suberr
				}
			}
		}

		if expr == nil { // should already have early-returned-with-error by now: if this message shows up, indicates earlier validations above are unacceptably not exhaustive
			return nil, errPos(toks[0], "not an expression: "+toks[0].String(), 0)
		} else if lastexpr == nil {
			lastexpr = expr
		} else {
			if ctortag, _ := lastexpr.(*ExprLitUInt); ctortag != nil {
				if ctorarity, _ := expr.(*ExprLitUInt); ctorarity != nil {
					lastexpr = Ct(ctortag.Val, ctorarity.Val)
					continue
				}
			}
			lastexpr = Ap(lastexpr, expr)
		}
	} // big for-loop
	return lastexpr, nil
}

func parseKeywordLet(tokens lex.Tokens) (IExpr, lex.Tokens, *Error) {
	toks := tokens[1:] // tokens[0] is `LET` keyword itself
	defstoks, bodytoks, numunclosed := toks.BreakOnIdent("IN", "LET")
	if numunclosed != 0 || (len(defstoks) == 0 && len(bodytoks) == 0) {
		return nil, nil, errPos(toks[0], "missing `IN` for some `LET`", 0)
	} else if len(defstoks) == 0 {
		return nil, nil, errPos(toks[0], "missing definitions between `LET` and `IN`", 0)
	} else if len(bodytoks) == 0 {
		return nil, nil, errPos(toks[0], "missing expression body following `IN`", 0)
	}

	bodyexpr, bodyerr := parseExpr(bodytoks)
	if bodyerr != nil {
		return nil, nil, bodyerr
	}

	if def0, lkwd := defstoks[0].Meta(), tokens[0].Meta(); def0.Line == lkwd.Line {
		def0.LineIndent += (def0.Column - lkwd.Column) // typically 4, ie. len("LET ")
	}
	defsyns, defserrs := parseDefs(defstoks)
	if len(defserrs) > 0 {
		return nil, nil, defserrs[0]
	}
	return &ExprLetIn{Body: bodyexpr, Defs: defsyns}, nil, nil
}

func parseKeywordCase(tokens lex.Tokens) (let IExpr, tail lex.Tokens, err *Error) {
	toks := tokens[1:] // tokens[0] is `CASE` keyword itself
	scruttoks, altstoks, numunclosed := toks.BreakOnIdent("OF", "CASE")
	if numunclosed != 0 || (len(scruttoks) == 0 && len(altstoks) == 0) {
		return nil, nil, errPos(toks[0], "missing `OF` for some `CASE`", 0)
	} else if len(scruttoks) == 0 {
		return nil, nil, errPos(toks[0], "missing scrutinee between `CASE` and `OF`", 0)
	} else if len(altstoks) == 0 {
		return nil, nil, errPos(toks[0], "missing `CASE` alternatives following `OF`", 0)
	}

	scrutexpr, scruterr := parseExpr(scruttoks)
	if scruterr != nil {
		return nil, nil, scruterr
	}
	caseof := &ExprCaseOf{Scrut: scrutexpr}

	return caseof, nil, errPos(tokens[0], "not yet implemented: `CASE` keyword", 4)
}
