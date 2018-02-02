package clsyn

import (
	"github.com/go-leap/dev/lex"
)

type SynDef struct {
	syn
	Name string
	Args []string
	Body IExpr
}

func ParseDefs(tokens []udevlex.Token) (defs []*SynDef, errs []error) {

	return
}
