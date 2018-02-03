package clsyn

import (
	"fmt"

	lex "github.com/go-leap/dev/lex"
)

type SynDef struct {
	syn
	Name string
	Args []string
	Body IExpr
}

func ParseDefs(srcFilePath string, tokens lex.Tokens) (defs []*SynDef, errs []*Error) {
	defs, errs = parseDefs(nil, tokens)
	for _, e := range errs {
		e.Pos.Filename = srcFilePath
	}
	return
}

func parseDefs(parent ISyn, tokens lex.Tokens) (defs []*SynDef, errs []*Error) {
	for len(tokens) > 0 {
		def, tail, deferr := parseDef(parent, tokens)
		if tokens = tail; deferr != nil {
			errs = append(errs, deferr)
		} else if len(errs) == 0 {
			defs = append(defs, def)
		}
	}
	if len(errs) > 0 {
		defs = nil
	}
	return
}

func parseDef(parent ISyn, tokens lex.Tokens) (*SynDef, lex.Tokens, *Error) {
	if len(tokens) < 3 {
		return nil, nil, errPos(lex.Pos(tokens, parent, ""), "not enough tokens to form a definition", 0)
	}

	tid, _ := tokens[0].(*lex.TokenIdent)
	if tid == nil {
		return nil, nil, errPos(tokens[0], fmt.Sprintf("expected identifier instead of `%s`", tokens[0]), len(tokens[0].String()))
	}

	i, def := 1, &SynDef{Name: tid.Token}
	def.syn.pos, def.syn.parent = tid.TokenMeta, parent

	// args up until `=`
	for ; i < len(tokens); i++ {
		if t, _ := tokens[i].(*lex.TokenOther); t != nil && t.Token == "=" {
			i++
			break
		} else if t, _ := tokens[i].(*lex.TokenIdent); t != nil {
			def.Args = append(def.Args, t.Token)
		} else {
			return nil, nil, errPos(tokens[i], fmt.Sprintf("expected argument name or `=` instead of `%s`", tokens[i]), len(tokens[i].String()))
		}
	}

	// body of definition after `=`
	expr, tail, err := parseExpr(def, tokens[i:])
	if err != nil {
		return nil, nil, err
	}
	def.Body = expr
	return def, tail, nil
}
