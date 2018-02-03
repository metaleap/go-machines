package clsyn

import (
	lex "github.com/go-leap/dev/lex"
)

type Keyword func(lex.Tokens) (IExpr, lex.Tokens, *Error)

var Keywords = map[string]Keyword{
	"let":  parseKeywordLet,
	"case": parseKeywordCase,
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

	i, def := 1, &SynDef{Name: tid.Token}
	def.syn.pos = tid.TokenMeta

	// args up until `=`
	for ; i < len(tokens); i++ {
		if t, _ := tokens[i].(*lex.TokenOther); t != nil && t.Token == "=" {
			i++
			break
		} else if t, _ := tokens[i].(*lex.TokenIdent); t != nil {
			def.Args = append(def.Args, t.Token)
		} else {
			return nil, nil, errPos(tokens[i], def.Name+": expected argument name or `=` instead of `"+tokens[i].String()+"`", len(tokens[i].String()))
		}
	}

	// body of definition after `=`
	body, tail := tokens[i:].BreakOnIndent()
	expr, exprerr := parseExpr(body)
	if def.Body = expr; exprerr != nil {
		exprerr.msg = def.Name + ": " + exprerr.msg
	}
	return def, tail, exprerr
}

func parseExpr(toks lex.Tokens) (IExpr, *Error) {
	var lastexpr IExpr
	for len(toks) > 0 {
		var expr IExpr

		if expr == nil { // LIT or IDENT or OP?
			if lit := parseLit(toks[0]); lit != nil {
				expr, toks = lit, toks[1:]
			} else if toth, _ := toks[0].(*lex.TokenOther); toth != nil {
				expr, toks = IdO(toth.Token), toks[1:]
			} else if tid, _ := toks[0].(*lex.TokenIdent); tid != nil {
				if keyword := Keywords[tid.Token]; keyword == nil {
					expr, toks = Id(tid.Token), toks[1:]
				} else if kx, kt, ke := keyword(toks[1:]); ke != nil {
					return nil, ke
				} else {
					expr, toks = kx, kt
				}
			}
		}

		if expr == nil { // LAMBDA?
			if tlam, _ := toks[0].(*lex.TokenOther); tlam != nil && tlam.Token == "\\" {
				lam := Ab(nil, nil)
				lamargs, lambody := toks.BreakOnOther("->")
				if len(lambody) == 0 {
					return nil, errPos(toks[0], "missing body for lambda expression", 0)
				}
				if len(lamargs) == 0 {
					return nil, errPos(toks[0], "missing argument(s) for lambda expression", 0)
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

func parseKeywordLet(toks lex.Tokens) (let IExpr, tail lex.Tokens, err *Error) {
	err = errPos(toks[0], "not yet supported: `let in` keyword", 0)
	return
}

func parseKeywordCase(toks lex.Tokens) (let IExpr, tail lex.Tokens, err *Error) {
	err = errPos(toks[0], "not yet supported: `case of` keyword", 0)
	return
}
