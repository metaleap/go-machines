# sapl
--
    import "github.com/metaleap/go-machines/sapl"

SAPL interpreter implementation following: **"Efficient Interpretation by
Transforming Data Types and Patterns to Functions"** (Jan Martin Jansen, Pieter
Koopman, Rinus Plasmeijer)

Divergence from the paper: NumArgs is not carried around with the Func Ref but
stored in the top-level-funcs array together with that func's expression.

"Non"-Parser loads from a JSON format: no need to expressly spec it out here,
it's under 40 LoC in `LoadFromJson` and `exprFromJson` funcs.

## Usage

#### type Expr

```go
type Expr interface{ String() string }
```


#### type ExprAppl

```go
type ExprAppl struct {
	Callee Expr
	Arg    Expr
}
```


#### func (ExprAppl) String

```go
func (me ExprAppl) String() string
```

#### type ExprArgRef

```go
type ExprArgRef int
```


#### func (ExprArgRef) String

```go
func (me ExprArgRef) String() string
```

#### type ExprFnRef

```go
type ExprFnRef int
```


#### func (ExprFnRef) String

```go
func (me ExprFnRef) String() string
```

#### type ExprNum

```go
type ExprNum int
```


#### func (ExprNum) String

```go
func (me ExprNum) String() string
```

#### type OpCode

```go
type OpCode int
```


```go
const (
	OpAdd OpCode = -1
	OpSub OpCode = -2
	OpMul OpCode = -3
	OpDiv OpCode = -4
	OpMod OpCode = -5
	OpEq  OpCode = -6
	OpLt  OpCode = -7
	OpGt  OpCode = -8
)
```

#### type Prog

```go
type Prog []TopDef
```


#### func  LoadFromJson

```go
func LoadFromJson(src []byte) Prog
```

#### func (Prog) Eval

```go
func (me Prog) Eval(expr Expr) (ret Expr, timeTaken time.Duration)
```

#### type TopDef

```go
type TopDef = struct {
	NumArgs int
	Expr    Expr
}
```
