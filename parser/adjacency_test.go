package parser

import (
	"fmt"

	gc "launchpad.net/gocheck"
)

type AdjacencySuite struct{}

var _ = gc.Suite(&AdjacencySuite{})

func (s *AdjacencySuite) TestStringFormatting(c *gc.C) {
	X, Y := &Cell{Index: 0, Token: "X"}, &Cell{Index: 1, Token: "Y"}
	a := &Adjacency{Head: Y, Tail: X, Position: -1}

	c.Check(fmt.Sprintf("%v", a), gc.Equals, "\"Y\"@1:-1 => \"X\"@0")
}
func (s *AdjacencySuite) TestHeadAndTailSide(c *gc.C) {
	X, Y := &Cell{Index: 0, Token: "X"}, &Cell{Index: 1, Token: "Y"}

	a := &Adjacency{Head: X, Tail: Y, Position: 1}
	c.Check(a.HeadSide(), gc.Equals, &X.Right)
	c.Check(a.TailSide(), gc.Equals, &Y.Left)

	a = &Adjacency{Head: Y, Tail: X, Position: -1}
	c.Check(a.HeadSide(), gc.Equals, &Y.Left)
	c.Check(a.TailSide(), gc.Equals, &X.Right)
}
