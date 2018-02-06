package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	core "github.com/metaleap/go-corelang"
	"github.com/metaleap/go-corelang/impl-01-tmplinst"
	coresyn "github.com/metaleap/go-corelang/syn"
)

func writeLn(s string) { _, _ = os.Stdout.WriteString(s + "\n") }

func main() {
	mod := &coresyn.SynMod{Defs: core.PreludeDefs}
	if !lexAndParse("from-const-srcMod-in.dummy-mod-src.go", srcMod, mod) {
		return
	}

	multiline, repl, pprint := "", bufio.NewScanner(os.Stdin), &core.InterpPrettyPrint{}
	for defname := range mod.Defs {
		writeLn(defname)
	}
	machine := climpl.CompileToMachine(mod)
	for repl.Scan() {
		if readln := strings.TrimSpace(repl.Text()); readln != "" {
			if readln == "…" && multiline != "" {
				readln, multiline = strings.TrimSpace(multiline), ""
			}
			if strings.HasSuffix(readln, "…") {
				multiline = readln[:len(readln)-len("…")] + "\n  "
			} else if multiline != "" {
				multiline += readln + "\n  "
			} else if !strings.Contains(readln, "=") { // will do until we introduce == / != / <= / >= / >>= etc ;)
				if readln == "*" || readln == "?" {
					for defname := range mod.Defs {
						writeLn(defname)
					}
				} else if strings.HasPrefix(readln, "!") {
					evalerr := machine.Eval(readln[1:])
					if evalerr != nil {
						println(evalerr.Error())
					} else {
						fmt.Printf("Reduced in %d steps to:\n%v\n", machine.Stats.NumStepsTaken, machine.Heap[machine.Stack[0]])
					}
				} else if def := mod.Defs[readln]; def == nil {
					println("not found: " + readln)
				} else {
					srcfmt, _ := pprint.Def(mod.Defs[readln])
					writeLn(srcfmt.(string))
				}
			} else if lexAndParse("<input>", readln, mod) {
				machine = climpl.CompileToMachine(mod)
				writeLn("all definition successfully parsed, enter its name to pretty-print it's syntax tree")
			}
		}
	}
}

func lexAndParse(filePath string, src string, mod *coresyn.SynMod) bool {
	defs, errs_parse := coresyn.LexAndParseDefs(filePath, src)

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
