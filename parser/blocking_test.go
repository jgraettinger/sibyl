package parser

import (
	gc "launchpad.net/gocheck"
)

type PartialBlockingSuite struct{}

var _ = gc.Suite(&PartialBlockingSuite{})

func (s *PartialBlockingSuite) TestName(c *gc.C) {
	c.Check(PartialBlocking(true).Name(), gc.Equals, "Partial Blocking")
}
func (s *PartialBlockingSuite) TestRestrictedDepthsBeyondUtterance(c *gc.C) {
	var m PartialBlocking
	chart := &Chart{input: &TokenArray{"Y", "Z"}}

	Y, _ := chart.NextCell()
	Z, _ := chart.NextCell()

	c.Check(m.RestrictedDepths(Y.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)
	c.Check(m.RestrictedDepths(Z.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)

	chart.UseAdjacency(Y.Right.OutboundAdjacency, 0)
	chart.UseAdjacency(Z.Left.OutboundAdjacency, 0)

	// Adjacencies to {begin} and {end} are partially blocked.
	c.Check(m.RestrictedDepths(Y.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_D0)
	c.Check(m.RestrictedDepths(Z.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_D0)
}
func (s *PartialBlockingSuite) TestRestrictedDepths(c *gc.C) {
	var m PartialBlocking
	chart := &Chart{input: &TokenArray{"W", "X", "Y", "Z"}}

	chart.NextCell()
	X, _ := chart.NextCell()
	Y, _ := chart.NextCell()

	// There are no inbound links; adjacencies between X & Y aren't blocked.
	c.Check(m.RestrictedDepths(X.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)
	c.Check(m.RestrictedDepths(Y.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_NONE)

	chart.UseAdjacency(X.Right.OutboundAdjacency, 0)
	chart.UseAdjacency(Y.Left.OutboundAdjacency, 0)
	chart.NextCell()

	// New adjacencies to W & Z are, however.
	c.Check(m.RestrictedDepths(X.Right.OutboundAdjacency),
		gc.Equals, RESTRICT_D0)
	c.Check(m.RestrictedDepths(Y.Left.OutboundAdjacency),
		gc.Equals, RESTRICT_D0)
}
