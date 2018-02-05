package clutil

import (
	"errors"
)

func Catch(err *error) {
	if errmsg, _ := recover().(string); errmsg != "" {
		*err = errors.New(errmsg)
	}
}
