package clsyn

func Id(name string) *ExprIdent    { return &ExprIdent{Val: name} }
func Lf(lit float64) *ExprLitFloat { return &ExprLitFloat{Val: lit} }
func Li(lit int64) *ExprLitInt     { return &ExprLitInt{Val: lit} }
func Lu(lit uint64) *ExprLitUInt   { return &ExprLitUInt{Val: lit} }
func Lr(lit rune) *ExprLitRune     { return &ExprLitRune{Val: lit} }
func Lt(lit string) *ExprLitText   { return &ExprLitText{Val: lit} }

type ExprIdent struct {
	exprAtomic
	Val string
}

type ExprLitFloat struct {
	exprAtomic
	Val float64
}

type ExprLitInt struct {
	exprAtomic
	Val int64
}

type ExprLitUInt struct {
	exprAtomic
	Val uint64
}

type ExprLitRune struct {
	exprAtomic
	Val rune
}

type ExprLitText struct {
	exprAtomic
	Val string
}
