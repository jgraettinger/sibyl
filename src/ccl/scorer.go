package ccl

import (
    "ccl/chart"
)

type Scorer interface {
	Score(adjacency *chart.Adjacency) (score float64, depth uint)
}
