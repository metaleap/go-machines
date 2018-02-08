package climpl

import (
	"github.com/metaleap/go-corelang/util"
)

type nodeLitUint uint64

type nodeAppl struct {
	Callee clutil.Addr
	Arg    clutil.Addr
}

type nodeGlobal struct {
	NumArgs int
	Code    code
}

type nodePointTo struct {
	Addr clutil.Addr
}
