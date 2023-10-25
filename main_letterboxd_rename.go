package main

import (
	"fmt"
	"github.com/schollz/closestmatch"
	"log"
	"os"
	"regexp"
)

var (
	idReggie = regexp.MustCompile(`^\d+-.*\.jpg`)
)

func mainLetterboxdRename() {
	if len(os.Args) != 4 {
		log.Fatal("You must provide an input directory and input filename")
	}

	inputDirectory := os.Args[2]
	inputFilename := os.Args[3]

	entries, err := os.ReadDir(inputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	list, err := getVotingListFromFile(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	bagSizes := []int{2}
	wordsToTest := []string{}

	for _, e := range entries {
		if idReggie.MatchString(e.Name()) {
			wordsToTest = append(wordsToTest, e.Name())
		}
	}

	cm := closestmatch.New(wordsToTest, bagSizes)

	filenameToItems := make(map[string][]VotingItem)

	for _, item := range list.Items {
		closestFilename := cm.Closest(item.Name)
		entries, _ := filenameToItems[closestFilename]
		entries = append(entries, item)
		filenameToItems[closestFilename] = entries
	}

	for filename, items := range filenameToItems {
		itemNames := []string{}
		for _, item := range items {
			itemNames = append(itemNames, item.Name)
		}

		cm2 := closestmatch.New(itemNames, bagSizes)

		closestItemName := cm2.Closest(filename)
		for _, item := range items {
			if item.Name == closestItemName {
				oldFullPath := fmt.Sprintf("%v/%v", inputDirectory, filename)
				newFullPath := fmt.Sprintf("%v/%v.jpg", inputDirectory, item.UniqueKey())
				err = os.Rename(oldFullPath, newFullPath)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
