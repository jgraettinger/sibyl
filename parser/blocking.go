package parser

type PartialBlocking bool

func (_ PartialBlocking) Name() string {
	return "Partial Blocking"
}

// Presence of a backward link into Head from a node which would be
// spanned by this adjacency, restricts depth 0 of the adjacency.
func (_ PartialBlocking) RestrictedDepths(a *Adjacency) DepthRestriction {
	if il := a.HeadSide().InboundLink; il != nil {
		if a.Tail == nil {
			// Adjacency is to {begin} or {end}. It must then span il.Head.
			return RESTRICT_D0
		}
		if a.Position > 0 && il.Head.Index < a.Tail.Index {
			return RESTRICT_D0
		}
		if a.Position < 0 && il.Head.Index > a.Tail.Index {
			return RESTRICT_D0
		}
	}
	return RESTRICT_NONE
}

type FullBlocking struct {
	seenD1 []bool
	leftToRight, rightToLeft []int
}

func (_ FullBlocking) Name() string {
	return "Full Blocking"
}

// No adjacency may be used which spans beyond a cell with a
// path to Head, which also has an outbound d=1 link.
func (blocking FullBlocking) RestrictedDepths(a *Adjacency) DepthRestriction {
	if bound, ok := blocking[adjacency.HeadSide()]; ok {
		if a.Tail == nil {
			// Adjacency is to {begin} or {end}. It must span the bound.
			return RESTRICT_ALL
		}
		if a.Position > 0 && bound < a.Tail.Index {
			return RESTRICT_ALL
		}
		if a.Position < 0 && bound > a.Tail.Index {
			return RESTRICT_ALL
		}
	}
	return RESTRICT_NONE
}

func (b FullBlocking) projectBound(index, bound int) {
	invariant(bound != index)
	if bound < index {
		// Because of parser incrementalness, bounds cannot become
		// tighter then the first bound applied to a side.
		invariant(b.rightToLeft[index] == -1 || b.rightToLeft[index] > bound)
		b.rightToLeft[index] = bound
	} else {
		invariant(b.leftToRight[index] == -1 || b.leftToRight[index] < bound)
		b.leftToRight[index] = bound
	}
}

func (blocking FullBlocking) Observe(chart *Chart, link *Link) {
	for len(chart.Cells) != len(blocking.seenD1) {
		blocking.seenD1 = append(blocking.seenD1, false)
		blocking.leftToRight = append(blocking.leftToRight, -1)
		blocking.rightToLeft = append(blocking.rightToLeft, -1)
	}
	head := link.Head

	// Is this the first d=1 link from this cell? If so, project
	// it as a new blocking bound along all existing paths.
	if link.Depth == 1 {
		if !blocking.seenD1[head.Index] {
			if ll := lastLink(head.Left.OutboundLinks); ll != nil {
				for i := (*ll.BoxedFurthestPath).Index; i != head.Index; i++ {
					b.projectBound(i, head.Index)
				}
			}
			if ll := lastLink(head.Right.OutboundLinks); ll != nil {
				for i := (*ll.BoxedFurthestPath).Index; i != head.Index; i-- {
					b.projectBound(i, head.Index)
				}
			}
			blocking.seenD1[head.Index] = true
		} else {
			b.projectBound(link.Tail.Index, head.Index)
		}
	} else if head




	}

	if link.Position > 0 {
		if blocking.seenD1[link.Head.Index] {
			invariant(link.Tail.Index == len(chart.Cells))
			blocking.rightToLeft[link.Tail.Index] = link.Head.Index
		} else if link.Depth == 1 {

			blocking.rightToLeft[link.Tail.Index] = link.Head.Index

			for i := (*link.Head.Left.BoxedFurthestPath).Index;
				i != link.Head.Index; i++ {
					blocking.leftToRight[i] = link.Head.Index
				}
			}
		}


		}
		if !blocking.seenD1[link.Head.Index] && link.Depth == 1 {
			blocking.seenD1[link.Head.Index] = true
			updateLeftToRight(chart, (*link.Head.Left.BoxedFurthestPath).Index)


		}
		if blocking.seenD1[link.Head.Index] {


		}
	}

	if blocking.seenD1[link.Head.Index] {
		// As we've already seen a D1 link from this cell, just
		// project the blocking bound to the new tail.
		if link.Position < 0 {
			blocking.leftToRight[link.Tail.Index] = link.Head.Index
		}
	}
	if link.Depth == 1 {
		// A new d=1 link represents a blocking bound which must be
		// projected along furthest path extents in both directions.
	}
}

func (blocking FullBlocking) projectLeftToRight(chart *Chart, link *Link) {
	if blocking.seenD1[link.Head.Index] {

	}
}
