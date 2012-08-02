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

func (chart *Chart) AddCell(token string) {
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
			adjacency.BlockedDepths = [2]bool{true, true}
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
		adjacency.BlockedDepths = [2]bool{true, true}
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
}

func (chart *Chart) Use(usedAdjacency *Adjacency, targetDepth uint) {
	invariant.IsFalse(usedAdjacency.BlockedDepths[targetDepth])

	forward := DirectionFromPosition(usedAdjacency.Position)
	cellFrom, cellTo := usedAdjacency.From, usedAdjacency.To

	log.Printf("Using adjacency %v at depth %v, direction %v",
		*usedAdjacency, targetDepth, forward)

	createdLink := NewLink(usedAdjacency, targetDepth)

	// Update the link path to reflect this new link.
	if forward == LEFT_TO_RIGHT {
		if lastOut := forward.OutboundLinks(cellFrom).Last(); lastOut == nil {
			// As this is cellFrom's first outbound link in this direction,
			// using this adjacency doesn't create a new link-path.
			if chainIn := forward.InboundLink(cellFrom); chainIn != nil {
				// Extend the existing path.
				createdLink.FurthestPath = chainIn.FurthestPath
			} else {
				// There is no existing path. Create a new one.
				createdLink.FurthestPath = new(BoxedCellPointer)
			}
		} else {
			// Adding a second outbound link creates a new link-path. Because
			// parsing is left-to-right, we expect that additional links will
			// be added to this path. We want antecedent nodes to be updated
			// with further extensions, while preserving the paths rooted at
			// the existing outbound link.
			createdLink.FurthestPath = lastOut.ForkFurthestPath()
		}
		// This update is visible from all previous links on the path.
		*createdLink.FurthestPath = createdLink.To
	} else {
		// Any successive links in this direction already exist. Propogate
		// an existing boxed path back to createdLink.
		if chainOut := forward.OutboundLinks(cellTo).Last(); chainOut != nil {
			log.Printf("%v has last outbound link %v to %v", cellTo,
				chainOut, (*Cell)(*chainOut.FurthestPath))
			createdLink.FurthestPath = chainOut.FurthestPath
		} else {
			createdLink.FurthestPath = new(BoxedCellPointer)
			*createdLink.FurthestPath = createdLink.To
		}
	}

	log.Printf("Created link %v with path to %v", createdLink,
		(*Cell)(*createdLink.FurthestPath))

	// One beyond the link path from cellFrom, is the next adjacency position.
	nextAdjacencyIndex := forward.Increment((*createdLink.FurthestPath).Index)
	log.Printf("nextAdjacencyIndex is %v", nextAdjacencyIndex)

	// 3.2.1 Connectedness: Use of this adjacency creates a new one,
	// spanning to exactly nextAdjacencyIndex (and no further).
	newAdjacency := new(Adjacency)
	newAdjacency.From = cellFrom
	newAdjacency.Position = forward.Increment(usedAdjacency.Position)

	// Replace the outbound adjacency, and add the new outbound link.
	forward.OutboundLinks(cellFrom).Add(createdLink)
	forward.SetOutboundAdjacency(cellFrom, newAdjacency)

	if createdLink.Depth == 0 {
		forward.SetLastOutboundLinkD0(cellFrom, createdLink)
	} else {
		forward.SetLastOutboundLinkD1(cellFrom, createdLink)
	}

	// Replace the inbound adjacency, and set the inbound link.
	forward.InboundAdjacencies(cellTo).Remove(usedAdjacency)
	chart.updateAdjacencyTo(newAdjacency, nextAdjacencyIndex)
	forward.SetInboundLink(cellTo, createdLink)

	// Examine other inbound adjacencies into cellTo.
	for adjacency := range forward.InboundAdjacencies(cellTo) {
		if adjacency.IsBlocked() {
			continue
		}

		if !forward.Less(adjacency.From.Index, cellFrom.Index) {
			// 3.2.2 Covering Links - Adjacencies which are fully covered by 
			// usedAdjacency are permanently blocked. No further attachments
			// are possible from cellFrom in the forward direction.
			log.Printf("Blocking covered %v", adjacency)
			adjacency.BlockedDepths = [2]bool{true, true}
		} else {
			// 3.2.1 Connectedness & 3.2.1 Minimality:
			// Non-covered adjacencies must be shifted through the link-path
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
	return
}

func (chart *Chart) updateAdjacencyTo(adjacency *Adjacency, toIndex int) {
	invariant.IsFalse(adjacency.IsBlocked())

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
