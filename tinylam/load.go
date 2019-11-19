package tinylam

import (
	"bytes"
	"strconv"
	"strings"
)

const (
	StdModuleName            = "std"
	StdRequiredDefs_true     = StdModuleName + "." + "True"
	StdRequiredDefs_false    = StdModuleName + "." + "False"
	StdRequiredDefs_tupCons  = StdModuleName + "." + "Pair"
	StdRequiredDefs_list     = StdModuleName + "." + "List"
	StdRequiredDefs_listCons = StdModuleName + "." + "ListLink"
	StdRequiredDefs_listNil  = StdModuleName + "." + "ListEnd"
)

type ctxParse struct {
	prog      *Prog
	srcs      map[string][]byte
	counter   int
	curModule struct{ name string }
	curTopDef struct {
		bracketsParens  map[string]string
		bracketsCurlies map[string]string
		bracketsSquares map[string]string
	}
}

type nodeLocInfo struct {
	srcLocModuleName string
	srcLocTopDefName string
	srcLocLineNr     int
}

func (me *nodeLocInfo) locInfo() *nodeLocInfo { return me }
func (me *nodeLocInfo) locStr() string {
	if me == nil {
		return ""
	}
	return "in '" + me.srcLocModuleName + "." + me.srcLocTopDefName + "', line " + strconv.Itoa(me.srcLocLineNr) + ": "
}

func (me *Prog) ParseModules(modules map[string][]byte) {
	ctx := ctxParse{prog: me, srcs: modules}
	ctx.curTopDef.bracketsParens, ctx.curTopDef.bracketsCurlies, ctx.curTopDef.bracketsSquares = make(map[string]string, 16), make(map[string]string, 2), make(map[string]string, 4) // reset every top-def, potentially needed earlier for type-spec top-defs
	if me.NumEvalSteps, me.TopDefs = 0, map[string]Expr{}; me.pseudoSumTypes == nil {
		me.pseudoSumTypes = map[string][]pseudoSumTypeCtor{}
	}

	for modulename, modulesrc := range modules {
		ctx.curModule.name = modulename
		modules[modulename] = ctx.gatherPseudoSumTypesAndBasedOnTheirDefsAppendToSrcs(ctx.rewriteStrLitsToIntLists(modulesrc))
	}
	for modulename, modulesrc := range modules {
		ctx.curModule.name = modulename
		module := ctx.parseModule(string(modulesrc))
		for topdefname, topdefbody := range module {
			me.TopDefs[modulename+"."+topdefname] = ctx.populateNames(topdefbody, make(map[string]int, 16), module, topdefname)
		}
	}
	me.exprBoolTrue, me.exprBoolFalse = me.TopDefs[StdRequiredDefs_true].(*ExprFunc), me.TopDefs[StdRequiredDefs_false].(*ExprFunc)
	me.exprListNil, me.exprListConsCtorBody = me.TopDefs[StdRequiredDefs_listNil].(*ExprFunc), me.TopDefs[StdRequiredDefs_listCons].(*ExprFunc).Body.(*ExprFunc).Body.(*ExprFunc).Body
	for instrname, instrcode := range instrs {
		me.TopDefs[StdModuleName+".//op"+instrname] = &ExprFunc{nil, "//" + instrname, &ExprCall{nil, &ExprName{nil, instrname, int(instrcode)}, &ExprName{nil, "//" + instrname, -1}}, -1}
	}
	for topdefqname, topdefbody := range me.TopDefs {
		me.TopDefs[topdefqname] = me.preResolveExprs(topdefbody, topdefqname, topdefbody)
	}
}

func (me *ctxParse) gatherPseudoSumTypesAndBasedOnTheirDefsAppendToSrcs(moduleSrc []byte) []byte {
	lines := strings.Split(string(moduleSrc), "\n")
	for l, i := len(lines), 0; i < l; i++ {
		if ln := lines[i]; len(ln) > 0 && ln[0] >= 'A' && ln[0] <= 'Z' {
			if idx := strings.Index(ln, ":="); idx > 0 && !strings.Contains(ln, " -> __") {
				if tparts, cparts := strings.Fields(ln[:idx]), strings.Split(me.extractBrackets(nil, strings.TrimSpace(ln[idx+2:]), ln, 1), " | "); len(tparts) == 1 && len(cparts) > 0 && len(cparts[0]) > 0 {
					lines[i] = "//" + lines[i]
					for i, cpart := range cparts {
						cpart = " " + cpart + " "
						for u, underscore := 0, strings.Index(cpart, " _ "); underscore > 0; u, underscore = u+1, strings.Index(cpart, " _ ") {
							cpart = cpart[:underscore] + " __" + strings.ToLower(tparts[0]) + strconv.Itoa(u) + "__" + cpart[underscore+2:]
						}
						cpart = strings.TrimSpace(cpart)
						cparts[i] = cpart
						str := cpart + " :="
						for _, ctorstr := range cparts {
							ctorstr += " "
							str += " __" + tparts[0] + "_Of_" + ctorstr[:strings.IndexByte(ctorstr, ' ')]
						}
						str += " -> __" + tparts[0] + "_Of_" + cpart
						lines = append(lines, str)
					}
					tqname, strcases := me.curModule.name+"."+tparts[0], ""
					for cidx, cpart := range cparts {
						parts := strings.Fields(cpart)
						me.prog.pseudoSumTypes[tqname] = append(me.prog.pseudoSumTypes[tqname], pseudoSumTypeCtor{parts[0], len(parts) - 1})
						if strcases += " caseOf" + parts[0]; len(parts) > 1 {
							for _, ctorarg := range parts[1:] {
								if strdtor := ctorarg + "Of" + tparts[0] + parts[0]; len(ctorarg) > 1 && ctorarg[0] != '_' {
									strdtor += " a" + tparts[0] + "Of" + parts[0] + " := a" + tparts[0] + "Of" + parts[0]
									for cidx2 := range cparts {
										if strdtor += " ("; cidx2 != cidx {
											strdtor += ")"
										} else {
											for _, ca := range parts[1:] {
												strdtor += ca + " "
											}
											strdtor += "-> " + ctorarg + ")"
										}
									}
									lines = append(lines, strdtor)
								}
							}
						}
					}
					str := tparts[0] + strcases + " scrutinee" + tparts[0] + " := scrutinee" + tparts[0] + strcases
					lines = append(lines, str)
				}
			}
		}
	}
	return []byte(strings.Join(lines, "\n") + "\n")
}

func (me *ctxParse) parseModule(src string) map[string]Expr {
	lines, module := strings.Split(src, "\n"), make(map[string]Expr, 32)
	for idx, last, i := 0, len(lines), len(lines)-1; i >= 0; i-- {
		if idx = strings.Index(lines[i], "//"); idx >= 0 {
			lines[i] = lines[i][:idx]
		}
		if nonempty := (len(lines[i]) > 0); i == 0 || (nonempty && lines[i][0] != ' ' && lines[i][0] != '\t') {
			if topdefname, topdefbody, firstln := me.parseTopDef(lines, i, last); topdefname != "" && topdefbody != nil {
				if module[topdefname] != nil || topdefname == "_" || strings.IndexByte(topdefname, '.') >= 0 {
					panic("in '" + me.curModule.name + "', line " + strconv.Itoa(i+1) + ": illegal or duplicate global def name '" + topdefname + "' in:\n" + firstln)
				}
				module[topdefname] = topdefbody
			}
			last = i
		}
	}
	return module
}

func (me *ctxParse) parseTopDef(lines []string, idxStart int, idxEnd int) (topDefName string, topDefBody Expr, firstLn string) {
	topDefName, me.curTopDef.bracketsParens, me.curTopDef.bracketsCurlies, me.curTopDef.bracketsSquares = "<unknown>", make(map[string]string, 16), make(map[string]string, 2), make(map[string]string, 4)
	var topdefargs []string
	for i, ln := range lines[idxStart:idxEnd] {
		if ln = strings.TrimSpace(ln); ln != "" {
			lnorig, loc := ln, &nodeLocInfo{me.curModule.name, topDefName, 1 + i + idxStart}
			for idx := strings.IndexByte(ln, '\''); idx >= 0 && idx < (len(ln)-1); idx = strings.IndexByte(ln, '\'') {
				if (idx+3) <= len(ln) && ln[idx+2] == '\'' {
					ln = ln[:idx] + " " + strconv.FormatUint(uint64(ln[idx+1]), 10) + " " + ln[idx+3:]
				} else {
					panic(loc.locStr() + "bad quoted-byte-literal " + ln[idx:idx+2] + " in line:\n" + lnorig)
				}
			}
			ln = me.extractBrackets(loc, ln, lnorig, 0)
			if idx := strings.Index(ln, ":="); idx < 0 {
				panic(loc.locStr() + "expected ':=' in:\n" + lnorig)
			} else if lhs, rhs := strings.TrimSpace(ln[:idx]), strings.TrimSpace(ln[idx+2:]); lhs == "" || rhs == "" {
				panic(loc.locStr() + "expected '<name/s> := <expr>' in:\n" + lnorig)
			} else if defsig := strings.Fields(lhs); firstLn == "" {
				firstLn, topDefName, loc.srcLocTopDefName, topdefargs = lnorig, defsig[0], defsig[0], defsig[1:]
				topDefBody = me.parseExpr(rhs, lnorig, loc)
			} else if localname := defsig[0]; localname == "_" || strings.IndexByte(localname, '.') >= 0 {
				panic(loc.locStr() + "illegal  local def name '" + localname + "' in:\n" + lnorig)
			} else {
				localbody := me.hoistArgs(me.parseExpr(rhs, lnorig, loc), defsig[1:])
				if 0 < localbody.replaceName(localname, "//recur3//"+localname) {
					localbody = me.rewriteForRecursion(localname, localbody, "recur")
				}
				topDefBody = &ExprCall{loc, &ExprFunc{loc, localname, topDefBody, -1}, localbody}
			}
		}
	}
	topDefBody = me.hoistArgs(topDefBody, topdefargs)
	if topDefBody != nil && 0 < topDefBody.replaceName(topDefName, topDefName) /* aka "refers to"*/ {
		topDefBody = me.rewriteForRecursion(topDefName, topDefBody, "Recur")
	}
	return
}

func (me *ctxParse) rewriteForRecursion(defName string, defBody Expr, dynNamePref string) Expr {
	return &ExprCall{defBody.locInfo(), &ExprFunc{defBody.locInfo(), "//" + dynNamePref + "1//" + defName, &ExprCall{defBody.locInfo(), &ExprName{defBody.locInfo(), "//" + dynNamePref + "1//" + defName, 0}, &ExprName{defBody.locInfo(), "//" + dynNamePref + "1//" + defName, 0}}, -1}, &ExprFunc{defBody.locInfo(), "//" + dynNamePref + "2//" + defName, defBody, -1}}
}

func (me *ctxParse) parseExpr(src string, locHintLn string, locInfo *nodeLocInfo) (expr Expr) {
	if idx := strings.Index(src, "? "); idx > 0 {
		if pos := strings.LastIndexByte(src[:idx], ' '); pos > 0 {
			tname := src[pos+1 : idx]
			if tname == "" {
				tname = StdRequiredDefs_list
			}
			ctors := me.prog.pseudoSumTypes[tname]
			if len(ctors) == 0 {
				ctors = me.prog.pseudoSumTypes[me.curModule.name+"."+tname]
			}
			if len(ctors) == 0 {
				ctors = me.prog.pseudoSumTypes[StdModuleName+"."+tname]
			}
			if len(ctors) > 0 {
				scases := strings.Split(src[idx+1:], " | ")
				mcases := make(map[string]string, len(scases))
				for _, scase := range scases {
					scase = strings.TrimSpace(scase)
					if idxcol := strings.Index(scase, "=>"); idxcol < 0 {
						panic(locInfo.locStr() + "for scrutinizing `" + tname + "`, expected `=>` in case:\n" + scase)
					} else {
						casename := strings.TrimSpace(scase[:idxcol])
						if tname == StdRequiredDefs_list {
							if casename == "" {
								panic(locInfo.locStr() + "fallback/default empty-case not supported for scrutinizing `" + tname + "`")
							} else if emptysquarebrackets, ok := me.curTopDef.bracketsSquares[casename]; ok && len(emptysquarebrackets) == 0 {
								casename = StdRequiredDefs_listNil
							} else if casename == ".." || (ok && emptysquarebrackets == ",") {
								casename = StdRequiredDefs_listCons
							}
						}
						mcases[casename[strings.LastIndexByte(casename, '.')+1:]] = strings.TrimSpace(scase[idxcol+2:])
					}
				}
				for _, ctor := range ctors {
					if mcases[ctor.name] == "" {
						if tname == StdRequiredDefs_list && (ctor.name == StdRequiredDefs_listNil || ctor.name == StdRequiredDefs_listNil[strings.IndexByte(StdRequiredDefs_listNil, '.')+1:]) {
							mcases[ctor.name] = StdRequiredDefs_listNil
						} else if mcases[""] == "" {
							panic(locInfo.locStr() + "for scrutinizing `" + tname + "`, a case for `" + ctor.name + "` is required")
						}
					}
				}
				if len(mcases) != len(ctors) && mcases[""] == "" {
					panic(locInfo.locStr() + "for scrutinizing `" + tname + "`, expected " + strconv.Itoa(len(ctors)) + " cases but found (effectively) " + strconv.Itoa(len(mcases)) + ", in:\n" + src)
				} else {
					src = src[:pos] + " "
					for _, ctor := range ctors {
						tmpname, casecode := "__case__of__"+ctor.name, mcases[ctor.name]
						if me.curTopDef.bracketsParens[tmpname] = casecode; casecode == "" {
							me.curTopDef.bracketsParens[tmpname] = strings.Repeat("_ -> ", ctor.arity) + mcases[""]
						}
						src += " " + tmpname
					}
				}
			}
		}
	}
	return me.parseExprToks(strings.Fields(src), locHintLn, locInfo)
}

func (me *ctxParse) parseExprToks(toks []string, locHintLn string, locInfo *nodeLocInfo) (expr Expr) {
	if len(toks) == 0 {
		panic(locInfo.locStr() + " expression expected before / after comma in:\n" + locHintLn)
	} else if tok, islambda, lamsplit := toks[0], 0, 0; len(toks) > 1 {
		me.counter++
		for i := range toks {
			if lamsplit == 0 && toks[i] == "->" {
				lamsplit = i
				break
			} else if l := len(toks[i]); lamsplit == 0 && toks[i][0] == '_' && toks[i] == strings.Repeat("_", l) {
				if toks[i] = "//lam//" + strconv.Itoa(l) + "//" + strconv.Itoa(me.counter); l > islambda {
					islambda = l
				}
			}
		}
		if args := toks[:lamsplit]; lamsplit > 0 {
			tupdtorpos := -1
			for i, tok := range toks[:lamsplit] {
				if "" != strings.TrimSpace(me.curTopDef.bracketsCurlies[tok]) {
					if tupdtorpos >= 0 {
						panic(locInfo.locStr() + "no multiple tuple-destructors per single lambda please:\n" + locHintLn)
					}
					tupdtorpos = i
				} else if tok[0] == '/' {
					me.counter, toks[i] = me.counter+1, "//"+strconv.Itoa(i)+"//"+strconv.Itoa(me.counter)
				}
			}
			if tupdtorpos >= 0 {
				tupdtor := me.curTopDef.bracketsCurlies[toks[tupdtorpos]]
				me.curTopDef.bracketsParens["__"+toks[tupdtorpos]] = tupdtor + " -> " + strings.Join(toks[lamsplit+1:], " ")
				toks = append(toks[:lamsplit+1], toks[tupdtorpos]+"__", "__"+toks[tupdtorpos])
				args[tupdtorpos], toks[tupdtorpos] = toks[tupdtorpos]+"__", toks[tupdtorpos]+"__"
			}
			expr = me.hoistArgs(me.parseExprToks(toks[lamsplit+1:], locHintLn, locInfo), args)
		} else if args = make([]string, islambda); islambda > 0 {
			for i := range args {
				args[i] = "//lam//" + strconv.Itoa(i+1) + "//" + strconv.Itoa(me.counter)
			}
			expr = me.hoistArgs(me.parseExprToks(toks, locHintLn, locInfo), args)
		} else {
			expr = &ExprCall{locInfo, me.parseExprToks(toks[:len(toks)-1], locHintLn, locInfo), me.parseExprToks(toks[len(toks)-1:], locHintLn, locInfo)}
		}
	} else if isnum, isneg := (tok[0] >= '0' && tok[0] <= '9'), tok[0] == '-' && len(tok) > 1; isnum || (isneg && tok[1] >= '0' && tok[1] <= '9') {
		if numint, err := strconv.ParseInt(tok, 0, 0); err != nil {
			panic(locInfo.locStr() + err.Error() + " in:\n" + locHintLn)
		} else {
			expr = &ExprLitNum{locInfo, int(numint)}
		}
	} else if subexpr, ok := me.curTopDef.bracketsParens[tok]; ok {
		if subexpr = strings.TrimSpace(subexpr); subexpr == "" {
			expr = &ExprCall{locInfo, &ExprName{locInfo, "ERR", int(instrERR)}, me.prog.newStr(true, locInfo, "forced crash via `()`!")}
		} else {
			expr = me.parseExpr(subexpr, locHintLn, locInfo)
		}
	} else if subexpr, ok = me.curTopDef.bracketsSquares[tok]; ok {
		expr = &ExprName{locInfo, StdRequiredDefs_listNil, 0}
		if items := strings.Split(strings.Trim(strings.TrimSpace(subexpr), ","), ","); len(items) > 0 && items[0] != "" {
			for i := len(items) - 1; i >= 0; i-- {
				expr = &ExprCall{locInfo, &ExprCall{locInfo, &ExprName{locInfo, StdRequiredDefs_listCons, 0}, me.parseExpr(items[i], locHintLn, locInfo)}, expr}
			}
		}
	} else if subexpr, ok = me.curTopDef.bracketsCurlies[tok]; ok {
		if items := strings.Split(strings.Trim(strings.TrimSpace(subexpr), ","), ","); len(items) == 0 || (len(items) == 1 && items[0] == "") {
			expr = &ExprName{locInfo, StdRequiredDefs_tupCons, 0}
		} else if len(items) == 1 {
			expr = &ExprCall{locInfo, &ExprName{locInfo, StdRequiredDefs_tupCons, 0}, me.parseExpr(items[0], locHintLn, locInfo)}
		} else {
			expr = me.parseExpr(items[len(items)-1], locHintLn, locInfo)
			for i := len(items) - 2; i >= 0; i-- {
				expr = &ExprCall{locInfo, &ExprCall{locInfo, &ExprName{locInfo, StdRequiredDefs_tupCons, 0}, me.parseExpr(items[i], locHintLn, locInfo)}, expr}
			}
		}
	} else {
		expr = &ExprName{locInfo, tok, int(instrs[tok])}
	}
	return
}

func (me *ctxParse) rewriteStrLitsToIntLists(src []byte) []byte {
	if bytes.IndexByte(src, 0) >= 0 {
		panic("NUL char in module source: " + me.curModule.name)
	}
	src = bytes.ReplaceAll(src, []byte{'\\', '"'}, []byte{0})
	for idx := bytes.IndexByte(src, '"'); idx > 0; idx = bytes.IndexByte(src, '"') {
		if pos := bytes.IndexByte(src[idx+1:], '"'); pos < 0 {
			panic("in '" + me.curModule.name + "': non-terminated string literal: " + string(src[idx:]))
		} else {
			src[idx], src[idx+1+pos] = '[', ']'
			pref, suff, inner := src[:idx+1], src[idx+1+pos:], make([]byte, 0, 4*(1+pos))
			for i := idx + 1; i < idx+1+pos; i++ {
				b := src[i]
				if b == 0 {
					b = '"'
				}
				inner = append(append(inner, strconv.FormatUint(uint64(b), 10)...), ',')
			}
			src = append(pref, append(inner, suff...)...)
		}
	}
	return bytes.ReplaceAll(src, []byte{0}, []byte{'\\', '"'})
}

func (*ctxParse) hoistArgs(expr Expr, argNames []string) Expr {
	if expr != nil && len(argNames) != 0 {
		for i, argname := len(argNames)-1, ""; i >= 0; i-- {
			if argname = argNames[i]; argname == "_" {
				argname = strconv.Itoa(i)
			}
			expr = &ExprFunc{expr.locInfo(), argname, expr, -1}
		}
	}
	return expr
}

func (me *ctxParse) populateNames(expr Expr, binders map[string]int, curModule map[string]Expr, locHintTopDefName string) Expr {
	const stdpref = StdModuleName + "."
	fixinstrval := func(expr Expr) {
		if name, _ := expr.(*ExprName); name != nil && name.idxOrInstr > 0 {
			name.NameVal, name.idxOrInstr = stdpref+"//op"+name.NameVal, 0
		}
	}
	switch it := expr.(type) {
	case *ExprCall:
		it.Callee = me.populateNames(it.Callee, binders, curModule, locHintTopDefName)
		it.CallArg = me.populateNames(it.CallArg, binders, curModule, locHintTopDefName)
		fixinstrval(it.CallArg)
		if fn, _ := it.Callee.(*ExprFunc); fn != nil && 1 == fn.replaceName(fn.ArgName, fn.ArgName) && 0 == it.CallArg.replaceName(fn.ArgName, fn.ArgName) {
			nope, fnames, cnames := false, fn.namesDeclared(), it.CallArg.namesDeclared()
			for _, fname := range fnames {
				if nope = (0 != it.CallArg.replaceName(fname, fname)); nope {
					break
				}
				for _, cname := range cnames {
					if nope = (cname == fname); nope {
						break
					}
				}
			}
			if !nope {
				expr = fn.Body.rewriteName(fn.ArgName, it.CallArg.rewriteName(fn.ArgName, nil))
				return me.populateNames(expr, binders, curModule, locHintTopDefName)
			}
		}
	case *ExprName:
		if it.NameVal == locHintTopDefName {
			return me.populateNames(&ExprCall{it.locInfo(), &ExprName{it.locInfo(), "//Recur2//" + it.NameVal, 0}, &ExprName{it.locInfo(), "//Recur2//" + it.NameVal, 0}}, binders, curModule, locHintTopDefName)
		} else if strings.HasPrefix(it.NameVal, "//recur3//") {
			it.NameVal = it.NameVal[len("//recur3//"):]
			return me.populateNames(&ExprCall{it.locInfo(), &ExprName{it.locInfo(), "//recur2//" + it.NameVal, 0}, &ExprName{it.locInfo(), "//recur2//" + it.NameVal, 0}}, binders, curModule, locHintTopDefName)
		}
		if posdot := strings.LastIndexByte(it.NameVal, '.'); posdot > 0 && nil == me.srcs[it.NameVal[:posdot]] && (nil == me.srcs[stdpref+it.NameVal[:posdot]] || 0 != binders[it.NameVal[:posdot]]) {
			dotpath := strings.Split(it.NameVal, ".") // desugar a.b.c into (c (b a))
			var ret Expr = &ExprName{it.nodeLocInfo, dotpath[0], 0}
			for i := 1; i < len(dotpath); i++ {
				ret = &ExprCall{it.nodeLocInfo, &ExprName{it.nodeLocInfo, dotpath[i], int(instrs[dotpath[i]])}, ret}
			}
			return me.populateNames(ret, binders, curModule, locHintTopDefName)
		} else if it.idxOrInstr == 0 && posdot < 0 { // neither a prim-instr-op-code, nor an already-qualified cross-module reference
			if it.idxOrInstr = binders[it.NameVal]; it.idxOrInstr > 0 {
				it.idxOrInstr = -it.idxOrInstr // mark as referring to a local / arg (De Bruijn index but negative)
			} else if _, topdefexists := curModule[it.NameVal]; topdefexists {
				it.NameVal = me.curModule.name + "." + it.NameVal // mark as referring to a global in the current module
			} else {
				it.NameVal = stdpref + it.NameVal // mark as referring to a global in std
			}
		}
	case *ExprFunc:
		if _, topdefexists := curModule[it.ArgName]; topdefexists || binders[it.ArgName] != 0 {
			panic("in '" + me.curModule.name + "." + locHintTopDefName + "', line " + strconv.Itoa(it.srcLocLineNr) + ": local name '" + it.ArgName + "' already taken: " + it.String())
		}
		for k, v := range binders {
			binders[k] = v + 1
		}
		binders[it.ArgName] = 1
		it.Body = me.populateNames(it.Body, binders, curModule, locHintTopDefName)
		fixinstrval(it.Body)
		it.numArgUses = it.Body.replaceName(it.ArgName, it.ArgName)
		delete(binders, it.ArgName) // must delete, not just zero (because of our map-ranging incrs/decrs)
		for k, v := range binders {
			binders[k] = v - 1
		}
		if call, _ := it.Body.(*ExprCall); call != nil && locHintTopDefName != "main" {
			if arg, _ := call.CallArg.(*ExprName); arg != nil && arg.NameVal == it.ArgName && 0 == call.Callee.replaceName(it.ArgName, it.ArgName) {
				return me.populateNames(call.Callee.rewriteName(it.ArgName, nil), binders, curModule, locHintTopDefName)
			}
		}
	}
	return expr
}

func (me *ctxParse) extractBrackets(loc *nodeLocInfo, ln string, lnOrig string, needLegalName int) string {
	for str, b, idx, m, ip, ic, is := "", byte(0), 0, map[string]string(nil), strings.IndexByte(ln, ')'), strings.IndexByte(ln, '}'), strings.IndexByte(ln, ']'); ip > 0 || ic > 0 || is > 0; me.counter, ip, ic, is = me.counter+1, strings.IndexByte(ln, ')'), strings.IndexByte(ln, '}'), strings.IndexByte(ln, ']') {
		if p, c, s := ip > 0 && (ic <= 0 || ip < ic) && (is <= 0 || ip < is), ic > 0 && (ip <= 0 || ic < ip) && (is <= 0 || ic < is), is > 0 && (ic <= 0 || is < ic) && (ip <= 0 || is < ip); p {
			str, b, idx, m = "//bp//", '(', ip, me.curTopDef.bracketsParens
		} else if c {
			str, b, idx, m = "//bc//", '{', ic, me.curTopDef.bracketsCurlies
		} else if s {
			str, b, idx, m = "//bs//", '[', is, me.curTopDef.bracketsSquares
		}
		if name, pos := str+strconv.Itoa(me.counter), strings.LastIndexByte(ln[:idx], b); pos < 0 {
			panic(loc.locStr() + "missing opening '" + string(b) + "' bracket in:\n" + lnOrig)
		} else {
			if needLegalName > 0 {
				needLegalName, name = needLegalName+1, "_"+strings.Replace(ln[pos+1:idx], " ", "_", -1)+"_"+strconv.Itoa(needLegalName)
			}
			ln, m[name] = ln[:pos]+" "+name+" "+ln[idx+1:], ln[pos+1:idx]
		}
	}
	return ln
}

func (me *Prog) preResolveExprs(expr Expr, topDefQName string, topDefBody Expr) Expr {
	switch it := expr.(type) {
	case *ExprFunc:
		it.Body = me.preResolveExprs(it.Body, topDefQName, topDefBody)
	case *ExprCall:
		it.CallArg, it.Callee = me.preResolveExprs(it.CallArg, topDefQName, topDefBody), me.preResolveExprs(it.Callee, topDefQName, topDefBody)
		if call, _ := it.Callee.(*ExprCall); call != nil {
			if numlit, _ := it.CallArg.(*ExprLitNum); numlit != nil && call.ifConstNumArithOpInstrThenPreCalcInto(numlit, it) {
				return numlit
			}
		} else if fn, _ := it.Callee.(*ExprFunc); fn != nil && fn.isIdentity() {
			return it.CallArg
		}
	case *ExprName:
		if it.idxOrInstr <= 0 {
			const stdpref = StdModuleName + "."
			topdefbody := me.TopDefs[it.NameVal]
			if topdefbody == nil {
				topdefbody = me.TopDefs[stdpref+it.NameVal]
			}
			if it.idxOrInstr == 0 {
				if topdefbody == nil && strings.HasPrefix(it.NameVal, stdpref) && strings.LastIndexByte(it.NameVal, '.') == len(StdModuleName) {
					needle := strings.TrimPrefix(it.NameVal, stdpref)
					for name, expr := range me.TopDefs {
						if strings.HasPrefix(name, stdpref) && strings.HasSuffix(name, "."+needle) {
							topdefbody = expr
							break
						}
					}
				}
				if topdefbody == nil {
					panic("in '" + topDefQName + "', line " + strconv.Itoa(it.srcLocLineNr) + ": name '" + it.NameVal + "' unresolvable in:\n" + topDefBody.String())
				} else if it.NameVal == topDefQName {
					panic(it.locStr() + "NEW BUG in `Prog.preResolveExprs` for top-level def '" + it.NameVal + "' recursion")
				} else if name, _ := topdefbody.(*ExprName); name != nil {
					return me.preResolveExprs(name, it.NameVal, name)
				} else {
					return topdefbody
				}
			} else if topdefbody != nil {
				panic("in '" + topDefQName + "', line " + strconv.Itoa(it.srcLocLineNr) + ": local name '" + it.NameVal + "' already taken (no shadowing allowed)")
			}
		}
	}
	return expr
}

func (me *ExprFunc) isIdentity() bool {
	name, ok := me.Body.(*ExprName)
	return ok && name.idxOrInstr == -1
}

func (me *ExprCall) ifConstNumArithOpInstrThenPreCalcInto(rhs *ExprLitNum, parent *ExprCall) (ok bool) {
	if name, _ := me.Callee.(*ExprName); name != nil && name.idxOrInstr > 0 {
		if lhs, _ := me.CallArg.(*ExprLitNum); lhs != nil {
			if instr := instr(name.idxOrInstr); instr < instrEQ {
				ok, rhs.nodeLocInfo, rhs.NumVal = true, parent.nodeLocInfo, int(instr.callCalc(parent.nodeLocInfo, valNum(lhs.NumVal), valNum(rhs.NumVal)))
			}
		}
	}
	return
}
