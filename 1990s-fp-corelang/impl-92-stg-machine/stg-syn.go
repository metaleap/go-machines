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

func (this *syn) setUpd(...func(string) *synBinding) {}

func (*syn) taggedSyn() {}

type synMod struct {
	syn
	Binds synBindings
}

func (this *synMod) setUpd(resolvers ...func(string) *synBinding) {
	this.Binds.setUpd()
}

type synBindings []*synBinding

func (this synBindings) byName(name string) *synBinding {
	for _, bind := range this {
		if bind.Name == name {
			return bind
		}
	}
	return nil
}

func (this synBindings) setUpd(resolvers ...func(string) *synBinding) {
	r := append(resolvers, this.byName)
	for _, bind := range this {
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

func (this *synBinding) setUpd(resolvers ...func(string) *synBinding) {
	if this.LamForm.Upd = true; len(this.LamForm.Args) > 0 {
		this.LamForm.Upd = false
	} else if call, iscall := this.LamForm.Body.(*synExprCall); iscall {
		for _, r := range resolvers {
			if def := r(call.Callee.Name); def != nil && len(def.LamForm.Args) != len(call.Args) {
				this.LamForm.Upd = false
				break
			}
		}
	} else if _, isctor := this.LamForm.Body.(*synExprCtor); isctor {
		this.LamForm.Upd = false
	}
	this.LamForm.Body.setUpd(resolvers...)
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

func (this *synExprLet) setUpd(resolvers ...func(string) *synBinding) {
	this.Binds.setUpd(resolvers...)
	this.Body.setUpd(append(resolvers, this.Binds.byName)...)
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

func (this *synExprCaseOf) setUpd(resolvers ...func(string) *synBinding) {
	this.Scrut.setUpd(resolvers...)
	for _, alt := range this.Alts {
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

func (this *synCaseAlt) setUpd(resolvers ...func(string) *synBinding) {
	this.Body.setUpd(resolvers...)
}
