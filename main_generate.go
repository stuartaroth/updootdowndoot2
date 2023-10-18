package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func mainGenerate() {
	log.Println("starting generate...")

	generatedDir := "generated_web"
	err := cleanGeneratedDir(generatedDir)
	if err != nil {
		log.Fatal(err)
	}

	indexTemplate, err := getTemplateString("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	itemTemplate, err := getTemplateString("templates/item.html")
	if err != nil {
		log.Fatal(err)
	}

	linkTemplate, err := getTemplateString("templates/link.html")
	if err != nil {
		log.Fatal(err)
	}

	listTemplate, err := getTemplateString("templates/list.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(indexTemplate), len(itemTemplate), len(linkTemplate), len(listTemplate))

	jsonDir := "json_data"
	lists, err := getVotingLists(jsonDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, list := range lists {
		err = generateListHtml(list, listTemplate, itemTemplate, generatedDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateListHtml(list VotingList, listTemplate string, itemTemplate string, generatedDir string) error {
	itemTemplates := []string{}
	for _, item := range list.Items {
		localItemTemplate := itemTemplate
		localItemTemplate = strings.Replace(localItemTemplate, "$name", item.Name, -1)
		localItemTemplate = strings.Replace(localItemTemplate, "$year", item.Year, -1)
		localItemTemplate = strings.Replace(localItemTemplate, "$link", item.Link, -1)
		localItemTemplate = strings.Replace(localItemTemplate, "$uuid", item.Uuid, -1)
		itemTemplates = append(itemTemplates, localItemTemplate)
	}

	joinedItemTemplates := strings.Join(itemTemplates, "\n")

	localListTemplate := listTemplate
	localListTemplate = strings.Replace(localListTemplate, "$title", list.Name, -1)
	localListTemplate = strings.Replace(localListTemplate, "$items", joinedItemTemplates, -1)

	outputFilename := list.Name
	outputFilename = strings.Replace(outputFilename, " ", "-", -1)
	outputFilename = strings.ToLower(outputFilename)
	outputFilename = fmt.Sprintf("%v.html", outputFilename)
	fullPath := fmt.Sprintf("%v/%v", generatedDir, outputFilename)
	return os.WriteFile(fullPath, []byte(localListTemplate), 066)
}

func getVotingLists(jsonDir string) ([]VotingList, error) {
	lists := []VotingList{}
	distinctNames := make(map[string]bool)

	entries, err := os.ReadDir(jsonDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			fullPath := fmt.Sprintf("%v/%v", jsonDir, e.Name())
			bits, err := os.ReadFile(fullPath)
			if err != nil {
				return lists, err
			}

			var list VotingList
			err = json.Unmarshal(bits, &list)
			if err != nil {
				return lists, err
			}

			lists = append(lists, list)
			distinctNames[list.Name] = true
		}
	}

	if len(lists) != len(distinctNames) {
		return lists, errors.New("each list must have a distinct name")
	}

	return lists, nil
}

func getTemplateString(filename string) (string, error) {
	bits, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(bits), nil
}

func cleanGeneratedDir(generatedDir string) error {
	entries, err := os.ReadDir(generatedDir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Name() == ".do" {
			continue
		}

		fullPath := fmt.Sprintf("%v/%v", generatedDir, e.Name())
		if e.IsDir() {
			err = os.RemoveAll(fullPath)
			if err != nil {
				return err
			}
		} else {
			err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
