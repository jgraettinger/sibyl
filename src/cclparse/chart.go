package cclparse

import (
	"fmt"
	"invariant"
	"strings"
)

const (
	LEFT  = 0
	RIGHT = 1
)

type Cell struct {
	Index int
	Token string

	Outbound [2]AdjacencyList
	Inbound  [2]AdjacencySet
}

type Chart struct {
	cells      []*Cell
	endInbound AdjacencySet
}

func NewChart() (chart *Chart) {
	chart = new(Chart)
	chart.endInbound = make(AdjacencySet)
	return chart
}

type Adjacency struct {
	From *Cell // inclusive
	To   *Cell // inclusive; may be nil

	// the potential argument attachment
	//  position this adjacency reflects
	Position int

	// link properties
	Used  bool
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
	*list = append(*list, adjacency)
	invariant.Equal(len(*list), iabs(adjacency.Position))
}

func (chart *Chart) AddCell(token string) {

	var prevCell, nextCell *Cell

	if len(chart.cells) > 0 {
		prevCell = chart.cells[len(chart.cells)-1]
	}

	nextCell = new(Cell)
	nextCell.Index = len(chart.cells)
	nextCell.Token = token
	nextCell.Inbound[LEFT] = make(AdjacencySet)
	nextCell.Inbound[RIGHT] = make(AdjacencySet)
	chart.cells = append(chart.cells, nextCell)

	// update all adjacencies to {end}, to be adjacent to nextCell
	for adjacency := range chart.endInbound {
		delete(chart.endInbound, adjacency)

		adjacency.To = nextCell
		nextCell.Inbound[LEFT].Add(adjacency)
	}

	// add nextCell => prevCell adjacency 
	{
		adjacency := new(Adjacency)
		adjacency.From = nextCell
		adjacency.To = prevCell
		adjacency.Position = -1

		nextCell.Outbound[LEFT].Add(adjacency)
		if prevCell != nil {
			prevCell.Inbound[RIGHT].Add(adjacency)
		}
	}

	// add nextCell => {end} adjacency
	{
		adjacency := new(Adjacency)
		adjacency.From = nextCell
		adjacency.Position = 1

		nextCell.Outbound[RIGHT].Add(adjacency)
		chart.endInbound.Add(adjacency)
	}
}

func (chart *Chart) AsGraphviz() string {

	parts := []string{
		"digraph {",
		"  rankdir=LR;",
		"  tok_begin [label=\"{begin}\"];"}

	renderAdjacency := func(adjacency *Adjacency, left bool) string {
		var style, label, to string

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

        var weight int
		if adjacency.To != nil {
			to = fmt.Sprintf("tok_%d", adjacency.To.Index)
            weight = len(chart.cells) - iabs(adjacency.From.Index - adjacency.To.Index)
		} else if left {
			to = "tok_begin"
            weight = len(chart.cells) - (adjacency.From.Index + 1)
		} else {
			to = "tok_end"
            weight = len(chart.cells) - (len(chart.cells) - adjacency.From.Index)
		}

        var constraint bool
        if adjacency.Position < 0 {
            constraint = false
        } else {
            constraint = true
        }

		return fmt.Sprintf("  tok_%d -> %v [label=\"%v\",style=\"%v\",constraint=%v,weight=%v]",
			adjacency.From.Index, to, label, style, constraint, weight)
	}

	for index, cell := range chart.cells {
		parts = append(parts, fmt.Sprintf("  tok_%d [label=\"%v\",shape=\"box\"];",
			index, cell.Token))

		for _, adj := range cell.Outbound[LEFT] {
			parts = append(parts, renderAdjacency(adj, true))
		}
		for _, adj := range cell.Outbound[RIGHT] {
			parts = append(parts, renderAdjacency(adj, false))
		}
	}
	parts = append(parts, "  tok_end [label=\"{end}\"];", "}")
	return strings.Join(parts, "\n")
}
