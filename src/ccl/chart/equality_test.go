package chart

import (
	"testing"
)

func checkEquality(t *testing.T, adjacency *Adjacency, e DepthRestriction) {
	if r := adjacency.EqualityRestriction(); r != e {
		t.Errorf("Expected equality %v on %v, but got %v",
			e, adjacency, r)
	}
}

func TestEqualityNeighborCycle(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("X", "Y", "Z")

	X := chart.NextCell()
	Y := chart.NextCell()
	chart.Use(X.OutboundAdjacency[RIGHT], 0)
	checkEquality(t, Y.OutboundAdjacency[LEFT], RESTRICT_D1)

	Z := chart.NextCell()
	chart.Use(Z.OutboundAdjacency[LEFT], 1)
	checkEquality(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
}

func TestEquality_LongLink(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	// This excercies the first condition of equality, where
	// the back link-path of a potential 'long' link already
	// exists, and constrains available depths of a 'long'
	// adjacency which would complete a cycle.

	V := chart.NextCell()
	W := chart.NextCell()
	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 0)

	X := chart.NextCell()
	chart.Use(X.OutboundAdjacency[LEFT], 0)

	// Would complete a cycle via d=0 path from X -> V.
	checkEquality(t, V.OutboundAdjacency[RIGHT], RESTRICT_D1)

	Y := chart.NextCell()
	chart.Use(X.OutboundAdjacency[RIGHT], 1)

	Z := chart.NextCell()
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)

	// Would complete a cycle via d=1 path from X -> Z.
	checkEquality(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}

func TestEqualityBackPathCompletion(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("T", "U", "V", "W", "X", "Y", "Z")

	// This excercises the second condition of equality, where
	// a long-link already exists and we're considering creating
	// the first link in back-path towards the long-link head.

	T, U := chart.NextCell(), chart.NextCell()
	chart.Use(T.OutboundAdjacency[RIGHT], 0)
	chart.Use(U.OutboundAdjacency[LEFT], 0)

	_ = chart.NextCell()
	chart.Use(T.OutboundAdjacency[RIGHT], 1)

	W := chart.NextCell()
	chart.Use(W.OutboundAdjacency[LEFT], 0)

	// Use of T => W restricts W's adjacency to U.
	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_NONE)
	chart.Use(T.OutboundAdjacency[RIGHT], 1)
	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_D0)

	_ = chart.NextCell()
	chart.Use(W.OutboundAdjacency[RIGHT], 0)

	Y := chart.NextCell()
	// Leave W's adjacency to Y unused.

	Z := chart.NextCell()
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)

	// Use of Z => W restricts W's adjacency to Y.
	checkEquality(t, W.OutboundAdjacency[RIGHT], RESTRICT_NONE)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)
	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_D0)
}
