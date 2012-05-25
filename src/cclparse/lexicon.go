package cclparse

import (
    "io"
    "os"
    "encoding/json"
)

type LabelType uint8

type AdjacencyPointKey struct {
    Token string
    Position int
}
type LabelWeight struct {
    ClassWeight float32
    AdjacencyWeight float32
}
type AdjacencyPoint struct {
    AdjacencyPointKey

    UpdateCount uint64
    Stop uint64

    InRaw float32
    Out float32
    In float32

    Labels map[string]LabelWeight
}

type AdjacencyPoints map[*AdjacencyPointKey]*AdjacencyPoint

func NewLexiconFromJson(input io.ReadCloser) (
        lexicon AdjacencyPoints, err error) {

    var file io.ReadCloser
    lexicon = make(AdjacencyPoints)

    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    for {
        var adjPoint AdjacencyPoint

        if err = decoder.Decode(&adjPoint); err == io.EOF {
            break
        } else if err != nil {
            return
        }
        lexicon[&adjPoint.AdjacencyPointKey] = &adjPoint
    }
    err = nil
    return
}

