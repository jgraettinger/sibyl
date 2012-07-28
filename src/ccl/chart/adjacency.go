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

	BlockedDepths [2]bool
}

func (adjacency *Adjacency) IsBlocked() bool {
    // TODO: Get rid of this?
    return adjacency.BlockedDepths[0] && adjacency.BlockedDepths[1]
}

func (adjacency *Adjacency) ToIndex(chart *Chart) int {
    if adjacency.To != nil {
        return adjacency.To.Index
    }
    if adjacency.Position < 0 {
        return -1
    }
    return len(chart.Cells)
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
	if adjacency.BlockedDepths[0] {
		properties += ", (Blocked 0)"
	}
	if adjacency.BlockedDepths[1] {
		properties += ", (Blocked 1)"
	}

    return fmt.Sprintf("Adjacency<%v:%d => %v%v>", adjacency.From,
        adjacency.Position, adjacency.To, properties)
}
