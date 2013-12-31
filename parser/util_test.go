package parser

import (
	gc "launchpad.net/gocheck"
)

type UtilSuite struct{}

var _ = gc.Suite(&UtilSuite{})

func (s *UtilSuite) TestInvariant(c *gc.C) {
	invariant(true)
	c.Check(func() {invariant(false)}, gc.PanicMatches,
		"Internal consistency check failed.")
}
