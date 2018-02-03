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

func ParseDefs(srcFilePath string, parent ISyn, tokens []lex.IToken) (defs []*SynDef, errs []*Error) {
	for len(tokens) > 0 {
		if def, tail, deferr := parseDef(parent, tokens); deferr != nil {
			deferr.Pos.Filename = srcFilePath
			defs, tokens, errs = nil, nil, append(errs, deferr)
		} else if len(errs) == 0 {
			defs, tokens = append(defs, def), tail
		}
	}
	return
}

func parseDef(parent ISyn, tokens []lex.IToken) (*SynDef, []lex.IToken, *Error) {
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
	for insig := true; insig && i < len(tokens); i++ {
		if t, _ := tokens[i].(*lex.TokenOther); t != nil && t.Token == "=" {
			insig = false // dont break, still want to inc i
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
