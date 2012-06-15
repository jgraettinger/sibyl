package cclparse

import (
	"encoding/json"
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

type LabelWeights map[Label]float64

type AdjacencyStatistics struct {
	AdjacencyPoint

	Count uint64
	Stop  uint64
	InRaw int64

	Out float64
	In  float64

	LabelWeights LabelWeights
}

func NewAdjacencyStatistics(point AdjacencyPoint) (
	this *AdjacencyStatistics) {

	this = new(AdjacencyStatistics)
	this.AdjacencyPoint = point
	this.LabelWeights = make(map[Label]float64)
	return
}

func (this *AdjacencyStatistics) InRawNorm() float64 {
	if this == nil {
		return 0
	}
	return float64(this.InRaw) / float64(this.Count)
}
func (this *AdjacencyStatistics) OutNorm() float64 {
	if this == nil {
		return 0
	}
	return this.Out / float64(this.Count)
}
func (this *AdjacencyStatistics) InNorm() float64 {
	if this == nil {
		return 0
	}
	return this.In / float64(this.Count)
}
func (this *AdjacencyStatistics) LabelWeightNorm(label Label) float64 {
	if this == nil {
		return 0
	}
	return this.LabelWeights[label] / float64(this.Count)
}
func (this *AdjacencyStatistics) HasLargeStop() bool {
	if this == nil {
		return true
	}
	for _, weight := range this.LabelWeights {
		if weight > float64(this.Stop) {
			return false
		}
	}
	return true
}

func bestMatchingLabel(xOut, yIn *AdjacencyStatistics) (
	bestLabel Label, bestWeight float64) {

	invariant.NotNil(xOut)
	invariant.NotNil(yIn)

	for label, weight := range xOut.LabelWeights {
		if weight <= float64(xOut.Stop) {
			continue
		}

		if label.Type == ADJACENCY && label.Token == yIn.Token {
			// l = (y, 1); corresponding weight of yIn is defined to be 1
			weight /= float64(xOut.Count)
		} else {
			weight = fmin(
				xOut.LabelWeightNorm(label), yIn.LabelWeightNorm(label.Flip()))
		}

		if weight > bestWeight {
			bestWeight = weight
			bestLabel = label
		}
	}
	return
}

func (this *AdjacencyStatistics) update(lexicon Lexicon, token string) {

	this.Count += 1

	// tick label reflecting direct adjacency with token
	directLabel := Label{ADJACENCY, token}
	this.LabelWeights[directLabel] += 1

	// lookup this adjacency's inverse
	inverse := lexicon[AdjacencyPoint{token, -this.Sign()}]

	if inverse != nil {
		// for other labels of the inverse adjacency, add normalized weight
		//  contributions to corresonding flipped labels. Very reminiscent of the
		//  power method for finding the largest eigen-value of a sparse matrix...
		norm := float64(inverse.Count)
		for label, weight := range inverse.LabelWeights {
			if label != directLabel {
				this.LabelWeights[label.Flip()] += weight / norm
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
		//  and shouldn't be guarded by LabelWeights

		if len(inverse.LabelWeights) != 0 {
			// a link from inverse to this is unlikely
			this.InRaw += -1
		}
	} else {
		flippedInverse := lexicon[AdjacencyPoint{token, this.Sign()}]

		if flippedInverse == nil || flippedInverse.HasLargeStop() {
			this.InRaw += 1
		}
	}
	invariant.NotNil(inverse)

	// calculate smoothed Out & In updates
	this.Out += inverse.InRawNorm()
	this.In += inverse.OutNorm()

}

func (this *AdjacencyStatistics) fold(other *AdjacencyStatistics) {

	invariant.Equal(this.AdjacencyPoint, other.AdjacencyPoint)

	this.Count += other.Count
	this.Stop += other.Stop
	this.InRaw += other.InRaw
	this.Out += other.Out
	this.In += other.In

	for label, weight := range other.LabelWeights {
		this.LabelWeights[label] += weight
	}
}

func (this *AdjacencyStatistics) updateFromBlocking() {
	this.Count += 1
	this.Stop += 1
}

func (this *AdjacencyPoint) String() string {
	return fmt.Sprintf("%#v %d", this.Token, this.Position)
}

func (this *AdjacencyStatistics) String() string {
	parts := []string{
		fmt.Sprintf("Point: %v", this.AdjacencyPoint.String()),
		fmt.Sprintf("\tCount %v Stop %v In* %v Out %v In %v",
			this.Count, this.Stop, this.InRaw, this.Out, this.In)}

	for label, weight := range this.LabelWeights {
		parts = append(parts, fmt.Sprintf("\t%v: %v", label, weight))
	}
	return strings.Join(parts, "\n")
}

func (self *AdjacencyStatistics) MarshalJSON() ([]byte, error) {

	type AsJson struct {
		Token        string
		Position     int
		Count        uint64
		Stop         uint64
		InRaw        int64
		Out          float64
		In           float64
		LabelWeights LabelWeights
	}

	asJson := AsJson{
		Token:        self.Token,
		Position:     self.Position,
		Count:        self.Count,
		Stop:         self.Stop,
		InRaw:        self.InRaw,
		Out:          self.Out,
		In:           self.In,
		LabelWeights: self.LabelWeights,
	}
    return json.Marshal(asJson)
}

func (labelWeights LabelWeights) MarshalJSON() (result []byte, err error) {
	stringWeights := make(map[string]float64)

	for label, weight := range labelWeights {
		stringWeights[label.String()] = weight
	}
	result, err = json.Marshal(stringWeights)
	return
}
