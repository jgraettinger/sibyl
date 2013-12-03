package chart

import (
	"fmt"
)

type CellSide struct {
	// The singular unused outbound adjacency, and potentially multiple
	// unused inbound adjacencies.
	OutboundAdjacency  *Adjacency

	// Adjacencies which have been used to form links. Multiple outbound
	// adjacencies may have been used, but only one inbound one can be.
	OutboundLinks []*Link
	InboundLink   *Link
}

type Cell struct {
	Index int
	Token Token

	Left, Right CellSide
}

func lastCell(list []*Cell) *Cell {
	if l := len(list); l != 0 {
		return list[l-1]
	}
	return nil
}

func (cell *Cell) String() string {
	return fmt.Sprintf("%#v@%d", cell.Token, cell.Index)
}
