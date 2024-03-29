package chart

import (
	"fmt"

	"github.com/dademurphy/sibyl/invariant"
)

type DepthRestriction int

const (
	RESTRICT_NONE = 0
	RESTRICT_D0   = 1
	RESTRICT_D1   = 2
	RESTRICT_ALL  = 3
)

type Adjacency struct {
	// Head is always non-nil, but tail may be nil to represent
	// an adjacency beyond the sentence extents.
	Head, Tail *Cell

	// The argument attachment position this adjacency reflects
	Position int

	SpansPunctuation bool

	// Denotes this adjacency is 'covered' by a link spanning
	// Head & Tail. Covered adjacencies cannot be used (see 3.2.2).
	CoveredByLink bool
}

func (adjacency *Adjacency) HeadSide() *Side {
	if adjacency.Position < 0 {
		return &adjacency.Head.Left
	} else {
		return &adjacency.Head.Right
	}
}
func (adjacency *Adjacency) TailSide() *Side {
	if adjacency.Position < 0 {
		return &adjacency.Tail.Right
	} else {
		return &adjacency.Tail.Left
	}
}

func (adjacency *Adjacency) RestrictedDepths(chart *Chart) (
	restrict DepthRestriction) {
	restrict |= adjacency.PunctuationRestriction()
	restrict |= adjacency.CoveredLinkRestriction()
	restrict |= adjacency.MontonicityRestriction()
	restrict |= adjacency.BlockingRestriction()
	restrict |= adjacency.EqualityRestriction()
	restrict |= adjacency.ResolutionRestriction(chart)
	return
}

func (adjacency *Adjacency) IsMoveable() bool {
	var restrict DepthRestriction
	restrict |= adjacency.PunctuationRestriction()
	restrict |= adjacency.CoveredLinkRestriction()
	restrict |= adjacency.MontonicityRestriction()
	restrict |= adjacency.BlockingRestriction()
	return restrict != RESTRICT_ALL
}

func (adjacency *Adjacency) PunctuationRestriction() DepthRestriction {
	if adjacency.SpansPunctuation {
		return RESTRICT_ALL
	}
	return RESTRICT_NONE
}
func (adjacency *Adjacency) CoveredLinkRestriction() DepthRestriction {
	if adjacency.CoveredByLink {
		return RESTRICT_ALL
	}
	return RESTRICT_NONE
}

type AdjacencySet map[*Adjacency]bool

func (set AdjacencySet) Add(adjacency *Adjacency) {
	set[adjacency] = true
}

func (set AdjacencySet) Remove(adjacency *Adjacency) {
	_, present := set[adjacency]
	invariant.IsTrue(present)
	delete(set, adjacency)
}

func (adjacency *Adjacency) String() string {

	properties := ""

	if r := adjacency.PunctuationRestriction(); r != RESTRICT_NONE {
		properties += fmt.Sprintf(", (punc %v)", r)
	}
	if r := adjacency.CoveredLinkRestriction(); r != RESTRICT_NONE {
		properties += fmt.Sprintf(", (cov %v)", r)
	}
	if r := adjacency.MontonicityRestriction(); r != RESTRICT_NONE {
		properties += fmt.Sprintf(", (mont %v)", r)
	}
	if r := adjacency.BlockingRestriction(); r != RESTRICT_NONE {
		properties += fmt.Sprintf(", (block %v)", r)
	}
	if r := adjacency.EqualityRestriction(); r != RESTRICT_NONE {
		properties += fmt.Sprintf(", (eq %v)", r)
	}
	return fmt.Sprintf("Adjacency<%v:%d => %v%v>", adjacency.Head,
		adjacency.Position, adjacency.Tail, properties)
}

func (restrict DepthRestriction) String() string {
	if restrict == RESTRICT_NONE {
		return "NONE"
	} else if restrict == RESTRICT_D0 {
		return "D0"
	} else if restrict == RESTRICT_D1 {
		return "D1"
	}
	return "D0+1"
}

func (restrict DepthRestriction) Allows(depth uint) bool {
	invariant.IsTrue(depth == 0 || depth == 1)
	if depth == 0 {
		return restrict&RESTRICT_D0 == 0
	}
	return restrict&RESTRICT_D1 == 0
}
