package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
)

const webhook = "https://tat-ru-bot.herokuapp.com/"
const yandexApi = "dict.1.1.20171024T175215Z.d79c6c40e3a0bf31.0f44341ac31440368c75d3e143c641ab1a7acec6"
const telegramToken = "384640172:AAFOh_vCuFizDclHRxjpsY0SGoAtlsSCHs4"
const helpMessage = "Укажите направление перевода:\n" +
	"/rutat - русско-татарский\n/tatru - татарско-русский"
const commandRuTat = "rutat"
const commandTatRu = "tatru"
const wordNotFound = "Слово не найдено в словаре"
const tatRu = "tt-ru"
const ruTat = "ru-tt"

func main() {
	user := make(map[int]string)
	telegram(user)
}

func telegram(userState map[int]string) {
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
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
		go executeCommand(update, bot, userState)
	}
}
func executeCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, userState map[int]string) {
	var msg tgbotapi.MessageConfig
	log.Println("User lock", userState)
	var command, ok = userState[update.Message.From.ID]
	log.Println("User unlock", userState)

	msg, command = defineCommand(ok, update, msg, userState, command)

	switch command {
	case commandRuTat:
		if !update.Message.IsCommand() {
			inputMsg := strings.Split(update.Message.Text, " ")
			// берем первое слово всегда
			translatedWord := translate(inputMsg[0], ruTat)
			if len(translatedWord) == 0 {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, wordNotFound)
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord, ", "))
			}
		}
	case commandTatRu:
		if !update.Message.IsCommand() {
			inputMsg := strings.Split(update.Message.Text, " ")
			// берем первое слово всегда
			translatedWord := translate(inputMsg[0], tatRu)
			if len(translatedWord) == 0 {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, wordNotFound)
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord, ", "))
			}
		}
	case "start":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	}
	bot.Send(msg)
}
func defineCommand(ok bool, update tgbotapi.Update, msg tgbotapi.MessageConfig,
	userState map[int]string, command string) (tgbotapi.MessageConfig, string) {
	if !ok {
		if !update.Message.IsCommand() {
			log.Println("It's not command")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
		} else {
			userState[update.Message.From.ID] = update.Message.Command()

			if update.Message.Command() == commandRuTat || update.Message.Command() == commandTatRu {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите слово для перевода")
			}
		}
	} else if update.Message.IsCommand() {
		newCommand := userState[update.Message.From.ID]

		if update.Message.Command() != newCommand &&
			update.Message.Command() == commandRuTat || update.Message.Command() == commandTatRu {
			log.Println("The user comamnd update.")
			userState[update.Message.From.ID] = update.Message.Command()
			command = update.Message.Command()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите слово для перевода")
		}
	}
	return msg, command
}

func translate(msg string, dictionary string) []string {

	resp, err := http.Get("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=" +
		yandexApi + "&lang=" + dictionary + "&text=" + msg)

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
		Tr []struct {
			Text string `json:"text"`
			Pos  string `json:"pos"`
			Mean []struct {
				Text string `json:"text"`
			} `json:"mean"`
		} `json:"tr"`
	} `json:"def"`
}

type UserState struct {
	value map[int]string
}
