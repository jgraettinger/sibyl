package chart

import (
	"fmt"
	"testing"
)

// Use table-driven testing to verify expectations about
// graph adjacency, covering, and link path post-conditions.
type expect struct {
	c *Cell

	// Must be declared in order of increasing absolute position.
	lLinkOut, rLinkOut expectLinks
	lAdjOut, rAdjOut   *Cell

	// Link-path extents implicit from lLinkOut & rLinkOut
	// don't need to be declared. Longer paths do.
	lPathD0, lPathD1 *Cell
	rPathD0, rPathD1 *Cell

	// These are filled in implicitly by verify(), based
	// on outgoing links & adjacencies of other cells.
	linkIn        [2]*Cell
	adjacenciesIn [2]map[*Cell]bool
}

type expectLink struct {
	depth uint
	cell  *Cell
}
type expectLinks []expectLink

func verify(t *testing.T, expectations ...expect) {

	expectMap := make(map[*Cell]*expect)
	for ind := range expectations {
		expectMap[expectations[ind].c] = &expectations[ind]
		expectations[ind].adjacenciesIn[0] = make(map[*Cell]bool)
		expectations[ind].adjacenciesIn[1] = make(map[*Cell]bool)
	}

	// Pass 1: Verify declared outbound links, adjacencies, and
	// link-paths. Collect implicit inbound links & adjacencies.
	checkSideOut := func(cell *Cell, dir Direction,
		eLinks expectLinks, eAdjacency, ePathD0, ePathD1 *Cell) {

		linkCount := len(*dir.OutboundLinks(cell))
		if len(eLinks) != linkCount {
			t.Error("cell link arity mismatch ", cell, eLinks)
			return
		}

		for i, eLink := range eLinks {
			link := (*dir.OutboundLinks(cell))[i]

			if link.Position != (i+1)*dir.PositionSign() ||
				link.Depth != eLink.depth ||
				link.To != eLink.cell {
				t.Error("link mismatch ", link, eLink)
			}
			// Mark the expected inbound link.
			expectMap[link.To].linkIn[dir.SideIn()] = cell
		}

		adjacency := dir.OutboundAdjacency(cell)
		if adjacency.Position != (1+linkCount)*dir.PositionSign() ||
			adjacency.To != eAdjacency {
			t.Error("adjacency mismatch ", adjacency, eAdjacency)
		}
		if adjacency.To != nil {
			if _, ok := expectMap[adjacency.To]; !ok {
				t.Error("saw adjacency to cell w/o expect{} ", adjacency)
				return
			}
			// Mark the expected inbound adjacency.
			expectMap[adjacency.To].adjacenciesIn[dir.SideIn()][cell] = true
		}

		// Don't require trivial link-paths to be declared.
		if l := dir.LastOutboundLinkD0(cell); ePathD0 == nil && l == nil {
		} else if ePathD0 == nil && *l.FurthestPath == l.To {
		} else if ePathD0 != nil && l != nil && *l.FurthestPath == ePathD0 {
		} else {
			t.Error("depth=0 link path mismatch ", ePathD0, l)
		}

		if l := dir.LastOutboundLinkD1(cell); ePathD1 == nil && l == nil {
		} else if ePathD1 == nil && *l.FurthestPath == l.To {
		} else if ePathD1 != nil && l != nil && *l.FurthestPath == ePathD1 {
		} else {
			t.Error("depth=1 link path mismatch ", ePathD1, l)
		}
	}

	for cell, eCell := range expectMap {
		checkSideOut(cell, LEFT_TO_RIGHT,
			eCell.rLinkOut, eCell.rAdjOut, eCell.rPathD0, eCell.rPathD1)
		checkSideOut(cell, RIGHT_TO_LEFT,
			eCell.lLinkOut, eCell.lAdjOut, eCell.lPathD0, eCell.lPathD1)
	}

	// Pass 2: Verify inbound links & adjacencies implied by expectations.
	checkSideIn := func(eCell *expect, dir Direction) {
		eLink := eCell.linkIn[dir.SideIn()]
		if link := dir.InboundLink(eCell.c); eLink == nil && link == nil {
		} else if eLink != nil && link != nil && link.From == eLink {
		} else {
			t.Error("inbound link mismatch ", link, eLink)
		}

		eAdjacencies := eCell.adjacenciesIn[dir.SideIn()]
		for adjacency := range dir.InboundAdjacencies(eCell.c) {
			if !eAdjacencies[adjacency.From] {
				t.Error("unexpected inbound adjacency ", eCell, adjacency)
			}
			delete(eAdjacencies, adjacency.From)
		}
		if len(eAdjacencies) != 0 {
			t.Error("remaining unobserved inbound adjacencies ", eAdjacencies)
		}
	}

	for _, eCell := range expectMap {
		checkSideIn(eCell, LEFT_TO_RIGHT)
		checkSideIn(eCell, RIGHT_TO_LEFT)
	}
}

func TestChart_AdjacencyUpdate(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("W", "X", "Y", "Z")

	// With one cell, W's adjacent to {begin} and {end}.
	W := chart.nextCell()
	verify(t, expect{c: W})

	// Adding a new cell updates W's right-side adjacency.
	X := chart.nextCell()
	verify(t, expect{c: W, rAdjOut: X}, expect{c: X, lAdjOut: W})

	// Use W & X's adjacencies, creating new ones.
	chart.use(W.OutboundAdjacency[RIGHT], 0)
	if a := W.OutboundAdjacency[RIGHT]; a.Position != 2 || a.To != nil {
		t.Error("Expected a new adjacency to chart {end}")
	}
	chart.use(X.OutboundAdjacency[LEFT], 0)
	if a := X.OutboundAdjacency[LEFT]; a.Position != -2 || a.To != nil {
		t.Error("Expected a new adjacency to chart {begin}")
	}

	// Adding Y updates X & W's {end} adjacencies.
	Y := chart.nextCell()
	verify(t,
		expect{c: W, rLinkOut: expectLinks{{0, X}}, rAdjOut: Y},
		expect{c: X, lLinkOut: expectLinks{{0, W}}, rAdjOut: Y},
		expect{c: Y, lAdjOut: X})

	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 0)

	// Connectedness implies adding Z doesn't affect X & W's adjacencies.
	verify(t,
		expect{c: W, rLinkOut: expectLinks{{0, X}}, rAdjOut: Y},
		expect{c: X, lLinkOut: expectLinks{{0, W}}, rAdjOut: Y},
		expect{c: Y, lAdjOut: X, rAdjOut: Z},
		expect{c: Z, lLinkOut: expectLinks{{0, Y}}, lAdjOut: X})

	if chart.nextCell() != nil {
		t.Error("Should have run out of input")
	}
}

func TestChart_StoppingPunctuation(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("W", "X", ";", "Y", "Z")

	W := chart.nextCell()
	X := chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	Y := chart.nextCell()
	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 0)

	// Punctuation is skipped when creating cells.
	if chart.nextCell() != nil {
		t.Error("Should have run out of input")
	}

	// Adjacencies are across stopping punctiation.
	verify(t,
		expect{c: W, rLinkOut: expectLinks{{0, X}}, rAdjOut: Y},
		expect{c: X, lAdjOut: W, rAdjOut: Y},
		expect{c: Y, lAdjOut: X, rAdjOut: Z},
		expect{c: Z, lLinkOut: expectLinks{{0, Y}}, lAdjOut: X})

	fmt.Print(chart)
	// However, adjacencies spanning stopping punctuation are marked as such.
	if !W.OutboundAdjacency[RIGHT].SpansPunctuation ||
		!X.OutboundAdjacency[RIGHT].SpansPunctuation ||
		!Y.OutboundAdjacency[LEFT].SpansPunctuation ||
		!Z.OutboundAdjacency[LEFT].SpansPunctuation {
		t.Error("expected adjacencies to span punctuation")
	}
}

func TestChart_SimpleLinkPaths(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("W", "X", "Y", "Z")

	W := chart.nextCell()
	X := chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 1)

	Y := chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 0)
	chart.use(Y.OutboundAdjacency[LEFT], 0)

	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 1)

	verify(t,
		expect{c: W, rLinkOut: expectLinks{{1, X}}, rAdjOut: Z, rPathD1: Y},
		expect{c: X, rLinkOut: expectLinks{{0, Y}}, rAdjOut: Z, lAdjOut: W},
		expect{c: Y, lLinkOut: expectLinks{{0, X}}, lAdjOut: W, rAdjOut: Z},
		expect{c: Z, lLinkOut: expectLinks{{1, Y}}, lAdjOut: W, lPathD1: X})
}

func TestChart_BranchingLinkPathForward(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V := chart.nextCell()
	W := chart.nextCell()
	chart.use(V.OutboundAdjacency[RIGHT], 0)

	X := chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	Y := chart.nextCell()
	chart.use(X.OutboundAdjacency[RIGHT], 0)

	Z := chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 1)

	verify(t,
		expect{c: V, rLinkOut: expectLinks{{0, W}}, rPathD0: Z},
		expect{c: W, rLinkOut: expectLinks{{0, X}, {1, Z}},
			rPathD0: Y, rPathD1: Z, lAdjOut: V},
		expect{c: X, rLinkOut: expectLinks{{0, Y}}, lAdjOut: W, rAdjOut: Z},
		expect{c: Y, lAdjOut: X, rAdjOut: Z},
		expect{c: Z, lAdjOut: Y})
}

func TestChart_BranchingLinkPathBackward(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("V", "W", "X", "Y", "Z")

	V := chart.nextCell()
	W := chart.nextCell()
	X := chart.nextCell()
	chart.use(X.OutboundAdjacency[LEFT], 0)

	Y := chart.nextCell()
	chart.use(Y.OutboundAdjacency[LEFT], 0)
	chart.use(Y.OutboundAdjacency[LEFT], 1)

	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 0)

	verify(t,
		expect{c: V, rAdjOut: W},
		expect{c: W, lAdjOut: V, rAdjOut: X},
		expect{c: X, lLinkOut: expectLinks{{0, W}}, lAdjOut: V, rAdjOut: Y},
		expect{c: Y, lLinkOut: expectLinks{{0, X}, {1, V}},
			lPathD0: W, lPathD1: V, rAdjOut: Z},
		expect{c: Z, lLinkOut: expectLinks{{0, Y}}, lPathD0: V})
}

func TestChart_CoveredAdjacencies(t *testing.T) {
	chart := NewChart()
	chart.AddTokens("W", "X", "Y", "Z")

	W := chart.nextCell()
	X := chart.nextCell()
	chart.use(W.OutboundAdjacency[RIGHT], 0)

	Y := chart.nextCell()
	if a := X.OutboundAdjacency[RIGHT]; a.To != Y || a.CoveredByLink {
		t.Error("expected non-covered adjacency X => Y ", a)
	}

	// Use of W => Y covers the X => Y adjacency. It isn't moved.
	chart.use(W.OutboundAdjacency[RIGHT], 0)
	if a := X.OutboundAdjacency[RIGHT]; a.To != Y || !a.CoveredByLink {
		t.Error("expected covered adjacency X => Y ", a)
	}

	Z := chart.nextCell()
	chart.use(Z.OutboundAdjacency[LEFT], 0)

	if a := Y.OutboundAdjacency[LEFT]; a.To != X || a.CoveredByLink {
		t.Error("expected non-covered adjacency Y => X", a)
	}

	// Use of Z => X covers the Y => X adjacency. It isn't moved.
	chart.use(Z.OutboundAdjacency[LEFT], 0)
	if a := Y.OutboundAdjacency[LEFT]; a.To != X || !a.CoveredByLink {
		t.Error("expected non-covered adjacency Y => X", a)
	}

	verify(t,
		expect{c: W, rLinkOut: expectLinks{{0, X}, {0, Y}}, rAdjOut: Z},
		expect{c: X, lAdjOut: W, rAdjOut: Y},
		expect{c: Y, lAdjOut: X, rAdjOut: Z},
		expect{c: Z, lLinkOut: expectLinks{{0, Y}, {0, X}}, lAdjOut: W})
}
