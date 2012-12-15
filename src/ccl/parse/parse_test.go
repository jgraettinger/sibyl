package parse

import (
    "log"
	. "ccl/chart"
	"ccl/graphviz"
	. "ccl/mocks"
	"os"
	"testing"
)

func TestBasicLinking(t *testing.T) {

    log.SetFlags(log.Lshortfile)

	fixtures := NewParseFixtures()
	fixtures.AddUsed("the", "cat")
	fixtures.AddUsed("cat", "the")
	fixtures.AddUsed("ate", "cat")
	fixtures.AddUsed("ate", "a")
	fixtures.AddUsed("a", "mouse")
	fixtures.AddUsed("mouse", "a")

	parseAndValidate(fixtures, []string{"the", "cat", "ate", "a", "mouse"}, t)
	parseAndValidate(fixtures, []string{"mouse", "a", "ate", "cat", "the"}, t)
}

func TestMontonicity(t *testing.T) {

    // Because of the x =1=> y link, x => z must have depth 1.
    fixtures := NewParseFixtures()
    fixtures.AddUsed("x", "y").WithScoreDepth(1).AtDepth(1)
    fixtures.AddUsed("x", "z").WithScoreDepth(0).AtDepth(1)

    parseAndValidate(fixtures, []string{"x", "y", "z"}, t)
    parseAndValidate(fixtures, []string{"z", "y", "x"}, t)
}

func TestMinimality(t *testing.T) {

	fixtures := NewParseFixtures()
	fixtures.AddUsed("x", "y")
	fixtures.AddUsed("y", "z")
    // x's adjacency to z is moved when using the y => z link
	fixtures.AddNotAdjacent("x", "z")

	parseAndValidate(fixtures, []string{"x", "y", "z"}, t)
	parseAndValidate(fixtures, []string{"z", "y", "x"}, t)
}

func TestConnectedness(t *testing.T) {

    fixtures := NewParseFixtures()
    fixtures.AddNotUsed("x", "y").WithScore(0)
    fixtures.AddNotUsed("y", "z").WithScore(0)
    // As the x => y adjacency isn't used,
    // x cannot be adjacent to z.
    fixtures.AddNotAdjacent("x", "z").WithScore(1)

    parseAndValidate(fixtures, []string{"x", "y", "z"}, t)
    parseAndValidate(fixtures, []string{"z", "y", "x"}, t)
}

func TestBlockingXYZ(t *testing.T) {

    fixtures := NewParseFixtures()
    fixtures.AddUsed("y", "x")
    fixtures.AddUsed("x", "y")
    fixtures.AddUsed("y", "z").WithScoreDepth(0).AtDepth(1)

    parseAndValidate(fixtures, []string{"y", "x", "z"}, t)
    parseAndValidate(fixtures, []string{"z", "x", "y"}, t)
}

func TestBlockingYWXZ(t *testing.T) {

    fixtures := NewParseFixtures()
    fixtures.AddUsed("x", "y").AtDepth(1)
    fixtures.AddUsed("w", "x")
    fixtures.AddUsed("x", "w")
    fixtures.AddBlocked("w", "z")

    parseAndValidate(fixtures, []string{"y", "w", "x", "z"}, t)
}

func TestBlockingWXYZ(t *testing.T) {

    fixtures := NewParseFixtures()
    fixtures.AddUsed("x", "y").WithScoreDepth(1)
    fixtures.AddUsed("w", "x")
    fixtures.AddUsed("x", "w")
    fixtures.AddBlocked("w", "y")
    fixtures.AddNotAdjacent("w", "z")

    parseAndValidate(fixtures, []string{"w", "x", "y", "z"}, t)
}

/*
func TestFoobar(t *testing.T) {

    w -0> x
    x -0> y
    y -1> z

    both:
    y -0> x // passes blocking to y
    x -0> w // passes blocking to w, blocks w -> z

    x -0> w // d=0 blocked of w -> z
    y -0> x // y's d=1 propogates to fuly block x -> z
}

func TestEqualityFoo(t *testing.T) {

    w, x, y, z

    w -0> x
    w -1> y
    y -0> z
    z -1> w

    // this should *pass*, despite looking like a cycle might be possible

}
*/

func TestEqualityXY(t *testing.T) {

    fixtures := NewParseFixtures()
    fixtures.AddUsed("x", "y").WithScoreDepth(1)
    fixtures.AddUsed("y", "x").AtDepth(1)

    parseAndValidate(fixtures, []string{"x", "y"}, t)
}

func TestEqualityXXYZ(t *testing.T) {

    fixtures := NewParseFixtures()
    fixtures.AddUsed("x", "x'")
    fixtures.AddUsed("x'", "x")
    fixtures.AddUsed("y", "x'")
    fixtures.AddUsed("x'", "y")
    fixtures.AddUsed("z", "y").WithScoreDepth(1)
    fixtures.AddUsed("x", "z").AtDepth(1)

    parseAndValidate(fixtures, []string{"x", "x'", "y", "z"}, t)
}

func parseAndValidate(fixtures ParseFixtures,
	utterance []string, t *testing.T) {

    if t.Failed() {
        return
    }

	parser := NewParser(fixtures)
	parser.ParseUtterance(utterance)

	fixtures.Validate(parser.Chart, t)
    graphvizOut(parser.Chart)
}

func graphvizOut(chart *Chart) {

	if graphOut, err := os.Create("/tmp/chart.graphviz"); err != nil {
		panic(err)
	} else {
		graphOut.Write([]byte(graphviz.RenderChart(chart)))
		graphOut.Close()
	}

}
