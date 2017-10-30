package main

import ("log"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"sync"
)

const webhook = "https://tat-ru-bot.herokuapp.com/"
const yandexApi = "dict.1.1.20171024T175215Z.d79c6c40e3a0bf31.0f44341ac31440368c75d3e143c641ab1a7acec6"
const telegramToken = "384640172:AAFOh_vCuFizDclHRxjpsY0SGoAtlsSCHs4"
const helpMessage = "Укажите направление перевода:\n" +
	"/rutat - русско-татарский\n/tatru - татарско-русский"
const commandRuTat = "rutat"
const commandTatRu = "tatru"
const wordNotFound = "Слово не найдено.\nВозможно стоит переключить словарь?\n"+
	"/rutat - русско-татарский\n/tatru - татарско-русский"
const tatRu = "tt-ru"
const ruTat = "ru-tt"

func main() {
	var userState *UserState
	telegram(userState)
}

func telegram(userState *UserState) {
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
func executeCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, userState *UserState) {
	var msg tgbotapi.MessageConfig
	log.Println("User lock", userState)
	userState.m.Lock()
	var command, ok = userState.value[update.Message.From.ID]
	userState.m.Unlock()
	log.Println("User unlock", userState)

	msg, command = preDefineCommand(ok, update, msg, userState, command)

	switch command {
	case commandRuTat:
		inputMsg := strings.Split(update.Message.Text, " ")
		// берем первое слово всегда
		translatedWord := translate(inputMsg[0], ruTat)
		if len(translatedWord) == 0 {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, wordNotFound)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord, ", "))
		}
	case commandTatRu:
		inputMsg := strings.Split(update.Message.Text, " ")
		// берем первое слово всегда
		translatedWord := translate(inputMsg[0], tatRu)
		if len(translatedWord) == 0 {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, wordNotFound)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(translatedWord, ", "))
		}
	case "start":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	}
	bot.Send(msg)
}
func preDefineCommand(ok bool, update tgbotapi.Update, msg tgbotapi.MessageConfig,
	userState *UserState, command string) (tgbotapi.MessageConfig, string) {
	if !ok {
		if !update.Message.IsCommand() {
			log.Println("It's not command")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
		} else {
			userState.m.Lock()
			userState.value[update.Message.From.ID] = update.Message.Command()
			userState.m.Unlock()

			if update.Message.Command() == commandRuTat || update.Message.Command() == commandTatRu {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите слово для перевода")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
			}
		}
	} else if update.Message.IsCommand() {
		userState.m.RLock()
		newCommand := userState.value[update.Message.From.ID]
		userState.m.RUnlock()

		if update.Message.Command() != newCommand &&
			update.Message.Command() == commandRuTat || update.Message.Command() == commandTatRu {
			log.Println("The user comamnd update.")
			userState.m.Lock()
			userState.value[update.Message.From.ID] = update.Message.Command()
			userState.m.Unlock()
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
		Tr   []struct {
			Text string `json:"text"`
			Pos  string `json:"pos"`
			Mean []struct {
				Text string `json:"text"`
			} `json:"mean"`
		} `json:"tr"`
	} `json:"def"`
}

type UserState struct {
	m sync.RWMutex
	value map[int]string
}