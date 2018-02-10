package climpl

import (
	"github.com/metaleap/go-corelang/util"
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
