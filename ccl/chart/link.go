package chart


type Link struct {
	// Never nil
	Head, Tail *Cell

	Position int
	Depth    uint

	// BoxedFurthestPath *Cell instances are allocated and shared
	// among linear chains of links, and always point to the last
	// reachable cell along this path. Sharing a *Cell allows
	// the furthest cell to be updated for multiple links in O(1).
	BoxedFurthestPath **Cell
}

func NewLink(adjacency *Adjacency, depth uint) *Link {
	invariant(adjacency.Head != nil)
	invariant(adjacency.Tail != nil)

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
func (link *Link) forkBoxedFurthestPath() **Cell {
	forward := DirectionFromPosition(link.Position)

	oldPath, newPath := link.BoxedFurthestPath, new(*Cell)
	*newPath = *oldPath

	next := link
	for next != nil {
		// Replace the boxed path instance along this link path
		next.BoxedFurthestPath = newPath
		next = forward.HeadSide(next.To).lastOutboundLink
	}
	return oldPath
}
