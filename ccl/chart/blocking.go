package chart

import (
	"log"
)

// Presence of a backward link into Head from a node which would be fully
// covered by this adjacency, restricts depth 0 of the adjacency.
func PartialBlockingConstraint(adjacency *Adjacency) DepthRestriction {
	if backLink := adjacency.HeadSide().InboundLink; backLink != nil &&
		adjacency.Covers(backLink.Head) {
		return RESTRICT_D0
	}
	return RESTRICT_NONE
}

func (adjacency *Adjacency) BlockingRestriction() (restrict DepthRestriction) {
	forward := DirectionFromPosition(adjacency.Position)

	// Presence of a backward link into Head which isn't from To, blocks d=0.
	if backLink := adjacency.HeadSide().InboundLink; backLink != nil &&
		backLink.Head != adjacency.Tail {
		log.Printf("Blocking d=0 because of %v", backLink)
		restrict |= RESTRICT_D0
	}

	// No adjacency may be used beyond a cell with a path to Head,
	// which also has an outbound d=1 link.
	if bound := adjacency.HeadSide().FullyBlockedAfter; bound != nil &&
		forward.Less(*bound, adjacency.Tail.Index) {
		log.Printf("Block-bound of %v fully blocks", bound)
		restrict |= RESTRICT_ALL
	}
	return
}

func updateBlocking(chart *Chart, link *Link) {
	forward := DirectionFromPosition(link.Position)
	backward := forward.Flip()

	// Blocking requires that no link may span beyond a node having an outbound
	// d=1 link (in any direction) and also having a path to the adjacency head.

	// Retrieve the most recent link from the head, but in the other direction.
	// Due to montonicity if any links are d=1, than the last one is as well.
	backLink := backward.HeadSide(link.Head).lastOutboundLink

	if link.Depth == 1 {
		// A new d=1 link represents a blocking bound which must be
		// projected along furthest path extents in both directions.
		projectBlocking(chart, link.BoxedFurthestPath.Index, link.Head.Index)
		if backLink != nil && backLink.Depth == 0 {
			// If backlink d=1, then the bound at this head must have
			// already has already been projected in that direction.
			projectBlocking(chart, backLink.BoxedFurthestPath.Index, link.Head.Index)
		}
	} else if backLink != nil && backLink.Depth == 1 {
		// Link head is a blocking bound, due to d=1 link in other direction.
		projectBlocking(chart, link.BoxedFurthestPath.Index, link.Head.Index)
	} else if bound := backward.HeadSide(link.Head).FullyBlockedAfter; bound != nil {
		// Propogate an earlier constraint along the path extended by link.
		projectBlocking(chart, link.BoxedFurthestPath.Index, bound)
	}
}

func projectBlocking(chart *Chart, link *Link, blockIndex int) {
	forward := DirectionFromPosition(link.Position)
	backward := forward.Flip()

	log.Printf("Projecting blocking bound %v along %v", blockIndex, link)

	begin := link.To.Index
	end := forward.Increment((*link.FurthestPath).Index)

	for index := begin; index != end; index = forward.Increment(index) {
		cell := chart.Cells[index]

		if backward.HasFullyBlockedAfter(cell) &&
			!backward.Less(blockIndex, backward.FullyBlockedAfter(cell)) {
			log.Printf("Breaking out of blocking projection at %v", cell)
			break
		}

		log.Printf("Projecting blocking bound %v to %v", blockIndex, cell)
		backward.SetFullyBlockedAfter(cell, blockIndex)
	}
}
