package cclparse

import (
    "io"
    "os"
    "encoding/json"
    "invariant"
)

type AdjacencyPoint struct {
    Token string
    Position int
}
type LabelWeight struct {
    ClassWeight float32
    AdjacencyWeight float32
}
type AdjacencyStatistics struct {
    AdjacencyPoint

    UpdateCount uint64
    Stop uint64

    InRaw float32
    Out float32
    In float32

    Labels map[string]LabelWeight
}

type Lexicon map[*AdjacencyPoint]*AdjacencyStatistics

func NewLexiconFromJson(input io.ReadCloser) (
        lexicon Lexicon, err error) {

    var file io.ReadCloser
    lexicon = make(Lexicon)

    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    for {
        var adjStats AdjacencyStatistics

        if err = decoder.Decode(&adjStats); err == io.EOF {
            break
        } else if err != nil {
            return
        }
        lexicon[&adjStats.AdjacencyPoint] = &adjStats
    }
    err = nil
    return
}

func linkWeight(out, in *AdjacencyStatistics) float32 {

    invariant.NotNil(out)

    var bestLabel string
    var bestIsClass bool

    for token, weights := range(out.Labels) {

        


    }



}
