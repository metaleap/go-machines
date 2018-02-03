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

func parseExpr(parent ISyn, tokens []lex.IToken) (IExpr, []lex.IToken, *Error) {
	if len(tokens) == 0 {
		return nil, nil, errPos(lex.Pos(nil, parent, ""), "not enough tokens to form an expression", 0)
	}
	expr := parseLit(tokens[0])
	if expr == nil {
		switch t := tokens[0].(type) {
		case *lex.TokenIdent:
			expr = Id(t.Token)
		case *lex.TokenOther:
			expr = Id(t.Token)
		}
	}
	tail := tokens
	if expr != nil {
		tail = tokens[1:]
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
