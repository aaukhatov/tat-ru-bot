package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const cmdRuTat = "русско-татарский"
const cmdTatRu = "татарча-русча"
const cmdAbout = "about"
const cmdStart = "start"
const aboutMsg = `Arthur Aukhatov - @aaukhatov ©️`
const wordNotFoundMsg = "🤷‍♀️ Слово не найдено в словаре ☹️"
const helloMsg = `Привет, %s %s! 😊🙌
Я бот умеющий делать перевод слов.
- татарско-русский
- русско-татарский

Сначала выберите направление перевода, затем пишите слово.
`

func executeCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, userState map[int]string, c *chat) {
	var msg tgbotapi.MessageConfig

	if update.Message.Command() == cmdAbout {
		userState[update.Message.From.ID] = cmdAbout
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, aboutMsg)
		msg.ParseMode = "Markdown"
	}

	if update.Message.Command() == cmdStart {
		userState[update.Message.From.ID] = cmdStart
		msg = tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf(helloMsg, update.Message.From.FirstName, update.Message.From.LastName))
		msg.ParseMode = "Markdown"
	}

	if update.Message.Text == cmdRuTat {
		userState[update.Message.From.ID] = cmdRuTat
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "✏️ Введите слово для перевода 😊")
	}

	if update.Message.Text == cmdTatRu {
		userState[update.Message.From.ID] = cmdTatRu
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "✏️ Тәрҗемә итергә сүзеңне яз 😊")
	}

	command := userState[update.Message.From.ID]

	if update.Message.Text != command {
		switch command {
		case cmdRuTat:
			inputMsg := strings.Split(update.Message.Text, " ")
			// always get the first word
			translatedWord := translate(inputMsg[0], ruTat)
			msg = newTelegramMessage(update.Message.Chat.ID, translatedWord)

		case cmdTatRu:
			inputMsg := strings.Split(update.Message.Text, " ")
			// always get the first word
			translatedWord := translate(inputMsg[0], tatRu)
			msg = newTelegramMessage(update.Message.Chat.ID, translatedWord)
		}
	}

	c.botReponse <- msg
}

func newTelegramMessage(chatID int64, translatedWord []string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	if len(translatedWord) == 0 {
		msg = tgbotapi.NewMessage(chatID, wordNotFoundMsg)
	} else {
		msg = tgbotapi.NewMessage(chatID, strings.Join(translatedWord, ", "))
	}

	return msg
}

type userState struct {
	value map[int]string
}

type chat struct {
	userMsg    chan string
	botReponse chan tgbotapi.MessageConfig
}

func newChat() *chat {
	return &chat{
		userMsg:    make(chan string),
		botReponse: make(chan tgbotapi.MessageConfig),
	}
}

func (c *chat) run(bot *tgbotapi.BotAPI) {
	for {
		select {
		case userMsg := <-c.userMsg:
			log.Printf("[INFO] userMsg: %s", userMsg)

		case botReponse := <-c.botReponse:
			log.Printf("[INFO] botReponse: %s", botReponse.Text)
			botReponse.ReplyMarkup = tgbotapi.NewReplyKeyboard(replyButton())
			bot.Send(botReponse)
		}
	}
}

func replyButton() []tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(cmdTatRu),
		tgbotapi.NewKeyboardButton(cmdRuTat))
}
