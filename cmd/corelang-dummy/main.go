package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/go-leap/dev/lex"
	core "github.com/metaleap/go-corelang"
	coresyn "github.com/metaleap/go-corelang/syn"
)

func writeLn(s string) { _, _ = os.Stdout.WriteString(s + "\n") }

func main() {
	mod := &coresyn.SynMod{Defs: core.PreludeDefs}
	if e := lexAndParse("", srcMod, mod); e != nil {
		println(e.Error())
	}

	repl, pprint := bufio.NewScanner(os.Stdin), &core.InterpPrettyPrint{}
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

func lexAndParse(filePath string, src string, mod *coresyn.SynMod) error {
	if filePath == "" {
		filePath = "dummy-mod-src.go"
	}
	lexed, errs_lex := udevlex.Lex(filePath, src, "(", ")")
	for _, e := range errs_lex {
		return e
	}

	defs, errs_parse := coresyn.ParseDefs(filePath, lexed.SansComments())
	for _, def := range defs {
		if mod.Defs[def.Name] != nil {
			println("Redefined: " + def.Name)
		}
		mod.Defs[def.Name] = def
	}

	for _, e := range errs_parse {
		return e
	}
	return nil
}
