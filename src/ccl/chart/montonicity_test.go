package chart

import (
	assert "invariant"
	"testing"
)

func TestMontonicity(t *testing.T) {
	chart := NewChart()
	Y := chart.AddCell("Y")
	Z := chart.AddCell("Z")

	assert.IsFalse(Y.OutboundAdjacency[RIGHT].MontonicityRestricted)
	assert.IsFalse(Z.OutboundAdjacency[LEFT].MontonicityRestricted)

	chart.Use(Y.OutboundAdjacency[RIGHT], 1)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)

	// Outbound adjacencies created through linking have depth 0 blocked.
	assert.IsTrue(Y.OutboundAdjacency[RIGHT].MontonicityRestricted)
	assert.IsTrue(Z.OutboundAdjacency[LEFT].MontonicityRestricted)
}
