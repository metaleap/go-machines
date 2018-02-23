package corelang

import (
	. "github.com/metaleap/go-machines/1990s-fp-corelang/syn"
)

var (
	PreludeDefs = map[string]*SynDef{
		// id x = x
		"id": {TopLevel: true, Name: "id", Args: []string{"x"},
			Body: Id("x")},

		// k0 x y = x
		"k0": {TopLevel: true, Name: "k0", Args: []string{"x", "y"},
			Body: Id("x")},

		// k1 x y = y
		"k1": {TopLevel: true, Name: "k1", Args: []string{"x", "y"},
			Body: Id("y")},

		// subst f g x = f x (g x)
		"subst": {TopLevel: true, Name: "subst", Args: []string{"f", "g", "x"},
			Body: Ap(Ap(Id("f"), Id("x")), Ap(Id("g"), Id("x")))},

		// comp f g x = f (g x)
		"comp": {TopLevel: true, Name: "comp", Args: []string{"f", "g", "x"},
			Body: Ap(Id("f"), Ap(Id("g"), Id("x")))},

		// comp2 f = comp f f
		"comp2": {TopLevel: true, Name: "comp2", Args: []string{"f"},
			Body: Ap(Ap(Id("comp"), Id("f")), Id("f"))},
	}
)
