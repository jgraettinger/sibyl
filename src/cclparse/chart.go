package cclparse

import (
	"fmt"
	"strings"
	"invariant"
)

type Cell struct {
	Index uint
	Token string

	Outbound struct{Left, Right AdjacencyList}
	Inbound  struct{Left, Right AdjacencySet}
}

type Chart struct {
	cells []*Cell
	endInbound AdjacencySet
}

func NewChart() (chart *Chart) {
    chart = new(Chart)
    chart.endInbound = make(AdjacencySet)
}

type Adjacency struct {

	From *Cell // inclusive
	To   *Cell // inclusive; may be nil

	// the potential argument attachment
	//  position this adjacency reflects
	Position int

	// link properties
	Used bool
	Depth uint

	Blocked bool
}

type AdjacencySet map[*Adjacency]bool
type AdjacencyList []*Adjacency

func (set AdjacencySet) Add(adjacency *Adjacency) {
	set[adjacency] = true
}

func (set AdjacencySet) Remove(adjacency *Adjacency) {
	_, present := set[adjacency]
	invariant.IsTrue(present)
	delete(set, adjacency)
}

func (list *AdjacencyList) Add(adjacency *Adjacency) {
    list = append(list, adjacency)
    invariant.IsEqual(len(*list), iabs(adjacency.Position))
}

func (chart *Chart) AddCell(token string) {

	var prevCell, nextCell *Cell

	if len(*chart) > 0 {
		prevCell = (*chart)[len(*chart) - 1]
	}

	nextCell = new(Cell)
	nextCell.Index = (uint)(len(chart.cells))
	nextCell.Token = token
	nextCell.Inbound.Left = make(AdjacencySet)
	nextCell.Inbound.Right = make(AdjacencySet)
	chart.cells = append(chart.cells, nextCell)

    // update all adjacencies to {end}, to be adjacent to nextCell
    for adjacency := range(chart.endInbound) {
        delete(chart.endInbound, adjacency)

        adjacency.To = nextCell
        nextCell.Inbound.Left.Add(adjacency)
    }

    // add nextCell => prevCell adjacency 
    {
		adjacency := new(Adjacency)
		adjacency.From = nextCell
		adjacency.To = prevCell
		adjacency.Position = -1

		nextCell.Outbound.Left.Add(adjacency)
		prevCell.Inbound.Right.Add(adjacency)
	}

    // add nextCell => {end} adjacency
    {
        adjacency := new(Adjacency)
        adjacency.From = nextCell
        adjacency.Position = 1

        nextCell.Outbound.Right.Add(adjacency)
        chart.endInbound.Add(adjacency)
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

