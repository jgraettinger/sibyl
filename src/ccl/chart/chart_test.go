package chart

import (
	assert "invariant"
	"testing"
)

type Expect struct {
	index int

	leftAdjacency  *Cell
	leftBlocked    bool
	rightAdjacency *Cell
	rightBlocked   bool

	leftPathD0, leftPathD1   *Cell
	rightPathD0, rightPathD1 *Cell
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

	checkAdjacency := func(adj *Adjacency, cell *Cell, blocked bool) {
		if adj.To != cell {
			t.Errorf("expected adjacency to %v, not %v", cell, adj)
		}
		if adj.IsBlocked() != blocked {
			t.Errorf("expected blocking %v", blocked)
		}
	}

	for _, e := range expectations {
		cell := chart.Cells[e.index]

		checkAdjacency(cell.OutboundAdjacency[LEFT],
			e.leftAdjacency, e.leftBlocked)
		checkAdjacency(cell.OutboundAdjacency[RIGHT],
			e.rightAdjacency, e.rightBlocked)

		checkLink(cell.LastOutboundLinkD0[LEFT], e.leftPathD0)
		checkLink(cell.LastOutboundLinkD1[LEFT], e.leftPathD1)
		checkLink(cell.LastOutboundLinkD0[RIGHT], e.rightPathD0)
		checkLink(cell.LastOutboundLinkD1[RIGHT], e.rightPathD1)
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
	assert.IsFalse(chart.Cells[W].OutboundAdjacency[RIGHT].IsBlocked())
	assert.IsFalse(chart.Cells[X].OutboundAdjacency[LEFT].IsBlocked())

	// Link W -0> X creates a new right adjacency of w.
	chart.Use(chart.Cells[W].OutboundAdjacency[RIGHT], 0)

	// Punctionation immediately blocks adjacency to the next added cell.
	chart.StoppingPunctuation()

	chart.AddCell("Y")
	// Both X & Y have an adjacency to y.
	assert.Equal(chart.Cells[W].OutboundAdjacency[RIGHT].To, chart.Cells[Y])
	assert.Equal(chart.Cells[X].OutboundAdjacency[RIGHT].To, chart.Cells[Y])
	assert.Equal(len(chart.Cells[Y].InboundAdjacencies[LEFT]), 2)

	// But all adjacencies are blocked.
	assert.IsTrue(chart.Cells[W].OutboundAdjacency[RIGHT].IsBlocked())
	assert.IsTrue(chart.Cells[X].OutboundAdjacency[RIGHT].IsBlocked())
	assert.IsTrue(chart.Cells[Y].OutboundAdjacency[LEFT].IsBlocked())
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
		{W, nil, false, c[Z], false, nil, nil, nil, c[Y]},
		{X, c[W], false, c[Z], false, nil, nil, c[Y], nil},
		{Y, c[W], false, c[Z], false, c[X], nil, nil, nil},
		{Z, c[W], false, nil, false, nil, c[X], nil, nil},
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
		{V, nil, false, nil, false, nil, nil, c[Z], nil},
		{W, c[V], false, nil, false, nil, nil, c[Y], c[Z]},
		{X, c[W], false, c[Z], true, nil, nil, c[Y], nil},
		{Y, c[X], false, c[Z], true, nil, nil, nil, nil},
		{Z, c[Y], false, nil, false, nil, nil, nil, nil},
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
		{V, nil, false, c[W], false, nil, nil, nil, nil},
		{W, c[V], true, c[X], false, nil, nil, nil, nil},
		{X, c[V], true, c[Y], false, c[W], nil, nil, nil},
		{Y, nil, false, c[Z], false, c[W], c[V], nil, nil},
		{Z, nil, false, nil, false, c[V], nil, nil, nil},
	}
	checkExpectations(t, chart, expectations)
}

func TestChart_Montonicity(t *testing.T) {
	chart := NewChart()
	Y, Z := 0, 1

	chart.AddCell("Y")
	chart.AddCell("Z")
	chart.Use(chart.Cells[Y].OutboundAdjacency[RIGHT], 1)
	chart.Use(chart.Cells[Z].OutboundAdjacency[LEFT], 1)

	// Outbound adjacencies created through linking have depth 0 blocked.
	assert.IsTrue(chart.Cells[Y].OutboundAdjacency[RIGHT].BlockedDepths[0])
	assert.IsTrue(chart.Cells[Z].OutboundAdjacency[LEFT].BlockedDepths[0])
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

	// Use table-driven testing to verify expectations about
	// remaining graph adjacencies and link paths.
	c := chart.Cells
	expectations := []Expect{
		{T, nil, false, c[Y], false, nil, nil, c[X], nil},
		{U, c[T], true, c[Y], false, nil, nil, c[W], c[X]},
		{V, c[T], true, c[X], true, c[U], nil, c[W], nil},
		{W, c[V], true, c[X], true, nil, nil, nil, nil},
		{X, c[W], true, c[Y], false, nil, nil, nil, nil},
		{Y, nil, false, nil, false, c[U], c[T], nil, c[Z]},
		{Z, c[Y], false, nil, false, nil, nil, nil, nil},
	}
	checkExpectations(t, chart, expectations)

}

/*
func assertBlocked(t *testing.T, adjacency *Adjacency) {
	if !adjacency.IsBlocked() {
		t.Errorf("expected %v to be blocked", adjacency)
	}
}
func assertUsed(t *testing.T, adjacency *Adjacency) {
	if !adjacency.Used {
		t.Errorf("expected %v to be used", adjacency)
	}
}
func assertNotUsed(t *testing.T, adjacency *Adjacency) {
	if adjacency.Used {
		t.Errorf("expected %v to be unused", adjacency)
	}
}
func assertNil(t *testing.T, adjacency *Adjacency) {
	if adjacency != nil {
		t.Errorf("didn't expect %v to exist", adjacency)
	}
}
func assertPaths(t *testing.T, cell *Cell, dir Direction,
	pathBegin, pathEndD0, pathEndD1 int) {

	if cell.PathBegin[dir.Side()] != pathBegin {
		t.Errorf("PathBegin mismatch: %v vs %v (%v)",
			pathBegin, cell.PathBegin[dir.Side()], cell)
	}
	if cell.PathEndD0[dir.Side()] != pathEndD0 {
		t.Errorf("PathEndD0 mismatch: %v vs %v (%v)",
			pathEndD0, cell.PathEndD0[dir.Side()], cell)
	}
	if cell.PathEndD1[dir.Side()] != pathEndD1 {
		t.Errorf("PathEndD1 mismatch: %v vs %v (%v)",
			pathEndD1, cell.PathEndD1[dir.Side()], cell)
	}
}

func buildChart(utterance string) (chart *Chart) {
	chart = NewChart()
	for _, token := range strings.Split(utterance, " ") {
		chart.AddCell(token)
	}
	return
}

func adj(chart *Chart, from, to int) *Adjacency {
	forward := DirectionFromPosition(to - from)

	for adjacency := range forward.Inbound(chart.Cells[to]) {
		if adjacency.From.Index == from {
			return adjacency
		}
	}
	return nil
}
*/
