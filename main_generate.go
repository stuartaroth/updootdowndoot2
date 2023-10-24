package main

import (
	"encoding/json"
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

	imageTemplate, err := getTemplateString("templates/image.html")
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
	allKeys := make(map[string]bool)

	for _, list := range lists {
		outputFilename := getOutputFilename(list.Name)
		linkFilenames[outputFilename] = list.Name
		err = generateListHtml(list, listTemplate, itemTemplate, imageTemplate, generatedDir, outputFilename)
		if err != nil {
			log.Fatal(err)
		}

		addUniqueKeys(list, allKeys)
	}

	err = generateIndexHtml(linkFilenames, indexTemplate, linkTemplate, generatedDir)
	if err != nil {
		log.Fatal(err)
	}

	err = generateAllKeysJson(generatedDir, allKeys)
	if err != nil {
		log.Fatal(err)
	}
}

func addUniqueKeys(list VotingList, allKeys map[string]bool) {
	for _, item := range list.Items {
		allKeys[item.UniqueKey()] = true
	}
}

func generateAllKeysJson(generatedDir string, allKeys map[string]bool) error {
	localKeys := []string{}

	for key, _ := range allKeys {
		localKeys = append(localKeys, key)
	}

	sort.Strings(localKeys)

	bits, err := json.MarshalIndent(localKeys, "", "    ")
	if err != nil {
		return err
	}

	fullPath := fmt.Sprintf("%v/%v", generatedDir, "all-keys.json")

	err = os.WriteFile(fullPath, bits, 0666)
	if err != nil {
		return err
	}

	return nil
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

func generateListHtml(list VotingList, listTemplate string, itemTemplate string, imageTemplate string, generatedDir string, outputFilename string) error {
	err := checkVotingList(list)
	if err != nil {
		return err
	}

	itemTemplates := []string{}
	for _, item := range list.Items {
		imageContent := getImageContent(imageTemplate, item, generatedDir)

		localItemTemplate := itemTemplate
		localItemTemplate = strings.Replace(localItemTemplate, "$image", imageContent, -1)
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

func getImageContent(imageTemplate string, item VotingItem, generatedDir string) string {
	uniqueKey := item.UniqueKey()
	fullPath := fmt.Sprintf("%v/%v/%v.jpg", generatedDir, "images", uniqueKey)
	_, err := os.Stat(fullPath)
	if err != nil {
		return ""
	}

	localImageTemplate := imageTemplate
	localImageTemplate = strings.Replace(localImageTemplate, "$uniqueKey", uniqueKey, -1)
	return localImageTemplate
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
		if !e.IsDir() {
			err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		} else if e.Name() != "images" {
			err = os.RemoveAll(fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
