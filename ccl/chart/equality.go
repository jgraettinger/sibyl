package chart

import (
	"log"

	"github.com/dademurphy/sibyl/invariant"
)

func (adjacency *Adjacency) EqualityRestriction() (restrict DepthRestriction) {
	forward := DirectionFromPosition(adjacency.Position)
	backward := forward.Flip()

	if adjacency.To == nil {
		return
	}

	// Is there a link-path back to head from tail? The depth of the first
	// link on that path constrains the allowable depth of the adjacency.
	if adjacency.To.D0LinkPathReaches(adjacency.From) {
		log.Print("A d=0 path back-to-head exists")
		restrict = RESTRICT_D1
	} else if adjacency.To.D1LinkPathReaches(adjacency.From) {
		log.Print("A d=1 path back-to-head exists")
		restrict = RESTRICT_D0
	}

	// Is there a cell Z, such that Z directly links to head and this
	// adjacency would be the first link on a completed path reaching Z?
	// If so, the Z => head link constrains the depth of this adjacency.
	if backLink := backward.InboundLink(adjacency.From); backLink != nil &&
		forward.Less(adjacency.To.Index, backLink.From.Index) &&
		adjacency.To.LinkPathReaches(backLink.From) {

		log.Printf("Eq-restricting link is %v", backLink)

		// If adjacency.To has a link-path to backLink.From, than the partial
		// blocking restriction constrains backLink to d=1.
		invariant.IsTrue(backLink.Depth == 1)
		restrict = RESTRICT_D0
	}
	invariant.IsFalse(restrict == RESTRICT_ALL)
	return
}
