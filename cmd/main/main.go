package main

import (
	"log"
	"github.com/petersizovdev/tg-pkadmin/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7918464444:AAFMleSbTQjwlE_ggIrL6bn5uXTABbv4Brg")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	botService := services.NewBotService(bot)

	for update := range updates {
		botService.HandleUpdate(update)
	}
}