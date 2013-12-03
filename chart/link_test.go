package chart

import (
	gc "launchpad.net/gocheck"
)

type LinkSuite struct{}

var _ = gc.Suite(&LinkSuite{})

func (s *LinkSuite) TestBoxedPathUpdateLeftToRight(c *gc.C) {
	T, W, X, Y, Z := &Cell{Index: 0}, &Cell{Index: 1}, &Cell{Index: 2},
		&Cell{Index: 3}, &Cell{Index: 4}

	createAndUpdate := func(head, tail *Cell) *Link {
		link := &Link{
			Head:     head,
			Tail:     tail,
			Position: len(head.Right.OutboundLinks) + 1}

		updateBoxedPathLeftToRight(link)
		link.appendTo(&head.Right.OutboundLinks)
		tail.Left.InboundLink = link
		return link
	}
	linkTW := createAndUpdate(T, W)
	c.Check(*linkTW.BoxedFurthestPath, gc.Equals, W)

	linkWX := createAndUpdate(W, X)
	c.Check(*linkTW.BoxedFurthestPath, gc.Equals, X)
	c.Check(*linkWX.BoxedFurthestPath, gc.Equals, X)

	linkXY := createAndUpdate(X, Y)
	c.Check(*linkTW.BoxedFurthestPath, gc.Equals, Y)
	c.Check(*linkWX.BoxedFurthestPath, gc.Equals, Y)
	c.Check(*linkXY.BoxedFurthestPath, gc.Equals, Y)

	linkWZ := createAndUpdate(W, Z)
	c.Check(*linkTW.BoxedFurthestPath, gc.Equals, Z)
	c.Check(*linkWX.BoxedFurthestPath, gc.Equals, Y)
	c.Check(*linkXY.BoxedFurthestPath, gc.Equals, Y)
	c.Check(*linkWZ.BoxedFurthestPath, gc.Equals, Z)
}
func (s *LinkSuite) TestBoxedPathUpdateRightToLeft(c *gc.C) {
	W, X, Y, Z := &Cell{Index: 0}, &Cell{Index: 1}, &Cell{Index: 2},
		&Cell{Index: 3}

	createAndUpdate := func(head, tail *Cell) *Link {
		link := &Link{
			Head:     head,
			Tail:     tail,
			Position: len(head.Right.OutboundLinks) - 1}

		updateBoxedPathRightToLeft(link)
		link.appendTo(&head.Left.OutboundLinks)
		tail.Right.InboundLink = link
		return link
	}
	linkYX := createAndUpdate(Y, X)
	c.Check(*linkYX.BoxedFurthestPath, gc.Equals, X)

	linkYW := createAndUpdate(Y, W)
	c.Check(*linkYX.BoxedFurthestPath, gc.Equals, X)
	c.Check(*linkYW.BoxedFurthestPath, gc.Equals, W)

	linkZY := createAndUpdate(Z, Y)
	c.Check(*linkYX.BoxedFurthestPath, gc.Equals, X)
	c.Check(*linkYW.BoxedFurthestPath, gc.Equals, W)
	c.Check(*linkZY.BoxedFurthestPath, gc.Equals, W)
}
