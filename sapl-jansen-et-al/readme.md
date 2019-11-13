# sapl
--
    import "github.com/metaleap/go-machines/sapl-jansen-et-al"


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
type ExprFnRef struct {
	NumArgs int
	Idx     int
}
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
type Prog []Expr
```


#### func  LoadFromJson

```go
func LoadFromJson(src []byte) Prog
```

#### func (Prog) Eval

```go
func (me Prog) Eval(expr Expr) Expr
```
