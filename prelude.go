package corelang

const (
	prelDefId    = `id x = x`
	prelDefK0    = `k0 x y = x`
	prelDefK1    = `k1 x y = y`
	prelDefSubst = `subst f g x = f x (g x)`
	prelDefComp  = `comp f g x = f (g x)`
	prelDefComp2 = `comp2 f = comp f f`
)
