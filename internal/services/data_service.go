package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/petersizovdev/tg-pkadmin/internal/models"
	"github.com/petersizovdev/tg-pkadmin/internal/utils"
)

type DataService struct{}

func NewDataService() *DataService {
	return &DataService{}
}

func (s *DataService) GetProjects(url string) ([]models.Project, error) {
	data, statusCode, err := s.getDataFromURL(url)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned status code: %d", statusCode)
	}

	var projects []models.Project
	err = json.Unmarshal(data, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *DataService) GetCategories(url string) ([]models.Category, error) {
	data, statusCode, err := s.getDataFromURL(url)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned status code: %d", statusCode)
	}

	var categories []models.Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *DataService) GetProjectDetails(url string) (string, error) {
	data, statusCode, err := s.getDataFromURL(url)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned status code: %d", statusCode)
	}

	var project models.Project
	err = json.Unmarshal(data, &project)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("ID: %s\nSlug: %s\nНазвание проекта: %s\nОписание проекта: %s\nАннотация: %s\nВремя затраченное на проект: %s\nКатегории: %s\nКарточка: %s\nDesktop изображение: %s\nMobile изображение: %s\nCreatedAt: %s\nUpdatedAt: %s\n\n",
		project.ID, project.Slug, project.Name, project.Description, project.Annotation, project.Time, utils.GetCategoriesString(project.Categories), project.CardImgSrc, project.BannerDesktop, project.BannerMobile, project.CreatedAt, project.UpdatedAt)

	return response, nil
}

func (s *DataService) GetCategoryDetails(url string) (string, error) {
	data, statusCode, err := s.getDataFromURL(url)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned status code: %d", statusCode)
	}

	var category models.Category
	err = json.Unmarshal(data, &category)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("ID: %s\nНазвание: %s\nCreatedAt: %s\nUpdatedAt: %s\n\n",
		category.ID, category.Name, category.CreatedAt, category.UpdatedAt)

	return response, nil
}

func (s *DataService) getDataFromURL(url string) ([]byte, int, error) {
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