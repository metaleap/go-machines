package clsyn

func Id(name string) *ExprIdent                { return &ExprIdent{Val: name} }
func IdO(name string, lone bool) *ExprIdent    { return &ExprIdent{Val: name, OpLike: true, OpLone: lone} }
func Lf(lit float64) *ExprLitFloat             { return &ExprLitFloat{Val: lit} }
func Lu(lit uint64, origBase int) *ExprLitUInt { return &ExprLitUInt{Val: lit, Base: origBase} }
func Lr(lit rune) *ExprLitRune                 { return &ExprLitRune{Val: lit} }
func Lt(lit string) *ExprLitText               { return &ExprLitText{Val: lit} }

type ExprIdent struct {
	exprAtomic
	Val    string
	OpLike bool
	OpLone bool
}

type ExprLitFloat struct {
	exprAtomic
	Val float64
}

type ExprLitUInt struct {
	exprAtomic
	Base int
	Val  uint64
}

type ExprLitRune struct {
	exprAtomic
	Val rune
}

type ExprLitText struct {
	exprAtomic
	Val string
}
