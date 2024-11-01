package utils

import (
	"fmt"
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