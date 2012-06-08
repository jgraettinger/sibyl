package cclparse

import (
//	"fmt"
//	"invariant"
)

func (chart Chart) AddLink(lexicon Lexicon) bool {
	if len(chart.cells) < 2 {
		return false
	}

	var bestAdjacency *Adjacency
	var bestWeight float64
	var bestDepth uint8

	check := func(adjacency *Adjacency) {
		if adjacency.Used || adjacency.Blocked || adjacency.To == nil {
			return
		}

		weight, depth := lexicon.Score(adjacency)

		if weight > bestWeight {
			bestWeight = weight
			bestDepth = depth
			bestAdjacency = adjacency
		}
	}

	index := len(chart.cells) - 1

	// Prefer direct adjacency between the last & second-to-last
	//  cells, if either of these adjacencies has positive weight
	// Due to the incremental nature of the parser, this is sufficient
	//  to satisfy the general preference for direct adjacency,
	//  as any other linkable direct adjacencies have already been linked
	check(chart.cells[index].Outbound[LEFT][0])
	check(chart.cells[index-1].Outbound[RIGHT][0])

	if bestAdjacency != nil {
		chart.Use(bestAdjacency)
		return true
	}

	// check all inbound adjacencies to the last cell
	for adjacency := range chart.cells[index].Inbound[LEFT] {
		check(adjacency)
	}
	// check outbound adjacency from the last cell
	check(chart.cells[index].Outbound[LEFT][len(chart.cells[index].Outbound[LEFT])-1])

	if bestAdjacency != nil {
		chart.Use(bestAdjacency)
		return true
	}
	return false
}

func (chart Chart) Use(usedAdjacency *Adjacency) {

	// inbound adjacencies to From are /moved/ through
	//  the link path created by using this adjacency
	moveIndex := usedAdjacency.To.Index + isign(usedAdjacency.Position)

	var side int
	if usedAdjacency.Position < 0 {
		side = RIGHT
	} else {
		side = LEFT
	}

	usedAdjacency.Used = true

	// using this adjacency creates a new one
	newAdjacency := new(Adjacency)
	newAdjacency.From = usedAdjacency.From
	newAdjacency.Position = usedAdjacency.Position + isign(usedAdjacency.Position)

	if moveIndex == -1 {
		newAdjacency.To = nil
	} else if moveIndex == len(chart.cells) {
		newAdjacency.To = nil
		chart.endInbound.Add(newAdjacency)
	} else {
		newAdjacency.To = chart.cells[moveIndex]
		chart.cells[moveIndex].Inbound[side].Add(newAdjacency)
	}

	newAdjacency.From.Outbound[side].Add(newAdjacency)

	// Other adjacencies into To on this side are /moved/
	//  through the link-path created by using this adjacency

	for adjacency := range usedAdjacency.To.Inbound[side] {
		if adjacency.Blocked || adjacency.Used {
			continue
		}

		usedAdjacency.To.Inbound[side].Remove(adjacency)

		if moveIndex == -1 {
			adjacency.To = nil
		} else if moveIndex == len(chart.cells) {
			adjacency.To = nil
			chart.endInbound.Add(adjacency)
		} else {
			adjacency.To = chart.cells[moveIndex]
			chart.cells[moveIndex].Inbound[side].Add(adjacency)
		}
	}
}
