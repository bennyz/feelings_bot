package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	botToken := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/start" {
			bot.Send(getReplyMessage(update.Message.Chat.ID, update.Message.Text))
		}

	}
}

func getReplyMessage(chatID int64, message string) tgbotapi.MessageConfig {
	commands := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Happy"),
			tgbotapi.NewKeyboardButton("Sad"),
		},
	)

	msg := tgbotapi.NewMessage(chatID, "How are you feeling?")
	msg.ReplyMarkup = commands
	return msg
}
