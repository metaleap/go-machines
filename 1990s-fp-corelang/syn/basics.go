package clsyn

import (
	lex "github.com/go-leap/dev/lex"
)

type ISyn interface {
	init(lex.Tokens)
	isExpr() bool
	FreeVars(map[string]bool, ...map[string]bool)
	Pos() *lex.Pos
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

func (this *syn) init(toks lex.Tokens) { this.toks = toks }

func (*syn) isExpr() bool { return false }

func (this *syn) Pos() *lex.Pos { return &this.toks[0].Pos }

func (this *syn) Toks() lex.Tokens { return this.toks }

type expr struct{ syn }

func (*expr) isExpr() bool { return true }

type exprAtomic struct{ expr }

func (*exprAtomic) IsAtomic() bool { return true }

type exprComp struct{ expr }

func (*exprComp) IsAtomic() bool { return false }

type Error struct {
	msg      string
	Pos      lex.Pos
	RangeLen int
}

func errPos(pos *lex.Pos, msg string, rangeLen int) *Error {
	return &Error{Pos: *pos, msg: msg, RangeLen: rangeLen}
}

func errTok(tok *lex.Token, msg string) *Error {
	return errPos(&tok.Pos, msg, len(tok.String()))
}

func (this *Error) Error() string { return this.msg }
