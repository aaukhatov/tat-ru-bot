package main

import ("log"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
)

const webhook = "asd"

func main() {
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
	go http.ListenAndServe(":" + port, nil)

	for update:= range updates {
		var msg tgbotapi.MessageConfig
		log.Println("Received text", update.Message.Text)
		switch update.Message.Text {
		case "get translate":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "This translate a word")
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "The command not supported")
		}

		bot.Send(msg)
	}
}
