package chart
/*
import (
	"invariant"
	"log"
)

func UpdateEquality(chart *Chart, usedAdjacency, newAdjacency *Adjacency, movedAdjacencies AdjacencyList) {

	forward := DirectionFromPosition(usedAdjacency.Position)
	backward := forward.Flip()

	cellFrom, cellTo := usedAdjacency.From, usedAdjacency.To

	// Cycle restriction Step 1: If a depth mismatch occurs in and
	// out of a cell, constrain other links from completing a cycle.

	// Look for a used backward link into cellFrom of mistmatched depth.
	if inbound := backward.UsedInbound(cellFrom); inbound != nil &&
		inbound.UsedDepth != usedAdjacency.UsedDepth {
		// We must prevent a cycle from forming through these two links.
		inequalityCycleRestriction(chart, inbound, usedAdjacency)
	}

	// Look for a used backward link out of cellTo of mismatched depth.
	for _, outbound := range *backward.Outbound(cellTo) {
		if outbound.Used && outbound.UsedDepth != usedAdjacency.UsedDepth {
			// We must prevent a cycle from formint through these links.
			inequalityCycleRestriction(chart, usedAdjacency, outbound)
		}
	}

	// Step 2: Project any previously set equality bound on cellFrom
	// through the link path opened by usedAdjacency. Because projection
	// is inclusive to cellFrom, this will also handle newAdjacency.
	if usedAdjacency.From.HasEqualityPathBound[forward.Side()] {
		projectEqualityPathBound(chart, usedAdjacency.From, forward,
			usedAdjacency.From.EqualityPathBound[forward.Side()])
	}

	// Step 3: Update equality blocking of moved adjacencies.
	for _, adjacency := range movedAdjacencies {
		if adjacency.From.HasEqualityPathBound[forward.Side()] &&
			adjacency.From.EqualityPathBound[forward.Side()] == adjacency.ToIndex(chart) {
			// This adjacency would complete an inequality-creating cycle.
			adjacency.EqualityBlocked = [2]bool{true, true}
		} else {
			// Otherwise Moving the adjacency clears current eqaulity restrictions.
			adjacency.EqualityBlocked = [2]bool{false, false}
		}
	}

	// Back-track over cells which have just had their link-path extended.
	for index := usedAdjacency.From.Index; ; index = forward.Decrement(index) {
		cell := chart.Cells[index]

		// Examine *backward* unused adjacencies, looking for potential long-
		// links that cell has a path back too, and which thus need to be
		// constrainted to the depth of cell's first link in the link-path.
		for adjacency := range backward.Inbound(cell) {
			if forward.Less(forward.PathEnd(cell), adjacency.From.Index) {
				// There isn't a path from cell to this adjacency's From.
				// Equality imposes no restriction on it's depth.
				continue
			}

			var requiredDepth uint
			if forward.Less(forward.PathEndD0(cell), adjacency.From.Index) {
				// Reachable through a path beginning with a d=1 link.
                invariant.IsFalse(forward.Less(forward.PathEndD1(cell),
                    adjacency.From.Index))
				requiredDepth = 1
			} else {
				// Reachable through a path beginning with a d=0 link.
				requiredDepth = 0
			}

			log.Printf("Equality is constraining back-%v to depth %v",
				adjacency, requiredDepth)

			if adjacency.Used {
				invariant.Equal(adjacency.UsedDepth, requiredDepth)
			} else if requiredDepth == 0 {
				adjacency.EqualityBlocked[1] = true
			} else if requiredDepth == 1 {
				adjacency.EqualityBlocked[0] = true
			}
		}

		// Examine *forward* unused adjacencies, looking for potential short-
		// links which have a used *backwards* long-link, which would be
		// reachable by link-path were the short-link to be used.
		for adjacency := range forward.Inbound(cell) {
			backInbound := backward.UsedInbound(adjacency.From)

			if backInbound == nil ||
				forward.Less(forward.PathEnd(cell), backInbound.From.Index) {
				// No used inbound, or it's not covered by the potential
				// link path if this adjacency were to be used.
				continue
			}

			requiredDepth := backInbound.UsedDepth

			log.Printf("Equality is constraining forward-%v to depth %v due to back-%v",
				adjacency, requiredDepth, backInbound)

			if adjacency.Used {
				invariant.Equal(adjacency.UsedDepth, requiredDepth)
			} else if requiredDepth == 0 {
				adjacency.EqualityBlocked[1] = true
			} else if requiredDepth == 1 {
				adjacency.EqualityBlocked[0] = true
			}
		}

		if index == forward.PathBegin(usedAdjacency.From) {
			// Reached the beginning of the link path.
			break
		}
	}
}

func projectEqualityPathBound(chart *Chart, root *Cell,
	forward Direction, pathBound int) {

	// Update cells in a link-path from root (inclusive).
	beginIndex := root.Index
	endIndex := forward.Increment(forward.PathEnd(root))

	for index := beginIndex; index != endIndex; index = forward.Increment(index) {

		cell := chart.Cells[index]

		log.Printf("Setting equality bound %v on %v", cell, pathBound)
		if cell.HasEqualityPathBound[forward.Side()] {
			log.Printf("Warning: equality bound %v already set",
				cell.EqualityPathBound[forward.Side()])
		}
		cell.EqualityPathBound[forward.Side()] = pathBound
		cell.HasEqualityPathBound[forward.Side()] = true

		adjacency := forward.Outbound(cell).Current()
		invariant.IsFalse(adjacency.Used)
		// Connectedness guarantees the adjacency can't span pathBound.
		invariant.IsFalse(forward.Less(pathBound, adjacency.To.Index))

		// Invalidate adjacencies to pathBound, as they'd complete a cycle.
		if adjacency.To.Index == pathBound {
			log.Printf("Blocking %v to prevent an equality-violating cycle", adjacency)
			adjacency.EqualityBlocked = [2]bool{true, true}
		}
	}
}

// Preconditions:
//  - adjacencyIn.Used && adjacencyOut.Used
//  - adjacencyIn.To == adjacencyOut.From
//  - adjacencyIn.Depth != adjacencyOut.Depth
//
// To satisfy the equality condition, we can't allow a cycle to form through
// these adjacencies. For purposes of discussion, one of these adjacencies
// will be over a longer span than the other: we must constrain paths
// rooted by the shorter adjaceny from forming a link back to the root of the
// longer adjacency, as this would complete a cycle in violation of equality.
//
// Precondition: cellEnd has a used outbound backward adjacency to cellBegin,
// and a used inbound forward adjacency from a node spaned by cellBegin &
// cellEnd, where these two adjacencies are of opposing depths.
func inequalityCycleRestriction(chart *Chart, linkIn, linkOut *Adjacency) {

	invariant.IsTrue(linkIn.Used)
	invariant.IsTrue(linkOut.Used)
	invariant.Equal(linkIn.To, linkOut.From)
	invariant.NotEqual(linkIn.UsedDepth, linkOut.UsedDepth)

	var cellBegin, cellEnd *Cell
	var longLink, shortLink *Adjacency
	var violationAtPathEnd bool

	if DirectionFromPosition(linkIn.Position).Less(linkIn.From.Index, linkOut.To.Index) {
		// linkIn is the long adjacency.
		longLink, shortLink = linkIn, linkOut
		cellBegin, cellEnd = longLink.From, longLink.To
		// shortLink is at the beginning of the path.
		violationAtPathEnd = false
	} else {
		// linkOut is the long adjacency
		longLink, shortLink = linkOut, linkIn
		cellBegin, cellEnd = longLink.To, longLink.From
		// shortLink is at the end of the path.
		violationAtPathEnd = true
	}

	forward := DirectionFromPosition(shortLink.Position)

	// Index of cell nearest to cellBegin which has a path to cellEnd.
	pathBound := forward.PathBegin(cellEnd)

	var projectFrom *Cell
	if violationAtPathEnd == false && longLink.UsedDepth == 1 {
		invariant.Equal(shortLink.UsedDepth, 0)

		// While we must prevent a cycle forming through shortLink, it's still
		// possible for a cycle to form via another d=1 link from cellBegin.
		projectFrom = shortLink.To
	} else {
		// Under no circumstances may a cycle form.
		projectFrom = cellBegin
	}

	projectEqualityPathBound(chart, projectFrom, forward, pathBound)
}
*/
