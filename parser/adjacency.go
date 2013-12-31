package parser

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
func (a *Adjacency) appendTo(list *[]*Adjacency) {
	*list = append(*list, a)
}

func (a *Adjacency) String() string {
	return fmt.Sprintf("%v:%d => %v", a.Head, a.Position, a.Tail)
}

func (a *Adjacency) HeadSide() *CellSide {
	if a.Position < 0 {
		return &a.Head.Left
	}
	return &a.Head.Right
}
func (a *Adjacency) TailSide() *CellSide {
	if a.Position < 0 {
		return &a.Tail.Right
	}
	return &a.Tail.Left
}
