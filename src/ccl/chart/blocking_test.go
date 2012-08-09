package chart

import (
	"testing"
)

func buildFixture() (chart *Chart, V, W, X, Y, Z *Cell) {
	chart = NewChart()
	V = chart.AddCell("V")
	W = chart.AddCell("W")
	X = chart.AddCell("X")
	Y = chart.AddCell("Y")
	Z = chart.AddCell("Z")
	return
}

func checkBlocking(t *testing.T, adjacency *Adjacency,
    expect DepthRestriction) {

    if r := adjacency.BlockingRestriction(); r != expect {
        t.Errorf("Expected blocking %v on %v, but got %v",
            expect, adjacency, r)
    }
}

func TestBlocking_PartialOne(t *testing.T) {
	chart, _, _, X, Y, _ := buildFixture()

	chart.Use(X.OutboundAdjacency[RIGHT], 0)
	chart.Use(Y.OutboundAdjacency[LEFT], 0)
	checkBlocking(t, X.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkBlocking(t, Y.OutboundAdjacency[LEFT], RESTRICT_D0)
}
func TestBlocking_PartialTwo(t *testing.T) {
	chart, _, _, X, Y, _ := buildFixture()

	chart.Use(Y.OutboundAdjacency[LEFT], 0)
	chart.Use(X.OutboundAdjacency[RIGHT], 0)
	checkBlocking(t, X.OutboundAdjacency[RIGHT], RESTRICT_D0)
	checkBlocking(t, Y.OutboundAdjacency[LEFT], RESTRICT_D0)
}
func TestBlocking_FullNeighbor(t *testing.T) {
	chart, V, W, _, Y, Z := buildFixture()

	chart.Use(V.OutboundAdjacency[RIGHT], 1)
	chart.Use(W.OutboundAdjacency[LEFT], 1)
	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_ALL)

	chart.Use(Y.OutboundAdjacency[RIGHT], 1)
	chart.Use(Z.OutboundAdjacency[LEFT], 1)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}
func TestBlocking_FullProjected(t *testing.T) {
	chart, V, W, _, Y, Z := buildFixture()

	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 0)
	chart.Use(W.OutboundAdjacency[RIGHT], 1)
	checkBlocking(t, V.OutboundAdjacency[RIGHT], RESTRICT_ALL)

	chart.Use(Y.OutboundAdjacency[LEFT], 1)
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	checkBlocking(t, Z.OutboundAdjacency[LEFT], RESTRICT_ALL)
}
