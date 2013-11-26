package chart

import (
	. "ccl/util"
	"invariant"
)

type Link struct {
	// Never nil
	From, To *Cell

	Position int
	Depth    uint

	// FurthestPath instances are shared among linear chains of links,
	// and always point to the last reachable cell along this path.
	FurthestPath *BoxedCellPointer
}

type LinkList []*Link

func (list *LinkList) Add(link *Link) {
	*list = append(*list, link)
	invariant.Equal(len(*list), Iabs(link.Position))
}

func (list LinkList) Last() *Link {
	if length := len(list); length != 0 {
		return list[length-1]
	}
	return nil
}

func NewLink(adjacency *Adjacency, depth uint) *Link {
	invariant.NotNil(adjacency.From)
	invariant.NotNil(adjacency.To)

	link := new(Link)
	link.From = adjacency.From
	link.To = adjacency.To
	link.Position = adjacency.Position
	link.Depth = depth

	return link
}

// Replaces the boxed furthest path along this ajacency path.
// Returns the current boxed path, which can be further updated
// with longer paths without altering this path.
func (link *Link) ForkFurthestPath() *BoxedCellPointer {
	forward := DirectionFromPosition(link.Position)

	oldPath, newPath := link.FurthestPath, new(BoxedCellPointer)
	*newPath = *oldPath

	next := link
	for next != nil {
		// Replace the boxed path instance along this link path
		next.FurthestPath = newPath
		next = forward.OutboundLinks(next.To).Last()
	}
	return oldPath
}
