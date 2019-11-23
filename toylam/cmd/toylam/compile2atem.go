package main

import (
	"io/ioutil"
	"os"
	"strconv"

	. "github.com/metaleap/go-machines/toylam"
)

type ctxCompileToAtem struct {
	prog    Prog
	outProg []atemDef
	done    map[Expr]int
}

type atemDef struct {
	info string
	args []int
	body atemExpr
}
type atemExpr interface{ String() string }
type atemFuncRef int
type atemArgRef int
type atemNumInt int
type atemCall struct {
	callee  atemExpr
	callArg atemExpr
}

func (me atemNumInt) String() string  { return strconv.Itoa(int(me)) }
func (me atemArgRef) String() string  { return "\"" + strconv.Itoa(int(me)) + "\"" }
func (me atemFuncRef) String() string { return "[" + strconv.Itoa(int(me)) + "]" }
func (me atemCall) String() string    { return "[" + me.callee.String() + ", " + me.callArg.String() + "]" }

func (me *ctxCompileToAtem) do(mainTopDefQName string, outJsonFilePath string) {
	me.done = make(map[Expr]int, len(me.prog.TopDefs))
	me.compileFuncDef("same", me.prog.TopDefs["same"])

	ioutil.WriteFile(outJsonFilePath, []byte(me.String()), os.ModePerm)
}

func (me *ctxCompileToAtem) compileFuncDef(info string, body Expr) (idx int) {
	if have, ok := me.done[body]; ok {
		return have
	}
	switch it := body.(type) {
	case *ExprFunc:

	default:
		panic(it)
	}
	me.done[body] = idx
	return
}

func (me *ctxCompileToAtem) compileExpr() {}

func (me *ctxCompileToAtem) String() (outJson string) {
	outJson = "[ "
	for i, def := range me.outProg {
		if i > 0 {
			outJson += ", "
		}
		outJson += def.String() + ",\n"
	}
	return outJson + "]\n"
}
func (me atemDef) String() (outJson string) {
	outJson = "[ ["
	for i, a := range me.args {
		if i > 0 {
			outJson += ","
		}
		outJson += strconv.Itoa(a)
	}
	if outJson += "], "; me.info == "" {
		outJson += me.body.String()
	} else {
		outJson += " { \"_\": " + strconv.Quote(me.info) + ", \"\": " + me.body.String() + " }"
	}
	return outJson + " ]"
}
