package climpl

import (
	"github.com/metaleap/go-machines/1990s-fp-corelang/util"
)

type nodeInt int

type nodeAppl struct {
	Callee clutil.Addr
	Arg    clutil.Addr
}

type nodeGlobal struct {
	NumArgs int
	Code    code
}

type nodeIndirection struct {
	Addr clutil.Addr
}

type nodeCtor struct {
	Tag   int
	Items []clutil.Addr
}
