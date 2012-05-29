package cclparse

import (
	"invariant"
)

type AdjacencyPoint struct {
    Token string
    Position int
}

type AdjacencyStatistics struct {
    AdjacencyPoint

    Count uint64
    Stop uint64

    InRaw float64
    Out float64
    In float64

    LabelWeights map[Label]float64
}

func NewAdjacencyStatistics(token string, position int) (
		this *AdjacencyStatistics) {

	this = new(AdjacencyStatistics)
	this.Token = token
	this.Position = position
	this.LabelWeights = make(map[Label]float64)
	return
}

func (this AdjacencyPoint) Sign() int {
	if this.Position < 0 {
		return -1
	}
	return 1
}

func (this *AdjacencyStatistics) InRawNorm() float64 {
	return this.InRaw / float64(this.Count)
}
func (this *AdjacencyStatistics) OutNorm() float64 {
	return this.Out / float64(this.Count)
}
func (this *AdjacencyStatistics) InNorm() float64 {
	return this.In / float64(this.Count)
}

func bestMatchingLabel(xOut, yIn *AdjacencyStatistics) (
		bestLabel Label, bestWeight float64) {

    invariant.NotNil(xOut)
    invariant.NotNil(yIn)

    for label, weight := range(xOut.LabelWeights) {
        if weight <= float64(xOut.Stop) {
			continue
        }

		if label.Type == ADJACENCY && label.Token == yIn.Token {
			// l = (y, 1); corresponding weight of yIn is defined to be 1
			weight /= float64(xOut.Count)
		} else {
			weight = fmin(weight / float64(xOut.Count),
				yIn.LabelWeights[label.Flip()] / float64(yIn.Count))
		}

		if weight > bestWeight {
			bestWeight = weight
			bestLabel = label
		}
    }
    return
}

