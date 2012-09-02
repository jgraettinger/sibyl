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

	X := chart.nextCell()
	Y := chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 0)
	checkEquality(t, Y.OutboundAdjacency[LEFT], RESTRICT_D1)

	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 1)
	checkEquality(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
}

func TestEquality_LongLink(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	// This excercies the first condition of equality, where
	// the back link-path of a potential 'long' link already
	// exists, and constrains available depths of a 'long'
	// adjacency which would complete a cycle.

	V := chart.nextCell()
	W := chart.nextCell()
	chart.use(V.OutboundAdjacency[RIGHT], 0)
	chart.use(W.OutboundAdjacency[LEFT], 0)

	X := chart.nextCell()
	chart.use(X.OutboundAdjacency[LEFT], 0)

	// Would complete a cycle via d=0 path from X -> V.
	checkEquality(t, V.OutboundAdjacency[RIGHT], RESTRICT_D1)

	Y := chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 1)

	Z := chart.nextCell()
	chart.use(Y.OutboundAdjacency[RIGHT], 0)
	chart.use(Z.OutboundAdjacency[LEFT], 0)

	// Would complete a cycle via d=1 path from X -> Z.
	checkEquality(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}

func TestEqualityBackPathCompletion(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("T", "U", "V", "W", "X", "Y", "Z")

	// This excercises the second condition of equality, where
	// a long-link already exists and we're considering creating
	// the first link in back-path towards the long-link head.

	T, U := chart.nextCell(), chart.nextCell()
	chart.use(T.OutboundAdjacency[RIGHT], 0)
	chart.use(U.OutboundAdjacency[LEFT], 0)

	_ = chart.nextCell()
	chart.use(T.OutboundAdjacency[RIGHT], 1)

	W := chart.nextCell()
	chart.use(W.OutboundAdjacency[LEFT], 0)

	// Use of T => W restricts W's adjacency to U.
	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_NONE)
	chart.use(T.OutboundAdjacency[RIGHT], 1)
	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_D0)

	_ = chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	Y := chart.nextCell()
	// Leave W's adjacency to Y unused.

	Z := chart.nextCell()
	chart.use(Y.OutboundAdjacency[RIGHT], 0)
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	chart.use(Z.OutboundAdjacency[LEFT], 1)

	// Use of Z => W restricts W's adjacency to Y.
	checkEquality(t, W.OutboundAdjacency[RIGHT], RESTRICT_NONE)
	chart.use(Z.OutboundAdjacency[LEFT], 1)
	checkEquality(t, W.OutboundAdjacency[LEFT], RESTRICT_D0)
}
