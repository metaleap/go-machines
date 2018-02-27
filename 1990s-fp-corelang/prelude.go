package corelang

import (
	. "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
)

var (
	PreludeDefs = []*SynDef{
		// id x = x
		{Name: "id", TopLevel: true, Args: []string{"x"},
			Body: Id("x")},

		// k0 x y = x
		{Name: "k0", TopLevel: true, Args: []string{"x", "y"},
			Body: Id("x")},

		// k1 x y = y
		{Name: "k1", TopLevel: true, Args: []string{"x", "y"},
			Body: Id("y")},

		// subst f g x = f x (g x)
		{Name: "subst", TopLevel: true, Args: []string{"f", "g", "x"},
			Body: Ap(Ap(Id("f"), Id("x")), Ap(Id("g"), Id("x")))},

		// comp f g x = f (g x)
		{Name: "comp", TopLevel: true, Args: []string{"f", "g", "x"},
			Body: Ap(Id("f"), Ap(Id("g"), Id("x")))},

		// comp2 f = comp f f
		{Name: "comp2", TopLevel: true, Args: []string{"f"},
			Body: Ap(Ap(Id("comp"), Id("f")), Id("f"))},
	}
)
