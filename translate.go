package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// DicResult represents translation result
type DicResult struct {
	Head struct {
	} `json:"head"`
	Def []struct {
		Text string `json:"text"`
		Pos  string `json:"pos"`
		Tr   []struct {
			Text string `json:"text"`
			Pos  string `json:"pos"`
			Mean []struct {
				Text string `json:"text"`
			} `json:"mean"`
		} `json:"tr"`
	} `json:"def"`
}

func translate(msg string, dictionary string) []string {

	resp, err := http.Get("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=" +
		os.Getenv("YANDEX_API_TOKEN") + "&lang=" + dictionary + "&text=" + msg)

	if err != nil {
		log.Println("YANDEX_API_TOKEN ", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var translatedWord *DicResult
	var translatedResponse []string
	json.Unmarshal(body, &translatedWord)

	for _, def := range translatedWord.Def {
		for _, tr := range def.Tr {
			translatedResponse = append(translatedResponse, tr.Text)
		}
	}

	return translatedResponse
}
