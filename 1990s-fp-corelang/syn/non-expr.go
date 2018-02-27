package clsyn

type SynMod struct {
	Defs []*SynDef
}

func (me *SynMod) Def(name string) *SynDef {
	if i := me.IndexOf(name); i > -1 {
		return me.Defs[i]
	}
	return nil
}

func (me *SynMod) IndexOf(name string) int {
	for i, def := range me.Defs {
		if def.Name == name {
			return i
		}
	}
	return -1
}

type SynDef struct {
	syn
	Name     string
	Args     []string
	Body     IExpr
	TopLevel bool
}

type SynCaseAlt struct {
	syn
	Tag   string
	Binds []string
	Body  IExpr
}
