package main

import (
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strconv"
)

func handle(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		handleMessages(update, bot)
	}
	if update.CallbackQuery != nil {
		handleCallbackQuery(update, bot)
	}
}

func handleMessages(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message.Chat.Type != "private" {
		log.Printf("Received message from non-private chat: %v", update.Message.Chat.Type)
		return
	}

	user := User{
		ID:     update.Message.From.ID,
		ChatID: update.Message.Chat.ID,
	}

	// Handle commands
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			startCommand(user, bot)
		}
	}

}

func startCommand(user User, bot *tgbotapi.BotAPI) {
	startMsg := tgbotapi.NewMessage(user.ChatID, "Hey! Choose your lang:")
	startMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("English", "lang_en"),
			tgbotapi.NewInlineKeyboardButtonData("Русский", "lang_ru"),
		),
	)
	bot.Send(startMsg)
}

func handleCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := User{
		ID:     update.CallbackQuery.From.ID,
		ChatID: update.CallbackQuery.Message.Chat.ID,
	}

	callback := update.CallbackQuery

	var channelsMessage tgbotapi.MessageConfig
	var generateMessage tgbotapi.MessageConfig

	switch callback.Data {
	case "lang_en":
		channelsMessage = tgbotapi.NewMessage(user.ChatID, "Please subscribe to these channels:")
		channelsMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(enChannels()...)
		generateMessage = tgbotapi.NewMessage(user.ChatID, "After subscribing, launch the Mini-App by clicking the blue «Get codes» button at the bottom")

	case "lang_ru":
		channelsMessage = tgbotapi.NewMessage(user.ChatID, "Пожалуйста, подпишитесь на эти каналы:")
		channelsMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(ruChannels()...)
		generateMessage = tgbotapi.NewMessage(user.ChatID, "После подписки запустите Mini-App по синей кнопке «Get codes» внизу")
	}

	bot.Send(channelsMessage)
	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	bot.Send(callbackResponse)
	bot.Send(generateMessage)
}

func enChannels() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{{
		tgbotapi.NewInlineKeyboardButtonURL("TON Games", "https://t.me/+y6-lsFUJR305ZTli")}, {
		tgbotapi.NewInlineKeyboardButtonURL("Hamster Playground Codes", "https://t.me/+pGzy5ciQuUQ2YTAy")},
	}
}

func ruChannels() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{{
		tgbotapi.NewInlineKeyboardButtonURL("TON Игры", "https://t.me/+nAh9Lkl15RMyMWYy")}, {
		tgbotapi.NewInlineKeyboardButtonURL("Чат", "https://t.me/+AhrrvNt9WepkOGQy")}, {
		tgbotapi.NewInlineKeyboardButtonURL("Hamster Playground Codes", "https://t.me/+pGzy5ciQuUQ2YTAy")}, {
		tgbotapi.NewInlineKeyboardButtonURL("Бизон Хиггса", "https://t.me/+7yw-yZPsCp5iMzg6")},
	}

}

func checkSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "userId and userLang are required", http.StatusBadRequest)
		return
	}
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)

	enSubscriptionChannels := getEnChannels()
	ruSubscriptionChannels := getRuChannels()

	enSubscribed := true
	for _, channel := range enSubscriptionChannels {
		if !checkSubscriptionChannel(channel, userIdInt) {
			enSubscribed = false
		}
	}

	ruSubscribed := true
	for _, channel := range ruSubscriptionChannels {
		if !checkSubscriptionChannel(channel, userIdInt) {
			ruSubscribed = false
		}
	}

	subscribed := false
	if ruSubscribed || enSubscribed {
		subscribed = true
	}

	response := SubscriptionCheckResponse{
		Subscribed: subscribed,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkSubscriptionChannel(channelChatID int64, chatID int64) bool {
	memberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: channelChatID,
			UserID: chatID,
		},
	}

	chatMember, _ := bot.GetChatMember(memberConfig)

	return chatMember.Status == "member" || chatMember.Status == "administrator" || chatMember.Status == "creator"
}

func getEnChannels() []int64 {
	var channelIDs []int64

	// Define environment variable keys
	keys := []string{
		"TELEGRAM_PLAYGROUND_CHANNEL_ID",
		"TELEGRAM_EN_CHANNEL_ID",
	}

	// Populate map with channel IDs from environment variables
	for _, key := range keys {
		id := os.Getenv(key)
		channelID, _ := strconv.ParseInt(id, 10, 64)
		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs
}

func getRuChannels() []int64 {
	var channelIDs []int64

	// Define environment variable keys
	keys := []string{
		"TELEGRAM_MAIN_CHANNEL_ID",
		"TELEGRAM_PLAYGROUND_CHANNEL_ID",
		"TELEGRAM_MAIN_CHAT_ID",
		"TELEGRAM_BIZON_CHANNEL_ID",
	}

	// Populate map with channel IDs from environment variables
	for _, key := range keys {
		id := os.Getenv(key)
		channelID, _ := strconv.ParseInt(id, 10, 64)
		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs
}
