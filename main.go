package main

import (
	"flag"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"os"
)

const webhook = "https://tat-ru-bot.herokuapp.com/"
const helpMessage = "Укажите направление перевода:\n" +
	"/rutat - русско-татарский\n/tatru - татарско-русский"
const commandRuTat = "rutat"
const commandTatRu = "tatru"
const wordNotFound = "Слово не найдено в словаре"
const tatRu = "tt-ru"
const ruTat = "ru-tt"

func main() {
	var isHeroku = flag.Bool("heroku", false, "Heroku mode.")
	flag.Parse()

	user := make(map[int]string)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic("TELEGRAM_API_TOKEN ", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	var updates tgbotapi.UpdatesChannel
	
	if *isHeroku {
		log.Printf("Bot has been started on Heroku")
		port := os.Getenv("PORT")
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhook))
		if err != nil {
			log.Fatal(err)
		}

		go http.ListenAndServe(":"+port, nil)
		updates = bot.ListenForWebhook("/")
	} else {
		log.Printf("Bot has been started on Local")
		bot.Debug = true
		bot.RemoveWebhook()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			log.Println(err)
		}
	}

	handleUpdates(updates, bot, user)
}

func handleUpdates(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, user map[int]string) {
	for {
		select {
		case update := <-updates:
			go executeCommand(update, bot, user)
		}
	}
}
