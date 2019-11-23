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
	me.compileTopDef("same")

	ioutil.WriteFile(outJsonFilePath, []byte(me.outProg.String()), os.ModePerm)
}

func (me *ctxCompileToAtem) compileTopDef(name string) {

}

func (me *ctxCompileToAtem) compileFuncDef(argName string, expr Expr) (idx int) {
	if have, ok := me.done[expr]; ok {
		return have
	}

	me.done[expr] = idx
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

func (me *ctxCompileToAtem) dissectFunc(it *ExprFunc) (innerBody Expr, args []int, freeVars map[int]string) {
	var fns []*ExprFunc
	var numargs int
	for fn := it; fn != nil; fn, _ = fn.Body.(*ExprFunc) {
		if _, done := me.done[fn.Body]; done {
			panic(fn.ArgName)
		}
		fns, numargs, innerBody = append(fns, fn), numargs+1, fn.Body
	}
	args, freeVars = make([]int, numargs), make(map[int]string, 2)
	Walk(innerBody, func(expr Expr) {
		if name, _ := expr.(*ExprName); name != nil && name.IdxOrInstr < 0 {
			if idx := len(args) + name.IdxOrInstr; idx >= 0 && idx < len(args) {
				args[idx] = args[idx] + 1
			} else {
				freeVars[name.IdxOrInstr] = name.NameVal
			}
		}
	})
	return
}
