package climpl

import (
	"fmt"
)

type iSyn interface {
	fmt.Stringer
	setUpd(...func(string) *synBinding)
	taggedSyn()
}

type syn struct{}

func (me *syn) setUpd(...func(string) *synBinding) {}

func (*syn) taggedSyn() {}

type synMod struct {
	syn
	Binds synBindings
}

func (me *synMod) setUpd(resolvers ...func(string) *synBinding) {
	me.Binds.setUpd()
}

type synBindings []*synBinding

func (me synBindings) byName(name string) *synBinding {
	for _, bind := range me {
		if bind.Name == name {
			return bind
		}
	}
	return nil
}

func (me synBindings) setUpd(resolvers ...func(string) *synBinding) {
	r := append(resolvers, me.byName)
	for _, bind := range me {
		bind.setUpd(r...)
	}
}

type synBinding struct {
	syn
	Name    string
	LamForm struct {
		Free []*synExprAtomIdent
		Args []*synExprAtomIdent
		Body iSynExpr
		Upd  bool
	}
}

func (me *synBinding) setUpd(resolvers ...func(string) *synBinding) {
	if me.LamForm.Upd = true; len(me.LamForm.Args) > 0 {
		me.LamForm.Upd = false
	} else if call, iscall := me.LamForm.Body.(*synExprCall); iscall {
		for _, r := range resolvers {
			if def := r(call.Callee.Name); def != nil && len(def.LamForm.Args) != len(call.Args) {
				me.LamForm.Upd = false
				break
			}
		}
	} else if _, isctor := me.LamForm.Body.(*synExprCtor); isctor {
		me.LamForm.Upd = false
	}
	me.LamForm.Body.setUpd(resolvers...)
}

type iSynExpr interface {
	iSyn
	taggedSynExpr()
}

type synExpr struct{ syn }

func (*synExpr) taggedSynExpr() {}

type iSynExprAtom interface {
	iSynExpr
	taggedSynExprAtom()
}

type synExprAtom struct{ synExpr }

func (*synExprAtom) taggedSynExprAtom() {}

type synExprAtomIdent struct {
	synExprAtom
	Name string
}

type iSynExprAtomLit interface {
	iSynExprAtom
	exprAtomLit()
}

type synExprAtomLit struct{ synExprAtom }

func (*synExprAtomLit) exprAtomLit() {}

type synExprAtomLitFloat struct {
	synExprAtomLit
	Lit float64
}

type synExprAtomLitUInt struct {
	synExprAtomLit
	Lit uint64
}

type synExprAtomLitRune struct {
	synExprAtomLit
	Lit rune
}

type synExprAtomLitText struct {
	synExprAtomLit
	Lit string
}

type synExprLet struct {
	synExpr
	Binds synBindings
	Body  iSynExpr
	Rec   bool
}

func (me *synExprLet) setUpd(resolvers ...func(string) *synBinding) {
	me.Binds.setUpd(resolvers...)
	me.Body.setUpd(append(resolvers, me.Binds.byName)...)
}

type synExprCall struct {
	synExpr
	Callee *synExprAtomIdent
	Args   []iSynExprAtom
}

type synExprCtor struct {
	synExpr
	Tag  *synExprAtomIdent
	Args []iSynExprAtom
}

type synExprPrimOp struct {
	synExpr
	PrimOp string
	Left   iSynExprAtom
	Right  iSynExprAtom
}

type synExprCaseOf struct {
	synExpr
	Scrut iSynExpr
	Alts  []*synCaseAlt
}

func (me *synExprCaseOf) setUpd(resolvers ...func(string) *synBinding) {
	me.Scrut.setUpd(resolvers...)
	for _, alt := range me.Alts {
		alt.setUpd(resolvers...)
	}
}

type synCaseAlt struct {
	Ctor struct {
		Tag  *synExprAtomIdent
		Vars []*synExprAtomIdent
	}
	Atom iSynExprAtom
	Body iSynExpr
}

func (me *synCaseAlt) setUpd(resolvers ...func(string) *synBinding) {
	me.Body.setUpd(resolvers...)
}
