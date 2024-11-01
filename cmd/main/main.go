package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/petersizovdev/tg-pkadmin/internal/services"
	"github.com/petersizovdev/tg-pkadmin/internal/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	telegramToken := utils.GetTelegramToken()

	bot, err := tgbotapi.NewBotAPI(telegramToken)
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