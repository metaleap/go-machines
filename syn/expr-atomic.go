package clsyn

func Id(name string) *ExprIdent                { return &ExprIdent{Name: name} }
func Op(name string, lone bool) *ExprIdent     { return &ExprIdent{Name: name, OpLike: true, OpLone: lone} }
func Lf(lit float64) *ExprLitFloat             { return &ExprLitFloat{Lit: lit} }
func Lu(lit uint64, origBase int) *ExprLitUInt { return &ExprLitUInt{Lit: lit, Base: origBase} }
func Lr(lit rune) *ExprLitRune                 { return &ExprLitRune{Lit: lit} }
func Lt(lit string) *ExprLitText               { return &ExprLitText{Lit: lit} }

type ExprIdent struct {
	exprAtomic
	Name   string
	OpLike bool
	OpLone bool
}

type ExprLitFloat struct {
	exprAtomic
	Lit float64
}

type ExprLitUInt struct {
	exprAtomic
	Base int
	Lit  uint64
}

type ExprLitRune struct {
	exprAtomic
	Lit rune
}

type ExprLitText struct {
	exprAtomic
	Lit string
}
