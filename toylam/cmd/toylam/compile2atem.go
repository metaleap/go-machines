package main

import (
	"io/ioutil"
	"os"

	"github.com/metaleap/atmo/atem"
	. "github.com/metaleap/go-machines/toylam"
)

var opInstrMappings = map[Instr]atem.OpCode{
	InstrADD: atem.OpAdd,
	InstrDIV: atem.OpDiv,
	InstrEQ:  atem.OpEq,
	InstrGT:  atem.OpGt,
	InstrLT:  atem.OpLt,
	InstrMOD: atem.OpMod,
	InstrMUL: atem.OpMul,
	InstrSUB: atem.OpSub,
}

type ctxCompileToAtem struct {
	prog    *Prog
	outProg atem.Prog
	done    map[Expr]int
}

func (me *ctxCompileToAtem) do(mainTopDefQName string, outJsonFilePath string) {
	me.done = make(map[Expr]int, len(me.prog.TopDefs))
	me.compileFuncDef("same", me.prog.TopDefs["same"])

	ioutil.WriteFile(outJsonFilePath, []byte(me.outProg.String()), os.ModePerm)
}

func (me *ctxCompileToAtem) compileFuncDef(argName string, body Expr) (idx int) {
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

func (me *ctxCompileToAtem) compileExpr(expr Expr) atem.Expr {
	switch it := expr.(type) {
	case *ExprLitNum:
		return atem.ExprNumInt(it.NumVal)
	case *ExprCall:
		return atem.ExprCall{Callee: me.compileExpr(it.Callee), Arg: me.compileExpr(it.CallArg)}
	case *ExprName:
		if instr := Instr(it.IdxOrInstr); it.IdxOrInstr < 0 {
			return atem.ExprArgRef(it.IdxOrInstr)
		} else if op, ok := opInstrMappings[instr]; ok {
			return atem.ExprFuncRef(op)
		}
	case *ExprFunc:
		return atem.ExprFuncRef(me.compileFuncDef(it.ArgName, it.Body))
	}
	panic(expr)
}
