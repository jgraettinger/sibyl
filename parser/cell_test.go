package parser

import (
	"fmt"

	gc "launchpad.net/gocheck"
)

type CellSuite struct{}

var _ = gc.Suite(&CellSuite{})

func (s *CellSuite) TestStringFormatting(c *gc.C) {
	cell := &Cell{Index: 1, Token: "X"}
	c.Check(fmt.Sprintf("%v", cell), gc.Equals, "\"X\"@1")
}

func (s *CellSuite) TestLastCell(c *gc.C) {
	X, Y := &Cell{Index: 0, Token: "X"}, &Cell{Index: 1, Token: "Y"}
	var cells []*Cell
	c.Check(lastCell(cells), gc.IsNil)
	cells = []*Cell{X, Y}
	c.Check(lastCell(cells), gc.Equals, Y)
}
