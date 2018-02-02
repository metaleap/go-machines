package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/metaleap/go-corelang"
	"github.com/metaleap/go-corelang/syn"
)

func main() {
	mod := clsyn.SynMod{Defs: corelang.PreludeDefs}

	repl, pprint := bufio.NewScanner(os.Stdin), &corelang.InterpPrettyPrint{}
	for repl.Scan() {
		if errscan := repl.Err(); errscan != nil {
			panic(errscan)
		} else if readln := strings.TrimSpace(repl.Text()); readln != "" {
			if !strings.Contains(readln, "=") {
				if readln == "*" {
					for defname := range mod.Defs {
						println(defname)
					}
				} else if def := mod.Defs[readln]; def == nil {
					println("not found: " + readln)
				} else {
					srcfmt, _ := pprint.Def(mod.Defs[readln])
					println(srcfmt.(string))
				}
			} else {
				println("coming soon..")
			}
		}
	}
}
