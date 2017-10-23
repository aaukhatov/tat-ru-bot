package main

import ("log"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
)

const webhook = "https://tat-ru-bot.herokuapp.com/"
const tat_dictionary = "http://tatpoisk.net/dict/"

func main() {
	translate("")
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
		log.Println("Received text", update.Message.Text)
		switch update.Message.Text {
		case "test":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "This translate a word")
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "The command not supported")
		}

		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func translate(msg string) {
	resp, err := http.Get("http://tatpoisk.net/dict/мин")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	z := html.NewTokenizer(resp.Body)
	fmt.Println(z.Next())
	/*for {
		tt := z.Next()
		fmt.Println(tt)
	}*/
	result := string(body)
	fmt.Println(result)

}
