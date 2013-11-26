package mocks

import (
	. "ccl/chart"
	"testing"
)

type AdjacencyFixtureKey struct {
	From, To string
}

type AdjacencyFixture struct {
	AdjacencyFixtureKey

	// Defaults to 1.0 if nil.
	ScoreValue *float64
	// Defaults to 0 if nil.
	ScoreDepth *uint

	// If nil, any position is allowed.
	ExpectPosition *int
	// If nil, any depth is allowed.
	ExpectDepth *uint

	ExpectUsed         bool
	ExpectBlocked      bool
	ExpectNotAdjacent  bool
	ExpectStoppingPunc bool
}

func (f *AdjacencyFixture) WithScore(score float64) *AdjacencyFixture {
	f.ScoreValue = &score
	return f
}
func (f *AdjacencyFixture) WithScoreDepth(depth uint) *AdjacencyFixture {
	f.ScoreDepth = &depth
	return f
}

func (f *AdjacencyFixture) AtDepth(depth uint) *AdjacencyFixture {
	f.ExpectDepth = &depth
	return f
}

func (f *AdjacencyFixture) Score() (score float64, depth uint) {
	score, depth = 1.0, 0

	if f == nil {
		score = 0.0
		return
	}
	if f.ScoreValue != nil {
		score = *f.ScoreValue
	}
	if f.ScoreDepth != nil {
		depth = *f.ScoreDepth
	}
	return
}

func (f *AdjacencyFixture) Validate(adjacency *Adjacency, t *testing.T) {
	if f == nil {
		if adjacency.Used {
			t.Error("No fixture for observed adjacency", adjacency)
		}
		return
	}

    if f.ExpectNotAdjacent {
        t.Error("Didn't expect", adjacency)
    }
	if f.ExpectUsed && !adjacency.Used {
		t.Error("Expected to be Used", adjacency)
	}
	if f.ExpectBlocked && !adjacency.IsBlocked() {
		t.Error("Expected to be Blocked", adjacency)
	}
	if f.ExpectStoppingPunc && !adjacency.StoppingPunc {
		t.Error("Expected StoppingPunc blocking", adjacency)
	}
	if f.ExpectPosition != nil && *f.ExpectPosition != adjacency.Position {
		t.Error("Expected position %d: %s", *f.ExpectPosition, adjacency)
	}
	if f.ExpectDepth != nil && *f.ExpectDepth != adjacency.UsedDepth {
		t.Errorf("Expected depth %d: %s", *f.ExpectDepth, adjacency)
	}
}
