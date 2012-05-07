package cclparse

import (
	"fmt"
)

func invariant(check bool, a ...interface{}) {
	if check {
		return
	}
	if len(a) != 0 {
		errf := a[0].(string)

		panic(fmt.Sprintf(errf, a[1:]...))
	} else {
		panic("Invariant check failed")
	}
}

