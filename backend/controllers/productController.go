package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Category struct {
	ID   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type Product struct {
	ID          int      `gorm:"column:id" json:"id"`
	Name        string   `gorm:"column:name" json:"name"`
	Code        string   `gorm:"column:code" json:"code"`
	Description *string  `gorm:"column:description" json:"description"`
	Stock       int      `gorm:"column:stock" json:"stock"`
	Price       float32  `gorm:"column:price" json:"price"`
	Photo       *string  `gorm:"column:photo" json:"photo"`
	CategoryID  *uint    `gorm:"column:category_id" json:"category_id"`
	Category    Category `json:"category"`
	DeletedAt   *gorm.DeletedAt
}

func IndexProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	page := c.Query("page", "1")
	new_page, _ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page, _ := strconv.Atoi(per_page)

	soft_deleted := c.Query("soft_deleted", "")

	search := c.Query("search", "")

	var resultCount int64

	var modelCount models.Product

	queryResultCount := database.Model(&modelCount)
	queryResultCount.Select("id")

	if search != "" {
		queryResultCount.Where(
			"(name LIKE ? OR CAST(price as text) LIKE ?) OR code LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if soft_deleted != "" {
		queryResultCount.Unscoped()
	}

	queryResultCount.Count(&resultCount)

	count_total_page := float64(resultCount) / float64(new_per_page)
	total_page := int(math.Ceil(count_total_page))
	limitStart := (new_page - 1) * new_per_page

	result := []Product{}

	query := database.Debug().Preload("Category")
	if search != "" {
		query.Where(
			"(name LIKE ? OR CAST(price as text) LIKE ?) OR code LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if soft_deleted != "" {
		query.Unscoped()
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

func GetCodeProduct(c *fiber.Ctx) error {
	now := time.Now()
	code := strconv.Itoa(now.Year()) + strconv.Itoa(int(now.Month())) + strconv.Itoa(now.Day()) + "-" + strconv.Itoa(now.Nanosecond())

	return c.Status(200).JSON(fiber.Map{
		"code": code,
	})
}

func StoreProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	product := new(requests.ProductRequest)

	if errParser := c.BodyParser(product); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	product.ValidateData()

	errorPhotoValidation := product.ValidatePhoto()
	if errorPhotoValidation != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errorPhotoValidation.Error(),
		})
	}

	errorCategoryValidation := product.ValidateExistsCategory(database)
	if errorCategoryValidation != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errorCategoryValidation.Error(),
		})
	}

	errValidate := requests.ValidateStruct(product)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	price, _ := strconv.ParseFloat(product.Price, 32)
	stock, _ := strconv.ParseInt(product.Stock, 10, 32)
	categoryID, _ := strconv.ParseUint(product.CategoryID, 10, 32)
	realCategory := uint(categoryID)
	realCategoryData := &realCategory

	productModel := models.Product{
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Stock:       int(stock),
		Price:       float32(price),
		CategoryID:  realCategoryData,
	}

	errCreateProduct := database.Create(&productModel).Error

	if errCreateProduct != nil {
		fmt.Println(errCreateProduct.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func ShowProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resultData := Product{}
	queryData := database.Preload("Category")
	queryData.Where("id = ?", id)
	queryData.First(&resultData)

	if resultData.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": resultData,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	product := new(requests.ProductRequest)

	if errParser := c.BodyParser(product); errParser != nil {
		fmt.Println(errParser.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	product.ValidateData()

	errorPhotoValidation := product.ValidatePhoto()
	if errorPhotoValidation != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errorPhotoValidation.Error(),
		})
	}

	errorCategoryValidation := product.ValidateExistsCategory(database)
	if errorCategoryValidation != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errorCategoryValidation.Error(),
		})
	}

	errValidate := requests.ValidateStruct(product)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	queryData := database.Model(&models.Product{})
	queryData.Select("name", "code", "description", "stock", "price", "photo", "category_id")
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	price, _ := strconv.ParseFloat(product.Price, 32)
	stock, _ := strconv.ParseInt(product.Stock, 10, 32)
	categoryID, _ := strconv.ParseUint(product.CategoryID, 10, 32)
	realCategory := uint(categoryID)
	realCategoryData := &realCategory

	queryData.Updates(&models.Product{
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Stock:       int(stock),
		Price:       float32(price),
		CategoryID:  realCategoryData,
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

func DestroyProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	queryData := database.Model(&models.Product{})
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Delete(&models.Product{})

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
