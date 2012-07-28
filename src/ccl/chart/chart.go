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
		// As this is cellFrom's first outbound link in this direction,
		// using this adjacency doesn't create a new link-path. We can
		// forward an existing path (or create one if this begins a new path).
		if lastOut := forward.OutboundLinks(cellFrom).Last(); lastOut == nil {
			if chainIn := forward.InboundLink(cellFrom); chainIn != nil {
				createdLink.FurthestPath = chainIn.FurthestPath
			} else {
				createdLink.FurthestPath = new(BoxedCellPointer)
			}
			// This update is visible from all previous links on the path.
			*createdLink.FurthestPath = createdLink.To
		} else {
			// In this direction we expect that successive links will be added
			// to this path. It's more efficient to fork & isolate the present
			// path from future updates.
			createdLink.FurthestPath = lastOut.ForkFurthestPath()
		}
	} else {
		// Any successive links in this direction already exist. Propogate
		// an existing boxed path back to createdLink.
		if chainOut := forward.OutboundLinks(cellTo).Last(); chainOut != nil {
			createdLink.FurthestPath = chainOut.FurthestPath
		} else {
			createdLink.FurthestPath = new(BoxedCellPointer)
			*createdLink.FurthestPath = createdLink.To
		}
	}

	// One beyond the link path from cellFrom, is the next adjacency position.
	nextAdjacencyIndex := forward.Increment((*createdLink.FurthestPath).Index)
	log.Printf("nextAdjacencyIndex is %v", nextAdjacencyIndex)

	// Use of this adjacency creates a new one, spanning to nextAdjacencyIndex.
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

	// 3.2.1 Connectedness & Minimality:
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
			// Non-covered adjacencies must be shifted through the link-path
			// opened by usedAdjacency. This ensures minimality (as the
			// adjacency would otherwise be redudant with usedAdjacency),
			// and connectedness (as an adjacency is opened through to
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
