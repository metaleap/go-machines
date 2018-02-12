package clsyn

import (
	"github.com/metaleap/go-corelang/util"
)

func LookupEnvFrom(defs []*SynDef, globals clutil.Env, argsEnv map[string]int, otherNames []string) (lookupEnv map[string]bool) {
	lookupEnv = make(map[string]bool, len(defs)+len(otherNames))
	for _, def := range defs {
		lookupEnv[def.Name] = true
	}
	for name := range globals {
		lookupEnv[name] = true
	}
	for name := range argsEnv {
		lookupEnv[name] = true
	}
	for _, name := range otherNames {
		lookupEnv[name] = true
	}
	return
}
