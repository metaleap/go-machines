package clutil

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Stats struct {
	NumSteps int
	NumAppls int
}

func Catch(err *error) {
	if maybeerr := recover(); maybeerr != nil {
		switch e := maybeerr.(type) {
		case string:
			*err = errors.New(e)
		case error:
			*err = e
		default:
			*err = fmt.Errorf("%T", e)
		}
	}
}

// these just to occasionally compare compiled perf of factorial with our interpreters

func _fac(n int) int {
	forcesNumToBeUnpredictablyNonConstishByGoRun := os.Getenv("DOESNT_REALLY_EXIST")
	num, _ := strconv.ParseInt(forcesNumToBeUnpredictablyNonConstishByGoRun, 10, 64)
	n = int(num) + n
	timestarted := time.Now()
	n = fac(n)
	timetaken := time.Now().Sub(timestarted)
	fmt.Printf("%v", timetaken) // always around 140-199ns for fac(15)
	return n
}

func fac(n int) int { // cmp vs Core (gMachine interpreter) where (fac 15) always around 430-1354µs — so approx ~3000-10000x slower
	if n == 0 {
		return 1
	}
	return n * fac(n-1)
}
