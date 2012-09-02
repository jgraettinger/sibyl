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

	X, Y := chart.nextCell(), chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 0)
	chart.use(Y.OutboundAdjacency[LEFT], 0)

	Z := chart.nextCell()
	// Flip linking order from that of the X,Y case.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	chart.use(Y.OutboundAdjacency[RIGHT], 0)

	checkBlocking(t, X.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkBlocking(t, Y.OutboundAdjacency[LEFT], RESTRICT_D0)
	checkBlocking(t, Y.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
}

func TestFullBlockingNeighbor(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("X", "Y", "Z")

	X, Y := chart.nextCell(), chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 1)
	checkFullBlockBound(t, Y, X)

	chart.use(Y.OutboundAdjacency[LEFT], 1)
	checkFullBlockBound(t, X, Y)

	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 1)
	checkFullBlockBound(t, Y, Z)

	chart.use(Y.OutboundAdjacency[RIGHT], 1)
	checkFullBlockBound(t, Z, Y)

	checkBlocking(t, X.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Y.OutboundAdjacency[LEFT], RESTRICT_ALL)
	checkBlocking(t, Y.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}

func TestFullBlockingProjectedForwardLink(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V, W := chart.nextCell(), chart.nextCell()
	chart.use(V.OutboundAdjacency[RIGHT], 0)
	chart.use(W.OutboundAdjacency[LEFT], 0)

	_ = chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 1)
	checkFullBlockBound(t, V, W)

	Y := chart.nextCell()
	chart.use(Y.OutboundAdjacency[LEFT], 1)

	Z := chart.nextCell()
	chart.use(Y.OutboundAdjacency[RIGHT], 0)
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkFullBlockBound(t, Z, Y)

	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}

func TestFullBlockingProjectedBackwardLink(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V, W := chart.nextCell(), chart.nextCell()
	chart.use(V.OutboundAdjacency[RIGHT], 0)
	chart.use(W.OutboundAdjacency[LEFT], 0)

	X := chart.nextCell()
	chart.use(X.OutboundAdjacency[LEFT], 1)
	checkFullBlockBound(t, V, X)

	// A V => X link is still possible, but not beyond.
	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_D0)
	chart.use(V.OutboundAdjacency[RIGHT], 1)

	Y := chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 1)

	Z := chart.nextCell()
	chart.use(Y.OutboundAdjacency[RIGHT], 0)
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkFullBlockBound(t, Z, X)

	// A Z => X link is still possible, but not beyond.
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_D0)
	chart.use(Z.OutboundAdjacency[LEFT], 1)

	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_ALL)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}
