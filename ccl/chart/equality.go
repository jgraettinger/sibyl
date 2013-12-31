package chart

import (
	"log"
)

func (adjacency *Adjacency) EqualityRestriction() (restrict DepthRestriction) {
	if adjacency.To == nil {
		return
	}
	forward := DirectionFromPosition(adjacency.Position)

	// Look for a link-path back to the head from tail. The depth of the
	// first link on that path constrains the allowable depth of the adjacency.
	if backReach := adjacency.TailSide().FurthestD0Path(); !forward.Less(
		adjacency.Head.Index, backReach) {
		log.Print("A d=0 path back-to-head exists")
		restrict = RESTRICT_D1
	} else if backReach = adjacency.TailSide().FurthestPath(); !forward.Less(
		adjacency.Head.Index, backReach) {
		log.Print("A d=1 path back-to-head exists")
		restrict = RESTRICT_D1
	}

	// Is there a cell Z such that Z is beyond the adjacency tail in the
	// forward direction, and where Z also links directly back to head?
	// If so, is there a forward link path from the tail to Z, such that
	// this adjacency would form the first link on a completed path
	// to Z from the head? If so, the Z => head back-link then constrains
	// the depth of this adjacency.
	if backLink := adjacency.Head.InboundLink; backLink != nil &&
		forward.Less(adjacency.Tail.Index, backLink.Head.Index) {

		if reach := forward.HeadSide(adjacency.Tail).FurthestPath(); !forward.Less(
			reach, backLink.Head) {

			log.Printf("Eq-restricting link is %v", backLink)
			// If tail has a link-path to Z, then the partial blocking
			// restruction must have constrained backLink to d=1.
			invariant(backLink.Depth == 1)
			restrict = RESTRICT_D0
		}
	}
	invariant(restrict != RESTRICT_ALL)
	return
}
