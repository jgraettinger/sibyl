package cclparse
/*
import (
	"invariant"
	"strings"
    "ccl/parse"
)


type Lexicon map[AdjacencyPoint]*AdjacencyStatistics

func NewLexicon() Lexicon {
	return make(Lexicon)
}

func (lexicon Lexicon) linkWeight(xOut, yIn *AdjacencyStatistics) (
	linkWeight float64, linkDepth uint8) {

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

func (lexicon Lexicon) Score(adjacency *parse.Adjacency) (
	linkWeight float64, linkDepth uint8) {

	if adjacency.To == nil {
		// empty adjacency
		return
	}

	yIn := lexicon[AdjacencyPoint{adjacency.To.Token,
		-isign(adjacency.Position)}]
	if yIn == nil {
		return
	}

	// check each adjacency point of From, starting at the adjacency
	//  position and working backwards to position 1
	position := adjacency.Position
	for ; linkWeight == 0 && position != 0; position -= isign(position) {

		xOut := lexicon[AdjacencyPoint{adjacency.From.Token, position}]
		if xOut == nil {
			continue
		}
		linkWeight, linkDepth = lexicon.linkWeight(xOut, yIn)
	}
	return
}

func (lexicon Lexicon) Learn(chart *parse.Chart) {

	var deltas []*AdjacencyStatistics

	// Closure which adds to 'deltas' the appropriate
	// lexicon update for an argument Adjacency.
	update := func(adjacency *Adjacency) {

		point := AdjacencyPoint{adjacency.From.Token, adjacency.Position}
		delta := NewAdjacencyStatistics(point)

		if adjacency.To == nil || adjacency.StoppingPunc {
			delta.updateFromBlocking()
		} else {
			delta.update(lexicon, adjacency.To.Token)
		}

		deltas = append(deltas, delta)
	}

	// Compute lexicon update deltas for every outbound link
	// of every cell in the chart. By computing deltas rather
	// than directly updating the lexicon, we isolate updates
	// from early cells from affecting updates from later ones.
	for _, cell := range chart.cells {
		for _, adjacency := range cell.Outbound[LEFT] {

			// TODO HACK : I don't think this should be here. To match cclparse behavior
			if !adjacency.Used && adjacency.Position <= -2 && adjacency.To != nil {
				if _, okay := lexicon[AdjacencyPoint{adjacency.To.Token, 1}]; !okay {
					continue
				}
			}

			update(adjacency)
		}
		for _, adjacency := range cell.Outbound[RIGHT] {
			update(adjacency)
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
*/
