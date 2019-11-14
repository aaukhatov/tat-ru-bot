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
	client := &http.Client{}
	yandexAPI := os.Getenv("YANDEX_API_TOKEN")

	if yandexAPI == "" {
		log.Println("[ERROR] YANDEX_API_TOKEN is requered parameter")
	}

	req, err := http.NewRequest("POST", "https://translate.yandex.net/api/v1.5/tr.json/translate?key="+
	yandexAPI+"&lang="+dictionary, bytes.NewBuffer(text))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		log.Println("[ERROR] Something went wrong", err)
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
