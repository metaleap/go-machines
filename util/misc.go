package clutil

import (
	"errors"
	"fmt"
)

type Stats struct {
	NumSteps int
	NumAppls int
}

func Catch(err *error) {
	switch e := recover().(type) {
	case string:
		*err = errors.New(e)
	case error:
		*err = e
	default:
		if e != nil {
			*err = fmt.Errorf("%T", e)
		}
	}
}
