package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
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

	linkFilenames := make(map[string]string)

	for _, list := range lists {
		outputFilename := getOutputFilename(list.Name)
		linkFilenames[outputFilename] = list.Name
		err = generateListHtml(list, listTemplate, itemTemplate, generatedDir, outputFilename)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = generateIndexHtml(linkFilenames, indexTemplate, linkTemplate, generatedDir)
	if err != nil {
		log.Fatal(err)
	}
}

func generateIndexHtml(linkFilenames map[string]string, indexTemplate string, linkTemplate string, generatedDir string) error {

	linkKeys := []string{}
	for key, _ := range linkFilenames {
		linkKeys = append(linkKeys, key)
	}

	sort.Strings(linkKeys)

	links := []string{}
	for i := range linkKeys {
		key := linkKeys[i]
		value := linkFilenames[key]
		localLinkTemplate := linkTemplate
		localLinkTemplate = strings.Replace(localLinkTemplate, "$link", key, -1)
		localLinkTemplate = strings.Replace(localLinkTemplate, "$name", value, -1)
		links = append(links, localLinkTemplate)
	}

	joinedLinks := strings.Join(links, "\n")
	localIndexTemplate := indexTemplate
	localIndexTemplate = strings.Replace(localIndexTemplate, "$links", joinedLinks, -1)

	fullPath := fmt.Sprintf("%v/%v", generatedDir, "index.html")
	return os.WriteFile(fullPath, []byte(localIndexTemplate), 066)
}

func getOutputFilename(name string) string {
	outputFilename := name
	outputFilename = strings.Replace(outputFilename, " ", "-", -1)
	outputFilename = strings.ToLower(outputFilename)
	outputFilename = fmt.Sprintf("%v.html", outputFilename)
	return outputFilename
}

func generateListHtml(list VotingList, listTemplate string, itemTemplate string, generatedDir string, outputFilename string) error {
	err := checkVotingList(list)
	if err != nil {
		return err
	}

	itemTemplates := []string{}
	for _, item := range list.Items {
		localItemTemplate := itemTemplate
		localItemTemplate = strings.Replace(localItemTemplate, "$name", item.Name, -1)
		localItemTemplate = strings.Replace(localItemTemplate, "$year", item.Year, -1)
		localItemTemplate = strings.Replace(localItemTemplate, "$link", item.Link, -1)
		localItemTemplate = strings.Replace(localItemTemplate, "$uniqueKey", item.UniqueKey(), -1)
		itemTemplates = append(itemTemplates, localItemTemplate)
	}

	joinedItemTemplates := strings.Join(itemTemplates, "\n")

	localListTemplate := listTemplate
	localListTemplate = strings.Replace(localListTemplate, "$title", list.Name, -1)
	localListTemplate = strings.Replace(localListTemplate, "$items", joinedItemTemplates, -1)

	fullPath := fmt.Sprintf("%v/%v", generatedDir, outputFilename)
	return os.WriteFile(fullPath, []byte(localListTemplate), 066)
}

func checkVotingList(list VotingList) error {
	uniqueNames := make(map[string]bool)
	for _, item := range list.Items {
		uniqueNames[item.UniqueKey()] = true
	}

	if len(uniqueNames) != len(list.Items) {
		return errors.New("list needs to have unique name and year items")
	}

	return nil
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
			list, err := getVotingListFromFile(fullPath)
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
