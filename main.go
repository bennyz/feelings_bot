package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"encoding/json"
	"strings"
	"io/ioutil"
	"bytes"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/peterhellberg/giphy"
)

type indicoResponse struct {
	Results map[string]interface{}
}

var indicoKey = os.Getenv("INDICO_KEY")

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func main() {
	botToken := os.Getenv("TELEGRAM_TOKEN")
	log.Printf("Key: [%s]", botToken)
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
		} else {
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
	text = resolveEmotion(text)

	g := giphy.DefaultClient
	s, _ := g.Search([]string{text})

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	imageIndex := r.Intn(len(s.Data))
	return s.Data[imageIndex].URL
}

func resolveEmotion(text string) string {
	textLower := strings.ToLower(text)
	if (textLower == "sad" || textLower == "happy") {
		return text;
	}

	return fetchEmotion(text)
}

func fetchEmotion(text string) string {
	const url = "https://apiv2.indico.io/emotion"
	var jsonStr = []byte(`{"data": "` + text + `", "top_n": 1}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-ApiKey", indicoKey)
	req.Header.Set("client-lib", "golang")
	req.Header.Set("version-number", "0.9.0")
	req.Header.Set("Accept", "text/plain")

	resp, err := netClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	log.Println("response status: ", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	var jsonResponse indicoResponse
	json.Unmarshal(body, &jsonResponse)

	results := jsonResponse.Results

	key := text
	for k := range results {
		key=k
		break
	}

	log.Printf("Detected emotion: %s", key)
	return key
}
