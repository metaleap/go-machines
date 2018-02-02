package corelang

import (
	. "github.com/metaleap/go-corelang/syn"
)

var (
	PreludeDefs = map[string]*Def{
		// id x = x
		"id": {Name: "id", Args: []string{"x"},
			Body: Id("x")},

		// k0 x y = x
		"k0": {Name: "k0", Args: []string{"x", "y"},
			Body: Id("x")},

		// k1 x y = y
		"k1": {Name: "k1", Args: []string{"x", "y"},
			Body: Id("y")},

		// subst f g x = f x (g x)
		"subst": {Name: "subst", Args: []string{"f", "g", "x"},
			Body: Ap(Ap(Id("f"), Id("x")), Ap(Id("g"), Id("x")))},

		// comp f g x = f (g x)
		"comp": {Name: "comp", Args: []string{"f", "g", "x"},
			Body: Ap(Id("f"), Ap(Id("g"), Id("x")))},

		// comp2 f = comp f f
		"comp2": {Name: "comp2", Args: []string{"f"},
			Body: Ap(Ap(Id("comp"), Id("f")), Id("f"))},
	}
)
