package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func mainAnalyze() {
	if len(os.Args) != 5 {
		log.Fatal("You must provide an input directory, reference data file, and an output prefix")
	}

	inputDirectory := os.Args[2]
	votingListFilename := os.Args[3]
	outputPrefix := os.Args[4]
	log.Println(outputPrefix)

	inputData, err := getAllInputData(inputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(inputData)

	list, err := getVotingListFromFile(votingListFilename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(list)
}

func getVotingListFromFile(filename string) (VotingList, error) {
	var vl VotingList
	bits, err := os.ReadFile(filename)
	if err != nil {
		return vl, err
	}

	err = json.Unmarshal(bits, &vl)
	if err != nil {
		return vl, err
	}

	return vl, nil
}

func getAllInputData(inputDirectory string) (map[string][]string, error) {
	mappy := make(map[string][]string)

	entries, err := os.ReadDir(inputDirectory)
	if err != nil {
		return mappy, err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.Contains(e.Name(), ".json") {
			continue
		}

		fullPath := fmt.Sprintf("%v/%v", inputDirectory, e.Name())
		bits, err := os.ReadFile(fullPath)
		if err != nil {
			return mappy, err
		}

		var jsonStringArray []string
		err = json.Unmarshal(bits, &jsonStringArray)
		if err != nil {
			return mappy, err
		}

		mappy[e.Name()] = jsonStringArray
	}

	return mappy, nil
}
