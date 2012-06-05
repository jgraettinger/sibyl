package cclparse

import (
	//"fmt"
    "invariant"
)

func fmin(a, b float64) float64 {
	if a < b { return a }
	return b
}

func fabs(a float64) float64 {
	if a < 0 { return -a }
	return a
}

func iabs(a int) int {
    if a < 0 { return -a }
    return a
}

func isign(a int) int {
    invariant.NotEqual(a, 0)
    if a < 0 { return -1 }
    return 1
}
