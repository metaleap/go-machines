package climpl

import (
	"strconv"
	"strings"
)

var uglyHackyIndent int

func (this synMod) String() (s string) {
	for i := range this.Binds {
		if uglyHackyIndent = 0; i > 0 {
			s += "\n\n"
		}
		s += this.Binds[i].String()
	}
	return
}

func (this synBinding) String() (s string) {
	if s = this.Name; this.LamForm.Upd {
		s += " ¤"
	} else {
		s += " Ø"
	}
	if len(this.LamForm.Free) > 0 {
		s += " ‹"
		for i := range this.LamForm.Free {
			if i > 0 {
				s += ","
			}
			s += this.LamForm.Free[i].String()
		}
		s += "›"
	}
	s += " = \\"
	for i := range this.LamForm.Args {
		s += " " + this.LamForm.Args[i].String()
	}
	s += " -> "
	if _, isatomic := this.LamForm.Body.(iSynExprAtom); isatomic {
		s += this.LamForm.Body.String()
	} else {
		uglyHackyIndent++
		s += "\n" + strings.Repeat("  ", uglyHackyIndent) + this.LamForm.Body.String()
		uglyHackyIndent--
	}
	return
}

func (this synExprAtomIdent) String() string    { return this.Name }
func (this synExprAtomLitFloat) String() string { return strconv.FormatFloat(this.Lit, 'g', -1, 64) }
func (this synExprAtomLitUInt) String() string  { return strconv.FormatUint(this.Lit, 10) }
func (this synExprAtomLitRune) String() string  { return strconv.QuoteRune(this.Lit) }
func (this synExprAtomLitText) String() string  { return strconv.Quote(this.Lit) }

func (this synExprLet) String() (s string) {
	s = "LET"
	if this.Rec {
		s += " REC"
	}
	uglyHackyIndent++
	for i := range this.Binds {
		s += "\n" + strings.Repeat("  ", uglyHackyIndent) + this.Binds[i].String()
	}
	uglyHackyIndent--
	s += "\n" + strings.Repeat("  ", uglyHackyIndent) + "IN\n"
	uglyHackyIndent++
	s += strings.Repeat("  ", uglyHackyIndent) + this.Body.String()
	uglyHackyIndent--
	return
}

func (this synExprCall) String() (s string) {
	s = "(" + this.Callee.String()
	for i := range this.Args {
		s += " " + this.Args[i].String()
	}
	s += ")"
	return
}

func (this synExprCtor) String() (s string) {
	s += "‹" + this.Tag.String()
	for i := range this.Args {
		s += " " + this.Args[i].String()
	}
	s += "›"
	return
}

func (this synExprPrimOp) String() string {
	return "(" + this.Left.String() + " " + this.PrimOp + " " + this.Right.String() + ")"
}

func (this synExprCaseOf) String() (s string) {
	s = "CASE " + this.Scrut.String() + " OF"
	uglyHackyIndent++
	for i := range this.Alts {
		s += "\n" + strings.Repeat("  ", uglyHackyIndent) + this.Alts[i].String()
	}
	uglyHackyIndent--
	return
}

func (this synCaseAlt) String() (s string) {
	if this.Atom != nil {
		s = this.Atom.String()
	} else if this.Ctor.Tag.Name != "" {
		s = "‹" + this.Ctor.Tag.String()
		for i := range this.Ctor.Vars {
			s += " " + this.Ctor.Vars[i].String()
		}
		s += "›"
	} else {
		s = "_"
	}
	s += " -> " + this.Body.String()
	return
}
