package clsyn

import (
	"github.com/metaleap/go-machines/util"
)

func NewLookupEnv(defs []*SynDef, globals clutil.Env, argsEnv map[string]int, otherNames []string) (me map[string]bool) {
	me = make(map[string]bool, len(defs)+len(globals)+len(argsEnv)+len(otherNames))
	for _, def := range defs {
		me[def.Name] = true
	}
	for name := range globals {
		me[name] = true
	}
	for name := range argsEnv {
		me[name] = true
	}
	for _, name := range otherNames {
		me[name] = true
	}
	return
}
