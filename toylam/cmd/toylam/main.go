package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	tl "github.com/metaleap/go-machines/toylam"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		os.Stdout.WriteString("USAGE: toylam [--lazy] [<file-with-main>]\n")
		return
	}
	argpos, lazyeval := 1, len(os.Args) > 1 && os.Args[1] == "--lazy"
	if lazyeval {
		argpos = 2
	}

	var srcfilepath, srcdirpath string
	if argpos >= len(os.Args) {
		srcdirpath, argpos = ".", argpos-1
	} else {
		srcfilepath = os.Args[argpos]
		if stat, _ := os.Stat(srcfilepath); stat != nil && stat.IsDir() {
			srcdirpath, srcfilepath = srcfilepath, ""
		} else {
			srcdirpath = filepath.Dir(srcfilepath)
		}
	}

	files, _ := ioutil.ReadDir(srcdirpath)
	modules := make(map[string][]byte, len(files))
	for _, file := range files {
		if curfilepath := filepath.Join(srcdirpath, file.Name()); !file.IsDir() {
			if idxdot := strings.LastIndexByte(file.Name(), '.'); (curfilepath == srcfilepath) || (idxdot > 0 && file.Name()[idxdot:] == ".tl") {
				if src, err := ioutil.ReadFile(curfilepath); err == nil {
					modules[file.Name()[:idxdot]] = src
				} else {
					panic(err)
				}
			}
		}
	}
	var srcfilename, srcfileext string
	if len(srcfilepath) != 0 {
		srcfilename, srcfileext = filepath.Base(srcfilepath), filepath.Ext(srcfilepath)
	}
	if len(modules) == 0 {
		panic("neither `" + srcfilename + "` nor any other `.tl` source files found in: " + srcdirpath)
	}

	prog, maintopdefqname, strdividerline := tl.Prog{LazyEval: lazyeval}, srcfilename[:len(srcfilename)-len(srcfileext)]+".main", "\n────────────────────────────────────────────────────────────────────────────────"
	prog.ParseModules(modules, tl.ParseOpts{})
	prog.OnInstrMSG = func(msg string, val tl.Value) { println("LOG: " + msg + "\t" + prog.Value(val).String()) }
	if maintopdefbody := prog.TopDefs[maintopdefqname]; maintopdefbody != nil {
		if retval := prog.RunAsMain(maintopdefbody, os.Args[argpos+1:]); retval != nil {
			if bytes, ok := tl.ValueBytes(retval); ok {
				_, _ = os.Stdout.Write(append(bytes, '\n'))
			} else {
				_, _ = os.Stdout.WriteString(retval.String() + "\n")
			}
		}
		os.Exit(0)
	} else if srcfilename != "" {
		panic("no such global top-level def: " + maintopdefqname)
	}
	{ /* REPL */
		os.Stdout.WriteString("Ctrl+C to quit this REPL." + strdividerline + "\n")
		readln, eval := bufio.NewScanner(os.Stdin), func(ln string) (retval tl.Value, err interface{}) {
			defer func() { err = recover() }()
			modules["<repl>"] = []byte("<input> := " + ln)
			prog.ParseModules(modules, tl.ParseOpts{}) // _technically_ very inefficient to reload-it-all on every single input but "works smoothly enough for me for now" --- the goal of toylam was to stay "tiny in terms of LoCs". which we already failed. even whackier is that the original `ParseModules` already did rewrite sources in our `module` map.
			val := prog.Eval(prog.TopDefs["<repl>.<input>"], nil)
			retval = prog.Value(val)
			println("STEPS", prog.NumEvalSteps)
			return
		}
		for readln.Scan() {
			if input := strings.TrimSpace(readln.Text()); input != "" {
				if retval, err := eval(input); err == nil {
					os.Stdout.WriteString(retval.String() + strdividerline + "\n")
				} else if errval, ok := err.(tl.Value); !ok {
					os.Stderr.WriteString(fmt.Sprintf("%v%s\n", err, strdividerline))
				} else {
					os.Stderr.WriteString(fmt.Sprintf("ERR: %s%s\n", errval, strdividerline))
				}
			}
		}
	}
}

func stdinReadSplitterBy(sep byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.IndexByte(data, sep); i >= 0 {
			advance, token = i+1, data[0:i]
		} else if atEOF {
			advance, token = len(data), data
		}
		return
	}
}
