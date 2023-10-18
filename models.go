package main

type VotingList struct {
	Name  string       `json:"name"`
	Items []VotingItem `json:"items"`
}

type VotingItem struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Year string `json:"year"`
	Link string `json:"link"`
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
