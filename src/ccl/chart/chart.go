package chart

import (
	"invariant"
	"log"
)

const (
	LEFT  = 0
	RIGHT = 1
)

type Chart struct {
	nextTokens []Token
	Cells      []*Cell

	endInbound AdjacencySet

	minimalViolation      *Cell
	minimalViolationDepth uint
}

func NewChart() (chart *Chart) {
	chart = new(Chart)
	chart.endInbound = make(AdjacencySet)
	return chart
}

func (chart *Chart) AddTokens(tokens ...Token) {
	chart.nextTokens = append(chart.nextTokens, tokens...)
}

func (chart *Chart) nextCell() *Cell {
	var stoppingPunc bool
	var token Token

	// Pop the next token to process.
	for {
		if l := len(chart.nextTokens); l == 0 ||
			(l == 1 && chart.nextTokens[0].IsStoppingPunctuaction()) {
			return nil
		}

		token = chart.nextTokens[0]
		chart.nextTokens = chart.nextTokens[1:len(chart.nextTokens)]

		if token.IsStoppingPunctuaction() {
			stoppingPunc = true
		} else {
			break
		}
	}

	var prevCell *Cell
	if len(chart.Cells) > 0 {
		prevCell = chart.Cells[len(chart.Cells)-1]
	}

	nextCell := NewCell(len(chart.Cells), token)
	chart.Cells = append(chart.Cells, nextCell)

	// Update all current adjacencies to {end},
	// to instead be adjacent to nextCell.
	for adjacency := range chart.endInbound {
		delete(chart.endInbound, adjacency)

		if stoppingPunc {
			adjacency.SpansPunctuation = true
		}

		adjacency.To = nextCell
		nextCell.InboundAdjacencies[LEFT].Add(adjacency)
	}

	// Add nextCell => prevCell adjacency.
	adjacency := new(Adjacency)
	adjacency.From = nextCell
	adjacency.To = prevCell
	adjacency.Position = -1

	if stoppingPunc {
		adjacency.SpansPunctuation = true
	}

	nextCell.OutboundAdjacency[LEFT] = adjacency
	if prevCell != nil {
		prevCell.InboundAdjacencies[RIGHT].Add(adjacency)
	} else {
		// Adjacent to beginning-of-utterance
	}

	// Add nextCell => {end} adjacency.
	adjacency = new(Adjacency)
	adjacency.From = nextCell
	adjacency.Position = 1

	nextCell.OutboundAdjacency[RIGHT] = adjacency
	chart.endInbound.Add(adjacency)
	return nextCell
}

func (chart *Chart) CurrentCell() *Cell {
	if l := len(chart.Cells); l != 0 {
		return chart.Cells[l-1]
	}
	return nil
}

func (chart *Chart) use(adjacency *Adjacency, depth uint) {
	// Used adjacencies must be From or To the current (last) cell.
	{
		invariant.NotNil(adjacency.To)

		cell := chart.CurrentCell()
		invariant.IsTrue(adjacency.From == cell || adjacency.To == cell)
		invariant.IsTrue(chart.minimalViolation == nil ||
			adjacency.From == cell || adjacency.From == chart.minimalViolation)
	}

	forward := DirectionFromPosition(adjacency.Position)
	head, tail := adjacency.From, adjacency.To

	log.Printf("Using adjacency %v at depth %v, direction %v",
		adjacency, depth, forward)

	link := NewLink(adjacency, depth)

	// Update the link path to reflect this new link.
	if forward == LEFT_TO_RIGHT {
		invariant.IsNil(forward.OutboundLinks(tail).Last())

		if lastOut := forward.OutboundLinks(head).Last(); lastOut == nil {
			// As this is head's first outbound link in this direction,
			// using this adjacency doesn't create a new link-path.
			if chainIn := forward.InboundLink(head); chainIn != nil {
				// Extend the existing path.
				link.FurthestPath = chainIn.FurthestPath
			} else {
				// There is no existing path. Create a new one.
				link.FurthestPath = new(BoxedCellPointer)
			}
		} else {
			// Adding a second outbound link creates a new link-path. Because
			// parsing is left-to-right, we expect that additional links will
			// be added to this path. We want antecedent nodes to be updated
			// with further extensions, while preserving the paths rooted at
			// the existing outbound link.
			link.FurthestPath = lastOut.ForkFurthestPath()
		}
		// This update is visible from all previous links on the path.
		*link.FurthestPath = link.To
	} else {
		invariant.IsNil(forward.InboundLink(head))

		// Any successive links in this direction already exist. Propogate
		// an existing boxed path backwards to link.
		if chainOut := forward.OutboundLinks(tail).Last(); chainOut != nil {
			log.Printf("%v has last outbound link %v to %v", tail,
				chainOut, (*Cell)(*chainOut.FurthestPath))
			link.FurthestPath = chainOut.FurthestPath
		} else {
			link.FurthestPath = new(BoxedCellPointer)
			*link.FurthestPath = link.To
		}
	}

	log.Printf("Created link %v with path to %v", link,
		(*Cell)(*link.FurthestPath))

	// One beyond the link path from head, is the next adjacency position.
	nextAdjacencyIndex := forward.Increment((*link.FurthestPath).Index)
	log.Printf("nextAdjacencyIndex is %v", nextAdjacencyIndex)

	// If the adjacency from the furthest path spans punctation, any
	// adjacency to nextAdjacencyIndex in this direction must as well.
	spansPunctuation := forward.OutboundAdjacency(
		*link.FurthestPath).SpansPunctuation

	// 3.2.1 Connectedness: Use of this adjacency creates a new one,
	// spanning to exactly nextAdjacencyIndex (and no further).
	newAdjacency := new(Adjacency)
	newAdjacency.From = head
	newAdjacency.Position = forward.Increment(adjacency.Position)
	newAdjacency.SpansPunctuation = spansPunctuation

	// Replace head's outbound adjacency, and add the new outbound link.
	forward.OutboundLinks(head).Add(link)
	forward.SetOutboundAdjacency(head, newAdjacency)

	if link.Depth == 0 {
		forward.SetLastOutboundLinkD0(head, link)
	} else {
		forward.SetLastOutboundLinkD1(head, link)
	}

	// Replace tail's inbound adjacency, and set the inbound link.
	forward.InboundAdjacencies(tail).Remove(adjacency)
	chart.updateAdjacencyTo(newAdjacency, nextAdjacencyIndex)
	forward.SetInboundLink(tail, link)

	// Examine other inbound adjacencies into tail.
	for otherAdjacency := range forward.InboundAdjacencies(tail) {
		if !forward.Less(otherAdjacency.From.Index, head.Index) {
			// 3.2.2 Covering Links - Adjacencies which are fully covered by
			// adjacency are permanently blocked. No further attachments
			// are possible from cellFrom in the forward direction.
			log.Printf("Adjacency covered: %v", otherAdjacency)
			otherAdjacency.CoveredByLink = true
		}

		if otherAdjacency.IsMoveable() {
			// 3.2.1 Connectedness & 3.2.1 Minimality:
			// Moveable adjacencies must be shifted through the link-path
			// opened by adjacency. This ensures minimality (as the
			// adjacency would otherwise be redudant with adjacency), and
			// connectedness (as an adjacency is opened through to
			// nextAdjacencyIndex, and no further).
			log.Printf("Moving %v.To: %v => %v", otherAdjacency,
				tail.Index, nextAdjacencyIndex)

			otherAdjacency.SpansPunctuation = spansPunctuation ||
				otherAdjacency.SpansPunctuation

			forward.InboundAdjacencies(tail).Remove(otherAdjacency)
			chart.updateAdjacencyTo(otherAdjacency, nextAdjacencyIndex)
		}
	}
	updateBlocking(chart, link)
	updateResolution(chart, link)
	return
}

func (chart *Chart) updateAdjacencyTo(adjacency *Adjacency, toIndex int) {
	if toIndex == -1 {
		invariant.IsTrue(adjacency.Position < 0)
		adjacency.To = nil
	} else if toIndex == len(chart.Cells) {
		invariant.IsTrue(adjacency.Position > 0)
		adjacency.To = nil
		chart.endInbound.Add(adjacency)
	} else {
		adjacency.To = chart.Cells[toIndex]

		if adjacency.Position < 0 {
			chart.Cells[toIndex].InboundAdjacencies[RIGHT].Add(adjacency)
		} else {
			chart.Cells[toIndex].InboundAdjacencies[LEFT].Add(adjacency)
		}
	}
}
