package clsyn

import (
	"strings"

	lex "github.com/go-leap/dev/lex"
	"github.com/go-leap/str"
)

type Keyword func(lex.Tokens) (IExpr, lex.Tokens, *Error)

var keywords = map[string]Keyword{}

func init() {
	RegisterKeyword("LET", parseKeywordLet)
	RegisterKeyword("CASE", parseKeywordCase)
	lex.RestrictedWhitespace = true
	lex.StandaloneSeps = []string{"(", ")"}
}

func RegisterKeyword(triggerWord string, keyword Keyword) string {
	if triggerWord = strings.ToUpper(strings.TrimSpace(triggerWord)); triggerWord != "" && keyword != nil && keywords[triggerWord] == nil {
		keywords[triggerWord] = keyword
		return triggerWord
	}
	return ""
}

func Lex(srcFilePath string, src string) (lex.Tokens, []*lex.Error) {
	return lex.Lex(srcFilePath, strings.NewReader(src), 0, 0, len(src)/6)
}

func LexAndParseDefs(srcFilePath string, src string) ([]*SynDef, []*Error) {
	toks, lexerrs := Lex(srcFilePath, src)
	if len(lexerrs) == 0 {
		return ParseDefs(srcFilePath, toks.SansComments(nil))
	}

	errs := make([]*Error, len(lexerrs))
	for i, lexerr := range lexerrs {
		errs[i] = errPos(&lexerr.Pos, lexerr.Error(), 0)
	}
	return nil, errs
}

func ParseDefs(srcFilePath string, tokens lex.Tokens) (defs []*SynDef, errs []*Error) {
	defs, errs = parseDefs(tokens, true)
	// for _, def := range defs {
	// 	freevars := map[string]bool{}
	// 	def.FreeVars(freevars, NewLookupEnv(defs, nil, nil, nil))
	// 	for name := range freevars {
	// 		errs = append(errs, errTok(&def.Toks()[0], "undefined: "+name))
	// 	}
	// }
	for _, e := range errs {
		e.Pos.Filename = srcFilePath
	}
	return
}

func parseDefs(tokens lex.Tokens, topLevel bool) (defs []*SynDef, errs []*Error) {
	for len(tokens) > 0 {
		def, tail, deferr := parseDef(tokens)
		if tokens = tail; deferr != nil {
			errs = append(errs, deferr)
		} else {
			def.TopLevel, defs = topLevel, append(defs, def)
		}
	}
	return
}

func parseDef(tokens lex.Tokens) (*SynDef, lex.Tokens, *Error) {
	if tokens[0].Kind() != lex.TOKEN_IDENT {
		return nil, nil, errTok(&tokens[0], "expected identifier instead of `"+tokens[0].String()+"`")
	} else if len(tokens) == 1 {
		return nil, nil, errTok(&tokens[0], tokens[0].Str+": expected argument name(s) or `=` next")
	} else if len(tokens) == 2 {
		return nil, nil, errTok(&tokens[1], tokens[0].Str+": expected definition body next")
	}

	toks, tail := tokens[1:].BreakOnIndent(tokens[0].Meta.LineIndent)
	if len(toks) < 2 {
		return nil, nil, errTok(&tokens[0], tokens[0].Str+": incomplete definition (possibly mal-indentation)")
	}

	i, def := 0, &SynDef{Name: tokens[0].Str}
	def.init(toks)

	// args up until `=`
	for inargs := true; inargs && i < len(toks); i++ {
		if tkind := toks[i].Kind(); tkind == lex.TOKEN_OTHER && toks[i].Str == "=" {
			inargs = false
		} else if tkind == lex.TOKEN_IDENT {
			def.Args = append(def.Args, toks[i].Str)
		} else {
			return nil, tail, errTok(&toks[i], def.Name+": expected argument name or `=` instead of `"+toks[i].String()+"`")
		}
	}

	// body of definition after `=`
	bodytoks := toks[i:]
	if len(bodytoks) == 0 {
		return nil, tail, errTok(&toks[len(toks)-1], def.Name+": missing body of definition")
	}
	expr, exprerr := parseExpr(toks[i:])
	if def.Body = expr; exprerr != nil {
		exprerr.msg = def.Name + ": " + exprerr.msg
	}
	return def, tail, exprerr
}

func parseExpr(toks lex.Tokens) (IExpr, *Error) {
	var prevexpr IExpr

	for len(toks) > 0 {
		var thisexpr IExpr
		var thistoks lex.Tokens // always set together with thisexpr

		// LAMBDA?
		if toks[0].Kind() == lex.TOKEN_OTHER && toks[0].Str == "\\" {
			if toks = toks[1:]; len(toks) == 0 {
				return nil, errTok(&toks[0], "expected complete lambda abstraction")
			}
			lamargs, lambody := toks.BreakOnOther("->")
			if len(lamargs) == 0 {
				return nil, errTok(&toks[0], "missing argument(s) for lambda expression")
			} else if len(lambody) == 0 {
				return nil, errTok(&toks[0], "missing body for lambda expression")
			}
			lam := Ab(nil, nil)
			for i := 0; i < len(lamargs); i++ {
				if lamargs[i].Kind() == lex.TOKEN_IDENT {
					lam.Args = append(lam.Args, lamargs[i].Str)
				} else {
					return nil, errTok(&lamargs[i], "expected `->` or identifier for lambda argument instead of `"+lamargs[i].String()+"`")
				}
			}
			lamexpr, lamerr := parseExpr(lambody)
			if lam.Body = lamexpr; lamerr != nil {
				return nil, lamerr
			}
			thistoks, toks, thisexpr = toks, nil, lam
		}

		if thisexpr == nil { // single-token cases: LIT or OP or IDENT/KEYWORD?
			switch toks[0].Kind() {
			case lex.TOKEN_FLOAT:
				thistoks, toks, thisexpr = toks[:1], toks[1:], Lf(toks[0].Float)
			case lex.TOKEN_UINT:
				thistoks, toks, thisexpr = toks[:1], toks[1:], Lu(toks[0].Uint, toks[0].UintBase())
			case lex.TOKEN_RUNE:
				thistoks, toks, thisexpr = toks[:1], toks[1:], Lr(toks[0].Rune())
			case lex.TOKEN_STR:
				thistoks, toks, thisexpr = toks[:1], toks[1:], Lt(toks[0].Str)
			case lex.TOKEN_OTHER: // any operator/separator/punctuation sequence other than "(" and ")"
				thistoks, toks, thisexpr = toks[:1], toks[1:], Op(toks[0].Str, len(toks) == 1)
			case lex.TOKEN_IDENT:
				if keyword := keywords[toks[0].Str]; keyword == nil || len(toks) == 1 {
					thistoks, toks, thisexpr = toks[:1], toks[1:], Id(toks[0].Str)
				} else if kx, kt, ke := keyword(toks); ke != nil {
					return nil, ke
				} else {
					thistoks, toks, thisexpr = toks[:len(toks)-len(kt)], kt, kx
				}
			}
		}

		if thisexpr == nil { // PARENSED SUB-EXPR?
			if toks[0].Kind() == lex.TOKEN_SEP && toks[0].Str == "(" {
				sub, subtail, numunclosed := toks.Sub("(", ")")
				if numunclosed != 0 {
					return nil, errTok(&toks[0], "unclosed parentheses in current indent level")
				} else if len(sub) == 0 {
					return nil, errTok(&toks[0], "empty or mis-matched parentheses")
				} else if subexpr, suberr := parseExpr(sub); suberr == nil {
					thistoks, toks, thisexpr = subexpr.Toks(), subtail, subexpr
				} else {
					return nil, suberr
				}
			}
		}

		if thisexpr == nil { // should already have early-returned-with-error by now: if this message shows up, indicates earlier validations above are unacceptably non-exhaustive
			return nil, errTok(&toks[0], "not an expression: "+toks[0].String())
		} else if thisexpr.init(thistoks); prevexpr == nil {
			prevexpr = thisexpr
		} else {
			// at this point, the only sensible way in corelang to joint prev and cur expr is by application:

			// special case, ctor? any appl form akin to (intlit intlit) is parsed as: Ctor{tag,arity} instead of application
			if ctortag, _ := prevexpr.(*ExprIdent); ctortag != nil && (ustr.BeginsUpper(ctortag.Name) || ctortag.Name == "_") {
				if ctorarity, _ := thisexpr.(*ExprLitUInt); ctorarity != nil {
					prevexpr = Ct(ctortag.Name, ctorarity.Lit)
					prevexpr.init(append(ctortag.toks, ctorarity.toks...)) // TODO: see comment below
					continue
				}
			}

			bothtoks := append(prevexpr.Toks(), thisexpr.Toks()...) // TODO: not nice --- so far, for all syns except Ap and Ct, we could do without extra allocations, reusing the single incoming Tokens slice via sub-slices
			// special case, infix op? any appl infix form of (expr op) is flipped to prefix form (op expr) --- precedence/associativity dont exist in corelang and are simply forced via parens — auto-inserting them via precedence etc being a matter of a later higher-level desugarer
			if exop, _ := thisexpr.(*ExprIdent); exop != nil && exop.OpLike && !exop.OpLone {
				prevexpr = Ap(thisexpr, prevexpr)
			} else if _, isid := prevexpr.(*ExprIdent); (!isid) && prevexpr.IsAtomic() {
				return nil, errTok(&prevexpr.Toks()[0], "atomic literal "+prevexpr.Toks()[0].String()+" cannot be applied like a function")
			} else {
				// default case: apply aka. (prev cur)
				prevexpr = Ap(prevexpr, thisexpr)
			}
			prevexpr.init(bothtoks)
		}
	} // big for-loop
	return prevexpr, nil
}

func parseKeywordLet(tokens lex.Tokens) (IExpr, lex.Tokens, *Error) {
	isrec, toks := false, tokens[1:] // tokens[0] is `LET` keyword itself

	if toks[0].Kind() == lex.TOKEN_IDENT && toks[0].Str == "REC" {
		isrec, toks = true, toks[1:]
	}

	defstoks, bodytoks, numunclosed := toks.BreakOnIdent("IN", "LET")
	if nodef, nobod := len(defstoks) == 0, len(bodytoks) == 0; (nodef && nobod) || numunclosed != 0 {
		return nil, nil, errTok(&toks[0], "a `LET` is missing a corresponding `IN`")
	} else if nodef {
		return nil, nil, errTok(&toks[0], "missing definitions between `LET` and `IN`")
	} else if nobod {
		return nil, nil, errTok(&toks[0], "missing expression body following `IN`")
	}

	bodyexpr, bodyerr := parseExpr(bodytoks)
	if bodyerr != nil {
		return nil, nil, bodyerr
	}

	if def0 := &defstoks[0].Meta; def0.Line == tokens[0].Meta.Line { // first def on same line as LET?
		def0.LineIndent = def0.Column
	} else if isrec && def0.Line == tokens[1].Meta.Line { // or on same line as REC?
		def0.LineIndent = def0.Column
	}
	defsyns, deferrs := parseDefs(defstoks, false)
	if len(deferrs) > 0 {
		return nil, nil, deferrs[0]
	}

	letin := &ExprLetIn{Body: bodyexpr, Defs: defsyns, Rec: isrec}
	letin.init(tokens)
	return letin, nil, nil
}

func parseKeywordCase(tokens lex.Tokens) (IExpr, lex.Tokens, *Error) {
	toks := tokens[1:] // tokens[0] is `CASE` keyword itself

	scruttoks, altstoks, numunclosed := toks.BreakOnIdent("OF", "CASE")
	if numunclosed != 0 || (len(scruttoks) == 0 && len(altstoks) == 0) {
		return nil, nil, errTok(&toks[0], "a `CASE` is missing a corresponding `OF`")
	} else if len(scruttoks) == 0 {
		return nil, nil, errTok(&toks[0], "missing scrutinee between `CASE` and `OF`")
	} else if len(altstoks) == 0 {
		return nil, nil, errTok(&toks[0], "missing `CASE` alternatives following `OF`")
	}

	scrutexpr, scruterr := parseExpr(scruttoks)
	if scruterr != nil {
		return nil, nil, scruterr
	}

	if alt0 := &altstoks[0].Meta; alt0.Line == tokens[0].Meta.Line {
		alt0.LineIndent = alt0.Column
	}
	altsyns, alterrs := parseKeywordCaseAlts(altstoks)
	if len(alterrs) > 0 {
		return nil, nil, alterrs[0]
	}

	caseof := &ExprCaseOf{Scrut: scrutexpr}
	caseof.init(tokens)
	caseof.Alts = altsyns
	return caseof, nil, nil
}

func parseKeywordCaseAlts(tokens lex.Tokens) (alts []*SynCaseAlt, errs []*Error) {
	for len(tokens) > 0 {
		alt, tail, alterr := parseKeywordCaseAlt(tokens)
		if tokens = tail; alterr != nil {
			errs = append(errs, alterr)
		} else {
			alts = append(alts, alt)
		}
	}
	return
}

func parseKeywordCaseAlt(tokens lex.Tokens) (*SynCaseAlt, lex.Tokens, *Error) {
	if tokens[0].Kind() != lex.TOKEN_IDENT {
		return nil, nil, errTok(&tokens[0], "expected constructor tag instead of `"+tokens[0].String()+"`")
	} else if len(tokens) == 1 {
		return nil, nil, errTok(&tokens[0], "expected name(s) or `->` next")
	} else if len(tokens) == 2 {
		return nil, nil, errTok(&tokens[1], "expected `CASE`-alternative body next")
	}

	toks, tail := tokens[1:].BreakOnIndent(tokens[0].Meta.LineIndent)
	if len(toks) < 2 {
		return nil, nil, errTok(&tokens[0], "incomplete `CASE` alternative (possibly mal-indentation)")
	}

	i, alt := 0, &SynCaseAlt{Tag: tokens[0].Str}
	alt.init(toks)

	// binds up until `->`
	for inbinds := true; inbinds && i < len(toks); i++ {
		if tkind := toks[i].Kind(); tkind == lex.TOKEN_OTHER && toks[i].Str == "->" {
			inbinds = false
		} else if tkind == lex.TOKEN_IDENT {
			alt.Binds = append(alt.Binds, toks[i].Str)
		} else {
			return nil, nil, errTok(&toks[i], "expected identifier or `->` instead of `"+toks[i].String()+"`")
		}
	}

	// body of case-alternative after `->`
	bodytoks := toks[i:]
	if len(bodytoks) == 0 {
		return nil, nil, errTok(&toks[len(toks)-1], "missing body of `CASE` alternative")
	}
	expr, exprerr := parseExpr(toks[i:])
	alt.Body = expr
	return alt, tail, exprerr
}
