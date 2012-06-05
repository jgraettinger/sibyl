package cclparse

import (
//	"fmt"
//	"invariant"
)

func (chart Chart) AddLink(lexicon *Lexicon) bool {
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

    uint index = len(chart.cells) - 1

    // Prefer direct adjacency between the last & second-to-last
    //  cells, if either of these adjacencies has positive weight
    // Due to the incremental nature of the parser, this is sufficient
    //  to satisfy the general preference for direct adjacency,
    //  as any other linkable direct adjacencies have already been linked
    check(chart.cells[index].Outbound.Left[0])
    check(chart.cells[index - 1].Outbound.Right[0])

    if bestAdjacency {
        chart.Use(bestAdjacency)
        return true
    }

    // check all inbound adjacencies to the last cell
    for adjacency := range chart.cells[index].Inbound.Left {
        check(adjacency)
    }

    if bestAdjacency {
        chart.Use(bestAdjacency)
        return true
    }
    return false
}

