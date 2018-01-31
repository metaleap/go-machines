package corelang

import (
	lex "github.com/go-leap/dev/lex"
)

type parse struct {
	syn  iSyn
	tail []lex.Token
}

type parser func([]lex.Token) []*parse
