package main

import (
	"bufio"
	"os"
	"strings"

	core "github.com/metaleap/go-corelang"
	coresyn "github.com/metaleap/go-corelang/syn"
)

func writeLn(s string) { _, _ = os.Stdout.WriteString(s + "\n") }

func main() {
	mod := &coresyn.SynMod{Defs: core.PreludeDefs}
	if e := lexAndParse("", srcMod, mod); e != nil {
		println(e.Error())
	}

	multiline, repl, pprint := "", bufio.NewScanner(os.Stdin), &core.InterpPrettyPrint{}
	for defname := range mod.Defs {
		writeLn(defname)
	}
	for repl.Scan() {
		if errscan := repl.Err(); errscan != nil {
			panic(errscan)
		} else if readln := strings.TrimRight(repl.Text(), " \t\r\n\v\b"); readln != "" {
			if readln == "…" && multiline != "" {
				readln, multiline = strings.TrimSpace(multiline), ""
			}
			if readln != "…" && multiline == "" && strings.HasSuffix(readln, "…") {
				multiline = readln[:len(readln)-len("…")] + "\n  "
			} else if multiline != "" {
				multiline += readln + "\n  "
			} else if !strings.Contains(readln, "=") {
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
			} else if e := lexAndParse("<input>", readln, mod); e != nil {
				println(e.Error())
			}
		}
	}
}

func lexAndParse(filePath string, src string, mod *coresyn.SynMod) error {
	if filePath == "" {
		filePath = "dummy-mod-src.go"
	}

	defs, errs_parse := coresyn.LexAndParseDefs(filePath, src)
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
