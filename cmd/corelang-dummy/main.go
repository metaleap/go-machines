package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	multiline, repl, pprint := "", bufio.NewScanner(os.Stdin), &corelang.InterpPrettyPrint{}
	for defname := range mod.Defs {
		writeLn(defname)
	}
	var machine clutil.IMachine = climpl.CompileToMachine(mod)
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
					val, numsteps, evalerr := machine.Eval(readln[1:])
					if evalerr != nil {
						println(evalerr.Error())
					} else {
						fmt.Printf("Reduced in %d steps to:\n%v\n", numsteps, val)
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
