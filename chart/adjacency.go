package chart

import (
	"fmt"
)

type Adjacency struct {
	// Head is always non-nil, but Tail may be nil to represent
	// an adjacency to the current utterance end.
	Head, Tail *Cell

	// The argument attachment position this adjacency reflects
	Position int
}

func (adjacency *Adjacency) HeadSide() *CellSide {
	if adjacency.Position < 0 {
		return &adjacency.Head.Left
	} else {
		return &adjacency.Head.Right
	}
}
func (adjacency *Adjacency) TailSide() *CellSide {
	if adjacency.Position < 0 {
		return &adjacency.Tail.Right
	} else {
		return &adjacency.Tail.Left
	}
}
func (adjacency *Adjacency) appendTo(list *[]*Adjacency) {
	*list = append(*list, adjacency)
}

func (adjacency *Adjacency) String() string {
	return fmt.Sprintf("Adjacency<%v:%d => %v>",
		adjacency.Head, adjacency.Position, adjacency.Tail)
}
