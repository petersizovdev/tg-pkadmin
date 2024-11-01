package services

import (
	"fmt"
	"log"
	"strings"
	"github.com/petersizovdev/tg-pkadmin/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
	bot *tgbotapi.BotAPI
	dataService *DataService
	userStates map[int64]models.UserState
}

func NewBotService(bot *tgbotapi.BotAPI) *BotService {
	return &BotService{
		bot: bot,
		dataService: NewDataService(),
		userStates: make(map[int64]models.UserState),
	}
}

func (s *BotService) HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil { // Обработка сообщений
		s.handleMessage(update.Message)
	} else if update.CallbackQuery != nil { // Обработка инлайн-кнопок
		s.handleCallbackQuery(update.CallbackQuery)
	}
}

func (s *BotService) handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		s.handleCommand(message)
	} else if message.Text == "Проекты" {
		s.sendProjectList(message.Chat.ID)
	} else if message.Text == "Категории" {
		s.sendCategoryList(message.Chat.ID)
	}
}

func (s *BotService) handleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	switch message.Command() {
	case "start":
		msg.Text = "Привет! Выбери опцию:"
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Проекты"),
				tgbotapi.NewKeyboardButton("Категории"),
			),
		)
	default:
		msg.Text = "Я не знаю такой команды :("
	}
	s.bot.Send(msg)
}

func (s *BotService) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	if _, err := s.bot.Request(callback); err != nil {
		log.Println(err)
	}

	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "")
	if strings.HasPrefix(callbackQuery.Data, "project_") {
		slug := strings.TrimPrefix(callbackQuery.Data, "project_")
		url := fmt.Sprintf("http://185.178.47.249:1001/projects?slug=%s", slug)
		response, err := s.dataService.GetProjectDetails(url)
		if err != nil {
			log.Println(err)
			msg.Text = "Ошибка при получении данных."
		} else {
			msg.Text = response
			msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
					},
				},
			}
			s.userStates[callbackQuery.Message.Chat.ID] = models.UserState{
				CurrentState: "details",
				PreviousData: callbackQuery.Message.Text,
				PreviousMarkup: callbackQuery.Message.ReplyMarkup,
			}
		}
	} else if strings.HasPrefix(callbackQuery.Data, "category_") {
		categoryID := strings.TrimPrefix(callbackQuery.Data, "category_")
		url := fmt.Sprintf("http://185.178.47.249:1001/categories?id=%s", categoryID)
		response, err := s.dataService.GetCategoryDetails(url)
		if err != nil {
			log.Println(err)
			msg.Text = "Ошибка при получении данных."
		} else {
			msg.Text = response
			msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
					},
				},
			}
			s.userStates[callbackQuery.Message.Chat.ID] = models.UserState{
				CurrentState: "details",
				PreviousData: callbackQuery.Message.Text,
				PreviousMarkup: callbackQuery.Message.ReplyMarkup,
			}
		}
	} else if callbackQuery.Data == "back" {
		userState, exists := s.userStates[callbackQuery.Message.Chat.ID]
		if exists && userState.CurrentState == "details" {
			msg.Text = userState.PreviousData
			msg.ReplyMarkup = userState.PreviousMarkup
			delete(s.userStates, callbackQuery.Message.Chat.ID)
		} else {
			msg.Text = "Невозможно вернуться назад."
		}
	}
	s.bot.Send(msg)
}

func (s *BotService) sendProjectList(chatID int64) {
	url := "http://185.178.47.249:1001/projects"
	projects, err := s.dataService.GetProjects(url)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении данных.")
		s.bot.Send(msg)
		return
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, project := range projects {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(project.Name, fmt.Sprintf("project_%s", project.Slug)),
		))
	}

	msg := tgbotapi.NewMessage(chatID, "Выберите проект:")
	msg.ReplyMarkup = inlineKeyboard
	s.bot.Send(msg)
}

func (s *BotService) sendCategoryList(chatID int64) {
	url := "http://185.178.47.249:1001/categories"
	categories, err := s.dataService.GetCategories(url)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении данных.")
		s.bot.Send(msg)
		return
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, category := range categories {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, fmt.Sprintf("category_%s", category.ID)),
		))
	}

	msg := tgbotapi.NewMessage(chatID, "Выберите категорию:")
	msg.ReplyMarkup = inlineKeyboard
	s.bot.Send(msg)
}