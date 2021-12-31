package controllers

import (
	"fmt"
	"math"
	"strconv"

	"api-gofiber/pos/helpers"
	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func IndexUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	page := c.Query("page", "1")
	new_page, _ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page, _ := strconv.Atoi(per_page)

	soft_deleted := c.Query("soft_deleted", "")

	search := c.Query("search", "")

	var resultCount int64

	var modelCount models.User

	queryResultCount := database.Model(&modelCount)
	queryResultCount.Select("id")

	if search != "" {
		queryResultCount.Where("name LIKE ?", "%"+search+"%")
	}

	if soft_deleted != "" {
		queryResultCount.Unscoped()
	}

	queryResultCount.Count(&resultCount)

	count_total_page := float64(resultCount) / float64(new_per_page)
	total_page := int(math.Ceil(count_total_page))
	limitStart := (new_page - 1) * new_per_page

	result := []map[string]interface{}{}

	var queryCount models.User
	query := database.Model(&queryCount)
	query.Select("name", "id", "email")
	if search != "" {
		query.Where("name LIKE ?", "%"+search+"%")
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

func StoreUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	user := new(requests.UserAddRequest)

	if errParser := c.BodyParser(user); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	user.ValidateData()

	errValidate := requests.ValidateStruct(user)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	hash, errHash := helpers.HashPassword(user.Password)

	if errHash != nil {
		fmt.Println(errHash.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	userModel := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hash,
	}

	errCreateUser := database.Create(&userModel).Error

	if errCreateUser != nil {
		fmt.Println(errCreateUser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func ShowUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resultData := map[string]interface{}{}

	queryData := database.Model(&models.User{})
	queryData.Select("name", "id", "email")
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

func UpdateUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	user := new(requests.UserUpdateRequest)

	if errParser := c.BodyParser(user); errParser != nil {
		fmt.Println(errParser.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	user.ValidateData()

	errValidate := requests.ValidateStruct(user)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	queryData := database.Model(&models.User{})
	queryData.Select("name", "email")
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	if user.Password == nil {
		queryData.Updates(&models.User{
			Name:  user.Name,
			Email: user.Email,
		})
	} else {
		hash, errHash := helpers.HashPassword(*user.Password)

		if errHash != nil {
			fmt.Println(errHash.Error())

			return c.Status(500).JSON(fiber.Map{
				"message": "Terjadi Kesalahan",
			})
		}

		queryData.Updates(&models.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: hash,
		})
	}

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

func DestroyUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	queryData := database.Model(&models.User{})
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Delete(&models.User{})

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
