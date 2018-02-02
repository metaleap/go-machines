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

func LexedTokensToTopLevelChunks(tokens []udevlex.Token) (topLevelTokenChunks [][]udevlex.Token) {
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

func ParseDefs(srcFilePath string, topLevelTokenChunks [][]udevlex.Token) (defs []*SynDef, errs []*Error) {
	for _, topleveltokenchunk := range topLevelTokenChunks {
		if def, deferr := ParseDef(srcFilePath, topleveltokenchunk); deferr != nil {
			defs, errs = nil, append(errs, deferr)
		} else if len(errs) == 0 {
			defs = append(defs, def)
		}
	}
	return
}

func ParseDef(srcFilePath string, tokens []udevlex.Token) (*SynDef, *Error) {
	if len(tokens) < 3 {
		return nil, errPos(&tokens[0].Meta().Position, "not enough tokens to form a definition")
	}
	tid, _ := tokens[0].(*udevlex.TokenIdent)
	if tid == nil {
		return nil, errPos(&tokens[0].Meta().Position, fmt.Sprintf("expected identifier instead of `%s`", tokens[0].String()[4:]))
	}
	def := &SynDef{Name: tid.Token}
	for i, insig := 1, true; insig; i++ {
		if t, _ := tokens[i].(*udevlex.TokenOther); t != nil && t.Token == "=" {
			insig = false
		} else if t, _ := tokens[i].(*udevlex.TokenIdent); t != nil {
			def.Args = append(def.Args, t.Token)
		} else {
			return nil, errPos(&tokens[0].Meta().Position, fmt.Sprintf("expected argument name or `=` instead of `%s`", tokens[i].String()[4:]))
		}
	}

	return def, nil
}
