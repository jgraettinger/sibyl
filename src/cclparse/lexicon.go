package cclparse

import (
	"invariant"
	"strings"
)

type Lexicon map[AdjacencyPoint]*AdjacencyStatistics

func NewLexicon() Lexicon {
	return make(Lexicon)
}

func (lexicon Lexicon) linkWeight(xOut, yIn *AdjacencyStatistics,
) (linkWeight float64, linkDepth uint8) {

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

	if label.IsClass() && prototype.out > 0 {
		linkWeight = fmin(labelWeight, prototype.OutNorm())
		return
	}

	if label.IsAdjacency() {
		if prototype.in > 0 {
			linkWeight = fmin(labelWeight, prototype.InNorm())
		} else if float64(prototype.inRaw) > fabs(prototype.in) {
			linkWeight = fmin(labelWeight, prototype.InRawNorm())
		}

		if linkWeight != 0 {
			if float64(prototype.inRaw) < 0 && prototype.out <= 0 {
				linkDepth = 1
			}
			return
		}
	}

	if prototype.out <= 0 && prototype.in <= 0 && (label.IsAdjacency() || prototype.out == 0) {
		linkWeight = labelWeight
		return
	}
	return
}

func (lexicon Lexicon) Score(adjacency *Adjacency
    ) (linkWeight float64, linkDepth uint8) {

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
    for ; linkWeight == 0 && position != 0; position - isign(position) {

        xOut := lexicon[AdjacencyPoint{adjacency.From.Token, position}]
        if xOut == nil {
            continue
        }
        linkWeight, linkDepth = linkWeight(xOut, yIn)
    }
    return
}

func (this Lexicon) Learn(chart *Chart) {

    var deltas []*AdjacencyStatistics

	update := func(adjacency *Adjacency) {

		point := AdjacencyPoint{adjacency.From.Token, adjacency.Position}
        delta := NewAdjacencyStatistics(point)

		if adjacency.To == nil {
			delta.updateFromBlocking()
		} else {
			delta.update(this, adjacency.To.Token)
		}

        deltas = append(deltas, delta)
	}

	for _, cell := range chart.cells {
		for _, adjacency := range cell.Outbound.Left {
			update(adjacency)
		}
		for _, adjacency := range cell.Outbound.Right {
			update(adjacency)
		}
	}

    for _, delta := range deltas {
        if stats, found := this[delta.AdjacencyPoint]; !found {
            this[delta.AdjacencyPoint] = delta
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
