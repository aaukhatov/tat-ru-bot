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

func executeCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, translationChat *chat) {
	var msg tgbotapi.MessageConfig

	if update.Message.Command() == cmdAbout {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, aboutMsg)
		msg.ParseMode = "Markdown"
	}

	if update.Message.Command() == cmdStart {
		translationChat.userState[update.Message.From.ID] = cmdStart
		msg = tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf(helloMsg, update.Message.From.FirstName, update.Message.From.LastName))
		msg.ParseMode = "Markdown"
	}

	if update.Message.Text == cmdRuTat {
		translationChat.userState[update.Message.From.ID] = cmdRuTat
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "✏️ Введите слово для перевода 😊")
	}

	if update.Message.Text == cmdTatRu {
		translationChat.userState[update.Message.From.ID] = cmdTatRu
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "✏️ Тәрҗемә итергә сүзеңне яз 😊")
	}

	command := translationChat.userState[update.Message.From.ID]

	if update.Message.Text != command {
		switch command {
		case cmdRuTat:
			translatedWord := translate(update.Message.Text, ruTat)
			msg = newTelegramMessage(update.Message.Chat.ID, translatedWord)

		case cmdTatRu:
			translatedWord := translate(update.Message.Text, tatRu)
			msg = newTelegramMessage(update.Message.Chat.ID, translatedWord)
		}
	}

	translationChat.botReponse <- msg
}

func newTelegramMessageByArray(chatID int64, translatedWord []string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	if len(translatedWord) == 0 {
		msg = tgbotapi.NewMessage(chatID, wordNotFoundMsg)
	} else {
		msg = tgbotapi.NewMessage(chatID, strings.Join(translatedWord, ", "))
	}

	return msg
}

func newTelegramMessage(chatID int64, translatedWord string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	if len(translatedWord) == 0 {
		msg = tgbotapi.NewMessage(chatID, wordNotFoundMsg)
	} else {
		msg = tgbotapi.NewMessage(chatID, translatedWord)
	}
	return msg
}

type chat struct {
	botReponse chan tgbotapi.MessageConfig
	userState  map[int]string
}

func newChat() *chat {
	return &chat{
		botReponse: make(chan tgbotapi.MessageConfig),
		userState: make(map[int]string),
	}
}

func (c *chat) run(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case userMsg := <-updates:
			log.Printf("[INFO] Received user message '%s'", userMsg.Message.Text)
			go executeCommand(userMsg, bot, c)

		case botReponse := <-c.botReponse:
			log.Printf("[INFO] Sent a reponse '%s'", botReponse.Text)
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
