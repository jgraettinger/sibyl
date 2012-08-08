package chart

import (
	assert "invariant"
	"testing"
)

// Use table-driven testing to verify expectations about
// graph adjacency, covering, and link path post-conditions.
type Expect struct {
	index int

	lAdj, rAdj           *Cell
	lUnusable, rUnusable bool

	lPathD0, lPathD1 *Cell
	rPathD0, rPathD1 *Cell
}

func checkExpectations(t *testing.T, chart *Chart, expectations []Expect) {
	checkLink := func(link *Link, expectedPath *Cell) {
		if link == nil {
			if expectedPath != nil {
				t.Errorf("Have nil link but expected path to %v", expectedPath)
			} else {
				// Both are nil
			}
		} else {
			if expectedPath == nil {
				t.Errorf("Expected nil link, but have %v", link)
			} else if *link.FurthestPath != expectedPath {
				t.Errorf("Expected path to %v, not %v to path %v",
					expectedPath, link, *link.FurthestPath)
			}
		}
	}

	checkAdjacency := func(adj *Adjacency, cell *Cell, unusable bool) {
		if adj.To != cell {
			t.Errorf("expected adjacency to %v, not %v", cell, adj)
		}
		if !adj.IsUsable() != unusable {
			t.Errorf("expected unusable %v", unusable)
		}
	}

	for _, e := range expectations {
		cell := chart.Cells[e.index]

		checkAdjacency(cell.OutboundAdjacency[LEFT], e.lAdj, e.lUnusable)
		checkAdjacency(cell.OutboundAdjacency[RIGHT], e.rAdj, e.rUnusable)

		checkLink(cell.LastOutboundLinkD0[LEFT], e.lPathD0)
		checkLink(cell.LastOutboundLinkD1[LEFT], e.lPathD1)
		checkLink(cell.LastOutboundLinkD0[RIGHT], e.rPathD0)
		checkLink(cell.LastOutboundLinkD1[RIGHT], e.rPathD1)
	}
}

func TestChart_AddCell(t *testing.T) {
	chart := NewChart()
	W, X, Y := 0, 1, 2

	chart.AddCell("W")
	assert.Equal(chart.Cells[W].OutboundAdjacency[LEFT].Position, -1)
	assert.IsNil(chart.Cells[W].OutboundAdjacency[LEFT].To)
	assert.Equal(chart.Cells[W].OutboundAdjacency[RIGHT].Position, 1)
	assert.IsNil(chart.Cells[W].OutboundAdjacency[RIGHT].To)
	assert.Equal(len(chart.Cells[W].InboundAdjacencies[LEFT]), 0)
	assert.Equal(len(chart.Cells[W].InboundAdjacencies[RIGHT]), 0)

	// Adding a new cell updates W's right-side adjacency.
	chart.AddCell("X")
	assert.Equal(chart.Cells[W].OutboundAdjacency[RIGHT].To, chart.Cells[X])
	assert.Equal(chart.Cells[X].OutboundAdjacency[LEFT].To, chart.Cells[W])
	assert.Equal(len(chart.Cells[W].InboundAdjacencies[RIGHT]), 1)
	assert.Equal(len(chart.Cells[X].InboundAdjacencies[LEFT]), 1)

	// Both adjacencies are free to be used.
	assert.IsFalse(chart.Cells[W].OutboundAdjacency[RIGHT].SpansPunctuation)
	assert.IsFalse(chart.Cells[X].OutboundAdjacency[LEFT].SpansPunctuation)

	// Link W -0> X creates a new right adjacency of w.
	chart.Use(chart.Cells[W].OutboundAdjacency[RIGHT], 0)

	// Punctionation immediately blocks adjacency to the next added cell.
	chart.StoppingPunctuation()

	chart.AddCell("Y")
	// Both X & Y have an adjacency to y.
	assert.Equal(chart.Cells[W].OutboundAdjacency[RIGHT].To, chart.Cells[Y])
	assert.Equal(chart.Cells[X].OutboundAdjacency[RIGHT].To, chart.Cells[Y])
	assert.Equal(len(chart.Cells[Y].InboundAdjacencies[LEFT]), 2)

	// But all adjacencies have stopping punctuation set.
	assert.IsTrue(chart.Cells[W].OutboundAdjacency[RIGHT].SpansPunctuation)
	assert.IsTrue(chart.Cells[X].OutboundAdjacency[RIGHT].SpansPunctuation)
	assert.IsTrue(chart.Cells[Y].OutboundAdjacency[LEFT].SpansPunctuation)
}

func TestChart_SimpleLinkPaths(t *testing.T) {
	chart := NewChart()
	W, X, Y, Z := 0, 1, 2, 3

	chart.AddCell("W")
	chart.AddCell("X")
	chart.Use(chart.Cells[W].OutboundAdjacency[RIGHT], 1)

	chart.AddCell("Y")
	chart.Use(chart.Cells[X].OutboundAdjacency[RIGHT], 0)
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 0)

	chart.AddCell("Z")
	chart.Use(chart.Cells[Z].OutboundAdjacency[LEFT], 1)

	c := chart.Cells
	expectations := []Expect{
		{index: W, rAdj: c[Z], rPathD1: c[Y]},
		{index: X, lAdj: c[W], rAdj: c[Z], rPathD0: c[Y]},
		{index: Y, lAdj: c[W], rAdj: c[Z], lPathD0: c[X]},
		{index: Z, lAdj: c[W], lPathD1: c[X]},
	}
	checkExpectations(t, chart, expectations)
}

func TestChart_LinkPathForward(t *testing.T) {
	chart := NewChart()
	V, W, X, Y, Z := 0, 1, 2, 3, 4

	chart.AddCell("V")
	chart.AddCell("W")
	chart.Use(chart.Cells[V].OutboundAdjacency[RIGHT], 0)

	chart.AddCell("X")
	chart.Use(chart.Cells[W].OutboundAdjacency[RIGHT], 0)

	chart.AddCell("Y")
	chart.Use(chart.Cells[X].OutboundAdjacency[RIGHT], 0)

	chart.AddCell("Z")
	chart.Use(chart.Cells[W].OutboundAdjacency[RIGHT], 1)

	c := chart.Cells
	expectations := []Expect{
		{index: V, rPathD0: c[Z]},
		{index: W, lAdj: c[V], rPathD0: c[Y], rPathD1: c[Z]},
		{index: X, lAdj: c[W], rAdj: c[Z], rUnusable: true, rPathD0: c[Y]},
		{index: Y, lAdj: c[X], rAdj: c[Z], rUnusable: true},
		{index: Z, lAdj: c[Y]},
	}
	checkExpectations(t, chart, expectations)
}

func TestChart_LinkPathBackward(t *testing.T) {
	chart := NewChart()
	V, W, X, Y, Z := 0, 1, 2, 3, 4

	chart.AddCell("V")
	chart.AddCell("W")
	chart.AddCell("X")
	chart.Use(chart.Cells[X].OutboundAdjacency[LEFT], 0)

	chart.AddCell("Y")
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 0)
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 1)

	chart.AddCell("Z")
	chart.Use(chart.Cells[Z].OutboundAdjacency[LEFT], 0)

	c := chart.Cells
	expectations := []Expect{
		{index: V, rAdj: c[W]},
		{index: W, lAdj: c[V], lUnusable: true, rAdj: c[X]},
		{index: X, lAdj: c[V], lUnusable: true, rAdj: c[Y], lPathD0: c[W]},
		{index: Y, rAdj: c[Z], lPathD0: c[W], lPathD1: c[V]},
		{index: Z, lPathD0: c[V]},
	}
	checkExpectations(t, chart, expectations)
}

func TestChart_UseExtended(t *testing.T) {
	chart := NewChart()
	T, U, V, W, X, Y, Z := 0, 1, 2, 3, 4, 5, 6

	// Manually build up a parse structure, which excercises
	// multiple link depths, covering links, and link paths.
	chart.AddCell("T")
	chart.AddCell("U")
	chart.Use(chart.Cells[T].OutboundAdjacency[RIGHT], 0)

	chart.AddCell("V")
	chart.Use(chart.Cells[U].OutboundAdjacency[RIGHT], 0)
	chart.Use(chart.Cells[V].OutboundAdjacency[LEFT], 0)

	chart.AddCell("W")
	chart.Use(chart.Cells[V].OutboundAdjacency[RIGHT], 0)

	chart.AddCell("X")
	chart.Use(chart.Cells[U].OutboundAdjacency[RIGHT], 1)

	chart.AddCell("Y")
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 0)
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 0)
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 0)
	chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 1)

	chart.AddCell("Z")
	chart.Use(chart.Cells[Y].OutboundAdjacency[RIGHT], 1)

	c := chart.Cells
	expectations := []Expect{
		{index: T, rAdj: c[Y], rPathD0: c[X]},
		{index: U, lAdj: c[T], lUnusable: true, rAdj: c[Y],
			rPathD0: c[W], rPathD1: c[X]},
		{index: V, lAdj: c[T], lUnusable: true, rAdj: c[X], rUnusable: true,
			lPathD0: c[U], rPathD0: c[W]},
		{index: W, lAdj: c[V], lUnusable: true, rAdj: c[X], rUnusable: true},
		{index: X, lAdj: c[W], lUnusable: true, rAdj: c[Y]},
		{index: Y, lPathD0: c[U], lPathD1: c[T], rPathD1: c[Z]},
		{index: Z, lAdj: c[Y]},
	}
	checkExpectations(t, chart, expectations)
}
