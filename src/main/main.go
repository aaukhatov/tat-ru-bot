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
const telegram_token = "384640172:AAFOh_vCuFizDclHRxjpsY0SGoAtlsSCHs4"

func main() {
	var userState = make(map[int]string)
	telegram(userState)
}

func telegram(userState map[int]string) {
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI(telegram_token)
	if err != nil {
		log.Panic(err)
	}
	//var userState = make(map[int]string)
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
		log.Printf("Translated request: %s | command %s", update.Message.Text, update.Message.Command())

		var command, ok = userState[update.Message.From.ID]

		if !ok {
			log.Println("The user not found.", ok)
			if !update.Message.IsCommand() {
				log.Println("It's not command")
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите направление перевеода:\n" +
					"/rutat - русско-татарский\n/tatru - татарско-русский")
				bot.Send(msg)
				return
			}
			userState[update.Message.From.ID] = update.Message.Command()
		} else {
			log.Println("The user found.")
			if update.Message.IsCommand() &&
				update.Message.Command() != userState[update.Message.From.ID] {
				log.Println("The user comamnd update.")
				userState[update.Message.From.ID] = update.Message.Command()
				command = update.Message.Command()
			}
		}

		switch command {
		case "rutat":
			inputMsg := strings.Split(update.Message.Text, " ")
			// берем первое слово всегда
			translatedWord := translate(inputMsg[0], "ru-tt")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord,", "))
		case "tatru":
			inputMsg := strings.Split(update.Message.Text, " ")
			// берем первое слово всегда
			translatedWord := translate(inputMsg[0], "tt-ru")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord,", "))
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите направление перевода.")
		}

		bot.Send(msg)
	}
}

func translate(msg string, dictionary string) []string {

	resp, err := http.Get("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=" +
		yandex_api + "&lang=" + dictionary + "&text=" + msg)

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