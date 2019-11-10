package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
)

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

type UserState struct {
	value map[int]string
}