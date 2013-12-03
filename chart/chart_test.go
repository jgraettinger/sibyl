package chart

import (
	"io"
	"testing"

	gc "launchpad.net/gocheck"
)

// Package-level test setup.
func TestPackage(t *testing.T) { gc.TestingT(t) }

type ChartSuite struct{}

var _ = gc.Suite(&ChartSuite{})

func (s *ChartSuite) TestNextCell(c *gc.C) {
	chart := &Chart{input: &TokenArray{"X", "Y", "Z"}}
	c.Check(chart.Cells, gc.HasLen, 0)

	X, err := chart.NextCell()
	// "X" Cell was created properly.
	c.Check(err, gc.IsNil)
	c.Check(X.Index, gc.Equals, 0)
	c.Check(X.Token, gc.Equals, Token("X"))
	c.Check(chart.Cells, gc.DeepEquals, []*Cell{X})

	// Forward and backward adjacencies were created from X.
	c.Check(*X.Left.OutboundAdjacency, gc.Equals,
		Adjacency{Head: X, Tail: nil, Position: -1})
	c.Check(*X.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: X, Tail: nil, Position: 1})

	// X has an adjacency to {end} only.
	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{X.Right.OutboundAdjacency})
	c.Check(chart.leftToRightAdjacencies, gc.HasLen, 0)

	Y, err := chart.NextCell()
	c.Check(err, gc.IsNil)
	c.Check(Y.Index, gc.Equals, 1)
	c.Check(Y.Token, gc.Equals, Token("Y"))
	c.Check(chart.Cells, gc.DeepEquals, []*Cell{X, Y})

	// Forward and backward adjacencies were created from Y.
	c.Check(*Y.Left.OutboundAdjacency, gc.Equals,
		Adjacency{Head: Y, Tail: X, Position: -1})
	c.Check(*Y.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: Y, Tail: nil, Position: 1})

	// X's adjacency to {end} has been moved to Y, which is now adjacent to {end}.
	c.Check(*X.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: X, Tail: Y, Position: 1})
	c.Check(chart.leftToRightAdjacencies, gc.DeepEquals,
		[]*Adjacency{X.Right.OutboundAdjacency})
	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{Y.Right.OutboundAdjacency})

	Z, _ := chart.NextCell()
	c.Check(chart.Cells, gc.DeepEquals, []*Cell{X, Y, Z})

	// Y's adjacency to {end} has been moved to Z, which is now adjacent to {end}.
	// (X's unused, active adjacency to Y is no longer tracked).
	c.Check(*Y.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: Y, Tail: Z, Position: 1})
	c.Check(chart.leftToRightAdjacencies, gc.DeepEquals,
		[]*Adjacency{Y.Right.OutboundAdjacency})
	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{Z.Right.OutboundAdjacency})

	// End-of-sequence handling.
	shouldBeNil, err := chart.NextCell()
	c.Check(err, gc.Equals, io.EOF)
	c.Check(shouldBeNil, gc.IsNil)
}

func (s *ChartSuite) TestLinkLeftToRight(c *gc.C) {
	checkLink := func(l *Link, h, t, f *Cell, p, d int) {
		c.Check(l.Head, gc.Equals, h)
		c.Check(l.Tail, gc.Equals, t)
		c.Check(*l.BoxedFurthestPath, gc.Equals, f)
		c.Check(l.Position, gc.Equals, p)
		c.Check(l.Depth, gc.Equals, d)
	}
	chart := &Chart{input: &TokenArray{"W", "X", "Y", "Z"}}

	W, _ := chart.NextCell()
	X, _ := chart.NextCell()

	// Precondition: Only X is adjacent to {end}. W's adjacency to X is active.
	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{X.Right.OutboundAdjacency})
	c.Check(chart.leftToRightAdjacencies, gc.DeepEquals,
		[]*Adjacency{W.Right.OutboundAdjacency})

	// A link W => X is created and tracked.
	linkWX := chart.UseAdjacency(W.Right.OutboundAdjacency, 1)
	checkLink(linkWX, W, X, X, 1, 1)
	c.Check(W.Right.OutboundLinks, gc.DeepEquals, []*Link{linkWX})
	c.Check(X.Left.InboundLink, gc.Equals, linkWX)

	// A new adjacency from W to {end} was created.
	c.Check(*W.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: W, Tail: nil, Position: 2})
	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{X.Right.OutboundAdjacency, W.Right.OutboundAdjacency})
	c.Check(chart.leftToRightAdjacencies, gc.HasLen, 0)

	Y, _ := chart.NextCell()

	// Precondition: W & X are both adjacent to Y.
	c.Check(chart.leftToRightAdjacencies, gc.DeepEquals,
		[]*Adjacency{X.Right.OutboundAdjacency, W.Right.OutboundAdjacency})

	// A link X => Y is created and tracked.
	linkXY := chart.UseAdjacency(X.Right.OutboundAdjacency, 0)
	checkLink(linkXY, X, Y, Y, 1, 0)
	c.Check(X.Right.OutboundLinks, gc.DeepEquals, []*Link{linkXY})
	c.Check(Y.Left.InboundLink, gc.Equals, linkXY)

	// A new adjacency from X to {end} was created. W's was moved to {end}.
	c.Check(*X.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: X, Tail: nil, Position: 2})
	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{Y.Right.OutboundAdjacency, W.Right.OutboundAdjacency,
			X.Right.OutboundAdjacency})
	c.Check(chart.leftToRightAdjacencies, gc.HasLen, 0)

	Z, _ := chart.NextCell()

	// A link X => Z is created and tracked.
	linkXZ := chart.UseAdjacency(X.Right.OutboundAdjacency, 0)
	checkLink(linkXZ, X, Z, Z, 2, 0)
	c.Check(X.Right.OutboundLinks, gc.DeepEquals, []*Link{linkXY, linkXZ})
	c.Check(Z.Left.InboundLink, gc.Equals, linkXZ)

	// A new adjacency from X to {end} was created. W's was moved to {end}.
	// As Y's former adjacency has been covered, it no longer has one.
	c.Check(*X.Right.OutboundAdjacency, gc.Equals,
		Adjacency{Head: X, Tail: nil, Position: 3})
	c.Check(Y.Right.OutboundAdjacency, gc.IsNil)
	c.Check(chart.leftToRightAdjacencies, gc.HasLen, 0)

	c.Check(chart.endAdjacencies, gc.DeepEquals,
		[]*Adjacency{Z.Right.OutboundAdjacency, W.Right.OutboundAdjacency,
			X.Right.OutboundAdjacency})

	// Verify that furthest-path updates have propagated.
	c.Check(*linkWX.BoxedFurthestPath, gc.Equals, Z)
	c.Check(*linkXY.BoxedFurthestPath, gc.Equals, Y)
	c.Check(*linkXZ.BoxedFurthestPath, gc.Equals, Z)
}

func (s *ChartSuite) TestLinkRightToLeft(c *gc.C) {
	checkLink := func(l *Link, h, t, f *Cell, p, d int) {
		c.Check(l.Head, gc.Equals, h)
		c.Check(l.Tail, gc.Equals, t)
		c.Check(*l.BoxedFurthestPath, gc.Equals, f)
		c.Check(l.Position, gc.Equals, p)
		c.Check(l.Depth, gc.Equals, d)
	}
	chart := &Chart{input: &TokenArray{"W", "X", "Y", "Z"}}

	W, _ := chart.NextCell()
	X, _ := chart.NextCell()
	Y, _ := chart.NextCell()

	// A link Y => X is created and tracked.
	linkYX := chart.UseAdjacency(Y.Left.OutboundAdjacency, 0)
	checkLink(linkYX, Y, X, X, -1, 0)
	c.Check(Y.Left.OutboundLinks, gc.DeepEquals, []*Link{linkYX})
	c.Check(X.Right.InboundLink, gc.DeepEquals, linkYX)

	// A new adjacency from Y => W is created.
	c.Check(*Y.Left.OutboundAdjacency, gc.Equals,
		Adjacency{Head: Y, Tail: W, Position: -2})

	// A link Y => W is created and tracked.
	linkYW := chart.UseAdjacency(Y.Left.OutboundAdjacency, 1)
	checkLink(linkYW, Y, W, W, -2, 1)
	c.Check(Y.Left.OutboundLinks, gc.DeepEquals, []*Link{linkYX, linkYW})
	c.Check(W.Right.InboundLink, gc.DeepEquals, linkYW)

	// Y has an adjacency to {begin}.
	c.Check(*Y.Left.OutboundAdjacency, gc.Equals,
		Adjacency{Head: Y, Tail: nil, Position: -3})

	Z, _ := chart.NextCell()

	// A link Z => Y is created and tracked, with a path to W.
	linkZY := chart.UseAdjacency(Z.Left.OutboundAdjacency, 0)
	checkLink(linkZY, Z, Y, W, -1, 0)
	c.Check(Z.Left.OutboundLinks, gc.DeepEquals, []*Link{linkZY})
	c.Check(Y.Right.InboundLink, gc.DeepEquals, linkZY)

	// Z has an adjacency to {begin}.
	c.Check(*Z.Left.OutboundAdjacency, gc.Equals,
		Adjacency{Head: Z, Tail: nil, Position: -2})
}
