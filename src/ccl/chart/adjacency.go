package chart

import (
	"fmt"
	"invariant"
)

type Adjacency struct {
	From *Cell
	To   *Cell // inclusive; may be nil

	// The argument attachment position this adjacency reflects
	Position int

	SpansPunctuation bool

	// Denotes this adjacency is 'covered' by a link spanning
	// From & To. Covered adjacencies cannot be used (see 3.2.2).
	CoveredByLink bool

	MontonicityRestricted bool

	Blocking BlockingFlags
}

func (adjacency *Adjacency) IsUsable() bool {
	if !adjacency.SpansPunctuation &&
		!adjacency.CoveredByLink &&
		adjacency.Blocking&BLOCK_ALL == 0 {
		return true
	}
	return false
}

func (adjacency *Adjacency) IsMoveable() bool {
	if !adjacency.SpansPunctuation &&
		!adjacency.CoveredByLink &&
		adjacency.Blocking&BLOCK_ALL == 0 {
		return true
	}
	return false
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
	if adjacency.SpansPunctuation {
		properties += ", (punc)"
	}
	if adjacency.CoveredByLink {
		properties += ", (covered)"
	}
	if adjacency.MontonicityRestricted {
		properties += ", (montonicity)"
	}
	if adjacency.Blocking&BLOCK_D0 != 0 {
		properties += ", (partial blocking)"
	}
	if adjacency.Blocking&BLOCK_ALL != 0 {
		properties += ", (full blocking)"
	}

	return fmt.Sprintf("Adjacency<%v:%d => %v%v>", adjacency.From,
		adjacency.Position, adjacency.To, properties)
}
