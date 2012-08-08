package chart

import (
	assert "invariant"
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

func TestBlocking_PartialOne(t *testing.T) {
	chart, _, _, X, Y, _ := buildFixture()

	chart.Use(X.OutboundAdjacency[RIGHT], 0)
	chart.Use(Y.OutboundAdjacency[LEFT], 0)
	assert.IsTrue(X.OutboundAdjacency[RIGHT].Blocking == BLOCK_D0)
	assert.IsTrue(Y.OutboundAdjacency[LEFT].Blocking == BLOCK_D0)
}
func TestBlocking_PartialTwo(t *testing.T) {
	chart, _, _, X, Y, _ := buildFixture()

	chart.Use(Y.OutboundAdjacency[LEFT], 0)
	chart.Use(X.OutboundAdjacency[RIGHT], 0)
	assert.IsTrue(X.OutboundAdjacency[RIGHT].Blocking == BLOCK_D0)
	assert.IsTrue(Y.OutboundAdjacency[LEFT].Blocking == BLOCK_D0)
}
func TestBlocking_FullNeighbor(t *testing.T) {
	chart, V, W, _, Y, Z := buildFixture()

	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 1)
	assert.IsTrue(V.OutboundAdjacency[RIGHT].Blocking == BLOCK_D0|BLOCK_ALL)

	chart.Use(Y.OutboundAdjacency[RIGHT], 1)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	assert.IsTrue(Z.OutboundAdjacency[LEFT].Blocking == BLOCK_D0|BLOCK_ALL)
}
func TestBlocking_FullProjected(t *testing.T) {
	chart, V, W, _, Y, Z := buildFixture()

	chart.Use(V.OutboundAdjacency[RIGHT], 0)
	chart.Use(W.OutboundAdjacency[LEFT], 0)
	chart.Use(W.OutboundAdjacency[RIGHT], 1)
	assert.IsTrue(V.OutboundAdjacency[RIGHT].Blocking == BLOCK_D0|BLOCK_ALL)

	chart.Use(Y.OutboundAdjacency[LEFT], 1)
	chart.Use(Y.OutboundAdjacency[RIGHT], 0)
	chart.Use(Z.OutboundAdjacency[LEFT], 0)
	assert.IsTrue(Z.OutboundAdjacency[LEFT].Blocking == BLOCK_D0|BLOCK_ALL)
}
