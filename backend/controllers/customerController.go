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

func IndexCustomer(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	page := c.Query("page", "1")
	new_page, _ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page, _ := strconv.Atoi(per_page)

	soft_deleted := c.Query("soft_deleted", "")

	search := c.Query("search", "")

	var resultCount int64

	var modelCount models.Customer

	queryResultCount := database.Model(&modelCount)
	queryResultCount.Select("id")

	if search != "" {
		queryResultCount.Or("name LIKE ?", "%"+search+"%")
		queryResultCount.Or("phone LIKE ?", "%"+search+"%")
		queryResultCount.Or("address LIKE ?", "%"+search+"%")
	}

	if soft_deleted != "" {
		queryResultCount.Unscoped()
	}

	queryResultCount.Count(&resultCount)

	count_total_page := float64(resultCount) / float64(new_per_page)
	total_page := int(math.Ceil(count_total_page))
	limitStart := (new_page - 1) * new_per_page

	result := []map[string]interface{}{}

	var queryCount models.Customer
	query := database.Model(&queryCount)
	query.Select("name", "id", "phone", "email", "address")
	if search != "" {
		query.Or("name LIKE ?", "%"+search+"%")
		query.Or("phone LIKE ?", "%"+search+"%")
		query.Or("address LIKE ?", "%"+search+"%")
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
		"page":       new_page,
	})
}

func GetAllCustomer(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	result := []map[string]interface{}{}

	database.Model(&models.Customer{}).Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func StoreCustomer(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	customer := new(requests.CustomerRequest)

	if errParser := c.BodyParser(customer); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	customer.ValidateData()

	errValidate := requests.ValidateStruct(customer)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	customerModel := models.Customer{
		Name:    customer.Name,
		Email:   customer.Email,
		Address: customer.Address,
		Phone:   customer.Phone,
	}

	errCreateCustomer := database.Create(&customerModel).Error

	if errCreateCustomer != nil {
		fmt.Println(errCreateCustomer.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func ShowCustomer(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resultData := map[string]interface{}{}

	queryData := database.Model(&models.Customer{})
	queryData.Select("name", "id", "phone", "email", "address")
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

func UpdateCustomer(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	customer := new(requests.CustomerRequest)

	if errParser := c.BodyParser(customer); errParser != nil {
		fmt.Println(errParser.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	customer.ValidateData()

	errValidate := requests.ValidateStruct(customer)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	queryData := database.Model(&models.Customer{})
	queryData.Select("name", "phone", "email", "address")
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Updates(&models.Customer{
		Name:    customer.Name,
		Email:   customer.Email,
		Address: customer.Address,
		Phone:   customer.Phone,
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

func DestroyCustomer(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	queryData := database.Model(&models.Customer{})
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Delete(&models.Customer{})

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
