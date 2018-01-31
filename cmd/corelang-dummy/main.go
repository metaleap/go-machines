package main

import (
	"github.com/go-leap/dev/lex"
	"github.com/metaleap/go-corelang"
)

func main() {
	var interp corelang.IInterpreter = &corelang.InterpPrettyPrint{}
	for name, def := range corelang.PreludeDefs {
		println("\n\n" + name + ":\n")
		if result, err := interp.Def(def); err != nil {
			panic(err)
		} else {
			println(result.(string))
		}
	}

	_, errs := udevlex.Lex("dummy.foo", "foo 'x' \t\t\t    123 \t and \r\n \"str1\" /*and*/ `str2` *+ - 'd' <:> `\n` // noice")
	for _, e := range errs {
		println(e.Error())
	}

}
