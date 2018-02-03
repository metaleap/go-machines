package clsyn

import (
	"text/scanner"

	lex "github.com/go-leap/dev/lex"
)

type ISyn interface {
	isExpr() bool
	Pos() *lex.TokenMeta
}

type IExpr interface {
	ISyn
	IsAtomic() bool
}

type syn struct {
	pos    lex.TokenMeta
	root   *SynMod
	parent ISyn
}

func (*syn) isExpr() bool { return false }

func (me *syn) Pos() *lex.TokenMeta { return &me.pos }

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

func errPos(pos lex.IPos, msg string, rangeLen int) *Error {
	return &Error{Pos: pos.Pos().Position, msg: msg, RangeLen: rangeLen}
}

func (me *Error) Error() string { return me.msg }
