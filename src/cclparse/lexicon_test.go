package cclparse

import (
    "testing"
)

func TestFoo(t *testing.T) {
    lexicon, err := NewLexiconFromJson("/home/johng/sibyl/brown_lexicon.json")
    if err != nil {
        t.Error(err)
    }

    var xOut, yIn *AdjacencyStatistics
    var label Label
    var lWeight, kWeight float64
    var kDepth uint8

    xOut = lexicon[AdjacencyPoint{"i", 1}]
    yIn = lexicon[AdjacencyPoint{"know", -1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    xOut = lexicon[AdjacencyPoint{"know", -1}]
    yIn = lexicon[AdjacencyPoint{"i", 1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    xOut = lexicon[AdjacencyPoint{"know", 1}]
    yIn = lexicon[AdjacencyPoint{"the", -1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    xOut = lexicon[AdjacencyPoint{"the", -1}]
    yIn = lexicon[AdjacencyPoint{"know", 1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    xOut = lexicon[AdjacencyPoint{"the", 1}]
    yIn = lexicon[AdjacencyPoint{"boy", -1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    xOut = lexicon[AdjacencyPoint{"boy", -1}]
    yIn = lexicon[AdjacencyPoint{"the", 1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    /*
    xOut = lexicon[AdjacencyPoint{"boy", 1}]
    yIn = lexicon[AdjacencyPoint{"sleeps", -1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)
    */

    xOut = lexicon[AdjacencyPoint{"it", 1}]
    yIn = lexicon[AdjacencyPoint{"goes", -1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)

    xOut = lexicon[AdjacencyPoint{"goes", -1}]
    yIn = lexicon[AdjacencyPoint{"it", 1}]
    label, lWeight = bestMatchingLabel(xOut, yIn)
    kWeight, kDepth = lexicon.linkWeight(xOut, yIn)
    t.Logf("%#v %v %v %v", label, lWeight, kWeight, kDepth)
}

