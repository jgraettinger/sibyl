package ccl

import (
    "github.com/dademurphy/sibyl/invariant"
)

func Fmin(a, b float64) float64 {
	if a < b { return a }
	return b
}

func Fabs(a float64) float64 {
	if a < 0 { return -a }
	return a
}

func Iabs(a int) int {
    if a < 0 { return -a }
    return a
}

func Isign(a int) int {
    invariant.NotEqual(a, 0)
    if a < 0 { return -1 }
    return 1
}
