package mocks

import (
    "testing"
    "invariant"
    . "ccl/chart"
)

type ParseFixtures map[AdjacencyFixtureKey]*AdjacencyFixture

func NewParseFixtures() ParseFixtures {
    return make(ParseFixtures)
}

func (fixtures ParseFixtures) AddUsed(from, to string) *AdjacencyFixture {

    fixture := new(AdjacencyFixture)
    fixture.From = from
    fixture.To = to
    fixture.ExpectUsed = true

    _, present := fixtures[fixture.AdjacencyFixtureKey]
    invariant.IsFalse(present)

    fixtures[fixture.AdjacencyFixtureKey] = fixture
    return fixture
}

func (fixtures ParseFixtures) AddNotUsed(from, to string) *AdjacencyFixture {

    fixture := new(AdjacencyFixture)
    fixture.From = from
    fixture.To = to
    fixture.ExpectUsed = false

    _, present := fixtures[fixture.AdjacencyFixtureKey]
    invariant.IsFalse(present)

    fixtures[fixture.AdjacencyFixtureKey] = fixture
    return fixture
}

func (fixtures ParseFixtures) AddBlocked(from, to string) *AdjacencyFixture {

    fixture := new(AdjacencyFixture)
    fixture.From = from
    fixture.To = to
    fixture.ExpectBlocked = true

    _, present := fixtures[fixture.AdjacencyFixtureKey]
    invariant.IsFalse(present)

    fixtures[fixture.AdjacencyFixtureKey] = fixture
    return fixture
}

func (fixtures ParseFixtures) AddNotAdjacent(
    from, to string) *AdjacencyFixture {

    fixture := new(AdjacencyFixture)
    fixture.From = from
    fixture.To = to
    fixture.ExpectNotAdjacent = true

    _, present := fixtures[fixture.AdjacencyFixtureKey]
    invariant.IsFalse(present)

    fixtures[fixture.AdjacencyFixtureKey] = fixture
    return fixture
}

// ParseFixtures conforms to Scorer interface
func (fixtures ParseFixtures) Score(adjacency *Adjacency) (float64, uint) {
    invariant.NotNil(adjacency.From)
    invariant.NotNil(adjacency.To)

    key := AdjacencyFixtureKey{adjacency.From.Token, adjacency.To.Token}
    return fixtures[key].Score()
}

func (fixtures ParseFixtures) Validate(chart *Chart, t *testing.T) {

    // Copy fixtures map, so we can track observed adjacencies via deletion.
    fixturesCopy := NewParseFixtures()
    for key, value := range fixtures {
        fixturesCopy[key] = value
    }

    // Closure to validate expected properties of the adjacency.
    validate := func(adjacency *Adjacency) {
        if adjacency.To == nil {
            return
        }
        key := AdjacencyFixtureKey{adjacency.From.Token, adjacency.To.Token}
        fixturesCopy[key].Validate(adjacency, t)
        delete(fixturesCopy, key)
    }
    // Walk through adjacencies of the parsed chart.
    for _, cell := range chart.Cells {
        for _, adjacency := range cell.Outbound[LEFT] {
            validate(adjacency)
        }
        for _, adjacency := range cell.Outbound[RIGHT] {
            validate(adjacency)
        }
    }
    // Any un-observed but expected adjacency fixtures are an error.
    for _, fixture := range fixturesCopy {
        if !fixture.ExpectNotAdjacent {
            t.Error("Didn't see expected adjacency", *fixture)
        }
    }
}
