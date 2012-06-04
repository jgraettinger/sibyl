package cclparse

import (
    "invariant"
)

type Lexicon map[AdjacencyPoint]AdjacencyStatistics

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

	if label.IsClass() && prototype.Out > 0 {
		linkWeight = fmin(labelWeight, prototype.OutNorm())
		return
	}

	if label.IsAdjacency() {
		if prototype.In > 0 {
			linkWeight = fmin(labelWeight, prototype.InNorm())
		} else if prototype.InRaw > fabs(prototype.In) {
			linkWeight = fmin(labelWeight, prototype.InRawNorm())
		}

		if linkWeight != 0 {
			if prototype.InRaw < 0 && prototype.Out <= 0 {
				linkDepth = 1
			}
			return
		}
	}

	if prototype.Out <= 0 && prototype.In <= 0 && (
			label.IsAdjacency() || prototype.Out == 0) {
		linkWeight = labelWeight
		return
	}
	return
}

func (lexicon Lexicon) Update(adj *Adjacency) {

    var fromStats, toStats *AdjacencyStatistics

    fromStats = lexicon.Intern(AdjacencyPoint{adj.From.Token, adj.Position})

    stats.Count++

    if adj.To == nil { // || adj.To.IsPunctuation
        // track that this adjacency is unusable
        stats.Stop++
        return
    }

    // increment direct adjacency label count
    stats.LabelWeights[Label{ADJACENCY, adjacency.To.Token}] += 1

    sameSign := lexicon[AdjacencyPoint{adj.From.Token
    for label, weight := range(
        lexicon[AdjacencyPoint{adj.To.Token, -from.Sign()}) {

        // Add activation

    }

    to := AdjacencyPoint{adjacency.To.Token, -from.Sign()}


}

func (lexicon Lexicon) Learn(chart Chart) {

    for _, cell := range(chart) {

        if cell.Index == 0 {
            // blocked adjacency to start-of-sentence

        }

    }


}

