package clsyn

type SynMod struct {
	Defs []*SynDef
}

func (this *SynMod) Def(name string) *SynDef {
	if i := this.IndexOf(name); i > -1 {
		return this.Defs[i]
	}
	return nil
}

func (this *SynMod) IndexOf(name string) int {
	for i, def := range this.Defs {
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
