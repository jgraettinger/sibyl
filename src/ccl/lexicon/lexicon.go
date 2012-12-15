package lexicon

import (
	"ccl/chart"
	"invariant"
	"strings"
)

type Lexicon map[AdjacencyPoint]*AdjacencyStatistics

func New() Lexicon {
	return make(Lexicon)
}

func (lexicon Lexicon) Score(adjacency *chart.Adjacency) (
	linkWeight float64, linkDepth uint) {

	if adjacency.To == nil {
		// empty adjacency
		return
	}

	yIn := lexicon[AdjacencyPoint{string(adjacency.To.Token),
		-isign(adjacency.Position)}]
	if yIn == nil {
		return
	}

	// check each adjacency point of From, starting at the adjacency
	//  position and working backwards to position 1
	position := adjacency.Position
	for ; linkWeight == 0 && position != 0; position -= isign(position) {
		xOut := lexicon[AdjacencyPoint{string(adjacency.From.Token), position}]
		if xOut == nil {
			continue
		}
		linkWeight, linkDepth = lexicon.linkWeight(xOut, yIn)
	}
	return
}

func (lexicon Lexicon) linkWeight(xOut, yIn *AdjacencyStatistics) (
	linkWeight float64, linkDepth uint) {

	invariant.NotNil(xOut)
	invariant.NotNil(yIn)

	label, labelWeight := bestMatchingLabel(xOut, yIn)
	if labelWeight == 0 {
		return
	}

	// use the best matching label to determine the 'prototype'
	//  statistics to use in calculating link weight;
	// prototypes tend to be frequent words, and conceptually
	//  describe the relationship between x & y
	var prototype *AdjacencyStatistics
	if label.IsClass() {
		// label prototypes x along with tokens separating x from y
		prototype = lexicon[AdjacencyPoint{label.Token, xOut.Sign()}]
	} else {
		// label prototypes y in it's relationship with x
		prototype = lexicon[AdjacencyPoint{label.Token, -xOut.Sign()}]
	}

	if prototype == nil {
		return
	}

	if label.IsClass() && prototype.Out > 0 {
		linkWeight = fmin(labelWeight, prototype.OutNorm())
		return
	}

	if label.IsAdjacency() {
		if prototype.In > 0 {
			linkWeight = fmin(labelWeight, prototype.InNorm())
		} else if float64(prototype.InRaw) > fabs(prototype.In) {
			linkWeight = fmin(labelWeight, prototype.InRawNorm())
		}

		if linkWeight != 0 {
			if float64(prototype.InRaw) < 0 && prototype.Out <= 0 {
				linkDepth = 1
			}
			return
		}
	}

	if prototype.Out <= 0 && prototype.In <= 0 && (label.IsAdjacency() || prototype.Out == 0) {
		linkWeight = labelWeight
		return
	}
	return
}

func (lexicon Lexicon) Learn(ch *chart.Chart) {
	var deltas []*AdjacencyStatistics

	// Closure which adds to 'deltas' the appropriate
	// lexicon update for an argument Adjacency.
	updateFromAdjacency := func(adjacency *chart.Adjacency) {

		point := AdjacencyPoint{adjacency.From.Token.Str(), adjacency.Position}
		delta := NewAdjacencyStatistics(point)

		if adjacency.To == nil ||
			adjacency.RestrictedDepths(ch) == chart.RESTRICT_ALL {
			delta.updateFromBlocking()
		} else {
			delta.update(lexicon, adjacency.To.Token.Str())
		}
		deltas = append(deltas, delta)
	}

	updateFromLink := func(link *chart.Link) {
		point := AdjacencyPoint{link.From.Token.Str(), link.Position}
		delta := NewAdjacencyStatistics(point)

		delta.update(lexicon, link.To.Token.Str())
		deltas = append(deltas, delta)
	}

	// Compute lexicon update deltas for every outbound link
	// of every cell in the chart. By computing deltas rather
	// than directly updating the lexicon, we isolate updates
	// from early cells from affecting updates from later ones.
	for _, cell := range ch.Cells {

		updateFromAdjacency(cell.OutboundAdjacency[chart.LEFT])
		updateFromAdjacency(cell.OutboundAdjacency[chart.RIGHT])

		for _, link := range cell.OutboundLinks[chart.LEFT] {
			updateFromLink(link)
		}
		for _, link := range cell.OutboundLinks[chart.RIGHT] {
			updateFromLink(link)
		}
	}
	// Fold each delta into the lexicon.
	for _, delta := range deltas {
		if stats, found := lexicon[delta.AdjacencyPoint]; !found {
			lexicon[delta.AdjacencyPoint] = delta
		} else {
			stats.fold(delta)
		}
	}
}

func (this Lexicon) String() string {
	var parts []string
	for _, adjacency := range this {
		parts = append(parts, adjacency.String())
	}
	return strings.Join(parts, "\n")
}

func isign(a int) int {
	invariant.NotEqual(a, 0)
	if a < 0 {
		return -1
	}
	return 1
}

func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func fabs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
