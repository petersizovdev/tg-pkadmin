package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/petersizovdev/tg-pkadmin/internal/models"
)

func GetCategoriesString(categories []models.Category) string {
	var categoryNames []string
	for _, category := range categories {
		categoryNames = append(categoryNames, category.Name)
	}
	return fmt.Sprintf("[%s]", strings.Join(categoryNames, ", "))
}

func GetTelegramToken() string {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN environment variable is not set")
	}
	return token
}