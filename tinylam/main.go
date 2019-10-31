package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	var prog Prog
	if err == nil {
		srcs := make(map[string][]byte, len(files))
		for _, file := range files {
			if idxdot := strings.LastIndexByte(file.Name(), '.'); idxdot > 0 && file.Name()[idxdot:] == ".tl" && !file.IsDir() {
				if srcs[file.Name()[:idxdot]], err = ioutil.ReadFile(filepath.Join(os.Args[1], file.Name())); err != nil {
					panic(err)
				}
			}
		}
		prog = Load(srcs)
		jsonout := json.NewEncoder(os.Stdout)
		jsonout.SetIndent("", "  ")
		err = jsonout.Encode(prog)
	}
	if err != nil {
		panic(err)
	}
}
