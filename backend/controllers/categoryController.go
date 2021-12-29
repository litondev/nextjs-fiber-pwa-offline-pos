package controllers

import (
	"fmt"
	"math"
	"strconv"

	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func IndexCategory(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	page := c.Query("page", "1")
	new_page, _ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page, _ := strconv.Atoi(per_page)

	soft_deleted := c.Query("soft_deleted", "")

	search := c.Query("search", "")

	var resultCount int64

	var modelCount models.Category

	queryResultCount := database.Model(&modelCount)
	queryResultCount.Select("id")

	if search != "" {
		queryResultCount.Where("name LIKE ?", "%"+search+"%")
	}

	if soft_deleted != "" {
		queryResultCount.Where("deleted_at != ?", nil)
	}

	queryResultCount.Count(&resultCount)

	count_total_page := float64(resultCount) / float64(new_per_page)
	total_page := int(math.Ceil(count_total_page))
	limitStart := (new_page - 1) * new_per_page

	result := []map[string]interface{}{}

	var queryCount models.Category
	query := database.Model(&queryCount)
	query.Select("name", "id", "description")
	if search != "" {
		query.Where("name LIKE ?", "%"+search+"%")
	}
	if soft_deleted != "" {
		query.Where("deleted_at != ?", nil)
	}
	query.Order("id desc")
	query.Offset(limitStart)
	query.Limit(new_per_page)
	query.Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"data":       result,
		"per_page":   new_per_page,
		"total_page": total_page,
		"total_data": resultCount,
	})
}

func StoreCategory(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	category := new(requests.CategoryRequest)

	if errParser := c.BodyParser(category); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	category.ValidateData()

	errValidate := requests.ValidateStruct(category)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	categoryModel := models.Category{
		Name:        category.Name,
		Description: category.Description,
	}

	errCreateCategory := database.Create(&categoryModel).Error

	if errCreateCategory != nil {
		fmt.Println(errCreateCategory.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func ShowCategory(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resultData := map[string]interface{}{}

	queryData := database.Model(&models.Category{})
	queryData.Select("id", "name", "description")
	queryData.Where("id = ?", id)
	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": resultData,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	category := new(requests.CategoryRequest)

	if errParser := c.BodyParser(category); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	category.ValidateData()

	errValidate := requests.ValidateStruct(category)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	queryData := database.Model(&models.Category{})
	queryData.Select("name", "description")
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Updates(&models.Category{
		Name:        category.Name,
		Description: category.Description,
	})

	if queryData.Error != nil {
		fmt.Println(queryData.Error)

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func DestroyCategory(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	queryData := database.Model(&models.Category{})
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Delete(&models.Category{})

	if queryData.Error != nil {
		fmt.Println(queryData.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}
