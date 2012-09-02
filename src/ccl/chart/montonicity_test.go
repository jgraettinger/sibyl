package chart

import (
	"testing"
)

func checkMontonicity(t *testing.T, adjacency *Adjacency, e DepthRestriction) {
	if r := adjacency.MontonicityRestriction(); r != e {
		t.Errorf("Expected montonicity %v on %v, but got %v",
			e, adjacency, r)
	}
}

func TestMontonicity(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("Y", "Z")

	Y, Z := chart.nextCell(), chart.nextCell()
	checkMontonicity(t, Y.OutboundAdjacency[RIGHT], RESTRICT_NONE)
	checkMontonicity(t, Z.OutboundAdjacency[LEFT], RESTRICT_NONE)

	chart.use(Y.OutboundAdjacency[RIGHT], 1)
	chart.use(Z.OutboundAdjacency[LEFT], 1)

	// Outbound adjacencies created through linking have depth 0 blocked.
	checkMontonicity(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkMontonicity(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}
