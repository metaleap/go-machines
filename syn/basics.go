package clsyn

import (
	"github.com/go-leap/dev/lex"
)

type ISyn interface {
	isExpr() bool
	Pos() *udevlex.TokenMeta
}

type IExpr interface {
	ISyn
	IsAtomic() bool
}

type syn struct {
	pos    udevlex.TokenMeta
	root   *SynMod
	parent ISyn
}

func (*syn) isExpr() bool { return false }

func (me *syn) Pos() *udevlex.TokenMeta { return &me.pos }

type expr struct{ syn }

func (*expr) isExpr() bool { return true }

type exprAtomic struct{ expr }

func (*exprAtomic) IsAtomic() bool { return true }

type exprComp struct{ expr }

func (*exprComp) IsAtomic() bool { return false }
