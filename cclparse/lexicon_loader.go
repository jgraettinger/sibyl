package cclparse
/*
import (
    "io"
    "os"
    "encoding/json"
)

func NewLexiconFromJson(path string) (lexicon Lexicon, err error) {

    var file io.ReadCloser
    lexicon = make(Lexicon)

    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    for {
		var record struct {
            Token string
			Position int
			Update_count uint64
		    Stop uint64
		    InRaw float64
		    Out float64
		    In float64
		    Labels map[string]struct {Adjacency_Weight, Class_Weight float64}
		}

        if err = decoder.Decode(&record); err == io.EOF {
            break
        } else if err != nil {
            return
        }

		adjStats := NewAdjacencyStatistics(
			AdjacencyPoint{record.Token, record.Position})

		adjStats.count = record.Update_count
		adjStats.stop = record.Stop
		adjStats.inRaw = int64(record.InRaw)
		adjStats.out = record.Out
		adjStats.in = record.In

		for token, weights := range(record.Labels) {
			adjStats.labelWeights[Label{CLASS, token}] = weights.Class_Weight
			adjStats.labelWeights[Label{ADJACENCY, token}] =
				weights.Adjacency_Weight
		}

        lexicon[adjStats.AdjacencyPoint] = adjStats
    }
    err = nil
    return
}
*/
