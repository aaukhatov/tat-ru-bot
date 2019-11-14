package main

import (
	"flag"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"os"
)

const webhook = "https://tat-ru-bot.herokuapp.com/"

func main() {
	var isHeroku = flag.Bool("heroku", false, "Heroku mode.")
	flag.Parse()
	translationChat := newChat()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic("TELEGRAM_API_TOKEN ", err)
	}

	log.Printf("[INFO] Authorized on account %s", bot.Self.UserName)
	var updates tgbotapi.UpdatesChannel

	if *isHeroku {
		log.Printf("[INFO] Bot has been started on Heroku")
		port := os.Getenv("PORT")
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhook))
		if err != nil {
			log.Fatal(err)
		}

		go http.ListenAndServe(":"+port, nil)
		updates = bot.ListenForWebhook("/")
	} else {
		log.Printf("[DEBUG] Bot has been started on Local")
		bot.RemoveWebhook()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			log.Println(err)
		}
	}

	translationChat.run(bot, updates)
}
