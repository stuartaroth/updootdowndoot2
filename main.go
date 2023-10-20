package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 5 && os.Args[1] == "analyze" {
		mainAnalyze()
	} else if len(os.Args) == 2 && os.Args[1] == "generate" {
		mainGenerate()
	} else if len(os.Args) == 4 && os.Args[1] == "letterboxd" {
		mainLetterboxd()
	} else {
		fmt.Println("you done messed up")
	}
}
