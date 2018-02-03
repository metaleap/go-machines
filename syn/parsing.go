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
		return nil, nil, errPos(lex.Pos(nil, parent, ""), "not enough tokens to form an expression", 0)
	}

	var expr IExpr
	sub, subtail, numunclosed := toks.SubTokens("(", ")")
	if numunclosed > 0 {
		return nil, nil, errPos(toks[0], "unclosed parens in current indent level", 1)
	} else if len(sub) > 0 {
		subexpr, _, suberr := parseExpr(parent, sub)
		return subexpr, append(subtail, tail...), suberr
	}

	expr = parseLit(toks[0])
	if expr != nil {
		tail = append(toks[1:], tail...)
	}

	return expr, tail, nil
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
