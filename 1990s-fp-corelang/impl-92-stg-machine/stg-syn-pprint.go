package climpl

import (
	"strconv"
	"strings"
)

func (me synMod) String() (s string) {
	for i := range me.Binds {
		if uglyHackyIndent = 0; i > 0 {
			s += "\n\n"
		}
		s += me.Binds[i].String()
	}
	return
}

func (me synBinding) String() (s string) {
	if s = me.Name; me.LamForm.Upd {
		s += " UPD"
	}
	if len(me.LamForm.Free) > 0 {
		s += " ‹"
		for i := range me.LamForm.Free {
			if i > 0 {
				s += ","
			}
			s += me.LamForm.Free[i].String()
		}
		s += "›"
	}
	s += " = \\"
	for i := range me.LamForm.Args {
		s += " " + me.LamForm.Args[i].String()
	}
	s += " -> "
	if _, isatomic := me.LamForm.Body.(iSynExprAtom); isatomic {
		s += me.LamForm.Body.String()
	} else {
		uglyHackyIndent++
		s += "\n" + strings.Repeat("  ", uglyHackyIndent) + me.LamForm.Body.String()
		uglyHackyIndent--
	}
	return
}

func (me synExprAtomIdent) String() string { return me.Name }

func (me synExprAtomLitFloat) String() string { return strconv.FormatFloat(me.Lit, 'g', -1, 64) }

func (me synExprAtomLitUInt) String() string { return strconv.FormatUint(me.Lit, 10) }

func (me synExprAtomLitRune) String() string { return strconv.QuoteRune(me.Lit) }

func (me synExprAtomLitText) String() string { return strconv.Quote(me.Lit) }

func (me synExprLet) String() (s string) {
	s = "LET"
	if me.Rec {
		s += " REC"
	}
	uglyHackyIndent++
	for i := range me.Binds {
		s += "\n" + strings.Repeat("  ", uglyHackyIndent) + me.Binds[i].String()
	}
	uglyHackyIndent--
	s += "\n" + strings.Repeat("  ", uglyHackyIndent) + "IN\n"
	uglyHackyIndent++
	s += strings.Repeat("  ", uglyHackyIndent) + me.Body.String()
	uglyHackyIndent--
	return
}

func (me synExprCall) String() (s string) {
	s = "(" + me.Callee.String()
	for i := range me.Args {
		s += " " + me.Args[i].String()
	}
	s += ")"
	return
}

func (me synExprCtor) String() (s string) {
	s += "«" + me.Tag.String()
	for i := range me.Args {
		if me.Args[i] == nil {
			s += " NIL"
		} else {
			s += " " + me.Args[i].String()
		}
	}
	s += "»"
	return
}

func (me synExprPrimOp) String() string {
	return "(" + me.Left.String() + " " + me.PrimOp + " " + me.Right.String() + ")"
}

func (me synExprCaseOf) String() (s string) {
	s = "CASE " + me.Scrut.String() + " OF"
	uglyHackyIndent++
	for i := range me.Alts {
		s += "\n" + strings.Repeat("  ", uglyHackyIndent) + me.Alts[i].String()
	}
	uglyHackyIndent--
	return
}

func (me synCaseAlt) String() (s string) {
	if me.Atom != nil {
		s = me.Atom.String()
	} else if me.Ctor.Tag.Name != "" {
		s = "«" + me.Ctor.Tag.String()
		for i := range me.Ctor.Vars {
			s += " " + me.Ctor.Vars[i].String()
		}
		s += "»"
	} else {
		s = "_"
	}
	s += " -> " + me.Body.String()
	return
}
