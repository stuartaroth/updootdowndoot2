package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	validRatings = map[string]bool{
		"unseen":   true,
		"dislike":  true,
		"fine":     true,
		"like":     true,
		"love":     true,
		"terrible": true,
		"bad":      true,
		"neutral":  true,
		"good":     true,
		"great":    true,
	}
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

	uniqueKeyToRatingsToUsers := make(map[string]map[string][]string)
	uniqueKeyToItem := make(map[string]VotingItem)

	for _, item := range list.Items {
		uniqueKeyToRatingsToUsers[item.UniqueKey()] = map[string][]string{
			"unseen":   {},
			"dislike":  {},
			"fine":     {},
			"like":     {},
			"love":     {},
			"terrible": {},
			"bad":      {},
			"neutral":  {},
			"good":     {},
			"great":    {},
		}

		uniqueKeyToItem[item.UniqueKey()] = item
	}

	for userFilename, userRatings := range inputData {
		log.Println(userFilename)
		for _, s := range userRatings {
			uniqueKey, rating, err := getUniqueKeyAndRating(s)
			if err != nil {
				continue
			}

			uniqueKeyValue, exists := uniqueKeyToRatingsToUsers[uniqueKey]
			if !exists {
				continue
			}

			ratingValue, exists := uniqueKeyValue[rating]
			if !exists {
				continue
			}

			ratingValue = append(ratingValue, userFilename)
			uniqueKeyValue[rating] = ratingValue
			uniqueKeyToRatingsToUsers[uniqueKey] = uniqueKeyValue
		}
	}

	votingItemResults := []VotingItemResult{}

	for uniqueKey, ratingsToUsers := range uniqueKeyToRatingsToUsers {
		votingItem, exists := uniqueKeyToItem[uniqueKey]
		if !exists {
			continue
		}

		vir := VotingItemResult{
			VotingItem: votingItem,
			Unseen:     getSpecificRatingsToUsers(ratingsToUsers, "unseen"),
			Dislike:    getSpecificRatingsToUsers(ratingsToUsers, "dislike"),
			Fine:       getSpecificRatingsToUsers(ratingsToUsers, "fine"),
			Like:       getSpecificRatingsToUsers(ratingsToUsers, "like"),
			Love:       getSpecificRatingsToUsers(ratingsToUsers, "love"),
			Terrible:   getSpecificRatingsToUsers(ratingsToUsers, "terrible"),
			Bad:        getSpecificRatingsToUsers(ratingsToUsers, "bad"),
			Neutral:    getSpecificRatingsToUsers(ratingsToUsers, "neutral"),
			Good:       getSpecificRatingsToUsers(ratingsToUsers, "good"),
			Great:      getSpecificRatingsToUsers(ratingsToUsers, "great"),
		}

		votingItemResults = append(votingItemResults, vir)
	}

	for _, vir := range votingItemResults {
		log.Println(vir)
	}
}

func getSpecificRatingsToUsers(ratingsToUsers map[string][]string, key string) []string {
	ratings, exists := ratingsToUsers[key]
	if !exists {
		return []string{}
	}

	return ratings
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

func getUniqueKeyAndRating(s string) (string, string, error) {
	index := strings.LastIndex(s, "-")
	if index == -1 || index == len(s)-1 {
		return "", "", errors.New("entry missing hyphen or hypen at end")
	}

	uniqueKey := s[:index]
	rating := s[index+1:]
	_, exists := validRatings[rating]
	if !exists {
		return "", "", errors.New("invalid rating")
	}

	return uniqueKey, rating, nil
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
