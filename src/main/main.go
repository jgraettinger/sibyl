package main

import (
	"fmt"
    "os"
	"cclparse"
)

func main() {

    lexicon := cclparse.NewLexicon()

	chart := cclparse.NewChart()
	for _, token := range([]string{"this", "is", "a", "hello", "world", "sentence"}) {
		chart.AddCell(token)
	}
    lexicon.Learn(chart)

	chart = cclparse.NewChart()
	for _, token := range([]string{"that", "sentence", "is", "a", "new", "example"}) {
		chart.AddCell(token)
	}
    lexicon.Learn(chart)


	fmt.Println(chart.AsGraphviz())

    fmt.Fprintln(os.Stderr, lexicon)
}

