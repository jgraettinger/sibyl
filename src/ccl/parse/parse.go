package parse

import (
    //"log"
    "ccl"
    . "ccl/chart"
)

type Parser struct {
	Chart  *Chart
	Scorer ccl.Scorer
}

func NewParser(scorer ccl.Scorer) (parser *Parser) {
	return &Parser{NewChart(), scorer}
}

func (parser *Parser) ParseIncremental(token string) {
	parser.Chart.AddCell(token)
	for addBestLink(parser.Chart, parser.Scorer) {
	}
}

func (parser *Parser) ParseUtterance(tokens []string) {
	for _, token := range tokens {
		parser.ParseIncremental(token)
	}
}

func addBestLink(chart *Chart, scorer ccl.Scorer) bool {
	if len(chart.Cells) < 2 {
		return false
	}

	var bestScore float64
	var bestScoreDepth uint
	var bestAdjacency *Adjacency

	check := func(adjacency *Adjacency) {
		if adjacency.Used || adjacency.IsBlocked() || adjacency.To == nil {
			return
		}

		score, scoreDepth := scorer.Score(adjacency)

        //log.Printf("Adjacency %s has score %v depth %v", adjacency, score, scoreDepth)

		if score > bestScore && !adjacency.Blocked[scoreDepth] {
			bestScore = score
			bestScoreDepth = scoreDepth
			bestAdjacency = adjacency
		} else if score > bestScore && !adjacency.Blocked[1] {
			bestScore = score
			bestScoreDepth = 1
			bestAdjacency = adjacency
        }
	}

	index := len(chart.Cells) - 1

	// Prefer direct adjacency between the last & second-to-last
	//  cells, if either of these adjacencies has positive weight
	// Due to the incremental nature of the parser, this is sufficient
	//  to satisfy the general preference for direct adjacency,
	//  as any other linkable direct adjacencies have already been linked
	check(chart.Cells[index-1].Outbound[RIGHT][0])
	check(chart.Cells[index].Outbound[LEFT][0])

	if bestAdjacency != nil {
		chart.Use(bestAdjacency, bestScoreDepth)
		return true
	}

	// check all inbound adjacencies to the last cell
	for adjacency := range chart.Cells[index].Inbound[LEFT] {
		check(adjacency)
	}
	// check outbound adjacency from the last cell
	check(chart.Cells[index].Outbound[LEFT][len(chart.Cells[index].Outbound[LEFT])-1])

	if bestAdjacency != nil {
		chart.Use(bestAdjacency, bestScoreDepth)
		return true
	}
	return false
}

