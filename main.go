package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/peterhellberg/giphy"
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
		message := update.Message.Text
		chatID := update.Message.Chat.ID
		if message == "/start" {
			bot.Send(getReplyMessage(chatID,
				message,
				"How are you feeling?"))
		} else if message == "Sad" || message == "Happy" {
			bot.Send(getReplyMessage(chatID,
				message,
				getGIF(message)))
		}

	}
}

func getReplyMessage(chatID int64, message string, reply string) tgbotapi.MessageConfig {
	commands := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Happy"),
			tgbotapi.NewKeyboardButton("Sad"),
		},
	)

	msg := tgbotapi.NewMessage(chatID, reply)
	msg.ReplyMarkup = commands
	return msg
}

func getGIF(text string) string {
	g := giphy.DefaultClient
	s, _ := g.Search([]string{text})
	return s.Data[0].URL
}
