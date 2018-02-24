package climpl

type synPrimOp int

const (
	_ synPrimOp = iota
	SYN_PRIMOP_ADD
	SYN_PRIMOP_SUB
	SYN_PRIMOP_MUL
	SYN_PRIMOP_DIV
)

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

type synExprAtomLitFloat struct {
	synExprAtom
	Lit float64
}

type synExprAtomLitUInt struct {
	synExprAtom
	Lit uint64
}

type synExprAtomLitRune struct {
	synExprAtom
	Lit rune
}

type synExprAtomLitText struct {
	synExprAtom
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
	PrimOp synPrimOp
	Left   iSynExprAtom
	Right  iSynExprAtom
}

type synExprCaseOf struct {
	synExpr
	Scrut iSynExpr
}

type synCaseAlt struct {
	CtorTag   synExprAtomIdent
	CtorVars  []synExprAtomIdent
	PrimOrDef iSynExprAtom
}
