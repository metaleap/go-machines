package main

import (
	"strconv"
	"strings"
)

type IExpr interface {
}

type ExprLit interface{}
type ExprName string
type ExprCall struct {
	Callee IExpr
	Arg    IExpr
}
type ExprFunc struct {
	Arg  string
	Body IExpr
}
type Module map[string]IExpr
type Prog map[string]Module

var allTags = make(map[string]int, 32)

func Load(modules map[string][]byte) Prog {
	prog := make(Prog, len(modules))
	for name, src := range modules {
		prog[name] = parseModule(string(src))
	}
	return prog
}

func parseModule(src string) Module {
	topchunks := strings.Split(src, "\n\n")
	module := make(Module, len(topchunks))
	for _, chunk := range topchunks {
		if chunk = strings.TrimSpace(chunk); chunk != "" {
			topdefname, firstln, lines, localnamesdone := "", "", strings.Split(chunk, "\n"), make(map[string]bool, 8)
			var topdefbody IExpr
			for _, ln := range lines {
				if ln = strings.TrimSpace(ln); ln != "" && ln[0] != '#' {
					if sl, sr := strBreakOn(ln, ':'); firstln == "" {
						firstln, topdefname, topdefbody = ln, sl[0], parseExpr(sr, ln, sl[1:])
					} else if localnamesdone[sl[0]] {
						panic("duplicate local def name '" + sl[0] + "' in:\n" + ln)
					} else {
						localnamesdone[sl[0]], topdefbody = true, ExprCall{Callee: ExprFunc{Arg: sl[0], Body: topdefbody}, Arg: parseExpr(sr, ln, sl[1:])}
					}
				}
			}
			if module[topdefname] == nil {
				module[topdefname] = topdefbody
			} else {
				panic("duplicate global def name '" + topdefname + "' in:\n" + firstln)
			}
		}
	}
	return module
}

func parseExpr(toks []string, ln string, argNames []string) (expr IExpr) {
	if len(toks) > 1 {
		expr = ExprCall{
			Callee: parseExpr(toks[:len(toks)-1], ln, nil),
			Arg:    parseExpr(toks[len(toks)-1:], ln, nil),
		}
	} else if tok := toks[0]; len(tok) == 3 && tok[0] == '\'' && tok[2] == '\'' {
		expr = ExprLit(tok[1])
	} else if tok[0] >= 'A' && tok[0] <= 'Z' {
		tagval, ok := allTags[tok]
		if expr = ExprLit(tagval); !ok {
			expr, allTags[tok] = ExprLit(len(allTags)+1), len(allTags)+1
		}
	} else if isnum, isneg, isdot := (tok[0] >= '0' && tok[0] <= '9'), tok[0] == '-', tok[0] == '.'; isdot || ((isnum || isneg) && strings.IndexByte(tok, '.') > 0) {
		if f64, err := strconv.ParseFloat(tok, 64); err != nil {
			panic(err.Error() + " in:\n" + ln)
		} else {
			expr = ExprLit(f64)
		}
	} else if isnum || isneg {
		if i64, err := strconv.ParseInt(tok, 0, 64); err != nil {
			panic(err.Error() + " in:\n" + ln)
		} else {
			expr = ExprLit(i64)
		}
	} else {
		expr = ExprName(tok)
	}
	if len(argNames) != 0 {
		for i, argsdone := len(argNames)-1, make(map[string]bool, len(argNames)); i >= 0; i-- {
			if argsdone[argNames[i]] {
				panic("duplicate arg name '" + argNames[i] + "' in:\n" + ln)
			} else {
				argsdone[argNames[i]], expr = true, ExprFunc{Arg: argNames[i], Body: expr}
			}
		}
	}
	return
}

func strBreakOn(it string, sep byte) ([]string, []string) {
	if idx := strings.IndexByte(it, sep); idx <= 0 {
		panic("expected '" + string(sep) + "' in:\n" + it)
	} else if sl, sr := strings.TrimSpace(it[:idx]), strings.TrimSpace(it[idx+1:]); sl == "" || sr == "" {
		panic("expected sth. preceding and following the '" + string(sep) + "' in:\n" + it)
	} else {
		return strings.Fields(sl), strings.Fields(sr)
	}
}

func (me Prog) Run(moduleName string, globalDefName string, inStream []byte) (ret IExpr) {
	module := me[moduleName]
	if module == nil {
		panic("module '" + moduleName + "' not known")
	}
	if ret = module[globalDefName]; ret == nil {
		panic("global def '" + globalDefName + "' not known in module '" + moduleName + "'")
	}
	return
}
