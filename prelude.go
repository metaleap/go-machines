package corelang

var (
	PreludeDefs = map[string]*aDef{
		// id x = x
		"id": {Name: "id", Args: []string{"x"},
			Body: aSym("x")},

		// k0 x y = x
		"k0": {Name: "k0", Args: []string{"x", "y"},
			Body: aSym("x")},

		// k1 x y = y
		"k1": {Name: "k1", Args: []string{"x", "y"},
			Body: aSym("y")},

		// subst f g x = f x (g x)
		"subst": {Name: "subst", Args: []string{"f", "g", "x"},
			Body: aCall(aCall(aSym("f"), aSym("x")), aCall(aSym("g"), aSym("x")))},

		// comp f g x = f (g x)
		"comp": {Name: "comp", Args: []string{"f", "g", "x"},
			Body: aCall(aSym("f"), aCall(aSym("g"), aSym("x")))},

		// comp2 f = comp f f
		"comp2": {Name: "comp2", Args: []string{"f"},
			Body: aCall(aCall(aSym("comp"), aSym("f")), aSym("f"))},
	}
)
