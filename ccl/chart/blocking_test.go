package chart

import (
	"testing"
)

func checkBlocking(t *testing.T, adjacency *Adjacency, e DepthRestriction) {
	if r := adjacency.BlockingRestriction(); r != e {
		t.Errorf("Expected blocking %v on %v, but got %v",
			e, adjacency, r)
	}
}

func checkFullBlockBound(t *testing.T, cell *Cell, blockCell *Cell) {
	dir := DirectionFromPosition(blockCell.Index - cell.Index)
	if !dir.HasFullyBlockedAfter(cell) {
		t.Error("full blocking bound not set on ", cell)
	} else if dir.FullyBlockedAfter(cell) != blockCell.Index {
		t.Errorf("expected blocking of %v past %v (%v)", cell,
			blockCell, dir.FullyBlockedAfter(cell))
	}
}

func TestPartialBlocking(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("X", "Y", "Z")

	X, Y := chart.NextCell(), chart.NextCell()
	chart.Use(X.OutboundAdjacency[RIGHT], 0)
	chart.Use(Y.OutboundAdjacency[LEFT], 0)

	Z := chart.NextCell()
	// Flip linking order from that of the X,Y case.
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)

	checkBlocking(t, X.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkBlocking(t, Y.OutboundAdjacency[LEFT], RESTRICT_D0)
	checkBlocking(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}

func TestFullBlockingNeighbor(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("X", "Y", "Z")

	X, Y := chart.NextCell(), chart.NextCell()
	chart.Use(X.OutboundAdjacency[RIGHT], 1)
	checkFullBlockBound(t, Y, X)

	chart.Use(Y.OutboundAdjacency[LEFT], 1)
	checkFullBlockBound(t, X, Y)

	Z := chart.NextCell()
	chart.Use(Z.OutboundAdjacency[LEFT], 1)
	checkFullBlockBound(t, Y, Z)

	chart.Use(Y.OutboundAdjacency[RIGHT], 1)
	checkFullBlockBound(t, Z, Y)

	checkBlocking(t, X.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Y.OutboundAdjacency[LEFT], RESTRICT_ALL)
	checkBlocking(t, Y.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}

func TestFullBlockingProjectedForwardLink(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V, W := chart.NextCell(), chart.NextCell()
	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 0)

	X := chart.NextCell()
	chart.Use(W.OutboundAdjacency[RIGHT], 1)
	chart.Use(X.OutboundAdjacency[LEFT], 1) // Resolves violation
	checkFullBlockBound(t, V, W)

	Y := chart.NextCell()
	chart.Use(Y.OutboundAdjacency[LEFT], 1)
	chart.Use(X.OutboundAdjacency[RIGHT], 1) // Resolves violation

	Z := chart.NextCell()
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	checkFullBlockBound(t, Z, Y)

	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}

func TestFullBlockingProjectedBackwardLink(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V, W := chart.NextCell(), chart.NextCell()
	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 0)

	X := chart.NextCell()
	chart.Use(X.OutboundAdjacency[LEFT], 1)
	checkFullBlockBound(t, V, X)

	// A V => X link is still possible, but not beyond.
	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_D0)
	chart.Use(V.OutboundAdjacency[RIGHT], 1)

	Y := chart.NextCell()
	chart.Use(X.OutboundAdjacency[RIGHT], 1)

	Z := chart.NextCell()
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	checkFullBlockBound(t, Z, X)

	// A Z => X link is still possible, but not beyond.
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)

	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}
