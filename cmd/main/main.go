package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Project struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Annotation  string     `json:"annotation"`
	Time        string     `json:"time"`
	CardImgSrc  string     `json:"card_img_src"`
	BannerDesktop string   `json:"banner_desktop"`
	BannerMobile  string   `json:"banner_mobile"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	Categories  []Category `json:"categories"`
}

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserState struct {
	CurrentState string
	PreviousData string
	PreviousMarkup *tgbotapi.InlineKeyboardMarkup
}

var userStates = make(map[int64]UserState)

func main() {
	// Замените на токен вашего бота
	bot, err := tgbotapi.NewBotAPI("7918464444:AAFMleSbTQjwlE_ggIrL6bn5uXTABbv4Brg")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Обработка сообщений
			// Обработка команды /start
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
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
				bot.Send(msg)
			} else if update.Message.Text == "Проекты" {
				sendProjectList(bot, update.Message.Chat.ID)
			} else if update.Message.Text == "Категории" {
				sendCategoryList(bot, update.Message.Chat.ID)
			}
		} else if update.CallbackQuery != nil { // Обработка инлайн-кнопок
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
			}

			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
			if strings.HasPrefix(update.CallbackQuery.Data, "project_") {
				slug := strings.TrimPrefix(update.CallbackQuery.Data, "project_")
				url := fmt.Sprintf("http://185.178.47.249:1001/projects?slug=%s", slug)
				response, err := getProjectDetails(url)
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
					userStates[update.CallbackQuery.Message.Chat.ID] = UserState{
						CurrentState: "details",
						PreviousData: update.CallbackQuery.Message.Text,
						PreviousMarkup: update.CallbackQuery.Message.ReplyMarkup,
					}
				}
			} else if strings.HasPrefix(update.CallbackQuery.Data, "category_") {
				categoryID := strings.TrimPrefix(update.CallbackQuery.Data, "category_")
				url := fmt.Sprintf("http://185.178.47.249:1001/categories?id=%s", categoryID)
				response, err := getCategoryDetails(url)
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
					userStates[update.CallbackQuery.Message.Chat.ID] = UserState{
						CurrentState: "details",
						PreviousData: update.CallbackQuery.Message.Text,
						PreviousMarkup: update.CallbackQuery.Message.ReplyMarkup,
					}
				}
			} else if update.CallbackQuery.Data == "back" {
				userState, exists := userStates[update.CallbackQuery.Message.Chat.ID]
				if exists && userState.CurrentState == "details" {
					msg.Text = userState.PreviousData
					msg.ReplyMarkup = userState.PreviousMarkup
					delete(userStates, update.CallbackQuery.Message.Chat.ID)
				} else {
					msg.Text = "Невозможно вернуться назад."
				}
			}
			bot.Send(msg)
		}
	}
}

func sendProjectList(bot *tgbotapi.BotAPI, chatID int64) {
	url := "http://185.178.47.249:1001/projects"
	data, statusCode, err := getDataFromURL(url)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении данных.")
		bot.Send(msg)
		return
	}

	if statusCode != http.StatusOK {
		log.Printf("Server returned status code: %d", statusCode)
		msg := tgbotapi.NewMessage(chatID, "Сервер вернул ошибку.")
		bot.Send(msg)
		return
	}

	var projects []Project
	err = json.Unmarshal(data, &projects)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при обработке данных.")
		bot.Send(msg)
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
	bot.Send(msg)
}

func sendCategoryList(bot *tgbotapi.BotAPI, chatID int64) {
	url := "http://185.178.47.249:1001/categories"
	data, statusCode, err := getDataFromURL(url)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении данных.")
		bot.Send(msg)
		return
	}

	if statusCode != http.StatusOK {
		log.Printf("Server returned status code: %d", statusCode)
		msg := tgbotapi.NewMessage(chatID, "Сервер вернул ошибку.")
		bot.Send(msg)
		return
	}

	var categories []Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при обработке данных.")
		bot.Send(msg)
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
	bot.Send(msg)
}

func getProjectDetails(url string) (string, error) {
	data, statusCode, err := getDataFromURL(url)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned status code: %d", statusCode)
	}

	var project Project
	err = json.Unmarshal(data, &project)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("ID: %s\nSlug: %s\nНазвание проекта: %s\nОписание проекта: %s\nАннотация: %s\nВремя затраченное на проект: %s\nКатегории: %s\nКарточка: %s\nDesktop изображение: %s\nMobile изображение: %s\nCreatedAt: %s\nUpdatedAt: %s\n\n",
		project.ID, project.Slug, project.Name, project.Description, project.Annotation, project.Time, getCategoriesString(project.Categories), project.CardImgSrc, project.BannerDesktop, project.BannerMobile, project.CreatedAt, project.UpdatedAt)

	return response, nil
}

func getCategoryDetails(url string) (string, error) {
	data, statusCode, err := getDataFromURL(url)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned status code: %d", statusCode)
	}

	var category Category
	err = json.Unmarshal(data, &category)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("ID: %s\nНазвание: %s\nCreatedAt: %s\nUpdatedAt: %s\n\n",
		category.ID, category.Name, category.CreatedAt, category.UpdatedAt)

	return response, nil
}

func getDataFromURL(url string) ([]byte, int, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, response.StatusCode, nil
}

func getCategoriesString(categories []Category) string {
	var categoryNames []string
	for _, category := range categories {
		categoryNames = append(categoryNames, category.Name)
	}
	return fmt.Sprintf("[%s]", strings.Join(categoryNames, ", "))
}