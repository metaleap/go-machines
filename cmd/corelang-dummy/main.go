package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/metaleap/go-corelang"
	// "github.com/metaleap/go-corelang/impl-00-naive"
	// "github.com/metaleap/go-corelang/impl-01-tmplinst"
	"github.com/metaleap/go-corelang/impl-02-gmachine"
	"github.com/metaleap/go-corelang/syn"
	"github.com/metaleap/go-corelang/util"
)

func writeLn(s string) { _, _ = os.Stdout.WriteString(s + "\n") }

func main() {
	mod := &clsyn.SynMod{Defs: corelang.PreludeDefs}
	if !lexAndParse("from-const-srcMod-in.dummy-mod-src.go", srcMod, mod) {
		return
	}

	writeLn("module lexed and parsed, globals are:")
	for defname := range mod.Defs {
		writeLn("\t" + defname)
	}
	machine := recompile(mod)

	multiline, repl, pprint := "", bufio.NewScanner(os.Stdin), &corelang.SyntaxTreePrinter{}
	for repl.Scan() {
		if readln := strings.TrimSpace(repl.Text()); readln != "" {
			if readln == "…" && multiline != "" {
				readln, multiline = strings.TrimSpace(multiline), ""
			}
			switch {
			case strings.HasSuffix(readln, "…"):
				multiline = readln[:len(readln)-len("…")] + "\n  "
			case multiline != "":
				multiline += readln + "\n  "
			case !strings.Contains(readln, "="):
				if readln == "*" || readln == "?" {
					for defname := range mod.Defs {
						writeLn(defname)
					}
				} else if strings.HasPrefix(readln, "!") {
					defname, starttime := readln[1:], time.Now()
					val, stats, evalerr := machine.Eval(defname)
					timetaken := time.Now().Sub(starttime)
					if evalerr != nil {
						println(evalerr.Error())
					} else {
						fmt.Printf("Reduced in %v (%d appls / %d steps) to:\n%s\n", timetaken, stats.NumAppls, stats.NumSteps, machine.String(val))
					}
				} else if def := mod.Defs[readln]; def == nil {
					println("not found: " + readln)
				} else {
					srcfmt, _ := pprint.Def(mod.Defs[readln])
					writeLn(srcfmt)
				}
			case lexAndParse("<input>", readln, mod):
				machine = recompile(mod)
			}
		}
	}
}

func recompile(mod *clsyn.SynMod) clutil.IMachine {
	timestarted := time.Now()
	machine, errs := climpl.CompileToMachine(mod)
	timetaken := time.Now().Sub(timestarted)

	for _, err := range errs {
		println(err.Error())
	}
	fmt.Printf("module re-compiled in %s\n\n", timetaken)
	return machine
}

func lexAndParse(filePath string, src string, mod *clsyn.SynMod) bool {
	defs, errs_parse := clsyn.LexAndParseDefs(filePath, src)

	for _, def := range defs {
		if mod.Defs[def.Name] != nil {
			println("Redefined: " + def.Name)
		}
		mod.Defs[def.Name] = def
	}
	for _, e := range errs_parse {
		println(e.Error())
	}
	return len(errs_parse) == 0
}
