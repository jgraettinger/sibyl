package chart

import (
	assert "invariant"
	"testing"
)

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


func TestChart_SimpleLinkPaths(l *testing.T) {
    chart := NewChart()
    W, X, Y, Z = 0, 1, 2

    chart.AddCell("x")
    chart.AddCell("y")
    chart.Use(chart.Cells[X].OutboundAdjacency[RIGHT], 0)
    chart.Use(chart.Cells[Y].OutboundAdjacency[LEFT], 0)

    chart.AddCell("z")
    chart.Use(chart.Cells[Y].OutboundAdjacency[RIGHT], 0)
    chart.Use(chart.Cells[Z].OutboundAdjacency[LEFT], 0)

    assert.Equal(*chart.Cells[X].LastOutboundLinkD0[RIGHT].FurthestPath,
        chart.Cells[Z])
    assert.Equal(*chart.Cells[X].LastOutboundLinkD0[RIGHT].FurthestPath,
        chart.Cells[Z])

}


func TestChart_UseExtended(l *testing.T) {
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
	type Expect struct {
		index int

		leftAdjacency  *Cell
		leftBlocked    bool
		rightAdjacency *Cell
		rightBlocked   bool

		leftPathD0, leftPathD1   *Cell
		rightPathD0, rightPathD1 *Cell
	}

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

	for _, e := range expectations {
		cell := chart.Cells[e.index]

        checkAdjacency := func(adj *Adjacency, cell *Cell, blocked bool) {
            if adj.To != nil {
                l.Errorf("expected adjacency %v, not %v", cell, adj)
            }
            if adj.IsBlocked() != blocked {
                l.Errorf("expected blocking %v", blocked)
            }
        }

        checkAdjacency(cell.OutboundAdjacency[LEFT],
            e.leftAdjacency, e.leftBlocked)
        checkAdjacency(cell.OutboundAdjacency[RIGHT],
            e.rightAdjacency, e.rightBlocked)

		checkLink := func(cell *Cell, link *Link) {
			if cell == nil && link != nil {
				l.Errorf("non-nil link path: %v", link)
			} else if link == nil {
				l.Errorf("expected non-nil link from cell: %v", cell)
			} else {
				if *link.FurthestPath != cell {
					l.Errorf("expected path %v, not %v",
						cell, link.FurthestPath)
				}
			}
		}

        checkLink(e.leftPathD0, cell.LastOutboundLinkD0[LEFT])
        checkLink(e.leftPathD1, cell.LastOutboundLinkD1[LEFT])
        checkLink(e.rightPathD0, cell.LastOutboundLinkD0[RIGHT])
        checkLink(e.rightPathD1, cell.LastOutboundLinkD1[RIGHT])
	}
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
