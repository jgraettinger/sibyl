package chart

import (
	"log"
)

type BlockingFlags byte

const (
	BLOCK_NONE = 0
	BLOCK_D0   = 1
	BLOCK_ALL  = 2
)

func updateBlocking(chart *Chart, newLink *Link, newAdjacency *Adjacency) {
	forward := DirectionFromPosition(newLink.Position)
	backward := forward.Flip()

	cell := newLink.From

	// 3.2.1 Full blocking condition

	// Blocking requires that no link may span a d=1 link from a
	// node which has also has a path to the adjacency base.

	// Retrieve the most recent link from cell, in the other direction.
	// Due to montonicity if any links are d=1, than so is the last one.
	backLink := backward.OutboundLinks(cell).Last()

	// The head of a d=1 link is a blocking bound which must be
	// projected along furthest path extents in both directions.
	if newLink.Depth == 1 {
		projectBlocking(chart, newLink, cell.Index)
		if backLink != nil {
			projectBlocking(chart, backLink, cell.Index)
		}
	} else if backLink != nil && backLink.Depth == 1 {
		// Cell is a blocking bound to be projected along this link.
		projectBlocking(chart, newLink, cell.Index)
	} else if backward.HasFullyBlockedAfter(cell) {
		// Propogate an earlier constraint along the path extended by newLink.
		projectBlocking(chart, newLink,
			backward.FullyBlockedAfter(cell))
	}

	// If newAdjacency spans FullyBlockedAfter of cell, there must
	// must be a d=1 link which newAdjacency spans.
	if forward.HasFullyBlockedAfter(cell) {
		if bound := forward.FullyBlockedAfter(cell); newAdjacency.To == nil ||
			forward.Less(bound, newAdjacency.To.Index) {

			log.Printf("Immediately fully blocking new %v (blocked after %v)",
				newAdjacency, bound)
			newAdjacency.Blocking |= BLOCK_ALL
		}
	}

	// 3.2.1 Partial (d=0) blocking condition

	// Step 1: If there is a backward inbound link into
	// cell which is not from newAdjacency.To, than d=0 is blocked.
	if link := backward.InboundLink(cell); link != nil &&
		link.From != newAdjacency.To {
		log.Printf("Blocking d=0 of new adjacency because of link %v", link)
		newAdjacency.Blocking |= BLOCK_D0
	}

	// Step 2: Inversely, the creation of newLink blocks d=0 of a
	// current backward adjacency from newLink.To spanning beyond newLink.From.
	backAdjacency := backward.OutboundAdjacency(newLink.To)
	log.Printf("backAdjacency is %v", backAdjacency)
	if backAdjacency.To != newLink.From {
		log.Printf("Blocking d=0 of %v due to %v", backAdjacency, newLink)
		backAdjacency.Blocking |= BLOCK_D0
	}
}

// Projects blockIndex to all cells reachable on a path from link.
// Updates blocking of any adacencies affected by the new bound.
func projectBlocking(chart *Chart, link *Link, blockIndex int) {
	log.Printf("projecting blocking bound %v along %v", blockIndex, link)
	forward := DirectionFromPosition(link.Position)
	backward := forward.Flip()

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

		// Does this projection invalidate the current adjacency from cell?
		adjacency := backward.OutboundAdjacency(cell)
		if adjacency.To == nil ||
			backward.Less(blockIndex, adjacency.To.Index) {
			log.Printf("Fully blocking adjacency %v (%v)",
				adjacency, blockIndex)
			adjacency.Blocking |= BLOCK_ALL
		}
	}
}
