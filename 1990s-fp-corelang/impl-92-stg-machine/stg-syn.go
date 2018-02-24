package climpl

type synMod struct {
	Binds []synBinding
}

type synBinding struct {
	Name    string
	LamForm struct {
		Free []synExprAtomIdent
		Args []synExprAtomIdent
		Body iSynExpr
		Upd  bool
	}
}

type iSynExpr interface {
	expr()
}

type synExpr struct{}

func (synExpr) expr() {}

type iSynExprAtom interface {
	iSynExpr
	exprAtom()
}

type synExprAtom struct{ synExpr }

func (synExprAtom) exprAtom() {}

type synExprAtomIdent struct {
	synExprAtom
	Name string
}

type iSynExprAtomLit interface {
	iSynExprAtom
	exprAtomLit()
}

type synExprAtomLit struct{ synExprAtom }

func (synExprAtomLit) exprAtomLit() {}

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
	Binds []synBinding
	Body  iSynExpr
	Rec   bool
}

type synExprCall struct {
	synExpr
	Callee synExprAtomIdent
	Args   []iSynExprAtom
}

type synExprCtor struct {
	synExpr
	Tag  synExprAtomIdent
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
	Alts  []synCaseAlt
}

type synCaseAlt struct {
	Ctor struct {
		Tag  synExprAtomIdent
		Vars []synExprAtomIdent
	}
	Atom iSynExprAtom
	Body iSynExpr
}
