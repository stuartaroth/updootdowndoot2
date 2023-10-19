package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"
)

func mainLetterboxd() {
	if len(os.Args) != 4 {
		log.Fatal("You must provide an input and output file")
	}

	inputFilename := os.Args[2]
	outputFilename := os.Args[3]

	inputBits, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Read from", inputFilename)

	inputContents := string(inputBits)
	lines := strings.Split(inputContents, "\r\n")

	if len(lines) < 6 {
		log.Fatal("Needs at list one entry")
	}

	items := []VotingItem{}

	usableLines := lines[5:]

	rowRows := strings.Join(usableLines, "\r\n")
	r := csv.NewReader(strings.NewReader(rowRows))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		item := VotingItem{
			Name: record[1],
			Year: record[2],
			Link: record[3],
		}

		items = append(items, item)
	}

	sort.Sort(SortVotingItem(items))

	outputObject := VotingList{
		Name:  inputFilename,
		Items: items,
	}

	outputBits, err := json.MarshalIndent(outputObject, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(outputFilename, outputBits, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Wrote to", outputFilename)
}
