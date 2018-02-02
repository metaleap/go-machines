package clsyn

import (
	"text/scanner"

	"github.com/go-leap/dev/lex"
)

func errPos(pos *scanner.Position, msg string) *Error {
	if pos == nil {
		pos = &scanner.Position{Line: 1, Column: 1, Offset: 0, Filename: ""}
	}
	return &Error{Pos: *pos, msg: msg}
}

type Error struct {
	msg string
	Pos scanner.Position
}

func (me *Error) Error() string { return me.msg }

type parse struct {
	syn  ISyn
	tail []udevlex.Token
}

type parser interface {
	parse([]udevlex.Token) []parse
}

func parseLit(tokens []udevlex.Token) []parse {
	var syn ISyn
	switch t := tokens[0].(type) {
	case *udevlex.TokenFloat:
		syn = Lf(t.Token)
	case *udevlex.TokenInt:
		syn = Li(t.Token)
	case *udevlex.TokenUint:
		syn = Lu(t.Token)
	case *udevlex.TokenRune:
		syn = Lr(t.Token)
	case *udevlex.TokenStr:
		syn = Lt(t.Token)
	}
	if syn != nil {
		return []parse{{syn: syn, tail: tokens[1:]}}
	}
	return nil
}

func parseIdent(tokens []udevlex.Token) []parse {
	if t, _ := tokens[0].(*udevlex.TokenIdent); t != nil {
		return []parse{{syn: Id(t.Token), tail: tokens[1:]}}
	}
	return nil
}
