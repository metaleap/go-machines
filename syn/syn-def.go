package clsyn

import (
	"fmt"

	"github.com/go-leap/dev/lex"
)

type SynDef struct {
	syn
	Name string
	Args []string
	Body IExpr
}

func LexedTokensToTopLevelChunks(tokens []udevlex.IToken) (topLevelTokenChunks [][]udevlex.IToken) {
	var cur int
	for i, ln, l := 0, 1, len(tokens); i < l; i++ {
		if tpos := tokens[i].Meta(); i == l-1 {
			if tlc := tokens[cur:]; len(tlc) > 0 {
				topLevelTokenChunks = append(topLevelTokenChunks, tlc)
			}
		} else if tpos.LineIndent == 0 && tpos.Line != ln {
			if tlc := tokens[cur:i]; len(tlc) > 0 {
				topLevelTokenChunks = append(topLevelTokenChunks, tlc)
			}
			cur, ln = i, tpos.Line
		}
	}
	return
}

func ParseDefs(srcFilePath string, topLevelTokenChunks [][]udevlex.IToken) (defs []*SynDef, errs []*Error) {
	for _, topleveltokenchunk := range topLevelTokenChunks {
		if def, deferr := parseDef(nil, topleveltokenchunk); deferr != nil {
			deferr.Pos.Filename = srcFilePath
			defs, errs = nil, append(errs, deferr)
		} else if len(errs) == 0 {
			defs = append(defs, def)
		}
	}
	return
}

func parseDef(parent ISyn, tokens []udevlex.IToken) (*SynDef, *Error) {
	if len(tokens) < 3 {
		return nil, errPos(udevlex.Pos(tokens, parent, ""), "not enough tokens to form a definition", 0)
	}

	tid, _ := tokens[0].(*udevlex.TokenIdent)
	if tid == nil {
		return nil, errPos(tokens[0], fmt.Sprintf("expected identifier instead of `%s`", tokens[0]), len(tokens[0].String()))
	}

	i, def := 1, &SynDef{Name: tid.Token}
	def.syn.pos, def.syn.parent = tid.TokenMeta, parent

	// args up until `=`
	for insig := true; insig && i < len(tokens); i++ {
		if t, _ := tokens[i].(*udevlex.TokenOther); t != nil && t.Token == "=" {
			insig = false // dont break, still want to inc i
		} else if t, _ := tokens[i].(*udevlex.TokenIdent); t != nil {
			def.Args = append(def.Args, t.Token)
		} else {
			return nil, errPos(tokens[i], fmt.Sprintf("expected argument name or `=` instead of `%s`", tokens[i]), len(tokens[i].String()))
		}
	}

	// body of definition after `=`
	expr, tail, err := parseExpr(def, tokens[i:])
	if err != nil {
		return nil, err
	}
	if len(tail) > 0 {
		println("TODO..")
	}
	def.Body = expr
	return def, nil
}
