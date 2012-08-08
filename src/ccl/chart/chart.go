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
	Cells []*Cell

	endInbound   AdjacencySet
	stoppingPunc bool
}

func NewChart() (chart *Chart) {
	chart = new(Chart)
	chart.endInbound = make(AdjacencySet)
	return chart
}

func (chart *Chart) StoppingPunctuation() {
	chart.stoppingPunc = true
}

func (chart *Chart) AddCell(token string) *Cell {
	invariant.IsFalse(IsStoppingPunctuaction(token))

	var prevCell, nextCell *Cell
	if len(chart.Cells) > 0 {
		prevCell = chart.Cells[len(chart.Cells)-1]
	}

	nextCell = NewCell(len(chart.Cells), token)
	chart.Cells = append(chart.Cells, nextCell)

	// Update all current adjacencies to {end},
	// to instead be adjacent to nextCell.
	for adjacency := range chart.endInbound {
		delete(chart.endInbound, adjacency)

		if chart.stoppingPunc /* && adjacency.From == prevCell */ {
			// This direct adjacency is prohibited by punctuation.
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

	if chart.stoppingPunc {
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

	chart.stoppingPunc = false
	return nextCell
}

func (chart *Chart) Use(usedAdjacency *Adjacency, usedDepth uint) {
	invariant.IsTrue(usedAdjacency.IsUsable())

	forward := DirectionFromPosition(usedAdjacency.Position)
	cellFrom, cellTo := usedAdjacency.From, usedAdjacency.To

	log.Printf("Using adjacency %v at depth %v, direction %v",
		*usedAdjacency, usedDepth, forward)

	newLink := NewLink(usedAdjacency, usedDepth)

	// Update the link path to reflect this new link.
	if forward == LEFT_TO_RIGHT {
        invariant.IsNil(forward.OutboundLinks(cellTo).Last())

		if lastOut := forward.OutboundLinks(cellFrom).Last(); lastOut == nil {
			// As this is cellFrom's first outbound link in this direction,
			// using this adjacency doesn't create a new link-path.
			if chainIn := forward.InboundLink(cellFrom); chainIn != nil {
				// Extend the existing path.
				newLink.FurthestPath = chainIn.FurthestPath
			} else {
				// There is no existing path. Create a new one.
				newLink.FurthestPath = new(BoxedCellPointer)
			}
		} else {
			// Adding a second outbound link creates a new link-path. Because
			// parsing is left-to-right, we expect that additional links will
			// be added to this path. We want antecedent nodes to be updated
			// with further extensions, while preserving the paths rooted at
			// the existing outbound link.
			newLink.FurthestPath = lastOut.ForkFurthestPath()
		}
		// This update is visible from all previous links on the path.
		*newLink.FurthestPath = newLink.To
	} else {
        invariant.IsNil(forward.InboundLink(cellFrom))

		// Any successive links in this direction already exist. Propogate
		// an existing boxed path back to newLink.
		if chainOut := forward.OutboundLinks(cellTo).Last(); chainOut != nil {
			log.Printf("%v has last outbound link %v to %v", cellTo,
				chainOut, (*Cell)(*chainOut.FurthestPath))
			newLink.FurthestPath = chainOut.FurthestPath
		} else {
			newLink.FurthestPath = new(BoxedCellPointer)
			*newLink.FurthestPath = newLink.To
		}
	}

	log.Printf("Created link %v with path to %v", newLink,
		(*Cell)(*newLink.FurthestPath))

	// One beyond the link path from cellFrom, is the next adjacency position.
	nextAdjacencyIndex := forward.Increment((*newLink.FurthestPath).Index)
	log.Printf("nextAdjacencyIndex is %v", nextAdjacencyIndex)

	// 3.2.1 Connectedness: Use of this adjacency creates a new one,
	// spanning to exactly nextAdjacencyIndex (and no further).
	newAdjacency := new(Adjacency)
	newAdjacency.From = cellFrom
	newAdjacency.Position = forward.Increment(usedAdjacency.Position)

	// Replace cellFrom's outbound adjacency, and add the new outbound link.
	forward.OutboundLinks(cellFrom).Add(newLink)
	forward.SetOutboundAdjacency(cellFrom, newAdjacency)

	if newLink.Depth == 0 {
		forward.SetLastOutboundLinkD0(cellFrom, newLink)
	} else {
		forward.SetLastOutboundLinkD1(cellFrom, newLink)

		// 3.2.1 Montonicity: New adjacencies must have depth >= previous ones.
		newAdjacency.MontonicityRestricted = true
	}

	// Replace cellTo's inbound adjacency, and set the inbound link.
	forward.InboundAdjacencies(cellTo).Remove(usedAdjacency)
	chart.updateAdjacencyTo(newAdjacency, nextAdjacencyIndex)
	forward.SetInboundLink(cellTo, newLink)

	// Examine other inbound adjacencies into cellTo.
	for adjacency := range forward.InboundAdjacencies(cellTo) {
		if !forward.Less(adjacency.From.Index, cellFrom.Index) {
			// 3.2.2 Covering Links - Adjacencies which are fully covered by 
			// usedAdjacency are permanently blocked. No further attachments
			// are possible from cellFrom in the forward direction.
			log.Printf("Adjacency covered: %v", adjacency)
			adjacency.CoveredByLink = true
		}

		if adjacency.IsMoveable() {
			// 3.2.1 Connectedness & 3.2.1 Minimality:
			// Moveable adjacencies must be shifted through the link-path
			// opened by usedAdjacency. This ensures minimality (as the
			// adjacency would otherwise be redudant with usedAdjacency), and
			// connectedness (as an adjacency is opened through to
			// nextAdjacencyIndex, and no further).
			log.Printf("Moving %v.To: %v => %v", adjacency,
				cellTo.Index, nextAdjacencyIndex)

			forward.InboundAdjacencies(cellTo).Remove(adjacency)
			chart.updateAdjacencyTo(adjacency, nextAdjacencyIndex)
		}
	}

	updateBlocking(chart, newLink, newAdjacency)
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
