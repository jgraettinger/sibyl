package chart

import (
	"testing"
)

func checkEquality(t *testing.T, adjacency *Adjacency,
	expect DepthRestriction) {

	if r := adjacency.EqualityRestriction(); r != expect {
		t.Errorf("Expected equality %v on %v, but got %v",
			expect, adjacency, r)
	}
}

func TestEquality_NeighborCycle(t *testing.T) {
	chart, V, W, _, Y, Z := buildFixture()

	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)

	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_D1)
	checkEquality(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
}

func TestEquality_Extended(t *testing.T) {
	chart, V, W, X, Y, Z := buildFixture()

	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 0)
	chart.Use(X.OutboundAdjacency[LEFT], 0)
	chart.Use(X.OutboundAdjacency[RIGHT], 1)
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)

	// Would complete a cycle via d=0 path from X -> V.
	checkEquality(t, V.OutboundAdjacency[RIGHT], RESTRICT_D1)
	// Would complete a cycle via d=1 path from X -> Z.
	checkEquality(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}

func TestEquality_BackBranch(t *testing.T) {
	chart, _, W, X, _, Z := buildFixture()

	// This excercises the second condition of equality, where
	// a "long-link" already exists and we're considering creating
	// the first link in back-path towards the long-link head.

	chart.Use(W.OutboundAdjacency[RIGHT], 0)
	chart.Use(X.OutboundAdjacency[LEFT], 0)
	chart.Use(W.OutboundAdjacency[RIGHT], 1)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	chart.Use(W.OutboundAdjacency[RIGHT], 1)

	checkEquality(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}
