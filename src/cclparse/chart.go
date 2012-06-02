package cclparse

import (
	"fmt"
	"strings"
	"invariant"
)

type Cell struct {
	Index uint
	Token string

	Inbound  struct{Left, Right AdjacencySet}
	Outbound struct{Left, Right AdjacencySet}
}

type Chart []*Cell

type Adjacency struct {

	From *Cell // inclusive
	To   *Cell // inclusive

	// the potential argument attachment
	//  position this adjacency reflects
	Position int

	// link properties
	Used bool
	Depth uint

	Blocked bool
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

func NewChart() (chart *Chart) {
	return new(Chart)
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell<%d, %s, %v, %v>",
		c.Index, c.Token, c.Inbound, c.Outbound)
}

/*func (l *CoverLink) String() string {
	return fmt.Sprintf("CoverLink<%s (%d), %s (%d), %d>",
		l.From.Token, l.From.Index, l.To.Token, l.To.Index, l.Depth)
}*/

func (chart *Chart) AddCell(token string) {

	var prevCell, nextCell *Cell

	if len(*chart) > 0 {
		prevCell = (*chart)[len(*chart) - 1]
	}

	nextCell = new(Cell)
	nextCell.Index = (uint)(len(*chart))
	nextCell.Token = token
	nextCell.Inbound.Left = make(AdjacencySet)
	nextCell.Inbound.Right = make(AdjacencySet)
	nextCell.Outbound.Left = make(AdjacencySet)
	nextCell.Outbound.Right = make(AdjacencySet)
	*chart = append(*chart, nextCell)

	if prevCell == nil {
		return
	}

	// initialize direct left-to-right adjacency
	{
		adjacency := new(Adjacency)
		adjacency.From = prevCell
		adjacency.To = nextCell
		adjacency.Position = 1

		prevCell.Outbound.Right.Add(adjacency)
		nextCell.Inbound.Left.Add(adjacency)
	}

	// initialize direct right-to-left adjacency
	{
		adjacency := new(Adjacency)
		adjacency.From = nextCell
		adjacency.To = prevCell
		adjacency.Position = -1

		nextCell.Outbound.Left.Add(adjacency)
		prevCell.Inbound.Right.Add(adjacency)
	}

	// for all *used* adjacencies ending at prevCell,
	//  create a new adjacency with position + 1
	// (there should be no more than one such used adjacency)
	{
		foundUsed := false
		for adjacency := range(prevCell.Inbound.Left) {
			if adjacency.Used {

				// we expect to see at most one used adjacency
				//  on the left side of prevCell
				invariant.IsTrue(!foundUsed)
				foundUsed = true

				// this adjacency is used; create a new one	with position + 1
				newAdjacency := new(Adjacency)
				newAdjacency.From = adjacency.From
				newAdjacency.To = nextCell
				newAdjacency.Position = adjacency.Position + 1

				newAdjacency.From.Outbound.Right.Add(newAdjacency)
				nextCell.Inbound.Left.Add(adjacency)
			}
		}
	}
}

func (chart *Chart) AsGraphviz() string {

	parts := []string{"digraph {", "  rankdir=LR;"}

	renderAdjacency := func(adjacency *Adjacency) string {
		var style, label string

		if adjacency.Used {
			style = "bold"
		} else if adjacency.Blocked {
			style = "dotted"
		} else {
			style = "dashed"
		}

		if adjacency.Used {
			label = fmt.Sprintf("%d", adjacency.Depth)
		}

		return fmt.Sprintf("  tok_%d -> tok_%d [label=\"%v\",style=\"%v\"]",
			adjacency.From.Index, adjacency.To.Index, label, style)
	}

	for index, cell := range(*chart) {
		parts = append(parts, fmt.Sprintf("  tok_%d [label=\"%v\",shape=\"box\"];",
			index, cell.Token))

		for adj := range(cell.Outbound.Left) {
			parts = append(parts, renderAdjacency(adj))
		}
		for adj := range(cell.Outbound.Right) {
			parts = append(parts, renderAdjacency(adj))
		}
	}
	parts = append(parts, "}")
	return strings.Join(parts, "\n")
}

