package cclparse

import (
	//"fmt"
)

func fmin(a, b float64) float64 {
	if a < b { return a }
	return b
}

func fabs(a float64) float64 {
	if a < 0 { return -a }
	return a
}

