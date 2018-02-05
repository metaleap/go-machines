package clutil

import (
	"errors"
	"fmt"
)

func Catch(err *error) {
	switch e := recover().(type) {
	case string:
		*err = errors.New(e)
	case error:
		*err = e
	default:
		if e != nil {
			*err = fmt.Errorf("%v", e)
		}
	}
}
