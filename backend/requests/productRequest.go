package requests

import (
	"api-gofiber/pos/models"
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type ProductRequest struct {
	Name        string `form:"name" json:"name" xml:"name" validate:"required"`
	Code        string `form:"code" json:"code" xml:"code" validate:"required"`
	Description *string
	Stock       string `form:"stock" json:"stock" xml:"stock" validate:"required,number"`
	Price       string `form:"price" json:"price" xml:"price" validate:"required,numeric"`
	CategoryID  string `form:"category_id" json:"category_id" xml:"category_id"  validate:"required,number"`
}

func (requestData *ProductRequest) ValidateData() error {
	requestData.Name = html.EscapeString(strings.Trim(requestData.Name, " "))
	requestData.Code = html.EscapeString(strings.Trim(requestData.Code, " "))

	if requestData.Description != nil {
		des := html.EscapeString(strings.Trim(*requestData.Description, " "))
		requestData.Description = &des
	}

	return nil
}

func (requestData *ProductRequest) ValidatePhoto() error {
	return nil
}

func (requestData *ProductRequest) ValidateExistsCategory(database *gorm.DB) error {
	resultData := map[string]interface{}{}

	queryData := database.Model(&models.Category{})
	queryData.Select("id")
	queryData.Where("id = ?", requestData.CategoryID)
	queryData.First(&resultData)

	if len(resultData) == 0 {
		return errors.New("Category Tidak Ditemukan")
	}

	return nil
}
