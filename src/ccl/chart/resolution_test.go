package chart

import (
	"testing"
)

func checkMinimalViolation(t *testing.T, chart *Chart, c *Cell, d uint) {
	if chart.minimalViolation != c {
		t.Errorf("expected minimalViolation %v not %v",
			c, chart.minimalViolation)
	}
	if chart.minimalViolation != nil && chart.minimalViolationDepth != d {
		t.Error("expected violation depth mismatch")
	}
}

func checkResolution(t *testing.T, chart *Chart,
	adjacency *Adjacency, e DepthRestriction) {
	if r := adjacency.ResolutionRestriction(chart); r != e {
		t.Errorf("expected resolution %v on %v, but got %v",
			e, adjacency, r)
	}
}

func TestResolutionResolvedWithBackLinks(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V, W := chart.nextCell(), chart.nextCell()
	chart.use(V.OutboundAdjacency[RIGHT], 0)

	_ = chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	_ = chart.nextCell()
	chart.use(V.OutboundAdjacency[RIGHT], 1)

	Z := chart.nextCell()
	checkMinimalViolation(t, chart, nil, 0)

	// Creates a violation with minimal violation V.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkMinimalViolation(t, chart, V, 1)
	// Resolution restricts Z's left adjacency to 0, and V's to 1.
	checkResolution(t, chart, Z.OutboundAdjacency[LEFT], RESTRICT_D1)
	checkResolution(t, chart, V.OutboundAdjacency[RIGHT], RESTRICT_D0)

	// Moves the minimal violation to W, 0.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkMinimalViolation(t, chart, W, 0)
	checkResolution(t, chart, Z.OutboundAdjacency[LEFT], RESTRICT_D1)

	// Moves the minimal violation to V, 0.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkMinimalViolation(t, chart, V, 0)
	// Resolution now restricts V's adjacency to 1.
	checkResolution(t, chart, V.OutboundAdjacency[RIGHT], RESTRICT_D1)

	// Resolves the resolution violation.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkMinimalViolation(t, chart, nil, 0)
}

func TestResolutionResolvedWithForwardLinks(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("W", "X", "Y", "Z")

	W, _ := chart.nextCell(), chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	_ = chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	Z := chart.nextCell()
	checkMinimalViolation(t, chart, nil, 0)

	// Creates a violation with minimal resolution W.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkMinimalViolation(t, chart, W, 0)
	checkResolution(t, chart, W.OutboundAdjacency[RIGHT], RESTRICT_D1)
	checkResolution(t, chart, Z.OutboundAdjacency[LEFT], RESTRICT_D1)

	// Use of Z => W doesn't change the violation to W.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	checkMinimalViolation(t, chart, W, 0)

	// Use of W => Z resolves.
	chart.use(W.OutboundAdjacency[RIGHT], 0)
	checkMinimalViolation(t, chart, nil, 0)
}
