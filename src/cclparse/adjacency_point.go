package cclparse

import (
	"fmt"
	"invariant"
	"strings"
)

type AdjacencyPoint struct {
	Token    string
	Position int
}

func (this AdjacencyPoint) Sign() int {
	if this.Position < 0 {
		return -1
	}
	return 1
}

type AdjacencyStatistics struct {
	AdjacencyPoint

	count uint64
	stop  uint64
	inRaw int64

	out float64
	in  float64

	labelWeights map[Label]float64
}

func NewAdjacencyStatistics(point AdjacencyPoint) (
	this *AdjacencyStatistics) {

	this = new(AdjacencyStatistics)
	this.AdjacencyPoint = point
	this.labelWeights = make(map[Label]float64)
	return
}

func (this *AdjacencyStatistics) Count() float64 {
	if this == nil {
		return 0
	}
	return float64(this.count)
}
func (this *AdjacencyStatistics) Stop() float64 {
	if this == nil {
		return 0
	}
	return float64(this.stop)
}
func (this *AdjacencyStatistics) InRaw() float64 {
	if this == nil {
		return 0
	}
	return float64(this.inRaw)
}
func (this *AdjacencyStatistics) InRawNorm() float64 {
	if this == nil {
		return 0
	}
	return float64(this.inRaw) / float64(this.count)
}
func (this *AdjacencyStatistics) Out() float64 {
	if this == nil {
		return 0
	}
	return this.out
}
func (this *AdjacencyStatistics) OutNorm() float64 {
	if this == nil {
		return 0
	}
	return this.out / float64(this.count)
}
func (this *AdjacencyStatistics) In() float64 {
	if this == nil {
		return 0
	}
	return this.in
}
func (this *AdjacencyStatistics) InNorm() float64 {
	if this == nil {
		return 0
	}
	return this.in / float64(this.count)
}
func (this *AdjacencyStatistics) LabelWeightNorm(label Label) float64 {
	if this == nil {
		return 0
	}
	return this.labelWeights[label] / float64(this.count)
}
func (this *AdjacencyStatistics) HasLargeStop() bool {
	if this == nil {
		return true
	}
	for _, weight := range this.labelWeights {
		if weight > float64(this.stop) {
			return false
		}
	}
	return true
}

func bestMatchingLabel(xOut, yIn *AdjacencyStatistics) (
	bestLabel Label, bestWeight float64) {

	invariant.NotNil(xOut)
	invariant.NotNil(yIn)

	for label, weight := range xOut.labelWeights {
		if weight <= xOut.Stop() {
			continue
		}

		if label.Type == ADJACENCY && label.Token == yIn.Token {
			// l = (y, 1); corresponding weight of yIn is defined to be 1
			weight /= xOut.Count()
		} else {
			weight = fmin(weight/xOut.Count(),
				yIn.LabelWeightNorm(label.Flip()))
		}

		if weight > bestWeight {
			bestWeight = weight
			bestLabel = label
		}
	}
	return
}

func (this *AdjacencyStatistics) update(lexicon Lexicon, token string) {

	this.count += 1

	// tick label reflecting direct adjacency with token
	directLabel := Label{ADJACENCY, token}
	this.labelWeights[directLabel] += 1

	// lookup this adjacency's inverse
	inverse := lexicon[AdjacencyPoint{token, -this.Sign()}]

	if inverse != nil {
		// for other labels of the inverse adjacency, add normalized weight
		//  contributions to corresonding flipped labels. Very reminiscent of the
		//  power method for finding the largest eigen-value of a sparse matrix...
		norm := float64(inverse.count)
		for label, weight := range inverse.labelWeights {
			if label != directLabel {
				this.labelWeights[label.Flip()] += weight / norm
			}
		}
	}

	if iabs(this.Position) != 1 {
		return
	}

	// calculate bootstrap 'In*' statistic update
	if inverse == nil || inverse.HasLargeStop() {
		if inverse == nil {
			// all other updates are effectively zero
			return
		}

		// TODO HACK - I think this belongs before the return
		//  and shouldn't be guarded by labelWeights

		if len(inverse.labelWeights) != 0 {
			// a link from inverse to this is unlikely
			this.inRaw += -1
		}
	} else {
		flippedInverse := lexicon[AdjacencyPoint{token, this.Sign()}]

		if flippedInverse == nil || flippedInverse.HasLargeStop() {
			this.inRaw += 1
		}
	}
	invariant.NotNil(inverse)

	// calculate smoothed Out & In updates
	this.out += inverse.InRawNorm()
	this.in += inverse.OutNorm()

}

func (this *AdjacencyStatistics) fold(other *AdjacencyStatistics) {

	invariant.Equal(this.AdjacencyPoint, other.AdjacencyPoint)

	this.count += other.count
	this.stop += other.stop
	this.inRaw += other.inRaw
	this.out += other.out
	this.in += other.in

	for label, weight := range other.labelWeights {
		this.labelWeights[label] += weight
	}
}

func (this *AdjacencyStatistics) updateFromBlocking() {
	this.count += 1
	this.stop += 1
}

func (this *AdjacencyPoint) String() string {
	return fmt.Sprintf("%#v %d", this.Token, this.Position)
}

func (this *AdjacencyStatistics) String() string {
	parts := []string{
		fmt.Sprintf("Point: %v", this.AdjacencyPoint.String()),
		fmt.Sprintf("\tCount %v Stop %v In* %v Out %v In %v",
			this.count, this.stop, this.inRaw, this.out, this.in)}

	for label, weight := range this.labelWeights {
		parts = append(parts, fmt.Sprintf("\t%v: %v", label, weight))
	}
	return strings.Join(parts, "\n")
}
