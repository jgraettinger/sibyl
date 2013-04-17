package main

import (
	"bufio"
	"ccl/chart"
	"ccl/graphviz"
	"ccl/lexicon"
	"log"
	"net/http"
	"os"
	"strings"
)

func parseSentence(req *http.Request,
	w http.ResponseWriter, lexicon *Lexicon) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cht := chart.NewChart()
	for _, field := range req.Form["token"] {
		cht.AddTokens(chart.Token(field))
	}

	

}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v /path/to/lexicon.json", os.Args[1])
	}
	lexicon, err := lexicon.Open(os.Args[1])
	if err != nil {
		log.Panic(err)
	}

	for {
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		var tokens []chart.Token
		for _, field := range strings.Fields(line) {
			posIndex := strings.IndexRune(field, '/')
			tokens = append(tokens, chart.Token(field[:posIndex]))
		}
		if len(tokens) == 0 {
			continue
		}
		log.Printf("%v\n", tokens)

		chart := chart.NewChart()
		chart.AddTokens(tokens...)

		defer func() {
			if r := recover(); r != nil {
				log.Print(r)

				http.HandleFunc("/chart",
					func(w http.ResponseWriter, r *http.Request) {
						graphviz.RenderChartSvg(chart, w)
					})
				log.Fatal(http.ListenAndServe(":8080", nil))
			}
		}()

		for chart.NextCell() != nil {
			for {
				adjacency, depth, score := chart.BestAdjacency(lexicon)
				if adjacency == nil {
					break
				}
				log.Printf("Using %v@%d (%v)\n", adjacency, depth, score)
				chart.Use(adjacency, depth)
			}
		}
		lexicon.Learn(chart)
	}

	/*
		lexicon := cclparse.NewLexicon()

		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Panic(err)
		}

		input := bufio.NewReader(file)

		var chart *cclparse.Chart
		for {
			var tokens []string

			if line, isPrefix, err := input.ReadLine(); isPrefix {
				log.Panic("line too long")
			} else if err == io.EOF {
				break
			} else if err != nil {
				log.Panic(err)
			} else {
				tokens = strings.Split(strings.ToLower(string(line)), " ")
			}

			chart = cclparse.NewChart()
			for _, token := range tokens {
				chart.AddCell(token)
				for chart.AddLink(lexicon) {
				}
			}
			lexicon.Learn(chart)
		}

		if graphOut, err := os.Create("/tmp/chart.graphviz"); err != nil {
			log.Panic(err)
		} else {
			graphOut.Write([]byte(chart.AsGraphviz()))
			graphOut.Close()
		}

		type LexiconTuple struct {
			Token    string  `json:"token"`
			Position int     `json:"position"`
			Label    *string `json:"label"`
			Stat     *string `json:"stat"`
			Value    float64 `json:"weight"`
		}
		tuples := make(chan LexiconTuple)

		emitOutputTuples := func() {
			for _, adjPointStats := range lexicon {

				tuple := LexiconTuple{
					Token:    adjPointStats.Token,
					Position: adjPointStats.Position,
				}

				sptr := func(s string) *string { return &s }

				// Output adjacency point bootstrap statistics
				tuple.Stat, tuple.Value = sptr("count"), float64(adjPointStats.Count)
				tuples <- tuple
				tuple.Stat, tuple.Value = sptr("stop"), float64(adjPointStats.Stop)
				tuples <- tuple
				tuple.Stat, tuple.Value = sptr("in_raw"), float64(adjPointStats.InRaw)
				tuples <- tuple
				tuple.Stat, tuple.Value = sptr("out"), adjPointStats.Out
				tuples <- tuple
				tuple.Stat, tuple.Value = sptr("in"), adjPointStats.In
				tuples <- tuple

				tuple.Stat = nil

				for label, weight := range adjPointStats.LabelWeights.FilterToTopN(1000) {
					tuple.Label = sptr(label.String())
					tuple.Value = weight
					tuples <- tuple
				}
			}
			close(tuples)
		}
		go emitOutputTuples()

		for {
			// encoder := json.NewEncoder(os.Stdout)

			if tuple, okay := <-tuples; !okay {
				break
			} else {

				if tuple.Label != nil {
					fmt.Printf("%v\t%d\t\"%v\"\t%0.1f\n", tuple.Token,
						tuple.Position, *tuple.Label, tuple.Value)
				} else {
					fmt.Printf("%v\t%d\t%v\t%0.1f\n", tuple.Token,
						tuple.Position, *tuple.Stat, tuple.Value)
				}

				/*
				   if err := encoder.Encode(tuple); err != nil {
				       log.Panic(err)
				   }
				* /
			}
		}
	*/
}

