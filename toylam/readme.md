# toylam
--
    import "github.com/metaleap/go-machines/toylam"


## Usage

```go
const (
	StdModuleName            = "std"
	StdRequiredDefs_true     = StdModuleName + "." + "True"
	StdRequiredDefs_false    = StdModuleName + "." + "False"
	StdRequiredDefs_tupCons  = StdModuleName + "." + "Pair"
	StdRequiredDefs_list     = StdModuleName + "." + "List"
	StdRequiredDefs_listCons = StdModuleName + "." + "ListLink"
	StdRequiredDefs_listNil  = StdModuleName + "." + "ListEnd"
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

#### func  ValueOther

```go
func ValueOther(it Value) (string, bool)
```

#### func  Walk

```go
func Walk(expr Expr, visitor func(Expr))
```

#### type Expr

```go
type Expr interface {
	LocInfo() *Loc
	NamesDeclared() []string
	ReplaceName(string, string) int
	RewriteName(string, Expr) Expr
	String() string
}
```


#### type ExprCall

```go
type ExprCall struct {
	*Loc
	Callee  Expr
	CallArg Expr
}
```


#### func (*ExprCall) NamesDeclared

```go
func (me *ExprCall) NamesDeclared() []string
```

#### func (*ExprCall) ReplaceName

```go
func (me *ExprCall) ReplaceName(nameOld string, nameNew string) int
```

#### func (*ExprCall) RewriteName

```go
func (me *ExprCall) RewriteName(name string, with Expr) Expr
```

#### func (*ExprCall) String

```go
func (me *ExprCall) String() string
```

#### type ExprFunc

```go
type ExprFunc struct {
	*Loc
	ArgName string
	Body    Expr
}
```


#### func (*ExprFunc) NamesDeclared

```go
func (me *ExprFunc) NamesDeclared() []string
```

#### func (*ExprFunc) ReplaceName

```go
func (me *ExprFunc) ReplaceName(old string, new string) int
```

#### func (*ExprFunc) RewriteName

```go
func (me *ExprFunc) RewriteName(name string, with Expr) Expr
```

#### func (*ExprFunc) String

```go
func (me *ExprFunc) String() string
```

#### type ExprLitNum

```go
type ExprLitNum struct {
	*Loc
	NumVal int
}
```


#### func (*ExprLitNum) NamesDeclared

```go
func (me *ExprLitNum) NamesDeclared() []string
```

#### func (*ExprLitNum) ReplaceName

```go
func (me *ExprLitNum) ReplaceName(string, string) int
```

#### func (*ExprLitNum) RewriteName

```go
func (me *ExprLitNum) RewriteName(string, Expr) Expr
```

#### func (*ExprLitNum) String

```go
func (me *ExprLitNum) String() string
```

#### type ExprName

```go
type ExprName struct {
	*Loc
	NameVal    string
	IdxOrInstr int // if <0 then De Bruijn index, if >0 then instrCode
}
```


#### func (*ExprName) NamesDeclared

```go
func (me *ExprName) NamesDeclared() []string
```

#### func (*ExprName) ReplaceName

```go
func (me *ExprName) ReplaceName(nameOld string, nameNew string) (didReplace int)
```

#### func (*ExprName) RewriteName

```go
func (me *ExprName) RewriteName(name string, with Expr) Expr
```

#### func (*ExprName) String

```go
func (me *ExprName) String() string
```

#### type Instr

```go
type Instr int
```


```go
const (
	InstrADD Instr
	InstrMUL
	InstrSUB
	InstrDIV
	InstrMOD
	InstrMSG
	InstrERR
	InstrEQ
	InstrGT
	InstrLT
)
```

#### type Loc

```go
type Loc struct {
	ModuleName string
	TopDefName string
	LineNr     int
}
```


#### func (*Loc) LocInfo

```go
func (me *Loc) LocInfo() *Loc
```

#### func (*Loc) LocStr

```go
func (me *Loc) LocStr() string
```

#### type ParseOpts

```go
type ParseOpts struct {
	KeepRec       bool
	KeepNameRefs  bool
	KeepOpRefs    bool
	KeepSepLocals bool
}
```

for compilers or other syntax users. defaults to all-`false`s for our
interpreter in here

#### type Prog

```go
type Prog struct {
	LazyEval        bool
	TopDefs         map[string]Expr
	TopDefSepLocals map[string][]localDef
	OnInstrMSG      func(string, Value)
	NumEvalSteps    int
}
```


#### func (*Prog) Eval

```go
func (me *Prog) Eval(expr Expr, env Values) Value
```

#### func (*Prog) ParseModules

```go
func (me *Prog) ParseModules(modules map[string][]byte, opts ParseOpts)
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
