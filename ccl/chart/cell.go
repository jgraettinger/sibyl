package chart

import (
	"fmt"
)

type Side struct {
	// Marks the index of the closest node which has an outbound d=1
	// link, and also has a link-path back to this node. Nil if there
	// is none. This is used to enforce blocking: specifically, no
	// adjacency from this side may span beyond BlockedAfter.
	// (See section 3.2.1, condition 3)
	FullyBlockedAfter *int

	// The singular unused outbound adjacency, and potentially multiple
	// unused inbound adjacencies.
	OutboundAdjacency  *Adjacency
	InboundAdjacencies []*Adjacency

	// Adjacencies which have been used to form links. Multiple outbound
	// adjacencies may have been used, but only one inbound one can be.
	OutboundLinks []*Link
	InboundLink   *Link

	// TODO fix
	// Last-added d=0 & d=1 links, reflecting the furthest link of each depth
	// from this cell. Due to montonicity, the d=1 link will either be nil,
	// or will be a further adjancency than the d=0 link.
	lastOutboundD0Link *Link
	lastOutboundLink *Link
}

func (s *Side) FurthestD0Path() int {
	if s.lastOutboundD0Link != nil {
		return s.lastOutboundD0Link.FurthestPath.Index
	} else {
		return s.OutboundAdjacency.From.Index
	}
}
func (s *Side) FurthestPath() int {
	if  s.lastOutboundLink != nil {
		return s.lastOutboundLink.FurthestPath.Index
	} else {
		return s.OutboundAdjacency.From.Index
	}
}

type Cell struct {
	Index int
	Token Token

	Left, Right Side
}

func (cell *Cell) String() string {
	return fmt.Sprintf("%#v@%d", cell.Token, cell.Index)
}

// BoxedCellPointers are used to represent paths; an instance is shared
// by linear chains of links, and only allocated when a link path is forked
// due to the addition of a 2nd, 3rd, ... Nth outbound link from a cell.
//type BoxedCellPointer *Cell
