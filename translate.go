package main

import (
	"bytes"
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

// TranslationResult is represents Yandex Translate response
type TranslationResult struct {
	Code int    `json:"code"`
	Lang string `json:"lang"`
	Text []string `json:"text"`
	Message string `json: "message"`
}

func translate(msg string, dictionary string) string {
	text := []byte("text=" + msg)
	
	resp, err := http.Post("https://translate.yandex.net/api/v1.5/tr.json/translate?key="+
		os.Getenv("YANDEX_API_TOKEN")+"&lang="+dictionary, "application/x-www-form-urlencoded", bytes.NewBuffer(text))

	if err != nil {
		log.Println("YANDEX_API_TOKEN ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Could not read a response", err)
	}
	
	var translatedWord *TranslationResult

	err = json.Unmarshal(body, &translatedWord)
	if err != nil {
		log.Println("[ERROR] Could not be unmarshalled ", err)
	}

	return translatedWord.Text[0]
}
