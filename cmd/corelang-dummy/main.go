package main

import (
	"github.com/metaleap/go-corelang"
)

func main() {
	var interp corelang.IInterpreter = &corelang.InterpPrettyPrint{}
	for _, def := range corelang.PreludeDefs {
		println("\n\n")
		if result, err := interp.Def(def); err != nil {
			panic(err)
		} else {
			println(result.(string))
		}
	}
}
