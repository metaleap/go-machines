package main

import (
	"io/ioutil"
	"os"
	"strings"

	tl "github.com/metaleap/tinylam/go"
)

func main() {
	files, err := ioutil.ReadDir(".")
	var prog tl.Prog
	if err == nil {
		srcs := make(map[string][]byte, len(files))
		for _, file := range files {
			if idxdot := strings.LastIndexByte(file.Name(), '.'); idxdot > 0 && file.Name()[idxdot:] == ".tl" && !file.IsDir() {
				if srcs[file.Name()[:idxdot]], err = ioutil.ReadFile(file.Name()); err != nil {
					panic(err)
				}
			}
		}
		if len(os.Args) <= 1 {
			panic("missing argument: name of prog to run")
		}
		defqname := os.Args[1]
		prog = tl.Load(srcs, os.Args[2:])
		if strings.IndexByte(defqname, '.') < 0 {
			defqname = "appdemo." + defqname + ".main"
		}
		if def, ok := prog[defqname]; !ok {
			panic("unknown: " + defqname)
		} else {
			println(def.Str())
		}
		println("__________________")
		println(prog.Run(defqname, make([]byte, 321)).Str())
	}
	if err != nil {
		panic(err)
	}
}
