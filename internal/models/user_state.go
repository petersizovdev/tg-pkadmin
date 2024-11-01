package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UserState struct {
	CurrentState   string
	PreviousData   string
	PreviousMarkup *tgbotapi.InlineKeyboardMarkup
}