package parser

type Link struct {
	// Never nil
	Head, Tail *Cell

	Position int
	Depth    int

	// BoxedFurthestPath *Cell instances are allocated and shared
	// among linear chains of links, and always point to the last
	// reachable cell along this path. Sharing a *Cell allows
	// the furthest cell to be updated for multiple links in O(1).
	BoxedFurthestPath **Cell
}

func NewLink(adjacency *Adjacency, depth int) *Link {
	invariant(adjacency.Head != nil)
	invariant(adjacency.Tail != nil)
	invariant(depth == 0 || depth == 1)
	return &Link{
		Head: adjacency.Head,
		Tail: adjacency.Tail,
		Position: adjacency.Position,
		Depth: depth}
}

func (link *Link) appendTo(list *[]*Link) {
	*list = append(*list, link)
}

func lastLink(list []*Link) *Link {
	if l := len(list); l != 0 {
		return list[l-1]
	}
	return nil
}

func updateBoxedPathLeftToRight(link *Link) {
	// Update furthest-paths to reflect the new link
	var boxedPath **Cell
	if ll := lastLink(link.Head.Right.OutboundLinks); ll == nil {
		// As this is head's first outbound link in this direction,
		// using this adjacency doesn't create a new link-path, but
		// may extend an existing path from antecedent cells.
		if pathIn := link.Head.Left.InboundLink; pathIn != nil {
			// We'll update the existing link-path.
			boxedPath = pathIn.BoxedFurthestPath
		} else {
			// There is no existing path. Create a new one.
			boxedPath = new(*Cell)
		}
	} else {
		// Adding a second outbound link creates a new link-path. Because
		// parsing is left-to-right, we expect that additional links may
		// be added to this path. However, by minimality, the current
		// outbound link path may not be extended any further and can
		// be updated with a copied, 'frozen' boxed path. The existing
		// boxed path is used in the new link, so that antecedents continue
		// to be updated.
		boxedPath = ll.BoxedFurthestPath
		frozenPath := new(*Cell)
		*frozenPath = *boxedPath

		for ll != nil {
			// Replace the boxed path instance along this link path.
			ll.BoxedFurthestPath = frozenPath
			ll = lastLink(ll.Tail.Right.OutboundLinks)
		}
	}
	// This update is now visible from all previous links on the path.
	*boxedPath = link.Tail
	link.BoxedFurthestPath = boxedPath
}

func updateBoxedPathRightToLeft(link *Link) {
	// By the incremental nature of the parser, any successive links
	// of a right-to-left link must already exist. The current
	// furthest path can't be extended, and can just be copied.
	if ll := lastLink(link.Tail.Left.OutboundLinks); ll != nil {
		link.BoxedFurthestPath = ll.BoxedFurthestPath
	} else {
		link.BoxedFurthestPath = new(*Cell)
		*link.BoxedFurthestPath = link.Tail
	}
}
func (l *Link) HeadSide() *CellSide {
	if l.Position < 0 {
		return &l.Head.Left
	}
	return &l.Head.Right
}
func (l *Adjacency) TailSide() *CellSide {
	if l.Position < 0 {
		return &l.Tail.Right
	}
	return &l.Tail.Left
}
