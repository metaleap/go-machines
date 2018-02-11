package clsyn

import (
	"text/scanner"

	lex "github.com/go-leap/dev/lex"
)

type ISyn interface {
	init(lex.Tokens)
	isExpr() bool
	Pos() *lex.TokenMeta
	Toks() lex.Tokens
}

type IExpr interface {
	ISyn
	IsAtomic() bool
}

type syn struct {
	toks lex.Tokens
	// root   *SynMod
	// parent ISyn
}

func (me *syn) init(toks lex.Tokens) { me.toks = toks }

func (*syn) isExpr() bool { return false }

func (me *syn) Pos() *lex.TokenMeta { return &me.toks[0].Meta }

func (me *syn) Toks() lex.Tokens { return me.toks }

type expr struct{ syn }

func (*expr) isExpr() bool { return true }

type exprAtomic struct{ expr }

func (*exprAtomic) IsAtomic() bool { return true }

type exprComp struct{ expr }

func (*exprComp) IsAtomic() bool { return false }

type Error struct {
	msg      string
	Pos      scanner.Position
	RangeLen int
}

func errPos(pos *scanner.Position, msg string, rangeLen int) *Error {
	return &Error{Pos: *pos, msg: msg, RangeLen: rangeLen}
}

func errTok(tok *lex.Token, msg string) *Error {
	return errPos(&tok.Meta.Position, msg, len(tok.String()))
}

func (me *Error) Error() string { return me.msg }
