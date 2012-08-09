package chart

import (
	"log"
)

func (adjacency *Adjacency) BlockingRestriction() (restrict DepthRestriction) {
	forward := DirectionFromPosition(adjacency.Position)
	backward := forward.Flip()

	// Presence of a backward link into From which isn't from To, blocks d=0.
	backLink := backward.InboundLink(adjacency.From)
	if backLink != nil && backLink.From != adjacency.To {
		log.Printf("Blocking d=0 because of %v", backLink)
		restrict |= RESTRICT_D0
	}

	// Spanning a d=1 link of a cell having a path back to From, blocks all.
	if forward.HasFullyBlockedAfter(adjacency.From) {
		bound := forward.FullyBlockedAfter(adjacency.From)
		if forward.Less(bound, adjacency.To.Index) {
			log.Printf("Block-bound of %v fully blocks", bound)
			restrict |= RESTRICT_ALL
		}
	}
	return
}

func updateBlocking(chart *Chart, link *Link) {
	forward := DirectionFromPosition(link.Position)
	backward := forward.Flip()

	// Blocking requires that no link may span a d=1 link from a
	// node which has also has a path to the adjacency base.

	// Retrieve the most recent link from the head, but in the other direction.
	// Due to montonicity if any links are d=1, than the last one is as well.
	backLink := backward.OutboundLinks(link.From).Last()

	// A new d=1 link represents a blocking bound which must be
	// projected along furthest path extents in both directions.
	if link.Depth == 1 {
		projectBlocking(chart, link, link.From.Index)
		if backLink != nil {
			projectBlocking(chart, backLink, link.From.Index)
		}
	} else if backLink != nil && backLink.Depth == 1 {
		// Link head is a blocking bound, due to d=1 link in other direction.
		projectBlocking(chart, link, link.From.Index)
	} else if backward.HasFullyBlockedAfter(link.From) {
		// Propogate an earlier constraint along the path extended by link.
		projectBlocking(chart, link, backward.FullyBlockedAfter(link.From))
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
