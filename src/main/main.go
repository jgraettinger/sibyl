package main

import (
	"bufio"
	"cclparse"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

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

	for _, value := range lexicon {
		bytes, err := json.Marshal(value)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(bytes))
	}
}
