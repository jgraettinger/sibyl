package cclparse

import (
    "fmt"
    "testing"
)

func TestLoadLexicon(t *testing.T) {
    lexicon, err := NewLexiconFromJson("/home/johng/sibyl/brown_lexicon.json")

    if err != nil {
        t.Error(err)
    }
    fmt.Println(lexicon)
}
