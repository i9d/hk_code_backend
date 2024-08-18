package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var _mainChannelID = os.Getenv("TELEGRAM_MAIN_CHANNEL_ID")
var _enChannelID = os.Getenv("TELEGRAM_EN_CHANNEL_ID")
var _bizonChannelID = os.Getenv("TELEGRAM_BIZON_CHANNEL_ID")
var _playgroundChannelID = os.Getenv("TELEGRAM_PLAYGROUND_CHANNEL_ID")
var _mainChatID = os.Getenv("TELEGRAM_MAIN_CHAT_ID")

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			handle(update, bot)
		}
	}()

	http.HandleFunc("/check-subscription", checkSubscriptionHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var bot *tgbotapi.BotAPI

// SubscriptionCheckResponse Struct for the JSON response
type SubscriptionCheckResponse struct {
	Subscribed bool `json:"subscribed"`
}
