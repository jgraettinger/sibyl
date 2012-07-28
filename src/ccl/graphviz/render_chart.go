package graphviz

import (
	. "ccl/chart"
    . "ccl/util"
	"fmt"
    "strings"
)

func RenderChart(chart *Chart) string {

	parts := []string{
		"digraph {",
		"  rankdir=LR;",
		"  tok_begin [label=\"{begin}\"];"}

	renderAdjacency := func(adjacency *Adjacency, left bool) string {
		var style, label, to string

		if adjacency.Used {
			style = "bold"
		} else if adjacency.IsBlocked() {
			style = "dotted"
		} else {
			style = "dashed"
		}

		if adjacency.Used {
			label = fmt.Sprintf("%d", adjacency.UsedDepth)
		}

		var weight int
		if adjacency.To != nil {
			to = fmt.Sprintf("tok_%d", adjacency.To.Index)
			weight = len(chart.Cells) - Iabs(adjacency.From.Index-adjacency.To.Index)
		} else if left {
			to = "tok_begin"
			weight = len(chart.Cells) - (adjacency.From.Index + 1)
		} else {
			to = "tok_end"
			weight = len(chart.Cells) - (len(chart.Cells) - adjacency.From.Index)
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

	for index, cell := range chart.Cells {
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
