package parser

import (
	gc "launchpad.net/gocheck"
)

type MontonicitySuite struct{}

var _ = gc.Suite(&MontonicitySuite{})

func (s *MontonicitySuite) TestName(c *gc.C) {
	c.Check(Montonicity(true).Name(), gc.Equals, "Montonicity")
}
func (s *MontonicitySuite) TestRestrictedDepths(c *gc.C) {
	var m Montonicity
	chart := &Chart{input: &TokenArray{"Y", "Z"}}

	Y, _ := chart.NextCell()
	Z, _ := chart.NextCell()

	c.Check(m.RestrictedDepths(Y.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)
	c.Check(m.RestrictedDepths(Z.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)

	// Depth=0 links don't constrain under montonicity.
	linkYZ := chart.UseAdjacency(Y.Right.OutboundAdjacency, 0)
	linkZY := chart.UseAdjacency(Z.Left.OutboundAdjacency, 0)

	c.Check(m.RestrictedDepths(Y.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)
	c.Check(m.RestrictedDepths(Z.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)

	// Depth=1 links do constrain, restricting depth=0 links.
	linkYZ.Depth = 1
	linkZY.Depth = 1

	c.Check(m.RestrictedDepths(Y.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_D0)
	c.Check(m.RestrictedDepths(Z.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_D0)
}
