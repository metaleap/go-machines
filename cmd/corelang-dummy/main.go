package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/metaleap/go-machines"
	// "github.com/metaleap/go-machines/impl-01-tmplinst-buggy"
	"github.com/metaleap/go-machines/impl-02-gmachine-mark7"
	"github.com/metaleap/go-machines/syn"
	"github.com/metaleap/go-machines/util"
)

func writeLn(s string) { _, _ = os.Stdout.WriteString(s + "\n") }

func main() {
	fname, mod := "from `srcMod` in `dummy-mod-src.go`", &clsyn.SynMod{Defs: corelang.PreludeDefs}
	if !lexAndParse(fname, srcMod, mod) {
		return
	}

	writeLn("module " + fname + " lexed and parsed, globals are:\n")
	for defname := range mod.Defs {
		_, _ = os.Stdout.WriteString(" · " + defname)
	}
	writeLn("\n\n➜ enter any name to pretty-print the (parsed) AST")
	writeLn("\n➜ define new globals via `name = expr`, `name x y z = expr` etc (any amount of args is fine)")
	writeLn("\n➜ the following globals have no args:")
	for _, def := range mod.Defs {
		if len(def.Args) == 0 {
			_, _ = os.Stdout.WriteString(" · " + def.Name)
		}
	}
	writeLn("\n…and can be evaluated easily using `!<name>` or `?<name>`\n")
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
				} else if readln[0] == '!' || readln[0] == '?' {
					defname, starttime := readln[1:], time.Now()
					val, stats, evalerr := machine.Eval(defname)
					timetaken := time.Now().Sub(starttime)
					if evalerr != nil {
						println(evalerr.Error())
					} else {
						fmt.Printf("Reduced in %v (%d appls / %d steps / S%d / H%d) to:\n%s\n", timetaken, stats.NumAppls, stats.NumSteps, stats.MaxStack, stats.HeapSize, machine.String(val))
					}
				} else if def := mod.Defs[readln]; def == nil {
					println("not found: " + readln)
				} else {
					writeLn(pprint.Def(mod.Defs[readln]))
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
