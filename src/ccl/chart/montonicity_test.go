package chart

import (
	"testing"
)

func checkMontonicity(t *testing.T, adjacency *Adjacency,
	expect DepthRestriction) {

	if r := adjacency.MontonicityRestriction(); r != expect {
		t.Errorf("Expected montonicity %v on %v, but got %v",
			expect, adjacency, r)
	}
}

func TestMontonicity(t *testing.T) {
	chart := NewChart()
	Y := chart.AddCell("Y")
	Z := chart.AddCell("Z")

	checkMontonicity(t, Y.OutboundAdjacency[RIGHT], RESTRICT_NONE)
	checkMontonicity(t, Z.OutboundAdjacency[LEFT], RESTRICT_NONE)

	chart.Use(Y.OutboundAdjacency[RIGHT], 1)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)

	// Outbound adjacencies created through linking have depth 0 blocked.
	checkMontonicity(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkMontonicity(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}
