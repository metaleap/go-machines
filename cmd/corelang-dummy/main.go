package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/go-leap/dev/lex"
	"github.com/metaleap/go-corelang"
	"github.com/metaleap/go-corelang/syn"
)

func writeLn(s string) { _, _ = os.Stdout.WriteString(s + "\n") }

func main() {
	mod := &clsyn.SynMod{Defs: corelang.PreludeDefs}
	if e := lexAndParse(srcMod, mod); e != nil {
		panic(e)
	}

	repl, pprint := bufio.NewScanner(os.Stdin), &corelang.InterpPrettyPrint{}
	writeLn("Ready.")
	for repl.Scan() {
		if errscan := repl.Err(); errscan != nil {
			panic(errscan)
		} else if readln := strings.TrimSpace(repl.Text()); readln != "" {
			if !strings.Contains(readln, "=") {
				if readln == "*" {
					for defname := range mod.Defs {
						writeLn(defname)
					}
				} else if def := mod.Defs[readln]; def == nil {
					println("not found: " + readln)
				} else {
					srcfmt, _ := pprint.Def(mod.Defs[readln])
					writeLn(srcfmt.(string))
				}
			} else {
				println("coming soon..")
			}
		}
	}
}

func lexAndParse(src string, mod *clsyn.SynMod) error {
	lexed, errs_lex := udevlex.Lex("dummy-mod-src.go", "src")
	for _, e := range errs_lex {
		return e
	}

	defs, errs_parse := clsyn.ParseDefs(lexed)
	for _, e := range errs_parse {
		return e
	}

	for _, def := range defs {
		if mod.Defs[def.Name] != nil {
			println("Redefined: " + def.Name)
		}
		mod.Defs[def.Name] = def
	}
	return nil
}
