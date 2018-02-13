package climpl

import (
	"strconv"
)

type code []instr

func (me code) String() (s string) {
	s = "["
	for i, instr := range me {
		if i > 0 {
			s += " · "
		}
		s += instr.String()
	}
	return s + "]"
}

type instruction int

const (
	_ instruction = iota
	INSTR_UNWIND
	INSTR_PUSHGLOBAL
	INSTR_PUSHINT
	INSTR_PUSHARG
	INSTR_MAKEAPPL
	INSTR_SLIDE

	INSTR_UPDATE
	INSTR_POP
	INSTR_ALLOC

	INSTR_EVAL
	INSTR_PRIM_AR_ADD
	INSTR_PRIM_AR_SUB
	INSTR_PRIM_AR_MUL
	INSTR_PRIM_AR_DIV
	INSTR_PRIM_AR_NEG
	INSTR_PRIM_CMP_EQ
	INSTR_PRIM_CMP_NEQ
	INSTR_PRIM_CMP_LT
	INSTR_PRIM_CMP_LEQ
	INSTR_PRIM_CMP_GT
	INSTR_PRIM_CMP_GEQ
	INSTR_PRIM_COND

	INSTR_CTOR_PACK
	INSTR_CASE_JUMP
	INSTR_CASE_SPLIT

	INSTR_MARK7_PUSHINTVAL
	INSTR_MARK7_PUSHNODEINT
	INSTR_MARK7_MAKENODEINT
	INSTR_MARK7_MAKENODEBOOL
)

type instr struct {
	Op        instruction
	Int       int
	Name      string
	CondThen  code
	CondElse  code
	CtorArity int
	CaseJump  []code
}

func (me instr) String() string {
	switch me.Op {
	case INSTR_UNWIND:
		return "Unwd"
	case INSTR_PUSHGLOBAL:
		return "Push`" + me.Name
	case INSTR_PUSHINT:
		return "Push=" + strconv.Itoa(me.Int)
	case INSTR_PUSHARG:
		return "Push@" + strconv.Itoa(me.Int)
	case INSTR_SLIDE:
		return "Slide:" + strconv.Itoa(me.Int)
	case INSTR_MAKEAPPL:
		return "MkAp"
	case INSTR_UPDATE:
		return "Upd@" + strconv.Itoa(me.Int)
	case INSTR_POP:
		return "Pop@" + strconv.Itoa(me.Int)
	case INSTR_ALLOC:
		return "Alloc=" + strconv.Itoa(me.Int)
	case INSTR_EVAL:
		return "Eval"
	case INSTR_PRIM_AR_ADD:
		return "Add"
	case INSTR_PRIM_AR_SUB:
		return "Sub"
	case INSTR_PRIM_AR_MUL:
		return "Mul"
	case INSTR_PRIM_AR_DIV:
		return "Div"
	case INSTR_PRIM_AR_NEG:
		return "Neg"
	case INSTR_PRIM_CMP_EQ:
		return "Eq"
	case INSTR_PRIM_CMP_NEQ:
		return "NEq"
	case INSTR_PRIM_CMP_LT:
		return "Lt"
	case INSTR_PRIM_CMP_LEQ:
		return "LEq"
	case INSTR_PRIM_CMP_GT:
		return "Gt"
	case INSTR_PRIM_CMP_GEQ:
		return "GEq"
	case INSTR_PRIM_COND:
		return "Cond"
	case INSTR_CTOR_PACK:
		return "Ctor" + strconv.Itoa(me.Int) + ":" + strconv.Itoa(me.CtorArity)
	case INSTR_CASE_JUMP:
		return "CJmp#" + strconv.Itoa(len(me.CaseJump))
	case INSTR_CASE_SPLIT:
		return "CSpl=" + strconv.Itoa(me.Int)
	case INSTR_MARK7_PUSHINTVAL:
		return "Push:"
	case INSTR_MARK7_PUSHNODEINT:
		return "PushN"
	case INSTR_MARK7_MAKENODEINT:
		return "MkNI"
	case INSTR_MARK7_MAKENODEBOOL:
		return "MkNB"
	}
	return strconv.Itoa(int(me.Op))
}
