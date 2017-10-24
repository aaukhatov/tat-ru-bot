package main

import ("log"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
)

const webhook = "https://tat-ru-bot.herokuapp.com/"
const yandex_api = "dict.1.1.20171024T175215Z.d79c6c40e3a0bf31.0f44341ac31440368c75d3e143c641ab1a7acec6"

func main() {
	telegram()
}

func telegram() {
	log.Print("I'm a bot!")
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("384640172:AAFOh_vCuFizDclHRxjpsY0SGoAtlsSCHs4")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhook))
	if err != nil {
		log.Fatal(err)
	}
	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+port, nil)

	for update := range updates {
		var msg tgbotapi.MessageConfig
		log.Println("Translated request:", update.Message.Text)
		inputMsg := strings.Split(update.Message.Text, " ")
		// берем первое слово всегда
		translatedWord := translate(inputMsg[0])
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord,","))
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func translate(msg string) []string {

	resp, err := http.Get("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=" +
		yandex_api + "&lang=ru-tt&text=" + msg)

	if err != nil {
		log.Println(err)
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