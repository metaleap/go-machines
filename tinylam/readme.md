# tinylam
--
    import "github.com/metaleap/go-machines/tinylam"


## Usage

```go
const (
	StdModuleName               = "std"
	StdRequiredDefs_true        = StdModuleName + "." + "true"
	StdRequiredDefs_false       = StdModuleName + "." + "false"
	StdRequiredDefs_pairCons    = StdModuleName + "." + "__tlPairCons"
	StdRequiredDefs_listNil     = StdModuleName + "." + "__tlListNil"
	StdRequiredDefs_listIsNil   = StdModuleName + "." + "__tlIsNil"
	StdRequiredDefs_listIsntNil = StdModuleName + "." + "__tlIsntNil"
	StdRequiredDefs_tagCons     = StdModuleName + "." + "__tlTagCons"
)
```

#### func  ValueBool

```go
func ValueBool(it Value) (bool, bool)
```

#### func  ValueBytes

```go
func ValueBytes(it Value) ([]byte, bool)
```

#### func  ValueNum

```go
func ValueNum(it Value) (int, bool)
```

#### func  ValueTag

```go
func ValueTag(it Value) (string, bool)
```

#### type Expr

```go
type Expr interface {
	String() string
	// contains filtered or unexported methods
}
```


#### type ExprCall

```go
type ExprCall struct {
	Callee  Expr
	CallArg Expr
}
```


#### func (*ExprCall) String

```go
func (me *ExprCall) String() string
```

#### type ExprFunc

```go
type ExprFunc struct {
	ArgName string
	Body    Expr
}
```


#### func (*ExprFunc) String

```go
func (me *ExprFunc) String() string
```

#### type ExprLitNum

```go
type ExprLitNum struct {
	NumVal int
}
```


#### func (*ExprLitNum) String

```go
func (me *ExprLitNum) String() string
```

#### type ExprLitTag

```go
type ExprLitTag struct {
	TagVal string
}
```


#### func (*ExprLitTag) String

```go
func (me *ExprLitTag) String() string
```

#### type ExprName

```go
type ExprName struct {
	NameVal string
}
```


#### func (*ExprName) String

```go
func (me *ExprName) String() string
```

#### type Prog

```go
type Prog struct {
	LazyEval     bool
	TopDefs      map[string]Expr
	OnInstrMSG   func(string, Value)
	NumEvalSteps int
}
```


#### func (*Prog) Eval

```go
func (me *Prog) Eval(expr Expr, env Values) Value
```

#### func (*Prog) ParseModules

```go
func (me *Prog) ParseModules(modules map[string][]byte)
```

#### func (*Prog) RunAsMain

```go
func (me *Prog) RunAsMain(mainFuncExpr Expr, osProcArgs []string) (ret Value)
```

#### func (*Prog) Value

```go
func (me *Prog) Value(it Value) (retVal Value)
```

#### type Value

```go
type Value interface {
	String() string
	// contains filtered or unexported methods
}
```


#### type Values

```go
type Values []Value
```


#### func  ValueSlice

```go
func ValueSlice(it Value) (Values, bool)
```

#### func (Values) Copy

```go
func (me Values) Copy() Values
```
