package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reggie = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
)

type VotingList struct {
	Name  string       `json:"name"`
	Items []VotingItem `json:"items"`
}

type VotingItem struct {
	Name        string `json:"name"`
	Year        string `json:"year"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

func (item VotingItem) UniqueKey() string {
	return fmt.Sprintf("%v-%v", sanitize(item.Name), sanitize(item.Year))
}

func sanitize(input string) string {
	localInput := input
	localInput = reggie.ReplaceAllString(localInput, "")
	localInput = strings.Replace(localInput, " ", "-", -1)
	localInput = strings.ToLower(localInput)
	return localInput
}

type SortVotingItem []VotingItem

func (s SortVotingItem) Len() int {
	return len(s)
}

func (s SortVotingItem) Less(i, j int) bool {
	if s[i].Year != s[j].Year {
		return s[i].Year < s[j].Year
	}

	return s[i].Name < s[j].Name
}

func (s SortVotingItem) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
