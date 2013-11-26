package graphviz

import (
	"ccl/chart"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	//"strings"
)

func RenderChartSvg(cht *chart.Chart, w io.Writer) {

	canvas := svg.New(w)
	canvas.Start(100+100*len(cht.Cells), 100)

	for i, cell := range cht.Cells {
		canvas.Roundrect(100+100*i, 25, 75, 50, 10, 10, "fill:none;stroke:black")
		canvas.Text(110+100*i, 45, string(cell.Token), "")

		for j, link := range cell.OutboundLinks[chart.LEFT] {
			id := fmt.Sprintf("L%v_%v", i, j)
			from, to := 150 + 100 * i, 150 + 100 * link.To.Index
			canvas.Qbez(from, 75, (from + to)/2, 100, to, 75,
				"fill='none'", "stroke='black'", "id='"+id+"'")
			canvas.Textpath(fmt.Sprintf("    d=%v", link.Depth), "#"+id)
		}
		for j, link := range cell.OutboundLinks[chart.RIGHT] {
			id := fmt.Sprintf("R%v_%v", i, j)
			from, to := 150 + 100 * i, 150 + 100 * link.To.Index
			canvas.Qbez(from, 25, (from + to)/2, 0, to, 25,
				"fill='none'", "stroke='black'", "id='"+id+"'")
			canvas.Textpath(fmt.Sprintf("    d=%v", link.Depth), "#"+id)
		}
	}
	canvas.End()
}

/*
func RenderChart(ch *chart.Chart) string {

	parts := []string{
		"digraph {",
		"  rankdir=LR;",
		"  tok_begin [label=\"{begin}\"];"}

	renderAdjacency := func(adjacency *chart.Adjacency, left bool) string {
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
			weight = len(ch.Cells) - Iabs(adjacency.From.Index-adjacency.To.Index)
		} else if left {
			to = "tok_begin"
			weight = len(ch.Cells) - (adjacency.From.Index + 1)
		} else {
			to = "tok_end"
			weight = len(ch.Cells) - (len(ch.Cells) - adjacency.From.Index)
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

	for index, cell := range ch.Cells {
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
*/
