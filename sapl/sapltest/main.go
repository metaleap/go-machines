package main

import (
	"io/ioutil"
	"os"

	"github.com/metaleap/go-machines/sapl"
)

func main() {
	src, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	prog := sapl.LoadFromJson(src)
	println(prog.Eval(prog[len(prog)-1].Expr).String())
}
