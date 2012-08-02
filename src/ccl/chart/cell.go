package chart

import (
	"fmt"
)

type Cell struct {
	Index int
	Token string

	// Marks the index of the closest cell which has an outbound d=1
	// link, and also has a link-path back to this cell. Nil if there
	// is none. This is used to enforce blocking: specifically, no
	// adjacency from this node may span beyond BlockedAfter.
	// (See section 3.2.1, condition 3)
	// FullyBlockedAfter [2]*int

	// Active (non-linked) outbound adjacencies. Each side has exactly one.
	OutboundAdjacency [2]*Adjacency
	// Active (non-linked) inbound adjacencies.
	InboundAdjacencies [2]AdjacencySet

	OutboundLinks [2]LinkList
	InboundLink   [2]*Link

	// Last-added d=0 & d=1 links, reflecting the furthest link of each depth
	// from this cell. Due to montonicity, the d=1 link will either be nil,
	// or will be a further adjancency than the d=0 link.
	LastOutboundLinkD0 [2]*Link
	LastOutboundLinkD1 [2]*Link
}

// BoxedCellPointers are used to represent paths; an instance is shared
// by linear chains of links, and only allocated when a link path is forked
// due to the addition of a 2nd, 3rd, ... Nth outbound link from a cell.
type BoxedCellPointer *Cell

func NewCell(index int, token string) *Cell {
	cell := new(Cell)
	cell.Index = index
	cell.Token = token
	cell.InboundAdjacencies[LEFT] = make(AdjacencySet)
	cell.InboundAdjacencies[RIGHT] = make(AdjacencySet)
	return cell
}

func (cell *Cell) String() string {
	return fmt.Sprintf("%#v@%d", cell.Token, cell.Index)
}
