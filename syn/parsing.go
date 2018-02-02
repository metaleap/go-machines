package clsyn

import (
	"text/scanner"

	"github.com/go-leap/dev/lex"
)

func errPos(pos udevlex.IPos, msg string, rangeLen int) *Error {
	return &Error{Pos: pos.Pos().Position, msg: msg, RangeLen: rangeLen}
}

type Error struct {
	msg      string
	Pos      scanner.Position
	RangeLen int
}

func (me *Error) Error() string { return me.msg }

func parseExpr(parent ISyn, tokens []udevlex.IToken) (IExpr, []udevlex.IToken, *Error) {
	if len(tokens) == 0 {
		return nil, nil, errPos(udevlex.Pos(nil, parent, ""), "not enough tokens to form an expression", 0)
	}
	expr := parseLit(tokens[0])
	if expr == nil {
		switch t := tokens[0].(type) {
		case *udevlex.TokenIdent:
			expr = Id(t.Token)
		case *udevlex.TokenOther:
			expr = Id(t.Token)
		}
	}
	tail := tokens
	if expr != nil {
		tail = tokens[1:]
	}
	return expr, tail, nil
}

func parseLit(token udevlex.IToken) IExpr {
	switch t := token.(type) {
	case *udevlex.TokenFloat:
		return Lf(t.Token)
	case *udevlex.TokenInt:
		return Li(t.Token)
	case *udevlex.TokenUint:
		return Lu(t.Token)
	case *udevlex.TokenRune:
		return Lr(t.Token)
	case *udevlex.TokenStr:
		return Lt(t.Token)
	}
	return nil
}
