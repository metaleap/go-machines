# corelang
--
    import "github.com/metaleap/go-machines/1991-fpcorelang"


## Usage

```go
var (
	PreludeDefs = map[string]*SynDef{

		"id": {TopLevel: true, Name: "id", Args: []string{"x"},
			Body: Id("x")},

		"k0": {TopLevel: true, Name: "k0", Args: []string{"x", "y"},
			Body: Id("x")},

		"k1": {TopLevel: true, Name: "k1", Args: []string{"x", "y"},
			Body: Id("y")},

		"subst": {TopLevel: true, Name: "subst", Args: []string{"f", "g", "x"},
			Body: Ap(Ap(Id("f"), Id("x")), Ap(Id("g"), Id("x")))},

		"comp": {TopLevel: true, Name: "comp", Args: []string{"f", "g", "x"},
			Body: Ap(Id("f"), Ap(Id("g"), Id("x")))},

		"comp2": {TopLevel: true, Name: "comp2", Args: []string{"f"},
			Body: Ap(Ap(Id("comp"), Id("f")), Id("f"))},
	}
)
```

#### type SyntaxTreePrinter

```go
type SyntaxTreePrinter struct {
}
```


#### func (*SyntaxTreePrinter) Def

```go
func (me *SyntaxTreePrinter) Def(def *SynDef) string
```

#### func (*SyntaxTreePrinter) Expr

```go
func (me *SyntaxTreePrinter) Expr(expr IExpr) string
```

#### func (*SyntaxTreePrinter) Mod

```go
func (me *SyntaxTreePrinter) Mod(mod *SynMod) string
```
